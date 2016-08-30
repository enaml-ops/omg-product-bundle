package cloudfoundry

import (
	"github.com/codegangsta/cli"
	"github.com/enaml-ops/enaml"
	etcdlib "github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/etcd"
	etcdmetricslib "github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/etcd_metrics_server"
	"github.com/xchapter7x/lo"
)

//Etcd -
type Etcd struct {
	Config             *Config
	VMTypeName         string
	NetworkIPs         []string
	PersistentDiskType string
	Metron             *Metron
	StatsdInjector     *StatsdInjector
}

//NewEtcdPartition -
func NewEtcdPartition(c *cli.Context, config *Config) (igf InstanceGrouper) {
	igf = &Etcd{
		Config:             config,
		NetworkIPs:         c.StringSlice("etcd-machine-ip"),
		VMTypeName:         c.String("etcd-vm-type"),
		PersistentDiskType: c.String("etcd-disk-type"),
		Metron:             NewMetron(config),
		StatsdInjector:     NewStatsdInjector(c),
	}
	return
}

//ToInstanceGroup -
func (s *Etcd) ToInstanceGroup() (ig *enaml.InstanceGroup) {
	ig = &enaml.InstanceGroup{
		Name:      "etcd_server-partition",
		Instances: len(s.NetworkIPs),
		VMType:    s.VMTypeName,
		AZs:       s.Config.AZs,
		Stemcell:  s.Config.StemcellName,
		Jobs: []enaml.InstanceJob{
			s.newEtcdJob(),
			s.newEtcdMetricsServerJob(),
			s.Metron.CreateJob(),
			s.StatsdInjector.CreateJob(),
		},
		Networks: []enaml.Network{
			enaml.Network{Name: s.Config.NetworkName, StaticIPs: s.NetworkIPs},
		},
		PersistentDiskType: s.PersistentDiskType,
		Update: enaml.Update{
			MaxInFlight: 1,
			Serial:      false,
		},
	}
	return
}

func (s *Etcd) newEtcdJob() enaml.InstanceJob {
	return enaml.InstanceJob{
		Name:    "etcd",
		Release: "cf",
		Properties: &etcdlib.EtcdJob{
			Etcd: &etcdlib.Etcd{
				PeerRequireSsl: false,
				RequireSsl:     false,
				Machines:       s.NetworkIPs,
			},
		},
	}
}

func (s *Etcd) newEtcdMetricsServerJob() enaml.InstanceJob {
	return enaml.InstanceJob{
		Name:    "etcd_metrics_server",
		Release: "cf",
		Properties: &etcdmetricslib.EtcdMetricsServerJob{
			EtcdMetricsServer: &etcdmetricslib.EtcdMetricsServer{
				Nats: &etcdmetricslib.Nats{
					Username: s.Config.NATSUser,
					Password: s.Config.NATSPassword,
					Machines: s.Config.NATSMachines,
					Port:     s.Config.NATSPort,
				},
			},
		},
	}
}

//HasValidValues -
func (s *Etcd) HasValidValues() bool {

	lo.G.Debugf("checking '%s' for valid flags", "etcd")

	if len(s.NetworkIPs) <= 0 {
		lo.G.Debugf("could not find the correct number of network ips configured '%v' : '%v'", len(s.NetworkIPs), s.NetworkIPs)
	}
	if s.VMTypeName == "" {
		lo.G.Debugf("could not find a valid vmtypename '%v'", s.VMTypeName)
	}
	if s.PersistentDiskType == "" {
		lo.G.Debugf("could not find a valid PersistentDiskType '%v'", s.PersistentDiskType)
	}
	return (s.VMTypeName != "" &&
		len(s.NetworkIPs) > 0 &&
		s.PersistentDiskType != "")
}
