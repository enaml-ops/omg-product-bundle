package pmysql_test

import (
	"github.com/enaml-ops/enaml"

	"github.com/enaml-ops/omg-product-bundle/products/p-mysql/enaml-gen/cf-mysql-broker"
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
				plgn.BrokerIPs = controlIPs
				ig = NewCfMysqlBrokerPartition(plgn)
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
				ig = NewCfMysqlBrokerPartition(plgn)

			})

			It("then it should contain instance jobs for all required releases", func() {
				Ω(len(ig.Jobs)).Should(Equal(1))
				Ω(checkJobExists(ig.Jobs, "cf-mysql-broker")).Should(BeTrue(), "there should be a job for broker partition release")
			})
		})

		Describe("given a cf-mysql-broker job in the instance group", func() {
			Context("when given all valid flags and initialized properly", func() {
				var ig *enaml.InstanceGroup
				var jobProperties *cf_mysql_broker.CfMysqlBrokerJob
				var controlQuotaPass = "balkhaslkdhgasdg"
				var controlNetwork = "default"
				var controlBaseDomain = "blah.com"
				var controlHostDomain = "p-mysql.sys." + controlBaseDomain
				var controlAPIDomain = "https://api.sys." + controlBaseDomain
				var controlAuthUser = "lkaslkdfhlaksdf"
				var controlAuthPass = "lkaslkdalksdklklnasdgkn"
				var controlCookieSecret = "lkaslkdghalsdgh"
				var controlNatsPort = "4222"
				var controlNatsUser = "nats-user"
				var controlNatsPass = "nats-pass"
				var controlProxyIPs = []string{
					"1.0.0.5", "1.0.0.6",
				}
				var controlAddress = "address"
				var controlPort = "port"
				var controlTransport = "transport"

				BeforeEach(func() {
					plgn.NetworkName = controlNetwork
					plgn.BrokerQuotaEnforcerPassword = controlQuotaPass
					plgn.BaseDomain = controlBaseDomain
					plgn.BrokerAuthUsername = controlAuthUser
					plgn.BrokerAuthPassword = controlAuthPass
					plgn.BrokerCookieSecret = controlCookieSecret
					plgn.ProxyIPs = controlProxyIPs
					plgn.NatsUser = controlNatsUser
					plgn.NatsPassword = controlNatsPass
					plgn.NatsPort = controlNatsPort
					plgn.SyslogAddress = controlAddress
					plgn.SyslogPort = controlPort
					plgn.SyslogTransport = controlTransport
					ig = NewCfMysqlBrokerPartition(plgn)
					jobProperties = ig.GetJobByName("cf-mysql-broker").Properties.(*cf_mysql_broker.CfMysqlBrokerJob)
				})
				It("then it should have a valid broker element", func() {
					Ω(jobProperties.Broker.QuotaEnforcer.Password).Should(Equal(controlQuotaPass))
					Ω(jobProperties.Broker.QuotaEnforcer.Pause).Should(Equal(30))
				})
				It("then it should have a valid networks element", func() {
					Ω(jobProperties.Networks.BrokerNetwork).Should(Equal(controlNetwork))
				})
				It("then it should have a valid ssl-enabled element", func() {
					Ω(jobProperties.SslEnabled).Should(BeTrue())
				})
				It("then it should have a valid skip-ssl-validation element", func() {
					Ω(jobProperties.SkipSslValidation).Should(BeTrue())
				})
				It("then it should have a valid external-host element", func() {
					Ω(jobProperties.ExternalHost).Should(Equal(controlHostDomain), "route to access mysql from external source")
				})
				It("then it should have a valid cc-api-uri element", func() {
					Ω(jobProperties.CcApiUri).Should(Equal(controlAPIDomain), "cloud controller api endpoint for pcf installation")
				})
				It("then it should have a valid cookie-secret element", func() {
					Ω(jobProperties.CookieSecret).Should(Equal(controlCookieSecret))
				})
				It("then it should have a valid auth-username element", func() {
					Ω(jobProperties.AuthUsername).Should(Equal(controlAuthUser))
				})
				It("then it should have a valid auth-password element", func() {
					Ω(jobProperties.AuthPassword).Should(Equal(controlAuthPass))
				})
				It("then it should have a valid nats element", func() {
					Ω(jobProperties.Nats.Machines).Should(Equal(controlProxyIPs), "we should have a nats proxy ip list defined")
					Ω(jobProperties.Nats.Password).Should(Equal(controlNatsPass), "we should have a nats password defined")
					Ω(jobProperties.Nats.Port).Should(Equal(controlNatsPort), "we should have a nats port defined")
					Ω(jobProperties.Nats.User).Should(Equal(controlNatsUser), "we should have a nats user defined")
				})
				It("then it should have a valid syslog-aggregator element", func() {
					Ω(jobProperties.SyslogAggregator.Address).Should(Equal(controlAddress), "does not set a valid syslog address")
					Ω(jobProperties.SyslogAggregator.Port).Should(Equal(controlPort), "does not set a valid syslog port")
					Ω(jobProperties.SyslogAggregator.Transport).Should(Equal(controlTransport), "does not set a valid syslog transport")
				})
				XIt("then it should have a valid mysql-node element", func() { Ω(true).Should(BeFalse()) })
				XIt("then it should have a valid services element", func() { Ω(true).Should(BeFalse()) })
			})
		})
	})
})
