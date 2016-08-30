package cloudfoundry

import (
	"github.com/codegangsta/cli"
	"github.com/enaml-ops/enaml"
	"github.com/xchapter7x/lo"
)

// Consul -
type Consul struct {
	Config         *Config
	VMTypeName     string
	NetworkIPs     []string
	ConsulAgent    *ConsulAgent
	Metron         *Metron
	StatsdInjector *StatsdInjector
}

//NewConsulPartition -
func NewConsulPartition(c *cli.Context, config *Config) InstanceGrouper {
	return &Consul{
		Config:         config,
		NetworkIPs:     c.StringSlice("consul-ip"),
		VMTypeName:     c.String("consul-vm-type"),
		ConsulAgent:    NewConsulAgentServer(c, config),
		Metron:         NewMetron(c),
		StatsdInjector: NewStatsdInjector(c),
	}
}

//ToInstanceGroup -
func (s *Consul) ToInstanceGroup() (ig *enaml.InstanceGroup) {
	ig = &enaml.InstanceGroup{
		Name:      "consul-partition",
		Instances: len(s.NetworkIPs),
		VMType:    s.VMTypeName,
		AZs:       s.Config.AZs,
		Stemcell:  s.Config.StemcellName,
		Jobs: []enaml.InstanceJob{
			s.ConsulAgent.CreateJob(),
			s.Metron.CreateJob(),
			s.StatsdInjector.CreateJob(),
		},
		Networks: []enaml.Network{
			enaml.Network{Name: s.Config.NetworkName, StaticIPs: s.NetworkIPs},
		},
		Update: enaml.Update{
			MaxInFlight: 1,
			Serial:      true,
		},
	}
	return
}

//HasValidValues - Check if the datastructure has valid fields
func (s *Consul) HasValidValues() bool {

	lo.G.Debugf("checking '%s' for valid flags", "consul")

	if s.VMTypeName == "" {
		lo.G.Debugf("could not find a valid vmtypename '%v'", s.VMTypeName)
	}

	if len(s.NetworkIPs) <= 0 {
		lo.G.Debugf("could not find the correct number of networkips configured '%v' : '%v'", len(s.NetworkIPs), s.NetworkIPs)
	}

	return (s.VMTypeName != "" &&
		len(s.NetworkIPs) > 0 &&
		s.ConsulAgent.HasValidValues())
}
