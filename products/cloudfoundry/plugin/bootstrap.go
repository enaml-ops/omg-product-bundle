package cloudfoundry

import (
	"github.com/enaml-ops/enaml"
	bstraplib "github.com/enaml-ops/omg-product-bundle/products/cf-mysql/enaml-gen/bootstrap"
	"github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/plugin/config"
)

type bootstrap struct {
	Config *config.Config
}

func NewBootstrapPartition(config *config.Config) InstanceGroupCreator {
	return &bootstrap{
		Config: config,
	}
}

func (b *bootstrap) ToInstanceGroup() *enaml.InstanceGroup {
	return &enaml.InstanceGroup{
		Name:      "bootstrap",
		Instances: 1,
		VMType:    b.Config.ErrandVMType,
		Lifecycle: "errand",
		AZs:       b.Config.AZs,
		Stemcell:  b.Config.StemcellName,
		Networks: []enaml.Network{
			{Name: b.Config.NetworkName},
		},
		Update: enaml.Update{
			MaxInFlight: 1,
		},
		Jobs: []enaml.InstanceJob{
			{
				Name:    "bootstrap",
				Release: CFMysqlReleaseName,
				Properties: &bstraplib.BootstrapJob{
					ClusterIps:             b.Config.MySQLIPs,
					DatabaseStartupTimeout: 1200,
					BootstrapEndpoint: &bstraplib.BootstrapEndpoint{
						Username: b.Config.MySQLBootstrapUser,
						Password: b.Config.MySQLBootstrapPassword,
					},
				},
			},
		},
	}
}
