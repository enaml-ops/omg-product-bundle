package cloudfoundry

import (
	"fmt"

	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/doppler"
	"github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/syslog_drain_binder"
)

//Doppler -
type Doppler struct {
	Config         *Config
	Metron         *Metron
	StatsdInjector *StatsdInjector
}

//NewDopplerPartition -
func NewDopplerPartition(config *Config) InstanceGroupCreator {
	return &Doppler{
		Config:         config,
		Metron:         NewMetron(config),
		StatsdInjector: NewStatsdInjector(nil),
	}
}

//ToInstanceGroup -
func (s *Doppler) ToInstanceGroup() (ig *enaml.InstanceGroup) {
	ig = &enaml.InstanceGroup{
		Name:      "doppler-partition",
		Instances: len(s.Config.DopplerIPs),
		VMType:    s.Config.DopplerVMType,
		AZs:       s.Config.AZs,
		Stemcell:  s.Config.StemcellName,
		Jobs: []enaml.InstanceJob{
			s.createDopplerJob(),
			s.Metron.CreateJob(),
			s.createSyslogDrainBinderJob(),
			s.StatsdInjector.CreateJob(),
		},
		Networks: []enaml.Network{
			enaml.Network{Name: s.Config.NetworkName, StaticIPs: s.Config.DopplerIPs},
		},
		Update: enaml.Update{
			MaxInFlight: 1,
		},
	}
	return
}

func (s *Doppler) createDopplerJob() enaml.InstanceJob {
	return enaml.InstanceJob{
		Name:    "doppler",
		Release: "cf",
		Properties: &doppler.DopplerJob{
			Doppler: &doppler.Doppler{
				Zone: s.Config.DopplerZone,
				MessageDrainBufferSize: s.Config.DopplerMessageDrainBufferSize,
			},
			DopplerEndpoint: &doppler.DopplerEndpoint{
				SharedSecret: s.Config.DopplerSharedSecret,
			},
			Loggregator: &doppler.Loggregator{
				Etcd: &doppler.Etcd{
					Machines: s.Config.EtcdMachines,
				},
			},
		},
	}
}

func (s *Doppler) createSyslogDrainBinderJob() enaml.InstanceJob {
	return enaml.InstanceJob{
		Name:    "syslog_drain_binder",
		Release: "cf",
		Properties: &syslog_drain_binder.SyslogDrainBinderJob{
			Ssl: &syslog_drain_binder.Ssl{
				SkipCertVerify: s.Config.SkipSSLCertVerify,
			},
			SystemDomain: s.Config.SystemDomain,
			Cc: &syslog_drain_binder.Cc{
				BulkApiPassword: s.Config.CCBuilkAPIPassword,
				SrvApiUri:       fmt.Sprintf("https://api.%s", s.Config.SystemDomain),
			},
			Loggregator: &syslog_drain_binder.Loggregator{
				Etcd: &syslog_drain_binder.Etcd{
					Machines: s.Config.EtcdMachines,
				},
			},
		},
	}
}
