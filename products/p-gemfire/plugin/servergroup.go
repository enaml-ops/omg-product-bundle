package gemfire_plugin

import (
	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-product-bundle/products/p-gemfire/enaml-gen/server"
)

type ServerGroup struct {
	Locator       *LocatorGroup
	InstanceCount int
	NetworkName   string
	VMType        string
	Port          int
	VMMemory      int
	StaticIPs     []string
	DevRestPort   int
	DevRestActive bool
}

func NewServerGroup(networkname string, serverport int, instanceCount int, staticIPs []string, vmtype string, vmmemory int, devrestport int, devrestactive bool, locator *LocatorGroup) *ServerGroup {
	sg := new(ServerGroup)
	sg.DevRestPort = devrestport
	sg.DevRestActive = devrestactive
	sg.NetworkName = networkname
	sg.Locator = locator
	sg.StaticIPs = staticIPs
	sg.InstanceCount = instanceCount
	sg.VMType = vmtype
	sg.Port = serverport
	sg.VMMemory = vmmemory
	return sg
}

func (s *ServerGroup) getInstanceCount() int {
	if len(s.StaticIPs) > 0 {
		return len(s.StaticIPs)
	}
	return s.InstanceCount
}

func (s *ServerGroup) getNetwork() enaml.Network {
	network := enaml.Network{
		Name: s.NetworkName,
		Default: []interface{}{
			"dns",
			"gateway",
		},
	}

	if len(s.StaticIPs) > 0 {
		network.StaticIPs = s.StaticIPs
	}
	return network
}

func (s *ServerGroup) GetInstanceGroup() *enaml.InstanceGroup {
	instanceGroup := new(enaml.InstanceGroup)
	instanceGroup.Name = serverGroup
	instanceGroup.AddNetwork(s.getNetwork())
	instanceGroup.Instances = s.getInstanceCount()
	instanceGroup.VMType = s.VMType
	serverJob := &enaml.InstanceJob{
		Name:    serverJobName,
		Release: releaseName,
		Properties: server.ServerJob{
			Gemfire: &server.Gemfire{
				Locator: &server.Locator{
					Addresses: s.Locator.StaticIPs,
					Port:      s.Locator.Port,
				},
				Server: &server.Server{
					RestPort: s.Locator.RestPort,
					Port:     s.Port,
					VmMemory: s.VMMemory,
					DevRestApi: &server.DevRestApi{
						Port:   s.DevRestPort,
						Active: s.DevRestActive,
					},
				},
				ClusterTopology: &server.ClusterTopology{
					NumberOfLocators: len(s.Locator.StaticIPs),
					NumberOfServers:  s.getInstanceCount(),
				},
			},
		},
	}
	arpJob := &enaml.InstanceJob{
		Name:       arpCleanerJobName,
		Release:    releaseName,
		Properties: server.ServerJob{},
	}
	instanceGroup.AddJob(serverJob)
	instanceGroup.AddJob(arpJob)
	return instanceGroup
}
