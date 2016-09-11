package prabbitmq_test

import (
	"github.com/enaml-ops/enaml"
	rmqb "github.com/enaml-ops/omg-product-bundle/products/p-rabbitmq/enaml-gen/rabbitmq-broker"
	prabbitmq "github.com/enaml-ops/omg-product-bundle/products/p-rabbitmq/plugin"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("rabbitmq-broker partition", func() {
	Context("when initialized with a complete configuration", func() {
		var ig *enaml.InstanceGroup

		const (
			controlNetworkName          = "foundry-net"
			controlBrokerIP             = "1.2.3.4"
			controlPublicIP             = "5.6.7.8"
			controlBrokerPassword       = "brokerpass"
			controlServiceURL           = "10.0.0.1"
			controlServiceAdminPassword = "serviceadminpassword"
			controlSyslogAddress        = "1.2.3.4"
			controlSyslogPort           = 1234
			controlNATSPort             = 4333
			controlNATSPassword         = "natspassword"
			controlNATSIP               = "10.0.0.2"
		)

		BeforeEach(func() {
			p := new(prabbitmq.Plugin)
			c := &prabbitmq.Config{
				Network:              controlNetworkName,
				SystemDomain:         "sys.example.com",
				BrokerIP:             controlBrokerIP,
				PublicIP:             controlPublicIP,
				BrokerPassword:       controlBrokerPassword,
				ServiceURL:           controlServiceURL,
				ServiceAdminPassword: controlServiceAdminPassword,
				SyslogAddress:        controlSyslogAddress,
				SyslogPort:           controlSyslogPort,
				NATSPort:             controlNATSPort,
				NATSPassword:         controlNATSPassword,
				NATSMachines:         []string{controlNATSIP},
			}
			ig = p.NewRabbitMQBrokerPartition(c)
			Ω(ig).ShouldNot(BeNil())
		})

		It("should configure the instance group parameters", func() {
			Ω(ig.Name).Should(Equal("rabbitmq-broker-partition"))
			Ω(ig.Lifecycle).Should(Equal("service"))
			Ω(ig.Instances).Should(Equal(1))
			Ω(ig.Networks).Should(HaveLen(1))
			Ω(ig.Networks[0].Name).Should(Equal(controlNetworkName))
			Ω(ig.Networks[0].StaticIPs).Should(ConsistOf(controlBrokerIP))
			Ω(ig.Networks[0].Default).Should(ConsistOf("dns", "gateway"))
		})

		It("should configure the rabbitmq-broker job", func() {
			Ω(ig.Jobs).Should(HaveLen(1))
			Ω(ig.Jobs[0].Properties).ShouldNot(BeNil())
			Ω(ig.Jobs[0].Name).Should(Equal("rabbitmq-broker"))
			Ω(ig.Jobs[0].Release).Should(Equal(prabbitmq.CFRabbitMQReleaseName))

			props := ig.Jobs[0].Properties.(*rmqb.RabbitmqBrokerJob)
			Ω(props.RabbitmqBroker).ShouldNot(BeNil())

			Ω(props.RabbitmqBroker.Route).Should(Equal("pivotal-rabbitmq-broker"))
			Ω(props.RabbitmqBroker.Ip).Should(Equal(controlBrokerIP))
			Ω(props.RabbitmqBroker.CcEndpoint).Should(Equal("https://api.sys.example.com"))

			Ω(props.RabbitmqBroker.Rabbitmq).ShouldNot(BeNil())
			Ω(props.RabbitmqBroker.Rabbitmq.OperatorSetPolicy).ShouldNot(BeNil())
			Ω(props.RabbitmqBroker.Rabbitmq.OperatorSetPolicy.Enabled).Should(BeFalse())
			Ω(props.RabbitmqBroker.Rabbitmq.OperatorSetPolicy.PolicyName).Should(Equal("operator_set_policy"))
			Ω(props.RabbitmqBroker.Rabbitmq.OperatorSetPolicy.PolicyDefinition).Should(MatchJSON(`{"ha-mode": "exactly", "ha-params": 2, "ha-sync-mode": "automatic"}`))
			Ω(props.RabbitmqBroker.Rabbitmq.OperatorSetPolicy.PolicyPriority).Should(Equal(50))

			Ω(props.RabbitmqBroker.Rabbitmq.ManagementDomain).Should(Equal("pivotal-rabbitmq.sys.example.com"))
			Ω(props.RabbitmqBroker.Rabbitmq.Hosts).Should(ConsistOf(controlPublicIP))

			Ω(props.RabbitmqBroker.Rabbitmq.Administrator).ShouldNot(BeNil())
			Ω(props.RabbitmqBroker.Rabbitmq.Administrator.Username).Should(Equal("broker"))
			Ω(props.RabbitmqBroker.Rabbitmq.Administrator.Password).Should(Equal(controlBrokerPassword))

			Ω(props.RabbitmqBroker.Service).ShouldNot(BeNil())
			Ω(props.RabbitmqBroker.Service.Url).Should(Equal(controlServiceURL))
			Ω(props.RabbitmqBroker.Service.Username).Should(Equal("admin"))
			Ω(props.RabbitmqBroker.Service.Password).Should(Equal(controlServiceAdminPassword))

			Ω(props.RabbitmqBroker.Logging).ShouldNot(BeNil())
			Ω(props.RabbitmqBroker.Logging.Level).Should(Equal("info"))
			Ω(props.RabbitmqBroker.Logging.PrintStackTraces).Should(BeTrue())

			Ω(props.SyslogAggregator).ShouldNot(BeNil())
			Ω(props.SyslogAggregator.Address).Should(Equal(controlSyslogAddress))
			Ω(props.SyslogAggregator.Port).Should(Equal(controlSyslogPort))

			Ω(props.Cf).ShouldNot(BeNil())
			Ω(props.Cf.Domain).Should(Equal("sys.example.com"))

			Ω(props.Cf.Nats).ShouldNot(BeNil())
			Ω(props.Cf.Nats.Machines).Should(ConsistOf(controlNATSIP))
			Ω(props.Cf.Nats.Port).Should(Equal(controlNATSPort))
			Ω(props.Cf.Nats.Username).Should(Equal("nats"))
			Ω(props.Cf.Nats.Password).Should(Equal(controlNATSPassword))
		})

	})
})
