package config

// IP contains IP configuration for a Cloud Foundry deployment.
type IP struct {
	HAProxyIPs      []string `omg:"haproxy-ip,optional"`
	NFSIP           string   `omg:"nfs-ip"`
	MySQLIPs        []string `omg:"mysql-ip"`
	LoggregratorIPs []string `omg:"loggregator-traffic-controller-ip"`
	DopplerIPs      []string `omg:"doppler-ip"`
	EtcdMachines    []string `omg:"etcd-machine-ip"`
	DiegoCellIPs    []string `omg:"diego-cell-ip"`
	ConsulIPs       []string `omg:"consul-ip"`
	DiegoBrainIPs   []string `omg:"diego-brain-ip"`
	MySQLProxyIPs   []string `omg:"mysql-proxy-ip"`
	RouterMachines  []string `omg:"router-ip"`
	DiegoDBIPs      []string `omg:"diego-db-ip"`
	NATSMachines    []string `omg:"nats-machine-ip"`
}
