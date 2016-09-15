package cloudfoundry_test

import (
	//"fmt"

	"io/ioutil"

	. "github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/plugin"
	"github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/plugin/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v2"

	ccnglib "github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/cloud_controller_ng"
	"github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/consul_agent"
	"github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/route_registrar"
)

var _ = Describe("Cloud Controller Partition", func() {
	Context("When initialized with a complete set of arguments", func() {
		var cloudController InstanceGroupCreator

		BeforeEach(func() {

			config := &config.Config{
				NATSPort:           4333,
				SystemDomain:       "sys.yourdomain.com",
				AppDomains:         []string{"apps.yourdomain.com"},
				AllowSSHAccess:     true,
				NetworkName:        "foundry",
				SkipSSLCertVerify:  true,
				HostKeyFingerprint: "hostkeyfingerprint",
				SharePath:          "/var/vcap/nfs",
				DopplerZone:        "DopplerZoneguid",
				SupportAddress:     "http://support.pivotal.io",
				MinCliVersion:      "6.7.0",
				Secret:             config.Secret{},
				User:               config.User{},
				Certs:              &config.Certs{},
				InstanceCount:      config.InstanceCount{},
				IP:                 config.IP{},
				LoggregatorPort:    443,
			}
			config.CCInternalAPIUser = "internalapiuser"
			config.NATSMachines = []string{"10.0.0.4"}
			config.NATSUser = "natsuser"
			config.NATSPassword = "natspass"
			config.ConsulIPs = []string{"1.0.0.1", "1.0.0.2"}
			config.ConsulEncryptKeys = []string{"consulencryptionkey"}
			config.ConsulAgentCert = "consul-agent-cert"
			config.ConsulAgentKey = "consul-agent-key"
			config.ConsulServerCert = "consulservercert"
			config.ConsulServerKey = "consulserverkey"
			config.CloudControllerVMType = "ccvmtype"
			config.CloudControllerInstances = 1
			config.StagingUploadUser = "staginguser"
			config.StagingUploadPassword = "stagingpassword"
			config.CCBulkAPIPassword = "bulkapipassword"
			config.DbEncryptionKey = "dbencryptionkey"
			config.CCInternalAPIPassword = "internalapipassword"
			config.NFSIP = "10.0.0.19"
			config.MySQLProxyIPs = []string{"10.0.0.3"}
			config.CCDBUsername = "ccdbuser"
			config.CCDBPassword = "ccdbpass"
			config.JWTVerificationKey = "uaajwtkey"
			config.CCServiceDashboardsClientSecret = "ccdashboardsecret"
			config.CloudControllerUsernameLookupClientSecret = "usernamelookupsecret"
			config.CCRoutingClientSecret = "ccroutingsecret"
			config.DopplerSharedSecret = "metronsecret"

			cloudController = NewCloudControllerPartition(config)
		})

		It("then should not be nil", func() {
			Ω(cloudController).ShouldNot(BeNil())
		})

		It("then it should configure 1 instance by default", func() {
			ig := cloudController.ToInstanceGroup()
			Ω(ig.Instances).Should(Equal(1))
		})

		It("should have the name of the Network correctly set", func() {
			igf := cloudController.ToInstanceGroup()

			networks := igf.Networks
			Ω(len(networks)).Should(Equal(1))
			Ω(networks[0].Name).Should(Equal("foundry"))
		})

		It("should have 13 jobs under it", func() {
			igf := cloudController.ToInstanceGroup()
			jobs := igf.Jobs
			jobNames := []string{}
			for _, job := range jobs {
				jobNames = append(jobNames, job.Name)
			}
			Ω(len(jobs)).Should(Equal(14))
			Ω(jobNames).Should(ConsistOf("cloud_controller_ng",
				"consul_agent",
				"nfs_mounter",
				"metron_agent",
				"statsd-injector",
				"route_registrar",
				"go-buildpack",
				"binary-buildpack",
				"nodejs-buildpack",
				"ruby-buildpack",
				"php-buildpack",
				"python-buildpack",
				"java-offline-buildpack",
				"staticfile-buildpack"))
		})

		It("should have configured the route_registrar job", func() {
			igf := cloudController.ToInstanceGroup()
			job := igf.GetJobByName("route_registrar")
			Ω(job).ShouldNot(BeNil())
			Ω(job.Release).Should(Equal(CFReleaseName))

			props := job.Properties.(route_registrar.RouteRegistrarJob)
			routes := props.RouteRegistrar.Routes.([]map[string]interface{})
			Ω(routes).Should(HaveLen(1))
			Ω(routes[0]).Should(HaveKeyWithValue("name", "api"))
			Ω(routes[0]).Should(HaveKeyWithValue("port", 9022))
			Ω(routes[0]).Should(HaveKeyWithValue("registration_interval", "20s"))
			Ω(routes[0]).Should(HaveKey("tags"))
			Ω(routes[0]["tags"]).Should(HaveKeyWithValue("component", "CloudController"))
			Ω(routes[0]).Should(HaveKey("uris"))
			Ω(routes[0]["uris"]).Should(ConsistOf("api.sys.yourdomain.com"))
			nats := props.Nats
			Ω(nats).ShouldNot(BeNil())
			Ω(nats.User).Should(Equal("natsuser"))
			Ω(nats.Password).Should(Equal("natspass"))
			Ω(nats.Port).Should(Equal(4333))
			Ω(nats.Machines).Should(ConsistOf("10.0.0.4"))
		})

		It("should have configured the cloud_controller_ng job", func() {
			igf := cloudController.ToInstanceGroup()
			job := igf.GetJobByName("cloud_controller_ng")
			Ω(job).ShouldNot(BeNil())
			Ω(job.Release).Should(Equal(CFReleaseName))

			props := job.Properties.(*ccnglib.CloudControllerNgJob)
			Ω(props.AppSsh.HostKeyFingerprint).Should(Equal("hostkeyfingerprint"))
			Ω(props.Domain).Should(Equal("sys.yourdomain.com"))
			Ω(props.SystemDomain).Should(Equal("sys.yourdomain.com"))
			Ω(props.SystemDomainOrganization).Should(Equal("system"))
			Ω(props.Login.Url).Should(Equal("https://login.sys.yourdomain.com"))

			By("configuring CC")
			Ω(props.Cc.AllowedCorsDomains).Should(ConsistOf("https://login.sys.yourdomain.com"))
			Ω(props.Cc.AllowAppSshAccess).Should(BeTrue())
			Ω(props.Cc.DefaultToDiegoBackend).Should(BeTrue())
			Ω(props.Cc.ClientMaxBodySize).Should(Equal("1024M"))
			Ω(props.Cc.ExternalProtocol).Should(Equal("https"))
			Ω(props.Cc.LoggingLevel).Should(Equal("info"))
			Ω(props.Cc.MaximumHealthCheckTimeout).Should(Equal(600))
			Ω(props.Cc.StagingUploadUser).Should(Equal("staginguser"))
			Ω(props.Cc.StagingUploadPassword).Should(Equal("stagingpassword"))
			Ω(props.Cc.BulkApiUser).Should(BeNil())
			Ω(props.Cc.BulkApiPassword).Should(Equal("bulkapipassword"))
			Ω(props.Cc.InternalApiUser).Should(Equal("internalapiuser"))
			Ω(props.Cc.InternalApiPassword).Should(Equal("internalapipassword"))
			Ω(props.Cc.DbEncryptionKey).Should(Equal("dbencryptionkey"))
			Ω(props.Cc.DefaultRunningSecurityGroups).Should(ConsistOf("all_open"))
			Ω(props.Cc.DefaultStagingSecurityGroups).Should(ConsistOf("all_open"))
			Ω(props.Cc.DisableCustomBuildpacks).Should(BeFalse())
			Ω(props.Cc.ExternalHost).Should(Equal("api"))
			Ω(props.Cc.QuotaDefinitions).Should(HaveKey("default"))
			Ω(props.Cc.QuotaDefinitions).Should(HaveKey("runaway"))
			Ω(props.Cc.SecurityGroupDefinitions).Should(HaveLen(1))
			sg := props.Cc.SecurityGroupDefinitions.([]map[string]interface{})
			Ω(sg[0]).Should(HaveKeyWithValue("name", "all_open"))
			Ω(sg[0]).Should(HaveKey("rules"))
			Ω(sg[0]["rules"]).Should(HaveLen(1))
			rules := sg[0]["rules"].([]map[string]interface{})
			Ω(rules[0]).Should(HaveKeyWithValue("protocol", "all"))
			Ω(rules[0]).Should(HaveKeyWithValue("destination", "0.0.0.0-255.255.255.255"))
			stacks := props.Cc.Stacks.([]map[string]interface{})
			Ω(stacks).Should(HaveLen(2))
			Ω(stacks[0]).Should(HaveKeyWithValue("name", "cflinuxfs2"))
			Ω(stacks[1]).Should(HaveKeyWithValue("name", "windows2012R2"))
			Ω(props.Cc.UaaResourceId).Should(Equal("cloud_controller,cloud_controller_service_permissions"))

			expectedFog := &ccnglib.DefaultFogConnection{
				Provider:  "Local",
				LocalRoot: "/var/vcap/nfs/shared",
			}
			Ω(props.Cc.Buildpacks.BlobstoreType).Should(Equal("fog"))
			Ω(props.Cc.Droplets.BlobstoreType).Should(Equal("fog"))
			Ω(props.Cc.Packages.BlobstoreType).Should(Equal("fog"))
			Ω(props.Cc.ResourcePool.BlobstoreType).Should(Equal("fog"))
			Ω(props.Cc.Buildpacks.FogConnection).Should(Equal(expectedFog))
			Ω(props.Cc.Droplets.FogConnection).Should(Equal(expectedFog))
			Ω(props.Cc.Packages.FogConnection).Should(Equal(expectedFog))
			Ω(props.Cc.ResourcePool.FogConnection).Should(Equal(expectedFog))

			By("configuring CCDB roles")
			ccdb := props.Ccdb
			Ω(ccdb.DbScheme).Should(Equal("mysql"))
			Ω(ccdb.Port).Should(Equal(3306))
			Ω(ccdb.Address).Should(Equal("10.0.0.3"))

			roles := ccdb.Roles.([]map[string]interface{})
			Ω(roles).Should(HaveLen(1))
			Ω(roles[0]).Should(HaveKeyWithValue("name", "ccdbuser"))
			Ω(roles[0]).Should(HaveKeyWithValue("password", "ccdbpass"))
			Ω(roles[0]).Should(HaveKeyWithValue("tag", "admin"))

			dbs := ccdb.Databases.([]map[string]interface{})
			Ω(dbs).Should(HaveLen(1))
			Ω(dbs[0]).Should(HaveKeyWithValue("citext", true))
			Ω(dbs[0]).Should(HaveKeyWithValue("name", "ccdb"))
			Ω(dbs[0]).Should(HaveKeyWithValue("tag", "cc"))

			By("configuring UAA")
			Ω(props.Uaa).ShouldNot(BeNil())
			Ω(props.Uaa.Url).Should(Equal("https://uaa.sys.yourdomain.com"))
			Ω(props.Uaa.Jwt.VerificationKey).Should(Equal("uaajwtkey"))
			Ω(props.Uaa.Clients).ShouldNot(BeNil())
			Ω(props.Uaa.Clients.CcServiceDashboards).ShouldNot(BeNil())
			Ω(props.Uaa.Clients.CcServiceDashboards.Scope).Should(Equal("cloud_controller.write,openid,cloud_controller.read,cloud_controller_service_permissions.read"))
			Ω(props.Uaa.Clients.CcServiceDashboards.Secret).Should(Equal("ccdashboardsecret"))
			Ω(props.Uaa.Clients.CloudControllerUsernameLookup).ShouldNot(BeNil())
			Ω(props.Uaa.Clients.CloudControllerUsernameLookup.Secret).Should(Equal("usernamelookupsecret"))
			Ω(props.Uaa.Clients.CcRouting).ShouldNot(BeNil())
			Ω(props.Uaa.Clients.CcRouting.Secret).Should(Equal("ccroutingsecret"))

			By("configuring SSL")
			Ω(props.Ssl).ShouldNot(BeNil())
			Ω(props.Ssl.SkipCertVerify).Should(BeTrue())

			By("configuring the logger endpoint")
			Ω(props.LoggerEndpoint).ShouldNot(BeNil())
			Ω(props.LoggerEndpoint.Port).Should(Equal(443))

			By("configuring doppler")
			Ω(props.Doppler).ShouldNot(BeNil())
			Ω(props.Doppler.Port).Should(Equal(443))

			By("configuring nfs server")
			Ω(props.NfsServer).ShouldNot(BeNil())
			Ω(props.NfsServer.Address).Should(Equal("10.0.0.19"))
			Ω(props.NfsServer.SharePath).Should(Equal("/var/vcap/nfs"))

			By("configuring nats")
			Ω(props.Nats).ShouldNot(BeNil())
			Ω(props.Nats.User).Should(Equal("natsuser"))
			Ω(props.Nats.Password).Should(Equal("natspass"))
			Ω(props.Nats.Port).Should(Equal(4333))
			Ω(props.Nats.Machines).Should(ConsistOf("10.0.0.4"))

			By("setting installed buildpacks")
			buildpacks, err := ioutil.ReadFile("fixtures/install_buildpacks.yml")
			Ω(err).ShouldNot(HaveOccurred())
			yml, err := yaml.Marshal(props.Cc.InstallBuildpacks)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(yml).Should(MatchYAML(buildpacks))
		})

		It("should have configured the consul agent job", func() {
			igf := cloudController.ToInstanceGroup()
			job := igf.GetJobByName("consul_agent")
			Ω(job).ShouldNot(BeNil())
			Ω(job.Release).Should(Equal(CFReleaseName))

			props := job.Properties.(*consul_agent.ConsulAgentJob)
			Ω(props.Consul.Agent.Domain).Should(Equal("cf.internal"))
			Ω(props.Consul.Agent.Services).Should(HaveKey("cloud_controller_ng"))
		})

		It("should have NFS Mounter set as a job", func() {
			igf := cloudController.ToInstanceGroup()
			nfsMounter := igf.Jobs[2]
			Ω(nfsMounter.Name).Should(Equal("nfs_mounter"))
		})

		It("should have NFS Mounter details set properly", func() {
			igf := cloudController.ToInstanceGroup()

			b, _ := yaml.Marshal(igf)
			Ω(string(b)).Should(ContainSubstring("https://login.sys.yourdomain.com"))
		})

		XIt("should account for QuotaDefinitions structure", func() {
			igf := cloudController.ToInstanceGroup()
			Ω(igf.Jobs[0].Name).Should(Equal("cloud_controller_worker"))
			ccNg, typecasted := igf.Jobs[0].Properties.(*ccnglib.CloudControllerNgJob)
			Ω(typecasted).Should(BeTrue())

			_, quotaTypeCasted := ccNg.Cc.QuotaDefinitions.([]string)
			Ω(quotaTypeCasted).Should(BeFalse())
		})

		XIt("should account for InstallBuildPacks structure", func() {
			igf := cloudController.ToInstanceGroup()
			Ω(igf.Jobs[0].Name).Should(Equal("cloud_controller_worker"))
			ccNg, typecasted := igf.Jobs[0].Properties.(*ccnglib.CloudControllerNgJob)
			Ω(typecasted).Should(BeTrue())

			_, bpTypecasted := ccNg.Cc.InstallBuildpacks.([]string)
			Ω(bpTypecasted).Should(BeFalse())
		})

		XIt("should account for SecurityGroupDefinitions structure", func() {
			igf := cloudController.ToInstanceGroup()
			Ω(igf.Jobs[0].Name).Should(Equal("cloud_controller_worker"))
			ccNg, typecasted := igf.Jobs[0].Properties.(*ccnglib.CloudControllerNgJob)
			Ω(typecasted).Should(BeTrue())

			_, securityTypcasted := ccNg.Cc.SecurityGroupDefinitions.([]string)
			Ω(securityTypcasted).Should(BeFalse())
		})
	})
})
