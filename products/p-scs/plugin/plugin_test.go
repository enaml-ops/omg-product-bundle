package pscs_test

import (
	"io/ioutil"

	"github.com/enaml-ops/enaml"
	pscs "github.com/enaml-ops/omg-product-bundle/products/p-scs/plugin"
	logging "github.com/op/go-logging"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("p-scs plugin", func() {

	BeforeSuite(func() {
		// suppress logging for tests
		logging.SetBackend(logging.NewLogBackend(ioutil.Discard, "", 0))
	})

	Context("when generating a manifest with incomplete input", func() {
		var (
			p   *pscs.Plugin
			err error
		)

		BeforeEach(func() {
			p = new(pscs.Plugin)
			_, err = p.GetProduct([]string{"foo"}, []byte{})
		})
		It("should yield an error", func() {
			Ω(err).Should(HaveOccurred())
		})
	})

	Context("when called with valid flags", func() {
		var (
			p  *pscs.Plugin
			dm *enaml.DeploymentManifest
		)

		BeforeEach(func() {
			p = new(pscs.Plugin)
			manifestBytes, err := p.GetProduct([]string{"foo",
				"--vm-type", "asdf",
				"--az", "asdf",
				"--system-domain", "asdf",
				"--app-domain", "asdf",
				"--network", "asdf",
				"--admin-password", "asdf",
				"--uaa-admin-secret", "asdf",
			}, []byte{})
			Ω(err).ShouldNot(HaveOccurred())
			dm = enaml.NewDeploymentManifest(manifestBytes)
		})
		It("should have the correct releases", func() {
			hasRelease := func(name, version string) bool {
				for i := range dm.Releases {
					if dm.Releases[i].Name == name && dm.Releases[i].Version == version {
						return true
					}
				}
				return false
			}
			Ω(hasRelease(pscs.SpringCloudBrokerReleaseName, pscs.SpringCloudBrokerReleaseVersion)).Should(BeTrue())
		})

		It("should have the correct instance groups", func() {
			Ω(dm.GetInstanceGroupByName("register-service-broker")).ShouldNot(BeNil(), "we should have a register service broker group")
			Ω(dm.GetInstanceGroupByName("deploy-service-broker")).ShouldNot(BeNil(), "we should have a deploy service broker group")
			Ω(dm.GetInstanceGroupByName("destroy-service-broker")).ShouldNot(BeNil(), "we should have a destroy service broker group")
		})

		It("should set the update", func() {
			Ω(dm.Update.Canaries).Should(Equal(1), "we found at least 1 canary in list")
			Ω(dm.Update.CanaryWatchTime).Should(Equal("30000-300000"))
			Ω(dm.Update.UpdateWatchTime).Should(Equal("30000-300000"))
			Ω(dm.Update.MaxInFlight).Should(Equal(1))
			Ω(dm.Update.Serial).Should(BeTrue())
		})

		It("should configure the stemcell", func() {
			Ω(dm.Stemcells).Should(HaveLen(1), "we found at least 1 stemcell in list")
			Ω(dm.Stemcells[0].OS).Should(Equal(pscs.StemcellName))
			Ω(dm.Stemcells[0].Alias).Should(Equal(pscs.StemcellAlias))
			Ω(dm.Stemcells[0].Version).Should(Equal(pscs.StemcellVersion))
		})
	})

	Context("when inferring defaults from cloud config", func() {
		var (
			p  *pscs.Plugin
			dm *enaml.DeploymentManifest
		)

		BeforeEach(func() {
			p = new(pscs.Plugin)
			cc, err := ioutil.ReadFile("fixtures/cloudconfig.yml")
			Ω(err).ShouldNot(HaveOccurred())

			manifestBytes, err := p.GetProduct([]string{"foo", "--infer-from-cloud",
				"--vm-type", "xlarge",
				"--system-domain", "asdf",
				"--app-domain", "asdf",
				"--admin-password", "asdf",
				"--uaa-admin-secret", "asdf",
			}, cc)
			Ω(err).ShouldNot(HaveOccurred())
			dm = enaml.NewDeploymentManifest(manifestBytes)
		})

		It("should not overwrite flags that were provided on the command line", func() {
			for _, name := range []string{"register-service-broker", "deploy-service-broker", "destroy-service-broker"} {
				ig := dm.GetInstanceGroupByName(name)
				Ω(ig).ShouldNot(BeNil(), "couldn't find instance group "+name)
				Ω(ig.VMType).Should(Equal("xlarge"))
			}
		})

		It("should infer values from the cloud config", func() {
			for _, name := range []string{"register-service-broker", "deploy-service-broker", "destroy-service-broker"} {
				ig := dm.GetInstanceGroupByName(name)
				Ω(ig).ShouldNot(BeNil(), "couldn't find instance group "+name)

				By("inferring AZs")
				Ω(ig.AZs).Should(ConsistOf("z1", "z2"))

				By("inferring the network name")
				Ω(ig.Networks).Should(HaveLen(1))
				Ω(ig.Networks[0].Name).Should(Equal("privatenet"))
			}
		})
	})
})
