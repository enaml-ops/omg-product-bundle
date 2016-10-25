package cloudfoundry

import (
	"github.com/enaml-ops/enaml"
	ccworkerlib "github.com/enaml-ops/omg-product-bundle/products/oss_cf/enaml-gen/cloud_controller_worker"
	"github.com/enaml-ops/omg-product-bundle/products/oss_cf/plugin/config"
)

//CloudControllerWorkerPartition - Cloud Controller Worker Partition
type CloudControllerWorkerPartition struct {
	Config      *config.Config
	Metron      *Metron
	ConsulAgent *ConsulAgent
}

//NewCloudControllerWorkerPartition - Creating a New Cloud Controller Partition
func NewCloudControllerWorkerPartition(config *config.Config) InstanceGroupCreator {

	return &CloudControllerWorkerPartition{
		Config:      config,
		Metron:      NewMetron(config),
		ConsulAgent: NewConsulAgent([]string{}, config),
	}
}

//ToInstanceGroup - Convert CLoud Controller Partition to an Instance Group
func (s *CloudControllerWorkerPartition) ToInstanceGroup() (ig *enaml.InstanceGroup) {
	ig = &enaml.InstanceGroup{
		Name:      "cloud_controller_worker-partition",
		AZs:       s.Config.AZs,
		Instances: s.Config.CloudControllerWorkerInstances,
		VMType:    s.Config.CloudControllerWorkerVMType,
		Stemcell:  s.Config.StemcellName,
		Networks: []enaml.Network{
			enaml.Network{Name: s.Config.NetworkName},
		},
		Jobs: []enaml.InstanceJob{
			newCloudControllerWorkerJob(s),
			s.ConsulAgent.CreateJob(),
			CreateNFSMounterJob(s.Config),
			s.Metron.CreateJob(),
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
			Ssl: &ccworkerlib.Ssl{
				SkipCertVerify: c.Config.SkipSSLCertVerify,
			},
			Cc: &ccworkerlib.Cc{
				AllowAppSshAccess: c.Config.AllowSSHAccess,
				Buildpacks: &ccworkerlib.Buildpacks{
					BlobstoreType: "fog",
					FogConnection: map[string]string{
						"provider":   "Local",
						"local_root": "/var/vcap/nfs/shared",
					},
				},
				Droplets: &ccworkerlib.Droplets{
					BlobstoreType: "fog",
					FogConnection: map[string]string{
						"provider":   "Local",
						"local_root": "/var/vcap/nfs/shared",
					},
				},
				Packages: &ccworkerlib.Packages{
					BlobstoreType: "fog",
					FogConnection: map[string]string{
						"provider":   "Local",
						"local_root": "/var/vcap/nfs/shared",
					},
				},
				ResourcePool: &ccworkerlib.ResourcePool{
					BlobstoreType: "fog",
					FogConnection: map[string]string{
						"provider":   "Local",
						"local_root": "/var/vcap/nfs/shared",
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
				StagingUploadUser:            c.Config.StagingUploadUser,
				StagingUploadPassword:        c.Config.StagingUploadPassword,
				BulkApiPassword:              c.Config.CCBulkAPIPassword,
				InternalApiUser:              c.Config.CCInternalAPIUser,
				InternalApiPassword:          c.Config.CCInternalAPIPassword,
				DbEncryptionKey:              c.Config.DbEncryptionKey,
			},
			Ccdb: &ccworkerlib.Ccdb{
				Address: c.Config.MySQLProxyHost(),
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
						"name":     c.Config.CCDBUsername,
						"password": c.Config.CCDBPassword,
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
