package cloudfoundry

import (
	"github.com/codegangsta/cli"
	"github.com/enaml-ops/enaml"
	grtrlib "github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/gorouter"
	"github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/metron_agent"
	"github.com/enaml-ops/pluginlib/util"
	"github.com/xchapter7x/lo"
)

type gorouter struct {
	Config       *Config
	VMTypeName   string
	NetworkIPs   []string
	SSLCert      string
	SSLKey       string
	EnableSSL    bool
	ClientSecret string
	Loggregator  metron_agent.Loggregator
	RouterUser   string
	RouterPass   string
	MetronZone   string
	MetronSecret string
}

//NewGoRouterPartition -
func NewGoRouterPartition(c *cli.Context, config *Config) InstanceGrouper {
	cert, err := pluginutil.LoadResourceFromContext(c, "router-ssl-cert")
	if err != nil {
		lo.G.Fatalf("router cert: %s\n", err.Error())
	}
	key, err := pluginutil.LoadResourceFromContext(c, "router-ssl-key")
	if err != nil {
		lo.G.Fatalf("router key: %s\n", err.Error())
	}

	return &gorouter{
		Config:       config,
		EnableSSL:    c.Bool("router-enable-ssl"),
		NetworkIPs:   c.StringSlice("router-ip"),
		VMTypeName:   c.String("router-vm-type"),
		SSLCert:      cert,
		SSLKey:       key,
		ClientSecret: c.String("gorouter-client-secret"),
		RouterUser:   c.String("router-user"),
		RouterPass:   c.String("router-pass"),
		MetronZone:   c.String("metron-zone"),
		MetronSecret: c.String("metron-secret"),
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
		Instances: len(s.NetworkIPs),
		VMType:    s.VMTypeName,
		AZs:       s.Config.AZs,
		Stemcell:  s.Config.StemcellName,
		Jobs: []enaml.InstanceJob{
			s.newRouterJob(),
			s.newMetronJob(),
			s.newStatsdInjectorJob(),
		},
		Networks: []enaml.Network{
			enaml.Network{Name: s.Config.NetworkName, StaticIPs: s.NetworkIPs},
		},
		Update: enaml.Update{
			MaxInFlight: 1,
		},
	}
	return
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

	if len(s.NetworkIPs) <= 0 {
		lo.G.Debugf("could not find the correct number of network ips configured '%v' : '%v'", len(s.NetworkIPs), s.NetworkIPs)
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
	if s.SSLCert == "" {
		lo.G.Debugf("could not find a valid SSLCert '%v'", s.SSLCert)
	}
	if s.SSLKey == "" {
		lo.G.Debugf("could not find a valid SSLKey '%v'", s.SSLKey)
	}
	if s.Loggregator.Etcd.Machines != nil {
		lo.G.Debugf("could not find a valid Loggregator.Etcd.Machines '%v'", s.Loggregator.Etcd.Machines)
	}
	return (s.VMTypeName != "" &&
		s.MetronZone != "" &&
		s.MetronSecret != "" &&
		s.RouterPass != "" &&
		s.RouterUser != "" &&
		len(s.NetworkIPs) > 0 &&
		s.SSLCert != "" &&
		s.SSLKey != "" &&
		s.Loggregator.Etcd.Machines != nil)
}
