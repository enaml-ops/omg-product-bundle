package cloudfoundry

import (
	"github.com/codegangsta/cli"
	"github.com/enaml-ops/enaml"
	nfslib "github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/debian_nfs_server"
	"github.com/xchapter7x/lo"
)

//NFS -
type NFS struct {
	Config               *Config
	VMTypeName           string
	NetworkIPs           []string
	PersistentDiskType   string
	AllowFromNetworkCIDR []string
	Metron               *Metron
	StatsdInjector       *StatsdInjector
}

//NewNFSPartition -
func NewNFSPartition(c *cli.Context, config *Config) (igf InstanceGrouper) {
	igf = &NFS{
		Config:               config,
		NetworkIPs:           c.StringSlice("nfs-ip"),
		VMTypeName:           c.String("nfs-vm-type"),
		PersistentDiskType:   c.String("nfs-disk-type"),
		AllowFromNetworkCIDR: c.StringSlice("nfs-allow-from-network-cidr"),
		Metron:               NewMetron(c),
		StatsdInjector:       NewStatsdInjector(c),
	}
	return
}

//ToInstanceGroup -
func (s *NFS) ToInstanceGroup() (ig *enaml.InstanceGroup) {
	ig = &enaml.InstanceGroup{
		Name:      "nfs_server-partition",
		Instances: len(s.NetworkIPs),
		VMType:    s.VMTypeName,
		AZs:       s.Config.AZs,
		Stemcell:  s.Config.StemcellName,
		Jobs: []enaml.InstanceJob{
			s.newNFSJob(),
			s.Metron.CreateJob(),
			s.StatsdInjector.CreateJob(),
		},
		Networks: []enaml.Network{
			enaml.Network{Name: s.Config.NetworkName, StaticIPs: s.NetworkIPs},
		},
		PersistentDiskType: s.PersistentDiskType,
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
				AllowFromEntries: s.AllowFromNetworkCIDR,
			},
		},
	}
}

//HasValidValues -
func (s *NFS) HasValidValues() bool {

	lo.G.Debugf("checking '%s' for valid flags", "nfs")

	if len(s.NetworkIPs) <= 0 {
		lo.G.Debugf("could not find the correct number of network ips configured '%v' : '%v'", len(s.NetworkIPs), s.NetworkIPs)
	}
	if len(s.AllowFromNetworkCIDR) <= 0 {
		lo.G.Debugf("could not find the correct number of AllowFromNetworkCIDR configured '%v' : '%v'", len(s.AllowFromNetworkCIDR), s.AllowFromNetworkCIDR)
	}
	if s.VMTypeName == "" {
		lo.G.Debugf("could not find a valid vmtypename '%v'", s.VMTypeName)
	}
	if s.PersistentDiskType == "" {
		lo.G.Debugf("could not find a valid PersistentDiskType '%v'", s.PersistentDiskType)
	}

	return (s.VMTypeName != "" &&
		len(s.NetworkIPs) > 0 &&
		s.PersistentDiskType != "" &&
		len(s.AllowFromNetworkCIDR) > 0)
}
