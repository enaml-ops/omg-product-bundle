package cloudfoundry

import (
	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/haproxy"
)

// HAProxy -
type HAProxy struct {
	Config         *Config
	ConsulAgent    *ConsulAgent
	Metron         *Metron
	StatsdInjector *StatsdInjector
}

//NewHaProxyPartition -
func NewHaProxyPartition(config *Config) InstanceGroupCreator {
	return &HAProxy{
		Config:         config,
		ConsulAgent:    NewConsulAgent([]string{}, config),
		Metron:         NewMetron(config),
		StatsdInjector: NewStatsdInjector(nil),
	}
}

//ToInstanceGroup -
func (s *HAProxy) ToInstanceGroup() (ig *enaml.InstanceGroup) {
	if !s.Config.HAProxySkip {
		ig = &enaml.InstanceGroup{
			Name:      "ha_proxy-partition",
			Instances: len(s.Config.HAProxyIPs),
			VMType:    s.Config.HAProxyVMType,
			AZs:       s.Config.AZs,
			Stemcell:  s.Config.StemcellName,
			Jobs: []enaml.InstanceJob{
				s.createHAProxyJob(),
				s.ConsulAgent.CreateJob(),
				s.Metron.CreateJob(),
				s.StatsdInjector.CreateJob(),
			},
			Networks: []enaml.Network{
				enaml.Network{Name: s.Config.NetworkName, StaticIPs: s.Config.HAProxyIPs},
			},
			Update: enaml.Update{
				MaxInFlight: 1,
			},
		}
	}
	return
}

func (s *HAProxy) createHAProxyJob() enaml.InstanceJob {
	return enaml.InstanceJob{
		Name:    "haproxy",
		Release: "cf",
		Properties: &haproxy.HaproxyJob{
			RequestTimeoutInSeconds: 180,
			HaProxy: &haproxy.HaProxy{
				DisableHttp: true,
				SslPem:      s.Config.HAProxySSLPem,
			},
			Router: &haproxy.Router{
				Servers: &haproxy.Servers{
					Z1: s.Config.RouterMachines,
				},
			},
			Cc: &haproxy.Cc{
				AllowAppSshAccess: true,
			},
		},
	}
}
