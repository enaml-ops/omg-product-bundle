package vault

import (
	"fmt"
	"strings"

	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-product-bundle/products/vault/enaml-gen/consul"
	vaultlib "github.com/enaml-ops/omg-product-bundle/products/vault/enaml-gen/vault"
	"github.com/enaml-ops/pluginlib/cred"
	"github.com/enaml-ops/pluginlib/pcli"
	"github.com/enaml-ops/pluginlib/pluginutil"
	"github.com/enaml-ops/pluginlib/productv1"
	"github.com/xchapter7x/lo"
)

const (
	BoshVaultReleaseURL  = "https://bosh.io/d/github.com/cloudfoundry-community/vault-boshrelease?v=0.3.0"
	BoshVaultReleaseVer  = "0.3.0"
	BoshVaultReleaseSHA  = "bd1ae82104dcf36abf64875fc5a46e1661bf2eac"
	BoshConsulReleaseURL = "https://bosh.io/d/github.com/cloudfoundry-community/consul-boshrelease?v=20"
	BoshConsulReleaseVer = "20"
	BoshConsulReleaseSHA = "9a0591c6b4d88d7d756ea933e14ddf112d05f334"
	StemcellName         = "ubuntu-trusty"
	StemcellAlias        = "ubuntu-trusty"
	StemcellVersion      = "3263.8"
)

type jobBucket struct {
	JobName   string
	JobType   int
	Instances int
}

type Plugin struct {
	PluginVersion    string   `omg:"-"`
	NetworkName      string   `omg:"network"`
	IPs              []string `omg:"ip"`
	VMTypeName       string   `omg:"vm-type"`
	DiskTypeName     string   `omg:"disk-type"`
	AZs              []string `omg:"az"`
	StemcellName     string
	StemcellURL      string `omg:"stemcell-url,optional"`
	StemcellVersion  string `omg:"stemcell-ver"`
	StemcellSHA      string `omg:"stemcell-sha,optional"`
	VaultReleaseURL  string `omg:"vault-release-url"`
	VaultReleaseVer  string `omg:"vault-release-version"`
	VaultReleaseSHA  string `omg:"vault-release-sha"`
	ConsulReleaseURL string `omg:"consul-release-url"`
	ConsulReleaseVer string `omg:"consul-release-version"`
	ConsulReleaseSHA string `omg:"consul-release-sha"`
}

func (s *Plugin) GetFlags() (flags []pcli.Flag) {
	return []pcli.Flag{
		pcli.CreateStringSliceFlag("ip", "multiple static ips for each vault VM Node"),
		pcli.CreateStringSliceFlag("az", "list of AZ names to use"),
		pcli.CreateStringFlag("network", "the name of the network to use"),
		pcli.CreateStringFlag("vm-type", "name of your desired vm type"),
		pcli.CreateStringFlag("disk-type", "name of your desired disk type"),
		pcli.CreateStringFlag("stemcell-url", "the url of the stemcell you wish to use"),
		pcli.CreateStringFlag("stemcell-ver", "the version number of the stemcell you wish to use"),
		pcli.CreateStringFlag("stemcell-sha", "the sha of the stemcell you will use"),
		pcli.CreateStringFlag("stemcell-name", "the name of the stemcell you will use", s.GetMeta().Stemcell.Name),
		pcli.CreateStringFlag("vault-release-url", "vault bosh release url", BoshVaultReleaseURL),
		pcli.CreateStringFlag("vault-release-version", "vault bosh release version", BoshVaultReleaseVer),
		pcli.CreateStringFlag("vault-release-sha", "vault bosh release sha", BoshVaultReleaseSHA),
		pcli.CreateStringFlag("consul-release-url", "consul bosh release url", BoshConsulReleaseURL),
		pcli.CreateStringFlag("consul-release-version", "consul bosh release version", BoshConsulReleaseVer),
		pcli.CreateStringFlag("consul-release-sha", "consul bosh release sha", BoshConsulReleaseSHA),
	}
}

func (s *Plugin) GetMeta() product.Meta {
	return product.Meta{
		Name: "vault",
		Stemcell: enaml.Stemcell{
			Name:    StemcellName,
			Alias:   StemcellAlias,
			Version: StemcellVersion,
		},
		Releases: []enaml.Release{
			enaml.Release{
				Name:    "vault",
				Version: BoshVaultReleaseVer,
				URL:     BoshVaultReleaseURL,
				SHA1:    BoshVaultReleaseSHA,
			},
			enaml.Release{
				Name:    "consul",
				Version: BoshConsulReleaseVer,
				URL:     BoshConsulReleaseURL,
				SHA1:    BoshConsulReleaseSHA,
			},
		},
		Properties: map[string]interface{}{
			"version":        s.PluginVersion,
			"vault-release":  strings.Join([]string{BoshVaultReleaseURL, BoshVaultReleaseVer, BoshVaultReleaseSHA}, " / "),
			"consul-release": strings.Join([]string{BoshConsulReleaseURL, BoshConsulReleaseVer, BoshConsulReleaseSHA}, " / "),
		},
	}
}

func (s *Plugin) GetProduct(args []string, cloudConfig []byte, cs cred.Store) ([]byte, error) {
	c := pluginutil.NewContext(args, pluginutil.ToCliFlagArray(s.GetFlags()))

	err := pcli.UnmarshalFlags(s, c)
	if err != nil {
		return nil, err
	}

	err = s.cloudconfigValidation(enaml.NewCloudConfigManifest(cloudConfig))
	if err != nil {
		lo.G.Error("invalid settings for cloud config on target bosh: ", err)
		return nil, err
	}

	lo.G.Debug("context", c)
	var dm = new(enaml.DeploymentManifest)
	dm.SetName("vault")
	dm.AddRemoteRelease("vault", s.VaultReleaseVer, s.VaultReleaseURL, s.VaultReleaseSHA)
	dm.AddRemoteRelease("consul", s.ConsulReleaseVer, s.ConsulReleaseURL, s.ConsulReleaseSHA)
	dm.AddRemoteStemcell(s.StemcellName, s.StemcellName, s.StemcellVersion, s.StemcellURL, s.StemcellSHA)

	dm.AddInstanceGroup(s.NewVaultInstanceGroup())
	dm.Update = enaml.Update{
		MaxInFlight:     1,
		UpdateWatchTime: "30000-300000",
		CanaryWatchTime: "30000-300000",
		Serial:          false,
		Canaries:        1,
	}
	return dm.Bytes(), err
}

func (s *Plugin) NewVaultInstanceGroup() (ig *enaml.InstanceGroup) {
	ig = &enaml.InstanceGroup{
		Name:               "vault",
		Instances:          len(s.IPs),
		VMType:             s.VMTypeName,
		AZs:                s.AZs,
		Stemcell:           s.StemcellName,
		PersistentDiskType: s.DiskTypeName,
		Jobs: []enaml.InstanceJob{
			s.createVaultJob(),
			s.createConsulJob(),
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

func (s *Plugin) createVaultJob() enaml.InstanceJob {
	return enaml.InstanceJob{
		Name:    "vault",
		Release: "vault",
		Properties: &vaultlib.VaultJob{
			Vault: &vaultlib.Vault{
				Backend: &vaultlib.Backend{
					UseConsul: true,
				},
			},
		},
	}
}
func (s *Plugin) createConsulJob() enaml.InstanceJob {
	return enaml.InstanceJob{
		Name:    "consul",
		Release: "consul",
		Properties: &consul.ConsulJob{
			Consul: &consul.Consul{
				JoinHosts: s.IPs,
			},
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
	return
}
