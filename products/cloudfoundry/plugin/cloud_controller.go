package cloudfoundry

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/enaml-ops/enaml"
	ccnglib "github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/cloud_controller_ng"
	"github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/route_registrar"
	"github.com/xchapter7x/lo"
)

//CloudControllerPartition - Cloud Controller Partition
type CloudControllerPartition struct {
	Config                   *Config
	Instances                int
	VMTypeName               string
	AllowedCorsDomains       []string
	Metron                   *Metron
	ConsulAgent              *ConsulAgent
	StatsdInjector           *StatsdInjector
	NFSMounter               *NFSMounter
	StagingUploadUser        string
	StagingUploadPassword    string
	BulkAPIUser              string
	BulkAPIPassword          string
	DbEncryptionKey          string
	InternalAPIUser          string
	InternalAPIPassword      string
	HostKeyFingerprint       string
	SupportAddress           string
	MinCliVersion            string
	CCDBUsername             string
	CCDBPassword             string
	MySQLProxyIP             string
	UAAJWTVerificationKey    string
	CCServiceDashboardSecret string
	CCUsernameLookupSecret   string
	CCRoutingSecret          string
}

func NewCloudControllerPartition(c *cli.Context, config *Config) InstanceGrouper {
	var proxyIP string
	mysqlProxies := c.StringSlice("mysql-proxy-ip")
	if len(mysqlProxies) > 0 {
		proxyIP = mysqlProxies[0]
	}
	return &CloudControllerPartition{
		Config:                   config,
		Instances:                c.Int("cc-instances"),
		VMTypeName:               c.String("cc-vm-type"),
		Metron:                   NewMetron(c),
		ConsulAgent:              NewConsulAgent(c, []string{}, config),
		NFSMounter:               NewNFSMounter(c),
		StatsdInjector:           NewStatsdInjector(c),
		StagingUploadUser:        c.String("cc-staging-upload-user"),
		StagingUploadPassword:    c.String("cc-staging-upload-password"),
		BulkAPIUser:              c.String("cc-bulk-api-user"),
		BulkAPIPassword:          c.String("cc-bulk-api-password"),
		InternalAPIUser:          c.String("cc-internal-api-user"),
		InternalAPIPassword:      c.String("cc-internal-api-password"),
		DbEncryptionKey:          c.String("cc-db-encryption-key"),
		HostKeyFingerprint:       c.String("host-key-fingerprint"),
		SupportAddress:           c.String("support-address"),
		MinCliVersion:            c.String("min-cli-version"),
		CCDBUsername:             c.String("db-ccdb-username"),
		CCDBPassword:             c.String("db-ccdb-password"),
		MySQLProxyIP:             proxyIP,
		UAAJWTVerificationKey:    c.String("uaa-jwt-verification-key"),
		CCServiceDashboardSecret: c.String("cc-service-dashboards-client-secret"),
		CCUsernameLookupSecret:   c.String("cloud-controller-username-lookup-client-secret"),
		CCRoutingSecret:          c.String("cc-routing-client-secret"),
	}
}

//ToInstanceGroup - Convert CLoud Controller Partition to an Instance Group
func (s *CloudControllerPartition) ToInstanceGroup() (ig *enaml.InstanceGroup) {
	buildpacks := []string{"go-buildpack",
		"binary-buildpack",
		"nodejs-buildpack",
		"ruby-buildpack",
		"php-buildpack",
		"python-buildpack",
		"java-offline-buildpack",
		"staticfile-buildpack"}

	ig = &enaml.InstanceGroup{
		Name:      "cloud_controller-partition",
		AZs:       s.Config.AZs,
		Instances: s.Instances,
		VMType:    s.VMTypeName,
		Stemcell:  s.Config.StemcellName,
		Networks: []enaml.Network{
			enaml.Network{Name: s.Config.NetworkName},
		},
		Jobs: []enaml.InstanceJob{
			newCloudControllerNgJob(s),
			s.ConsulAgent.CreateJob(),
			s.NFSMounter.CreateJob(),
			s.Metron.CreateJob(),
			s.StatsdInjector.CreateJob(),
			newRouteRegistrarJob(s),
		},

		Update: enaml.Update{
			MaxInFlight: 1,
			Serial:      true,
		},
	}
	for _, buildpack := range buildpacks {
		ig.Jobs = append(ig.Jobs, enaml.InstanceJob{Name: buildpack, Release: CFReleaseName})
	}
	return
}

