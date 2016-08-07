package cloudfoundry

import (
	"github.com/codegangsta/cli"
	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/metron_agent"
	"github.com/xchapter7x/lo"
)

//NewMetron -
func NewMetron(c *cli.Context) (metron *Metron) {
	metron = &Metron{
		Zone:            c.String("metron-zone"),
		Secret:          c.String("metron-secret"),
		SyslogAddress:   c.String("syslog-address"),
		SyslogPort:      c.Int("syslog-port"),
		SyslogTransport: c.String("syslog-transport"),
		Loggregator: metron_agent.Loggregator{
			Etcd: &metron_agent.Etcd{
				Machines: c.StringSlice("etcd-machine-ip"),
			},
		},
	}
	if metron.SyslogTransport == "" {
		metron.SyslogTransport = "tcp"
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
				Transport: s.SyslogTransport,
				Address:   s.SyslogAddress,
				Port:      s.SyslogPort,
			},
			MetronAgent: &metron_agent.MetronAgent{
				Zone:       s.Zone,
				Deployment: DeploymentName,
			},
			MetronEndpoint: &metron_agent.MetronEndpoint{
				SharedSecret: s.Secret,
			},
			Loggregator: &s.Loggregator,
		},
	}
}

//HasValidValues -
func (s *Metron) HasValidValues() bool {
	lo.G.Debugf("checking '%s' for valid flags", "metron")
	if s.Zone == "" {
		lo.G.Debugf("could not find a valid Zone '%v'", s.Zone)
	}
	if s.Secret == "" {
		lo.G.Debugf("could not find a valid Secret '%v'", s.Secret)
	}
	return (s.Zone != "" && s.Secret != "")
}
