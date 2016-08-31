package cloudfoundry_test

import (
	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/acceptance-tests"
	. "github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/plugin"
	"github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/plugin/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("given the acceptance-tests partition", func() {

	Context("when initialized with valid config", func() {
		It("generates different instance group names for internet-less tests", func() {

			config := &config.Config{
				SystemDomain:  "sys.yourdomain.com",
				AppDomains:    []string{"apps.yourdomain.com"},
				AZs:           []string{"z1"},
				StemcellName:  "cool-ubuntu-animal",
				NetworkName:   "foundry-net",
				AdminPassword: "adminpass",
			}
			withInternet := NewAcceptanceTestsPartition(true, config).ToInstanceGroup()
			withoutInternet := NewAcceptanceTestsPartition(false, config).ToInstanceGroup()
			Ω(withInternet.Name).ShouldNot(Equal(withoutInternet.Name))
			Ω(withInternet.Jobs[0].Name).Should(Equal(withoutInternet.Jobs[0].Name))
		})
	})

	Context("when initialized with a complete set of arguments", func() {
		const includeInternetDependent = true
		var ig InstanceGroupCreator
		var dm *enaml.DeploymentManifest
		BeforeEach(func() {
			config := &config.Config{
				SystemDomain:      "sys.yourdomain.com",
				AppDomains:        []string{"apps.yourdomain.com"},
				AZs:               []string{"z1"},
				StemcellName:      "cool-ubuntu-animal",
				NetworkName:       "foundry-net",
				SkipSSLCertVerify: true,
				AdminPassword:     "adminpass",
				ErrandVMType:      "foo",
			}
			ig = NewAcceptanceTestsPartition(includeInternetDependent, config)
			dm = new(enaml.DeploymentManifest)
			dm.AddInstanceGroup(ig.ToInstanceGroup())
		})

		It("should have the correct VM type and lifecycle", func() {
			group := dm.GetInstanceGroupByName("acceptance-tests")
			Ω(group.Lifecycle).Should(Equal("errand"))
			Ω(group.VMType).Should(Equal("foo"))
		})

		It("should have a single instance", func() {
			group := dm.GetInstanceGroupByName("acceptance-tests")
			Ω(group.Instances).Should(Equal(1))
		})

		It("should have update max in flight 1", func() {
			group := dm.GetInstanceGroupByName("acceptance-tests")
			Ω(group.Update.MaxInFlight).Should(Equal(1))
		})

		It("should allow the user to configure the AZs", func() {
			group := dm.GetInstanceGroupByName("acceptance-tests")
			Ω(len(group.AZs)).Should(Equal(1))
			Ω(group.AZs[0]).Should(Equal("z1"))
		})

		It("should allow the user to configure the used stemcell", func() {
			group := dm.GetInstanceGroupByName("acceptance-tests")
			Ω(group.Stemcell).Should(Equal("cool-ubuntu-animal"))
		})

		It("should allow the user to configure the network to use", func() {
			group := dm.GetInstanceGroupByName("acceptance-tests")
			Ω(len(group.Networks)).Should(Equal(1))
			Ω(group.Networks[0].Name).Should(Equal("foundry-net"))
		})

		It("should have correctly configured the acceptance-tests job", func() {
			group := ig.ToInstanceGroup()
			job := group.GetJobByName("acceptance-tests")
			Ω(job.Release).Should(Equal(CFReleaseName))

			props := job.Properties.(*acceptance_tests.AcceptanceTestsJob)
			Ω(props.AcceptanceTests.Api).Should(Equal("https://api.sys.yourdomain.com"))
			Ω(props.AcceptanceTests.AppsDomain).Should(Equal("apps.yourdomain.com"))
			Ω(props.AcceptanceTests.AdminUser).Should(Equal("admin"))
			Ω(props.AcceptanceTests.AdminPassword).Should(Equal("adminpass"))
			Ω(props.AcceptanceTests.IncludeLogging).Should(BeTrue())
			Ω(props.AcceptanceTests.IncludeOperator).Should(BeTrue())
			Ω(props.AcceptanceTests.IncludeServices).Should(BeTrue())
			Ω(props.AcceptanceTests.IncludeSecurityGroups).Should(BeTrue())
			Ω(props.AcceptanceTests.SkipSslValidation).Should(BeTrue())
			Ω(props.AcceptanceTests.SkipRegex).Should(Equal("lucid64"))
			Ω(props.AcceptanceTests.JavaBuildpackName).Should(Equal("java_buildpack_offline"))

			Ω(props.AcceptanceTests.IncludeInternetDependent).Should(BeTrue())
		})
	})

	Context("when initialized with a complete set of arguments in internetless mode", func() {
		const includeInternetDependent = false
		var ig InstanceGroupCreator
		var dm *enaml.DeploymentManifest
		BeforeEach(func() {
			config := &config.Config{
				SystemDomain:      "sys.yourdomain.com",
				AppDomains:        []string{"apps.yourdomain.com"},
				AZs:               []string{"z1"},
				StemcellName:      "cool-ubuntu-animal",
				NetworkName:       "foundry-net",
				SkipSSLCertVerify: true,
				AdminPassword:     "adminpass",
			}
			ig = NewAcceptanceTestsPartition(includeInternetDependent, config)
			dm = new(enaml.DeploymentManifest)
			dm.AddInstanceGroup(ig.ToInstanceGroup())
		})

		It("should not be configured to include internet-dependent tests", func() {
			group := ig.ToInstanceGroup()
			job := group.GetJobByName("acceptance-tests")
			Ω(job.Release).Should(Equal(CFReleaseName))
			props := job.Properties.(*acceptance_tests.AcceptanceTestsJob)
			Ω(props.AcceptanceTests.Api).Should(Equal("https://api.sys.yourdomain.com"))
			Ω(props.AcceptanceTests.AppsDomain).Should(Equal("apps.yourdomain.com"))
			Ω(props.AcceptanceTests.AdminUser).Should(Equal("admin"))
			Ω(props.AcceptanceTests.AdminPassword).Should(Equal("adminpass"))
			Ω(props.AcceptanceTests.IncludeLogging).Should(BeTrue())
			Ω(props.AcceptanceTests.IncludeOperator).Should(BeTrue())
			Ω(props.AcceptanceTests.IncludeServices).Should(BeTrue())
			Ω(props.AcceptanceTests.IncludeSecurityGroups).Should(BeTrue())
			Ω(props.AcceptanceTests.SkipSslValidation).Should(BeTrue())
			Ω(props.AcceptanceTests.SkipRegex).Should(Equal("lucid64"))
			Ω(props.AcceptanceTests.JavaBuildpackName).Should(Equal("java_buildpack_offline"))

			Ω(props.AcceptanceTests.IncludeInternetDependent).Should(BeFalse())
		})
	})
})
