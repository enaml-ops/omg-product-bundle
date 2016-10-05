package gemfire_plugin

import (
	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-product-bundle/products/p-gemfire/enaml-gen/server"
)

type ServerGroup struct {
	StaticIPs   []string
	NetworkName string
}

func NewServerGroup(networkname string, staticips []string) *ServerGroup {
	sg := new(ServerGroup)
	sg.NetworkName = networkname
	sg.StaticIPs = staticips
	return sg
}

func (s *ServerGroup) GetInstanceGroup() *enaml.InstanceGroup {
	instanceGroup := new(enaml.InstanceGroup)
	instanceGroup.Name = serverGroup
	network := enaml.Network{
		Name:      s.NetworkName,
		StaticIPs: s.StaticIPs,
	}
	instanceGroup.AddNetwork(network)
	instanceGroup.Instances = len(s.StaticIPs)
	job := &enaml.InstanceJob{
		Name:    serverJobName,
		Release: releaseName,
		Properties: server.ServerJob{
			Gemfire: &server.Gemfire{
				Locator: &server.Locator{
					Addresses: s.StaticIPs,
				},
			},
		},
	}
	instanceGroup.AddJob(job)
	return instanceGroup
}