func newCloudControllerNgJob(c *CloudControllerPartition) enaml.InstanceJob {
	return enaml.InstanceJob{
		Name:    "cloud_controller_ng",
		Release: CFReleaseName,
		Properties: &ccnglib.CloudControllerNgJob{
			AppSsh: &ccnglib.AppSsh{
				HostKeyFingerprint: c.HostKeyFingerprint,
			},
			Domain:                   c.Config.SystemDomain,
			SystemDomain:             c.Config.SystemDomain,
			AppDomains:               c.Config.AppDomains,
			SystemDomainOrganization: "system",
			SupportAddress:           c.SupportAddress,
			Login: &ccnglib.Login{
				Url: fmt.Sprintf("https://login.%s", c.Config.SystemDomain),
			},
			Cc: &ccnglib.Cc{
				AllowedCorsDomains:    []string{fmt.Sprintf("https://login.%s", c.Config.SystemDomain)},
				AllowAppSshAccess:     c.Config.AllowSSHAccess,
				DefaultToDiegoBackend: true,
				Buildpacks: &ccnglib.Buildpacks{
					BlobstoreType: "fog",
					FogConnection: &ccnglib.DefaultFogConnection{
						Provider:  "Local",
						LocalRoot: "/var/vcap/nfs/shared",
					},
				},
				Droplets: &ccnglib.Droplets{
					BlobstoreType: "fog",
					FogConnection: &ccnglib.DefaultFogConnection{
						Provider:  "Local",
						LocalRoot: "/var/vcap/nfs/shared",
					},
				},
				Packages: &ccnglib.Packages{
					BlobstoreType: "fog",
					FogConnection: &ccnglib.DefaultFogConnection{
						Provider:  "Local",
						LocalRoot: "/var/vcap/nfs/shared",
					},
				},
				ResourcePool: &ccnglib.ResourcePool{
					BlobstoreType: "fog",
					FogConnection: &ccnglib.DefaultFogConnection{
						Provider:  "Local",
						LocalRoot: "/var/vcap/nfs/shared",
					},
				},
				ClientMaxBodySize:            "1024M",
				ExternalProtocol:             "https",
				LoggingLevel:                 "debug",
				MaximumHealthCheckTimeout:    600,
				StagingUploadUser:            c.StagingUploadUser,
				StagingUploadPassword:        c.StagingUploadPassword,
				BulkApiUser:                  c.BulkAPIUser,
				BulkApiPassword:              c.BulkAPIPassword,
				InternalApiUser:              c.InternalAPIUser,
				InternalApiPassword:          c.InternalAPIPassword,
				DbEncryptionKey:              c.DbEncryptionKey,
				DefaultRunningSecurityGroups: []string{"all_open"},
				DefaultStagingSecurityGroups: []string{"all_open"},
				DisableCustomBuildpacks:      false,
				ExternalHost:                 "api",
				InstallBuildpacks: []map[string]interface{}{
					map[string]interface{}{
						"name":    "staticfile_buildpack",
						"package": "staticfile-buildpack",
					},
					map[string]interface{}{
						"name":    "java_buildpack_offline",
						"package": "java-offline-buildpack",
					},
					map[string]interface{}{
						"name":    "ruby_buildpack",
						"package": "ruby-buildpack",
					},
					map[string]interface{}{
						"name":    "nodejs_buildpack",
						"package": "nodejs-buildpack",
					},
					map[string]interface{}{
						"name":    "go_buildpack",
						"package": "go-buildpack",
					},
					map[string]interface{}{
						"name":    "python_buildpack",
						"package": "python-buildpack",
					},
					map[string]interface{}{
						"name":    "php_buildpack",
						"package": "php-buildpack",
					},
					map[string]interface{}{
						"name":    "binary_buildpack",
						"package": "binary-buildpack",
					},
				},
				QuotaDefinitions: map[string]interface{}{
					"default": map[string]interface{}{
						"memory_limit":               10240,
						"total_services":             100,
						"non_basic_services_allowed": true,
						"total_routes":               1000,
						"trial_db_allowed":           true,
					},
					"runaway": map[string]interface{}{
						"memory_limit":               102400,
						"total_services":             -1,
						"total_routes":               1000,
						"non_basic_services_allowed": true,
					},
				},
				SecurityGroupDefinitions: []map[string]interface{}{
					map[string]interface{}{
						"name": "all_open",
						"rules": []map[string]interface{}{
							map[string]interface{}{
								"protocol":    "all",
								"destination": "0.0.0.0-255.255.255.255",
							},
						},
					},
				},
				Stacks: []map[string]interface{}{
					map[string]interface{}{
						"name":        "cflinuxfs2",
						"description": "Cloud Foundry Linux-based filesystem",
					},
					map[string]interface{}{
						"name":        "windows2012R2",
						"description": "Microsoft Windows / .NET 64 bit",
					},
				},
				UaaResourceId:            "cloud_controller,cloud_controller_service_permissions",
				MinCliVersion:            c.MinCliVersion,
				MinRecommendedCliVersion: c.MinCliVersion,
			},
			Ccdb: &ccnglib.Ccdb{
				Address:  c.MySQLProxyIP,
				Port:     3306,
				DbScheme: "mysql",
				Roles: []map[string]interface{}{
					map[string]interface{}{
						"tag":      "admin",
						"name":     c.CCDBUsername,
						"password": c.CCDBPassword,
					},
				},
				Databases: []map[string]interface{}{
					map[string]interface{}{
						"tag":    "cc",
						"name":   "ccdb",
						"citext": true,
					},
				},
			},
			Uaa: &ccnglib.Uaa{
				Url: fmt.Sprintf("https://uaa.%s", c.Config.SystemDomain),
				Jwt: &ccnglib.Jwt{
					VerificationKey: c.UAAJWTVerificationKey,
				},
				Clients: &ccnglib.Clients{
					CcServiceDashboards: &ccnglib.CcServiceDashboards{
						Scope:  "cloud_controller.write,openid,cloud_controller.read,cloud_controller_service_permissions.read",
						Secret: c.CCServiceDashboardSecret,
					},
					CloudControllerUsernameLookup: &ccnglib.CloudControllerUsernameLookup{
						Secret: c.CCUsernameLookupSecret,
					},
					CcRouting: &ccnglib.CcRouting{
						Secret: c.CCRoutingSecret,
					},
				},
			},
			Ssl: &ccnglib.Ssl{
				SkipCertVerify: c.Config.SkipSSLCertVerify,
			},
			LoggerEndpoint: &ccnglib.LoggerEndpoint{
				Port: 443,
			},
			Doppler: &ccnglib.Doppler{
				Port: 443,
			},
			NfsServer: &ccnglib.NfsServer{
				Address:   c.NFSMounter.NFSServerAddress,
				SharePath: "/var/vcap/nfs",
			},
			Nats: &ccnglib.Nats{
				User:     c.Config.NATSUser,
				Password: c.Config.NATSPassword,
				Port:     c.Config.NATSPort,
				Machines: c.Config.NATSMachines,
			},
		},
	}
}

