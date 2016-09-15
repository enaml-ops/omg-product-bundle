package cloudfoundry_test

import (
	"github.com/enaml-ops/enaml"
	pam "github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/push-apps-manager"
	. "github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/plugin"
	"github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/plugin/config"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("push-apps-manager", func() {
	Context("when initialized with a complete set of arguments", func() {
		var igc InstanceGroupCreator
		var ig *enaml.InstanceGroup
		var dm *enaml.DeploymentManifest

		const (
			controlPushAppsPassword   = "pushappspassword"
			controlPortalClientSecret = "portalclientsecret"
			controlSecretToken        = "secrettoken"
			controlProxyIP            = "10.1.2.3"
			controlConsoleDBUser      = "consoledbuser"
			controlConsoleDBPass      = "consoledbpass"
			controlAppUsageDBPass     = "appusagedbpassword"
			controlMySQLAdminPassword = "mysqladmin"
		)

		BeforeEach(func() {
			c := &config.Config{
				AZs:          []string{"z1"},
				StemcellName: "cool-ubuntu-animal",
				NetworkName:  "foundry-net",
				SystemDomain: "sys.example.com",
				Secret: config.Secret{
					PushAppsManagerPassword: controlPushAppsPassword,
					PortalClientSecret:      controlPortalClientSecret,
					AppsManagerSecretToken:  controlSecretToken,
					ConsoleDBPassword:       controlConsoleDBPass,
					AppUsageDBPassword:      controlAppUsageDBPass,
					MySQLAdminPassword:      controlMySQLAdminPassword,
				},
				User: config.User{
					ConsoleDBUserName: controlConsoleDBUser,
				},
				Certs:         &config.Certs{},
				InstanceCount: config.InstanceCount{},
				IP: config.IP{
					MySQLProxyIPs: []string{controlProxyIP},
				},
				SkipSSLCertVerify: true,
			}
			igc = NewPushAppsManager(c)
			ig = igc.ToInstanceGroup()

			dm = new(enaml.DeploymentManifest)
			dm.AddInstanceGroup(ig)
			Ω(dm.GetInstanceGroupByName("push-apps-manager")).ShouldNot(BeNil())
		})

		It("should have a single instance", func() {
			Ω(ig.Instances).Should(Equal(1))
		})

		It("should have update max in flight 1", func() {
			Ω(ig.Update.MaxInFlight).Should(Equal(1))
		})

		It("should have lifecycle errand", func() {
			Ω(ig.Lifecycle).Should(Equal("errand"))
		})

		It("should have a single AZ", func() {
			Ω(ig.AZs).Should(ConsistOf("z1"))
		})

		It("should have set the stemcell name", func() {
			Ω(ig.Stemcell).Should(Equal("cool-ubuntu-animal"))
		})

		It("should have set the network name", func() {
			Ω(ig.Networks).Should(HaveLen(1))
			Ω(ig.Networks[0].Name).Should(Equal("foundry-net"))
		})

		It("should have configured the push-apps-manager job", func() {
			Ω(ig.Jobs).Should(HaveLen(1))
			job := ig.GetJobByName("push-apps-manager")
			Ω(job).ShouldNot(BeNil())

			By("setting the release name")
			Ω(job.Release).Should(Equal(PushAppsReleaseName))

			By("setting the job properties")
			Ω(job.Properties).ShouldNot(BeNil())
			props, ok := job.Properties.(*pam.PushAppsManagerJob)
			Ω(ok).Should(BeTrue())

			By("setting the CF properties")
			Ω(props.Cf.ApiUrl).Should(Equal("https://api.sys.example.com"))
			Ω(props.Cf.AdminUsername).Should(Equal("push_apps_manager"))
			Ω(props.Cf.AdminPassword).Should(Equal(controlPushAppsPassword))
			Ω(props.Cf.SystemDomain).Should(Equal("sys.example.com"))

			By("configuring the service authentication")
			Ω(props.Services).ShouldNot(BeNil())
			Ω(props.Services.Authentication).ShouldNot(BeNil())
			Ω(props.Services.Authentication.CFCLIENTID).Should(Equal("portal"))
			Ω(props.Services.Authentication.CFCLIENTSECRET).Should(Equal(controlPortalClientSecret))
			Ω(props.Services.Authentication.CFUAASERVERURL).Should(Equal("https://uaa.sys.example.com"))
			Ω(props.Services.Authentication.CFLOGINSERVERURL).Should(Equal("https://login.sys.example.com"))

			By("configuring the environment")
			Ω(props.Env).ShouldNot(BeNil())
			Ω(props.Env.SecretToken).Should(Equal(controlSecretToken))
			Ω(props.Env.CfCcApiUrl).Should(Equal("https://api.sys.example.com"))
			Ω(props.Env.CfLoggregatorHttpUrl).Should(Equal("http://loggregator.sys.example.com"))
			Ω(props.Env.CfConsoleUrl).Should(Equal("https://apps.sys.example.com"))
			Ω(props.Env.CfNotificationsServiceUrl).Should(Equal("https://notifications.sys.example.com"))
			Ω(props.Env.UsageServiceHost).Should(Equal("https://app-usage.sys.example.com"))
			Ω(props.Env.BundleWithout).Should(Equal("test development hosted_only"))
			Ω(props.Env.EnableInternalUserStore).Should(BeFalse())
			Ω(props.Env.EnableNonAdminRoleManagement).Should(BeFalse())
			Ω(props.Env.GenericWhiteLabelConfigJson).ShouldNot(BeNil())
			Ω(props.Env.GenericWhiteLabelConfigJson.FooterText).ShouldNot(BeEmpty())

			By("configuring the databases")
			Ω(props.Databases).ShouldNot(BeNil())
			Ω(props.Databases.Console).ShouldNot(BeNil())
			Ω(props.Databases.Console.Ip).Should(Equal(controlProxyIP))
			Ω(props.Databases.Console.Username).Should(Equal("root"))
			Ω(props.Databases.Console.Password).Should(Equal(controlMySQLAdminPassword))
			Ω(props.Databases.Console.Adapter).Should(Equal("mysql"))
			Ω(props.Databases.Console.Port).Should(Equal(3306))
			Ω(props.Databases.AppUsageService).ShouldNot(BeNil())
			Ω(props.Databases.AppUsageService.Name).Should(Equal("app_usage_service"))
			Ω(props.Databases.AppUsageService.Ip).Should(Equal(controlProxyIP))
			Ω(props.Databases.AppUsageService.Port).Should(Equal(3306))
			Ω(props.Databases.AppUsageService.Username).Should(Equal("root"))
			Ω(props.Databases.AppUsageService.Password).Should(Equal(controlMySQLAdminPassword))

			By("configuring SSL")
			Ω(props.Ssl).ShouldNot(BeNil())
			Ω(props.Ssl.SkipCertVerify).Should(BeTrue())
			Ω(props.Ssl.HttpsOnlyMode).Should(BeTrue())
		})
	})
})
