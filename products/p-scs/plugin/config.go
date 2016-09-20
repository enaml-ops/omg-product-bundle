package pscs

import (
	"fmt"

	"github.com/enaml-ops/omg-cli/utils"

	cli "gopkg.in/urfave/cli.v2"
)

// Config is used as input for generating instance groups.
type Config struct {
	DeploymentName        string
	VMType                string
	AZs                   []string
	SystemDomain          string
	AppDomains            []string
	Network               string
	StemcellVersion       string
	SkipSSLVerify         bool
	BrokerUsername        string
	BrokerPassword        string
	WorkerClientSecret    string
	WorkerPassword        string
	InstancesPassword     string
	BrokerDashboardSecret string
	EncryptionKey         string
	CFAdminPassword       string
	UAAAdminClientSecret  string
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

	getStringSlice := func(flag string) []string {
		v := c.StringSlice(flag)
		if len(v) == 0 {
			missingFlags = append(missingFlags, flag)
		}
		return v
	}

	cfg := &Config{
		DeploymentName:        getString("deployment-name"),
		VMType:                getString("vm-type"),
		AZs:                   getStringSlice("az"),
		SystemDomain:          getString("system-domain"),
		AppDomains:            getStringSlice("app-domain"),
		Network:               getString("network"),
		StemcellVersion:       getString("stemcell-ver"),
		SkipSSLVerify:         c.Bool("skip-ssl-verify"),
		BrokerUsername:        getString("broker-username"),
		BrokerPassword:        getString("broker-password"),
		WorkerClientSecret:    getString("worker-client-secret"),
		WorkerPassword:        getString("worker-password"),
		InstancesPassword:     getString("instances-password"),
		BrokerDashboardSecret: getString("broker-dashboard-secret"),
		EncryptionKey:         getString("encryption-key"),
		CFAdminPassword:       getString("admin-password"),
		UAAAdminClientSecret:  getString("uaa-admin-secret"),
	}

	makePassword(&cfg.BrokerUsername)
	makePassword(&cfg.BrokerPassword)
	makePassword(&cfg.WorkerClientSecret)
	makePassword(&cfg.WorkerClientSecret)
	makePassword(&cfg.InstancesPassword)
	makePassword(&cfg.BrokerDashboardSecret)
	makePassword(&cfg.EncryptionKey)

	var err error
	if len(missingFlags) > 0 {
		err = fmt.Errorf("p-scs: missing flags: %#v", missingFlags)
	}
	return cfg, err
}

func makePassword(s *string) {
	if *s == generatePassword {
		*s = utils.NewPassword(16)
	}
}
