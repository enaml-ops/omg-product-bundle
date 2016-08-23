package cloudfoundry

import (
	"github.com/codegangsta/cli"
	"github.com/enaml-ops/enaml"
	ltc "github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/loggregator_trafficcontroller"
	"github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/route_registrar"
	"github.com/xchapter7x/lo"
)

func NewLoggregatorTrafficController(c *cli.Context) InstanceGrouper {
	return &loggregatorTrafficController{
		AZs:               c.StringSlice("az"),
		StemcellName:      c.String("stemcell-name"),
		NetworkName:       c.String("network"),
		NetworkIPs:        c.StringSlice("loggregator-traffic-controller-ip"),
		VMTypeName:        c.String("loggregator-traffic-controller-vmtype"),
		SystemDomain:      c.String("system-domain"),
		DopplerSecret:     c.String("doppler-client-secret"),
		SkipSSLCertVerify: c.BoolT("skip-cert-verify"),
		EtcdMachines:      c.StringSlice("etcd-machine-ip"),
		Metron:            NewMetron(c),
		Nats: &route_registrar.Nats{
			User:     c.String("nats-user"),
			Password: c.String("nats-pass"),
			Machines: c.StringSlice("nats-machine-ip"),
			Port:     4222,
		},
	}
}

func (l *loggregatorTrafficController) HasValidValues() bool {

	lo.G.Debugf("checking '%s' for valid flags", "loggregator")

	if len(l.AZs) <= 0 {
		lo.G.Debugf("could not find the correct number of AZs configured '%v' : '%v'", len(l.AZs), l.AZs)
	}
	if len(l.NetworkIPs) <= 0 {
		lo.G.Debugf("could not find the correct number of network ips configured '%v' : '%v'", len(l.NetworkIPs), l.NetworkIPs)
	}
	if len(l.EtcdMachines) <= 0 {
		lo.G.Debugf("could not find the correct number of EtcdMachines configured '%v' : '%v'", len(l.EtcdMachines), l.EtcdMachines)
	}
	if l.StemcellName == "" {
		lo.G.Debugf("could not find a valid stemcellname '%v'", l.StemcellName)
	}
	if l.VMTypeName == "" {
		lo.G.Debugf("could not find a valid vmtypename '%v'", l.VMTypeName)
	}
	if l.NetworkName == "" {
		lo.G.Debugf("could not find a valid NetworkName '%v'", l.NetworkName)
	}
	if l.SystemDomain == "" {
		lo.G.Debugf("could not find a valid SystemDomain '%v'", l.SystemDomain)
	}
	if l.DopplerSecret == "" {
		lo.G.Debugf("could not find a valid DopplerSecret '%v'", l.DopplerSecret)
	}

	return len(l.AZs) > 0 &&
		l.StemcellName != "" &&
		l.NetworkName != "" &&
		len(l.NetworkIPs) > 0 &&
		l.VMTypeName != "" &&
		l.SystemDomain != "" &&
		len(l.EtcdMachines) > 0 &&
		l.DopplerSecret != "" &&
		l.Metron.HasValidValues()
}

func (l *loggregatorTrafficController) ToInstanceGroup() *enaml.InstanceGroup {
	return &enaml.InstanceGroup{
		Name:      "loggregator_trafficcontroller-partition",
		AZs:       l.AZs,
		Stemcell:  l.StemcellName,
		VMType:    l.VMTypeName,
		Instances: len(l.NetworkIPs),

		Networks: []enaml.Network{
			{
				Name:      l.NetworkName,
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
			SystemDomain: l.SystemDomain,
			Cc: &ltc.Cc{
				SrvApiUri: prefixSystemDomain(l.SystemDomain, "api"),
			},
			Ssl: &ltc.Ssl{
				SkipCertVerify: l.SkipSSLCertVerify,
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
						"uris":                  []string{"doppler." + l.SystemDomain},
					},
					map[string]interface{}{
						"name":                  "loggregator",
						"port":                  8080,
						"registration_interval": "20s",
						"uris":                  []string{"loggregator." + l.SystemDomain},
					},
				},
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
