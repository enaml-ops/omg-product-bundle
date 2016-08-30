package cloudfoundry_test

import (
	"github.com/codegangsta/cli"
	. "github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/plugin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Config", func() {
	Context("when initialized WITHOUT a complete set of arguments", func() {
		It("then should return error", func() {
			plugin := new(Plugin)
			c := plugin.GetContext([]string{
				"cloudfoundry",
				"--az", "z1",
			})
			config, err := NewConfig(c)
			Ω(err).Should(HaveOccurred())
			Ω(config).Should(BeNil())
		})
	})
	Context("when initialized WITH a complete set of arguments", func() {
		var config *Config
		var err error
		var c *cli.Context
		BeforeEach(func() {
			plugin := new(Plugin)
			c = plugin.GetContext([]string{
				"cloudfoundry",
				"--az", "z1",
				"--network", "theNetwork",
				"--system-domain", "sys.domain",
				"--app-domain", "app.domain",
				"--nats-machine-ip", "10.0.0.10",
				"--nats-machine-ip", "10.0.0.11",
				"--mysql-bootstrap-password", "mysqlbootstrappwd",
				"--mysql-ip", "10.0.0.12",
				"--mysql-ip", "10.0.0.13",
			})
			config, err = NewConfig(c)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(config).ShouldNot(BeNil())
		})

		It("then az should be set", func() {
			Ω(config.AZs).Should(ConsistOf("z1"))
		})

		It("then network should be set", func() {
			Ω(config.NetworkName).Should(Equal("theNetwork"))
		})

		It("then system domain should be set", func() {
			Ω(config.SystemDomain).Should(Equal("sys.domain"))
		})

		It("then apps domain should be set", func() {
			Ω(config.AppDomains).Should(ConsistOf("app.domain"))
		})

		It("then nats ips should be set", func() {
			Ω(config.NATSMachines).Should(ConsistOf("10.0.0.10", "10.0.0.11"))
		})

		It("then mysql ips should be set", func() {
			Ω(config.MySQLIPs).Should(ConsistOf("10.0.0.12", "10.0.0.13"))
		})

		It("then apps domain should be set", func() {
			Ω(config.MySQLBootstrapPassword).Should(Equal("mysqlbootstrappwd"))
		})

	})
})
