package cloudfoundry

import (
	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/garden"
	"github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/rep"
	"github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/plugin/config"
)

type diegoCell struct {
	Config         *config.Config
	ConsulAgent    *ConsulAgent
	StatsdInjector *StatsdInjector
	Metron         *Metron
}

func NewDiegoCellPartition(config *config.Config) InstanceGroupCreator {

	return &diegoCell{
		Config:         config,
		ConsulAgent:    NewConsulAgent([]string{}, config),
		Metron:         NewMetron(config),
		StatsdInjector: NewStatsdInjector(nil),
	}
}

func (s *diegoCell) ToInstanceGroup() (ig *enaml.InstanceGroup) {
	ig = &enaml.InstanceGroup{
		Name:               "diego_cell-partition",
		Lifecycle:          "service",
		Instances:          len(s.Config.DiegoCellIPs),
		VMType:             s.Config.DiegoCellVMType,
		AZs:                s.Config.AZs,
		PersistentDiskType: s.Config.DiegoCellPersistentDiskType,
		Stemcell:           s.Config.StemcellName,
		Networks: []enaml.Network{
			{Name: s.Config.NetworkName, StaticIPs: s.Config.DiegoCellIPs},
		},
		Update: enaml.Update{
			MaxInFlight: 1,
		},
	}

	ig.AddJob(&enaml.InstanceJob{
		Name:       "rep",
		Release:    DiegoReleaseName,
		Properties: s.newRDiego(),
	})
	ig.AddJob(&enaml.InstanceJob{
		Name:       "consul_agent",
		Release:    CFReleaseName,
		Properties: s.ConsulAgent.CreateJob().Properties,
	})
	ig.AddJob(&enaml.InstanceJob{
		Name:       "cflinuxfs2-rootfs-setup",
		Release:    CFLinuxFSReleaseName,
		Properties: struct{}{},
	})
	ig.AddJob(&enaml.InstanceJob{
		Name:       "garden",
		Release:    GardenReleaseName,
		Properties: s.newGarden(),
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

func (s *diegoCell) newGarden() (gardenLinux *garden.GardenJob) {
	gardenLinux = &garden.GardenJob{
		Garden: &garden.Garden{
			AllowHostAccess:     false,
			PersistentImageList: []string{"/var/vcap/packages/rootfs_cflinuxfs2/rootfs"},
			NetworkPool:         "10.254.0.0/22",
			DenyNetworks:        []string{"0.0.0.0/0"},
			NetworkMtu:          1454,
		},
	}
	return
}

func (s *diegoCell) newRDiego() (rdiego *rep.RepJob) {
	rdiego = &rep.RepJob{
		Diego: &rep.Diego{
			Executor: &rep.Executor{
				PostSetupHook: `sh -c "rm -f /home/vcap/app/.java-buildpack.log /home/vcap/app/**/.java-buildpack.log"`,
				PostSetupUser: "root",
			},
			Rep: &rep.Rep{
				Bbs: &rep.Bbs{
					ApiLocation: defaultBBSAPILocation,
					CaCert:      s.Config.BBSCACert,
					ClientCert:  s.Config.BBSClientCert,
					ClientKey:   s.Config.BBSClientKey,
				},
				PreloadedRootfses: []string{
					"cflinuxfs2:/var/vcap/packages/cflinuxfs2/rootfs",
				},
				Zone: s.Config.MetronZone,
			},
		},
	}
	return
}
