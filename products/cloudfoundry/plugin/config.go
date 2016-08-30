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

	ErrandVMType                  string
	SmokeTestsPassword            string
	UAALoginProtocol              string
	LoggregratorIPs               []string
	LoggregratorVMType            string
	DopplerSecret                 string
	EtcdMachines                  []string
	MetronZone                    string
	MetronSecret                  string
	SyslogAddress                 string
	SyslogPort                    int
	SyslogTransport               string
	DopplerIPs                    []string
	DopplerVMType                 string
	DopplerZone                   string
	DopplerMessageDrainBufferSize int
	DopplerSharedSecret           string
	CCBuilkAPIPassword            string
	DiegoCellVMType               string
	DiegoCellPersistentDiskType   string
	DiegoCellIPs                  []string
	ConsulIPs                     []string
	BBSCACert                     string
	BBSClientCert                 string
	BBSClientKey                  string
	BBSServerCert                 string
	BBSServerKey                  string

	DiegoBrainVMType             string
	DiegoBrainPersistentDiskType string
	DiegoBrainIPs                []string
	BBSRequireSSL                bool
	CCUploaderJobPollInterval    int
	CCInternalAPIUser            string
	CCInternalAPIPassword        string
	CCBulkBatchSize              int
	CCFetchTimeout               int
	FSListenAddr                 string
	FSStaticDirectory            string
	FSDebugAddr                  string
	FSLogLevel                   string
	MetronPort                   int
	SSHProxyClientSecret         string
	CCExternalPort               int
	TrafficControllerURL         string
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

		LoggregratorIPs:    c.StringSlice("loggregator-traffic-controller-ip"),
		LoggregratorVMType: c.String("loggregator-traffic-controller-vmtype"),

		EtcdMachines: c.StringSlice("etcd-machine-ip"),

		MetronZone:      c.String("metron-zone"),
		MetronSecret:    c.String("metron-secret"),
		SyslogAddress:   c.String("syslog-address"),
		SyslogPort:      c.Int("syslog-port"),
		SyslogTransport: c.String("syslog-transport"),

		DopplerIPs:                    c.StringSlice("doppler-ip"),
		DopplerVMType:                 c.String("doppler-vm-type"),
		DopplerSecret:                 c.String("doppler-client-secret"),
		DopplerZone:                   c.String("doppler-zone"),
		DopplerMessageDrainBufferSize: c.Int("doppler-drain-buffer-size"),
		DopplerSharedSecret:           c.String("doppler-shared-secret"),
		CCBuilkAPIPassword:            c.String("cc-bulk-api-password"),
		DiegoCellVMType:               c.String("diego-cell-vm-type"),
		DiegoCellPersistentDiskType:   c.String("diego-cell-disk-type"),
		DiegoCellIPs:                  c.StringSlice("diego-cell-ip"),

		ConsulIPs: c.StringSlice("consul-ip"),

		DiegoBrainVMType:             c.String("diego-brain-vm-type"),
		DiegoBrainPersistentDiskType: c.String("diego-brain-disk-type"),
		DiegoBrainIPs:                c.StringSlice("diego-brain-ip"),
		BBSRequireSSL:                c.BoolT("bbs-require-ssl"),
		CCUploaderJobPollInterval:    c.Int("cc-uploader-poll-interval"),
		CCInternalAPIUser:            c.String("cc-internal-api-user"),
		CCInternalAPIPassword:        c.String("cc-internal-api-password"),
		CCFetchTimeout:               c.Int("cc-fetch-timeout"),
		CCBulkBatchSize:              c.Int("cc-bulk-batch-size"),
		FSListenAddr:                 c.String("fs-listen-addr"),
		FSStaticDirectory:            c.String("fs-static-dir"),
		FSDebugAddr:                  c.String("fs-debug-addr"),
		FSLogLevel:                   c.String("fs-log-level"),
		MetronPort:                   c.Int("metron-port"),
		SSHProxyClientSecret:         c.String("ssh-proxy-uaa-secret"),
		CCExternalPort:               c.Int("cc-external-port"),
		TrafficControllerURL:         c.String("traffic-controller-url"),
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
	bbsCaCert, err := pluginutil.LoadResourceFromContext(c, "bbs-server-ca-cert")
	if err != nil {
		return err
	}

	bbsClientCert, err := pluginutil.LoadResourceFromContext(c, "bbs-client-cert")
	if err != nil {
		return err
	}

	bbsClientKey, err := pluginutil.LoadResourceFromContext(c, "bbs-client-key")
	if err != nil {
		return err
	}

	bbsServerCert, err := pluginutil.LoadResourceFromContext(c, "bbs-server-cert")
	if err != nil {
		return err
	}

	bbsServerKey, err := pluginutil.LoadResourceFromContext(c, "bbs-server-key")
	if err != nil {
		return err
	}
	ca.ConsulCaCert = caCert
	ca.ConsulAgentCert = agentCert
	ca.ConsulServerCert = serverCert
	ca.ConsulAgentKey = agentKey
	ca.ConsulServerKey = serverKey
	ca.BBSCACert = bbsCaCert
	ca.BBSClientCert = bbsClientCert
	ca.BBSClientKey = bbsClientKey
	ca.BBSServerCert = bbsServerCert
	ca.BBSServerKey = bbsServerKey
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
