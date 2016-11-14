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
	"github.com/enaml-ops/pluginlib/product"
	"github.com/xchapter7x/lo"
)

const (
	defaultConcourseReleaseURL = "https://bosh.io/d/github.com/concourse/concourse?v=2.2.1"
	defaultConcourseReleaseSHA = "879d5cb45d12f173ff4c7912c7c7cdcd3e18c442"
	defaultConcourseReleaseVer = "2.2.1"
	defaultGardenReleaseURL    = "https://bosh.io/d/github.com/cloudfoundry-incubator/garden-runc-release?v=0.8.0"
	defaultGardenReleaseSHA    = "20e98ea84c8f4426bba00bbca17d931e27d3c07d"
	defaultGardenReleaseVer    = "0.8.0"

	concoursePassword      = "concourse-password"
	concourseUsername      = "concourse-username"
	externalURL            = "external-url"
	webIPs                 = "web-ip"
	networkName            = "network-name"
	az                     = "az"
	deploymentName         = "deployment-name"
	webVMType              = "web-vm-type"
	databaseVMType         = "database-vm-type"
	workerVMType           = "worker-vm-type"
	workerInstances        = "worker-instance-count"
	databaseStorageType    = "database-storage-type"
	postgresqlDbPwd        = "concourse-db-pwd"
	concourseReleaseURL    = "concourse-release-url"
	concourseReleaseSHA    = "concourse-release-sha"
	concourseReleaseVer    = "concourse-release-ver"
	gardenReleaseURL       = "garden-release-url"
	gardenReleaseSHA       = "garden-release-sha"
	gardenReleaseVer       = "garden-release-ver"
	stemcellAlias          = "stemcell-alias"
	stemcellOS             = "stemcell-os"
	stemcellVersion        = "stemcell-version"
	defaultStemcellAlias   = "trusty"
	defaultStemcellName    = "ubuntu-trusty"
	defaultStemcellVersion = "latest"
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

// ConcoursePlugin is an omg product plugin for deploying Concourse.
type ConcoursePlugin struct {
	PluginVersion string
	cfg           Config
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
	err := pcli.UnmarshalFlags(&s.cfg, c)
	if err != nil {
		return nil, err
	}

	makePassword(&s.cfg.PostgresPassword)
	makePassword(&s.cfg.ConcoursePassword)

	dm, err := s.newDeploymentManifest(cloudConfig)
	return dm.Bytes(), err
}

func makePassword(s *string) {
	if *s == "" {
		*s = pluginutil.NewPassword(20)
	}
}

func (s *ConcoursePlugin) newDeploymentManifest(cloudConfig []byte) (enaml.DeploymentManifest, error) {
	cd := concourse.NewDeployment()
	cd.DeploymentName = s.cfg.DeploymentName
	cd.ConcourseUserName = s.cfg.ConcourseUsername
	cd.ConcoursePassword = s.cfg.ConcoursePassword
	cd.ConcourseURL = s.cfg.ExternalURL
	cd.NetworkName = s.cfg.NetworkName
	cd.WebIPs = s.cfg.WebIPs
	cd.WebVMType = s.cfg.WebVMType
	cd.WorkerVMType = s.cfg.WorkerVMType
	cd.DatabaseVMType = s.cfg.DatabaseVMType
	cd.DatabaseStorageType = s.cfg.DatabaseStorageType
	cd.AZs = s.cfg.AZs
	cd.WorkerInstances = s.cfg.WorkerInstances
	cd.ConcourseReleaseURL = s.cfg.ConcourseReleaseURL
	cd.ConcourseReleaseSHA = s.cfg.ConcourseReleaseSHA
	cd.ConcourseReleaseVer = s.cfg.ConcourseReleaseVersion
	cd.StemcellAlias = s.cfg.StemcellAlias
	cd.StemcellOS = s.cfg.StemcellOS
	cd.StemcellVersion = s.cfg.StemcellVersion
	cd.GardenReleaseURL = s.cfg.GardenReleaseURL
	cd.GardenReleaseSHA = s.cfg.GardenReleaseSHA
	cd.GardenReleaseVer = s.cfg.GardenReleaseVersion

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
