package prabbitmq

import (
	"fmt"

	cli "gopkg.in/urfave/cli.v2"
)

// Config is used as input for generating instance groups.
type Config struct {
	DeploymentName  string
	Network         string
	StemcellVersion string
	ServerIPs       []string
	SyslogAddress   string
	SyslogPort      int
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
		DeploymentName:  getString("deployment-name"),
		Network:         getString("network"),
		StemcellVersion: getString("stemcell-ver"),
		SyslogAddress:   getString("syslog-address"),
		SyslogPort:      getInt("syslog-port"),
		ServerIPs:       getStringSlice("server-ip"),
	}

	var err error
	if len(missingFlags) > 0 {
		err = fmt.Errorf("prabbitmq: missing flags: %#v", missingFlags)
	}
	return cfg, err
}
