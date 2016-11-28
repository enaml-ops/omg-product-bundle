package gemfire_plugin

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-product-bundle/products/p-gemfire/enaml-gen/locator"
	"github.com/enaml-ops/omg-product-bundle/products/p-gemfire/enaml-gen/server"
	"github.com/enaml-ops/pluginlib/cred"
	"github.com/enaml-ops/pluginlib/pcli"
	"github.com/enaml-ops/pluginlib/pluginutil"
	"github.com/enaml-ops/pluginlib/productv1"
	"github.com/xchapter7x/lo"
)

type Plugin struct {
	Version string `omg:"-"`

	DeploymentName              string
	NetworkName                 string
	GemfireReleaseVer           string
	StemcellName                string
	StemcellVer                 string
	StemcellAlias               string
	AZs                         []string `omg:"az"`
	LocatorStaticIPs            []string `omg:"locator-static-ip"`
	ServerStaticIPs             []string `omg:"server-static-ip,optional"`
	ServerInstanceCount         int
	GemfireLocatorPort          int
	GemfireLocatorRestPort      int
	GemfireServerPort           int
	GemfireLocatorVMMemory      int    `omg:"gemfire-locator-vm-memory"`
	GemfireLocatorVMSize        string `omg:"gemfire-locator-vm-size"`
	GemfireServerVMSize         string `omg:"gemfire-server-vm-size"`
	GemfireServerVMMemory       int    `omg:"gemfire-server-vm-memory"`
	ServerDevRestAPIPort        int    `omg:"gemfire-dev-rest-api-port"`
	ServerDevActive             bool   `omg:"gemfire-dev-rest-api-active"`
	AuthnActive                 bool   `omg:"use-authn,optional"`
	SecurityClientAuthenticator string `omg:"security-client-authenticator,optional"`
	KeystoreRemotePath          string `omg:"keystore-remote-path,optional"`
	PublicKeyPass               string `omg:"public-key-pass,optional"`
	KeystoreLocalPath           string `omg:"keystore-local-path,optional"`
	SecurityJarLocalPath        string `omg:"security-jar-local-path,optional"`
}

func (p *Plugin) authnFlagsValid() bool {
	if p.AuthnActive {
		flagValues := []string{
			p.PublicKeyPass,
			p.KeystoreLocalPath,
			p.SecurityJarLocalPath,
		}

		for _, val := range flagValues {

			if val == "" {
				return false
			}
		}

		if _, err := os.Stat(p.KeystoreLocalPath); os.IsNotExist(err) {
			lo.G.Errorf("file does not exist: %v", p.KeystoreLocalPath)
			return false
		}

		if _, err := os.Stat(p.SecurityJarLocalPath); os.IsNotExist(err) {
			lo.G.Errorf("file does not exist: %v", p.SecurityJarLocalPath)
			return false
		}
	}
	return true
}

// GetProduct generates a BOSH deployment manifest for p-gemfire.
func (p *Plugin) GetProduct(args []string, cloudConfig []byte, cs cred.Store) ([]byte, error) {
	c := pluginutil.NewContext(args, pluginutil.ToCliFlagArray(p.GetFlags()))
	err := pcli.UnmarshalFlags(p, c)

	if err != nil {
		return nil, err
	}

	if !p.authnFlagsValid() {
		return nil, ActiveAuthNErr
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

	ltr := NewLocatorGroup(p.NetworkName, p.LocatorStaticIPs, p.GemfireLocatorPort, p.GemfireLocatorRestPort, p.GemfireLocatorVMMemory, p.GemfireLocatorVMSize)
	locatorInstanceGroup := ltr.GetInstanceGroup(p.getLocatorAuthn())
	locatorInstanceGroup.Stemcell = p.StemcellAlias
	locatorInstanceGroup.AZs = p.AZs
	deploymentManifest.AddInstanceGroup(locatorInstanceGroup)

	svr := NewServerGroup(p.NetworkName, p.GemfireServerPort, p.ServerInstanceCount, p.ServerStaticIPs, p.GemfireServerVMSize, p.GemfireServerVMMemory, p.ServerDevRestAPIPort, p.ServerDevActive, ltr)
	serverInstanceGroup := svr.GetInstanceGroup(p.getServerAuthn())
	serverInstanceGroup.Stemcell = p.StemcellAlias
	serverInstanceGroup.AZs = p.AZs
	deploymentManifest.AddInstanceGroup(serverInstanceGroup)
	return deploymentManifest.Bytes(), nil
}

func (p *Plugin) getLocatorAuthn() locator.Authn {
	authn := locator.Authn{}

	if p.AuthnActive {

		buf := new(bytes.Buffer)

		if b, err := ioutil.ReadFile(p.SecurityJarLocalPath); err == nil {
			encoder := base64.NewEncoder(base64.StdEncoding, buf)
			encoder.Write(b)
			encoder.Close()
			authn.SecurityJarBase64Bits = buf.String()
		}
	}
	return authn
}

func (p *Plugin) getServerAuthn() server.Authn {
	authn := server.Authn{}

	if p.AuthnActive {
		buf := new(bytes.Buffer)

		if b, err := ioutil.ReadFile(p.KeystoreLocalPath); err == nil {
			encoder := base64.NewEncoder(base64.StdEncoding, buf)
			encoder.Write(b)
			encoder.Close()
			authn.KeystoreBits = buf.String()
		}
		authn.Enabled = true
		authn.SecurityPublickeyPass = p.PublicKeyPass
		authn.SecurityKeystoreFilepath = p.KeystoreRemotePath
		authn.SecurityClientAuthenticator = p.SecurityClientAuthenticator
	}
	return authn
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
		pcli.Flag{
			FlagType: pcli.BoolFlag,
			Name:     "use-authn",
			Usage:    "activates authN for your gemfire deployment",
		},
		pcli.Flag{
			FlagType: pcli.StringFlag,
			Name:     "security-client-authenticator",
			Value:    SecurityClientAuthenticatorDefault,
			Usage:    "will populate: what should the value of the gemfire property for security-client-authenticator be - gemfire.authn.security_client_authenticator",
		},
		pcli.Flag{
			FlagType: pcli.StringFlag,
			Name:     "keystore-remote-path",
			Value:    KeystoreRemotePathDefault,
			Usage:    "will populate: path on remote system for your keystore - gemfire.authn.security_keystore_filepath",
		},
		pcli.Flag{
			FlagType: pcli.StringFlag,
			Name:     "public-key-pass",
			Usage:    "will populate: password for the given key - gemfire.authn.security_publickey_pass",
		},
		pcli.Flag{
			FlagType: pcli.StringFlag,
			Name:     "keystore-local-path",
			Usage:    "will populate: keystore file bits, base64 encoded - gemfire.authn.keystore_bits",
		},
		pcli.Flag{
			FlagType: pcli.StringFlag,
			Name:     "security-jar-local-path",
			Usage:    "will populate: base64 encoding of authentication security jar - gemfire.authn.security_jar_base64_bits",
		},
	}
}
