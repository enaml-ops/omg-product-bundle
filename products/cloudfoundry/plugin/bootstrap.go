package cloudfoundry

import (
	"github.com/codegangsta/cli"
	"github.com/enaml-ops/enaml"
	bstraplib "github.com/enaml-ops/omg-product-bundle/products/cf-mysql/enaml-gen/bootstrap"
	"github.com/xchapter7x/lo"
)

func NewBootstrapPartition(c *cli.Context, config *Config) InstanceGrouper {
	return &bootstrap{
		Config: config,
		VMType: c.String("bootstrap-vm-type"),
	}
}

func (b *bootstrap) ToInstanceGroup() *enaml.InstanceGroup {
	return &enaml.InstanceGroup{
		Name:      "bootstrap",
		Instances: 1,
		VMType:    b.VMType,
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

func (b *bootstrap) HasValidValues() bool {

	lo.G.Debugf("checking '%s' for valid flags", "bootstrap")

	return b.VMType != ""
}
