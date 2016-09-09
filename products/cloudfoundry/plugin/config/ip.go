package config

import "gopkg.in/urfave/cli.v2"

func RequiredIPSliceFlags() []string {
	return []string{
		"mysql-ip",
		"diego-cell-ip",
		"consul-ip",
		"doppler-ip",
		"mysql-proxy-ip",
		"loggregator-traffic-controller-ip",
		"nats-machine-ip",
		"etcd-machine-ip",
		"diego-brain-ip",
		"diego-db-ip",
		"router-ip",
	}
}
func RequiredIPFlags() []string {
	return []string{
		"nfs-ip",
	}
}

func NewIP(c *cli.Context) IP {
	return IP{
		MySQLIPs:        c.StringSlice("mysql-ip"),
		NFSIP:           c.String("nfs-ip"),
		DiegoCellIPs:    c.StringSlice("diego-cell-ip"),
		ConsulIPs:       c.StringSlice("consul-ip"),
		DopplerIPs:      c.StringSlice("doppler-ip"),
		MySQLProxyIPs:   c.StringSlice("mysql-proxy-ip"),
		LoggregratorIPs: c.StringSlice("loggregator-traffic-controller-ip"),
		NATSMachines:    c.StringSlice("nats-machine-ip"),
		HAProxyIPs:      c.StringSlice("haproxy-ip"),
		EtcdMachines:    c.StringSlice("etcd-machine-ip"),
		DiegoBrainIPs:   c.StringSlice("diego-brain-ip"),
		DiegoDBIPs:      c.StringSlice("diego-db-ip"),
		RouterMachines:  c.StringSlice("router-ip"),
	}
}

type IP struct {
	HAProxyIPs      []string
	NFSIP           string
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
}
