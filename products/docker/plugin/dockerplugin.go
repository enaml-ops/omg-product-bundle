package docker

import (
	"fmt"
	"io/ioutil"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-product-bundle/products/docker/enaml-gen/containers"
	"github.com/enaml-ops/omg-product-bundle/products/docker/enaml-gen/docker"
	"github.com/enaml-ops/pluginlib/pcli"
	"github.com/enaml-ops/pluginlib/product"
	"github.com/enaml-ops/pluginlib/util"
	"github.com/xchapter7x/lo"
)

const (
	BoshDockerReleaseURL = "https://bosh.io/d/github.com/cf-platform-eng/docker-boshrelease?v=28.0.1"
	BoshDockerReleaseVer = "28.0.1"
	BoshDockerReleaseSHA = "448eaa2f478dc8794933781b478fae02aa44ed6b"
)

type jobBucket struct {
	JobName   string
	JobType   int
	Instances int
}
type Plugin struct {
	PluginVersion      string
	DeploymentName     string
	Containers         interface{}
	NetworkName        string
	IPs                []string
	VMTypeName         string
	DiskTypeName       string
	AZs                []string
	StemcellName       string
	StemcellURL        string
	StemcellVersion    string
	StemcellSHA        string
	RegistryMirrors    []string
	InsecureRegistries []string
}

func (s *Plugin) GetFlags() (flags []pcli.Flag) {
	return []pcli.Flag{
		pcli.Flag{FlagType: pcli.StringFlag, Name: "deployment-name", Value: "docker", Usage: "the name bosh will use for this deployment"},
		pcli.Flag{FlagType: pcli.BoolFlag, Name: "infer-from-cloud", Usage: "setting this flag will attempt to pull as many defaults from your targetted bosh's cloud config as it can (vmtype, network, disk, etc)."},
		pcli.Flag{FlagType: pcli.StringSliceFlag, Name: "ip", Usage: "multiple static ips for each redis leader vm"},
		pcli.Flag{FlagType: pcli.StringSliceFlag, Name: "az", Usage: "list of AZ names to use"},
		pcli.Flag{FlagType: pcli.StringSliceFlag, Name: "insecure-registry", Usage: "Array of insecure registries (no certificate verification for HTTPS and enable HTTP fallback)"},
		pcli.Flag{FlagType: pcli.StringSliceFlag, Name: "registry-mirror", Usage: "Array of preferred Docker registry mirrors"},
		pcli.Flag{FlagType: pcli.StringFlag, Name: "network", Usage: "the name of the network to use"},
		pcli.Flag{FlagType: pcli.StringFlag, Name: "vm-type", Usage: "name of your desired vm type"},
		pcli.Flag{FlagType: pcli.StringFlag, Name: "disk-type", Usage: "name of your desired disk type"},
		pcli.Flag{FlagType: pcli.StringFlag, Name: "stemcell-url", Usage: "the url of the stemcell you wish to use"},
		pcli.Flag{FlagType: pcli.StringFlag, Name: "stemcell-ver", Usage: "the version number of the stemcell you wish to use"},
		pcli.Flag{FlagType: pcli.StringFlag, Name: "stemcell-sha", Usage: "the sha of the stemcell you will use"},
		pcli.Flag{FlagType: pcli.StringFlag, Name: "stemcell-name", Value: "trusty", Usage: "the name of the stemcell you will use"},
		pcli.Flag{FlagType: pcli.StringFlag, Name: "container-definition", Usage: "filepath to the container definition for your docker containers"},
	}
}

func (s *Plugin) GetMeta() product.Meta {
	return product.Meta{
		Name: "docker",
		Properties: map[string]interface{}{
			"version":        s.PluginVersion,
			"docker-release": strings.Join([]string{BoshDockerReleaseURL, BoshDockerReleaseVer, BoshDockerReleaseSHA}, " / "),
		},
	}
}

