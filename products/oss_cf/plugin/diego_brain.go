package cloudfoundry

import (
	"fmt"
	"strings"

	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-product-bundle/products/oss_cf/enaml-gen/auctioneer"
	"github.com/enaml-ops/omg-product-bundle/products/oss_cf/enaml-gen/file_server"
	"github.com/enaml-ops/omg-product-bundle/products/oss_cf/enaml-gen/route_emitter"
	"github.com/enaml-ops/omg-product-bundle/products/oss_cf/enaml-gen/ssh_proxy"
	"github.com/enaml-ops/omg-product-bundle/products/oss_cf/plugin/config"
	"github.com/enaml-ops/pluginlib/pluginutil"
	"github.com/xchapter7x/lo"
)

type diegoBrain struct {
	Config      *config.Config
	ConsulAgent *ConsulAgent
	Metron      *Metron
	Statsd      *StatsdInjector
}

func NewDiegoBrainPartition(config *config.Config) InstanceGroupCreator {

	return &diegoBrain{
		Config: config,

		ConsulAgent: NewConsulAgent([]string{}, config),
		Metron:      NewMetron(config),
		Statsd:      NewStatsdInjector(nil),
	}
}

func (d *diegoBrain) ToInstanceGroup() *enaml.InstanceGroup {
	ig := &enaml.InstanceGroup{
		Name:               "diego_brain-partition",
		Instances:          len(d.Config.DiegoBrainIPs),
		VMType:             d.Config.DiegoBrainVMType,
		AZs:                d.Config.AZs,
		PersistentDiskType: d.Config.DiegoBrainPersistentDiskType,
		Stemcell:           d.Config.StemcellName,
		Networks: []enaml.Network{
			{Name: d.Config.NetworkName, StaticIPs: d.Config.DiegoBrainIPs},
		},
		Update: enaml.Update{
			MaxInFlight: 1,
		},
	}
	consulJob := d.ConsulAgent.CreateJob()
	metronJob := d.Metron.CreateJob()
	statsdJob := d.Statsd.CreateJob()

	ig.AddJob(d.newAuctioneer())
	ig.AddJob(d.newFileServer())
	ig.AddJob(d.newRouteEmitter())
	ig.AddJob(d.newSSHProxy())
	ig.AddJob(&consulJob)
	ig.AddJob(&metronJob)
	ig.AddJob(&statsdJob)
	return ig
}

func (d *diegoBrain) newAuctioneer() *enaml.InstanceJob {
	return &enaml.InstanceJob{
		Name:    "auctioneer",
		Release: DiegoReleaseName,
		Properties: &auctioneer.AuctioneerJob{
			Diego: &auctioneer.Diego{
				Auctioneer: &auctioneer.Auctioneer{
					Bbs: &auctioneer.Bbs{
						ApiLocation: defaultBBSAPILocation,
						CaCert:      d.Config.BBSCACert,
						ClientCert:  d.Config.BBSClientCert,
						ClientKey:   d.Config.BBSClientKey,
					},
				},
			},
		},
	}
}

func (d *diegoBrain) newFileServer() *enaml.InstanceJob {
	return &enaml.InstanceJob{
		Name:    "file_server",
		Release: DiegoReleaseName,
		Properties: &file_server.FileServerJob{
			Diego: &file_server.Diego{},
		},
	}
}

func (d *diegoBrain) newRouteEmitter() *enaml.InstanceJob {
	return &enaml.InstanceJob{
		Name:    "route_emitter",
		Release: DiegoReleaseName,
		Properties: &route_emitter.RouteEmitterJob{
			Diego: &route_emitter.Diego{
				RouteEmitter: &route_emitter.RouteEmitter{
					Bbs: &route_emitter.Bbs{
						ApiLocation: defaultBBSAPILocation,
						CaCert:      d.Config.BBSCACert,
						ClientCert:  d.Config.BBSClientCert,
						ClientKey:   d.Config.BBSClientKey,
						RequireSsl:  d.Config.BBSRequireSSL,
					},
					Nats: &route_emitter.Nats{
						User:     d.Config.NATSUser,
						Password: d.Config.NATSPassword,
						Port:     d.Config.NATSPort,
						Machines: d.Config.NATSMachines,
					},
				},
			},
		},
	}
}

func (d *diegoBrain) newSSHProxy() *enaml.InstanceJob {
	_, privateKey, err := pluginutil.GenerateKeys()
	if err != nil {
		lo.G.Error("couldn't generate private key for SSH proxy")
		return nil
	}

	return &enaml.InstanceJob{
		Name:    "ssh_proxy",
		Release: DiegoReleaseName,
		Properties: &ssh_proxy.SshProxyJob{
			Diego: &ssh_proxy.Diego{
				Ssl: &ssh_proxy.Ssl{SkipCertVerify: d.Config.SkipSSLCertVerify},
				SshProxy: &ssh_proxy.SshProxy{
					Bbs: &ssh_proxy.Bbs{
						ApiLocation: defaultBBSAPILocation,
						CaCert:      d.Config.BBSCACert,
						ClientCert:  d.Config.BBSClientCert,
						ClientKey:   d.Config.BBSClientKey,
						RequireSsl:  d.Config.BBSRequireSSL,
					},
					Cc: &ssh_proxy.Cc{
						ExternalPort: d.Config.CCExternalPort,
					},
					EnableCfAuth:    d.Config.AllowSSHAccess,
					EnableDiegoAuth: d.Config.AllowSSHAccess,
					UaaSecret:       d.Config.SSHProxyClientSecret,
					UaaTokenUrl:     prefixSystemDomain(d.Config.SystemDomain, "uaa") + "/oauth/token",
					HostKey:         privateKey,
				},
			},
		},
	}
}

// prefixSystemDomain adds a prefix to the system domain.
// For example:
//     prefixSystemDomain("https://sys.yourdomain.com", "uaa")
// would return 'https://uaa.sys.yourdomain.com'.
func prefixSystemDomain(domain, prefix string) string {
	d := domain
	// strip leading https:// if necessary
	if strings.HasPrefix(d, "https://") {
		d = d[len("https://"):]
	}
	return fmt.Sprintf("https://%s.%s", prefix, d)
}
