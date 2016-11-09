package concourseplugin

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-product-bundle/products/concourse"
	"github.com/enaml-ops/pluginlib/cred"
	"github.com/enaml-ops/pluginlib/pcli"
	"github.com/enaml-ops/pluginlib/pluginutil"
	"github.com/enaml-ops/pluginlib/productv1"
	"github.com/xchapter7x/lo"
)

type ConcoursePlugin struct {
	PluginVersion string
}

const (
	defaultConcourseReleaseURL string = "https://bosh.io/d/github.com/concourse/concourse?v=2.2.1"
	defaultConcourseReleaseSHA string = "879d5cb45d12f173ff4c7912c7c7cdcd3e18c442"
	defaultConcourseReleaseVer string = "2.2.1"
	defaultGardenReleaseURL    string = "https://bosh.io/d/github.com/cloudfoundry-incubator/garden-runc-release?v=0.8.0"
	defaultGardenReleaseSHA    string = "20e98ea84c8f4426bba00bbca17d931e27d3c07d"
	defaultGardenReleaseVer    string = "0.8.0"

	concoursePassword      string = "concourse-password"
	concourseUsername      string = "concourse-username"
	externalURL            string = "external-url"
	webIPs                 string = "web-ip"
	networkName            string = "network-name"
	az                     string = "az"
	deploymentName         string = "deployment-name"
	webVMType              string = "web-vm-type"
	databaseVMType         string = "database-vm-type"
	workerVMType           string = "worker-vm-type"
	workerInstances        string = "worker-instance-count"
	databaseStorageType    string = "database-storage-type"
	postgresqlDbPwd        string = "concourse-db-pwd"
	concourseReleaseURL    string = "concourse-release-url"
	concourseReleaseSHA    string = "concourse-release-sha"
	concourseReleaseVer    string = "concourse-release-ver"
	gardenReleaseURL       string = "garden-release-url"
	gardenReleaseSHA       string = "garden-release-sha"
	gardenReleaseVer       string = "garden-release-ver"
	stemcellAlias          string = "stemcell-alias"
	stemcellOS             string = "stemcell-os"
	stemcellVersion        string = "stemcell-version"
	defaultStemcellAlias          = "trusty"
	defaultStemcellName           = "ubuntu-trusty"
	defaultStemcellVersion        = "latest"
)

// Config contains the configuration for a Concourse deployment.
type Config struct {
	DeploymentName      string
	ConcourseUsername   string
	ConcoursePassword   string   `omg:"concourse-password,optional"`
	ExternalURL         string   `omg:"external-url"`
	AZs                 []string `omg:"az"`
	NetworkName         string
	WebIPs              []string `omg:"web-ip"`
	WebVMType           string   `omg:"web-vm-type"`
	WorkerVMType        string   `omg:"worker-vm-type"`
	DatabaseVMType      string   `omg:"database-vm-type"`
	DatabaseStorageType string
	WorkerInstances     int    `omg:"worker-instance-count"`
	PostgresPassword    string `omg:"concourse-db-pwd,optional"`

	ConcourseReleaseURL     string `omg:"concourse-release-url"`
	ConcourseReleaseSHA     string `omg:"concourse-release-sha"`
	ConcourseReleaseVersion string `omg:"concourse-release-ver"`

	GardenReleaseURL     string `omg:"garden-release-url"`
	GardenReleaseSHA     string `omg:"garden-release-sha"`
	GardenReleaseVersion string `omg:"garden-release-ver"`

	StemcellAlias   string
	StemcellOS      string `omg:"stemcell-os"`
	StemcellVersion string
}

func (s *ConcoursePlugin) GetFlags() (flags []pcli.Flag) {
	flags = []pcli.Flag{
		pcli.Flag{FlagType: pcli.StringFlag, Name: deploymentName, Value: "concourse", Usage: "deployment name"},
		pcli.Flag{FlagType: pcli.StringFlag, Name: concourseUsername, Usage: "concourse user id"},
		pcli.Flag{FlagType: pcli.StringFlag, Name: concoursePassword, Usage: "concourse password"},
		pcli.Flag{FlagType: pcli.StringFlag, Name: externalURL, Usage: "URL to access concourse"},
		pcli.Flag{FlagType: pcli.StringSliceFlag, Name: az, Usage: "list of AZ names to use"},
		pcli.Flag{FlagType: pcli.StringFlag, Name: networkName, Usage: "the name of the network to use"},
		pcli.Flag{FlagType: pcli.StringSliceFlag, Name: webIPs, Usage: "ips for web jobs"},
		pcli.Flag{FlagType: pcli.StringFlag, Name: webVMType, Usage: "type of vm to use for web jobs"},
		pcli.Flag{FlagType: pcli.StringFlag, Name: workerVMType, Usage: "type of vm to use for worker jobs"},
		pcli.Flag{FlagType: pcli.StringFlag, Name: databaseVMType, Usage: "type of vm to use for database job"},
		pcli.Flag{FlagType: pcli.StringFlag, Name: databaseStorageType, Usage: "type of disk type for database job"},
		pcli.Flag{FlagType: pcli.IntFlag, Name: workerInstances, Value: "1", Usage: "number of worker instances"},
		pcli.Flag{FlagType: pcli.StringFlag, Name: postgresqlDbPwd, Usage: "password for postgres db"},

		pcli.Flag{FlagType: pcli.StringFlag, Name: concourseReleaseURL, Value: defaultConcourseReleaseURL, Usage: "release url for concourse bosh release"},
		pcli.Flag{FlagType: pcli.StringFlag, Name: concourseReleaseSHA, Value: defaultConcourseReleaseSHA, Usage: "release sha for concourse bosh release"},
		pcli.Flag{FlagType: pcli.StringFlag, Name: concourseReleaseVer, Value: defaultConcourseReleaseVer, Usage: "release version for concourse bosh release"},

		pcli.Flag{FlagType: pcli.StringFlag, Name: gardenReleaseURL, Value: defaultGardenReleaseURL, Usage: "release url for garden bosh release"},
		pcli.Flag{FlagType: pcli.StringFlag, Name: gardenReleaseSHA, Value: defaultGardenReleaseSHA, Usage: "release sha for garden bosh release"},
		pcli.Flag{FlagType: pcli.StringFlag, Name: gardenReleaseVer, Value: defaultGardenReleaseVer, Usage: "release version for garden bosh release"},

		pcli.Flag{FlagType: pcli.StringFlag, Name: stemcellAlias, Value: defaultStemcellAlias, Usage: "alias of stemcell"},
		pcli.Flag{FlagType: pcli.StringFlag, Name: stemcellOS, Value: defaultStemcellName, Usage: "os of stemcell"},
		pcli.Flag{FlagType: pcli.StringFlag, Name: stemcellVersion, Value: defaultStemcellVersion, Usage: "version of stemcell"},
	}
	return
}

