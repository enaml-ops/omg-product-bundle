package prabbitmq

import (
	"github.com/enaml-ops/pluginlib/pcli"
	"github.com/enaml-ops/pluginlib/pluginutil"

	cli "gopkg.in/urfave/cli.v2"
)

// Config is used as input for generating instance groups.
type Config struct {
	DeploymentName            string   `omg:"deployment-name"`
	AZs                       []string `omg:"az"`
	AdminPassword             string   `omg:"rabbit-admin-password"`
	SystemDomain              string
	ServiceAdminPassword      string
	PublicIP                  string `omg:"rabbit-public-ip"`
	Network                   string
	StemcellVersion           string   `omg:"stemcell-ver"`
	ServerIPs                 []string `omg:"rabbit-server-ip"`
	BrokerIP                  string   `omg:"rabbit-broker-ip"`
	BrokerPassword            string
	SyslogAddress             string
	SyslogPort                int
	NATSMachines              []string `omg:"nats-machine-ip"`
	NATSPort                  int      `omg:"nats-port"`
	NATSPassword              string   `omg:"nats-pass"`
	HAProxyStatsAdminPassword string   `omg:"haproxy-stats-password"`
	SystemServicesPassword    string
	SkipSSLVerify             bool     `omg:"skip-ssl-verify"`
	MetronZone                string   `omg:"doppler-zone"`
	MetronSecret              string   `omg:"doppler-shared-secret"`
	EtcdMachines              []string `omg:"etcd-machine-ip"`
	BrokerVMType              string   `omg:"rabbit-broker-vm-type"`
	ServerVMType              string   `omg:"rabbit-server-vm-type"`
	HAProxyVMType             string   `omg:"rabbit-haproxy-vm-type"`
}

func configFromContext(c *cli.Context) (*Config, error) {

	cfg := &Config{}
	err := pcli.UnmarshalFlags(cfg, c)
	if err != nil {
		return nil, err
	}

	makePassword(&cfg.AdminPassword)
	makePassword(&cfg.ServiceAdminPassword)
	makePassword(&cfg.BrokerPassword)
	makePassword(&cfg.NATSPassword)
	makePassword(&cfg.HAProxyStatsAdminPassword)

	return cfg, nil
}

func makePassword(s *string) {
	if *s == generatePassword {
		*s = pluginutil.NewPassword(16)
	}
}
