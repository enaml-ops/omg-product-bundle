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

	AllowSSHAccess     bool
	SkipSSLCertVerify  bool
	NATSUser           string
	NATSPassword       string
	NATSPort           int
	NATSMachines       []string
	JWTVerificationKey string
	ConsulCaCert       string
	ConsulEncryptKeys  []string
	ConsulAgentCert    string
	ConsulAgentKey     string
	ConsulServerCert   string
	ConsulServerKey    string

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
	CCBulkAPIPassword             string
	CCBulkAPIUser                 string
	DiegoCellVMType               string
	DiegoCellPersistentDiskType   string
	DiegoCellIPs                  []string
	ConsulIPs                     []string
	BBSCACert                     string
	BBSClientCert                 string
	BBSClientKey                  string
	BBSServerCert                 string
	BBSServerKey                  string

	DiegoBrainVMType                          string
	DiegoBrainPersistentDiskType              string
	DiegoBrainIPs                             []string
	BBSRequireSSL                             bool
	CCUploaderJobPollInterval                 int
	CCInternalAPIUser                         string
	CCInternalAPIPassword                     string
	CCBulkBatchSize                           int
	CCFetchTimeout                            int
	FSListenAddr                              string
	FSStaticDirectory                         string
	FSDebugAddr                               string
	FSLogLevel                                string
	MetronPort                                int
	SSHProxyClientSecret                      string
	CCExternalPort                            int
	TrafficControllerURL                      string
	RouterMachines                            []string
	EtcdServerCert                            string
	EtcdServerKey                             string
	EtcdClientCert                            string
	EtcdClientKey                             string
	EtcdPeerCert                              string
	EtcdPeerKey                               string
	MySQLProxyIPs                             []string
	UAAVMType                                 string
	UAAInstances                              int
	DiegoDBPassphrase                         string
	DiegoDBVMType                             string
	DiegoDBPersistentDiskType                 string
	DiegoDBIPs                                []string
	SAMLServiceProviderKey                    string
	SAMLServiceProviderCertificate            string
	JWTSigningKey                             string
	AdminSecret                               string
	UAADBUserName                             string
	UAADBPassword                             string
	PushAppsManagerPassword                   string
	SystemServicesPassword                    string
	SystemVerificationPassword                string
	OpentsdbFirehoseNozzleClientSecret        string
	IdentityClientSecret                      string
	LoginClientSecret                         string
	PortalClientSecret                        string
	AutoScalingServiceClientSecret            string
	SystemPasswordsClientSecret               string
	CCServiceDashboardsClientSecret           string
	GoRouterClientSecret                      string
	NotificationsClientSecret                 string
	NotificationsUIClientSecret               string
	CloudControllerUsernameLookupClientSecret string
	CCRoutingClientSecret                     string
	AppsMetricsClientSecret                   string
	AppsMetricsProcessingClientSecret         string
	SelfServiceLinksEnabled                   bool
	SignupsEnabled                            bool
	CompanyName                               string
	ProductLogo                               string
	SquareLogo                                string
	FooterLegalText                           string
	LDAPUrl                                   string
	LDAPUserDN                                string
	LDAPUserPassword                          string
	LDAPSearchBase                            string
	LDAPSearchFilter                          string
	LDAPMailAttributeName                     string
	LDAPEnabled                               bool

	CloudControllerWorkerInstances int
	CloudControllerWorkerVMType    string
	NFSServerAddress               string
	SharePath                      string
	StagingUploadUser              string
	StagingUploadPassword          string
	DbEncryptionKey                string
	InternalAPIUser                string
	InternalAPIPassword            string
	CCDBUsername                   string
	CCDBPassword                   string
	CloudControllerInstances       int
	CloudControllerVMType          string
	HostKeyFingerprint             string
	SupportAddress                 string
	MinCliVersion                  string
	ClockGlobalVMType              string
	HAProxySSLPem                  string
	HAProxySkip                    bool
	HAProxyVMType                  string
	HAProxyIPs                     []string
	MySQLProxyVMType               string
	MySQLProxyAPIUsername          string
	MySQLProxyAPIPassword          string
	MySQLProxyExternalHost         string
	RouterSSLCert                  string
	RouterSSLKey                   string
	RouterVMType                   string
	RouterEnableSSL                bool
	RouterUser                     string
	RouterPass                     string
	NFSVMType                      string
	NFSIPs                         []string
	NFSPersistentDiskType          string
	NFSAllowFromNetworkCIDR        []string
	EtcdVMType                     string
	EtcdPersistentDiskType         string
	NatsVMType                     string
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

		HostKeyFingerprint: c.String("host-key-fingerprint"),
		SupportAddress:     c.String("support-address"),
		MinCliVersion:      c.String("min-cli-version"),

		CloudControllerWorkerInstances: c.Int("cc-worker-instances"),
		CloudControllerWorkerVMType:    c.String("cc-worker-vm-type"),
		NFSServerAddress:               c.String("nfs-server-address"),
		SharePath:                      c.String("nfs-share-path"),
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
		CCBulkAPIPassword:             c.String("cc-bulk-api-password"),
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
		StagingUploadUser:            c.String("cc-staging-upload-user"),
		StagingUploadPassword:        c.String("cc-staging-upload-password"),
		CCBulkAPIUser:                c.String("cc-bulk-api-user"),
		DbEncryptionKey:              c.String("cc-db-encryption-key"),
		CCDBUsername:                 c.String("db-ccdb-username"),
		CCDBPassword:                 c.String("db-ccdb-password"),
		CloudControllerInstances:     c.Int("cc-instances"),
		CloudControllerVMType:        c.String("cc-vm-type"),

		DiegoDBVMType:                             c.String("diego-db-vm-type"),
		DiegoDBPersistentDiskType:                 c.String("diego-db-disk-type"),
		DiegoDBIPs:                                c.StringSlice("diego-db-ip"),
		DiegoDBPassphrase:                         c.String("diego-db-passphrase"),
		MySQLProxyIPs:                             c.StringSlice("mysql-proxy-ip"),
		UAAVMType:                                 c.String("uaa-vm-type"),
		UAAInstances:                              c.Int("uaa-instances"),
		SAMLServiceProviderKey:                    c.String("uaa-saml-service-provider-key"),
		SAMLServiceProviderCertificate:            c.String("uaa-saml-service-provider-cert"),
		JWTSigningKey:                             c.String("uaa-jwt-signing-key"),
		JWTVerificationKey:                        c.String("uaa-jwt-verification-key"),
		AdminSecret:                               c.String("uaa-admin-secret"),
		RouterMachines:                            c.StringSlice("router-ip"),
		UAADBUserName:                             c.String("db-uaa-username"),
		UAADBPassword:                             c.String("db-uaa-password"),
		PushAppsManagerPassword:                   c.String("push-apps-manager-password"),
		SystemServicesPassword:                    c.String("system-services-password"),
		SystemVerificationPassword:                c.String("system-verification-password"),
		OpentsdbFirehoseNozzleClientSecret:        c.String("opentsdb-firehose-nozzle-client-secret"),
		IdentityClientSecret:                      c.String("identity-client-secret"),
		LoginClientSecret:                         c.String("login-client-secret"),
		PortalClientSecret:                        c.String("portal-client-secret"),
		AutoScalingServiceClientSecret:            c.String("autoscaling-service-client-secret"),
		SystemPasswordsClientSecret:               c.String("system-passwords-client-secret"),
		CCServiceDashboardsClientSecret:           c.String("cc-service-dashboards-client-secret"),
		GoRouterClientSecret:                      c.String("gorouter-client-secret"),
		NotificationsClientSecret:                 c.String("notifications-client-secret"),
		NotificationsUIClientSecret:               c.String("notifications-ui-client-secret"),
		CloudControllerUsernameLookupClientSecret: c.String("cloud-controller-username-lookup-client-secret"),
		CCRoutingClientSecret:                     c.String("cc-routing-client-secret"),
		AppsMetricsClientSecret:                   c.String("apps-metrics-client-secret"),
		AppsMetricsProcessingClientSecret:         c.String("apps-metrics-processing-client-secret"),
		SelfServiceLinksEnabled:                   c.BoolT("uaa-enable-selfservice-links"),
		SignupsEnabled:                            c.BoolT("uaa-signups-enabled"),
		CompanyName:                               c.String("uaa-company-name"),
		ProductLogo:                               c.String("uaa-product-logo"),
		SquareLogo:                                c.String("uaa-square-logo"),
		FooterLegalText:                           c.String("uaa-footer-legal-txt"),
		LDAPUrl:                                   c.String("uaa-ldap-url"),
		LDAPUserDN:                                c.String("uaa-ldap-user-dn"),
		LDAPUserPassword:                          c.String("uaa-ldap-user-password"),
		LDAPSearchBase:                            c.String("uaa-ldap-search-base"),
		LDAPSearchFilter:                          c.String("uaa-ldap-search-filter"),
		LDAPMailAttributeName:                     c.String("uaa-ldap-mail-attributename"),
		LDAPEnabled:                               c.BoolT("uaa-ldap-enabled"),
		ClockGlobalVMType:                         c.String("clock-global-vm-type"),
		HAProxySkip:                               c.BoolT("skip-haproxy"),
		HAProxyIPs:                                c.StringSlice("haproxy-ip"),
		HAProxyVMType:                             c.String("haproxy-vm-type"),
		MySQLProxyVMType:                          c.String("mysql-proxy-vm-type"),
		MySQLProxyAPIUsername:                     c.String("mysql-proxy-api-username"),
		MySQLProxyAPIPassword:                     c.String("mysql-proxy-api-password"),
		MySQLProxyExternalHost:                    c.String("mysql-proxy-external-host"),
		RouterEnableSSL:                           c.Bool("router-enable-ssl"),
		RouterVMType:                              c.String("router-vm-type"),
		RouterUser:                                c.String("router-user"),
		RouterPass:                                c.String("router-pass"),
		NFSIPs:                                    c.StringSlice("nfs-ip"),
		NFSVMType:                                 c.String("nfs-vm-type"),
		NFSPersistentDiskType:                     c.String("nfs-disk-type"),
		NFSAllowFromNetworkCIDR:                   c.StringSlice("nfs-allow-from-network-cidr"),
		EtcdVMType:                                c.String("etcd-vm-type"),
		EtcdPersistentDiskType:                    c.String("etcd-disk-type"),
		NatsVMType:                                c.String("nats-vm-type"),
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

	etcdServerCert, err := pluginutil.LoadResourceFromContext(c, "etcd-server-cert")
	if err != nil {
		return err
	}

	etcdServerKey, err := pluginutil.LoadResourceFromContext(c, "etcd-server-key")
	if err != nil {
		return err
	}

	etcdClientCert, err := pluginutil.LoadResourceFromContext(c, "etcd-client-cert")
	if err != nil {
		return err
	}

	etcdClientKey, err := pluginutil.LoadResourceFromContext(c, "etcd-client-key")
	if err != nil {
		return err
	}

	etcdPeerCert, err := pluginutil.LoadResourceFromContext(c, "etcd-peer-cert")
	if err != nil {
		return err
	}

	etcdPeerKey, err := pluginutil.LoadResourceFromContext(c, "etcd-peer-key")
	if err != nil {
		return err
	}

	sslpem, err := pluginutil.LoadResourceFromContext(c, "haproxy-sslpem")
	if err != nil {
		return err
	}

	routerCert, err := pluginutil.LoadResourceFromContext(c, "router-ssl-cert")
	if err != nil {
		return err
	}
	routerKey, err := pluginutil.LoadResourceFromContext(c, "router-ssl-key")
	if err != nil {
		return err
	}
	ca.RouterSSLCert = routerCert
	ca.RouterSSLKey = routerKey
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
	ca.EtcdClientCert = etcdClientCert
	ca.EtcdClientKey = etcdClientKey
	ca.EtcdPeerCert = etcdPeerCert
	ca.EtcdPeerKey = etcdPeerKey
	ca.EtcdServerKey = etcdServerKey
	ca.EtcdServerCert = etcdServerCert
	ca.HAProxySSLPem = sslpem
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

func (c *Config) MySQLProxyHost() string {
	return c.MySQLProxyIPs[0]
}
