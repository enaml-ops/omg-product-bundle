package cloudfoundry

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/doppler"
	"github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/syslog_drain_binder"
	"github.com/xchapter7x/lo"
)

//Doppler -
type Doppler struct {
	Config                 *Config
	AZs                    []string
	StemcellName           string
	VMTypeName             string
	NetworkName            string
	NetworkIPs             []string
	Metron                 *Metron
	StatsdInjector         *StatsdInjector
	Zone                   string
	MessageDrainBufferSize int
	SharedSecret           string
	CCBuilkAPIPassword     string
	EtcdMachines           []string
}

//NewDopplerPartition -
func NewDopplerPartition(c *cli.Context, config *Config) InstanceGrouper {
	return &Doppler{
		Config:         config,
		NetworkIPs:     c.StringSlice("doppler-ip"),
		VMTypeName:     c.String("doppler-vm-type"),
		Metron:         NewMetron(c),
		StatsdInjector: NewStatsdInjector(c),
		Zone:           c.String("doppler-zone"),
		MessageDrainBufferSize: c.Int("doppler-drain-buffer-size"),
		SharedSecret:           c.String("doppler-shared-secret"),
		CCBuilkAPIPassword:     c.String("cc-bulk-api-password"),
		EtcdMachines:           c.StringSlice("etcd-machine-ip"),
	}
}

//ToInstanceGroup -
func (s *Doppler) ToInstanceGroup() (ig *enaml.InstanceGroup) {
	ig = &enaml.InstanceGroup{
		Name:      "doppler-partition",
		Instances: len(s.NetworkIPs),
		VMType:    s.VMTypeName,
		AZs:       s.Config.AZs,
		Stemcell:  s.Config.StemcellName,
		Jobs: []enaml.InstanceJob{
			s.createDopplerJob(),
			s.Metron.CreateJob(),
			s.createSyslogDrainBinderJob(),
			s.StatsdInjector.CreateJob(),
		},
		Networks: []enaml.Network{
			enaml.Network{Name: s.Config.NetworkName, StaticIPs: s.NetworkIPs},
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
				Zone: s.Zone,
				MessageDrainBufferSize: s.MessageDrainBufferSize,
			},
			DopplerEndpoint: &doppler.DopplerEndpoint{
				SharedSecret: s.SharedSecret,
			},
			Loggregator: &doppler.Loggregator{
				Etcd: &doppler.Etcd{
					Machines: s.EtcdMachines,
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
				BulkApiPassword: s.CCBuilkAPIPassword,
				SrvApiUri:       fmt.Sprintf("https://api.%s", s.Config.SystemDomain),
			},
			Loggregator: &syslog_drain_binder.Loggregator{
				Etcd: &syslog_drain_binder.Etcd{
					Machines: s.EtcdMachines,
				},
			},
		},
	}
}

//HasValidValues - Check if the datastructure has valid fields
func (s *Doppler) HasValidValues() bool {

	lo.G.Debugf("checking '%s' for valid flags", "doppler")

	if len(s.NetworkIPs) <= 0 {
		lo.G.Debugf("could not find the correct number of network ips configured '%v' : '%v'", len(s.NetworkIPs), s.NetworkIPs)
	}
	if len(s.EtcdMachines) <= 0 {
		lo.G.Debugf("could not find the correct number of EtcdMachines configured '%v' : '%v'", len(s.EtcdMachines), s.EtcdMachines)
	}
	if s.VMTypeName == "" {
		lo.G.Debugf("could not find a valid vmtypename '%v'", s.VMTypeName)
	}
	if s.Zone == "" {
		lo.G.Debugf("could not find a valid zone '%v'", s.Zone)
	}
	if s.MessageDrainBufferSize <= 0 {
		lo.G.Debugf("could not find a valid MessageDrainBufferSize '%v'", s.MessageDrainBufferSize)
	}
	if s.SharedSecret == "" {
		lo.G.Debugf("could not find a valid SharedSecret '%v'", s.SharedSecret)
	}
	if s.CCBuilkAPIPassword == "" {
		lo.G.Debugf("could not find a valid CCBuilkAPIPassword '%v'", s.CCBuilkAPIPassword)
	}
	return (s.VMTypeName != "" &&
		len(s.NetworkIPs) > 0 &&
		s.Zone != "" &&
		s.MessageDrainBufferSize > 0 &&
		s.SharedSecret != "" &&
		s.CCBuilkAPIPassword != "" &&
		len(s.EtcdMachines) > 0 &&
		s.Metron.HasValidValues() &&
		s.StatsdInjector.HasValidValues())
}
