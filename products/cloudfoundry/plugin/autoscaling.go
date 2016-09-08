package cloudfoundry

import (
	"fmt"

	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-cli/utils"
	das "github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/deploy-autoscaling"
	db "github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/destroy-broker"
	rb "github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/register-broker"
	ta "github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/test-autoscaling"

	"github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/plugin/config"
)

type (
	deployAutoscaling       struct{ *config.Config }
	registerAutoscaleBroker struct{ *config.Config }
	destroyAutoscaleBroker  struct{ *config.Config }
	autoscalingTests        struct{ *config.Config }
)

func NewDeployAutoscaling(c *config.Config) InstanceGroupCreator {
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

func NewAutoscaleDestroyBroker(c *config.Config) InstanceGroupCreator {
	return destroyAutoscaleBroker{c}
}

func (d destroyAutoscaleBroker) ToInstanceGroup() *enaml.InstanceGroup {
	return &enaml.InstanceGroup{
		Name:      "autoscaling-destroy-broker",
		Instances: 1,
		VMType:    d.ErrandVMType,
		Lifecycle: "errand",
		AZs:       d.AZs,
		Stemcell:  d.StemcellName,
		Networks: []enaml.Network{
			{Name: d.NetworkName},
		},
		Update: enaml.Update{
			MaxInFlight: 1,
		},
		Jobs: []enaml.InstanceJob{
			{
				Name:    "destroy-broker",
				Release: CFAutoscalingReleaseName,
				Properties: &db.DestroyBrokerJob{
					Autoscale: &db.Autoscale{
						Broker: &db.Broker{
							User:     d.AutoscaleBrokerUser,
							Password: d.AutoscaleBrokerPassword,
						},
						Cf: &db.Cf{
							AdminUser:     "admin",
							AdminPassword: d.AdminPassword,
						},
						Organization: "system",
						Space:        "autoscaling",
					},
					Domain: d.SystemDomain,
					Ssl: &db.Ssl{
						SkipCertVerify: d.SkipSSLCertVerify,
					},
				},
			},
		},
	}
}

func NewAutoscalingTests(c *config.Config) InstanceGroupCreator {
	return autoscalingTests{c}
}

func (a autoscalingTests) ToInstanceGroup() *enaml.InstanceGroup {
	return &enaml.InstanceGroup{
		Name:      "autoscaling-tests",
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
				Name:    "test-autoscaling",
				Release: CFAutoscalingReleaseName,
				Properties: &ta.TestAutoscalingJob{
					Autoscale: &ta.Autoscale{
						Cf: &ta.Cf{
							AdminUser:     "admin",
							AdminPassword: a.AdminPassword,
						},
					},
					Domain: a.SystemDomain,
					Ssl: &ta.Ssl{
						SkipCertVerify: a.SkipSSLCertVerify,
					},
				},
			},
		},
	}
}
