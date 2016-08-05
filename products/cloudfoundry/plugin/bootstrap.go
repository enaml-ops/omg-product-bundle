package cloudfoundry

import (
	"github.com/codegangsta/cli"
	"github.com/enaml-ops/enaml"
	bstraplib "github.com/enaml-ops/omg-product-bundle/products/cf-mysql/enaml-gen/bootstrap"
	"github.com/xchapter7x/lo"
)

func NewBootstrapPartition(c *cli.Context) InstanceGrouper {
	return &bootstrap{
		AZs:           c.StringSlice("az"),
		StemcellName:  c.String("stemcell-name"),
		NetworkName:   c.String("network"),
		MySQLIPs:      c.StringSlice("mysql-ip"),
		MySQLUser:     c.String("mysql-bootstrap-username"),
		MySQLPassword: c.String("mysql-bootstrap-password"),
	}
}

func (b *bootstrap) ToInstanceGroup() *enaml.InstanceGroup {
	return &enaml.InstanceGroup{
		Name:      "bootstrap",
		Instances: 1,
		VMType:    "errand",
		Lifecycle: "errand",
		AZs:       b.AZs,
		Stemcell:  b.StemcellName,
		Networks: []enaml.Network{
			{Name: b.NetworkName},
		},
		Update: enaml.Update{
			MaxInFlight: 1,
		},
		Jobs: []enaml.InstanceJob{
			{
				Name:    "bootstrap",
				Release: CFMysqlReleaseName,
				Properties: &bstraplib.Bootstrap{
					ClusterIps:             b.MySQLIPs,
					DatabaseStartupTimeout: 1200,
					BootstrapEndpoint: &bstraplib.BootstrapEndpoint{
						Username: b.MySQLUser,
						Password: b.MySQLPassword,
					},
				},
			},
		},
	}
}

func (b *bootstrap) HasValidValues() bool {

	lo.G.Debugf("checking '%s' for valid flags", "bootstrap")

	if len(b.AZs) <= 0 {
		lo.G.Debugf("could not find the correct number of AZs configured '%v' : '%v'", len(b.AZs), b.AZs)
	}

	if b.StemcellName == "" {
		lo.G.Debugf("could not find a valid stemcellname '%v'", b.StemcellName)
	}

	if b.NetworkName == "" {
		lo.G.Debugf("could not find a valid networkname '%v'", b.NetworkName)
	}

	if len(b.MySQLIPs) <= 0 {
		lo.G.Debugf("could not find the correct number of mysql ips '%v' : '%v'", len(b.MySQLIPs), b.MySQLIPs)
	}

	if b.MySQLUser == "" {
		lo.G.Debugf("could not find a valid mysql user '%v'", b.MySQLUser)
	}

	if b.MySQLPassword == "" {
		lo.G.Debugf("could not find a valid admin password '%v'", b.MySQLPassword)
	}

	return len(b.AZs) > 0 &&
		b.StemcellName != "" &&
		b.NetworkName != "" &&
		len(b.MySQLIPs) > 0 &&
		b.MySQLUser != "" &&
		b.MySQLPassword != ""
}
