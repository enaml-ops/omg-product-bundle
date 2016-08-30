package cloudfoundry

import (
	"github.com/codegangsta/cli"
	"github.com/enaml-ops/enaml"
	ltc "github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/loggregator_trafficcontroller"
	"github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/route_registrar"
	"github.com/xchapter7x/lo"
)

type loggregatorTrafficController struct {
	Config        *Config
	VMTypeName    string
	NetworkIPs    []string
	EtcdMachines  []string
	DopplerSecret string
	Metron        *Metron
}

func NewLoggregatorTrafficController(c *cli.Context, config *Config) InstanceGrouper {
	return &loggregatorTrafficController{
		Config:        config,
		NetworkIPs:    c.StringSlice("loggregator-traffic-controller-ip"),
		VMTypeName:    c.String("loggregator-traffic-controller-vmtype"),
		DopplerSecret: c.String("doppler-client-secret"),
		EtcdMachines:  c.StringSlice("etcd-machine-ip"),
		Metron:        NewMetron(c),
	}
}

func (l *loggregatorTrafficController) HasValidValues() bool {

	lo.G.Debugf("checking '%s' for valid flags", "loggregator")

	if len(l.NetworkIPs) <= 0 {
		lo.G.Debugf("could not find the correct number of network ips configured '%v' : '%v'", len(l.NetworkIPs), l.NetworkIPs)
	}
	if len(l.EtcdMachines) <= 0 {
		lo.G.Debugf("could not find the correct number of EtcdMachines configured '%v' : '%v'", len(l.EtcdMachines), l.EtcdMachines)
	}
	if l.VMTypeName == "" {
		lo.G.Debugf("could not find a valid vmtypename '%v'", l.VMTypeName)
	}
	if l.DopplerSecret == "" {
		lo.G.Debugf("could not find a valid DopplerSecret '%v'", l.DopplerSecret)
	}

	return len(l.NetworkIPs) > 0 &&
		l.VMTypeName != "" &&
		len(l.EtcdMachines) > 0 &&
		l.DopplerSecret != "" &&
		l.Metron.HasValidValues()
}

func (l *loggregatorTrafficController) ToInstanceGroup() *enaml.InstanceGroup {
	return &enaml.InstanceGroup{
		Name:      "loggregator_trafficcontroller-partition",
		AZs:       l.Config.AZs,
		Stemcell:  l.Config.StemcellName,
		VMType:    l.VMTypeName,
		Instances: len(l.NetworkIPs),

		Networks: []enaml.Network{
			{
				Name:      l.Config.NetworkName,
				StaticIPs: l.NetworkIPs,
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
			TrafficController: &ltc.TrafficController{
				Zone: l.Metron.Zone,
			},
			Doppler: &ltc.Doppler{
				UaaClientId: "doppler",
			},
			Loggregator: &ltc.Loggregator{
				Etcd: &ltc.Etcd{
					Machines: l.EtcdMachines,
				},
			},
			Uaa: &ltc.Uaa{
				Clients: &ltc.Clients{
					Doppler: &ltc.ClientsDoppler{
						Secret: l.DopplerSecret,
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
