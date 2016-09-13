package prabbitmq

import (
	"fmt"

	"github.com/enaml-ops/omg-cli/utils"

	cli "gopkg.in/urfave/cli.v2"
)

// Config is used as input for generating instance groups.
type Config struct {
	DeploymentName            string
	SystemDomain              string
	ServiceURL                string
	ServiceAdminPassword      string
	PublicIP                  string
	Network                   string
	StemcellVersion           string
	ServerIPs                 []string
	BrokerIP                  string
	BrokerPassword            string
	SyslogAddress             string
	SyslogPort                int
	NATSMachines              []string
	NATSPort                  int
	NATSPassword              string
	HAProxyStatsAdminPassword string
	SystemServicesPassword    string
	SkipSSLVerify             bool
	MetronZone                string
	MetronSecret              string
	EtcdMachines              []string
}

func configFromContext(c *cli.Context) (*Config, error) {
	var missingFlags []string

	getString := func(flag string) string {
		v := c.String(flag)
		if v == "" {
			missingFlags = append(missingFlags, flag)
		}
		return v
	}
	getInt := func(flag string) int {
		v := c.Int(flag)
		if v == 0 { // TODO is this okay?
			missingFlags = append(missingFlags, flag)
		}
		return v
	}
	getStringSlice := func(flag string) []string {
		v := c.StringSlice(flag)
		if len(v) == 0 {
			missingFlags = append(missingFlags, flag)
		}
		return v
	}

	cfg := &Config{
		DeploymentName:            getString("deployment-name"),
		ServiceURL:                getString("service-url"),
		ServiceAdminPassword:      getString("service-admin-password"),
		SystemDomain:              getString("system-domain"),
		Network:                   getString("network"),
		StemcellVersion:           getString("stemcell-ver"),
		ServerIPs:                 getStringSlice("server-ip"),
		BrokerIP:                  getString("broker-ip"),
		BrokerPassword:            getString("broker-password"),
		SyslogAddress:             getString("syslog-address"),
		SyslogPort:                getInt("syslog-port"),
		NATSMachines:              getStringSlice("nats-ip"),
		NATSPort:                  getInt("nats-port"),
		NATSPassword:              getString("nats-password"),
		HAProxyStatsAdminPassword: getString("haproxy-stats-password"),
		SystemServicesPassword:    getString("system-services-password"),
		SkipSSLVerify:             c.Bool("skip-ssl-verify"),
		MetronZone:                getString("doppler-zone"),
		MetronSecret:              getString("doppler-shared-secret"),
		EtcdMachines:              getStringSlice("etcd-machine"),
	}

	makePassword(&cfg.ServiceAdminPassword)
	makePassword(&cfg.BrokerPassword)
	makePassword(&cfg.NATSPassword)
	makePassword(&cfg.HAProxyStatsAdminPassword)

	var err error
	if len(missingFlags) > 0 {
		err = fmt.Errorf("prabbitmq: missing flags: %#v", missingFlags)
	}
	return cfg, err
}

func makePassword(s *string) {
	if *s == generatePassword {
		*s = utils.NewPassword(16)
	}
}
