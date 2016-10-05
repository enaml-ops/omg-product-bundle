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
}

func NewServerGroup(networkname string, serverport int, instanceCount int, vmtype string, locator *LocatorGroup) *ServerGroup {
	sg := new(ServerGroup)
	sg.NetworkName = networkname
	sg.Locator = locator
	sg.InstanceCount = instanceCount
	sg.VMType = vmtype
	return sg
}

func (s *ServerGroup) GetInstanceGroup() *enaml.InstanceGroup {
	instanceGroup := new(enaml.InstanceGroup)
	instanceGroup.Name = serverGroup
	network := enaml.Network{
		Name: s.NetworkName,
		Default: []interface{}{
			"dns",
			"gateway",
		},
	}
	instanceGroup.AddNetwork(network)
	instanceGroup.Instances = s.InstanceCount
	instanceGroup.VMType = s.VMType
	job := &enaml.InstanceJob{
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
				},
				ClusterTopology: &server.ClusterTopology{
					NumberOfLocators: len(s.Locator.StaticIPs),
					NumberOfServers:  s.InstanceCount,
				},
			},
		},
	}
	instanceGroup.AddJob(job)
	return instanceGroup
}
