package config_test

import (
	"github.com/enaml-ops/omg-product-bundle/products/oss_cf/plugin"
	. "github.com/enaml-ops/omg-product-bundle/products/oss_cf/plugin/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/xchapter7x/lo"
	"gopkg.in/urfave/cli.v2"
)

func BuildConfigContext() *cli.Context {
	plugin := new(cloudfoundry.Plugin)
	c := plugin.GetContext([]string{
		"cloudfoundry",
		"--az", "z1",
		"--network", "theNetwork",
		"--system-domain", "sys.yourdomain.com",
		"--app-domain", "app.domain",
		"--nats-machine-ip", "10.0.0.10",
		"--nats-machine-ip", "10.0.0.11",
		"--mysql-bootstrap-password", "mysqlbootstrappwd",
		"--mysql-ip", "10.0.0.12",
		"--mysql-ip", "10.0.0.13",
		"--syslog-address", "syslog-server",
		"--syslog-port", "10601",
		"--syslog-transport", "tcp",
		"--etcd-machine-ip", "10.0.1.2",
		"--etcd-machine-ip", "10.0.1.3",
		"--etcd-machine-ip", "10.0.1.4",
		"--nats-user", "natsuser",
		"--nats-pass", "natspass",
		"--nats-port", "4222",
		"--loggregator-traffic-controller-vmtype", "vmtype",
		"--consul-agent-cert", "value",
		"--consul-agent-key", "value",
		"--consul-server-cert", "value",
		"--consul-server-key", "value",
		"--bbs-server-ca-cert", "value",
		"--bbs-client-cert", "value",
		"--bbs-client-key", "value",
		"--bbs-server-cert", "value",
		"--bbs-server-key", "value",
		"--etcd-server-cert", "value",
		"--etcd-server-key", "value",
		"--etcd-client-cert", "value",
		"--etcd-client-key", "value",
		"--etcd-peer-cert", "value",
		"--etcd-peer-key", "value",
		"--uaa-saml-service-provider-key", "value",
		"--uaa-saml-service-provider-cert", "value",
		"--uaa-jwt-signing-key", "value",
		"--uaa-jwt-verification-key", "value",
		"--router-ssl-cert", "value",
		"--router-ssl-key", "value",
		"--diego-cell-disk-type", "disk",
		"--diego-brain-disk-type", "disk",
		"--diego-db-disk-type", "disk",
		"--nfs-disk-type", "disk",
		"--etcd-disk-type", "disk",
		"--mysql-disk-type", "disk",
		"--db-uaa-password", "secret",
		"--push-apps-manager-password", "secret",
		"--system-services-password", "secret",
		"--system-verification-password", "secret",
		"--opentsdb-firehose-nozzle-client-secret", "secret",
		"--identity-client-secret", "secret",
		"--login-client-secret", "secret",
		"--portal-client-secret", "secret",
		"--autoscaling-service-client-secret", "secret",
		"--system-passwords-client-secret", "secret",
		"--cc-service-dashboards-client-secret", "secret",
		"--gorouter-client-secret", "secret",
		"--notifications-client-secret", "secret",
		"--notifications-ui-client-secret", "secret",
		"--cloud-controller-username-lookup-client-secret", "secret",
		"--cc-routing-client-secret", "secret",
		"--apps-metrics-client-secret", "secret",
		"--apps-metrics-processing-client-secret", "secret",
		"--admin-password", "secret",
		"--smoke-tests-password", "secret",
		"--doppler-shared-secret", "secret",
		"--doppler-client-secret", "secret",
		"--cc-bulk-api-password", "secret",
		"--cc-internal-api-password", "secret",
		"--ssh-proxy-uaa-secret", "secret",
		"--cc-db-encryption-key", "secret",
		"--db-ccdb-password", "secret",
		"--diego-db-passphrase", "secret",
		"--uaa-admin-secret", "secret",
		"--router-pass", "secret",
		"--mysql-proxy-api-password", "secret",
		"--mysql-admin-password", "secret",
		"--db-console-password", "secret",
		"--cc-staging-upload-password", "secret",
		"--mysql-proxy-vm-type", "vm-type",
		"--clock-global-vm-type", "vm-type",
		"--cc-vm-type", "vm-type",
		"--diego-brain-vm-type", "vm-type",
		"--diego-cell-vm-type", "vm-type",
		"--doppler-vm-type", "vm-type",
		"--cc-worker-vm-type", "vm-type",
		"--errand-vm-type", "vm-type",
		"--etcd-vm-type", "vm-type",
		"--nats-vm-type", "vm-type",
		"--consul-vm-type", "vm-type",
		"--mysql-vm-type", "vm-type",
		"--diego-db-vm-type", "vm-type",
		"--uaa-vm-type", "vm-type",
		"--router-vm-type", "vm-type",
		"--nfs-vm-type", "vm-type",
		"--nfs-ip", "127.0.0.1",
		"--diego-cell-ip", "127.0.0.1",
		"--consul-ip", "127.0.0.1",
		"--doppler-ip", "127.0.0.1",
		"--mysql-proxy-ip", "127.0.0.1",
		"--loggregator-traffic-controller-ip", "127.0.0.1",
		"--diego-brain-ip", "127.0.0.1",
		"--diego-db-ip", "127.0.0.1",
		"--router-ip", "127.0.0.1",
		"--diego-cell-ip", "127.0.0.2",
		"--consul-ip", "127.0.0.2",
		"--doppler-ip", "127.0.0.2",
		"--mysql-proxy-ip", "127.0.0.2",
		"--loggregator-traffic-controller-ip", "127.0.0.2",
		"--diego-brain-ip", "127.0.0.2",
		"--diego-db-ip", "127.0.0.2",
		"--router-ip", "127.0.0.2",
		"--host-key-fingerprint", "value",
		"--doppler-zone", "value",
		"--cc-internal-api-user", "value",
		"--nfs-allow-from-network-cidr", "value",
		"--apps-manager-secret-token", "token",
		"--db-app_usage-password", "appusagepasssword",
		"--db-autoscale-password", "autoscaledbpassword",
		"--db-autoscale-username", "usernameautoscale",
		"--db-notifications-password", "notificationspassword",
	})
	return c
}

func BuildConfig() *Config {
	if config, err := NewConfig(BuildConfigContext()); err == nil {
		return config
	} else {
		lo.G.Error("Error parsing context:", err.Error())
		return nil
	}
}

var _ = Describe("Config", func() {
	Context("when initialized WITHOUT a complete set of arguments", func() {
		It("then should return error", func() {
			plugin := new(cloudfoundry.Plugin)
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
			c = BuildConfigContext()
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
			Ω(config.SystemDomain).Should(Equal("sys.yourdomain.com"))
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

		It("then mysql bootstrap password should be set", func() {
			Ω(config.MySQLBootstrapPassword).Should(Equal("mysqlbootstrappwd"))
		})

	})
})
