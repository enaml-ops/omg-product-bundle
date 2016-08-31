package cloudfoundry

import (
	"github.com/enaml-ops/enaml"
	nfslib "github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/debian_nfs_server"
)

//NFS -
type NFS struct {
	Config         *Config
	Metron         *Metron
	StatsdInjector *StatsdInjector
}

//NewNFSPartition -
func NewNFSPartition(config *Config) (igf InstanceGroupCreator) {
	igf = &NFS{
		Config:         config,
		Metron:         NewMetron(config),
		StatsdInjector: NewStatsdInjector(nil),
	}
	return
}

//ToInstanceGroup -
func (s *NFS) ToInstanceGroup() (ig *enaml.InstanceGroup) {
	ig = &enaml.InstanceGroup{
		Name:      "nfs_server-partition",
		Instances: len(s.Config.NFSIPs),
		VMType:    s.Config.NFSVMType,
		AZs:       s.Config.AZs,
		Stemcell:  s.Config.StemcellName,
		Jobs: []enaml.InstanceJob{
			s.newNFSJob(),
			s.Metron.CreateJob(),
			s.StatsdInjector.CreateJob(),
		},
		Networks: []enaml.Network{
			enaml.Network{Name: s.Config.NetworkName, StaticIPs: s.Config.NFSIPs},
		},
		PersistentDiskType: s.Config.NFSPersistentDiskType,
		Update: enaml.Update{
			MaxInFlight: 1,
			Serial:      true,
		},
	}
	return
}

func (s *NFS) newNFSJob() enaml.InstanceJob {

	return enaml.InstanceJob{
		Name:    "debian_nfs_server",
		Release: "cf",
		Properties: &nfslib.DebianNfsServerJob{
			NfsServer: &nfslib.NfsServer{
				AllowFromEntries: s.Config.NFSAllowFromNetworkCIDR,
			},
		},
	}
}
