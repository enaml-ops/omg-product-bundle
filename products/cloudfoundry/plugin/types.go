package cloudfoundry

import (
	"github.com/codegangsta/cli"
	"github.com/enaml-ops/enaml"
	etcdmetricslib "github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/etcd_metrics_server"
	grtrlib "github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/gorouter"
	"github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/metron_agent"
	natslib "github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/nats"
)

//VaultRotater an interface for rotating vault hashes values
type VaultRotater interface {
	RotateSecrets(hash string, secrets interface{}) error
}

// InstanceGrouper creates and validates InstanceGroups.
type InstanceGrouper interface {
	ToInstanceGroup() (ig *enaml.InstanceGroup)
	HasValidValues() bool
}

// InstanceGrouperFactory is a function that creates InstanceGroupers from CLI args.
type InstanceGrouperFactory func(*cli.Context) InstanceGrouper

type InstanceGrouperConfigFactory func(*cli.Context, *Config) InstanceGrouper

type gorouter struct {
	Instances    int
	AZs          []string
	StemcellName string
	VMTypeName   string
	NetworkName  string
	NetworkIPs   []string
	SSLCert      string
	SSLKey       string
	EnableSSL    bool
	ClientSecret string
	Nats         grtrlib.Nats
	Loggregator  metron_agent.Loggregator
	RouterUser   string
	RouterPass   string
	MetronZone   string
	MetronSecret string
}

//Etcd -
type Etcd struct {
	AZs                []string
	StemcellName       string
	VMTypeName         string
	NetworkName        string
	NetworkIPs         []string
	PersistentDiskType string
	Metron             *Metron
	StatsdInjector     *StatsdInjector
	Nats               *etcdmetricslib.Nats
}

//Metron -
type Metron struct {
	Zone            string
	Secret          string
	SyslogAddress   string
	SyslogPort      int
	SyslogTransport string
	Loggregator     metron_agent.Loggregator
}

//StatsdInjector -
type StatsdInjector struct {
}

//NFSMounter -
type NFSMounter struct {
	NFSServerAddress string
	SharePath        string
}

//NatsPartition -
type NatsPartition struct {
	AZs            []string
	StemcellName   string
	VMTypeName     string
	NetworkName    string
	NetworkIPs     []string
	Nats           natslib.NatsJob
	Metron         *Metron
	StatsdInjector *StatsdInjector
}

//NFS -
type NFS struct {
	AZs                  []string
	StemcellName         string
	VMTypeName           string
	NetworkName          string
	NetworkIPs           []string
	PersistentDiskType   string
	AllowFromNetworkCIDR []string
	Metron               *Metron
	StatsdInjector       *StatsdInjector
}

//UAAClient - Structure to represent map of client priviledges
type UAAClient struct {
	ID                   string      `yaml:"id,omitempty"`
	Secret               string      `yaml:"secret,omitempty"`
	Scope                string      `yaml:"scope,omitempty"`
	AuthorizedGrantTypes string      `yaml:"authorized-grant-types,omitempty"`
	Authorities          string      `yaml:"authorities,omitempty"`
	AutoApprove          interface{} `yaml:"autoapprove,omitempty"`
	Override             bool        `yaml:"override,omitempty"`
	RedirectURI          string      `yaml:"redirect-uri,omitempty"`
	AccessTokenValidity  int         `yaml:"access-token-validity,omitempty"`
	RefreshTokenValidity int         `yaml:"refresh-token-validity,omitempty"`
	ResourceIDs          string      `yaml:"resource_ids,omitempty"`
	Name                 string      `yaml:"name,omitempty"`
	AppLaunchURL         string      `yaml:"app-launch-url,omitempty"`
	ShowOnHomepage       bool        `yaml:"show-on-homepage,omitempty"`
	AppIcon              string      `yaml:"app-icon,omitempty"`
}

//Plugin -
type Plugin struct {
	PluginVersion string
}
