package prabbitmq

import (
	"github.com/enaml-ops/enaml"
	rmqs "github.com/enaml-ops/omg-product-bundle/products/p-rabbitmq/enaml-gen/rabbitmq-server"
)

func (p *Plugin) NewRabbitMQServerPartition(c *Config) *enaml.InstanceGroup {
	return &enaml.InstanceGroup{
		Name:      "rabbitmq-server-partition",
		Lifecycle: "service",
		Instances: len(c.ServerIPs),
		Networks: []enaml.Network{
			{Name: c.Network, StaticIPs: c.ServerIPs},
		},
		Jobs: []enaml.InstanceJob{
			{
				Name:    "rabbitmq-server",
				Release: CFRabbitMQReleaseName, // TODO CHECK ME
				Properties: &rmqs.RabbitmqServerJob{
					RabbitmqServer: &rmqs.RabbitmqServer{
						Ssl: &rmqs.Ssl{
							Verify:            false,
							VerificationDepth: 5,
							FailIfNoPeerCert:  false,
						},
						ClusterPartitionHandling: "pause_minority",
					},
					SyslogAggregator: &rmqs.SyslogAggregator{
						Address: c.SyslogAddress,
						Port:    c.SyslogPort,
					},
				},
			},
		},
	}
}
