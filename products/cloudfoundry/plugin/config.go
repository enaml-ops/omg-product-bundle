package cloudfoundry

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/enaml-ops/pluginlib/util"
	"github.com/xchapter7x/lo"
)

type Config struct {
	AZs           []string
	StemcellName  string
	NetworkName   string
	SystemDomain  string
	AppDomains    []string
	AdminPassword string

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
	ConsulCaCert             string
	ConsulEncryptKeys        []string
	ConsulAgentCert          string
	ConsulAgentKey           string
	ConsulServerCert         string
	ConsulServerKey          string

	ErrandVMType       string
	SmokeTestsPassword string
	UAALoginProtocol   string
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
		AdminPassword:          c.String("admin-password"),
		NATSUser:               c.String("nats-user"),
		NATSPassword:           c.String("nats-pass"),
		NATSPort:               c.Int("nats-port"),
		NATSMachines:           c.StringSlice("nats-machine-ip"),
		MySQLIPs:               c.StringSlice("mysql-ip"),
		MySQLBootstrapUser:     c.String("mysql-bootstrap-username"),
		MySQLBootstrapPassword: c.String("mysql-bootstrap-password"),
		ConsulEncryptKeys:      c.StringSlice("consul-encryption-key"),

		//boolean flags
		AllowSSHAccess:    c.Bool("allow-app-ssh-access"),
		SkipSSLCertVerify: c.BoolT("skip-cert-verify"),

		ErrandVMType:       c.String("errand-vm-type"),
		SmokeTestsPassword: c.String("smoke-tests-password"),
		UAALoginProtocol:   c.String("uaa-login-protocol"),
	}
	if err := config.loadSSL(c); err != nil {
		return nil, err
	}
	return config, nil
}

func (ca *Config) loadSSL(c *cli.Context) error {
	caCert, err := pluginutil.LoadResourceFromContext(c, "consul-server-ca-cert")
	if err != nil {
		return err
	}
	agentCert, err := pluginutil.LoadResourceFromContext(c, "consul-agent-cert")
	if err != nil {
		return err
	}
	agentKey, err := pluginutil.LoadResourceFromContext(c, "consul-agent-key")
	if err != nil {
		return err
	}
	serverCert, err := pluginutil.LoadResourceFromContext(c, "consul-server-cert")
	if err != nil {
		return err
	}
	serverKey, err := pluginutil.LoadResourceFromContext(c, "consul-server-key")
	if err != nil {
		return err
	}

	ca.ConsulCaCert = caCert
	ca.ConsulAgentCert = agentCert
	ca.ConsulServerCert = serverCert
	ca.ConsulAgentKey = agentKey
	ca.ConsulServerKey = serverKey
	return nil
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
