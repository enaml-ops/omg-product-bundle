package cloudfoundry

import (
	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-product-bundle/products/oss_cf/enaml-gen/metron_agent"
	"github.com/enaml-ops/omg-product-bundle/products/oss_cf/plugin/config"
)

//Metron -
type Metron struct {
	Config *config.Config
}

//NewMetron -
func NewMetron(config *config.Config) (metron *Metron) {
	metron = &Metron{
		Config: config,
	}
	return
}

//CreateJob -
func (s *Metron) CreateJob() enaml.InstanceJob {
	return enaml.InstanceJob{
		Name:    "metron_agent",
		Release: "cf",
		Properties: &metron_agent.MetronAgentJob{
			SyslogDaemonConfig: &metron_agent.SyslogDaemonConfig{
				Transport: s.Config.SyslogTransport,
				Address:   s.Config.SyslogAddress,
				Port:      s.Config.SyslogPort,
			},
			MetronAgent: &metron_agent.MetronAgent{
				Zone:       s.Config.DopplerZone,
				Deployment: DeploymentName,
			},
			MetronEndpoint: &metron_agent.MetronEndpoint{
				SharedSecret: s.Config.DopplerSharedSecret,
			},
			Loggregator: &metron_agent.Loggregator{
				Etcd: &metron_agent.LoggregatorEtcd{
					Machines: s.Config.EtcdMachines,
				},
			},
		},
	}
}
