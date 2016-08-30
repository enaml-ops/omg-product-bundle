package cloudfoundry

import (
	"github.com/codegangsta/cli"
	"github.com/enaml-ops/enaml"
	ccworkerlib "github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/cloud_controller_worker"
	"github.com/xchapter7x/lo"
)

//CloudControllerWorkerPartition - Cloud Controller Worker Partition
type CloudControllerWorkerPartition struct {
	Config                *Config
	Instances             int
	VMTypeName            string
	AllowedCorsDomains    []string
	Metron                *Metron
	ConsulAgent           *ConsulAgent
	StatsdInjector        *StatsdInjector
	NFSMounter            *NFSMounter
	StagingUploadUser     string
	StagingUploadPassword string
	BulkAPIUser           string
	BulkAPIPassword       string
	DbEncryptionKey       string
	InternalAPIUser       string
	InternalAPIPassword   string
	CCDBUsername          string
	CCDBPassword          string
	MySQLProxyIP          string
}

//NewCloudControllerWorkerPartition - Creating a New Cloud Controller Partition
func NewCloudControllerWorkerPartition(c *cli.Context, config *Config) InstanceGrouper {
	var proxyIP string
	mysqlProxies := c.StringSlice("mysql-proxy-ip")
	if len(mysqlProxies) > 0 {
		proxyIP = mysqlProxies[0]
	}

	return &CloudControllerWorkerPartition{
		Config:                config,
		Instances:             c.Int("cc-worker-instances"),
		VMTypeName:            c.String("cc-worker-vm-type"),
		Metron:                NewMetron(c),
		ConsulAgent:           NewConsulAgent(c, []string{}, config),
		NFSMounter:            NewNFSMounter(c),
		StatsdInjector:        NewStatsdInjector(c),
		StagingUploadUser:     c.String("cc-staging-upload-user"),
		StagingUploadPassword: c.String("cc-staging-upload-password"),
		BulkAPIUser:           c.String("cc-bulk-api-user"),
		BulkAPIPassword:       c.String("cc-bulk-api-password"),
		InternalAPIUser:       c.String("cc-internal-api-user"),
		InternalAPIPassword:   c.String("cc-internal-api-password"),
		DbEncryptionKey:       c.String("cc-db-encryption-key"),
		CCDBUsername:          c.String("db-ccdb-username"),
		CCDBPassword:          c.String("db-ccdb-password"),
		MySQLProxyIP:          proxyIP,
	}
}

//ToInstanceGroup - Convert CLoud Controller Partition to an Instance Group
func (s *CloudControllerWorkerPartition) ToInstanceGroup() (ig *enaml.InstanceGroup) {
	ig = &enaml.InstanceGroup{
		Name:      "cloud_controller_worker-partition",
		AZs:       s.Config.AZs,
		Instances: s.Instances,
		VMType:    s.VMTypeName,
		Stemcell:  s.Config.StemcellName,
		Networks: []enaml.Network{
			enaml.Network{Name: s.Config.NetworkName},
		},
		Jobs: []enaml.InstanceJob{
			newCloudControllerWorkerJob(s),
			s.ConsulAgent.CreateJob(),
			s.NFSMounter.CreateJob(),
			s.Metron.CreateJob(),
			s.StatsdInjector.CreateJob(),
		},
		Update: enaml.Update{
			MaxInFlight: 1,
			Serial:      true,
		},
	}
	return
}

func newCloudControllerWorkerJob(c *CloudControllerWorkerPartition) enaml.InstanceJob {
	return enaml.InstanceJob{
		Name:    "cloud_controller_worker",
		Release: CFReleaseName,
		Properties: &ccworkerlib.CloudControllerWorkerJob{
			Domain:                   c.Config.SystemDomain,
			SystemDomain:             c.Config.SystemDomain,
			AppDomains:               c.Config.AppDomains,
			SystemDomainOrganization: "system",
			Cc: &ccworkerlib.Cc{
				AllowAppSshAccess: c.Config.AllowSSHAccess,
				Buildpacks: &ccworkerlib.Buildpacks{
					BlobstoreType: "fog",
					FogConnection: &ccworkerlib.DefaultFogConnection{
						Provider:  "Local",
						LocalRoot: "/var/vcap/nfs/shared",
					},
				},
				Droplets: &ccworkerlib.Droplets{
					BlobstoreType: "fog",
					FogConnection: &ccworkerlib.DefaultFogConnection{
						Provider:  "Local",
						LocalRoot: "/var/vcap/nfs/shared",
					},
				},
				Packages: &ccworkerlib.Packages{
					BlobstoreType: "fog",
					FogConnection: &ccworkerlib.DefaultFogConnection{
						Provider:  "Local",
						LocalRoot: "/var/vcap/nfs/shared",
					},
				},
				ResourcePool: &ccworkerlib.ResourcePool{
					BlobstoreType: "fog",
					FogConnection: &ccworkerlib.DefaultFogConnection{
						Provider:  "Local",
						LocalRoot: "/var/vcap/nfs/shared",
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
						"non_basic_services_allowed": true,
						"total_routes":               1000,
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
				DefaultRunningSecurityGroups: []string{"all_open"},
				DefaultStagingSecurityGroups: []string{"all_open"},
				LoggingLevel:                 "debug",
				MaximumHealthCheckTimeout:    "600",
				StagingUploadUser:            c.StagingUploadUser,
				StagingUploadPassword:        c.StagingUploadPassword,
				BulkApiUser:                  c.BulkAPIUser,
				BulkApiPassword:              c.BulkAPIPassword,
				InternalApiUser:              c.InternalAPIUser,
				InternalApiPassword:          c.InternalAPIPassword,
				DbEncryptionKey:              c.DbEncryptionKey,
			},
			Ccdb: &ccworkerlib.Ccdb{
				Address: c.MySQLProxyIP,
				Databases: []map[string]interface{}{
					map[string]interface{}{
						"citext": true,
						"name":   "ccdb",
						"tag":    "cc",
					},
				},
				DbScheme: "mysql",
				Port:     3306,
				Roles: []map[string]interface{}{
					{
						"name":     c.CCDBUsername,
						"password": c.CCDBPassword,
						"tag":      "admin",
					},
				},
			},
			Nats: &ccworkerlib.Nats{
				User:     c.Config.NATSUser,
				Port:     c.Config.NATSPort,
				Password: c.Config.NATSPassword,
				Machines: c.Config.NATSMachines,
			},
		},
	}
}

//HasValidValues - Check if valid values has been populated
func (s *CloudControllerWorkerPartition) HasValidValues() bool {
	lo.G.Debugf("checking '%s' for valid flags", "cloud controller worker")

	if s.VMTypeName == "" {
		lo.G.Debugf("could not find a valid vmtypename '%v'", s.VMTypeName)
	}

	if s.Metron.Zone == "" {
		lo.G.Debugf("could not find a valid metron zone '%v'", s.Metron.Zone)
	}

	if s.Metron.Secret == "" {
		lo.G.Debugf("could not find a valid metron secret '%v'", s.Metron.Secret)
	}
	return (s.VMTypeName != "" &&
		s.Metron.Zone != "" &&
		s.Metron.Secret != "" &&
		s.NFSMounter.hasValidValues() &&
		s.ConsulAgent.HasValidValues())
}
