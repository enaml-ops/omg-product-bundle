package cloudfoundry

import (
	"github.com/enaml-ops/enaml"
	grtrlib "github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/gorouter"
	"github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/metron_agent"
)

type gorouter struct {
	Config *Config
}

//NewGoRouterPartition -
func NewGoRouterPartition(config *Config) InstanceGroupCreator {

	return &gorouter{
		Config: config,
	}
}

func (s *gorouter) ToInstanceGroup() (ig *enaml.InstanceGroup) {
	ig = &enaml.InstanceGroup{
		Name:      "router-partition",
		Instances: len(s.Config.RouterMachines),
		VMType:    s.Config.RouterVMType,
		AZs:       s.Config.AZs,
		Stemcell:  s.Config.StemcellName,
		Jobs: []enaml.InstanceJob{
			s.newRouterJob(),
			s.newMetronJob(),
			s.newStatsdInjectorJob(),
		},
		Networks: []enaml.Network{
			enaml.Network{Name: s.Config.NetworkName, StaticIPs: s.Config.RouterMachines},
		},
		Update: enaml.Update{
			MaxInFlight: 1,
		},
	}
	return
}

func (s *gorouter) newRouter() *grtrlib.Router {
	return &grtrlib.Router{
		EnableSsl:     s.Config.RouterEnableSSL,
		SecureCookies: false,
		SslKey:        s.Config.RouterSSLKey,
		SslCert:       s.Config.RouterSSLCert,
		Status: &grtrlib.Status{
			User:     s.Config.RouterUser,
			Password: s.Config.RouterPass,
		},
	}
}

func (s *gorouter) newStatsdInjectorJob() enaml.InstanceJob {
	return enaml.InstanceJob{
		Name:       "statsd-injector",
		Release:    "cf",
		Properties: make(map[interface{}]interface{}),
	}
}

func (s *gorouter) newRouterJob() enaml.InstanceJob {
	return enaml.InstanceJob{
		Name:    "gorouter",
		Release: "cf",
		Properties: &grtrlib.GorouterJob{
			RequestTimeoutInSeconds: 180,
			Nats: &grtrlib.Nats{
				User:     s.Config.NATSUser,
				Password: s.Config.NATSPassword,
				Machines: s.Config.NATSMachines,
				Port:     s.Config.NATSPort,
			},
			Router: s.newRouter(),
			Uaa: &grtrlib.Uaa{
				Ssl: &grtrlib.Ssl{
					Port: -1,
				},
				Clients: &grtrlib.Clients{
					Gorouter: &grtrlib.Gorouter{
						Secret: s.Config.GoRouterClientSecret,
					},
				},
			},
		},
	}
}

func (s *gorouter) newMetronJob() enaml.InstanceJob {
	return enaml.InstanceJob{
		Name:    "metron_agent",
		Release: "cf",
		Properties: &metron_agent.MetronAgentJob{
			SyslogDaemonConfig: &metron_agent.SyslogDaemonConfig{
				Transport: "tcp",
			},
			MetronAgent: &metron_agent.MetronAgent{
				Zone:       s.Config.MetronZone,
				Deployment: DeploymentName,
			},
			MetronEndpoint: &metron_agent.MetronEndpoint{
				SharedSecret: s.Config.MetronSecret,
			},
			Loggregator: &metron_agent.Loggregator{
				Etcd: &metron_agent.Etcd{
					Machines: s.Config.EtcdMachines,
				},
			},
		},
	}
}
