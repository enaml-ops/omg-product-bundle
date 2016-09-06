package pmysql_test

import (
	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-product-bundle/products/p-mysql/enaml-gen/acceptance-tests"
	. "github.com/enaml-ops/omg-product-bundle/products/p-mysql/plugin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("given NewAcceptanceTests func", func() {
	Describe("given a valid plugin object argument", func() {
		var plgn *Plugin
		BeforeEach(func() {
			plgn = new(Plugin)
		})

		Context("when the plugin sets monitoring ip values", func() {
			var ig *enaml.InstanceGroup
			BeforeEach(func() {
				ig = NewBrokerRegistrar(plgn)
			})

			It("then it should create vm instances for each given ip", func() {
				Ω(ig.Instances).Should(Equal(1), "does not create the correct number of instances")
			})

			It("then it should create a static ip block for the given ips", func() {
				Ω(len(ig.Networks)).Should(Equal(1))
			})
		})

		Context("when creating the acceptance-tests instancegroup", func() {
			var ig *enaml.InstanceGroup
			var jobProperties *acceptance_tests.AcceptanceTestsJob

			const (
				controlCFAdminPassword = "cfadmin"
				controlProxyUser       = "proxy"
				controlProxyPass       = "proxypass"
			)

			BeforeEach(func() {
				plgn.BaseDomain = "abc.example.com"
				plgn.CFAdminPassword = controlCFAdminPassword
				plgn.ProxyAPIUser = controlProxyUser
				plgn.ProxyAPIPass = controlProxyPass

				ig = NewAcceptanceTests(plgn)
				jobProperties = ig.GetJobByName("acceptance-tests").Properties.(*acceptance_tests.AcceptanceTestsJob)
			})

			It("then it should contain instance jobs for all required releases", func() {
				Ω(len(ig.Jobs)).Should(Equal(1))
				Ω(checkJobExists(ig.Jobs, "acceptance-tests")).Should(BeTrue(), "there should be a job for acceptance-tests release")
			})

			It("then it should setup acceptance tests job information", func() {
				Ω(jobProperties.TimeoutScale).Should(Equal(1))

				Ω(jobProperties.Cf).ShouldNot(BeNil(), "should have a properly intitialized cf")
				Ω(jobProperties.Cf.ApiUrl).Should(Equal("https://api.sys.abc.example.com"))
				Ω(jobProperties.Cf.AdminUsername).Should(Equal("admin"))
				Ω(jobProperties.Cf.AdminPassword).Should(Equal(controlCFAdminPassword))
				Ω(jobProperties.Cf.AppsDomain).Should(Equal("apps.abc.example.com"))
				Ω(jobProperties.Cf.SkipSslValidation).Should(BeTrue())

				Ω(jobProperties.Proxy).ShouldNot(BeNil(), "should have a properly intitialized proxy")
				Ω(jobProperties.Proxy.ExternalHost).Should(Equal("p-mysql.sys.abc.example.com"))
				Ω(jobProperties.Proxy.ApiUsername).Should(Equal(controlProxyUser))
				Ω(jobProperties.Proxy.ApiPassword).Should(Equal(controlProxyPass))

				Ω(jobProperties.Broker).ShouldNot(BeNil(), "should have a properly intitialized broker")
				Ω(jobProperties.Broker.Host).Should(Equal("p-mysql.sys.abc.example.com"))

				Ω(jobProperties.Service).ShouldNot(BeNil(), "should have a properly intitialized service")
				Ω(jobProperties.Service.Name).Should(Equal("p-mysql"))

				plans := jobProperties.Service.Plans.([]map[string]interface{})
				Ω(plans).Should(HaveLen(1))
				Ω(plans[0]).Should(HaveKeyWithValue("name", "100mb-dev"))
				Ω(plans[0]).Should(HaveKeyWithValue("max_storage_mb", 100))
			})
		})
	})
})
