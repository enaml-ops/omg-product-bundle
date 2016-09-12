package prabbitmq_test

import (
	"github.com/enaml-ops/enaml"
	ma "github.com/enaml-ops/omg-product-bundle/products/p-rabbitmq/enaml-gen/metron_agent"
	rmqs "github.com/enaml-ops/omg-product-bundle/products/p-rabbitmq/enaml-gen/rabbitmq-server"
	sm "github.com/enaml-ops/omg-product-bundle/products/p-rabbitmq/enaml-gen/service-metrics"
	prabbitmq "github.com/enaml-ops/omg-product-bundle/products/p-rabbitmq/plugin"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("RabbitMQ server partition", func() {
	Context("when initialized with a complete configuration", func() {
		var ig *enaml.InstanceGroup

		const (
			controlDeploymentName = "p-rabbitmq"
			controlNetworkName    = "foundry-net"
			controlSyslogAddress  = "1.2.3.4"
			controlSyslogPort     = 1234
			controlBrokerPassword = "brokerpassword"
			controlMetronZone     = "metronzone"
			controlMetronSecret   = "metronsharedsecret"
		)

		BeforeEach(func() {
			p := new(prabbitmq.Plugin)
			c := &prabbitmq.Config{
				DeploymentName: controlDeploymentName,
				ServerIPs:      []string{"10.0.1.2", "10.0.1.3"},
				Network:        controlNetworkName,
				SyslogAddress:  controlSyslogAddress,
				SyslogPort:     controlSyslogPort,
				BrokerPassword: controlBrokerPassword,
				MetronZone:     controlMetronZone,
				MetronSecret:   controlMetronSecret,
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
			job := ig.GetJobByName("rabbitmq-server")
			Ω(job).ShouldNot(BeNil())
			Ω(job.Properties).ShouldNot(BeNil())
			Ω(job.Name).Should(Equal("rabbitmq-server"))
			Ω(job.Release).Should(Equal(prabbitmq.CFRabbitMQReleaseName))

			props := job.Properties.(*rmqs.RabbitmqServerJob)
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

		It("should configure the metron_agent job", func() {
			job := ig.GetJobByName("metron_agent")
			Ω(job).ShouldNot(BeNil())
			Ω(job.Properties).ShouldNot(BeNil())
			Ω(job.Release).Should(Equal(prabbitmq.LoggregatorReleaseName))

			props := job.Properties.(*ma.MetronAgentJob)
			Ω(props.MetronAgent).ShouldNot(BeNil())
			Ω(props.MetronAgent.Deployment).Should(Equal(controlDeploymentName))
			Ω(props.MetronAgent.Zone).Should(Equal(controlMetronZone))

			Ω(props.MetronEndpoint).ShouldNot(BeNil())
			Ω(props.MetronEndpoint.SharedSecret).Should(Equal(controlMetronSecret))

			Ω(props.Loggregator).ShouldNot(BeNil())
			Ω(props.Loggregator.Etcd).ShouldNot(BeNil())
		})

		It("should configure the service-metrics job", func() {
			job := ig.GetJobByName("service-metrics")
			Ω(job).ShouldNot(BeNil())
			Ω(job.Properties).ShouldNot(BeNil())
			Ω(job.Release).Should(Equal(prabbitmq.ServiceMetricsReleaseName))

			props := job.Properties.(*sm.ServiceMetricsJob)
			Ω(props.ServiceMetrics).ShouldNot(BeNil())
			Ω(props.ServiceMetrics.ExecutionIntervalSeconds).Should(Equal(30))
			Ω(props.ServiceMetrics.Origin).Should(Equal(controlDeploymentName))
			Ω(props.ServiceMetrics.MetricsCommand).Should(Equal("/var/vcap/packages/rabbitmq-server-metrics/bin/rabbitmq-server-metrics"))
			Ω(props.ServiceMetrics.MetricsCommandArgs).Should(ConsistOf(
				"-erlangBinPath=/var/vcap/packages/erlang/bin/",
				"-rabbitmqCtlPath=/var/vcap/packages/rabbitmq-server/bin/rabbitmqctl",
				"-logPath=/var/vcap/sys/log/service-metrics/rabbitmq-server-metrics.log",
				"-rabbitmqUsername=rabbitadmin",
				"-rabbitmqPassword=rabbitadmin",
				"-rabbitmqApiEndpoint=http://127.0.0.1:15672",
			))
		})
	})
})
