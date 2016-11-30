package gemfire_plugin

import (
	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-product-bundle/products/p-gemfire/enaml-gen/locator"
)

type LocatorGroup struct {
	StaticIPs   []string
	NetworkName string
	Port        int
	RestPort    int
	VMMemory    int
	VMType      string
}

func NewLocatorGroup(networkname string, staticips []string, port, restport, vmmemory int, vmtype string) *LocatorGroup {
	lg := new(LocatorGroup)
	lg.NetworkName = networkname
	lg.StaticIPs = staticips
	lg.Port = port
	lg.RestPort = restport
	lg.VMMemory = vmmemory
	lg.VMType = vmtype
	return lg
}

func (s *LocatorGroup) GetInstanceGroup() *enaml.InstanceGroup {
	instanceGroup := new(enaml.InstanceGroup)
	instanceGroup.Name = locatorGroup
	network := enaml.Network{
		Name: s.NetworkName,
		Default: []interface{}{
			"dns",
			"gateway",
		},
		StaticIPs: s.StaticIPs,
	}
	instanceGroup.AddNetwork(network)
	instanceGroup.Instances = len(s.StaticIPs)
	locatorJob := &enaml.InstanceJob{
		Name:    locatorJobName,
		Release: releaseName,
		Properties: locator.LocatorJob{
			Gemfire: &locator.Gemfire{
				Locator: &locator.Locator{
					Addresses: s.StaticIPs,
					Port:      s.Port,
					RestPort:  s.RestPort,
					VmMemory:  s.VMMemory,
				},
				ClusterTopology: &locator.ClusterTopology{
					NumberOfLocators:    len(s.StaticIPs),
					MinNumberOfLocators: len(s.StaticIPs),
				},
			},
		},
	}
	arpJob := &enaml.InstanceJob{
		Name:       arpCleanerJobName,
		Release:    releaseName,
		Properties: locator.LocatorJob{},
	}
	instanceGroup.AddJob(locatorJob)
	instanceGroup.AddJob(arpJob)
	instanceGroup.VMType = s.VMType
	instanceGroup.Lifecycle = "service"
	return instanceGroup
}
