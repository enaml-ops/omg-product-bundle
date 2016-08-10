package dockerregistry

import (
	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-cli/plugins/products/bosh-init/enaml-gen/registry"
	"github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/debian_nfs_server"
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
	PublicIP                 string
	NFSServerType            string
	NFSDiskType              string
	NFSIP                    string
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
	return server
}

func (d *DockerRegistry) createRegistryJob() *enaml.InstanceJob {
	job := &enaml.InstanceJob{
		Name:       "registry",
		Release:    "docker-registry",
		Properties: &registry.RegistryJob{},
	}

	return job
}

func (d *DockerRegistry) CreateProxyInstanceGroup() *enaml.InstanceGroup {
	server := &enaml.InstanceGroup{}
	return server
}

func (d *DockerRegistry) CreateNFSServerInstanceGroup() *enaml.InstanceGroup {
	server := &enaml.InstanceGroup{
		Name:               "nfs-server",
		Instances:          1,
		Stemcell:           d.StemcellAlias,
		VMType:             d.NFSServerType,
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
		Name:    "nfs-server",
		Release: "docker-registry",
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
