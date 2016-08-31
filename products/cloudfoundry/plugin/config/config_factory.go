package config

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/xchapter7x/lo"
)

func NewConfig(c *cli.Context) (*Config, error) {

	//check for required fields
	if err := checkRequired(c); err != nil {
		return nil, err
	}
	config := &Config{
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
		SkipSSLCertVerify:             c.BoolT("skip-cert-verify"),
		UAALoginProtocol:              c.String("uaa-login-protocol"),
		MetronZone:                    c.String("metron-zone"),
		SyslogAddress:                 c.String("syslog-address"),
		SyslogPort:                    c.Int("syslog-port"),
		SyslogTransport:               c.String("syslog-transport"),
		DopplerZone:                   c.String("doppler-zone"),
		DopplerMessageDrainBufferSize: c.Int("doppler-drain-buffer-size"),
		BBSRequireSSL:                 c.BoolT("bbs-require-ssl"),
		CCUploaderJobPollInterval:     c.Int("cc-uploader-poll-interval"),
		CCInternalAPIUser:             c.String("cc-internal-api-user"),
		CCFetchTimeout:                c.Int("cc-fetch-timeout"),
		CCBulkBatchSize:               c.Int("cc-bulk-batch-size"),
		FSListenAddr:                  c.String("fs-listen-addr"),
		FSStaticDirectory:             c.String("fs-static-dir"),
		FSDebugAddr:                   c.String("fs-debug-addr"),
		FSLogLevel:                    c.String("fs-log-level"),
		MetronPort:                    c.Int("metron-port"),
		CCExternalPort:                c.Int("cc-external-port"),
		TrafficControllerURL:          c.String("traffic-controller-url"),
		SelfServiceLinksEnabled:       c.BoolT("uaa-enable-selfservice-links"),
		SignupsEnabled:                c.BoolT("uaa-signups-enabled"),
		CompanyName:                   c.String("uaa-company-name"),
		ProductLogo:                   c.String("uaa-product-logo"),
		SquareLogo:                    c.String("uaa-square-logo"),
		FooterLegalText:               c.String("uaa-footer-legal-txt"),
		LDAPUrl:                       c.String("uaa-ldap-url"),
		LDAPUserDN:                    c.String("uaa-ldap-user-dn"),
		LDAPSearchBase:                c.String("uaa-ldap-search-base"),
		LDAPSearchFilter:              c.String("uaa-ldap-search-filter"),
		LDAPMailAttributeName:         c.String("uaa-ldap-mail-attributename"),
		LDAPEnabled:                   c.BoolT("uaa-ldap-enabled"),
		HAProxySkip:                   c.BoolT("skip-haproxy"),
		MySQLProxyExternalHost:        c.String("mysql-proxy-external-host"),
		RouterEnableSSL:               c.Bool("router-enable-ssl"),
		NFSAllowFromNetworkCIDR:       c.StringSlice("nfs-allow-from-network-cidr"),
	}
	config.IP = NewIP(c)
	config.VMType = NewVMType(c)
	config.Disk = NewDisk(c)
	config.Secret = NewSecret(c)
	config.User = NewUser(c)
	if certs, err := NewCerts(c); err != nil {
		return nil, err
	} else {
		config.Certs = certs
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
