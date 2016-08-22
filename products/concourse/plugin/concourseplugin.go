package concourseplugin

import (
	"strings"

	"github.com/codegangsta/cli"
	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-cli/utils"
	"github.com/enaml-ops/omg-product-bundle/products/concourse"
	"github.com/enaml-ops/pluginlib/pcli"
	"github.com/enaml-ops/pluginlib/product"
	"github.com/enaml-ops/pluginlib/util"
	"github.com/xchapter7x/lo"
)

type ConcoursePlugin struct {
	PluginVersion string
}

const (
	defaultConcourseReleaseURL string = "https://bosh.io/d/github.com/concourse/concourse?v=1.6.0"
	defaultConcourseReleaseSHA string = "6fcd122fd95005e076956cf25a940f082cccb222"
	defaultConcourseReleaseVer string = "1.6.0"
	defaultGardenReleaseURL    string = "https://bosh.io/d/github.com/cloudfoundry-incubator/garden-runc-release?v=0.5.0"
	defaultGardenReleaseSHA    string = "eee2d6fc57a03f1d7e5d2625d617dda18704610f"
	defaultGardenReleaseVer    string = "0.5.0"

	concoursePassword   string = "concourse-password"
	concourseUsername   string = "concourse-username"
	externalURL         string = "external-url"
	webIPs              string = "web-ip"
	networkName         string = "network-name"
	az                  string = "az"
	deploymentName      string = "deployment-name"
	webVMType           string = "web-vm-type"
	databaseVMType      string = "database-vm-type"
	workerVMType        string = "worker-vm-type"
	workerInstances     string = "worker-instance-count"
	databaseStorageType string = "database-storage-type"
	postgresqlDbPwd     string = "concourse-db-pwd"
	concourseReleaseURL string = "concourse-release-url"
	concourseReleaseSHA string = "concourse-release-sha"
	concourseReleaseVer string = "concourse-release-ver"
	gardenReleaseURL    string = "garden-release-url"
	gardenReleaseSHA    string = "garden-release-sha"
	gardenReleaseVer    string = "garden-release-ver"
	stemcellAlias       string = "stemcell-alias"
	stemcellOS          string = "stemcell-os"
	stemcellVersion     string = "stemcell-version"
	tlsCert             string = "tls-cert"
	tlsKey              string = "tls-key"
)

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

		pcli.Flag{FlagType: pcli.StringFlag, Name: stemcellAlias, Value: "trusty", Usage: "alias of stemcell"},
		pcli.Flag{FlagType: pcli.StringFlag, Name: stemcellOS, Value: "ubuntu-trusty", Usage: "os of stemcell"},
		pcli.Flag{FlagType: pcli.StringFlag, Name: stemcellVersion, Value: "latest", Usage: "version of stemcell"},

		//pcli.Flag{FlagType: pcli.StringFlag, Name: tlsCert, Usage: "tls cert for concourse web"},
		//pcli.Flag{FlagType: pcli.StringFlag, Name: tlsKey, Usage: "tls key for concourse web"},
	}
	return
}

func (s *ConcoursePlugin) GetMeta() product.Meta {
	return product.Meta{
		Name: "concourse",
		Properties: map[string]interface{}{
			"version":           s.PluginVersion,
			"concourse-release": strings.Join([]string{defaultConcourseReleaseURL, defaultConcourseReleaseVer, defaultConcourseReleaseSHA}, " / "),
			"garden-release":    strings.Join([]string{defaultGardenReleaseURL, defaultGardenReleaseVer, defaultGardenReleaseSHA}, " / "),
		},
	}
}

func (s *ConcoursePlugin) GetProduct(args []string, cloudConfig []byte) (b []byte) {
	if len(cloudConfig) == 0 {
		lo.G.Debug("concourseplugin: empty cloud config")
		panic("cloud config cannot be empty")
	}
	c := pluginutil.NewContext(args, pluginutil.ToCliFlagArray(s.GetFlags()))
	dm := NewDeploymentManifest(c, cloudConfig)
	return dm.Bytes()
}

func NewDeploymentManifest(c *cli.Context, cloudConfig []byte) enaml.DeploymentManifest {
	var deployment = concourse.NewDeployment()
	deployment.DeploymentName = c.String(deploymentName)

	if c.IsSet(postgresqlDbPwd) {
		deployment.PostgresPassword = c.String(postgresqlDbPwd)
	} else {
		deployment.PostgresPassword = utils.NewPassword(20)
	}
	if c.IsSet(concoursePassword) {
		deployment.ConcoursePassword = c.String(concoursePassword)
	} else {
		deployment.ConcoursePassword = utils.NewPassword(20)
	}
	deployment.ConcourseUserName = c.String(concourseUsername)
	deployment.ConcourseURL = c.String(externalURL)
	deployment.NetworkName = c.String(networkName)
	deployment.WebIPs = c.StringSlice(webIPs)
	deployment.WebVMType = c.String(webVMType)
	deployment.WorkerVMType = c.String(workerVMType)
	deployment.DatabaseVMType = c.String(databaseVMType)
	deployment.DatabaseStorageType = c.String(databaseStorageType)
	deployment.AZs = c.StringSlice(az)
	deployment.WorkerInstances = c.Int(workerInstances)
	deployment.DeploymentName = c.String(deploymentName)
	deployment.ConcourseReleaseURL = c.String(concourseReleaseURL)
	deployment.ConcourseReleaseSHA = c.String(concourseReleaseSHA)
	deployment.ConcourseReleaseVer = c.String(concourseReleaseVer)
	deployment.StemcellAlias = c.String(stemcellAlias)
	deployment.StemcellOS = c.String(stemcellOS)
	deployment.StemcellVersion = c.String(stemcellVersion)
	deployment.GardenReleaseURL = c.String(gardenReleaseURL)
	deployment.GardenReleaseSHA = c.String(gardenReleaseSHA)
	deployment.GardenReleaseVer = c.String(gardenReleaseVer)

	/*
		TLS causing currently causing issues, need to investigate further
		if c.IsSet(tlsCert) && c.IsSet(tlsKey) {
			deployment.TLSCert = c.String(tlsCert)
			deployment.TLSKey = c.String(tlsKey)
		} else {
			if _, cert, key, err := utils.GenerateCert([]string{c.String(concourseIP)}); err != nil {
				panic(err.Error())
			} else {
				deployment.TLSCert = cert
				deployment.TLSKey = key
			}
		}*/

	if err := deployment.Initialize(cloudConfig); err != nil {
		panic(err.Error())
	}
	return deployment.GetDeployment()
}