func newRouteRegistrarJob(c *CloudControllerPartition) enaml.InstanceJob {
	return enaml.InstanceJob{
		Name:    "route_registrar",
		Release: CFReleaseName,
		Properties: route_registrar.RouteRegistrarJob{
			RouteRegistrar: &route_registrar.RouteRegistrar{
				Routes: []map[string]interface{}{
					map[string]interface{}{
						"name":                  "api",
						"port":                  9022,
						"registration_interval": "20s",
						"tags": map[string]interface{}{
							"component": "CloudController",
						},
						"uris": []string{fmt.Sprintf("api.%s", c.Config.SystemDomain)},
					},
				},
			},
			Nats: &route_registrar.Nats{
				User:     c.Config.NATSUser,
				Password: c.Config.NATSPassword,
				Port:     c.Config.NATSPort,
				Machines: c.Config.NATSMachines,
			},
		},
	}
}

//HasValidValues - Check if valid values has been populated
func (s *CloudControllerPartition) HasValidValues() bool {

	lo.G.Debugf("checking '%s' for valid flags", "cloud controller")

	if s.VMTypeName == "" {
		lo.G.Debugf("could not find a valid vmtypename '%v'", s.VMTypeName)
	}

	if s.Metron.Zone == "" {
		lo.G.Debugf("could not find a valid metron zone '%v'", s.Metron.Zone)
	}

	if s.Metron.Secret == "" {
		lo.G.Debugf("could not find a valid metron secret '%v'", s.Metron.Secret)
	}

	if s.MySQLProxyIP == "" {
		lo.G.Debug("missing mysql proxy IP")
	}

	return (s.VMTypeName != "" &&
		s.Metron.Zone != "" &&
		s.Metron.Secret != "" &&
		s.NFSMounter.hasValidValues() &&
		s.ConsulAgent.HasValidValues()) &&
		s.MySQLProxyIP != ""
}
