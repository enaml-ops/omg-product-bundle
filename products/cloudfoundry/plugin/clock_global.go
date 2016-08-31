package cloudfoundry

import (
	"gopkg.in/yaml.v2"

	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/cloud_controller_clock"
	"github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/cloud_controller_ng"
	"github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/plugin/config"
)

type clockGlobal struct {
	Metron *Metron
	Statsd *StatsdInjector
	Config *config.Config
}

func NewClockGlobalPartition(config *config.Config) InstanceGroupCreator {

	cg := &clockGlobal{
		Config: config,
		Metron: NewMetron(config),
		Statsd: NewStatsdInjector(nil),
	}
	return cg
}

func (c *clockGlobal) ToInstanceGroup() *enaml.InstanceGroup {
	ig := &enaml.InstanceGroup{
		Name:      "clock_global-partition",
		Instances: 1,
		VMType:    c.Config.ClockGlobalVMType,
		AZs:       c.Config.AZs,
		Stemcell:  c.Config.StemcellName,
		Networks: []enaml.Network{
			{Name: c.Config.NetworkName},
		},
		Update: enaml.Update{
			MaxInFlight: 1,
		},
	}

	metronJob := c.Metron.CreateJob()
	nfsJob := CreateNFSMounterJob(c.Config)
	statsdJob := c.Statsd.CreateJob()

	ccw := newCloudControllerNgJob(NewCloudControllerPartition(c.Config).(*CloudControllerPartition))
	props := ccw.Properties.(*cloud_controller_ng.CloudControllerNgJob)

	ig.AddJob(c.newCloudControllerClockJob(props))
	ig.AddJob(&metronJob)
	ig.AddJob(&nfsJob)
	ig.AddJob(&statsdJob)
	return ig
}

func (c *clockGlobal) newCloudControllerClockJob(ccng *cloud_controller_ng.CloudControllerNgJob) *enaml.InstanceJob {
	props := &cloud_controller_clock.CloudControllerClockJob{
		Domain:                   c.Config.SystemDomain,
		SystemDomain:             c.Config.SystemDomain,
		SystemDomainOrganization: "system",
		AppDomains:               c.Config.AppDomains,
		Cc:                       &cloud_controller_clock.Cc{},
		Ccdb: &cloud_controller_clock.Ccdb{
			Address:  c.Config.MySQLProxyHost(),
			Port:     3306,
			DbScheme: "mysql",
			Roles: []map[string]interface{}{
				{
					"name":     c.Config.CCDBUsername,
					"password": c.Config.CCDBPassword,
					"tag":      "admin",
				},
			},
			Databases: []map[string]interface{}{
				map[string]interface{}{
					"citext": true,
					"name":   "ccdb",
					"tag":    "cc",
				},
			},
		},
		Uaa: &cloud_controller_clock.Uaa{
			Url: prefixSystemDomain(c.Config.SystemDomain, "uaa"),
			Jwt: &cloud_controller_clock.Jwt{
				VerificationKey: c.Config.JWTVerificationKey,
			},
			Clients: &cloud_controller_clock.Clients{
				CcServiceDashboards: &cloud_controller_clock.CcServiceDashboards{
					Secret: c.Config.CCServiceDashboardsClientSecret,
				},
			},
		},
		LoggerEndpoint: &cloud_controller_clock.LoggerEndpoint{
			Port: 443,
		},
		Ssl: &cloud_controller_clock.Ssl{
			SkipCertVerify: c.Config.SkipSSLCertVerify,
		},
		Nats: &cloud_controller_clock.Nats{
			User:     c.Config.NATSUser,
			Password: c.Config.NATSPassword,
			Port:     c.Config.NATSPort,
			Machines: c.Config.NATSMachines,
		},
	}

	job := &enaml.InstanceJob{
		Name:       "cloud_controller_clock",
		Release:    CFReleaseName,
		Properties: props,
	}

	ccYaml, _ := yaml.Marshal(ccng.Cc)
	yaml.Unmarshal(ccYaml, props.Cc)

	props.Cc.QuotaDefinitions = map[string]interface{}{
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
	}
	props.Cc.SecurityGroupDefinitions = []map[string]interface{}{
		map[string]interface{}{"name": "all_open",
			"rules": []map[string]interface{}{
				map[string]interface{}{
					"protocol":    "all",
					"destination": "0.0.0.0-255.255.255.255",
				},
			},
		},
	}
	return job
}
