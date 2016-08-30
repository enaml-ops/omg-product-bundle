package cloudfoundry

import (
	"github.com/codegangsta/cli"
	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/nats"
	"github.com/xchapter7x/lo"
)

//NatsPartition -
type NatsPartition struct {
	Config         *Config
	VMTypeName     string
	Metron         *Metron
	StatsdInjector *StatsdInjector
}

//NewNatsPartition --
func NewNatsPartition(c *cli.Context, config *Config) (igf InstanceGrouper) {
	igf = &NatsPartition{
		Config:         config,
		VMTypeName:     c.String("nats-vm-type"),
		Metron:         NewMetron(c),
		StatsdInjector: NewStatsdInjector(c),
	}
	return
}

//ToInstanceGroup --
func (s *NatsPartition) ToInstanceGroup() (ig *enaml.InstanceGroup) {
	ig = &enaml.InstanceGroup{
		Name:      "nats-partition",
		Instances: len(s.Config.NATSMachines),
		VMType:    s.VMTypeName,
		AZs:       s.Config.AZs,
		Stemcell:  s.Config.StemcellName,
		Jobs: []enaml.InstanceJob{
			s.newNatsJob(),
			s.Metron.CreateJob(),
			s.StatsdInjector.CreateJob(),
		},
		Networks: []enaml.Network{
			enaml.Network{Name: s.Config.NetworkName, StaticIPs: s.Config.NATSMachines},
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
		Name:    "nats",
		Release: "cf",
		Properties: &nats.Nats{
			User:     s.Config.NATSUser,
			Password: s.Config.NATSPassword,
			Machines: s.Config.NATSMachines,
			Port:     s.Config.NATSPort,
		},
	}
}

//HasValidValues - Checks that fields in NatsPartition are valid
func (s *NatsPartition) HasValidValues() bool {
	lo.G.Debugf("checking '%s' for valid flags", "nats")

	if s.VMTypeName == "" {
		lo.G.Debugf("could not find a valid vmtypename '%v'", s.VMTypeName)
	}
	if s.Metron.Zone == "" {
		lo.G.Debugf("could not find a valid Metron.Zone '%v'", s.Metron.Zone)
	}
	if s.Metron.Secret == "" {
		lo.G.Debugf("could not find a valid Metron.Secret '%v'", s.Metron.Secret)
	}

	return (s.VMTypeName != "" &&
		s.Metron.Zone != "" &&
		s.Metron.Secret != "")
}
