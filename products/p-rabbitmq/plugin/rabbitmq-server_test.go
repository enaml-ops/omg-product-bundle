package prabbitmq_test

import (
	"github.com/enaml-ops/enaml"
	rmqs "github.com/enaml-ops/omg-product-bundle/products/p-rabbitmq/enaml-gen/rabbitmq-server"
	prabbitmq "github.com/enaml-ops/omg-product-bundle/products/p-rabbitmq/plugin"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("RabbitMQ server partition", func() {
	Context("when initialized with a complete configuration", func() {
		var ig *enaml.InstanceGroup

		const (
			controlNetworkName    = "foundry-net"
			controlSyslogAddress  = "1.2.3.4"
			controlSyslogPort     = 1234
			controlBrokerPassword = "brokerpassword"
		)

		BeforeEach(func() {
			p := new(prabbitmq.Plugin)
			c := &prabbitmq.Config{
				ServerIPs:      []string{"10.0.1.2", "10.0.1.3"},
				Network:        controlNetworkName,
				SyslogAddress:  controlSyslogAddress,
				SyslogPort:     controlSyslogPort,
				BrokerPassword: controlBrokerPassword,
			}
			ig = p.NewRabbitMQServerPartition(c)
			Ω(ig).ShouldNot(BeNil())
		})

		It("should configure the instance group parameters", func() {
			Ω(ig.Name).Should(Equal("rabbitmq-server-partition"))
			Ω(ig.Lifecycle).Should(Equal("service"))
			Ω(ig.Instances).Should(Equal(2))
			Ω(ig.Networks).Should(HaveLen(1))
			Ω(ig.Networks[0].Name).Should(Equal(controlNetworkName))
			Ω(ig.Networks[0].StaticIPs).Should(ConsistOf("10.0.1.2", "10.0.1.3"))
		})

		It("should configure the rabbitmq-server job", func() {
			Ω(ig.Jobs).Should(HaveLen(1))
			Ω(ig.Jobs[0].Properties).ShouldNot(BeNil())
			Ω(ig.Jobs[0].Name).Should(Equal("rabbitmq-server"))
			Ω(ig.Jobs[0].Release).Should(Equal(prabbitmq.CFRabbitMQReleaseName))

			props := ig.Jobs[0].Properties.(*rmqs.RabbitmqServerJob)
			Ω(props.RabbitmqServer).ShouldNot(BeNil())

			Ω(props.RabbitmqServer.Ssl).ShouldNot(BeNil())
			Ω(props.RabbitmqServer.Ssl.Verify).Should(BeFalse())
			Ω(props.RabbitmqServer.Ssl.VerificationDepth).Should(Equal(5))
			Ω(props.RabbitmqServer.Ssl.FailIfNoPeerCert).Should(BeFalse())

			Ω(props.RabbitmqServer.ClusterPartitionHandling).Should(Equal("pause_minority"))

			Ω(props.SyslogAggregator).ShouldNot(BeNil())
			Ω(props.SyslogAggregator.Address).Should(Equal(controlSyslogAddress))
			Ω(props.SyslogAggregator.Port).Should(Equal(controlSyslogPort))

			Ω(props.RabbitmqServer.Administrators).ShouldNot(BeNil())
			Ω(props.RabbitmqServer.Administrators.Management).ShouldNot(BeNil())
			Ω(props.RabbitmqServer.Administrators.Management.Username).Should(Equal("rabbitadmin"))
			Ω(props.RabbitmqServer.Administrators.Management.Password).Should(Equal("rabbitadmin"))
			Ω(props.RabbitmqServer.Administrators.Broker).ShouldNot(BeNil())
			Ω(props.RabbitmqServer.Administrators.Broker.Username).Should(Equal("broker"))
			Ω(props.RabbitmqServer.Administrators.Broker.Password).Should(Equal(controlBrokerPassword))

		})
	})
})
