package prabbitmq

import (
	"github.com/enaml-ops/enaml"
	ma "github.com/enaml-ops/omg-product-bundle/products/p-rabbitmq/enaml-gen/metron_agent"
	sm "github.com/enaml-ops/omg-product-bundle/products/p-rabbitmq/enaml-gen/service-metrics"

	rmqs "github.com/enaml-ops/omg-product-bundle/products/p-rabbitmq/enaml-gen/rabbitmq-server"
)

func (p *Plugin) NewRabbitMQServerPartition(c *Config) *enaml.InstanceGroup {
	return &enaml.InstanceGroup{
		Name:      "rabbitmq-server-partition",
		Lifecycle: "service",
		Stemcell:  StemcellAlias,
		VMType:    c.ServerVMType,
		AZs:       c.AZs,
		Instances: len(c.ServerIPs),
		Networks: []enaml.Network{
			{Name: c.Network, StaticIPs: c.ServerIPs},
		},
		Jobs: []enaml.InstanceJob{
			newRabbitMQServerJob(c),
			newMetronAgentJob(c),
			newServiceMetricsServerJob(c),
			enaml.InstanceJob{
				Name:       "rabbitmq-server-metrics",
				Release:    RabbitMQMetricsReleaseName,
				Properties: struct{}{},
			},
		},
	}
}

func newRabbitMQServerJob(c *Config) enaml.InstanceJob {
	return enaml.InstanceJob{
		Name:    "rabbitmq-server",
		Release: CFRabbitMQReleaseName,
		Properties: &rmqs.RabbitmqServerJob{
			RabbitmqServer: &rmqs.RabbitmqServer{
				Ssl: &rmqs.Ssl{
					Verify:            false,
					VerificationDepth: 5,
					FailIfNoPeerCert:  false,
				},
				ClusterPartitionHandling: "pause_minority",
				Administrators: &rmqs.Administrators{
					Management: &rmqs.Management{
						Username: "rabbitadmin",
						Password: c.AdminPassword,
					},
					Broker: &rmqs.Broker{
						Username: "broker",
						Password: c.BrokerPassword,
					},
				},
				Plugins: []string{"rabbitmq_management"},
			},
			SyslogAggregator: &rmqs.SyslogAggregator{
				Address: c.SyslogAddress,
				Port:    c.SyslogPort,
			},
		},
	}
}

func newMetronAgentJob(c *Config) enaml.InstanceJob {
	return enaml.InstanceJob{
		Name:    "metron_agent",
		Release: LoggregatorReleaseName,
		Properties: &ma.MetronAgentJob{
			MetronAgent: &ma.MetronAgent{
				Zone:       c.MetronZone,
				Deployment: c.DeploymentName,
			},
			MetronEndpoint: &ma.MetronEndpoint{
				SharedSecret: c.MetronSecret,
			},
			Loggregator: &ma.Loggregator{
				Etcd: &ma.Etcd{
					Machines: c.EtcdMachines,
				},
			},
		},
	}
}

func newServiceMetricsServerJob(c *Config) enaml.InstanceJob {
	return enaml.InstanceJob{
		Name:    "service-metrics",
		Release: ServiceMetricsReleaseName,
		Properties: &sm.ServiceMetricsJob{
			ServiceMetrics: &sm.ServiceMetrics{
				ExecutionIntervalSeconds: 30,
				Origin:         c.DeploymentName,
				MetricsCommand: "/var/vcap/packages/rabbitmq-server-metrics/bin/rabbitmq-server-metrics",
				MetricsCommandArgs: []string{
					"-erlangBinPath=/var/vcap/packages/erlang/bin/",
					"-rabbitmqCtlPath=/var/vcap/packages/rabbitmq-server/bin/rabbitmqctl",
					"-logPath=/var/vcap/sys/log/service-metrics/rabbitmq-server-metrics.log",
					"-rabbitmqUsername=rabbitadmin",
					"-rabbitmqPassword=" + c.AdminPassword,
					"-rabbitmqApiEndpoint=http://127.0.0.1:15672",
				},
			},
		},
	}
}
