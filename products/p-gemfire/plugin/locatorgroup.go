package gemfire_plugin

import (
	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-product-bundle/products/p-gemfire/enaml-gen/locator"
)

type LocatorGroup struct {
	StaticIPs   []string
	NetworkName string
}

func NewLocatorGroup(networkname string, staticips []string) *LocatorGroup {
	lg := new(LocatorGroup)
	lg.NetworkName = networkname
	lg.StaticIPs = staticips
	return lg
}

func (s *LocatorGroup) GetInstanceGroup() *enaml.InstanceGroup {
	instanceGroup := new(enaml.InstanceGroup)
	instanceGroup.Name = locatorGroup
	network := enaml.Network{
		Name:      s.NetworkName,
		StaticIPs: s.StaticIPs,
	}
	instanceGroup.AddNetwork(network)
	instanceGroup.Instances = len(s.StaticIPs)
	job := &enaml.InstanceJob{
		Name:    locatorJobName,
		Release: releaseName,
		Properties: locator.LocatorJob{
			Gemfire: &locator.Gemfire{
				Locator: &locator.Locator{
					Addresses: s.StaticIPs,
				},
			},
		},
	}
	instanceGroup.AddJob(job)
	return instanceGroup
}
