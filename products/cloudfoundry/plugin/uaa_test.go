package cloudfoundry_test

import (
	"github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/route_registrar"
	"github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/uaa"
	. "github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/plugin"
	"github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/plugin/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("UAA Partition", func() {
	Context("when initialized WITH a complete set of arguments", func() {
		var uaaPartition InstanceGroupCreator
		BeforeEach(func() {

			config := &config.Config{
				SystemDomain:            "sys.test.com",
				AZs:                     []string{"eastprod-1"},
				StemcellName:            "cool-ubuntu-animal",
				NetworkName:             "foundry-net",
				AllowSSHAccess:          true,
				NATSPort:                4222,
				DopplerZone:              "DopplerZoneguid",
				SyslogAddress:           "syslog-server",
				SyslogPort:              10601,
				SyslogTransport:         "tcp",
				LDAPEnabled:             true,
				LDAPUrl:                 "ldap://ldap.test.com",
				LDAPUserDN:              "userdn",
				LDAPSearchFilter:        "filter",
				LDAPSearchBase:          "base",
				LDAPMailAttributeName:   "mail",
				UAALoginProtocol:        "https",
				SelfServiceLinksEnabled: true,
				SignupsEnabled:          true,
				Secret:                  config.Secret{},
				User:                    config.User{},
				Certs:                   &config.Certs{},
				InstanceCount:           config.InstanceCount{},
				IP:                      config.IP{},
			}
			config.LDAPUserPassword = "userpwd"
			config.EtcdMachines = []string{"1.0.0.7", "1.0.0.8"}
			config.SAMLServiceProviderKey = "saml-key"
			config.SAMLServiceProviderCertificate = "saml-cert"
			config.JWTVerificationKey = "jwt-verificationkey"
			config.JWTSigningKey = "jwt-signingkey"
			config.NATSMachines = []string{"1.0.0.5", "1.0.0.6"}
			config.UAAVMType = "blah"
			config.UAAInstances = 1
			config.ConsulIPs = []string{"1.0.0.1", "1.0.0.2"}
			config.DopplerSharedSecret = "metronsecret"
			config.AdminSecret = "adminclientsecret"
			config.RouterMachines = []string{"1.0.0.1", "1.0.0.2"}
			config.MySQLProxyIPs = []string{"1.0.10.3", "1.0.10.4"}
			config.UAADBUserName = "uaa-db-user"
			config.UAADBPassword = "uaa-db-pwd"
			config.AdminPassword = "admin"
			config.PushAppsManagerPassword = "appsman"
			config.SmokeTestsPassword = "smoke"
			config.SystemServicesPassword = "sysservices"
			config.SystemVerificationPassword = "sysverification"
			config.OpentsdbFirehoseNozzleClientSecret = "opentsdb-firehose-nozzle-client-secret"
			config.IdentityClientSecret = "identity-client-secret"
			config.LoginClientSecret = "login-client-secret"
			config.PortalClientSecret = "portal-client-secret"
			config.AutoScalingServiceClientSecret = "autoscaling-service-client-secret"
			config.SystemPasswordsClientSecret = "system-passwords-client-secret"
			config.CCServiceDashboardsClientSecret = "cc-service-dashboards-client-secret"
			config.DopplerSecret = "doppler-client-secret"
			config.GoRouterClientSecret = "gorouter-client-secret"
			config.NotificationsClientSecret = "notifications-client-secret"
			config.NotificationsUIClientSecret = "notifications-ui-client-secret"
			config.CloudControllerUsernameLookupClientSecret = "cloud-controller-username-lookup-client-secret"
			config.CCRoutingClientSecret = "cc-routing-client-secret"
			config.SSHProxyClientSecret = "ssh-proxy-client-secret"
			config.AppsMetricsClientSecret = "apps-metrics-client-secret"
			config.AppsMetricsProcessingClientSecret = "apps-metrics-processing-client-secret"
			config.ConsulEncryptKeys = []string{"encyption-key"}
			config.ConsulCaCert = "ca-cert"
			config.ConsulAgentCert = "agent-cert"
			config.ConsulAgentKey = "agent-key"
			config.ConsulServerCert = "server-cert"
			config.ConsulServerKey = "server-key"
			config.NATSUser = "nats"
			config.NATSPassword = "pass"
			uaaPartition = NewUAAPartition(config)
		})
		It("then it should not configure static ips for uaaPartition", func() {
			ig := uaaPartition.ToInstanceGroup()
			network := ig.Networks[0]
			Ω(len(network.StaticIPs)).Should(Equal(0))
		})
		It("then it should have 1 instances", func() {
			ig := uaaPartition.ToInstanceGroup()
			Ω(ig.Instances).Should(Equal(1))
		})
		It("then it should allow the user to configure the AZs", func() {
			ig := uaaPartition.ToInstanceGroup()
			Ω(len(ig.AZs)).Should(Equal(1))
			Ω(ig.AZs[0]).Should(Equal("eastprod-1"))
		})

		It("then it should allow the user to configure vm-type", func() {
			ig := uaaPartition.ToInstanceGroup()
			Ω(ig.VMType).ShouldNot(BeEmpty())
			Ω(ig.VMType).Should(Equal("blah"))
		})

		It("then it should allow the user to configure network to use", func() {
			ig := uaaPartition.ToInstanceGroup()
			network := ig.GetNetworkByName("foundry-net")
			Ω(network).ShouldNot(BeNil())
		})

		It("then it should allow the user to configure the used stemcell", func() {
			ig := uaaPartition.ToInstanceGroup()
			Ω(ig.Stemcell).ShouldNot(BeEmpty())
			Ω(ig.Stemcell).Should(Equal("cool-ubuntu-animal"))
		})
		It("then it should have update max in-flight 1 and serial false", func() {
			ig := uaaPartition.ToInstanceGroup()
			Ω(ig.Update.MaxInFlight).Should(Equal(1))
			Ω(ig.Update.Serial).Should(Equal(false))
		})

		It("then it should then have 5 jobs", func() {
			ig := uaaPartition.ToInstanceGroup()
			Ω(len(ig.Jobs)).Should(Equal(5))
		})
		It("then it should then have uaa job with client secret", func() {
			ig := uaaPartition.ToInstanceGroup()
			job := ig.GetJobByName("uaa")
			Ω(job).ShouldNot(BeNil())
			props, _ := job.Properties.(*uaa.UaaJob)
			Ω(props.Uaa.Admin).ShouldNot(BeNil())
			Ω(props.Uaa.Admin.ClientSecret).Should(Equal("adminclientsecret"))
		})
		It("then it should then have uaa job with proxy configured", func() {
			ig := uaaPartition.ToInstanceGroup()
			job := ig.GetJobByName("uaa")
			Ω(job).ShouldNot(BeNil())
			props, _ := job.Properties.(*uaa.UaaJob)
			Ω(props.Uaa.Proxy).ShouldNot(BeNil())
			Ω(props.Uaa.Proxy.Servers).Should(ConsistOf("1.0.0.1", "1.0.0.2"))
		})
		It("then it should then have uaa job with UAADB", func() {
			ig := uaaPartition.ToInstanceGroup()
			job := ig.GetJobByName("uaa")
			Ω(job).ShouldNot(BeNil())
			props, _ := job.Properties.(*uaa.UaaJob)
			Ω(props.Uaadb).ShouldNot(BeNil())
			Ω(props.Uaadb.DbScheme).Should(Equal("mysql"))
			Ω(props.Uaadb.Port).Should(Equal(3306))
			Ω(props.Uaadb.Address).Should(Equal("1.0.10.3"))
			Ω(props.Uaadb.Roles).Should(HaveLen(1))
			role := props.Uaadb.Roles.([]map[string]interface{})[0]
			Ω(role).Should(HaveKeyWithValue("tag", "admin"))
			Ω(role).Should(HaveKeyWithValue("name", "uaa-db-user"))
			Ω(role).Should(HaveKeyWithValue("password", "uaa-db-pwd"))
		})

		It("then it should have uaa databases as an array", func() {
			ig := uaaPartition.ToInstanceGroup()
			job := ig.GetJobByName("uaa")
			Ω(job).ShouldNot(BeNil())
			props, _ := job.Properties.(*uaa.UaaJob)
			Ω(props.Uaadb).ShouldNot(BeNil())
			Ω(props.Uaadb.Databases).ShouldNot(BeNil())
			dbs := props.Uaadb.Databases.([]map[string]interface{})
			Ω(dbs).Should(HaveLen(1))
			Ω(dbs[0]["tag"]).Should(Equal("uaa"))
			Ω(dbs[0]["name"]).Should(Equal("uaa"))
		})

		It("then it should then have uaa job with Clients", func() {
			ig := uaaPartition.ToInstanceGroup()
			job := ig.GetJobByName("uaa")
			Ω(job).ShouldNot(BeNil())
			props, _ := job.Properties.(*uaa.UaaJob)
			Ω(props.Uaa.Clients).ShouldNot(BeNil())
			clientMap := props.Uaa.Clients.(map[string]UAAClient)
			Ω(len(clientMap)).Should(Equal(19))

		})
		It("then it should then have uaa job with SCIM", func() {
			ig := uaaPartition.ToInstanceGroup()
			job := ig.GetJobByName("uaa")
			Ω(job).ShouldNot(BeNil())
			props, _ := job.Properties.(*uaa.UaaJob)
			Ω(props.Uaa.Scim).ShouldNot(BeNil())
			Ω(props.Uaa.Scim.User).ShouldNot(BeNil())
			Ω(props.Uaa.Scim.User.Override).Should(BeTrue())
			Ω(props.Uaa.Scim.UseridsEnabled).Should(BeTrue())
			Ω(props.Uaa.Scim.Users).ShouldNot(BeNil())
			users := props.Uaa.Scim.Users.([]string)
			Ω(len(users)).Should(Equal(5))
		})
		It("then it should then have uaa job with valid login information", func() {
			ig := uaaPartition.ToInstanceGroup()
			job := ig.GetJobByName("uaa")
			Ω(job).ShouldNot(BeNil())
			props, _ := job.Properties.(*uaa.UaaJob)
			Ω(props.Domain).Should(Equal("sys.test.com"))
		})
		It("then it should then have uaa job with valid uaa information", func() {
			ig := uaaPartition.ToInstanceGroup()
			job := ig.GetJobByName("uaa")
			Ω(job).ShouldNot(BeNil())
			props, _ := job.Properties.(*uaa.UaaJob)
			Ω(props.Uaa).ShouldNot(BeNil())
			Ω(props.Uaa.CatalinaOpts).Should(Equal("-Xmx768m -XX:MaxPermSize=256m"))
			Ω(props.Uaa.RequireHttps).Should(BeTrue())
			Ω(props.Uaa.Url).Should(Equal("https://uaa.sys.test.com"))
			Ω(props.Uaa.Ssl).ShouldNot(BeNil())
			Ω(props.Uaa.Ssl.Port).Should(Equal(-1))

			Ω(props.Uaa.Authentication).ShouldNot(BeNil())
			Ω(props.Uaa.Authentication.Policy).ShouldNot(BeNil())
			Ω(props.Uaa.Authentication.Policy.LockoutAfterFailures).Should(Equal(5))

			Ω(props.Uaa.Password).ShouldNot(BeNil())
			Ω(props.Uaa.Password.Policy).ShouldNot(BeNil())
			Ω(props.Uaa.Password.Policy.MinLength).Should(Equal(0))
			Ω(props.Uaa.Password.Policy.RequireLowerCaseCharacter).Should(Equal(0))
			Ω(props.Uaa.Password.Policy.RequireUpperCaseCharacter).Should(Equal(0))
			Ω(props.Uaa.Password.Policy.RequireDigit).Should(Equal(0))
			Ω(props.Uaa.Password.Policy.RequireSpecialCharacter).Should(Equal(0))
			Ω(props.Uaa.Password.Policy.ExpirePasswordInMonths).Should(Equal(0))

			Ω(props.Uaa.Jwt).ShouldNot(BeNil())
			Ω(props.Uaa.Jwt.SigningKey).Should(Equal("jwt-signingkey"))
			Ω(props.Uaa.Jwt.VerificationKey).Should(Equal("jwt-verificationkey"))

			Ω(props.Uaa.Ldap).ShouldNot(BeNil())
			Ω(props.Uaa.Ldap.Enabled).Should(BeTrue())
			Ω(props.Uaa.Ldap.Url).Should(Equal("ldap://ldap.test.com"))
			Ω(props.Uaa.Ldap.UserDN).Should(Equal("userdn"))
			Ω(props.Uaa.Ldap.UserPassword).Should(Equal("userpwd"))
			Ω(props.Uaa.Ldap.SearchBase).Should(Equal("base"))
			Ω(props.Uaa.Ldap.SearchFilter).Should(Equal("filter"))
			Ω(props.Uaa.Ldap.MailAttributeName).Should(Equal("mail"))
			Ω(props.Uaa.Ldap.ProfileType).Should(Equal("search-and-bind"))
			Ω(props.Uaa.Ldap.SslCertificate).Should(Equal(""))
			Ω(props.Uaa.Ldap.SslCertificateAlias).Should(Equal(""))
			Ω(props.Uaa.Ldap.Groups).ShouldNot(BeNil())
			Ω(props.Uaa.Ldap.Groups.ProfileType).Should(Equal("no-groups"))
			Ω(props.Uaa.Ldap.Groups.SearchBase).Should(Equal(""))
			Ω(props.Uaa.Ldap.Groups.GroupSearchFilter).Should(Equal(""))
		})
		It("then it should then have uaa job with valid login information", func() {
			ig := uaaPartition.ToInstanceGroup()
			job := ig.GetJobByName("uaa")
			Ω(job).ShouldNot(BeNil())
			props, _ := job.Properties.(*uaa.UaaJob)
			Ω(props.Login).ShouldNot(BeNil())
			Ω(props.Login.SelfServiceLinksEnabled).Should(BeTrue())
			Ω(props.Login.SignupsEnabled).Should(BeTrue())
			Ω(props.Login.Protocol).Should(Equal("https"))
			Ω(props.Login.UaaBase).Should(Equal("https://uaa.sys.test.com"))
			Ω(props.Login.Branding).ShouldNot(BeNil())

			Ω(props.Login.Links).ShouldNot(BeNil())
			links := props.Login.Links
			Ω(links.Passwd).Should(Equal("https://login.sys.test.com/forgot_password"))
			Ω(links.Signup).Should(Equal("https://login.sys.test.com/create_account"))

			Ω(props.Login.Notifications).ShouldNot(BeNil())
			Ω(props.Login.Notifications.Url).Should(Equal("https://notifications.sys.test.com"))

			Ω(props.Login.Saml).ShouldNot(BeNil())
			Ω(props.Login.Saml.Entityid).Should(Equal("https://login.sys.test.com"))
			Ω(props.Login.Saml.SignRequest).Should(BeTrue())
			Ω(props.Login.Saml.WantAssertionSigned).Should(BeFalse())
			Ω(props.Login.Saml.ServiceProviderKey).Should(Equal("saml-key"))
			Ω(props.Login.Saml.ServiceProviderCertificate).Should(Equal("saml-cert"))

			Ω(props.Login.Logout).ShouldNot(BeNil())
			Ω(props.Login.Logout.Redirect).ShouldNot(BeNil())
			Ω(props.Login.Logout.Redirect.Url).Should(Equal("/login"))
			Ω(props.Login.Logout.Redirect.Parameter).ShouldNot(BeNil())
			Ω(props.Login.Logout.Redirect.Parameter.Disable).Should(BeFalse())
			Ω(props.Login.Logout.Redirect.Parameter.Whitelist).Should(ConsistOf("https://console.sys.test.com", "https://apps.sys.test.com"))
		})
		It("then it should then have route_registrar job", func() {
			ig := uaaPartition.ToInstanceGroup()
			job := ig.GetJobByName("route_registrar")
			Ω(job).ShouldNot(BeNil())
			props, _ := job.Properties.(*route_registrar.RouteRegistrarJob)
			Ω(props.Nats).ShouldNot(BeNil())
			Ω(props.Nats.User).Should(Equal("nats"))
			Ω(props.Nats.Password).Should(Equal("pass"))
			Ω(props.Nats.Port).Should(Equal(4222))
			Ω(props.Nats.Machines).Should(ConsistOf("1.0.0.5", "1.0.0.6"))
			Ω(props.RouteRegistrar.Routes).ShouldNot(BeNil())
			routes := props.RouteRegistrar.Routes.([]map[string]interface{})
			Ω(routes).Should(HaveLen(1))
			route := routes[0]
			Ω(route["name"]).Should(Equal("uaa"))
			Ω(route["port"]).Should(Equal(8080))
			Ω(route["registration_interval"]).Should(Equal("40s"))
			Ω(route["uris"]).Should(ConsistOf("uaa.sys.test.com", "*.uaa.sys.test.com", "login.sys.test.com", "*.login.sys.test.com"))
		})
	})
})
