package prabbitmq

import (
	"fmt"

	"github.com/enaml-ops/enaml"
	rmqh "github.com/enaml-ops/omg-product-bundle/products/p-rabbitmq/enaml-gen/rabbitmq-haproxy"
	sm "github.com/enaml-ops/omg-product-bundle/products/p-rabbitmq/enaml-gen/service-metrics"
)

func (p *Plugin) NewRabbitMQHAProxyPartition(c *Config) *enaml.InstanceGroup {
	return &enaml.InstanceGroup{
		Name:      "rabbitmq-haproxy-partition",
		Lifecycle: "service",
		Instances: 1,
		Stemcell:  StemcellAlias,
		VMType:    c.HAProxyVMType,
		AZs:       c.AZs,
		Networks: []enaml.Network{
			{
				Name:      c.Network,
				StaticIPs: []string{c.PublicIP},
				Default:   []interface{}{"dns", "gateway"},
			},
		},
		Jobs: []enaml.InstanceJob{
			newRabbitMQHAProxyJob(c),
			newMetronAgentJob(c),
			newServiceMetricsHAProxyJob(c),
		},
	}
}

func newRabbitMQHAProxyJob(c *Config) enaml.InstanceJob {
	return enaml.InstanceJob{
		Name:    "rabbitmq-haproxy",
		Release: CFRabbitMQReleaseName,
		Properties: &rmqh.RabbitmqHaproxyJob{
			RabbitmqHaproxy: &rmqh.RabbitmqHaproxy{
				Stats: &rmqh.Stats{
					Username: "admin",
					Password: c.HAProxyStatsAdminPassword,
				},
				ServerIps: c.ServerIPs,
				Ports:     "15672, 5672, 5671, 1883, 8883, 61613, 61614, 15674",
			},
			RabbitmqBroker: &rmqh.RabbitmqBroker{
				Rabbitmq: &rmqh.Rabbitmq{
					ManagementDomain: fmt.Sprintf("pivotal-rabbitmq.%s", c.SystemDomain),
					ManagementIp:     c.PublicIP,
				},
			},
			Cf: &rmqh.Cf{
				Nats: &rmqh.Nats{
					Username: "nats",
					Password: c.NATSPassword,
					Machines: c.NATSMachines,
					Port:     c.NATSPort,
				},
			},
			SyslogAggregator: &rmqh.SyslogAggregator{
				Address: c.SyslogAddress,
				Port:    c.SyslogPort,
			},
		},
	}
}

func newServiceMetricsHAProxyJob(c *Config) enaml.InstanceJob {
	return enaml.InstanceJob{
		Name:    "service-metrics",
		Release: ServiceMetricsReleaseName,
		Properties: &sm.ServiceMetricsJob{
			ServiceMetrics: &sm.ServiceMetrics{
				ExecutionIntervalSeconds: 30,
				Origin:         c.DeploymentName,
				MetricsCommand: "/var/vcap/packages/rabbitmq-haproxy-metrics/bin/rabbitmq-haproxy-metrics",
				MetricsCommandArgs: []string{
					"-haproxyNetwork=unix",
					"-haproxyAddress=/var/vcap/sys/run/rabbitmq-haproxy/haproxy.sock",
					"-logPath=/var/vcap/sys/log/service-metrics/rabbitmq-haproxy-metrics.log",
				},
			},
		},
	}
}
