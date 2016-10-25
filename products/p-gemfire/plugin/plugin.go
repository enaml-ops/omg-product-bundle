package gemfire_plugin

import (
	"fmt"
	"os"
	"strings"

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
	deploymentManifest.SetName(c.String("deployment-name"))
	deploymentManifest.AddRelease(enaml.Release{Name: releaseName, Version: c.String("gemfire-release-ver")})
	deploymentManifest.AddStemcell(enaml.Stemcell{
		OS:      c.String("stemcell-name"),
		Version: c.String("stemcell-ver"),
		Alias:   c.String("stemcell-alias"),
	})
	deploymentManifest.Update = enaml.Update{
		MaxInFlight:     1,
		UpdateWatchTime: "30000-300000",
		CanaryWatchTime: "30000-300000",
		Serial:          false,
		Canaries:        1,
	}

	azs := c.StringSlice("az")
	networkname := c.String("network-name")
	locatorstaticips := c.StringSlice("locator-static-ip")
	locatorport := c.Int("gemfire-locator-port")
	locatorrestport := c.Int("gemfire-locator-rest-port")
	locatorvmmemory := c.Int("gemfire-locator-vm-memory")
	locatorvmtype := c.String("gemfire-locator-vm-size")
	locator := NewLocatorGroup(networkname, locatorstaticips, locatorport, locatorrestport, locatorvmmemory, locatorvmtype)
	locatorInstanceGroup := locator.GetInstanceGroup()
	locatorInstanceGroup.Stemcell = c.String("stemcell-alias")
	locatorInstanceGroup.AZs = azs
	deploymentManifest.AddInstanceGroup(locatorInstanceGroup)

	serverport := c.Int("gemfire-server-port")
	servervmtype := c.String("gemfire-server-vm-size")
	serverInstanceCount := c.Int("server-instance-count")
	serverStaticIPs := c.StringSlice("server-static-ip")
	servervmmemory := c.Int("gemfire-server-vm-memory")
	serverdevrestport := c.Int("gemfire-dev-rest-api-port")
	serverdevrestactive := c.Bool("gemfire-dev-rest-api-active")
	server := NewServerGroup(networkname, serverport, serverInstanceCount, serverStaticIPs, servervmtype, servervmmemory, serverdevrestport, serverdevrestactive, locator)
	serverInstanceGroup := server.GetInstanceGroup()
	serverInstanceGroup.Stemcell = c.String("stemcell-alias")
	serverInstanceGroup.AZs = azs
	deploymentManifest.AddInstanceGroup(serverInstanceGroup)
	return deploymentManifest.Bytes(), nil
}

var requiredFlags = []string{
	"az",
	"network-name",
	"locator-static-ip",
	"gemfire-locator-vm-size",
	"gemfire-server-vm-size",
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

func makeEnvVarName(flagName string) string {
	return "OMG_" + strings.Replace(strings.ToUpper(flagName), "-", "_", -1)
}

func validate(flagName string, c *cli.Context) error {

	if c.IsSet(flagName) || os.Getenv(makeEnvVarName(flagName)) != "" {
		return nil
	}
	return fmt.Errorf("error: sorry you need to give me an `--%s`", flagName)
}

// GetMeta returns metadata about the p-rabbitmq product.
func (p *Plugin) GetMeta() product.Meta {
	return product.Meta{
		Name: "p-gemfire",
		Stemcell: enaml.Stemcell{
			Name:    defaultStemcellName,
			Alias:   defaultStemcellAlias,
			Version: defaultStemcellVersion,
		},
		Releases: []enaml.Release{
			enaml.Release{
				Name:    releaseName,
				Version: releaseVersion,
			},
		},
		Properties: map[string]interface{}{
			"version":              p.Version,
			"stemcell":             defaultStemcellVersion,
			"pivotal-gemfire-tile": "NOT COMPATIBLE WITH TILE RELEASES",
			"p-gemfire":            fmt.Sprintf("%s / %s", releaseName, releaseVersion),
			"description":          "this plugin is designed to work with a special p-gemfire release",
		},
	}
}

// GetFlags returns the CLI flags accepted by the plugin.
func (p *Plugin) GetFlags() []pcli.Flag {
	return []pcli.Flag{
		pcli.Flag{
			FlagType: pcli.StringFlag,
			Name:     "deployment-name",
			Value:    defaultDeploymentName,
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
			FlagType: pcli.StringSliceFlag,
			Name:     "server-static-ip",
			Usage:    "static IPs to assign to server VMs - this is optional, if non given bosh will assign IPs and create instances based on the InstanceCount flag value",
		},
		pcli.Flag{
			FlagType: pcli.IntFlag,
			Name:     "server-instance-count",
			Value:    defaultServerInstanceCount,
			Usage:    "the number of server instances you wish to deploy - if static ips are given this will be ignored",
		},
		pcli.Flag{
			FlagType: pcli.IntFlag,
			Name:     "gemfire-server-port",
			Value:    defaultServerPort,
			Usage:    "the port gemfire servers will listen on",
		},
		pcli.Flag{
			FlagType: pcli.IntFlag,
			Name:     "gemfire-locator-port",
			Value:    defaultLocatorPort,
			Usage:    "the port gemfire locators will listen on",
		},
		pcli.Flag{
			FlagType: pcli.IntFlag,
			Name:     "gemfire-dev-rest-api-port",
			Value:    defaultDevRestPort,
			Usage:    "this will set the port the dev rest api listens on, if active",
		},
		pcli.Flag{
			FlagType: pcli.BoolFlag,
			Name:     "gemfire-dev-rest-api-active",
			Value:    defaultDevRestActive,
			Usage:    "set to true to activate the dev rest api on server nodes",
		},
		pcli.Flag{
			FlagType: pcli.IntFlag,
			Name:     "gemfire-locator-vm-memory",
			Value:    defaultLocatorVMMemory,
			Usage:    "the amount of memory allocated by the locator process",
		},
		pcli.Flag{
			FlagType: pcli.IntFlag,
			Name:     "gemfire-server-vm-memory",
			Value:    defaultLocatorVMMemory,
			Usage:    "the amount of memory allocated by the server process",
		},
		pcli.Flag{
			FlagType: pcli.IntFlag,
			Name:     "gemfire-locator-rest-port",
			Value:    defaultLocatorRestPort,
			Usage:    "the port gemfire locators rest service will listen on",
		},
		pcli.Flag{
			FlagType: pcli.StringFlag,
			Name:     "gemfire-locator-vm-size",
			Usage:    "the vm size of gemfire locators",
		},
		pcli.Flag{
			FlagType: pcli.StringFlag,
			Name:     "gemfire-server-vm-size",
			Usage:    "the vm size of gemfire servers",
		},
		pcli.Flag{
			FlagType: pcli.StringFlag,
			Name:     "stemcell-name",
			Value:    p.GetMeta().Stemcell.Name,
			Usage:    "the name of the stemcell you with to use",
		},
		pcli.Flag{
			FlagType: pcli.StringFlag,
			Name:     "stemcell-alias",
			Value:    p.GetMeta().Stemcell.Alias,
			Usage:    "the name of the stemcell you with to use",
		},
		pcli.Flag{
			FlagType: pcli.StringFlag,
			Name:     "stemcell-ver",
			Value:    p.GetMeta().Stemcell.Version,
			Usage:    "the name of the stemcell you with to use",
		},
		pcli.Flag{
			FlagType: pcli.StringFlag,
			Name:     "gemfire-release-ver",
			Value:    releaseVersion,
			Usage:    "the version of the release to use for the deployment",
		},
	}
}
