package plugin

import (
	"fmt"
	"strings"

	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-product-bundle/products/minio"
	"github.com/enaml-ops/pluginlib/cred"
	"github.com/enaml-ops/pluginlib/pcli"
	"github.com/enaml-ops/pluginlib/pluginutil"
	"github.com/enaml-ops/pluginlib/productv1"
)

const (
	networkName     = "network-name"
	az              = "az"
	ip              = "ip"
	vmType          = "vm-type"
	diskType        = "disk-type"
	region          = "region"
	accessKey       = "access-key"
	secretKey       = "secret-key"
	minioReleaseURL = "minio-release-url"
	minioReleaseSHA = "minio-release-sha"
	minioReleaseVer = "minio-release-ver"
	stemcellAlias   = "stemcell-alias"
	stemcellOS      = "stemcell-os"
	stemcellVersion = "stemcell-version"
)

// Plugin is an omg product plugin for deploying Minio.
type Plugin struct {
	PluginVersion string
	cfg           minio.Config
}

func (s *Plugin) GetFlags() (flags []pcli.Flag) {
	flags = []pcli.Flag{
		pcli.Flag{FlagType: pcli.StringFlag, Name: az, Usage: "az name to use"},
		pcli.Flag{FlagType: pcli.StringFlag, Name: networkName, Usage: "the name of the network to use"},
		pcli.Flag{FlagType: pcli.StringFlag, Name: ip, Usage: "ip for minio"},
		pcli.Flag{FlagType: pcli.StringFlag, Name: vmType, Usage: "type of vm to use"},
		pcli.Flag{FlagType: pcli.StringFlag, Name: diskType, Usage: "type of disk to use"},
		pcli.Flag{FlagType: pcli.StringFlag, Name: minioReleaseURL, Value: minio.DefaultMinioReleaseURL, Usage: "release url for minio bosh release"},
		pcli.Flag{FlagType: pcli.StringFlag, Name: minioReleaseSHA, Value: minio.DefaultMinioReleaseSHA, Usage: "release sha for minio bosh release"},
		pcli.Flag{FlagType: pcli.StringFlag, Name: minioReleaseVer, Value: minio.DefaultMinioReleaseVer, Usage: "release version for minio bosh release"},
		pcli.Flag{FlagType: pcli.StringFlag, Name: stemcellAlias, Value: minio.DefaultStemcellAlias, Usage: "alias of stemcell"},
		pcli.Flag{FlagType: pcli.StringFlag, Name: stemcellOS, Value: minio.DefaultStemcellName, Usage: "os of stemcell"},
		pcli.Flag{FlagType: pcli.StringFlag, Name: stemcellVersion, Value: minio.DefaultStemcellVersion, Usage: "version of stemcell"},
		pcli.Flag{FlagType: pcli.StringFlag, Name: region, Value: minio.DefaultRegion, Usage: "region for s3 blobstore"},
		pcli.Flag{FlagType: pcli.StringFlag, Name: accessKey, Usage: "access key"},
		pcli.Flag{FlagType: pcli.StringFlag, Name: secretKey, Usage: "secret key"},
	}
	return
}

func (s *Plugin) GetMeta() product.Meta {
	return product.Meta{
		Name: "minio",
		Stemcell: enaml.Stemcell{
			Name:    minio.DefaultStemcellName,
			Alias:   minio.DefaultStemcellAlias,
			Version: minio.DefaultStemcellVersion,
		},
		Releases: []enaml.Release{
			enaml.Release{
				Name:    "minio",
				Version: minio.DefaultMinioReleaseVer,
				URL:     minio.DefaultMinioReleaseURL,
				SHA1:    minio.DefaultMinioReleaseSHA,
			},
		},
		Properties: map[string]interface{}{
			"version":       s.PluginVersion,
			"minio-release": strings.Join([]string{minio.DefaultMinioReleaseURL, minio.DefaultMinioReleaseVer, minio.DefaultMinioReleaseSHA}, " / "),
		},
	}
}

func (s *Plugin) GetProduct(args []string, cloudConfig []byte, cs cred.Store) ([]byte, error) {
	if len(cloudConfig) == 0 {
		return nil, fmt.Errorf("cloud config cannot be empty")
	}

	c := pluginutil.NewContext(args, pluginutil.ToCliFlagArray(s.GetFlags()))
	err := pcli.UnmarshalFlags(&s.cfg, c)
	if err != nil {
		return nil, err
	}

	dm, err := minio.NewDeployment(s.cfg).CreateDeploymentManifest(cloudConfig)
	return dm.Bytes(), err
}
