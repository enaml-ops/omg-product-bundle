package prabbitmq_test

import (
	"github.com/enaml-ops/enaml"
	rmqh "github.com/enaml-ops/omg-product-bundle/products/p-rabbitmq/enaml-gen/rabbitmq-haproxy"
	prabbitmq "github.com/enaml-ops/omg-product-bundle/products/p-rabbitmq/plugin"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("rabbitmq haproxy partition", func() {
	Context("when initialized with a complete configuration", func() {
		var ig *enaml.InstanceGroup

		const (
			controlNetworkName    = "foundry-net"
			controlSyslogAddress  = "1.2.3.4"
			controlSyslogPort     = 1234
			controlBrokerPassword = "brokerpassword"
			controlPublicIP       = "10.0.1.10"
			controlStatsPassword  = "haproxystatspassword"
			controlNATSPort       = 4333
			controlNATSPassword   = "natspassword"
			controlNATSIP         = "10.0.0.2"
		)

		BeforeEach(func() {
			p := new(prabbitmq.Plugin)
			c := &prabbitmq.Config{
				ServerIPs:                 []string{"10.0.1.2", "10.0.1.3"},
				Network:                   controlNetworkName,
				SyslogAddress:             controlSyslogAddress,
				SyslogPort:                controlSyslogPort,
				BrokerPassword:            controlBrokerPassword,
				PublicIP:                  controlPublicIP,
				SystemDomain:              "sys.example.com",
				HAProxyStatsAdminPassword: controlStatsPassword,
				NATSPort:                  controlNATSPort,
				NATSPassword:              controlNATSPassword,
				NATSMachines:              []string{controlNATSIP},
			}
			ig = p.NewRabbitMQHAProxyPartition(c)
			Ω(ig).ShouldNot(BeNil())
		})

		It("should configure the instance group parameters", func() {
			Ω(ig.Name).Should(Equal("rabbitmq-haproxy-partition"))
			Ω(ig.Lifecycle).Should(Equal("service"))
			Ω(ig.Instances).Should(Equal(1))
			Ω(ig.Networks).Should(HaveLen(1))
			Ω(ig.Networks[0].Name).Should(Equal(controlNetworkName))
			Ω(ig.Networks[0].StaticIPs).Should(ConsistOf(controlPublicIP))
			Ω(ig.Networks[0].Default).Should(ConsistOf("dns", "gateway"))
		})

		It("should configure the rabbitmq-haproxy job", func() {
			Ω(ig.Jobs).Should(HaveLen(1))
			Ω(ig.Jobs[0].Properties).ShouldNot(BeNil())
			Ω(ig.Jobs[0].Name).Should(Equal("rabbitmq-haproxy"))
			Ω(ig.Jobs[0].Release).Should(Equal(prabbitmq.CFRabbitMQReleaseName))

			props := ig.Jobs[0].Properties.(*rmqh.RabbitmqHaproxyJob)
			Ω(props).ShouldNot(BeNil())
			Ω(props.RabbitmqHaproxy).ShouldNot(BeNil())
			Ω(props.RabbitmqHaproxy.Stats).ShouldNot(BeNil())
			Ω(props.RabbitmqHaproxy.Stats.Username).Should(Equal("admin"))
			Ω(props.RabbitmqHaproxy.Stats.Password).Should(Equal(controlStatsPassword))
			Ω(props.RabbitmqHaproxy.ServerIps).Should(ConsistOf("10.0.1.2", "10.0.1.3"))
			Ω(props.RabbitmqHaproxy.Ports).Should(Equal("15672, 5672, 5671, 1883, 8883, 61613, 61614, 15674"))

			Ω(props.RabbitmqBroker).ShouldNot(BeNil())
			Ω(props.RabbitmqBroker.Rabbitmq).ShouldNot(BeNil())
			Ω(props.RabbitmqBroker.Rabbitmq.ManagementIp).Should(Equal(controlPublicIP))
			Ω(props.RabbitmqBroker.Rabbitmq.ManagementDomain).Should(Equal("pivotal-rabbitmq.sys.example.com"))

			Ω(props.Cf).ShouldNot(BeNil())
			Ω(props.Cf.Nats).ShouldNot(BeNil())
			Ω(props.Cf.Nats.Machines).Should(ConsistOf(controlNATSIP))
			Ω(props.Cf.Nats.Port).Should(Equal(controlNATSPort))
			Ω(props.Cf.Nats.Username).Should(Equal("nats"))
			Ω(props.Cf.Nats.Password).Should(Equal(controlNATSPassword))

			Ω(props.SyslogAggregator).ShouldNot(BeNil())
			Ω(props.SyslogAggregator.Address).Should(Equal(controlSyslogAddress))
			Ω(props.SyslogAggregator.Port).Should(Equal(controlSyslogPort))
		})
	})
})
