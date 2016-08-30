package cloudfoundry

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/xchapter7x/lo"
)

type Config struct {
	AZs          []string
	StemcellName string
	NetworkName  string
	SystemDomain string
	AppDomains   []string

	MySQLIPs               []string
	MySQLBootstrapUser     string
	MySQLBootstrapPassword string

	AllowSSHAccess           bool
	SkipSSLCertVerify        bool
	NATSUser                 string
	NATSPassword             string
	NATSPort                 int
	NATSMachines             []string
	CCDBUser                 string
	CCDBPassword             string
	JWTVerificationKey       string
	CCServiceDashboardSecret string
}

func NewConfig(c *cli.Context) (*Config, error) {

	//check for required fields
	if err := checkRequired(c); err != nil {
		return nil, err
	}
	config := &Config{
		AZs:                    c.StringSlice("az"),
		StemcellName:           c.String("stemcell-name"),
		NetworkName:            c.String("network"),
		SystemDomain:           c.String("system-domain"),
		AppDomains:             c.StringSlice("app-domain"),
		NATSUser:               c.String("nats-user"),
		NATSPassword:           c.String("nats-pass"),
		NATSPort:               c.Int("nats-port"),
		NATSMachines:           c.StringSlice("nats-machine-ip"),
		MySQLIPs:               c.StringSlice("mysql-ip"),
		MySQLBootstrapUser:     c.String("mysql-bootstrap-username"),
		MySQLBootstrapPassword: c.String("mysql-bootstrap-password"),

		//boolean flags
		AllowSSHAccess:    c.Bool("allow-app-ssh-access"),
		SkipSSLCertVerify: c.BoolT("skip-cert-verify"),
	}
	return config, nil
}
func checkRequired(c *cli.Context) error {
	invalidNames := []string{}
	invalidNames = append(invalidNames, checkRequiredStringFlags(c)...)
	invalidNames = append(invalidNames, checkRequiredStringSliceFlags(c)...)
	if len(invalidNames) > 0 {
		return fmt.Errorf("Sorry you need to provide %v flags to continue", invalidNames)
	}
	return nil
}

func checkRequiredStringFlags(c *cli.Context) []string {
	requiredFlags := []string{"stemcell-name", "network", "system-domain", "nats-user", "nats-pass", "nats-port", "mysql-bootstrap-username", "mysql-bootstrap-password"}
	invalidNames := []string{}
	for _, name := range requiredFlags {
		if c.String(name) == "" {
			invalidNames = append(invalidNames, name)
		} else {
			lo.G.Debug(name, "==>", c.String(name))
		}
	}
	return invalidNames
}

func checkRequiredStringSliceFlags(c *cli.Context) []string {
	requiredFlags := []string{"az", "app-domain", "nats-machine-ip", "mysql-ip"}
	invalidNames := []string{}
	for _, name := range requiredFlags {
		if len(c.StringSlice(name)) == 0 {
			invalidNames = append(invalidNames, name)
		} else {
			lo.G.Debug(name, "==>", c.StringSlice(name))
		}
	}
	return invalidNames
}
