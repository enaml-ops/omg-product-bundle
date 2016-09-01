package pmysql_test

import (
	"github.com/enaml-ops/enaml"
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
	})
})
