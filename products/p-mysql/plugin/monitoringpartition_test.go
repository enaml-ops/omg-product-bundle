package pmysql_test

import (
	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-product-bundle/products/p-mysql/enaml-gen/replication-canary"
	. "github.com/enaml-ops/omg-product-bundle/products/p-mysql/plugin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("given NewMonitoringPartition func", func() {
	Describe("given a valid plugin object argument", func() {
		var plgn *Plugin
		BeforeEach(func() {
			plgn = new(Plugin)
		})

		Context("when the plugin sets monitoring ip values", func() {
			var ig *enaml.InstanceGroup
			var controlIPs = []string{
				"1.0.0.1", "1.0.0.2", "1.0.0.3",
			}
			BeforeEach(func() {
				plgn.MonitoringIPs = controlIPs
				ig = NewMonitoringPartition(plgn)
			})

			It("then it should create vm instances for each given ip", func() {
				Ω(ig.Instances).Should(Equal(len(controlIPs)), "does not create the correct number of instances")
			})

			It("then it should create a static ip block for the given ips", func() {
				Ω(len(ig.Networks)).Should(Equal(1))
				Ω(ig.Networks[0].StaticIPs).Should(Equal(controlIPs), "does not create the correct network ip set")
			})
		})

		Context("when creating the monitoring-partition instancegroup", func() {
			var ig *enaml.InstanceGroup

			BeforeEach(func() {
				ig = NewMonitoringPartition(plgn)
			})

			It("then it should contain instance jobs for all required releases", func() {
				Ω(len(ig.Jobs)).Should(Equal(1))
				Ω(checkJobExists(ig.Jobs, "replication-canary")).Should(BeTrue(), "there should be a job for replication-canary release")
			})
		})

		Describe("given a valid replication-canary job", func() {
			Context("when configured properly in the instance group", func() {
				var ig *enaml.InstanceGroup
				var jobProperties *replication_canary.ReplicationCanaryJob
				var controlIPs = []string{
					"1.0.0.1", "1.0.0.2", "1.0.0.3",
				}
				var controlBaseDomain = "bleh.com"
				var controlSysDomain = "sys." + controlBaseDomain
				var controlAddress = "address"
				var controlPort = "port"
				var controlTransport = "transport"
				var controlEmail = "lkahslkdghasd@aol.com"
				var controlUaaSecret = "lkjasdlfka"
				var controlDBPass = "aslkhaoihsoiehbishdb"
				var controlMonitoringSecret = "baklsjdhglakshdlkgsadg"
				var controlAPIUser = "api-user"
				var controlAPIPass = "api-pass"

				BeforeEach(func() {
					plgn.SyslogAddress = controlAddress
					plgn.SyslogPort = controlPort
					plgn.SyslogTransport = controlTransport
					plgn.MonitoringIPs = controlIPs
					plgn.BaseDomain = controlBaseDomain
					plgn.NotificationRecipientEmail = controlEmail
					plgn.IPs = controlIPs
					plgn.UaaAdminClientSecret = controlUaaSecret
					plgn.SeededDBPassword = controlDBPass
					plgn.NotificationClientSecret = controlMonitoringSecret
					plgn.ProxyAPIUser = controlAPIUser
					plgn.ProxyAPIPass = controlAPIPass

					ig = NewMonitoringPartition(plgn)
					Ω(len(ig.Jobs)).Should(BeNumerically(">=", 1))
					jobProperties = ig.GetJobByName("replication-canary").Properties.(*replication_canary.ReplicationCanaryJob)
				})

				It("then it should have a complete set of properties", func() {
					Ω(*jobProperties).ShouldNot(Equal(replication_canary.ReplicationCanaryJob{}), "object should have properly initialized canary job property values")
				})

				It("then it should have a valid domain property", func() {
					Ω(jobProperties.Domain).Should(Equal(controlSysDomain))
				})

				It("then it should have a valid SyslogAggregator property", func() {
					Ω(jobProperties.SyslogAggregator.Address).Should(Equal(controlAddress), "does not set a valid syslog address")
					Ω(jobProperties.SyslogAggregator.Port).Should(Equal(controlPort), "does not set a valid syslog port")
					Ω(jobProperties.SyslogAggregator.Transport).Should(Equal(controlTransport), "does not set a valid syslog transport")
				})

				It("then it should have a valid MysqlMonitoring property", func() {
					Ω(jobProperties.MysqlMonitoring.RecipientEmail).Should(Equal(controlEmail), "the monitoring recipient email goes here")
					Ω(jobProperties.MysqlMonitoring.NotifyOnly).Should(BeTrue(), "only notify value")
					Ω(jobProperties.MysqlMonitoring.ReplicationCanary.PollFrequency).Should(Equal(30), "how frequently should we poll")
					Ω(jobProperties.MysqlMonitoring.ReplicationCanary.WriteReadDelay).Should(Equal(20), "delay for write/read")
					Ω(jobProperties.MysqlMonitoring.ReplicationCanary.SwitchboardCount).Should(Equal(2), "count of switchboards")
					Ω(jobProperties.MysqlMonitoring.ReplicationCanary.ClusterIps).Should(Equal(controlIPs), "the cluster ip addresses")
					Ω(jobProperties.MysqlMonitoring.ReplicationCanary.UaaAdminClientSecret).Should(Equal(controlUaaSecret), "uaa admin client secret")
					Ω(jobProperties.MysqlMonitoring.ReplicationCanary.CanaryUsername).Should(Equal("repcanary"), "username for the canary")
					Ω(jobProperties.MysqlMonitoring.ReplicationCanary.CanaryPassword).Should(Equal(controlDBPass), "password for the canary")
					Ω(jobProperties.MysqlMonitoring.ReplicationCanary.NotificationsClientUsername).Should(Equal("mysql-monitoring"), "notifications client user")
					Ω(jobProperties.MysqlMonitoring.ReplicationCanary.NotificationsClientSecret).Should(Equal(controlMonitoringSecret), "notifications client secret")
					Ω(jobProperties.MysqlMonitoring.ReplicationCanary.SwitchboardUsername).Should(Equal(controlAPIUser), "switchboard user")
					Ω(jobProperties.MysqlMonitoring.ReplicationCanary.SwitchboardPassword).Should(Equal(controlAPIPass), "swtichboard user's password")
				})
			})
		})
	})
})
