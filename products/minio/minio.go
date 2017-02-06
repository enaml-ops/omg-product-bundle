package minio

import (
	"fmt"

	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-product-bundle/products/minio/enaml-gen/minio-server"
	"github.com/xchapter7x/lo"
	"gopkg.in/yaml.v2"
)

var (
	DefaultMinioReleaseURL = "https://github.com/Pivotal-Field-Engineering/minio-release/releases/download/v1/minio-1.tgz"
	DefaultMinioReleaseSHA = "0629ca0118749c539c9dc4ac457af411198c274a"
	DefaultMinioReleaseVer = "1"
	DefaultStemcellAlias   = "trusty"
	DefaultStemcellName    = "ubuntu-trusty"
	DefaultStemcellVersion = "latest"
	DefaultRegion          = "us-east-1"
)

// Config contains the configuration for a Minio deployment.
type Config struct {
	AZ                  string `omg:"az"`
	NetworkName         string `omg:"network-name"`
	IP                  string `omg:"ip"`
	VMType              string `omg:"vm-type"`
	DiskType            string `omg:"disk-type"`
	MinioReleaseURL     string `omg:"minio-release-url"`
	MinioReleaseSHA     string `omg:"minio-release-sha"`
	MinioReleaseVersion string `omg:"minio-release-ver"`
	StemcellAlias       string
	StemcellOS          string `omg:"stemcell-os"`
	StemcellVersion     string
	Region              string `omg:"region"`
	AccessKey           string `omg:"access-key"`
	SecretKey           string `omg:"secret-key"`
}

type Deployment struct {
	config Config
}

func NewConfig() Config {
	return Config{
		MinioReleaseURL:     DefaultMinioReleaseURL,
		MinioReleaseSHA:     DefaultMinioReleaseSHA,
		MinioReleaseVersion: DefaultMinioReleaseVer,
		StemcellAlias:       DefaultStemcellAlias,
		StemcellOS:          DefaultStemcellName,
		StemcellVersion:     DefaultStemcellVersion,
	}
}

func NewDeployment(config Config) *Deployment {
	return &Deployment{
		config: config,
	}
}

func (d *Deployment) CloudConfigValidation(data []byte) error {
	lo.G.Debug("Cloud Config:", string(data))
	c := &enaml.CloudConfigManifest{}
	if err := yaml.Unmarshal(data, &c); err != nil {
		return err
	}

	if !c.ContainsAZName(d.config.AZ) {
		return fmt.Errorf("AZ [%s] is not defined as a AZ in cloud config", d.config.AZ)
	}

	if !c.ContainsVMType(d.config.VMType) {
		return fmt.Errorf("VMType[%s] is not defined as a VMType in cloud config", d.config.VMType)
	}
	if !c.ContainsDiskType(d.config.DiskType) {
		return fmt.Errorf("DiskType[%s] is not defined as a DiskType in cloud config", d.config.DiskType)
	}
	return nil
}

func (d *Deployment) CreateDeploymentManifest(cloudConfig []byte) (enaml.DeploymentManifest, error) {
	manifest := enaml.DeploymentManifest{}
	if err := d.CloudConfigValidation(cloudConfig); err != nil {
		return manifest, err
	}

	manifest.SetName("minio")
	manifest.AddRelease(enaml.Release{
		Name:    "minio",
		URL:     d.config.MinioReleaseURL,
		SHA1:    d.config.MinioReleaseSHA,
		Version: d.config.MinioReleaseVersion,
	})
	manifest.AddStemcell(enaml.Stemcell{
		Alias:   d.config.StemcellAlias,
		OS:      d.config.StemcellOS,
		Version: d.config.StemcellVersion,
	})
	manifest.SetUpdate(d.CreateUpdate())
	manifest.AddInstanceGroup(d.CreateMinioServer())
	return manifest, nil
}

func (d *Deployment) CreateMinioServer() *enaml.InstanceGroup {
	ig := &enaml.InstanceGroup{
		Name:               "minio-server",
		Instances:          1,
		VMType:             d.config.VMType,
		PersistentDiskType: d.config.DiskType,
		AZs:                []string{d.config.AZ},
		Stemcell:           d.config.StemcellAlias,
	}
	ig.AddNetwork(enaml.Network{
		Name:      d.config.NetworkName,
		StaticIPs: []string{d.config.IP},
	})
	ig.AddJob(&enaml.InstanceJob{
		Name:    "minio-server",
		Release: "minio",
		Properties: minio_server.MinioServerJob{
			Region: d.config.Region,
			Credential: &minio_server.Credential{
				Accesskey: d.config.AccessKey,
				Secretkey: d.config.SecretKey,
			},
		},
	})
	return ig
}

//CreateUpdate -
func (d *Deployment) CreateUpdate() enaml.Update {
	return enaml.Update{
		Canaries:        1,
		MaxInFlight:     3,
		Serial:          false,
		CanaryWatchTime: "1000-60000",
		UpdateWatchTime: "1000-60000",
	}
}
