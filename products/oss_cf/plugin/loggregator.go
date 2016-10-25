package cloudfoundry

import (
	"github.com/enaml-ops/enaml"
	ltc "github.com/enaml-ops/omg-product-bundle/products/oss_cf/enaml-gen/loggregator_trafficcontroller"
	"github.com/enaml-ops/omg-product-bundle/products/oss_cf/enaml-gen/route_registrar"
	"github.com/enaml-ops/omg-product-bundle/products/oss_cf/plugin/config"
)

type loggregatorTrafficController struct {
	Config *config.Config
	Metron *Metron
}

func NewLoggregatorTrafficController(config *config.Config) InstanceGroupCreator {
	return &loggregatorTrafficController{
		Config: config,
		Metron: NewMetron(config),
	}
}

func (l *loggregatorTrafficController) ToInstanceGroup() *enaml.InstanceGroup {
	return &enaml.InstanceGroup{
		Name:      "loggregator_trafficcontroller-partition",
		AZs:       l.Config.AZs,
		Stemcell:  l.Config.StemcellName,
		VMType:    l.Config.LoggregratorVMType,
		Instances: len(l.Config.LoggregratorIPs),

		Networks: []enaml.Network{
			{
				Name:      l.Config.NetworkName,
				StaticIPs: l.Config.LoggregratorIPs,
			},
		},
		Update: enaml.Update{
			MaxInFlight: 1,
		},
		Jobs: []enaml.InstanceJob{
			l.createLoggregatorTrafficControllerJob(),
			l.Metron.CreateJob(),
			l.createRouteRegistrarJob(),
			l.createStatsdInjectorJob(),
		},
	}
}

func (l *loggregatorTrafficController) createLoggregatorTrafficControllerJob() enaml.InstanceJob {
	return enaml.InstanceJob{
		Name:    "loggregator_trafficcontroller",
		Release: CFReleaseName,
		Properties: &ltc.LoggregatorTrafficcontrollerJob{
			SystemDomain: l.Config.SystemDomain,
			Cc: &ltc.Cc{
				SrvApiUri: prefixSystemDomain(l.Config.SystemDomain, "api"),
			},
			Ssl: &ltc.Ssl{
				SkipCertVerify: l.Config.SkipSSLCertVerify,
			},
			TrafficController: &ltc.TrafficController{},
			Doppler:           &ltc.Doppler{},
			Loggregator: &ltc.Loggregator{
				Etcd: &ltc.LoggregatorEtcd{
					Machines: l.Config.EtcdMachines,
				},
			},
			Uaa: &ltc.Uaa{
				Clients: &ltc.Clients{
					Doppler: &ltc.ClientsDoppler{
						Secret: l.Config.DopplerSecret,
					},
				},
			},
		},
	}
}

func (l *loggregatorTrafficController) createRouteRegistrarJob() enaml.InstanceJob {
	return enaml.InstanceJob{
		Name:    "route_registrar",
		Release: CFReleaseName,
		Properties: &route_registrar.RouteRegistrarJob{
			RouteRegistrar: &route_registrar.RouteRegistrar{
				Routes: []map[string]interface{}{
					map[string]interface{}{
						"name":                  "doppler",
						"port":                  8081,
						"registration_interval": "20s",
						"uris":                  []string{"doppler." + l.Config.SystemDomain},
					},
					map[string]interface{}{
						"name":                  "loggregator",
						"port":                  8080,
						"registration_interval": "20s",
						"uris":                  []string{"loggregator." + l.Config.SystemDomain},
					},
				},
			},
			Nats: &route_registrar.Nats{
				User:     l.Config.NATSUser,
				Password: l.Config.NATSPassword,
				Port:     l.Config.NATSPort,
				Machines: l.Config.NATSMachines,
			},
		},
	}
}

func (l *loggregatorTrafficController) createStatsdInjectorJob() enaml.InstanceJob {
	return enaml.InstanceJob{
		Name:       "statsd-injector",
		Release:    CFReleaseName,
		Properties: struct{}{},
	}
}
