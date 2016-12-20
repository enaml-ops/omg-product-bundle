package config

import (
	"fmt"

	"github.com/xchapter7x/lo"
	"gopkg.in/urfave/cli.v2"
)

func NewConfig(c *cli.Context) (*Config, error) {
	//check for required fields
	if err := checkRequired(c); err != nil {
		return nil, err
	}
	config := &Config{
		DeploymentName:                c.String("deployment-name"),
		AZs:                           c.StringSlice("az"),
		StemcellName:                  c.String("stemcell-name"),
		NetworkName:                   c.String("network"),
		SystemDomain:                  c.String("system-domain"),
		AppDomains:                    c.StringSlice("app-domain"),
		NATSPort:                      c.Int("nats-port"),
		HostKeyFingerprint:            c.String("host-key-fingerprint"),
		SupportAddress:                c.String("support-address"),
		MinCliVersion:                 c.String("min-cli-version"),
		SharePath:                     c.String("nfs-share-path"),
		AllowSSHAccess:                c.Bool("allow-app-ssh-access"),
		SkipSSLCertVerify:             c.Bool("skip-cert-verify"),
		UAALoginProtocol:              c.String("uaa-login-protocol"),
		SyslogAddress:                 c.String("syslog-address"),
		SyslogPort:                    c.Int("syslog-port"),
		SyslogTransport:               c.String("syslog-transport"),
		DopplerZone:                   c.String("doppler-zone"),
		DopplerMessageDrainBufferSize: c.Int("doppler-drain-buffer-size"),
		BBSRequireSSL:                 c.Bool("bbs-require-ssl"),
		CCUploaderJobPollInterval:     c.Int("cc-uploader-poll-interval"),
		CCExternalPort:                c.Int("cc-external-port"),
		SelfServiceLinksEnabled:       c.Bool("uaa-enable-selfservice-links"),
		SignupsEnabled:                c.Bool("uaa-signups-enabled"),
		CompanyName:                   c.String("uaa-company-name"),
		ProductLogo:                   c.String("uaa-product-logo"),
		SquareLogo:                    c.String("uaa-square-logo"),
		FooterLegalText:               c.String("uaa-footer-legal-txt"),
		LDAPUrl:                       c.String("uaa-ldap-url"),
		LDAPUserDN:                    c.String("uaa-ldap-user-dn"),
		LDAPSearchBase:                c.String("uaa-ldap-search-base"),
		LDAPSearchFilter:              c.String("uaa-ldap-search-filter"),
		LDAPMailAttributeName:         c.String("uaa-ldap-mail-attributename"),
		LDAPEnabled:                   c.Bool("uaa-ldap-enabled"),
		HAProxySkip:                   c.Bool("skip-haproxy"),
		MySQLProxyExternalHost:        c.String("mysql-proxy-external-host"),
		RouterEnableSSL:               c.Bool("router-enable-ssl"),
		NFSAllowFromNetworkCIDR:       c.StringSlice("nfs-allow-from-network-cidr"),
		LoggregatorPort:               c.Int("loggregator-port"),
	}
	config.IP = NewIP(c)
	config.VMType = NewVMType(c)
	config.Disk = NewDisk(c)
	config.Secret = NewSecret(c)
	config.User = NewUser(c)
	config.InstanceCount = NewInstanceCount(c)

	if config.MySQLProxyExternalHost == "" {
		config.MySQLProxyExternalHost = config.MySQLProxyIPs[0]
	}

	if certs, err := NewCerts(c); err != nil {
		return nil, err
	} else {
		config.Certs = certs
	}
	return config, nil
}

