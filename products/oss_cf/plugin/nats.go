package cloudfoundry

import (
	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-product-bundle/products/oss_cf/enaml-gen/nats"
	"github.com/enaml-ops/omg-product-bundle/products/oss_cf/plugin/config"
)

//NatsPartition -
type NatsPartition struct {
	Config         *config.Config
	Metron         *Metron
	StatsdInjector *StatsdInjector
}

//NewNatsPartition --
func NewNatsPartition(config *config.Config) (igf InstanceGroupCreator) {
	igf = &NatsPartition{
		Config:         config,
		Metron:         NewMetron(config),
		StatsdInjector: NewStatsdInjector(nil),
	}
	return
}

//ToInstanceGroup --
func (s *NatsPartition) ToInstanceGroup() (ig *enaml.InstanceGroup) {
	ig = &enaml.InstanceGroup{
		Name:      "nats-partition",
		Instances: len(s.Config.NATSMachines),
		VMType:    s.Config.NatsVMType,
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
		Properties: &nats.NatsJob{
			Nats: &nats.Nats{
				User:     s.Config.NATSUser,
				Password: s.Config.NATSPassword,
				Machines: s.Config.NATSMachines,
				Port:     s.Config.NATSPort,
			},
		},
	}
}
