package prabbitmq

import (
	"fmt"

	"github.com/enaml-ops/enaml"
	br "github.com/enaml-ops/omg-product-bundle/products/p-rabbitmq/enaml-gen/broker-registrar"
)

func (p *Plugin) NewRabbitMQBrokerRegistrar(c *Config) *enaml.InstanceGroup {
	return &enaml.InstanceGroup{
		Name:      "broker-registrar",
		Lifecycle: "errand",
		Instances: 1,
		VMType:    c.BrokerVMType,
		AZs:       c.AZs,
		Stemcell:  StemcellAlias,
		Networks: []enaml.Network{
			{
				Name:    c.Network,
				Default: []interface{}{"dns", "gateway"},
			},
		},
		Jobs: []enaml.InstanceJob{
			{
				Name:    "broker-registrar",
				Release: CFRabbitMQReleaseName,
				Properties: &br.BrokerRegistrarJob{
					Broker: &br.Broker{
						Name:     "p-rabbitmq",
						Host:     fmt.Sprintf("pivotal-rabbitmq-broker.%s", c.SystemDomain),
						Username: "admin",
						Password: c.ServiceAdminPassword,
					},
					Cf: &br.Cf{
						ApiUrl:            fmt.Sprintf("https://api.%s", c.SystemDomain),
						AdminUsername:     "system_services",
						AdminPassword:     c.SystemServicesPassword,
						SkipSslValidation: c.SkipSSLVerify,
					},
				},
			},
		},
	}
}
