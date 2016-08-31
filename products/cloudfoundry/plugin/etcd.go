package cloudfoundry

import (
	"github.com/enaml-ops/enaml"
	etcdlib "github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/etcd"
	etcdmetricslib "github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/etcd_metrics_server"
	"github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/plugin/config"
)

//Etcd -
type Etcd struct {
	Config         *config.Config
	Metron         *Metron
	StatsdInjector *StatsdInjector
}

//NewEtcdPartition -
func NewEtcdPartition(config *config.Config) (igf InstanceGroupCreator) {
	igf = &Etcd{
		Config:         config,
		Metron:         NewMetron(config),
		StatsdInjector: NewStatsdInjector(nil),
	}
	return
}

//ToInstanceGroup -
func (s *Etcd) ToInstanceGroup() (ig *enaml.InstanceGroup) {
	ig = &enaml.InstanceGroup{
		Name:      "etcd_server-partition",
		Instances: len(s.Config.EtcdMachines),
		VMType:    s.Config.EtcdVMType,
		AZs:       s.Config.AZs,
		Stemcell:  s.Config.StemcellName,
		Jobs: []enaml.InstanceJob{
			s.newEtcdJob(),
			s.newEtcdMetricsServerJob(),
			s.Metron.CreateJob(),
			s.StatsdInjector.CreateJob(),
		},
		Networks: []enaml.Network{
			enaml.Network{Name: s.Config.NetworkName, StaticIPs: s.Config.EtcdMachines},
		},
		PersistentDiskType: s.Config.EtcdPersistentDiskType,
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
				Machines:       s.Config.EtcdMachines,
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
