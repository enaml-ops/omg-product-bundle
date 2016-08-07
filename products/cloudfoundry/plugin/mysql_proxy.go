package cloudfoundry

import (
	"github.com/codegangsta/cli"
	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-product-bundle/products/cf-mysql/enaml-gen/proxy"
	"github.com/xchapter7x/lo"
)

//NewMySQLProxyPartition -
func NewMySQLProxyPartition(c *cli.Context) InstanceGrouper {

	return &MySQLProxy{
		AZs:          c.StringSlice("az"),
		StemcellName: c.String("stemcell-name"),
		NetworkIPs:   c.StringSlice("mysql-proxy-ip"),
		NetworkName:  c.String("network"),
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
		Nats: &proxy.Nats{
			User:     c.String("nats-user"),
			Password: c.String("nats-pass"),
			Machines: c.StringSlice("nats-machine-ip"),
			Port:     4222,
		},
	}
}

//ToInstanceGroup -
func (s *MySQLProxy) ToInstanceGroup() (ig *enaml.InstanceGroup) {
	ig = &enaml.InstanceGroup{
		Name:      "mysql_proxy-partition",
		Instances: len(s.NetworkIPs),
		VMType:    s.VMTypeName,
		AZs:       s.AZs,
		Stemcell:  s.StemcellName,
		Jobs: []enaml.InstanceJob{
			s.newMySQLProxyJob(),
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

func (s *MySQLProxy) newMySQLProxyJob() enaml.InstanceJob {
	return enaml.InstanceJob{
		Name:    "proxy",
		Release: "cf-mysql",
		Properties: &proxy.Proxy{
			ApiUsername:      s.APIUsername,
			ApiPassword:      s.APIPassword,
			ExternalHost:     s.ExternalHost,
			ClusterIps:       s.ClusterIPs,
			SyslogAggregator: s.SyslogAggregator,
			Nats:             s.Nats,
		},
	}
}

//HasValidValues -
func (s *MySQLProxy) HasValidValues() bool {

	lo.G.Debugf("checking '%s' for valid flags", "mysqlproxy")

	if len(s.AZs) <= 0 {
		lo.G.Debugf("could not find the correct number of AZs configured '%v' : '%v'", len(s.AZs), s.AZs)
	}
	if len(s.NetworkIPs) <= 0 {
		lo.G.Debugf("could not find the correct number of network ips configured '%v' : '%v'", len(s.NetworkIPs), s.NetworkIPs)
	}
	if len(s.ClusterIPs) <= 0 {
		lo.G.Debugf("could not find the correct number of ClusterIPs configured '%v' : '%v'", len(s.ClusterIPs), s.ClusterIPs)
	}
	if s.StemcellName == "" {
		lo.G.Debugf("could not find a valid stemcellname '%v'", s.StemcellName)
	}
	if s.VMTypeName == "" {
		lo.G.Debugf("could not find a valid vmtypename '%v'", s.VMTypeName)
	}
	if s.NetworkName == "" {
		lo.G.Debugf("could not find a valid NetworkName '%v'", s.NetworkName)
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
	return (len(s.AZs) > 0 &&
		s.StemcellName != "" &&
		s.VMTypeName != "" &&
		s.NetworkName != "" &&
		len(s.NetworkIPs) > 0 &&
		s.ExternalHost != "" &&
		s.APIPassword != "" &&
		s.APIUsername != "" &&
		len(s.ClusterIPs) > 0)
}
