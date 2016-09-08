package cloudfoundry

import (
	"fmt"

	"github.com/enaml-ops/enaml"
	pam "github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/push-apps-manager"
	"github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/plugin/config"
)

type pushAppsManager struct {
	Config *config.Config
}

func NewPushAppsManager(c *config.Config) InstanceGroupCreator {
	return &pushAppsManager{Config: c}
}

func (p *pushAppsManager) ToInstanceGroup() *enaml.InstanceGroup {
	return &enaml.InstanceGroup{
		Name:      "push-apps-manager",
		Instances: 1,
		VMType:    p.Config.ErrandVMType,
		Lifecycle: "errand",
		AZs:       p.Config.AZs,
		Stemcell:  p.Config.StemcellName,
		Networks: []enaml.Network{
			{Name: p.Config.NetworkName},
		},
		Update: enaml.Update{
			MaxInFlight: 1,
		},
		Jobs: []enaml.InstanceJob{
			{
				Name:    "push-apps-manager",
				Release: PushAppsReleaseName,
				Properties: &pam.PushAppsManagerJob{
					Cf: &pam.Cf{
						ApiUrl:        fmt.Sprintf("https://api.%s", p.Config.SystemDomain),
						AdminUsername: "push_apps_manager",
						AdminPassword: p.Config.PushAppsManagerPassword,
						SystemDomain:  p.Config.SystemDomain,
					},
					Services: &pam.Services{
						Authentication: &pam.Authentication{
							CFCLIENTID:       "portal",
							CFCLIENTSECRET:   p.Config.PortalClientSecret,
							CFUAASERVERURL:   fmt.Sprintf("https://uaa.%s", p.Config.SystemDomain),
							CFLOGINSERVERURL: fmt.Sprintf("https://login.%s", p.Config.SystemDomain),
						},
					},
					Env: &pam.Env{
						SecretToken:                  p.Config.AppsManagerSecretToken,
						CfCcApiUrl:                   fmt.Sprintf("https://api.%s", p.Config.SystemDomain),
						CfLoggregatorHttpUrl:         fmt.Sprintf("http://loggregator.%s", p.Config.SystemDomain),
						CfConsoleUrl:                 fmt.Sprintf("https://apps.%s", p.Config.SystemDomain),
						CfNotificationsServiceUrl:    fmt.Sprintf("https://notifications.%s", p.Config.SystemDomain),
						UsageServiceHost:             fmt.Sprintf("https://app-usage.%s", p.Config.SystemDomain),
						BundleWithout:                "test development hosted_only",
						EnableInternalUserStore:      false,
						EnableNonAdminRoleManagement: false,
					},
					Databases: &pam.Databases{
						Console: &pam.Console{
							Ip:       p.Config.MySQLProxyHost(),
							Username: p.Config.ConsoleDBUserName,
							Password: p.Config.ConsoleDBPassword,
							Adapter:  "mysql",
							Port:     3306,
						},
						AppUsageService: &pam.DatabasesAppUsageService{
							Name:     "app_usage_service",
							Ip:       p.Config.MySQLProxyHost(),
							Port:     3306,
							Username: "app_usage",
							Password: p.Config.AppUsageDBPassword,
						},
					},
					Ssl: &pam.Ssl{
						SkipCertVerify: p.Config.SkipSSLCertVerify,
						HttpsOnlyMode:  true,
					},
				},
			},
		},
	}
}
