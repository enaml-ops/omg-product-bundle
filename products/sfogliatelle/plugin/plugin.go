package sfogliatelle

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/pluginlib/pcli"
	"github.com/enaml-ops/pluginlib/pluginutil"
	"github.com/enaml-ops/pluginlib/product"
	"github.com/imdario/mergo"
	"gopkg.in/urfave/cli.v2"
	"gopkg.in/yaml.v2"
)

type Plugin struct {
	Version string
	Source  *os.File
}

// GetProduct generates a BOSH deployment manifest for sfogliatelle.
func (p *Plugin) GetProduct(args []string, cloudConfig []byte) ([]byte, error) {
	var deploymentManifest = enaml.NewDeploymentManifestFromFile(p.Source)
	c := pluginutil.NewContext(args, pluginutil.ToCliFlagArray(p.GetFlags()))

	if err := checkRequiredFields(c); err != nil {
		return nil, err
	}
	layerFile, err := ioutil.ReadFile(c.String("layer-file"))

	if err != nil {
		return nil, err
	}

	if c.IsSet("instance-group-name") && c.IsSet("job-name") {
		p.layerJob(deploymentManifest, layerFile, c.String("instance-group-name"), c.String("job-name"))

	} else if c.IsSet("instance-group-name") && !c.IsSet("job-name") {
		p.layerInstanceGroup(deploymentManifest, layerFile, c.String("instance-group-name"))

	} else {
		layer := enaml.NewDeploymentManifest(layerFile)
		mergo.MergeWithOverwrite(deploymentManifest, layer)
	}
	return deploymentManifest.Bytes(), nil
}

func (p *Plugin) layerInstanceGroup(deploymentManifest *enaml.DeploymentManifest, layerFile []byte, instanceGroupName string) {
	layerGroup := new(enaml.InstanceGroup)
	yaml.Unmarshal(layerFile, layerGroup)
	iGroup := deploymentManifest.GetInstanceGroupByName(instanceGroupName)

	if iGroup == nil {
		deploymentManifest.AddInstanceGroup(layerGroup)
	} else {
		mergo.MergeWithOverwrite(iGroup, layerGroup)
	}
}

func (p *Plugin) layerJob(deploymentManifest *enaml.DeploymentManifest, layerFile []byte, instanceGroupName, jobName string) {
	layerJob := new(enaml.InstanceJob)
	yaml.Unmarshal(layerFile, layerJob)
	iGroup := deploymentManifest.GetInstanceGroupByName(instanceGroupName)

	if iGroup == nil {
		deploymentManifest.AddInstanceGroup(&enaml.InstanceGroup{
			Jobs: []enaml.InstanceJob{
				*layerJob,
			},
		})
	} else {

		for i, _ := range iGroup.Jobs {
			if iGroup.Jobs[i].Name == jobName {
				mergo.MergeWithOverwrite(&(iGroup.Jobs[i]), layerJob)
			}
		}
	}
}

var requiredFlags = []string{"layer-file"}

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
		Name: "sfogliatelle",
		Properties: map[string]interface{}{
			"version":     p.Version,
			"description": "This plugin is meant to facilitate a easy transition from the world of yaml to the world of enaml. it will overlay or if not called in a chain, just use the given manifest after unmarshalling it into an validated object structure",
		},
	}
}

// GetFlags returns the CLI flags accepted by the plugin.
func (p *Plugin) GetFlags() []pcli.Flag {
	return []pcli.Flag{
		pcli.Flag{
			FlagType: pcli.StringFlag,
			Name:     "layer-file",
			Usage:    "the path to the yaml overlay file",
		},
		pcli.Flag{
			FlagType: pcli.StringFlag,
			Name:     "job-name",
			Usage:    "if given a jobname, it will overlay the given yaml as a job object only",
		},
		pcli.Flag{
			FlagType: pcli.StringFlag,
			Name:     "instance-group-name",
			Usage:    "if given a groupname, it will overlay the given yaml as that intancegroup object",
		},
	}
}
