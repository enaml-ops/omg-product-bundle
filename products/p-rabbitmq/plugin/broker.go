package prabbitmq

import (
	"fmt"

	"github.com/enaml-ops/enaml"
	rmqb "github.com/enaml-ops/omg-product-bundle/products/p-rabbitmq/enaml-gen/rabbitmq-broker"
	sm "github.com/enaml-ops/omg-product-bundle/products/p-rabbitmq/enaml-gen/service-metrics"
)

func (p *Plugin) NewRabbitMQBrokerPartition(c *Config) *enaml.InstanceGroup {
	return &enaml.InstanceGroup{
		Name:      "rabbitmq-broker-partition",
		Lifecycle: "service",
		Instances: 1,
		Stemcell:  StemcellAlias,
		VMType:    c.BrokerVMType,
		AZs:       c.AZs,
		Networks: []enaml.Network{
			{
				Name:      c.Network,
				StaticIPs: []string{c.BrokerIP},
				Default:   []interface{}{"dns", "gateway"},
			},
		},
		Jobs: []enaml.InstanceJob{
			newRabbitMQBrokerJob(c),
			newMetronAgentJob(c),
			newServiceMetricsBrokerJob(c),
			enaml.InstanceJob{
				Name:       "rabbitmq-broker-metrics",
				Release:    RabbitMQMetricsReleaseName,
				Properties: struct{}{},
			},
		},
	}
}

func newRabbitMQBrokerJob(c *Config) enaml.InstanceJob {
	return enaml.InstanceJob{
		Name:    "rabbitmq-broker",
		Release: CFRabbitMQReleaseName,
		Properties: &rmqb.RabbitmqBrokerJob{
			RabbitmqBroker: &rmqb.RabbitmqBroker{
				Route:      "pivotal-rabbitmq-broker",
				Ip:         c.BrokerIP,
				CcEndpoint: fmt.Sprintf("https://api.%s", c.SystemDomain),
				Rabbitmq: &rmqb.Rabbitmq{
					OperatorSetPolicy: &rmqb.OperatorSetPolicy{
						Enabled:          false,
						PolicyName:       "operator_set_policy",
						PolicyDefinition: `{"ha-mode": "exactly", "ha-params": 2, "ha-sync-mode": "automatic"}`,
						PolicyPriority:   50,
					},
					ManagementDomain: fmt.Sprintf("pivotal-rabbitmq.%s", c.SystemDomain),
					Hosts:            []string{c.PublicIP},
					Administrator: &rmqb.Administrator{
						Username: "broker",
						Password: c.BrokerPassword,
					},
				},
				Service: &rmqb.Service{
					Url:      c.BrokerIP,
					Username: "admin",
					Password: c.ServiceAdminPassword,
				},
				Logging: &rmqb.Logging{
					Level:            "info",
					PrintStackTraces: true,
				},
			},
			SyslogAggregator: &rmqb.SyslogAggregator{
				Address: c.SyslogAddress,
				Port:    c.SyslogPort,
			},
			Cf: &rmqb.Cf{
				Domain: c.SystemDomain,
				Nats: &rmqb.Nats{
					Machines: c.NATSMachines,
					Port:     c.NATSPort,
					Username: "nats",
					Password: c.NATSPassword,
				},
			},
		},
	}
}

func newServiceMetricsBrokerJob(c *Config) enaml.InstanceJob {
	return enaml.InstanceJob{
		Name:    "service-metrics",
		Release: ServiceMetricsReleaseName,
		Properties: &sm.ServiceMetricsJob{
			ServiceMetrics: &sm.ServiceMetrics{
				ExecutionIntervalSeconds: 30,
				Origin:         c.DeploymentName,
				MetricsCommand: "/var/vcap/packages/rabbitmq-broker-metrics/heartbeat.sh",
				MetricsCommandArgs: []string{
					"admin",
					c.ServiceAdminPassword,
				},
			},
		},
	}
}
