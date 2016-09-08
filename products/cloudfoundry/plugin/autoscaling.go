package cloudfoundry

import (
	"fmt"

	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-cli/utils"
	das "github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/deploy-autoscaling"
	rb "github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/register-broker"

	"github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/plugin/config"
)

type (
	deployAutoscaling       struct{ *config.Config }
	registerAutoscaleBroker struct{ *config.Config }
)

func NewDepoyAutoscaling(c *config.Config) InstanceGroupCreator {
	return deployAutoscaling{c}
}

func (a deployAutoscaling) ToInstanceGroup() *enaml.InstanceGroup {
	return &enaml.InstanceGroup{
		Name:      "autoscaling",
		Instances: 1,
		VMType:    a.ErrandVMType,
		Lifecycle: "errand",
		AZs:       a.AZs,
		Stemcell:  a.StemcellName,
		Networks: []enaml.Network{
			{Name: a.NetworkName},
		},
		Update: enaml.Update{
			MaxInFlight: 1,
		},
		Jobs: []enaml.InstanceJob{
			{
				Name:    "deploy-autoscaling",
				Release: CFAutoscalingReleaseName,
				Properties: &das.DeployAutoscalingJob{
					Autoscale: &das.Autoscale{
						Broker: &das.Broker{
							User:     a.AutoscaleBrokerUser,
							Password: a.AutoscaleBrokerPassword,
						},
						Cf: &das.Cf{
							AdminUser:     "admin",
							AdminPassword: a.AdminPassword,
						},
						InstanceCount: 1,
						Database: &das.Database{
							Url: fmt.Sprintf("mysql://%s:%s@%s:3306/autoscale", a.AutoscaleDBUser, a.AutoscaleDBPassword, a.MySQLProxyHost()),
						},
						EncryptionKey:     utils.NewPassword(16),
						EnableDiego:       true,
						NotificationsHost: fmt.Sprintf("https://notifications.%s", a.SystemDomain),
						Organization:      "system",
						Space:             "autoscaling",
					},
					Domain: a.SystemDomain,
					Ssl: &das.Ssl{
						SkipCertVerify: a.SkipSSLCertVerify,
					},
					Uaa: &das.Uaa{
						Clients: &das.Clients{
							AutoscalingService: &das.AutoscalingService{
								Secret: a.AutoScalingServiceClientSecret,
							},
						},
					},
				},
			},
		},
	}
}

func NewAutoscaleRegisterBroker(c *config.Config) InstanceGroupCreator {
	return registerAutoscaleBroker{c}
}

func (a registerAutoscaleBroker) ToInstanceGroup() *enaml.InstanceGroup {
	return &enaml.InstanceGroup{
		Name:      "autoscaling-register-broker",
		Instances: 1,
		VMType:    a.ErrandVMType,
		Lifecycle: "errand",
		AZs:       a.AZs,
		Stemcell:  a.StemcellName,
		Networks: []enaml.Network{
			{Name: a.NetworkName},
		},
		Update: enaml.Update{
			MaxInFlight: 1,
		},
		Jobs: []enaml.InstanceJob{
			{
				Name:    "register-broker",
				Release: CFAutoscalingReleaseName,
				Properties: &rb.RegisterBrokerJob{
					AppDomains: a.AppDomains,
					Autoscale: &rb.Autoscale{
						Broker: &rb.Broker{
							User:     a.AutoscaleBrokerUser,
							Password: a.AutoscaleBrokerPassword,
						},
						Cf: &rb.Cf{
							AdminUser:     "admin",
							AdminPassword: a.AdminPassword,
						},
						Organization: "system",
						Space:        "autoscaling",
					},
					Domain: a.SystemDomain,
					Ssl: &rb.Ssl{
						SkipCertVerify: a.SkipSSLCertVerify,
					},
				},
			},
		},
	}
}
