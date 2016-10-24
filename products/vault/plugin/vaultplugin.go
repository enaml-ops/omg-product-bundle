package vault

import (
	"fmt"
	"strings"

	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-product-bundle/products/vault/enaml-gen/consul"
	vaultlib "github.com/enaml-ops/omg-product-bundle/products/vault/enaml-gen/vault"
	"github.com/enaml-ops/pluginlib/pcli"
	"github.com/enaml-ops/pluginlib/pluginutil"
	"github.com/enaml-ops/pluginlib/product"
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
	PluginVersion    string
	NetworkName      string
	IPs              []string
	VMTypeName       string
	DiskTypeName     string
	AZs              []string
	StemcellName     string
	StemcellURL      string
	StemcellVersion  string
	StemcellSHA      string
	VaultReleaseURL  string
	VaultReleaseVer  string
	VaultReleaseSHA  string
	ConsulReleaseURL string
	ConsulReleaseVer string
	ConsulReleaseSHA string
}

func (s *Plugin) GetFlags() (flags []pcli.Flag) {
	return []pcli.Flag{
		pcli.Flag{FlagType: pcli.StringSliceFlag, Name: "ip", Usage: "multiple static ips for each vault VM Node"},
		pcli.Flag{FlagType: pcli.StringSliceFlag, Name: "az", Usage: "list of AZ names to use"},
		pcli.Flag{FlagType: pcli.StringFlag, Name: "network", Usage: "the name of the network to use"},
		pcli.Flag{FlagType: pcli.StringFlag, Name: "vm-type", Usage: "name of your desired vm type"},
		pcli.Flag{FlagType: pcli.StringFlag, Name: "disk-type", Usage: "name of your desired disk type"},
		pcli.Flag{FlagType: pcli.StringFlag, Name: "stemcell-url", Usage: "the url of the stemcell you wish to use"},
		pcli.Flag{FlagType: pcli.StringFlag, Name: "stemcell-ver", Usage: "the version number of the stemcell you wish to use"},
		pcli.Flag{FlagType: pcli.StringFlag, Name: "stemcell-sha", Usage: "the sha of the stemcell you will use"},
		pcli.Flag{FlagType: pcli.StringFlag, Name: "stemcell-name", Value: s.GetMeta().Stemcell.Name, Usage: "the name of the stemcell you will use"},
		pcli.Flag{FlagType: pcli.StringFlag, Name: "vault-release-url", Value: BoshVaultReleaseURL, Usage: "vault bosh release url"},
		pcli.Flag{FlagType: pcli.StringFlag, Name: "vault-release-version", Value: BoshVaultReleaseVer, Usage: "vault bosh release version"},
		pcli.Flag{FlagType: pcli.StringFlag, Name: "vault-release-sha", Value: BoshVaultReleaseSHA, Usage: "vault bosh release sha"},
		pcli.Flag{FlagType: pcli.StringFlag, Name: "consul-release-url", Value: BoshConsulReleaseURL, Usage: "consul bosh release url"},
		pcli.Flag{FlagType: pcli.StringFlag, Name: "consul-release-version", Value: BoshConsulReleaseVer, Usage: "consul bosh release version"},
		pcli.Flag{FlagType: pcli.StringFlag, Name: "consul-release-sha", Value: BoshConsulReleaseSHA, Usage: "consul bosh release sha"},
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

func (s *Plugin) GetProduct(args []string, cloudConfig []byte) (b []byte, err error) {
	c := pluginutil.NewContext(args, pluginutil.ToCliFlagArray(s.GetFlags()))

	s.IPs = c.StringSlice("ip")
	s.AZs = c.StringSlice("az")
	s.NetworkName = c.String("network")
	s.StemcellName = c.String("stemcell-name")
	s.StemcellVersion = c.String("stemcell-ver")
	s.StemcellSHA = c.String("stemcell-sha")
	s.StemcellURL = c.String("stemcell-url")
	s.VMTypeName = c.String("vm-type")
	s.DiskTypeName = c.String("disk-type")
	s.VaultReleaseURL = c.String("vault-release-url")
	s.VaultReleaseVer = c.String("vault-release-version")
	s.VaultReleaseSHA = c.String("vault-release-sha")
	s.ConsulReleaseURL = c.String("consul-release-url")
	s.ConsulReleaseVer = c.String("consul-release-version")
	s.ConsulReleaseSHA = c.String("consul-release-sha")

	if err = s.flagValidation(); err != nil {
		lo.G.Error("invalid arguments: ", err)
		return nil, err
	}

	if err = s.cloudconfigValidation(enaml.NewCloudConfigManifest(cloudConfig)); err != nil {
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
