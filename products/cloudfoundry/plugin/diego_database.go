package cloudfoundry

import (
	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/bbs"
	"github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/etcd"
	"github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/plugin/config"
)

const diegoDatabaseIGName = "diego_database-partition"

type diegoDatabase struct {
	Config         *config.Config
	ConsulAgent    *ConsulAgent
	StatsdInjector *StatsdInjector
	Metron         *Metron
}

func NewDiegoDatabasePartition(config *config.Config) InstanceGroupCreator {

	return &diegoDatabase{
		Config:         config,
		ConsulAgent:    NewConsulAgent([]string{"bbs", "etcd"}, config),
		Metron:         NewMetron(config),
		StatsdInjector: NewStatsdInjector(nil),
	}
}

func (s *diegoDatabase) ToInstanceGroup() (ig *enaml.InstanceGroup) {
	ig = &enaml.InstanceGroup{
		Name:               diegoDatabaseIGName,
		Lifecycle:          "service",
		Instances:          len(s.Config.DiegoDBIPs),
		VMType:             s.Config.DiegoDBVMType,
		AZs:                s.Config.AZs,
		PersistentDiskType: s.Config.DiegoDBPersistentDiskType,
		Stemcell:           s.Config.StemcellName,
		Networks: []enaml.Network{
			{Name: s.Config.NetworkName, StaticIPs: s.Config.DiegoDBIPs},
		},
		Update: enaml.Update{
			MaxInFlight: 1,
			Serial:      true,
		},
	}

	ig.AddJob(&enaml.InstanceJob{
		Name:       "etcd",
		Release:    EtcdReleaseName,
		Properties: s.newEtcd(),
	})
	ig.AddJob(&enaml.InstanceJob{
		Name:       "bbs",
		Release:    DiegoReleaseName,
		Properties: s.newBBS(),
	})
	ig.AddJob(&enaml.InstanceJob{
		Name:       "consul_agent",
		Release:    CFReleaseName,
		Properties: s.ConsulAgent.CreateJob().Properties,
	})
	ig.AddJob(&enaml.InstanceJob{
		Name:       "statsd-injector",
		Release:    CFReleaseName,
		Properties: s.StatsdInjector.CreateJob().Properties,
	})
	ig.AddJob(&enaml.InstanceJob{
		Name:       "metron_agent",
		Release:    CFReleaseName,
		Properties: s.Metron.CreateJob().Properties,
	})
	return
}

func (s *diegoDatabase) newBBS() (dbdiego *bbs.BbsJob) {
	var keyname = "key1"
	return &bbs.BbsJob{
		Diego: &bbs.Diego{
			Bbs: &bbs.Bbs{
				RequireSsl:     false,
				CaCert:         s.Config.BBSCACert,
				ServerCert:     s.Config.BBSServerCert,
				ServerKey:      s.Config.BBSServerKey,
				ActiveKeyLabel: keyname,
				EncryptionKeys: []map[string]string{
					{
						"label":      keyname,
						"passphrase": s.Config.DiegoDBPassphrase,
					},
				},
				Auctioneer: &bbs.Auctioneer{
					ApiUrl: "http://auctioneer.service.cf.internal:9016",
				},
				Etcd: s.newBBSEtcd(),
			},
		},
	}
}

func (s *diegoDatabase) newEtcd() *etcd.EtcdJob {
	return &etcd.EtcdJob{
		Etcd: &etcd.Etcd{
			CaCert:                 s.Config.BBSCACert,
			ServerCert:             s.Config.EtcdServerCert,
			ServerKey:              s.Config.EtcdServerKey,
			ClientCert:             s.Config.EtcdClientCert,
			ClientKey:              s.Config.EtcdClientKey,
			PeerCaCert:             s.Config.BBSCACert,
			PeerCert:               s.Config.EtcdPeerCert,
			PeerKey:                s.Config.EtcdPeerKey,
			AdvertiseUrlsDnsSuffix: "etcd.service.cf.internal",
			Machines:               []string{"etcd.service.cf.internal"},
			Cluster: []map[string]interface{}{
				{
					"name":      diegoDatabaseIGName,
					"instances": len(s.Config.DiegoDBIPs),
				},
			},
		},
	}
}

func (s *diegoDatabase) newBBSEtcd() (dbetcd *bbs.Etcd) {
	dbetcd = &bbs.Etcd{
		CaCert:     s.Config.BBSCACert,
		ClientCert: s.Config.EtcdClientCert,
		ClientKey:  s.Config.EtcdClientKey,
		Machines:   []string{"etcd.service.cf.internal"},
	}
	return
}
