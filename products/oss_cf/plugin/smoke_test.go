package cloudfoundry_test

import (
	"github.com/enaml-ops/omg-product-bundle/products/oss_cf/enaml-gen/smoke-tests"
	. "github.com/enaml-ops/omg-product-bundle/products/oss_cf/plugin"
	"github.com/enaml-ops/omg-product-bundle/products/oss_cf/plugin/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Smoke test errand", func() {
	Context("when initialized WITH a complete set of arguments", func() {
		var smokeErrand InstanceGroupCreator
		BeforeEach(func() {
			config := &config.Config{
				StemcellName:      "cool-ubuntu-animal",
				AZs:               []string{"eastprod-1"},
				NetworkName:       "foundry-net",
				SystemDomain:      "sys.test.com",
				AppDomains:        []string{"apps.test.com"},
				UAALoginProtocol:  "https",
				SkipSSLCertVerify: true,
				Secret:            config.Secret{},
				User:              config.User{},
				Certs:             &config.Certs{},
				InstanceCount:     config.InstanceCount{},
				IP:                config.IP{},
			}
			config.ErrandVMType = "blah"
			config.SmokeTestsPassword = "password"
			smokeErrand = NewSmokeErrand(config)
		})
		It("then it should have 1 instances", func() {
			ig := smokeErrand.ToInstanceGroup()
			Ω(ig.Instances).Should(Equal(1))
		})
		It("then it should allow the user to configure the AZs", func() {
			ig := smokeErrand.ToInstanceGroup()
			Ω(len(ig.AZs)).Should(Equal(1))
			Ω(ig.AZs[0]).Should(Equal("eastprod-1"))
		})

		It("then it should allow the user to configure vm-type", func() {
			ig := smokeErrand.ToInstanceGroup()
			Ω(ig.VMType).ShouldNot(BeEmpty())
			Ω(ig.VMType).Should(Equal("blah"))
		})

		It("then it should allow the user to configure network to use", func() {
			ig := smokeErrand.ToInstanceGroup()
			network := ig.GetNetworkByName("foundry-net")
			Ω(network).ShouldNot(BeNil())
		})

		It("then it should allow the user to configure the used stemcell", func() {
			ig := smokeErrand.ToInstanceGroup()
			Ω(ig.Stemcell).ShouldNot(BeEmpty())
			Ω(ig.Stemcell).Should(Equal("cool-ubuntu-animal"))
		})
		It("then it should have update max in-flight 1", func() {
			ig := smokeErrand.ToInstanceGroup()
			Ω(ig.Update.MaxInFlight).Should(Equal(1))
			Ω(ig.Update.Serial).Should(Equal(false))
		})
		It("then it should have lifecycle of errand", func() {
			ig := smokeErrand.ToInstanceGroup()
			Ω(ig.Lifecycle).Should(Equal("errand"))
		})

		It("then it should then have 1 jobs", func() {
			ig := smokeErrand.ToInstanceGroup()
			Ω(len(ig.Jobs)).Should(Equal(1))
		})
		It("then it should then have smoke-tests job", func() {
			ig := smokeErrand.ToInstanceGroup()
			job := ig.GetJobByName("smoke-tests")
			Ω(job).ShouldNot(BeNil())
			props, _ := job.Properties.(*smoke_tests.SmokeTestsJob)
			Ω(props.SmokeTests.AppsDomain).Should(Equal("apps.test.com"))
			Ω(props.SmokeTests.Api).Should(Equal("https://api.sys.test.com"))
			Ω(props.SmokeTests.Org).Should(Equal("CF_SMOKE_TEST_ORG"))
			Ω(props.SmokeTests.Space).Should(Equal("CF_SMOKE_TEST_SPACE"))
			Ω(props.SmokeTests.User).Should(Equal("smoke_tests"))
			Ω(props.SmokeTests.Password).Should(Equal("password"))
			Ω(props.SmokeTests.UseExistingOrg).Should(BeFalse())
			Ω(props.SmokeTests.UseExistingSpace).Should(BeFalse())
			Ω(props.SmokeTests.SkipSslValidation).Should(BeTrue())
		})
	})
})
