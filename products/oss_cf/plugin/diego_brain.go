package cloudfoundry

import (
	"fmt"
	"strings"

	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-product-bundle/products/oss_cf/enaml-gen/auctioneer"
	"github.com/enaml-ops/omg-product-bundle/products/oss_cf/enaml-gen/cc_uploader"
	"github.com/enaml-ops/omg-product-bundle/products/oss_cf/enaml-gen/converger"
	"github.com/enaml-ops/omg-product-bundle/products/oss_cf/enaml-gen/file_server"
	"github.com/enaml-ops/omg-product-bundle/products/oss_cf/enaml-gen/nsync"
	"github.com/enaml-ops/omg-product-bundle/products/oss_cf/enaml-gen/route_emitter"
	"github.com/enaml-ops/omg-product-bundle/products/oss_cf/enaml-gen/ssh_proxy"
	"github.com/enaml-ops/omg-product-bundle/products/oss_cf/enaml-gen/stager"
	"github.com/enaml-ops/omg-product-bundle/products/oss_cf/enaml-gen/tps"
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
	ig.AddJob(d.newCCUploader())
	ig.AddJob(d.newConverger())
	ig.AddJob(d.newFileServer())
	ig.AddJob(d.newNsync())
	ig.AddJob(d.newRouteEmitter())
	ig.AddJob(d.newSSHProxy())
	ig.AddJob(d.newStager())
	ig.AddJob(d.newTPS())
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

func (d *diegoBrain) newCCUploader() *enaml.InstanceJob {
	return &enaml.InstanceJob{
		Name:    "cc_uploader",
		Release: DiegoReleaseName,
		Properties: &cc_uploader.CcUploaderJob{
			Capi: &cc_uploader.Capi{
				CcUploader: &cc_uploader.CcUploader{
					Cc: &cc_uploader.Cc{
						JobPollingIntervalInSeconds: d.Config.CCUploaderJobPollInterval,
					},
				},
			},
			Diego: &cc_uploader.Diego{
				Ssl: &cc_uploader.Ssl{SkipCertVerify: d.Config.SkipSSLCertVerify},
			},
		},
	}
}

func (d *diegoBrain) newConverger() *enaml.InstanceJob {
	return &enaml.InstanceJob{
		Name:    "converger",
		Release: DiegoReleaseName,
		Properties: &converger.ConvergerJob{
			Diego: &converger.Diego{
				Converger: &converger.Converger{
					Bbs: &converger.Bbs{
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
			Diego: &file_server.Diego{
				FileServer: &file_server.FileServer{},
			},
		},
	}
}

func (d *diegoBrain) newNsync() *enaml.InstanceJob {
	return &enaml.InstanceJob{
		Name:    "nsync",
		Release: DiegoReleaseName,
		Properties: &nsync.NsyncJob{
			Capi: &nsync.Capi{
				Nsync: &nsync.Nsync{
					Cc: &nsync.Cc{
						BaseUrl:                  prefixSystemDomain(d.Config.SystemDomain, "api"),
						BasicAuthUsername:        d.Config.CCInternalAPIUser,
						BasicAuthPassword:        d.Config.CCInternalAPIPassword,
						PollingIntervalInSeconds: d.Config.CCUploaderJobPollInterval,
					},
					Bbs: &nsync.Bbs{
						ApiLocation: defaultBBSAPILocation,
						CaCert:      d.Config.BBSCACert,
						ClientCert:  d.Config.BBSClientCert,
						ClientKey:   d.Config.BBSClientKey,
					},
				},
			},
			Diego: &nsync.Diego{
				Ssl: &nsync.Ssl{SkipCertVerify: d.Config.SkipSSLCertVerify},
			},
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

func (d *diegoBrain) newStager() *enaml.InstanceJob {
	return &enaml.InstanceJob{
		Name:    "stager",
		Release: CFReleaseName,
		Properties: &stager.StagerJob{
			Capi: &stager.Capi{
				Stager: &stager.Stager{
					Bbs: &stager.Bbs{
						ApiLocation: defaultBBSAPILocation,
						CaCert:      d.Config.BBSCACert,
						ClientCert:  d.Config.BBSClientCert,
						ClientKey:   d.Config.BBSClientKey,
						RequireSsl:  d.Config.BBSRequireSSL,
					},
					Cc: &stager.Cc{
						BasicAuthUsername: d.Config.CCInternalAPIUser,
						BasicAuthPassword: d.Config.CCInternalAPIPassword,
						ExternalPort:      d.Config.CCExternalPort,
					},
				},
			},
			Diego: &stager.Diego{
				Ssl: &stager.Ssl{SkipCertVerify: d.Config.SkipSSLCertVerify},
			},
		},
	}
}

func (d *diegoBrain) newTPS() *enaml.InstanceJob {
	return &enaml.InstanceJob{
		Name:    "tps",
		Release: CFReleaseName,
		Properties: &tps.TpsJob{
			Capi: &tps.Capi{
				Tps: &tps.Tps{
					TrafficControllerUrl: fmt.Sprintf("wss://doppler.%s:%d", d.Config.SystemDomain, d.Config.LoggregatorPort),
					Bbs: &tps.Bbs{
						ApiLocation: defaultBBSAPILocation,
						CaCert:      d.Config.BBSCACert,
						ClientCert:  d.Config.BBSClientCert,
						ClientKey:   d.Config.BBSClientKey,
						RequireSsl:  d.Config.BBSRequireSSL,
					},
					Cc: &tps.Cc{
						BasicAuthUsername: d.Config.CCInternalAPIUser,
						BasicAuthPassword: d.Config.CCInternalAPIPassword,
						ExternalPort:      d.Config.CCExternalPort,
					},
				},
			},
			Diego: &tps.Diego{
				Ssl: &tps.Ssl{SkipCertVerify: d.Config.SkipSSLCertVerify},
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
