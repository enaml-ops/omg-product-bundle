package cloudfoundry

import (
	"fmt"

	"github.com/enaml-ops/enaml"
	ccnglib "github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/cloud_controller_ng"
	"github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/route_registrar"
	"github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/plugin/config"
)

//CloudControllerPartition - Cloud Controller Partition
type CloudControllerPartition struct {
	Config         *config.Config
	Metron         *Metron
	ConsulAgent    *ConsulAgent
	StatsdInjector *StatsdInjector
}

func NewCloudControllerPartition(config *config.Config) InstanceGroupCreator {

	return &CloudControllerPartition{
		Config:         config,
		Metron:         NewMetron(config),
		ConsulAgent:    NewConsulAgent([]string{}, config),
		StatsdInjector: NewStatsdInjector(nil),
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
		Instances: s.Config.CloudControllerInstances,
		VMType:    s.Config.CloudControllerVMType,
		Stemcell:  s.Config.StemcellName,
		Networks: []enaml.Network{
			enaml.Network{Name: s.Config.NetworkName},
		},
		Jobs: []enaml.InstanceJob{
			newCloudControllerNgJob(s),
			s.ConsulAgent.CreateJob(),
			CreateNFSMounterJob(s.Config),
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
		ig.Jobs = append(ig.Jobs, enaml.InstanceJob{
			Name:       buildpack,
			Release:    CFReleaseName,
			Properties: make(map[interface{}]interface{}),
		})
	}
	return
}

func newCloudControllerNgJob(c *CloudControllerPartition) enaml.InstanceJob {
	return enaml.InstanceJob{
		Name:    "cloud_controller_ng",
		Release: CFReleaseName,
		Properties: &ccnglib.CloudControllerNgJob{
			AppSsh: &ccnglib.AppSsh{
				HostKeyFingerprint: c.Config.HostKeyFingerprint,
			},
			Domain:                   c.Config.SystemDomain,
			SystemDomain:             c.Config.SystemDomain,
			AppDomains:               c.Config.AppDomains,
			SystemDomainOrganization: "system",
			SupportAddress:           c.Config.SupportAddress,
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
				StagingUploadUser:            c.Config.StagingUploadUser,
				StagingUploadPassword:        c.Config.StagingUploadPassword,
				BulkApiUser:                  c.Config.CCBulkAPIUser,
				BulkApiPassword:              c.Config.CCBulkAPIPassword,
				InternalApiUser:              c.Config.CCInternalAPIUser,
				InternalApiPassword:          c.Config.CCInternalAPIPassword,
				DbEncryptionKey:              c.Config.DbEncryptionKey,
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
				MinCliVersion:            c.Config.MinCliVersion,
				MinRecommendedCliVersion: c.Config.MinCliVersion,
			},
			Ccdb: &ccnglib.Ccdb{
				Address:  c.Config.MySQLProxyHost(),
				Port:     3306,
				DbScheme: "mysql",
				Roles: []map[string]interface{}{
					map[string]interface{}{
						"tag":      "admin",
						"name":     c.Config.CCDBUsername,
						"password": c.Config.CCDBPassword,
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
					VerificationKey: c.Config.JWTVerificationKey,
				},
				Clients: &ccnglib.Clients{
					CcServiceDashboards: &ccnglib.CcServiceDashboards{
						Scope:  "cloud_controller.write,openid,cloud_controller.read,cloud_controller_service_permissions.read",
						Secret: c.Config.CCServiceDashboardsClientSecret,
					},
					CloudControllerUsernameLookup: &ccnglib.CloudControllerUsernameLookup{
						Secret: c.Config.CloudControllerUsernameLookupClientSecret,
					},
					CcRouting: &ccnglib.CcRouting{
						Secret: c.Config.CCRoutingClientSecret,
					},
				},
			},
			Ssl: &ccnglib.Ssl{
				SkipCertVerify: c.Config.SkipSSLCertVerify,
			},
			LoggerEndpoint: &ccnglib.LoggerEndpoint{
				Port: c.Config.LoggregatorPort,
			},
			Doppler: &ccnglib.Doppler{
				Port: c.Config.LoggregatorPort,
			},
			NfsServer: &ccnglib.NfsServer{
				Address:   c.Config.NFSIP,
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