func (s *Plugin) setContainerDefinitionFromFile(filename string) interface{} {
	var res []interface{}
	if b, e := ioutil.ReadFile(filename); e == nil {
		yaml.Unmarshal(b, &res)

	} else {
		lo.G.Fatalf("you have not given a valid path to a container definition file: %v", filename)
	}
	return res
}

func (s *Plugin) GetProduct(args []string, cloudConfig []byte) (b []byte) {
	var err error
	c := pluginutil.NewContext(args, pluginutil.ToCliFlagArray(s.GetFlags()))
	flgs := s.GetFlags()
	InferFromCloudDecorate(flagsToInferFromCloudConfig, cloudConfig, args, flgs)
	s.Containers = s.setContainerDefinitionFromFile(c.String("container-definition"))
	s.IPs = c.StringSlice("ip")
	s.AZs = c.StringSlice("az")
	s.InsecureRegistries = c.StringSlice("insecure-registry")
	s.RegistryMirrors = c.StringSlice("registry-mirror")
	s.DeploymentName = c.String("deployment-name")
	s.NetworkName = c.String("network")
	s.StemcellName = c.String("stemcell-name")
	s.StemcellVersion = c.String("stemcell-ver")
	s.StemcellSHA = c.String("stemcell-sha")
	s.StemcellURL = c.String("stemcell-url")
	s.VMTypeName = c.String("vm-type")
	s.DiskTypeName = c.String("disk-type")

	if err = s.flagValidation(); err != nil {
		lo.G.Error("invalid arguments: ", err)
		lo.G.Fatal("exiting due to invalid args")
	}

	if err = s.cloudconfigValidation(enaml.NewCloudConfigManifest(cloudConfig)); err != nil {
		lo.G.Error("invalid settings for cloud config on target bosh: ", err)
		lo.G.Fatal("your deployment is not compatible with your cloud config, exiting")
	}
	lo.G.Debug("context", c)
	var dm = new(enaml.DeploymentManifest)
	dm.SetName("docker")
	dm.AddRemoteRelease("docker", BoshDockerReleaseVer, BoshDockerReleaseURL, BoshDockerReleaseSHA)
	dm.AddRemoteStemcell(s.StemcellName, s.StemcellName, s.StemcellVersion, s.StemcellURL, s.StemcellSHA)

	dm.AddInstanceGroup(s.NewDockerInstanceGroup())
	dm.Update = enaml.Update{
		MaxInFlight:     1,
		UpdateWatchTime: "30000-300000",
		CanaryWatchTime: "30000-300000",
		Serial:          false,
		Canaries:        1,
	}
	return dm.Bytes()
}

func (s *Plugin) NewDockerInstanceGroup() (ig *enaml.InstanceGroup) {
	ig = &enaml.InstanceGroup{
		Name:               s.DeploymentName,
		Instances:          len(s.IPs),
		VMType:             s.VMTypeName,
		AZs:                s.AZs,
		Stemcell:           s.StemcellName,
		PersistentDiskType: s.DiskTypeName,
		Jobs: []enaml.InstanceJob{
			s.createDockerJob(),
			s.createContainersJob(),
		},
		Networks: []enaml.Network{
			enaml.Network{Name: s.NetworkName, StaticIPs: s.IPs},
		},
		Update: enaml.Update{
			MaxInFlight: 1,
		},
	}
	return
}

func (s *Plugin) createDockerJob() enaml.InstanceJob {
	return enaml.InstanceJob{
		Name:    "docker",
		Release: "docker",
		Properties: &docker.DockerJob{
			Docker: &docker.Docker{
				RegistryMirrors:    s.RegistryMirrors,
				InsecureRegistries: s.InsecureRegistries,
			},
		},
	}
}

func (s *Plugin) createContainersJob() enaml.InstanceJob {
	return enaml.InstanceJob{
		Name:    "containers",
		Release: "docker",
		Properties: &containers.ContainersJob{
			Containers: s.Containers,
		},
	}
}