func (s *ConcoursePlugin) GetMeta() product.Meta {
	return product.Meta{
		Name: "concourse",
		Stemcell: enaml.Stemcell{
			Name:    defaultStemcellName,
			Alias:   defaultStemcellAlias,
			Version: defaultStemcellVersion,
		},
		Releases: []enaml.Release{
			enaml.Release{
				Name:    "garden-runc",
				Version: defaultGardenReleaseVer,
				URL:     defaultGardenReleaseURL,
				SHA1:    defaultGardenReleaseSHA,
			},
			enaml.Release{
				Name:    "concourse",
				Version: defaultConcourseReleaseVer,
				URL:     defaultConcourseReleaseURL,
				SHA1:    defaultConcourseReleaseSHA,
			},
		},
		Properties: map[string]interface{}{
			"version":           s.PluginVersion,
			"concourse-release": strings.Join([]string{defaultConcourseReleaseURL, defaultConcourseReleaseVer, defaultConcourseReleaseSHA}, " / "),
			"garden-release":    strings.Join([]string{defaultGardenReleaseURL, defaultGardenReleaseVer, defaultGardenReleaseSHA}, " / "),
		},
	}
}

func (s *ConcoursePlugin) GetProduct(args []string, cloudConfig []byte, cs cred.Store) ([]byte, error) {
	if len(cloudConfig) == 0 {
		return nil, fmt.Errorf("cloud config cannot be empty")
	}

	c := pluginutil.NewContext(args, pluginutil.ToCliFlagArray(s.GetFlags()))
	cfg := &Config{}
	err := pcli.UnmarshalFlags(cfg, c)
	if err != nil {
		return nil, err
	}

	makePassword(&cfg.PostgresPassword)
	makePassword(&cfg.ConcoursePassword)

	dm, err := NewDeploymentManifest(cfg, cloudConfig)
	return dm.Bytes(), err
}

func makePassword(s *string) {
	if *s == "" {
		*s = pluginutil.NewPassword(20)
	}
}

func NewDeploymentManifest(c *Config, cloudConfig []byte) (enaml.DeploymentManifest, error) {
	cd := concourse.NewDeployment()
	cd.DeploymentName = c.DeploymentName
	cd.ConcourseUserName = c.ConcourseUsername
	cd.ConcourseURL = c.ExternalURL
	cd.NetworkName = c.NetworkName
	cd.WebIPs = c.WebIPs
	cd.WebVMType = c.WebVMType
	cd.WorkerVMType = c.WorkerVMType
	cd.DatabaseVMType = c.DatabaseVMType
	cd.DatabaseStorageType = c.DatabaseStorageType
	cd.AZs = c.AZs
	cd.WorkerInstances = c.WorkerInstances
	cd.ConcourseReleaseURL = c.ConcourseReleaseURL
	cd.ConcourseReleaseSHA = c.ConcourseReleaseSHA
	cd.ConcourseReleaseVer = c.ConcourseReleaseVersion
	cd.StemcellAlias = c.StemcellAlias
	cd.StemcellOS = c.StemcellOS
	cd.StemcellVersion = c.StemcellVersion
	cd.GardenReleaseURL = c.GardenReleaseURL
	cd.GardenReleaseSHA = c.GardenReleaseSHA
	cd.GardenReleaseVer = c.GardenReleaseVersion

	url, err := url.Parse(cd.ConcourseURL)
	if err != nil {
		return enaml.DeploymentManifest{}, fmt.Errorf("concourse-url invalid: %v", err)
	}
	if url.Scheme == "" {
		return enaml.DeploymentManifest{}, errors.New("concourse-url missing scheme")
	}

	if err = cd.Initialize(cloudConfig); err != nil {
		lo.G.Error(err.Error())
	}
	return cd.GetDeployment(), err
}
