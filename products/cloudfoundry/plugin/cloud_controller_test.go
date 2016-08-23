package cloudfoundry_test

import (
	//"fmt"

	. "github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/plugin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v2"

	ccnglib "github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/cloud_controller_ng"
	"github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/route_registrar"
)

var _ = Describe("Cloud Controller Partition", func() {
	Context("When initialized with a complete set of arguments", func() {
		var cloudController InstanceGrouper

		BeforeEach(func() {
			plugin := new(Plugin)
			c := plugin.GetContext([]string{
				"cloudfoundry",
				"--consul-ip", "1.0.0.1",
				"--consul-ip", "1.0.0.2",
				"--az", "az",
				"--stemcell-name", "stemcell",
				"--consul-encryption-key", "consulencryptionkey",
				"--consul-ca-cert", "consul-ca-cert",
				"--consul-agent-cert", "consul-agent-cert",
				"--consul-agent-key", "consul-agent-key",
				"--consul-server-cert", "consulservercert",
				"--consul-server-key", "consulserverkey",
				"--cc-vm-type", "ccvmtype",
				"--network", "foundry",
				"--host-key-fingerprint", "hostkeyfingerprint",
				"--cc-staging-upload-user", "staginguser",
				"--cc-staging-upload-password", "stagingpassword",
				"--cc-bulk-api-user", "bulkapiuser",
				"--cc-bulk-api-password", "bulkapipassword",
				"--cc-db-encryption-key", "dbencryptionkey",
				"--cc-internal-api-user", "internalapiuser",
				"--cc-internal-api-password", "internalapipassword",
				"--system-domain", "sys.yourdomain.com",
				"--app-domain", "apps.yourdomain.com",
				"--allow-app-ssh-access",
				"--nfs-server-address", "10.0.0.19",
				"--nfs-share-path", "/var/vcap/nfs",
				"--metron-secret", "metronsecret",
				"--metron-zone", "metronzoneguid",
				"--host-key-fingerprint", "hostkeyfingerprint",
				"--support-address", "http://support.pivotal.io",
				"--min-cli-version", "6.7.0",
				"--mysql-proxy-ip", "10.0.0.3",
				"--db-ccdb-username", "ccdbuser",
				"--db-ccdb-password", "ccdbpass",
			})

			cloudController = NewCloudControllerPartition(c)
		})

		It("then should not be nil", func() {
			Ω(cloudController).ShouldNot(BeNil())
		})

		It("should have valid values", func() {
			Ω(cloudController.HasValidValues()).Should(BeTrue())
		})

		It("should have the name of the Network correctly set", func() {
			igf := cloudController.ToInstanceGroup()

			networks := igf.Networks
			Ω(len(networks)).Should(Equal(1))
			Ω(networks[0].Name).Should(Equal("foundry"))
		})

		It("should have 6 jobs under it", func() {
			igf := cloudController.ToInstanceGroup()
			jobs := igf.Jobs
			Ω(len(jobs)).Should(Equal(6))
		})

		It("should have configured the route_registrar job", func() {
			igf := cloudController.ToInstanceGroup()
			job := igf.GetJobByName("route_registrar")
			Ω(job).ShouldNot(BeNil())
			Ω(job.Release).Should(Equal(CFReleaseName))

			props := job.Properties.(route_registrar.RouteRegistrar)
			routes := props.Routes.([]map[string]interface{})
			Ω(routes).Should(HaveLen(1))
			Ω(routes[0]).Should(HaveKeyWithValue("name", "api"))
			Ω(routes[0]).Should(HaveKeyWithValue("port", 9022))
			Ω(routes[0]).Should(HaveKeyWithValue("registration_interval", "20s"))
			Ω(routes[0]).Should(HaveKey("tags"))
			Ω(routes[0]["tags"]).Should(HaveKeyWithValue("component", "CloudController"))
			Ω(routes[0]).Should(HaveKey("uris"))
			Ω(routes[0]["uris"]).Should(ConsistOf("api.sys.yourdomain.com"))
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

			Ω(props.Cc.AllowedCorsDomains).Should(ConsistOf("https://login.sys.yourdomain.com"))
			Ω(props.Cc.AllowAppSshAccess).Should(BeTrue())
			Ω(props.Cc.DefaultToDiegoBackend).Should(BeTrue())
			Ω(props.Cc.ClientMaxBodySize).Should(Equal("1024M"))
			Ω(props.Cc.ExternalProtocol).Should(Equal("https"))
			Ω(props.Cc.LoggingLevel).Should(Equal("debug"))
			Ω(props.Cc.MaximumHealthCheckTimeout).Should(Equal(600))
			Ω(props.Cc.StagingUploadUser).Should(Equal("staginguser"))
			Ω(props.Cc.StagingUploadPassword).Should(Equal("stagingpassword"))
			Ω(props.Cc.BulkApiUser).Should(Equal("bulkapiuser"))
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
