package config

import "github.com/codegangsta/cli"

func NewIP(c *cli.Context) IP {
	return IP{
		MySQLIPs:         c.StringSlice("mysql-ip"),
		NFSIPs:           c.StringSlice("nfs-ip"),
		NFSServerAddress: c.String("nfs-server-address"),
		DiegoCellIPs:     c.StringSlice("diego-cell-ip"),
		ConsulIPs:        c.StringSlice("consul-ip"),
		DopplerIPs:       c.StringSlice("doppler-ip"),
		MySQLProxyIPs:    c.StringSlice("mysql-proxy-ip"),
		LoggregratorIPs:  c.StringSlice("loggregator-traffic-controller-ip"),
		NATSMachines:     c.StringSlice("nats-machine-ip"),
		HAProxyIPs:       c.StringSlice("haproxy-ip"),
		EtcdMachines:     c.StringSlice("etcd-machine-ip"),
		DiegoBrainIPs:    c.StringSlice("diego-brain-ip"),
		DiegoDBIPs:       c.StringSlice("diego-db-ip"),
		RouterMachines:   c.StringSlice("router-ip"),
	}
}

type IP struct {
	HAProxyIPs      []string
	NFSIPs          []string
	MySQLIPs        []string
	LoggregratorIPs []string
	DopplerIPs      []string
	EtcdMachines    []string
	DiegoCellIPs    []string
	ConsulIPs       []string
	DiegoBrainIPs   []string
	MySQLProxyIPs   []string
	RouterMachines  []string
	DiegoDBIPs      []string
	NATSMachines    []string

	//Duplicate of NFSIPs
	NFSServerAddress string
}
