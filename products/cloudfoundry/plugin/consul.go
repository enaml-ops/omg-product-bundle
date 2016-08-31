package cloudfoundry

import "github.com/enaml-ops/enaml"

// Consul -
type Consul struct {
	Config         *Config
	ConsulAgent    *ConsulAgent
	Metron         *Metron
	StatsdInjector *StatsdInjector
}

//NewConsulPartition -
func NewConsulPartition(config *Config) InstanceGroupCreator {
	return &Consul{
		Config:         config,
		ConsulAgent:    NewConsulAgentServer(config),
		Metron:         NewMetron(config),
		StatsdInjector: NewStatsdInjector(nil),
	}
}

//ToInstanceGroup -
func (s *Consul) ToInstanceGroup() (ig *enaml.InstanceGroup) {
	ig = &enaml.InstanceGroup{
		Name:      "consul-partition",
		Instances: len(s.Config.ConsulIPs),
		VMType:    s.Config.ConsulVMType,
		AZs:       s.Config.AZs,
		Stemcell:  s.Config.StemcellName,
		Jobs: []enaml.InstanceJob{
			s.ConsulAgent.CreateJob(),
			s.Metron.CreateJob(),
			s.StatsdInjector.CreateJob(),
		},
		Networks: []enaml.Network{
			enaml.Network{Name: s.Config.NetworkName, StaticIPs: s.Config.ConsulIPs},
		},
		Update: enaml.Update{
			MaxInFlight: 1,
			Serial:      true,
		},
	}
	return
}
