package cloudfoundry_test

import (
	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-product-bundle/products/cf-mysql/enaml-gen/bootstrap"
	. "github.com/enaml-ops/omg-product-bundle/products/oss_cf/plugin"
	"github.com/enaml-ops/omg-product-bundle/products/oss_cf/plugin/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("given a bootstrap partition", func() {

	Context("when initialized with a complete set of arguments", func() {
		var ig InstanceGroupCreator
		var dm *enaml.DeploymentManifest
		BeforeEach(func() {
			config := &config.Config{
				AZs:           []string{"z1"},
				StemcellName:  "cool-ubuntu-animal",
				NetworkName:   "foundry-net",
				Secret:        config.Secret{},
				User:          config.User{},
				Certs:         &config.Certs{},
				InstanceCount: config.InstanceCount{},
				IP:            config.IP{},
			}
			config.MySQLIPs = []string{"10.0.0.26", "10.0.0.27", "10.0.0.28"}
			config.MySQLBootstrapUser = "user"
			config.MySQLBootstrapPassword = "pass"
			config.ErrandVMType = "foo"
			ig = NewBootstrapPartition(config)

			dm = new(enaml.DeploymentManifest)
			dm.AddInstanceGroup(ig.ToInstanceGroup())
		})

		It("should have the correct VM type and lifecycle", func() {
			group := dm.GetInstanceGroupByName("bootstrap")
			Ω(group.Lifecycle).Should(Equal("errand"))
			Ω(group.VMType).Should(Equal("foo"))
		})

		It("should have a single instance", func() {
			group := dm.GetInstanceGroupByName("bootstrap")
			Ω(group.Instances).Should(Equal(1))
		})

		It("should have update max in flight 1", func() {
			group := dm.GetInstanceGroupByName("bootstrap")
			Ω(group.Update.MaxInFlight).Should(Equal(1))
		})

		It("should allow the user to configure the AZs", func() {
			group := dm.GetInstanceGroupByName("bootstrap")
			Ω(len(group.AZs)).Should(Equal(1))
			Ω(group.AZs[0]).Should(Equal("z1"))
		})

		It("should allow the user to configure the used stemcell", func() {
			group := dm.GetInstanceGroupByName("bootstrap")
			Ω(group.Stemcell).Should(Equal("cool-ubuntu-animal"))
		})

		It("should allow the user to configure the network to use", func() {
			group := dm.GetInstanceGroupByName("bootstrap")
			Ω(len(group.Networks)).Should(Equal(1))
			Ω(group.Networks[0].Name).Should(Equal("foundry-net"))
		})

		It("should have a valid bootstrap job", func() {
			group := dm.GetInstanceGroupByName("bootstrap")
			job := group.GetJobByName("bootstrap")
			Ω(job.Release).Should(Equal(CFMysqlReleaseName))

			props := job.Properties.(*bootstrap.BootstrapJob)
			Ω(props.ClusterIps).Should(ConsistOf("10.0.0.26", "10.0.0.27", "10.0.0.28"))
			Ω(props.DatabaseStartupTimeout).Should(Equal(1200))
			Ω(props.BootstrapEndpoint.Username).Should(Equal("user"))
			Ω(props.BootstrapEndpoint.Password).Should(Equal("pass"))
		})

	})
})
