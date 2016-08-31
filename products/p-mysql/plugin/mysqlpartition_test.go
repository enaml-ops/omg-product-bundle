package pmysql_test

import (
	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-product-bundle/products/p-mysql/enaml-gen/mysql"
	"github.com/enaml-ops/omg-product-bundle/products/p-mysql/enaml-gen/send-email"
	"github.com/enaml-ops/omg-product-bundle/products/p-mysql/enaml-gen/streaming-mysql-backup-tool"
	. "github.com/enaml-ops/omg-product-bundle/products/p-mysql/plugin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("given NewMysqlPartition func", func() {
	Describe("given a valid plugin object argument", func() {
		var plgn *Plugin
		BeforeEach(func() {
			plgn = new(Plugin)
		})

		Context("when creating the mysql-partition instancegroup", func() {
			var ig *enaml.InstanceGroup

			BeforeEach(func() {
				ig = NewMysqlPartition(plgn)
			})

			It("then it should contain instance jobs for all required releases", func() {
				Ω(len(ig.Jobs)).Should(Equal(3))
				Ω(checkJobExists(ig.Jobs, "mysql")).Should(BeTrue(), "there should be a job for cf-mysql release")
				Ω(checkJobExists(ig.Jobs, "streaming-mysql-backup-tool")).Should(BeTrue(), "there should be a job for mysql-backup release ")
				Ω(checkJobExists(ig.Jobs, "send-email")).Should(BeTrue(), "there should be a job for mysql-monitoring")
			})
		})

		Describe("given a mysql-backup release w/ instancejob streaming-mysql-backup-tool", func() {
			var ig *enaml.InstanceGroup
			var jobProperties *streaming_mysql_backup_tool.StreamingMysqlBackupToolJob

			Context("when initialized properly", func() {
				BeforeEach(func() {
					ig = NewMysqlPartition(plgn)
					Ω(len(ig.Jobs)).Should(BeNumerically(">=", 1))
					jobProperties = ig.GetJobByName("streaming-mysql-backup-tool").Properties.(*streaming_mysql_backup_tool.StreamingMysqlBackupToolJob)
				})
				It("then it should contain a valid job definition", func() {
					Ω(*jobProperties.CfMysqlBackup).ShouldNot(Equal(streaming_mysql_backup_tool.CfMysqlBackup{}), "cfmysql-backup should be properly populated")
					Ω(*jobProperties.CfMysql).ShouldNot(Equal(streaming_mysql_backup_tool.CfMysql{}), "cfmysql should be properly populated")
				})
			})
			Context("when we have a properly configured cf-mysql block", func() {
				var controlAdminPass = "admin-pass"

				BeforeEach(func() {
					plgn.AdminPassword = controlAdminPass
					ig = NewMysqlPartition(plgn)
					Ω(len(ig.Jobs)).Should(BeNumerically(">=", 1))
					jobProperties = ig.GetJobByName("streaming-mysql-backup-tool").Properties.(*streaming_mysql_backup_tool.StreamingMysqlBackupToolJob)
				})
				It("then it should contain valid mysql creds", func() {
					Ω(jobProperties.CfMysql.Mysql.AdminPassword).Should(Equal(controlAdminPass))
				})
			})

			Context("when we have a properly configured cf-mysql-backup block", func() {
				var controlBackupPass = "bu-pass"
				var controlBackupUser = "bu-user"

				BeforeEach(func() {
					plgn.BackupEndpointUser = controlBackupUser
					plgn.BackupEndpointPassword = controlBackupPass
					ig = NewMysqlPartition(plgn)
					Ω(len(ig.Jobs)).Should(BeNumerically(">=", 1))
					jobProperties = ig.GetJobByName("streaming-mysql-backup-tool").Properties.(*streaming_mysql_backup_tool.StreamingMysqlBackupToolJob)
				})
				It("then it should contain endpoint credentials", func() {
					Ω(jobProperties.CfMysqlBackup.EndpointCredentials.Username).Should(Equal(controlBackupUser), "we should set via a flag to the plugin a backup rest endpoint user")
					Ω(jobProperties.CfMysqlBackup.EndpointCredentials.Password).Should(Equal(controlBackupPass), "we should set via a flag to the plugin a backup rest endpoint password")
				})
			})
		})

		Describe("given a mysql-monitoring release w/ instancejob send-mail", func() {
			var ig *enaml.InstanceGroup
			var jobProperties *send_email.SendEmailJob
			var controlDomainArg = "bleh.blah.com"
			var controlDomain = "sys." + controlDomainArg

			Context("when initialized properly", func() {
				BeforeEach(func() {
					plgn.BaseDomain = controlDomainArg
					ig = NewMysqlPartition(plgn)
					Ω(len(ig.Jobs)).Should(BeNumerically(">=", 1))
					jobProperties = ig.GetJobByName("send-email").Properties.(*send_email.SendEmailJob)
				})
				It("then it should contain a valid job definition", func() {
					Ω(jobProperties.Ssl.SkipCertVerify).Should(BeTrue(), "we should be setting ssl cert verify to true")
					Ω(jobProperties.Domain).Should(Equal(controlDomain), "we should have a valid constructed sys domain")
					Ω(jobProperties.MysqlMonitoring).ShouldNot(BeNil(), "mysql monitoring should def not be nil")
					Ω(*jobProperties.MysqlMonitoring).ShouldNot(Equal(send_email.MysqlMonitoring{}), "mysql monitoring should def not be nil")
				})
			})

			Context("when a mysql-monitoring section is defined", func() {
				var jobProperties *send_email.SendEmailJob
				var controlEmail = "email@org.com"
				var controlUaaSecret = "uaa-secret"
				var controlClientSecret = "client-secret"
				BeforeEach(func() {
					plgn.NotificationRecipientEmail = controlEmail
					plgn.NotificationClientSecret = controlClientSecret
					plgn.UaaAdminClientSecret = controlUaaSecret
					ig := NewMysqlPartition(plgn)
					Ω(len(ig.Jobs)).Should(BeNumerically(">=", 1))
					jobProperties = ig.GetJobByName("send-email").Properties.(*send_email.SendEmailJob)
				})
				It("then it should contain all required fields", func() {
					Ω(jobProperties.MysqlMonitoring.RecipientEmail).Should(Equal(controlEmail), "we should be setting the notification email based on user flag input")
					Ω(jobProperties.MysqlMonitoring.AdminClient.Secret).Should(Equal(controlUaaSecret), "we should be setting the uaa notification secret from user input")
					Ω(jobProperties.MysqlMonitoring.Client.Secret).Should(Equal(controlClientSecret), "we should be setting the notification client secret from user input")
				})
			})
		})

		Describe("given a cf-mysql release w/ instancejob mysql", func() {

			Context("when the plugin gives credential and healthcheck details", func() {
				var ig *enaml.InstanceGroup
				var jobProperties *mysql.MysqlJob
				var controlGaleraHealthcheckUsername = "galera-healthcheck-username"
				var controlGaleraHealthcheckPassword = "galera-healthcheck-password"
				var controlGaleraHealthcheckDBPassword = "galera-healthcheck-db-password"
				var controlClusterHealthPassword = "cluster-health-password"

				BeforeEach(func() {
					plgn.GaleraHealthcheckUsername = controlGaleraHealthcheckUsername
					plgn.GaleraHealthcheckPassword = controlGaleraHealthcheckPassword
					plgn.GaleraHealthcheckDBPassword = controlGaleraHealthcheckDBPassword
					plgn.ClusterHealthPassword = controlClusterHealthPassword
					ig = NewMysqlPartition(plgn)
					Ω(len(ig.Jobs)).Should(BeNumerically(">=", 1))
					jobProperties = ig.GetJobByName("mysql").Properties.(*mysql.MysqlJob)
				})

				It("then it should not be nil", func() {
					Ω(jobProperties.CfMysql).ShouldNot(BeNil(), "cfmysql should contain endpoint and credential information for your mysql")
				})

				It("then it should contain healthcheck config info", func() {
					m := jobProperties.CfMysql.Mysql
					Ω(m.GaleraHealthcheck.EndpointUsername).Should(Equal(controlGaleraHealthcheckUsername), "we should have a username for our galera healthcheck endpoint")
					Ω(m.GaleraHealthcheck.EndpointPassword).Should(Equal(controlGaleraHealthcheckPassword), "we should have a password for our galera healthcheck endpoint")
					Ω(m.GaleraHealthcheck.DbPassword).Should(Equal(controlGaleraHealthcheckDBPassword), "we should have a db password for our galera healthcheck")
					Ω(m.ClusterHealth.Password).Should(Equal(controlClusterHealthPassword), "we should have a password fo our cluster health")
				})
			})

			Context("when the plugin sets a syslog endpoint", func() {
				var ig *enaml.InstanceGroup
				var controlAddress = "address"
				var controlPort = "port"
				var controlTransport = "transport"
				var jobProperties *mysql.MysqlJob

				BeforeEach(func() {
					plgn.SyslogAddress = controlAddress
					plgn.SyslogPort = controlPort
					plgn.SyslogTransport = controlTransport
					ig = NewMysqlPartition(plgn)
					Ω(len(ig.Jobs)).Should(BeNumerically(">=", 1))
					jobProperties = ig.GetJobByName("mysql").Properties.(*mysql.MysqlJob)
				})

				It("then it should create a syslog address record", func() {
					Ω(jobProperties.SyslogAggregator.Address).Should(Equal(controlAddress), "does not set a valid syslog address")
				})

				It("then it should create a syslog port record", func() {
					Ω(jobProperties.SyslogAggregator.Port).Should(Equal(controlPort), "does not set a valid syslog port")
				})

				It("then it should create a syslog transport record", func() {
					Ω(jobProperties.SyslogAggregator.Transport).Should(Equal(controlTransport), "does not set a valid syslog transport")
				})
			})

			Context("when the plugin sets passwords", func() {
				var ig *enaml.InstanceGroup
				var controlSeedPass = "seed-pass"
				var controlAdminPass = "admin-pass"
				var jobProperties *mysql.MysqlJob

				BeforeEach(func() {
					plgn.SeededDBPassword = controlSeedPass
					plgn.AdminPassword = controlAdminPass
					ig = NewMysqlPartition(plgn)
					Ω(len(ig.Jobs)).Should(BeNumerically(">=", 1))
					jobProperties = ig.GetJobByName("mysql").Properties.(*mysql.MysqlJob)
				})

				It("then it should create a record using the given seed pass", func() {
					Ω(jobProperties.AdminPassword).Should(Equal(controlAdminPass), "does not set a valid admin pass")
				})

				It("then it should create a record using the given admin pass", func() {
					Ω(jobProperties.SeededDatabases.([]map[string]string)[0]).Should(HaveKeyWithValue("password", controlSeedPass), "does not set a valid seeded pass")
				})
			})

			Context("when the plugin sets mysql ip values", func() {
				var ig *enaml.InstanceGroup
				var controlIPs = []string{
					"1.0.0.1", "1.0.0.2", "1.0.0.3",
				}
				BeforeEach(func() {
					plgn.IPs = controlIPs
					ig = NewMysqlPartition(plgn)
				})

				It("then it should create vm instances for each given ip", func() {
					Ω(ig.Instances).Should(Equal(len(controlIPs)), "does not create the correct number of instances")
				})

				It("then it should create a static ip block for the given ips", func() {
					Ω(len(ig.Networks)).Should(Equal(1))
					Ω(ig.Networks[0].StaticIPs).Should(Equal(controlIPs), "does not create the correct network ip set")
				})
			})
		})
	})
})
