package cloudfoundry

import (
	"github.com/codegangsta/cli"
	"github.com/enaml-ops/enaml"
	grtrlib "github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/gorouter"
	"github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/metron_agent"
	"github.com/enaml-ops/pluginlib/util"
	"github.com/xchapter7x/lo"
)

const natsPort = 4222

//NewGoRouterPartition -
func NewGoRouterPartition(c *cli.Context) InstanceGrouper {
	cert, err := pluginutil.LoadResourceFromContext(c, "router-ssl-cert")
	if err != nil {
		lo.G.Fatalf("router cert: %s\n", err.Error())
	}
	key, err := pluginutil.LoadResourceFromContext(c, "router-ssl-key")
	if err != nil {
		lo.G.Fatalf("router key: %s\n", err.Error())
	}

	return &gorouter{
		Instances:    len(c.StringSlice("router-ip")),
		AZs:          c.StringSlice("az"),
		EnableSSL:    c.Bool("router-enable-ssl"),
		StemcellName: c.String("stemcell-name"),
		NetworkIPs:   c.StringSlice("router-ip"),
		NetworkName:  c.String("network"),
		VMTypeName:   c.String("router-vm-type"),
		SSLCert:      cert,
		SSLKey:       key,
		ClientSecret: c.String("gorouter-client-secret"),
		RouterUser:   c.String("router-user"),
		RouterPass:   c.String("router-pass"),
		MetronZone:   c.String("metron-zone"),
		MetronSecret: c.String("metron-secret"),
		Nats: grtrlib.Nats{
			User:     c.String("nats-user"),
			Password: c.String("nats-pass"),
			Machines: c.StringSlice("nats-machine-ip"),
		},
		Loggregator: metron_agent.Loggregator{
			Etcd: &metron_agent.Etcd{
				Machines: c.StringSlice("etcd-machine-ip"),
			},
		},
	}
}

func (s *gorouter) ToInstanceGroup() (ig *enaml.InstanceGroup) {
	ig = &enaml.InstanceGroup{
		Name:      "router-partition",
		Instances: s.Instances,
		VMType:    s.VMTypeName,
		AZs:       s.AZs,
		Stemcell:  s.StemcellName,
		Jobs: []enaml.InstanceJob{
			s.newRouterJob(),
			s.newMetronJob(),
			s.newStatsdInjectorJob(),
		},
		Networks: []enaml.Network{
			enaml.Network{Name: s.NetworkName, StaticIPs: s.NetworkIPs},
		},
		Update: enaml.Update{
			MaxInFlight: 1,
		},
	}
	return
}

func (s *gorouter) newNats() *grtrlib.Nats {
	s.Nats.Port = natsPort
	return &s.Nats
}

func (s *gorouter) newRouter() *grtrlib.Router {
	return &grtrlib.Router{
		EnableSsl:     s.EnableSSL,
		SecureCookies: false,
		SslKey:        s.SSLKey,
		SslCert:       s.SSLCert,
		Status: &grtrlib.Status{
			User:     s.RouterUser,
			Password: s.RouterPass,
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
			Nats:   s.newNats(),
			Router: s.newRouter(),
			Uaa: &grtrlib.Uaa{
				Ssl: &grtrlib.Ssl{
					Port: -1,
				},
				Clients: &grtrlib.Clients{
					Gorouter: &grtrlib.Gorouter{
						Secret: s.ClientSecret,
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
				Zone:       s.MetronZone,
				Deployment: DeploymentName,
			},
			MetronEndpoint: &metron_agent.MetronEndpoint{
				SharedSecret: s.MetronSecret,
			},
			Loggregator: &s.Loggregator,
		},
	}
}

func (s *gorouter) HasValidValues() bool {

	lo.G.Debugf("checking '%s' for valid flags", "gorouter")

	if len(s.AZs) <= 0 {
		lo.G.Debugf("could not find the correct number of AZs configured '%v' : '%v'", len(s.AZs), s.AZs)
	}
	if len(s.NetworkIPs) <= 0 {
		lo.G.Debugf("could not find the correct number of network ips configured '%v' : '%v'", len(s.NetworkIPs), s.NetworkIPs)
	}
	if s.StemcellName == "" {
		lo.G.Debugf("could not find a valid stemcellname '%v'", s.StemcellName)
	}
	if s.VMTypeName == "" {
		lo.G.Debugf("could not find a valid vmtypename '%v'", s.VMTypeName)
	}
	if s.MetronZone == "" {
		lo.G.Debugf("could not find a valid MetronZone '%v'", s.MetronZone)
	}
	if s.MetronSecret == "" {
		lo.G.Debugf("could not find a valid MetronSecret '%v'", s.MetronSecret)
	}
	if s.RouterPass == "" {
		lo.G.Debugf("could not find a valid RouterPass '%v'", s.RouterPass)
	}
	if s.RouterUser == "" {
		lo.G.Debugf("could not find a valid RouterUser '%v'", s.RouterUser)
	}
	if s.NetworkName == "" {
		lo.G.Debugf("could not find a valid NetworkName '%v'", s.NetworkName)
	}
	if s.SSLCert == "" {
		lo.G.Debugf("could not find a valid SSLCert '%v'", s.SSLCert)
	}
	if s.SSLKey == "" {
		lo.G.Debugf("could not find a valid SSLKey '%v'", s.SSLKey)
	}
	if s.Nats.User == "" {
		lo.G.Debugf("could not find a valid Nats.User '%v'", s.Nats.User)
	}
	if s.Nats.Password == "" {
		lo.G.Debugf("could not find a valid Nats.Password '%v'", s.Nats.Password)
	}
	if s.Nats.Machines != nil {
		lo.G.Debugf("could not find a valid Nats.Machines '%v'", s.Nats.Machines)
	}
	if s.Loggregator.Etcd.Machines != nil {
		lo.G.Debugf("could not find a valid Loggregator.Etcd.Machines '%v'", s.Loggregator.Etcd.Machines)
	}
	return (len(s.AZs) > 0 &&
		s.StemcellName != "" &&
		s.VMTypeName != "" &&
		s.MetronZone != "" &&
		s.MetronSecret != "" &&
		s.RouterPass != "" &&
		s.RouterUser != "" &&
		s.NetworkName != "" &&
		len(s.NetworkIPs) > 0 &&
		s.SSLCert != "" &&
		s.SSLKey != "" &&
		s.Nats.User != "" &&
		s.Nats.Password != "" &&
		s.Nats.Machines != nil &&
		s.Loggregator.Etcd.Machines != nil)
}
