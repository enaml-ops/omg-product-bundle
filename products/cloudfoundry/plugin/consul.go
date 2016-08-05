package cloudfoundry

import (
	"github.com/codegangsta/cli"
	"github.com/enaml-ops/enaml"
	"github.com/xchapter7x/lo"
)

//NewConsulPartition -
func NewConsulPartition(c *cli.Context) InstanceGrouper {
	return &Consul{
		AZs:            c.StringSlice("az"),
		StemcellName:   c.String("stemcell-name"),
		NetworkIPs:     c.StringSlice("consul-ip"),
		NetworkName:    c.String("network"),
		VMTypeName:     c.String("consul-vm-type"),
		ConsulAgent:    NewConsulAgentServer(c),
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
		AZs:       s.AZs,
		Stemcell:  s.StemcellName,
		Jobs: []enaml.InstanceJob{
			s.ConsulAgent.CreateJob(),
			s.Metron.CreateJob(),
			s.StatsdInjector.CreateJob(),
		},
		Networks: []enaml.Network{
			enaml.Network{Name: s.NetworkName, StaticIPs: s.NetworkIPs},
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

	if len(s.AZs) <= 0 {
		lo.G.Debugf("could not find the correct number of AZs configured '%v' : '%v'", len(s.AZs), s.AZs)
	}

	if s.StemcellName == "" {
		lo.G.Debugf("could not find a valid stemcellname '%v'", s.StemcellName)
	}

	if s.NetworkName == "" {
		lo.G.Debugf("could not find a valid networkname '%v'", s.NetworkName)
	}

	if s.VMTypeName == "" {
		lo.G.Debugf("could not find a valid vmtypename '%v'", s.VMTypeName)
	}

	if len(s.NetworkIPs) <= 0 {
		lo.G.Debugf("could not find the correct number of networkips configured '%v' : '%v'", len(s.NetworkIPs), s.NetworkIPs)
	}

	if len(s.ConsulAgent.EncryptKeys) <= 0 {
		lo.G.Debugf("could not find the correct number of encrypt keys configured '%v' : '%v'", len(s.ConsulAgent.EncryptKeys), s.ConsulAgent.EncryptKeys)
	}
	return (len(s.AZs) > 0 &&
		s.StemcellName != "" &&
		s.VMTypeName != "" &&
		s.NetworkName != "" &&
		len(s.NetworkIPs) > 0 &&
		len(s.ConsulAgent.EncryptKeys) > 0 &&
		s.ConsulAgent.HasValidValues())
}
