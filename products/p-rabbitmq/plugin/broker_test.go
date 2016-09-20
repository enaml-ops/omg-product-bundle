package prabbitmq_test

import (
	"github.com/enaml-ops/enaml"
	ma "github.com/enaml-ops/omg-product-bundle/products/p-rabbitmq/enaml-gen/metron_agent"
	rmqb "github.com/enaml-ops/omg-product-bundle/products/p-rabbitmq/enaml-gen/rabbitmq-broker"
	sm "github.com/enaml-ops/omg-product-bundle/products/p-rabbitmq/enaml-gen/service-metrics"
	prabbitmq "github.com/enaml-ops/omg-product-bundle/products/p-rabbitmq/plugin"
	yaml "gopkg.in/yaml.v2"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("rabbitmq-broker partition", func() {
	Context("when initialized with a complete configuration", func() {
		var ig *enaml.InstanceGroup

		const (
			controlDeploymentName       = "p-rabbitmq"
			controlNetworkName          = "foundry-net"
			controlAdminPassword        = "rabbitadminpassword"
			controlBrokerIP             = "1.2.3.4"
			controlPublicIP             = "5.6.7.8"
			controlBrokerPassword       = "brokerpass"
			controlServiceAdminPassword = "serviceadminpassword"
			controlSyslogAddress        = "1.2.3.4"
			controlSyslogPort           = 1234
			controlNATSPort             = 4333
			controlNATSPassword         = "natspassword"
			controlNATSIP               = "10.0.0.2"
			controlMetronZone           = "metronzone"
			controlMetronSecret         = "metronsharedsecret"
			controlVMType               = "small"
			controlAZ                   = "az1"
		)

		BeforeEach(func() {
			p := new(prabbitmq.Plugin)
			c := &prabbitmq.Config{
				DeploymentName:       controlDeploymentName,
				AZs:                  []string{controlAZ},
				Network:              controlNetworkName,
				SystemDomain:         "sys.example.com",
				AdminPassword:        controlAdminPassword,
				BrokerIP:             controlBrokerIP,
				PublicIP:             controlPublicIP,
				BrokerPassword:       controlBrokerPassword,
				ServiceAdminPassword: controlServiceAdminPassword,
				SyslogAddress:        controlSyslogAddress,
				SyslogPort:           controlSyslogPort,
				NATSPort:             controlNATSPort,
				NATSPassword:         controlNATSPassword,
				NATSMachines:         []string{controlNATSIP},
				MetronZone:           controlMetronZone,
				MetronSecret:         controlMetronSecret,
				BrokerVMType:         controlVMType,
			}
			ig = p.NewRabbitMQBrokerPartition(c)
			Ω(ig).ShouldNot(BeNil())
		})

		It("should configure the instance group parameters", func() {
			Ω(ig.Name).Should(Equal("rabbitmq-broker-partition"))
			Ω(ig.Lifecycle).Should(Equal("service"))
			Ω(ig.VMType).Should(Equal(controlVMType))
			Ω(ig.Instances).Should(Equal(1))
			Ω(ig.Networks).Should(HaveLen(1))
			Ω(ig.AZs).Should(ConsistOf(controlAZ))
			Ω(ig.Stemcell).Should(Equal(prabbitmq.StemcellAlias))
			Ω(ig.Networks[0].Name).Should(Equal(controlNetworkName))
			Ω(ig.Networks[0].StaticIPs).Should(ConsistOf(controlBrokerIP))
			Ω(ig.Networks[0].Default).Should(ConsistOf("dns", "gateway"))
		})

		It("should configure the rabbitmq-broker job", func() {
			job := ig.GetJobByName("rabbitmq-broker")
			Ω(job).ShouldNot(BeNil())
			Ω(job.Properties).ShouldNot(BeNil())
			Ω(job.Release).Should(Equal(prabbitmq.CFRabbitMQReleaseName))

			props := job.Properties.(*rmqb.RabbitmqBrokerJob)
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
			Ω(props.RabbitmqBroker.Service.Url).Should(Equal(controlBrokerIP))
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
			Ω(props.ServiceMetrics.MetricsCommand).Should(Equal("/var/vcap/packages/rabbitmq-broker-metrics/heartbeat.sh"))
			Ω(props.ServiceMetrics.MetricsCommandArgs).Should(ConsistOf("admin", controlServiceAdminPassword))
		})

		It("should configure the metrics job", func() {
			job := ig.GetJobByName("rabbitmq-broker-metrics")
			Ω(job).ShouldNot(BeNil())
			Ω(job.Release).Should(Equal(prabbitmq.RabbitMQMetricsReleaseName))

			b, err := yaml.Marshal(job.Properties)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(b).Should(MatchYAML(`{}`))
		})
	})
})
