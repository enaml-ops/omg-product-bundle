package cloudfoundry

import (
	"github.com/enaml-ops/enaml"
	nfslib "github.com/enaml-ops/omg-product-bundle/products/oss_cf/enaml-gen/debian_nfs_server"
	"github.com/enaml-ops/omg-product-bundle/products/oss_cf/plugin/config"
)

//NFS -
type NFS struct {
	Config         *config.Config
	Metron         *Metron
	StatsdInjector *StatsdInjector
}

//NewNFSPartition -
func NewNFSPartition(config *config.Config) (igf InstanceGroupCreator) {
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
		Instances: 1,
		VMType:    s.Config.NFSVMType,
		AZs:       s.Config.AZs,
		Stemcell:  s.Config.StemcellName,
		Jobs: []enaml.InstanceJob{
			s.newNFSJob(),
			s.Metron.CreateJob(),
			s.StatsdInjector.CreateJob(),
		},
		Networks: []enaml.Network{
			enaml.Network{Name: s.Config.NetworkName, StaticIPs: []string{s.Config.NFSIP}},
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
