package prabbitmq_test

import (
	"github.com/enaml-ops/enaml"
	ma "github.com/enaml-ops/omg-product-bundle/products/p-rabbitmq/enaml-gen/metron_agent"
	rmqh "github.com/enaml-ops/omg-product-bundle/products/p-rabbitmq/enaml-gen/rabbitmq-haproxy"
	sm "github.com/enaml-ops/omg-product-bundle/products/p-rabbitmq/enaml-gen/service-metrics"
	prabbitmq "github.com/enaml-ops/omg-product-bundle/products/p-rabbitmq/plugin"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("rabbitmq haproxy partition", func() {
	Context("when initialized with a complete configuration", func() {
		var ig *enaml.InstanceGroup

		const (
			controlDeploymentName = "p-rabbitmq"
			controlNetworkName    = "foundry-net"
			controlSyslogAddress  = "1.2.3.4"
			controlSyslogPort     = 1234
			controlBrokerPassword = "brokerpassword"
			controlPublicIP       = "10.0.1.10"
			controlStatsPassword  = "haproxystatspassword"
			controlNATSPort       = 4333
			controlNATSPassword   = "natspassword"
			controlNATSIP         = "10.0.0.2"
			controlMetronZone     = "metronzone"
			controlMetronSecret   = "metronsharedsecret"
		)

		BeforeEach(func() {
			p := new(prabbitmq.Plugin)
			c := &prabbitmq.Config{
				DeploymentName:            controlDeploymentName,
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
				MetronZone:                controlMetronZone,
				MetronSecret:              controlMetronSecret,
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
			job := ig.GetJobByName("rabbitmq-haproxy")
			Ω(job).ShouldNot(BeNil())
			Ω(job.Properties).ShouldNot(BeNil())
			Ω(job.Name).Should(Equal("rabbitmq-haproxy"))
			Ω(job.Release).Should(Equal(prabbitmq.CFRabbitMQReleaseName))

			props := job.Properties.(*rmqh.RabbitmqHaproxyJob)
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
