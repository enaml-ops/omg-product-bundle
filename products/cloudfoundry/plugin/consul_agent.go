package cloudfoundry

import (
	"github.com/codegangsta/cli"
	"github.com/enaml-ops/enaml"
	consullib "github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/consul_agent"
	"github.com/xchapter7x/lo"
)

//ConsulAgent -
type ConsulAgent struct {
	Config     *Config
	NetworkIPs []string
	Mode       string
	Services   []string
}

//NewConsulAgent -
func NewConsulAgent(c *cli.Context, services []string, config *Config) *ConsulAgent {
	ca := &ConsulAgent{
		Config:     config,
		NetworkIPs: c.StringSlice("consul-ip"),
		Services:   services,
	}
	return ca
}

//NewConsulAgentServer -
func NewConsulAgentServer(c *cli.Context, config *Config) *ConsulAgent {
	ca := &ConsulAgent{
		Config:     config,
		NetworkIPs: c.StringSlice("consul-ip"),
		Mode:       "server",
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
				CaCert:      s.Config.ConsulCaCert,
				AgentCert:   s.Config.ConsulAgentCert,
				AgentKey:    s.Config.ConsulAgentKey,
				ServerCert:  s.Config.ConsulServerCert,
				ServerKey:   s.Config.ConsulServerKey,
				Agent: &consullib.Agent{
					Domain: "cf.internal",
					Mode:   s.getMode(),
					Servers: &consullib.Servers{
						Lan: s.NetworkIPs,
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

//HasValidValues -
func (s *ConsulAgent) HasValidValues() bool {

	lo.G.Debugf("checking '%s' for valid flags", "consul agent")

	if len(s.NetworkIPs) <= 0 {
		lo.G.Debugf("could not find the correct number of networkips configured '%v' : '%v'", len(s.NetworkIPs), s.NetworkIPs)
	}

	return len(s.NetworkIPs) > 0
}
