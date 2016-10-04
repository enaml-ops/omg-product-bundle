package gemfire_plugin

import (
	"fmt"

	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/pluginlib/pcli"
	"github.com/enaml-ops/pluginlib/pluginutil"
	"github.com/enaml-ops/pluginlib/product"
	"gopkg.in/urfave/cli.v2"
)

type Plugin struct {
	Version string
}

// GetProduct generates a BOSH deployment manifest for p-rabbitmq.
func (p *Plugin) GetProduct(args []string, cloudConfig []byte) ([]byte, error) {
	var deploymentManifest = new(enaml.DeploymentManifest)
	c := pluginutil.NewContext(args, pluginutil.ToCliFlagArray(p.GetFlags()))
	if err := checkRequiredFields(c); err != nil {
		return nil, err
	}
	networkname := c.String("network-name")
	staticips := c.StringSlice("locator-static-ip")
	locatorInstanceGroup := NewLocatorGroup(networkname, staticips)
	deploymentManifest.AddInstanceGroup(locatorInstanceGroup.GetInstanceGroup())
	return deploymentManifest.Bytes(), nil
}

var requiredFlags = []string{
	"az",
	"network-name",
	"locator-static-ip",
	"server-instance-count",
}

func checkRequiredFields(c *cli.Context) error {
	for _, flagname := range requiredFlags {
		err := validate(flagname, c)
		if err != nil {
			return err
		}
	}
	return nil
}

func validate(flagName string, c *cli.Context) error {
	if !c.IsSet("az") {
		return fmt.Errorf("error: sorry you need to give me an AZ")
	}
	return nil
}

// GetMeta returns metadata about the p-rabbitmq product.
func (p *Plugin) GetMeta() product.Meta {
	return product.Meta{
		Name:       "p-gemfire",
		Properties: map[string]interface{}{},
	}
}

// GetFlags returns the CLI flags accepted by the plugin.
func (p *Plugin) GetFlags() []pcli.Flag {
	return []pcli.Flag{
		pcli.Flag{
			FlagType: pcli.StringFlag,
			Name:     "deployment-name",
			Value:    "p-gemfire",
			Usage:    "the name bosh will use for this deployment",
		},
		pcli.Flag{
			FlagType: pcli.StringSliceFlag,
			Name:     "az",
			Usage:    "the list of Availability Zones where you wish to deploy gemfire",
		},
		pcli.Flag{
			FlagType: pcli.StringFlag,
			Name:     "network-name",
			Usage:    "the network where you wish to deploy locators and servers",
		},
		pcli.Flag{
			FlagType: pcli.StringSliceFlag,
			Name:     "locator-static-ip",
			Usage:    "static IPs to assign to locator VMs",
		},
		pcli.Flag{
			FlagType: pcli.IntFlag,
			Name:     "server-instance-count",
			Usage:    "the number of server instances you wish to deploy",
		},
	}
}
