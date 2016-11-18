package gemfire_plugin

import (
	"fmt"
	"strings"

	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/pluginlib/cred"
	"github.com/enaml-ops/pluginlib/pcli"
	"github.com/enaml-ops/pluginlib/pluginutil"
	"github.com/enaml-ops/pluginlib/productv1"
)

type Plugin struct {
	Version string `omg:"-"`

	DeploymentName         string
	NetworkName            string
	GemfireReleaseVer      string
	StemcellName           string
	StemcellVer            string
	StemcellAlias          string
	AZs                    []string `omg:"az"`
	LocatorStaticIPs       []string `omg:"locator-static-ip"`
	ServerStaticIPs        []string `omg:"server-static-ip,optional"`
	ServerInstanceCount    int
	GemfireLocatorPort     int
	GemfireLocatorRestPort int
	GemfireServerPort      int
	GemfireLocatorVMMemory int    `omg:"gemfire-locator-vm-memory"`
	GemfireLocatorVMSize   string `omg:"gemfire-locator-vm-size"`
	GemfireServerVMSize    string `omg:"gemfire-server-vm-size"`
	GemfireServerVMMemory  int    `omg:"gemfire-server-vm-memory"`
	ServerDevRestAPIPort   int    `omg:"gemfire-dev-rest-api-port"`
	ServerDevActive        bool   `omg:"gemfire-dev-rest-api-active"`
}

// GetProduct generates a BOSH deployment manifest for p-gemfire.
func (p *Plugin) GetProduct(args []string, cloudConfig []byte, cs cred.Store) ([]byte, error) {
	c := pluginutil.NewContext(args, pluginutil.ToCliFlagArray(p.GetFlags()))
	err := pcli.UnmarshalFlags(p, c)
	if err != nil {
		return nil, err
	}

	deploymentManifest := new(enaml.DeploymentManifest)
	deploymentManifest.SetName(p.DeploymentName)
	deploymentManifest.AddRelease(enaml.Release{Name: releaseName, Version: p.GemfireReleaseVer})
	deploymentManifest.AddStemcell(enaml.Stemcell{
		OS:      p.StemcellName,
		Version: p.StemcellVer,
		Alias:   p.StemcellAlias,
	})
	deploymentManifest.Update = enaml.Update{
		MaxInFlight:     1,
		UpdateWatchTime: "30000-300000",
		CanaryWatchTime: "30000-300000",
		Serial:          false,
		Canaries:        1,
	}

	locator := NewLocatorGroup(p.NetworkName, p.LocatorStaticIPs, p.GemfireLocatorPort, p.GemfireLocatorRestPort, p.GemfireLocatorVMMemory, p.GemfireLocatorVMSize)
	locatorInstanceGroup := locator.GetInstanceGroup()
	locatorInstanceGroup.Stemcell = p.StemcellAlias
	locatorInstanceGroup.AZs = p.AZs
	deploymentManifest.AddInstanceGroup(locatorInstanceGroup)

	server := NewServerGroup(p.NetworkName, p.GemfireServerPort, p.ServerInstanceCount, p.ServerStaticIPs, p.GemfireServerVMSize, p.GemfireServerVMMemory, p.ServerDevRestAPIPort, p.ServerDevActive, locator)
	serverInstanceGroup := server.GetInstanceGroup()
	serverInstanceGroup.Stemcell = p.StemcellAlias
	serverInstanceGroup.AZs = p.AZs
	deploymentManifest.AddInstanceGroup(serverInstanceGroup)
	return deploymentManifest.Bytes(), nil
}

func makeEnvVarName(flagName string) string {
	return "OMG_" + strings.Replace(strings.ToUpper(flagName), "-", "_", -1)
}

// GetMeta returns metadata about the p-gemfire product.
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
