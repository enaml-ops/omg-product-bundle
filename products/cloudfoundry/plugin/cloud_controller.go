package cloudfoundry

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/enaml-ops/enaml"
	ccnglib "github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/cloud_controller_ng"
	"github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/route_registrar"
	"github.com/xchapter7x/lo"
)

func NewCloudControllerPartition(c *cli.Context) InstanceGrouper {
	var proxyIP string
	mysqlProxies := c.StringSlice("mysql-proxy-ip")
	if len(mysqlProxies) > 0 {
		proxyIP = mysqlProxies[0]
	}
	return &CloudControllerPartition{
		AZs:                      c.StringSlice("az"),
		Instances:                c.Int("cc-instances"),
		VMTypeName:               c.String("cc-vm-type"),
		StemcellName:             c.String("stemcell-name"),
		NetworkName:              c.String("network"),
		SystemDomain:             c.String("system-domain"),
		AppDomains:               c.StringSlice("app-domain"),
		AllowAppSSHAccess:        c.Bool("allow-app-ssh-access"),
		Metron:                   NewMetron(c),
		ConsulAgent:              NewConsulAgent(c, []string{}),
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
		SkipSSLCertVerify:        c.BoolT("skip-cert-verify"),
		NATSUser:                 c.String("nats-user"),
		NATSPass:                 c.String("nats-pass"),
		NATSPort:                 c.Int("nats-port"),
		NATSMachines:             c.StringSlice("nats-machine-ip"),
	}
}

//ToInstanceGroup - Convert CLoud Controller Partition to an Instance Group
func (s *CloudControllerPartition) ToInstanceGroup() (ig *enaml.InstanceGroup) {
	ig = &enaml.InstanceGroup{
		Name:      "cloud_controller-partition",
		AZs:       s.AZs,
		Instances: s.Instances,
		VMType:    s.VMTypeName,
		Stemcell:  s.StemcellName,
		Networks: []enaml.Network{
			enaml.Network{Name: s.NetworkName},
		},
		Jobs: []enaml.InstanceJob{
			newCloudControllerNgWorkerJob(s),
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
	return
}

func newCloudControllerNgWorkerJob(c *CloudControllerPartition) enaml.InstanceJob {
	return enaml.InstanceJob{
		Name:    "cloud_controller_ng",
		Release: CFReleaseName,
		Properties: &ccnglib.CloudControllerNgJob{
			AppSsh: &ccnglib.AppSsh{
				HostKeyFingerprint: c.HostKeyFingerprint,
			},
			Domain:                   c.SystemDomain,
			SystemDomain:             c.SystemDomain,
			AppDomains:               c.AppDomains,
			SystemDomainOrganization: "system",
			SupportAddress:           c.SupportAddress,
			Login: &ccnglib.Login{
				Url: fmt.Sprintf("https://login.%s", c.SystemDomain),
			},
			Cc: &ccnglib.Cc{
				AllowedCorsDomains:    []string{fmt.Sprintf("https://login.%s", c.SystemDomain)},
				AllowAppSshAccess:     c.AllowAppSSHAccess,
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
				InstallBuildpacks:            []string{},
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
				Url: fmt.Sprintf("https://uaa.%s", c.SystemDomain),
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
				SkipCertVerify: c.SkipSSLCertVerify,
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
				User:     c.NATSUser,
				Password: c.NATSPass,
				Port:     c.NATSPort,
				Machines: c.NATSMachines,
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
						"uris": []string{fmt.Sprintf("api.%s", c.SystemDomain)},
					},
				},
			},
			Nats: &route_registrar.Nats{
				User:     c.NATSUser,
				Password: c.NATSPass,
				Port:     c.NATSPort,
				Machines: c.NATSMachines,
			},
		},
	}
}

//HasValidValues - Check if valid values has been populated
func (s *CloudControllerPartition) HasValidValues() bool {

	lo.G.Debugf("checking '%s' for valid flags", "cloud controller")

	if len(s.AZs) <= 0 {
		lo.G.Debugf("could not find the correct number of AZs configured '%v' : '%v'", len(s.AZs), s.AZs)
	}

	if s.StemcellName == "" {
		lo.G.Debugf("could not find a valid stemcellname '%v'", s.StemcellName)
	}

	if s.NetworkName == "" {
		lo.G.Debugf("could not find a valid networkname '%v'", s.NetworkName)
	}

	if len(s.AppDomains) <= 0 {
		lo.G.Debugf("could not find the correct number of app domains configured '%v' : '%v'", len(s.AppDomains), s.AppDomains)
	}

	if s.SystemDomain == "" {
		lo.G.Debugf("could not find a valid system domain '%v'", s.SystemDomain)
	}

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

	return (len(s.AZs) > 0 &&
		s.StemcellName != "" &&
		s.VMTypeName != "" &&
		s.Metron.Zone != "" &&
		s.Metron.Secret != "" &&
		s.NetworkName != "" &&
		s.SystemDomain != "" &&
		len(s.AppDomains) > 0 &&
		s.NFSMounter.hasValidValues() &&
		s.ConsulAgent.HasValidValues()) &&
		s.MySQLProxyIP != ""
}
