package cloudfoundry

import (
	"github.com/codegangsta/cli"
	"github.com/enaml-ops/enaml"
	natslib "github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/nats"
	"github.com/xchapter7x/lo"
)

//NewNatsPartition --
func NewNatsPartition(c *cli.Context) (igf InstanceGrouper) {
	igf = &NatsPartition{
		AZs:          c.StringSlice("az"),
		StemcellName: c.String("stemcell-name"),
		NetworkIPs:   c.StringSlice("nats-machine-ip"),
		NetworkName:  c.String("network"),
		VMTypeName:   c.String("nats-vm-type"),
		Metron:       NewMetron(c),
		Nats: natslib.NatsJob{
			Nats: &natslib.Nats{
				User:     c.String("nats-user"),
				Password: c.String("nats-pass"),
				Machines: c.StringSlice("nats-machine-ip"),
				Port:     4222,
			},
		},
		StatsdInjector: NewStatsdInjector(c),
	}
	return
}

//ToInstanceGroup --
func (s *NatsPartition) ToInstanceGroup() (ig *enaml.InstanceGroup) {
	ig = &enaml.InstanceGroup{
		Name:      "nats-partition",
		Instances: len(s.NetworkIPs),
		VMType:    s.VMTypeName,
		AZs:       s.AZs,
		Stemcell:  s.StemcellName,
		Jobs: []enaml.InstanceJob{
			s.newNatsJob(),
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

func (s *NatsPartition) newNatsJob() enaml.InstanceJob {
	return enaml.InstanceJob{
		Name:       "nats",
		Release:    "cf",
		Properties: s.Nats,
	}
}

//HasValidValues - Checks that fields in NatsPartition are valid
func (s *NatsPartition) HasValidValues() bool {
	lo.G.Debugf("checking '%s' for valid flags", "nats")

	if len(s.AZs) <= 0 {
		lo.G.Debugf("could not find the correct number of AZs configured '%v' : '%v'", len(s.AZs), s.AZs)
	}
	if len(s.NetworkIPs) <= 0 {
		lo.G.Debugf("could not find the correct number of network ips configured '%v' : '%v'", len(s.NetworkIPs), s.NetworkIPs)
	}
	if s.StemcellName == "" {
		lo.G.Debugf("could not find a valid stemcellname '%v'", s.StemcellName)
	}
	if s.VMTypeName == "" {
		lo.G.Debugf("could not find a valid vmtypename '%v'", s.VMTypeName)
	}
	if s.NetworkName == "" {
		lo.G.Debugf("could not find a valid NetworkName '%v'", s.NetworkName)
	}
	if s.Metron.Zone == "" {
		lo.G.Debugf("could not find a valid Metron.Zone '%v'", s.Metron.Zone)
	}
	if s.Metron.Secret == "" {
		lo.G.Debugf("could not find a valid Metron.Secret '%v'", s.Metron.Secret)
	}
	if s.Nats.Nats.User == "" {
		lo.G.Debugf("could not find a valid Nats.Nats.User '%v'", s.Nats.Nats.User)
	}
	if s.Nats.Nats.Password == "" {
		lo.G.Debugf("could not find a valid Nats.Nats.Password '%v'", s.Nats.Nats.Password)
	}
	if s.Nats.Nats.Machines == "" {
		lo.G.Debugf("could not find a valid Nats.Nats.Machines '%v'", s.Nats.Nats.Machines)
	}

	return (len(s.AZs) > 0 &&
		s.StemcellName != "" &&
		s.VMTypeName != "" &&
		s.Metron.Zone != "" &&
		s.Metron.Secret != "" &&
		s.NetworkName != "" &&
		len(s.NetworkIPs) > 0 &&
		s.Nats.Nats.User != "" &&
		s.Nats.Nats.Password != "" &&
		s.Nats.Nats.Machines != nil)
}
