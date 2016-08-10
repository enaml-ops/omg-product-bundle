package plugin

import (
	"github.com/codegangsta/cli"
	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/pluginlib/pcli"
	"github.com/enaml-ops/pluginlib/product"
	"github.com/enaml-ops/pluginlib/util"
	"github.com/xchapter7x/lo"
)

type Plugin struct {
}

func (p *Plugin) GetFlags() (flags []pcli.Flag) {
	return nil
}

func (p *Plugin) GetMeta() product.Meta {
	return product.Meta{
		Name: "docker-registry",
	}
}

func (p *Plugin) GetProduct(args []string, cloudConfig []byte) (b []byte) {
	if len(cloudConfig) == 0 {
		lo.G.Debug("plugin: empty cloud config")
		panic("cloud config cannot be empty")
	}
	c := pluginutil.NewContext(args, pluginutil.ToCliFlagArray(p.GetFlags()))
	dm := NewDeploymentManifest(c, cloudConfig)
	return dm.Bytes()
}

func NewDeploymentManifest(c *cli.Context, cloudConfig []byte) enaml.DeploymentManifest {
	return enaml.DeploymentManifest{}
}
