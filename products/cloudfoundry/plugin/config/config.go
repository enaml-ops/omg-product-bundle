package config

type Config struct {
	AZs                           []string
	StemcellName                  string
	NetworkName                   string
	SystemDomain                  string
	AppDomains                    []string
	AllowSSHAccess                bool
	SkipSSLCertVerify             bool
	NATSPort                      int
	UAALoginProtocol              string
	MetronZone                    string
	SyslogAddress                 string
	SyslogPort                    int
	SyslogTransport               string
	DopplerZone                   string
	DopplerMessageDrainBufferSize int
	BBSRequireSSL                 bool
	CCUploaderJobPollInterval     int
	CCInternalAPIUser             string
	CCBulkBatchSize               int
	CCFetchTimeout                int
	FSListenAddr                  string
	FSStaticDirectory             string
	FSDebugAddr                   string
	FSLogLevel                    string
	MetronPort                    int
	CCExternalPort                int
	TrafficControllerURL          string
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
	*Certs
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