func (s *Plugin) cloudconfigValidation(cloudConfig *enaml.CloudConfigManifest) (err error) {
	lo.G.Debug("running cloud config validation")

	for _, vmtype := range cloudConfig.VMTypes {
		err = fmt.Errorf("vm size %s does not exist in cloud config. options are: %v", s.VMTypeName, cloudConfig.VMTypes)
		if vmtype.Name == s.VMTypeName {
			err = nil
			break
		}
	}

	for _, disktype := range cloudConfig.DiskTypes {
		err = fmt.Errorf("disk size %s does not exist in cloud config. options are: %v", s.DiskTypeName, cloudConfig.DiskTypes)
		if disktype.Name == s.DiskTypeName {
			err = nil
			break
		}
	}

	for _, net := range cloudConfig.Networks {
		err = fmt.Errorf("network %s does not exist in cloud config. options are: %v", s.NetworkName, cloudConfig.Networks)
		if net.(map[interface{}]interface{})["name"] == s.NetworkName {
			err = nil
			break
		}
	}

	if len(cloudConfig.VMTypes) == 0 {
		err = fmt.Errorf("no vm sizes found in cloud config")
	}

	if len(cloudConfig.DiskTypes) == 0 {
		err = fmt.Errorf("no disk sizes found in cloud config")
	}

	if len(cloudConfig.Networks) == 0 {
		err = fmt.Errorf("no networks found in cloud config")
	}
	return
}

func (s *Plugin) flagValidation() (err error) {
	lo.G.Debug("validating given flags")

	if len(s.IPs) <= 0 {
		err = fmt.Errorf("no `ip` given")
	}

	if s.Containers == nil {
		err = fmt.Errorf("no valid container definition in given file")
	}

	if len(s.AZs) <= 0 {
		err = fmt.Errorf("no `az` given")
	}

	if s.NetworkName == "" {
		err = fmt.Errorf("no `network-name` given")
	}

	if s.VMTypeName == "" {
		err = fmt.Errorf("no `vm-type` given")
	}
	if s.DiskTypeName == "" {
		err = fmt.Errorf("no `disk-type` given")
	}

	if s.StemcellVersion == "" {
		err = fmt.Errorf("no `stemcell-ver` given")
	}
	return
}

func InferFromCloudDecorate(inferFlagMap map[string][]string, cloudConfig []byte, args []string, flgs []pcli.Flag) {
	c := pluginutil.NewContext(args, pluginutil.ToCliFlagArray(flgs))

	if c.Bool("infer-from-cloud") {
		ccinf := pluginutil.NewCloudConfigInferFromBytes(cloudConfig)
		setAllInferredFlagDefaults(inferFlagMap["disktype"], ccinf.InferDefaultDiskType(), flgs)
		setAllInferredFlagDefaults(inferFlagMap["vmtype"], ccinf.InferDefaultVMType(), flgs)
		setAllInferredFlagDefaults(inferFlagMap["az"], ccinf.InferDefaultAZ(), flgs)
		setAllInferredFlagDefaults(inferFlagMap["network"], ccinf.InferDefaultNetwork(), flgs)
	}
}

func setAllInferredFlagDefaults(matchlist []string, defaultvalue string, flgs []pcli.Flag) {

	for _, match := range matchlist {
		setFlagDefault(match, defaultvalue, flgs)
	}
}

func setFlagDefault(flagname, defaultvalue string, flgs []pcli.Flag) {
	for idx, flg := range flgs {

		if flg.Name == flagname {
			flgs[idx].Value = defaultvalue
		}
	}
}

var flagsToInferFromCloudConfig = map[string][]string{
	"disktype": []string{
		"disk-type",
	},
	"vmtype": []string{
		"vm-type",
	},
	"az":      []string{"az"},
	"network": []string{"network"},
}