type Config struct {
	DeploymentName                string
	AZs                           []string
	StemcellName                  string
	NetworkName                   string
	SystemDomain                  string
	AppDomains                    []string
	AllowSSHAccess                bool
	SkipSSLCertVerify             bool
	NATSPort                      int
	UAALoginProtocol              string
	SyslogAddress                 string
	SyslogPort                    int
	SyslogTransport               string
	DopplerZone                   string
	DopplerMessageDrainBufferSize int
	BBSRequireSSL                 bool
	CCUploaderJobPollInterval     int
	CCExternalPort                int
	SelfServiceLinksEnabled       bool
	SignupsEnabled                bool
	CompanyName                   string
	ProductLogo                   string
	SquareLogo                    string
	FooterLegalText               string
	LDAPUrl                       string
	LDAPUserDN                    string
	LDAPSearchBase                string
	LDAPSearchFilter              string
	LDAPMailAttributeName         string
	LDAPEnabled                   bool
	SharePath                     string
	HostKeyFingerprint            string
	SupportAddress                string
	MinCliVersion                 string
	HAProxySkip                   bool
	MySQLProxyExternalHost        string
	RouterEnableSSL               bool
	NFSAllowFromNetworkCIDR       []string
	LoggregatorPort               int
	*Certs
	IP
	VMType
	Disk
	Secret
	InstanceCount
	User
}

func RequiredStringFlags() []string {
	return []string{
		"stemcell-name",
		"network",
		"system-domain",
		"host-key-fingerprint",
		"support-address",
		"min-cli-version",
		"nfs-share-path",
		"uaa-login-protocol",
		"doppler-zone",
		"uaa-company-name",
		"uaa-product-logo",
		"uaa-square-logo",
		"uaa-footer-legal-txt",
	}
}

func RequiredStringSliceFlags() []string {
	return []string{
		"az",
		"app-domain",
		"nfs-allow-from-network-cidr",
	}
}

func RequiredIntFlags() []string {
	return []string{
		"nats-port",
		"doppler-drain-buffer-size",
		"cc-uploader-poll-interval",
		"cc-external-port",
		"loggregator-port",
	}
}
func (c *Config) MySQLProxyHost() string {
	return c.MySQLProxyIPs[0]
}

func checkRequired(c *cli.Context) error {
	invalidNames := []string{}
	invalidNames = append(invalidNames, checkRequiredStringFlags(c, RequiredStringFlags())...)
	invalidNames = append(invalidNames, checkRequiredStringSliceFlags(c, RequiredStringSliceFlags())...)
	invalidNames = append(invalidNames, checkRequiredIntFlags(c, RequiredIntFlags())...)
	invalidNames = append(invalidNames, checkRequiredStringFlags(c, RequiredCertFlags())...)
	invalidNames = append(invalidNames, checkRequiredStringFlags(c, RequiredDiskFlags())...)
	invalidNames = append(invalidNames, checkRequiredIntFlags(c, RequiredInstanceCountFlags())...)
	invalidNames = append(invalidNames, checkRequiredStringFlags(c, RequiredSecretFlags())...)
	invalidNames = append(invalidNames, checkRequiredStringFlags(c, RequiredUserFlags())...)
	invalidNames = append(invalidNames, checkRequiredStringFlags(c, RequiredVMTypeFlags())...)
	invalidNames = append(invalidNames, checkRequiredStringFlags(c, RequiredIPFlags())...)
	invalidNames = append(invalidNames, checkRequiredStringSliceFlags(c, RequiredIPSliceFlags())...)
	if len(invalidNames) > 0 {
		return fmt.Errorf("Sorry you need to provide %v flags to continue", invalidNames)
	}
	return nil
}

func checkRequiredIntFlags(c *cli.Context, requiredFlags []string) []string {
	invalidNames := []string{}
	for _, name := range requiredFlags {
		if c.Int(name) == 0 {
			invalidNames = append(invalidNames, name)
		} else {
			lo.G.Debug(name, "==>", c.String(name))
		}
	}
	return invalidNames
}

func checkRequiredStringFlags(c *cli.Context, requiredFlags []string) []string {
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

func checkRequiredStringSliceFlags(c *cli.Context, requiredFlags []string) []string {
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
