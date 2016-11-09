package pscs

import (
	"github.com/enaml-ops/pluginlib/pcli"
	"github.com/enaml-ops/pluginlib/pluginutil"

	cli "gopkg.in/urfave/cli.v2"
)

// Config is used as input for generating instance groups.
type Config struct {
	DeploymentName        string
	VMType                string   `omg:"vm-type"`
	AZs                   []string `omg:"az"`
	SystemDomain          string
	AppDomains            []string `omg:"app-domain"`
	Network               string
	StemcellVersion       string `omg:"stemcell-ver"`
	SkipSSLVerify         bool   `omg:"skip-ssl-verify"`
	BrokerUsername        string
	BrokerPassword        string
	WorkerClientSecret    string
	WorkerPassword        string
	InstancesPassword     string
	BrokerDashboardSecret string
	EncryptionKey         string
	CFAdminPassword       string `omg:"admin-password"`
	UAAAdminClientSecret  string `omg:"uaa-admin-secret"`
}

func configFromContext(c *cli.Context) (*Config, error) {
	cfg := &Config{}
	err := pcli.UnmarshalFlags(cfg, c)
	if err != nil {
		return nil, err
	}

	makePassword(&cfg.BrokerUsername)
	makePassword(&cfg.BrokerPassword)
	makePassword(&cfg.WorkerClientSecret)
	makePassword(&cfg.WorkerPassword)
	makePassword(&cfg.InstancesPassword)
	makePassword(&cfg.BrokerDashboardSecret)
	makePassword(&cfg.EncryptionKey)

	return cfg, nil
}

func makePassword(s *string) {
	if *s == generatePassword {
		*s = pluginutil.NewPassword(16)
	}
}
