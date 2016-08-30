package cloudfoundry

import (
	"github.com/codegangsta/cli"
	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/garden"
	"github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/rep"
	"github.com/xchapter7x/lo"
)

type diegoCell struct {
	Config             *Config
	context            *cli.Context
	VMTypeName         string
	PersistentDiskType string
	NetworkIPs         []string
	ConsulAgent        *ConsulAgent
	StatsdInjector     *StatsdInjector
	Metron             *Metron
	DiegoBrain         *diegoBrain
}

func NewDiegoCellPartition(c *cli.Context, config *Config) InstanceGrouper {

	return &diegoCell{
		context:            c,
		Config:             config,
		VMTypeName:         c.String("diego-cell-vm-type"),
		PersistentDiskType: c.String("diego-cell-disk-type"),
		NetworkIPs:         c.StringSlice("diego-cell-ip"),
		ConsulAgent:        NewConsulAgent(c, []string{}, config),
		Metron:             NewMetron(c),
		StatsdInjector:     NewStatsdInjector(c),
		DiegoBrain:         NewDiegoBrainPartition(c, config).(*diegoBrain),
	}
}

func (s *diegoCell) ToInstanceGroup() (ig *enaml.InstanceGroup) {
	ig = &enaml.InstanceGroup{
		Name:               "diego_cell-partition",
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

func (s *diegoCell) HasValidValues() bool {
	lo.G.Debugf("checking %v for valid flags", "diego cell")
	validStrings := hasValidStringFlags(s.context, []string{
		"stemcell-name",
		"diego-cell-vm-type",
		"diego-cell-disk-type",
		"network",
	})
	return validStrings
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
					CaCert:      s.DiegoBrain.BBSCACert,
					ClientCert:  s.DiegoBrain.BBSClientCert,
					ClientKey:   s.DiegoBrain.BBSClientKey,
				},
				PreloadedRootfses: []string{
					"cflinuxfs2:/var/vcap/packages/cflinuxfs2/rootfs",
				},
				Zone: s.Metron.Zone,
			},
		},
	}
	return
}
