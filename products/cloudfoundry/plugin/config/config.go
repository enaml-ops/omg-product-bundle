package config

import (
	"github.com/enaml-ops/pluginlib/pcli"
	"gopkg.in/urfave/cli.v2"
)

func NewConfig(c *cli.Context) (*Config, error) {
	config := &Config{}
	err := pcli.UnmarshalFlags(config, c)
	if err != nil {
		return nil, err
	}

	secret, err := NewSecret(c)
	if err != nil {
		return nil, err
	}
	config.Secret = *secret
	config.User = NewUser(c)

	if config.MySQLProxyExternalHost == "" {
		config.MySQLProxyExternalHost = config.MySQLProxyIPs[0]
	}

	certs, err := NewCerts(c)
	if err != nil {
		return nil, err
	}
	config.Certs = certs
	return config, nil
}

type Config struct {
	DeploymentName                string   `omg:"deployment-name,optional"`
	AZs                           []string `omg:"az"`
	StemcellName                  string
	NetworkName                   string `omg:"network"`
	SystemDomain                  string
	AppDomains                    []string `omg:"app-domain"`
	AllowSSHAccess                bool     `omg:"allow-app-ssh-access"`
	SkipSSLCertVerify             bool     `omg:"skip-cert-verify"`
	NATSPort                      int      `omg:"nats-port"`
	UAALoginProtocol              string   `omg:"uaa-login-protocol"`
	SyslogAddress                 string   `omg:"syslog-address,optional"`
	SyslogPort                    int      `omg:"syslog-port,optional"`
	SyslogTransport               string   `omg:"syslog-transport,optional"`
	DopplerZone                   string   `omg:"doppler-zone,optional"`
	DopplerMessageDrainBufferSize int      `omg:"doppler-drain-buffer-size"`
	BBSRequireSSL                 bool     `omg:"bbs-require-ssl"`
	CCUploaderJobPollInterval     int      `omg:"cc-uploader-poll-interval"`
	CCExternalPort                int      `omg:"cc-external-port"`
	SelfServiceLinksEnabled       bool     `omg:"uaa-enable-selfservice-links"`
	SignupsEnabled                bool     `omg:"uaa-signups-enabled"`
	CompanyName                   string   `omg:"uaa-company-name"`
	ProductLogo                   string   `omg:"uaa-product-logo"`
	SquareLogo                    string   `omg:"uaa-square-logo"`
	FooterLegalText               string   `omg:"uaa-footer-legal-txt"`
	LDAPUrl                       string   `omg:"uaa-ldap-url,optional"`
	LDAPUserDN                    string   `omg:"uaa-ldap-user-dn,optional"`
	LDAPSearchBase                string   `omg:"uaa-ldap-search-base,optional"`
	LDAPSearchFilter              string   `omg:"uaa-ldap-search-filter,optional"`
	LDAPMailAttributeName         string   `omg:"uaa-ldap-mail-attributename,optional"`
	LDAPEnabled                   bool     `omg:"uaa-ldap-enabled,optional"`
	SharePath                     string   `omg:"nfs-share-path"`
	HostKeyFingerprint            string
	SupportAddress                string   `omg:"support-address,optional"`
	MinCliVersion                 string   `omg:"min-cli-version"`
	HAProxySkip                   bool     `omg:"skip-haproxy"`
	MySQLProxyExternalHost        string   `omg:"mysql-proxy-external-host,optional"`
	RouterEnableSSL               bool     `omg:"router-enable-ssl"`
	NFSAllowFromNetworkCIDR       []string `omg:"nfs-allow-from-network-cidr"`
	LoggregatorPort               int
	*Certs                        `omg:"-"` // certs are a special case that we handle manually for now
	IP
	VMType
	Disk
	Secret
	InstanceCount
	User
}

func (c *Config) MySQLProxyHost() string {
	return c.MySQLProxyIPs[0]
}
