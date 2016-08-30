package cloudfoundry

import (
	"github.com/codegangsta/cli"
	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/bbs"
	"github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/etcd"
	"github.com/enaml-ops/pluginlib/util"
	"github.com/xchapter7x/lo"
)

const diegoDatabaseIGName = "diego_database-partition"

type diegoDatabase struct {
	context            *cli.Context
	Config             *Config
	Passphrase         string
	VMTypeName         string
	PersistentDiskType string
	NetworkIPs         []string
	EtcdServerCert     string
	EtcdServerKey      string
	EtcdClientCert     string
	EtcdClientKey      string
	EtcdPeerCert       string
	EtcdPeerKey        string
	ConsulAgent        *ConsulAgent
	StatsdInjector     *StatsdInjector
	Metron             *Metron
}

func NewDiegoDatabasePartition(c *cli.Context, config *Config) InstanceGrouper {

	etcdServerCert, err := pluginutil.LoadResourceFromContext(c, "etcd-server-cert")
	if err != nil {
		lo.G.Fatalf("etcd server cert: %s\n", err.Error())
	}

	etcdServerKey, err := pluginutil.LoadResourceFromContext(c, "etcd-server-key")
	if err != nil {
		lo.G.Fatalf("etcd server key: %s\n", err.Error())
	}

	etcdClientCert, err := pluginutil.LoadResourceFromContext(c, "etcd-client-cert")
	if err != nil {
		lo.G.Fatalf("etcd client cert: %s\n", err.Error())
	}

	etcdClientKey, err := pluginutil.LoadResourceFromContext(c, "etcd-client-key")
	if err != nil {
		lo.G.Fatalf("etcd client key: %s\n", err.Error())
	}

	etcdPeerCert, err := pluginutil.LoadResourceFromContext(c, "etcd-peer-cert")
	if err != nil {
		lo.G.Fatalf("etcd peer cert: %s\n", err.Error())
	}

	etcdPeerKey, err := pluginutil.LoadResourceFromContext(c, "etcd-peer-key")
	if err != nil {
		lo.G.Fatalf("etcd peer key: %s\n", err.Error())
	}

	return &diegoDatabase{
		context:            c,
		Config:             config,
		VMTypeName:         c.String("diego-db-vm-type"),
		PersistentDiskType: c.String("diego-db-disk-type"),
		NetworkIPs:         c.StringSlice("diego-db-ip"),
		Passphrase:         c.String("diego-db-passphrase"),
		EtcdServerCert:     etcdServerCert,
		EtcdServerKey:      etcdServerKey,
		EtcdClientCert:     etcdClientCert,
		EtcdClientKey:      etcdClientKey,
		EtcdPeerCert:       etcdPeerCert,
		EtcdPeerKey:        etcdPeerKey,
		ConsulAgent:        NewConsulAgent([]string{"bbs", "etcd"}, config),
		Metron:             NewMetron(config),
		StatsdInjector:     NewStatsdInjector(c),
	}
}

func (s *diegoDatabase) ToInstanceGroup() (ig *enaml.InstanceGroup) {
	ig = &enaml.InstanceGroup{
		Name:               diegoDatabaseIGName,
		Lifecycle:          "service",
		Instances:          len(s.NetworkIPs),
		VMType:             s.VMTypeName,
		AZs:                s.Config.AZs,
		PersistentDiskType: s.PersistentDiskType,
		Stemcell:           s.Config.StemcellName,
		Networks: []enaml.Network{
			{Name: s.Config.NetworkName, StaticIPs: s.NetworkIPs},
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

func (s *diegoDatabase) HasValidValues() bool {
	lo.G.Debugf("checking %v for valid flags", "diego database")
	validStrings := hasValidStringFlags(s.context, []string{
		"bbs-server-ca-cert",
		"bbs-server-cert",
		"bbs-server-key",
		"etcd-server-cert",
		"etcd-server-key",
		"etcd-client-cert",
		"etcd-client-key",
		"etcd-peer-cert",
		"etcd-peer-key",
		"system-domain",
		"stemcell-name",
		"diego-db-vm-type",
		"diego-db-disk-type",
		"network",
		"diego-db-passphrase",
	})
	return validStrings
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
						"passphrase": s.Passphrase,
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
			ServerCert:             s.EtcdServerCert,
			ServerKey:              s.EtcdServerKey,
			ClientCert:             s.EtcdClientCert,
			ClientKey:              s.EtcdClientKey,
			PeerCaCert:             s.Config.BBSCACert,
			PeerCert:               s.EtcdPeerCert,
			PeerKey:                s.EtcdPeerKey,
			AdvertiseUrlsDnsSuffix: "etcd.service.cf.internal",
			Machines:               []string{"etcd.service.cf.internal"},
			Cluster: []map[string]interface{}{
				{
					"name":      diegoDatabaseIGName,
					"instances": len(s.NetworkIPs),
				},
			},
		},
	}
}

func (s *diegoDatabase) newBBSEtcd() (dbetcd *bbs.Etcd) {
	dbetcd = &bbs.Etcd{
		CaCert:     s.Config.BBSCACert,
		ClientCert: s.EtcdClientCert,
		ClientKey:  s.EtcdClientKey,
		Machines:   []string{"etcd.service.cf.internal"},
	}
	return
}
