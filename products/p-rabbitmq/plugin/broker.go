package prabbitmq

import (
	"fmt"

	"github.com/enaml-ops/enaml"
	rmqb "github.com/enaml-ops/omg-product-bundle/products/p-rabbitmq/enaml-gen/rabbitmq-broker"
)

func (p *Plugin) NewRabbitMQBrokerPartition(c *Config) *enaml.InstanceGroup {
	return &enaml.InstanceGroup{
		Name:      "rabbitmq-broker-partition",
		Lifecycle: "service",
		Instances: 1,
		Networks: []enaml.Network{
			{
				Name:      c.Network,
				StaticIPs: []string{c.BrokerIP},
				Default:   []interface{}{"dns", "gateway"},
			},
		},
		Jobs: []enaml.InstanceJob{
			{
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
							Url:      c.ServiceURL,
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
			},
		},
	}
}
