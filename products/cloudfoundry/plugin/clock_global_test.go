package cloudfoundry_test

import (
	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/cloud_controller_clock"
	. "github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/plugin"
	"github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/plugin/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v2"
)

var _ = Describe("given a clock_global partition", func() {
	Context("when initialized with a complete set of arguments", func() {
		var ig InstanceGroupCreator
		var deploymentManifest *enaml.DeploymentManifest

		BeforeEach(func() {
			config := &config.Config{
				AZs:               []string{"eastprod-1"},
				StemcellName:      "cool-ubuntu-animal",
				NetworkName:       "foundry-net",
				SkipSSLCertVerify: false,
				AllowSSHAccess:    true,
				SystemDomain:      "sys.test.com",
				AppDomains:        []string{"apps.test.com"},
				NATSPort:          4333,
				MetronZone:        "metronzoneguid",
				SharePath:         "/var/vcap/nfs",
				Secret:            config.Secret{},
				User:              config.User{},
				Certs:             &config.Certs{},
				InstanceCount:     config.InstanceCount{},
				IP:                config.IP{},
			}
			config.CCInternalAPIUser = "internalapiuser"
			config.NFSIP = "10.0.0.19"
			config.NATSUser = "nats"
			config.NATSPassword = "pass"
			config.NATSMachines = []string{"1.0.0.5", "1.0.0.6"}
			config.ClockGlobalVMType = "vmtype"
			config.CloudControllerVMType = "ccvmtype"
			config.StagingUploadUser = "staginguser"
			config.StagingUploadPassword = "stagingpassword"
			config.CCBulkAPIUser = "bulkapiuser"
			config.CCBulkAPIPassword = "bulkapipassword"
			config.DbEncryptionKey = "dbencryptionkey"
			config.CCInternalAPIPassword = "internalapipassword"
			config.CCServiceDashboardsClientSecret = "ccsecret"
			config.MetronSecret = "metronsecret"
			config.ConsulCaCert = "consul-ca-cert"
			config.ConsulAgentCert = "consul-agent-cert"
			config.ConsulAgentKey = "consul-agent-key"
			config.ConsulServerCert = "consulservercert"
			config.ConsulServerKey = "consulserverkey"
			config.ConsulIPs = []string{"1.0.0.1", "1.0.0.2"}
			config.ConsulEncryptKeys = []string{"consulencryptionkey"}
			config.MySQLProxyIPs = []string{"1.0.10.3", "1.0.10.4"}
			config.CCDBUsername = "ccdb-user"
			config.CCDBPassword = "ccdb-password"
			config.JWTVerificationKey = "jwt-verificationkey"

			ig = NewClockGlobalPartition(config)
			deploymentManifest = new(enaml.DeploymentManifest)
			deploymentManifest.AddInstanceGroup(ig.ToInstanceGroup())
		})

		It("then it should allow the user to configure the AZs", func() {
			group := deploymentManifest.GetInstanceGroupByName("clock_global-partition")
			Ω(len(group.AZs)).Should(Equal(1))
			Ω(group.AZs[0]).Should(Equal("eastprod-1"))
		})

		It("then it should have update max in flight 1", func() {
			group := deploymentManifest.GetInstanceGroupByName("clock_global-partition")
			Ω(group.Update.MaxInFlight).Should(Equal(1))
		})

		It("then it should allow the user to configure the used stemcell", func() {
			group := deploymentManifest.GetInstanceGroupByName("clock_global-partition")
			Ω(group.Stemcell).Should(Equal("cool-ubuntu-animal"))
		})

		It("then it should allow the user to configure the network to use", func() {
			group := deploymentManifest.GetInstanceGroupByName("clock_global-partition")
			Ω(len(group.Networks)).Should(Equal(1))
			Ω(group.Networks[0].Name).Should(Equal("foundry-net"))
		})

		It("then it should allow the user to configure the VM type", func() {
			group := deploymentManifest.GetInstanceGroupByName("clock_global-partition")
			Ω(group.VMType).Should(Equal("vmtype"))
		})

		It("then it should have a single instance", func() {
			group := deploymentManifest.GetInstanceGroupByName("clock_global-partition")
			Ω(len(group.AZs)).Should(Equal(1))
		})

		It("should have correctly configured the cloud controller clock", func() {
			group := deploymentManifest.GetInstanceGroupByName("clock_global-partition")
			job := group.GetJobByName("cloud_controller_clock")
			Ω(job.Release).Should(Equal(CFReleaseName))
			props := job.Properties.(*cloud_controller_clock.CloudControllerClockJob)
			Ω(props.Domain).Should(Equal("sys.test.com"))
			Ω(props.SystemDomain).Should(Equal("sys.test.com"))
			Ω(props.SystemDomainOrganization).Should(Equal("system"))

			ad := props.AppDomains.([]string)
			Ω(len(ad)).Should(Equal(1))
			Ω(ad[0]).Should(Equal("apps.test.com"))

			Ω(props.Cc.AllowAppSshAccess).Should(BeTrue())
			Ω(props.Cc.Buildpacks.BlobstoreType).Should(Equal("fog"))
			Ω(props.Cc.Droplets.BlobstoreType).Should(Equal("fog"))
			Ω(props.Cc.Packages.BlobstoreType).Should(Equal("fog"))
			Ω(props.Cc.ResourcePool.BlobstoreType).Should(Equal("fog"))

			Ω(props.Cc.Droplets.FogConnection).Should(HaveKeyWithValue("provider", "Local"))
			Ω(props.Cc.Packages.FogConnection).Should(HaveKeyWithValue("provider", "Local"))
			Ω(props.Cc.ResourcePool.FogConnection).Should(HaveKeyWithValue("provider", "Local"))
			Ω(props.Cc.Droplets.FogConnection).Should(HaveKeyWithValue("local_root", "/var/vcap/nfs/shared"))
			Ω(props.Cc.Packages.FogConnection).Should(HaveKeyWithValue("local_root", "/var/vcap/nfs/shared"))
			Ω(props.Cc.ResourcePool.FogConnection).Should(HaveKeyWithValue("local_root", "/var/vcap/nfs/shared"))

			Ω(props.Cc.LoggingLevel).Should(Equal("debug"))
			Ω(props.Cc.MaximumHealthCheckTimeout).Should(Equal(600))
			Ω(props.Cc.StagingUploadUser).Should(Equal("staginguser"))
			Ω(props.Cc.StagingUploadPassword).Should(Equal("stagingpassword"))
			Ω(props.Cc.BulkApiUser).Should(Equal("bulkapiuser"))
			Ω(props.Cc.BulkApiPassword).Should(Equal("bulkapipassword"))
			Ω(props.Cc.DbEncryptionKey).Should(Equal("dbencryptionkey"))
			Ω(props.Cc.InternalApiUser).Should(Equal("internalapiuser"))
			Ω(props.Cc.InternalApiPassword).Should(Equal("internalapipassword"))

			quotaDefs := `default:
  memory_limit: 10240
  total_services: 100
  non_basic_services_allowed: true
  total_routes: 1000
  trial_db_allowed: true
runaway:
  memory_limit: 102400
  total_services: -1
  total_routes: 1000
  non_basic_services_allowed: true
`
			b, _ := yaml.Marshal(props.Cc.QuotaDefinitions)
			Ω(string(b)).Should(MatchYAML(quotaDefs))

			sgDefs := `- name: all_open
  rules:
    - protocol: all
      destination: 0.0.0.0-255.255.255.255
`
			b, _ = yaml.Marshal(props.Cc.SecurityGroupDefinitions)
			Ω(string(b)).Should(MatchYAML(sgDefs))

			Ω(props.Ccdb.Address).Should(Equal("1.0.10.3"))
			Ω(props.Ccdb.Port).Should(Equal(3306))
			Ω(props.Ccdb.DbScheme).Should(Equal("mysql"))

			Ω(props.Ccdb.Roles).Should(HaveLen(1))
			role := props.Ccdb.Roles.([]map[string]interface{})[0]
			Ω(role).Should(HaveKeyWithValue("tag", "admin"))
			Ω(role).Should(HaveKeyWithValue("name", "ccdb-user"))
			Ω(role).Should(HaveKeyWithValue("password", "ccdb-password"))

			Ω(props.Ccdb.Databases).Should(HaveLen(1))
			db := props.Ccdb.Databases.([]map[string]interface{})[0]
			Ω(db).Should(HaveKeyWithValue("tag", "cc"))
			Ω(db).Should(HaveKeyWithValue("name", "ccdb"))
			Ω(db).Should(HaveKeyWithValue("citext", true))

			Ω(props.Uaa.Url).Should(Equal("https://uaa.sys.test.com"))
			Ω(props.Uaa.Jwt).ShouldNot(BeNil())
			Ω(props.Uaa.Jwt.VerificationKey).Should(Equal("jwt-verificationkey"))

			Ω(props.Uaa.Clients).ShouldNot(BeNil())
			Ω(props.Uaa.Clients.CcServiceDashboards.Secret).Should(Equal("ccsecret"))

			Ω(props.LoggerEndpoint.Port).Should(Equal(443))
			Ω(props.Ssl.SkipCertVerify).Should(BeFalse())

			Ω(props.Nats.User).Should(Equal("nats"))
			Ω(props.Nats.Password).Should(Equal("pass"))
			Ω(props.Nats.Port).Should(Equal(4333))
			Ω(props.Nats.Machines).Should(ConsistOf("1.0.0.5", "1.0.0.6"))
		})
	})
})
