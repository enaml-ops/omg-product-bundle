package cloudfoundry

import (
	"github.com/enaml-ops/enaml"
	consullib "github.com/enaml-ops/omg-product-bundle/products/oss_cf/enaml-gen/consul_agent"
	"github.com/enaml-ops/omg-product-bundle/products/oss_cf/plugin/config"
)

//ConsulAgent -
type ConsulAgent struct {
	Config   *config.Config
	Mode     string
	Services []string
}

//NewConsulAgent -
func NewConsulAgent(services []string, config *config.Config) *ConsulAgent {
	ca := &ConsulAgent{
		Config:   config,
		Services: services,
	}
	return ca
}

//NewConsulAgentServer -
func NewConsulAgentServer(config *config.Config) *ConsulAgent {
	ca := &ConsulAgent{
		Config: config,
		Mode:   "server",
	}
	return ca
}

//CreateJob - Create the yaml job structure
func (s *ConsulAgent) CreateJob() enaml.InstanceJob {

	serviceMap := make(map[string]map[string]string)
	for _, serviceName := range s.Services {
		serviceMap[serviceName] = make(map[string]string)
	}

	return enaml.InstanceJob{
		Name:    "consul_agent",
		Release: "cf",
		Properties: &consullib.ConsulAgentJob{
			Consul: &consullib.Consul{
				EncryptKeys: s.Config.ConsulEncryptKeys,
				CaCert:      s.Config.BBSCACert,
				AgentCert:   s.Config.ConsulAgentCert,
				AgentKey:    s.Config.ConsulAgentKey,
				ServerCert:  s.Config.ConsulServerCert,
				ServerKey:   s.Config.ConsulServerKey,
				Agent: &consullib.Agent{
					Domain: "cf.internal",
					Mode:   s.getMode(),
					Servers: &consullib.Servers{
						Lan: s.Config.ConsulIPs,
					},
					Services: serviceMap,
				},
			},
		},
	}
}

func (s *ConsulAgent) getMode() interface{} {
	if s.Mode != "" {
		return s.Mode
	}
	return nil
}
