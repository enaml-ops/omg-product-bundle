package dockerregistry_test

import (
	"io/ioutil"

	"github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/debian_nfs_server"
	. "github.com/enaml-ops/omg-product-bundle/products/dockerregistry"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v2"
)

var _ = Describe("docker-registry Deployment", func() {
	controlAzs := []string{"z1"}
	controlNetworkName := "private"
	controlRelease := "docker-registry"
	controlStemcellAlias := "trusty"
	controlNFSServerType := "large"
	controlNFSDiskType := "large"
	controlNFSIP := "10.0.0.10"

	controlRegistryServerType := "medium"
	controlRegistryIPs := []string{"10.0.0.8", "10.0.0.9"}
	var deployment DockerRegistry
	BeforeEach(func() {
		deployment = DockerRegistry{
			AZs:            controlAzs,
			NetworkName:    controlNetworkName,
			NFSServerType:  controlNFSServerType,
			NFSDiskType:    controlNFSDiskType,
			NFSIP:          controlNFSIP,
			StemcellAlias:  controlStemcellAlias,
			RegistryIPs:    controlRegistryIPs,
			RegistryVMType: controlRegistryServerType,
		}
	})
	Describe("Given a CreateNFSServerInstanceGroup", func() {
		Context("when calling", func() {
			It("then we should return a valid enaml.InstanceGroup", func() {
				ig := deployment.CreateNFSServerInstanceGroup()
				Ω(ig.Instances).Should(Equal(1))
				Ω(ig.AZs).Should(ConsistOf(controlAzs))
				Ω(ig.VMType).Should(Equal(controlNFSServerType))
				Ω(ig.PersistentDiskType).Should(Equal(controlNFSDiskType))
				Ω(len(ig.Networks)).Should(Equal(1))
				Ω(ig.Networks[0].StaticIPs).Should(ConsistOf([]string{controlNFSIP}))
				Ω(ig.Networks[0].Name).Should(Equal(controlNetworkName))
				Ω(ig.Stemcell).Should(Equal(controlStemcellAlias))
				Ω(len(ig.Jobs)).Should(Equal(1))
				job := ig.Jobs[0]
				Ω(job.Name).Should(Equal("nfs-server"))
				Ω(job.Release).Should(Equal(controlRelease))
				jobProperties := job.Properties.(*debian_nfs_server.DebianNfsServerJob)
				Ω(jobProperties.NfsServer).ShouldNot(BeNil())
				Ω(jobProperties.NfsServer.AllowFromEntries).Should(ConsistOf(controlRegistryIPs))
			})
			It("then should return valid YML", func() {
				ig := deployment.CreateNFSServerInstanceGroup()
				bytes, err := ioutil.ReadFile("fixtures/nfs-instancegroup.yml")
				Ω(err).ShouldNot(HaveOccurred())
				igYml, err := yaml.Marshal(ig)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(igYml).Should(MatchYAML(bytes))
			})
		})

		Describe("Given a CreateRegistryInstanceGroup", func() {
			Context("when calling", func() {
				It("then we should return a valid enaml.InstanceGroup", func() {
					ig := deployment.CreateRegistryInstanceGroup()
					Ω(ig).ShouldNot(BeNil())
				})
				It("then should return valid YML", func() {
					ig := deployment.CreateRegistryInstanceGroup()
					bytes, err := ioutil.ReadFile("fixtures/registry-instancegroup.yml")
					Ω(err).ShouldNot(HaveOccurred())
					igYml, err := yaml.Marshal(ig)
					Ω(err).ShouldNot(HaveOccurred())
					Ω(igYml).Should(MatchYAML(bytes))
				})
			})
		})

		Describe("Given a CreateProxyInstanceGroup", func() {
			Context("when calling", func() {
				It("then we should return a valid enaml.InstanceGroup", func() {
					ig := deployment.CreateProxyInstanceGroup()
					Ω(ig).ShouldNot(BeNil())
				})
			})
		})

		Describe("Given a CreateUpdate", func() {
			Context("when calling", func() {
				It("then we should return a valid enaml.Update", func() {
					bytes, err := ioutil.ReadFile("fixtures/update.yml")
					Ω(err).ShouldNot(HaveOccurred())
					updateYml, err := yaml.Marshal(deployment.CreateUpdate())
					Ω(err).ShouldNot(HaveOccurred())
					Ω(updateYml).Should(MatchYAML(bytes))
				})
			})
		})
	})
})
