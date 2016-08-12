package dockerregistry

import (
	"fmt"

	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/debian_nfs_server"
	"github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/nfs_mounter"
	"github.com/enaml-ops/omg-product-bundle/products/dockerregistry/enaml-gen/proxy"
	"github.com/enaml-ops/omg-product-bundle/products/dockerregistry/enaml-gen/registry"
	"github.com/xchapter7x/lo"
	"gopkg.in/yaml.v2"
)

const (
	sharePath           = "/var/vcap/nfs"
	rootPath            = "/var/vcap/nfs/shared"
	registryReleaseName = "docker-registry"
)

type DockerRegistry struct {
	DeploymentName           string
	DockerRegistryReleaseVer string
	DockerRegistryReleaseURL string
	DockerRegistryReleaseSHA string
	StemcellVersion          string
	StemcellAlias            string
	StemcellOS               string
	AZs                      []string
	NetworkName              string
	RegistryVMType           string
	RegistryIPs              []string
	ProxyVMType              string
	ProxyIPs                 []string
	PublicIP                 []string
	ProxyCert                string
	ProxyCertKey             string
	NFSServerVMType          string
	NFSDiskType              string
	NFSIP                    string
	Secret                   string
}

func (d *DockerRegistry) CreateDeploymentManifest(cloudConfig []byte) (*enaml.DeploymentManifest, error) {
	lo.G.Debug("Cloud Config:", string(cloudConfig))
	c := &enaml.CloudConfigManifest{}
	if err := yaml.Unmarshal(cloudConfig, &c); err != nil {
		return nil, err
	}
	if err := d.doCloudConfigValidation(c); err != nil {
		return nil, err
	}
	manifest := &enaml.DeploymentManifest{}
	manifest.SetName(d.DeploymentName)
	manifest.AddRelease(enaml.Release{
		Name:    registryReleaseName,
		URL:     d.DockerRegistryReleaseURL,
		SHA1:    d.DockerRegistryReleaseSHA,
		Version: d.DockerRegistryReleaseVer,
	})
	manifest.AddStemcell(enaml.Stemcell{
		Alias:   d.StemcellAlias,
		OS:      d.StemcellOS,
		Version: d.StemcellVersion,
	})

	manifest.SetUpdate(d.CreateUpdate())
	manifest.AddInstanceGroup(d.CreateNFSServerInstanceGroup())
	manifest.AddInstanceGroup(d.CreateRegistryInstanceGroup())
	manifest.AddInstanceGroup(d.CreateProxyInstanceGroup())
	return manifest, nil
}

func (d *DockerRegistry) doCloudConfigValidation(cloudConfigManifest *enaml.CloudConfigManifest) (err error) {

	for _, azName := range d.AZs {
		if !cloudConfigManifest.ContainsAZName(azName) {
			err = fmt.Errorf("AZ [%s] is not defined as a AZ in cloud config", azName)
			return
		}
	}

	if !cloudConfigManifest.ContainsVMType(d.RegistryVMType) {
		err = fmt.Errorf("RegistryVMType[%s] is not defined as a VMType in cloud config", d.RegistryVMType)
		return
	}
	if !cloudConfigManifest.ContainsVMType(d.ProxyVMType) {
		err = fmt.Errorf("ProxyVMType[%s] is not defined as a VMType in cloud config", d.ProxyVMType)
		return
	}
	if !cloudConfigManifest.ContainsVMType(d.NFSServerVMType) {
		err = fmt.Errorf("NFSServerType[%s] is not defined as a VMType in cloud config", d.NFSServerVMType)
		return
	}
	if !cloudConfigManifest.ContainsDiskType(d.NFSDiskType) {
		err = fmt.Errorf("NFSDiskType[%s] is not defined as a DiskType in cloud config", d.NFSDiskType)
		return
	}
	return
}
func (d *DockerRegistry) CreateRegistryInstanceGroup() *enaml.InstanceGroup {
	server := &enaml.InstanceGroup{
		Name:      "registry",
		Instances: len(d.RegistryIPs),
		Stemcell:  d.StemcellAlias,
		VMType:    d.RegistryVMType,
		AZs:       d.AZs,
	}
	server.AddNetwork(enaml.Network{
		Name:      d.NetworkName,
		StaticIPs: d.RegistryIPs,
	})
	server.AddJob(d.createRegistryJob())
	server.AddJob(d.createNFSMounterJob())
	return server
}

func (d *DockerRegistry) createRegistryJob() *enaml.InstanceJob {
	job := &enaml.InstanceJob{
		Name:    "registry",
		Release: registryReleaseName,
		Properties: &registry.RegistryJob{
			Docker: &registry.Docker{
				Registry: &registry.Registry{
					Root:   rootPath,
					Cookie: d.Secret,
				},
			},
		},
	}
	return job
}

func (d *DockerRegistry) createNFSMounterJob() *enaml.InstanceJob {
	job := &enaml.InstanceJob{
		Name:    "nfs_mounter",
		Release: registryReleaseName,
		Properties: &nfs_mounter.NfsMounterJob{
			NfsServer: &nfs_mounter.NfsServer{
				Address:   d.NFSIP,
				SharePath: sharePath,
			},
		},
	}

	return job
}

func (d *DockerRegistry) CreateProxyInstanceGroup() *enaml.InstanceGroup {
	server := &enaml.InstanceGroup{
		Name:      "proxy",
		Instances: len(d.ProxyIPs),
		Stemcell:  d.StemcellAlias,
		VMType:    d.ProxyVMType,
		AZs:       d.AZs,
	}
	server.AddNetwork(enaml.Network{
		Name:      d.NetworkName,
		StaticIPs: d.ProxyIPs,
	})
	server.AddJob(d.createProxyJob())
	return server
}

func (d *DockerRegistry) createProxyJob() *enaml.InstanceJob {
	job := &enaml.InstanceJob{
		Name:    "proxy",
		Release: registryReleaseName,
		Properties: &proxy.ProxyJob{
			Docker: &proxy.Docker{
				Proxy: &proxy.Proxy{
					Backend: &proxy.Backend{
						Hosts: d.RegistryIPs,
					},
					Ssl: &proxy.Ssl{
						Cert: d.ProxyCert,
						Key:  d.ProxyCertKey,
					},
				},
			},
		},
	}

	return job
}

func (d *DockerRegistry) CreateNFSServerInstanceGroup() *enaml.InstanceGroup {
	server := &enaml.InstanceGroup{
		Name:               "nfs-server",
		Instances:          1,
		Stemcell:           d.StemcellAlias,
		VMType:             d.NFSServerVMType,
		PersistentDiskType: d.NFSDiskType,
		AZs:                d.AZs,
	}
	server.AddNetwork(enaml.Network{
		Name:      d.NetworkName,
		StaticIPs: []string{d.NFSIP},
	})
	server.AddJob(d.createNFSServerJob())

	return server
}

func (d *DockerRegistry) createNFSServerJob() *enaml.InstanceJob {
	job := &enaml.InstanceJob{
		Name:    "debian_nfs_server",
		Release: registryReleaseName,
		Properties: &debian_nfs_server.DebianNfsServerJob{
			NfsServer: &debian_nfs_server.NfsServer{
				AllowFromEntries: d.RegistryIPs,
			},
		},
	}

	return job
}

func (d *DockerRegistry) CreateUpdate() (update enaml.Update) {
	update = enaml.Update{
		Canaries:        1,
		MaxInFlight:     3,
		CanaryWatchTime: "1000-60000",
		UpdateWatchTime: "1000-60000",
	}

	return
}
