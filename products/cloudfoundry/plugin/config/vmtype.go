package config

import "gopkg.in/urfave/cli.v2"

func RequiredVMTypeFlags() []string {
	return []string{
		"mysql-proxy-vm-type",
		"clock-global-vm-type",
		"cc-vm-type",
		"diego-brain-vm-type",
		"diego-cell-vm-type",
		"doppler-vm-type",
		"loggregator-traffic-controller-vmtype",
		"cc-worker-vm-type",
		"errand-vm-type",
		"etcd-vm-type",
		"nats-vm-type",
		"consul-vm-type",
		"mysql-vm-type",
		"diego-db-vm-type",
		"uaa-vm-type",
		"router-vm-type",
		"nfs-vm-type",
	}
}

func NewVMType(c *cli.Context) VMType {
	return VMType{
		HAProxyVMType:               c.String("haproxy-vm-type"),
		MySQLProxyVMType:            c.String("mysql-proxy-vm-type"),
		ClockGlobalVMType:           c.String("clock-global-vm-type"),
		CloudControllerVMType:       c.String("cc-vm-type"),
		DiegoBrainVMType:            c.String("diego-brain-vm-type"),
		DiegoCellVMType:             c.String("diego-cell-vm-type"),
		DopplerVMType:               c.String("doppler-vm-type"),
		LoggregratorVMType:          c.String("loggregator-traffic-controller-vmtype"),
		CloudControllerWorkerVMType: c.String("cc-worker-vm-type"),
		ErrandVMType:                c.String("errand-vm-type"),
		EtcdVMType:                  c.String("etcd-vm-type"),
		NatsVMType:                  c.String("nats-vm-type"),
		ConsulVMType:                c.String("consul-vm-type"),
		MySQLVMType:                 c.String("mysql-vm-type"),
		DiegoDBVMType:               c.String("diego-db-vm-type"),
		UAAVMType:                   c.String("uaa-vm-type"),
		RouterVMType:                c.String("router-vm-type"),
		NFSVMType:                   c.String("nfs-vm-type"),
	}
}

type VMType struct {
	NatsVMType                  string
	ConsulVMType                string
	MySQLVMType                 string
	EtcdVMType                  string
	ClockGlobalVMType           string
	MySQLProxyVMType            string
	HAProxyVMType               string
	RouterVMType                string
	NFSVMType                   string
	CloudControllerVMType       string
	CloudControllerWorkerVMType string
	DiegoDBVMType               string
	UAAVMType                   string
	DiegoCellVMType             string
	DiegoBrainVMType            string
	DopplerVMType               string
	ErrandVMType                string
	LoggregratorVMType          string
}
