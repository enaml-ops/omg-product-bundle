package pmysql_test

import (
	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-product-bundle/products/p-mysql/enaml-gen/proxy"
	. "github.com/enaml-ops/omg-product-bundle/products/p-mysql/plugin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("given NewProxyPartition func", func() {
	Describe("given a valid plugin object argument", func() {
		var plgn *Plugin
		BeforeEach(func() {
			plgn = new(Plugin)
		})

		Context("when creating the proxy-partition instancegroup", func() {
			var ig *enaml.InstanceGroup

			BeforeEach(func() {
				ig = NewProxyPartition(plgn)
			})

			It("then it should contain instance jobs for all required releases", func() {
				Ω(len(ig.Jobs)).Should(Equal(1))
				Ω(checkJobExists(ig.Jobs, "proxy")).Should(BeTrue(), "there should be a job for proxy release")
			})
		})

		Describe("given a proxy instance job", func() {
			Context("when configured properly", func() {

				var ig *enaml.InstanceGroup
				var jobProperties *proxy.ProxyJob
				var controlNatsPort = "4222"
				var controlNatsUser = "nats-user"
				var controlNatsPass = "nats-pass"
				var controlBaseDomain = "bleh.blah.com"
				var controlExternalHost = "p-mysql.sys." + controlBaseDomain
				var controlProxyIPs = []string{
					"1.0.0.5", "1.0.0.6",
				}
				var controlIPs = []string{
					"1.0.0.1", "1.0.0.2", "1.0.0.3",
				}
				var controlNatsIPs = []string{
					"1.0.0.7", "1.0.0.8", "1.0.0.9",
				}
				var controlAddress = "address"
				var controlPort = "port"
				var controlTransport = "transport"
				var controlAPIUser = "api-user"
				var controlAPIPass = "api-pass"

				BeforeEach(func() {
					plgn.SyslogAddress = controlAddress
					plgn.SyslogPort = controlPort
					plgn.SyslogTransport = controlTransport
					plgn.BaseDomain = controlBaseDomain
					plgn.ProxyIPs = controlProxyIPs
					plgn.IPs = controlIPs
					plgn.NatsIPs = controlNatsIPs
					plgn.NatsUser = controlNatsUser
					plgn.NatsPassword = controlNatsPass
					plgn.NatsPort = controlNatsPort
					plgn.ProxyAPIUser = controlAPIUser
					plgn.ProxyAPIPass = controlAPIPass
					ig = NewProxyPartition(plgn)
					Ω(len(ig.Jobs)).Should(BeNumerically(">=", 1))
					jobProperties = ig.GetJobByName("proxy").Properties.(*proxy.ProxyJob)
				})
				It("then it should create a valid properties definition", func() {
					Ω(jobProperties.ExternalHost).ShouldNot(BeNil(), "we should have a external host defined")
					Ω(jobProperties.ClusterIps).ShouldNot(BeNil(), "we should have a cluster ips defined")
					Ω(jobProperties.Nats).ShouldNot(BeNil(), "we should have a nats defined")
					Ω(jobProperties.SyslogAggregator).ShouldNot(BeNil(), "we should have a syslog aggregator defined")
					Ω(jobProperties.Proxy).ShouldNot(BeNil(), "we should have a proxy defined")
				})

				It("then it should have a valid external host", func() {
					Ω(jobProperties.ExternalHost).Should(Equal(controlExternalHost), "we should have a valid host for the mysql proxies to be accessible on")
				})

				It("then it should have a valid list of cluster ips", func() {
					Ω(jobProperties.ClusterIps).Should(Equal(controlIPs), "we should have a cluster ips defined")
				})

				It("then it should have a valid nats configuration", func() {
					Ω(jobProperties.Nats.Machines).Should(Equal(controlNatsIPs), "we should have a nats proxy ip list defined")
					Ω(jobProperties.Nats.Password).Should(Equal(controlNatsPass), "we should have a nats password defined")
					Ω(jobProperties.Nats.Port).Should(Equal(controlNatsPort), "we should have a nats port defined")
					Ω(jobProperties.Nats.User).Should(Equal(controlNatsUser), "we should have a nats user defined")
				})

				It("then it should have a valid syslog aggregator configuration", func() {
					Ω(jobProperties.SyslogAggregator.Address).Should(Equal(controlAddress), "does not set a valid syslog address")
					Ω(jobProperties.SyslogAggregator.Port).Should(Equal(controlPort), "does not set a valid syslog port")
					Ω(jobProperties.SyslogAggregator.Transport).Should(Equal(controlTransport), "does not set a valid syslog transport")
				})

				It("then it should have a valid proxy configuration", func() {
					Ω(jobProperties.Proxy.ProxyIps).Should(Equal(controlProxyIPs), "we should have a proxy ip list defined")
					Ω(jobProperties.Proxy.ApiPassword).Should(Equal(controlAPIPass), "we should have a proxy pass defined")
					Ω(jobProperties.Proxy.ApiUsername).Should(Equal(controlAPIUser), "we should have a proxy user defined")
				})
			})
		})

		Context("when the plugin sets proxy ip values", func() {
			var ig *enaml.InstanceGroup
			var controlIPs = []string{
				"1.0.0.1", "1.0.0.2", "1.0.0.3",
			}
			BeforeEach(func() {
				plgn.ProxyIPs = controlIPs
				ig = NewProxyPartition(plgn)
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
