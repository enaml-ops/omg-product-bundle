package cloudfoundry

import (
	"github.com/codegangsta/cli"
	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-product-bundle/products/cf-mysql/enaml-gen/proxy"
	"github.com/xchapter7x/lo"
)

//MySQLProxy -
type MySQLProxy struct {
	Config           *Config
	VMTypeName       string
	NetworkIPs       []string
	ExternalHost     string
	APIUsername      string
	APIPassword      string
	ClusterIPs       []string
	SyslogAggregator *proxy.SyslogAggregator
}

//NewMySQLProxyPartition -
func NewMySQLProxyPartition(c *cli.Context, config *Config) InstanceGrouper {

	return &MySQLProxy{
		Config:       config,
		NetworkIPs:   c.StringSlice("mysql-proxy-ip"),
		VMTypeName:   c.String("mysql-proxy-vm-type"),
		APIUsername:  c.String("mysql-proxy-api-username"),
		APIPassword:  c.String("mysql-proxy-api-password"),
		ExternalHost: c.String("mysql-proxy-external-host"),
		ClusterIPs:   c.StringSlice("mysql-ip"),
		SyslogAggregator: &proxy.SyslogAggregator{
			Address:   c.String("syslog-address"),
			Port:      c.Int("syslog-port"),
			Transport: c.String("syslog-transport"),
		},
	}
}

//ToInstanceGroup -
func (s *MySQLProxy) ToInstanceGroup() (ig *enaml.InstanceGroup) {
	ig = &enaml.InstanceGroup{
		Name:      "mysql_proxy-partition",
		Instances: len(s.NetworkIPs),
		VMType:    s.VMTypeName,
		AZs:       s.Config.AZs,
		Stemcell:  s.Config.StemcellName,
		Jobs: []enaml.InstanceJob{
			s.newMySQLProxyJob(),
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

func (s *MySQLProxy) newMySQLProxyJob() enaml.InstanceJob {
	return enaml.InstanceJob{
		Name:    "proxy",
		Release: "cf-mysql",
		Properties: &proxy.ProxyJob{
			Proxy: &proxy.Proxy{
				ApiUsername: s.APIUsername,
				ApiPassword: s.APIPassword,
				ProxyIps:    s.NetworkIPs,
			},
			ExternalHost:     s.ExternalHost,
			ClusterIps:       s.ClusterIPs,
			SyslogAggregator: s.SyslogAggregator,
			Nats: &proxy.Nats{
				User:     s.Config.NATSUser,
				Password: s.Config.NATSPassword,
				Machines: s.Config.NATSMachines,
				Port:     s.Config.NATSPort,
			},
		},
	}
}

//HasValidValues -
func (s *MySQLProxy) HasValidValues() bool {

	lo.G.Debugf("checking '%s' for valid flags", "mysqlproxy")

	if len(s.NetworkIPs) <= 0 {
		lo.G.Debugf("could not find the correct number of network ips configured '%v' : '%v'", len(s.NetworkIPs), s.NetworkIPs)
	}
	if len(s.ClusterIPs) <= 0 {
		lo.G.Debugf("could not find the correct number of ClusterIPs configured '%v' : '%v'", len(s.ClusterIPs), s.ClusterIPs)
	}
	if s.VMTypeName == "" {
		lo.G.Debugf("could not find a valid vmtypename '%v'", s.VMTypeName)
	}
	if s.ExternalHost == "" {
		lo.G.Debugf("could not find a valid ExternalHost '%v'", s.ExternalHost)
	}
	if s.APIPassword == "" {
		lo.G.Debugf("could not find a valid APIPassword '%v'", s.APIPassword)
	}
	if s.APIUsername == "" {
		lo.G.Debugf("could not find a valid APIUsername '%v'", s.APIUsername)
	}
	return (s.VMTypeName != "" &&
		len(s.NetworkIPs) > 0 &&
		s.ExternalHost != "" &&
		s.APIPassword != "" &&
		s.APIUsername != "" &&
		len(s.ClusterIPs) > 0)
}
