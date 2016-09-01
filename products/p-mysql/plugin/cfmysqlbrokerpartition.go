package pmysql

import "github.com/enaml-ops/enaml"

func NewCfMysqlBrokerPartition(plgn *Plugin) *enaml.InstanceGroup {
	return &enaml.InstanceGroup{
		Name:               "cf-mysql-broker-partition",
		Lifecycle:          "service",
		Instances:          len(plgn.BrokerIPs),
		VMType:             plgn.VMTypeName,
		AZs:                plgn.AZs,
		Stemcell:           plgn.StemcellName,
		PersistentDiskType: plgn.DiskTypeName,
		Jobs: []enaml.InstanceJob{
			newBrokerJob(plgn),
		},
		Networks: []enaml.Network{
			enaml.Network{
				Name:      plgn.NetworkName,
				StaticIPs: plgn.BrokerIPs,
				Default:   []interface{}{"dns", "gateway"},
			},
		},
		Update: enaml.Update{
			MaxInFlight: 1,
		},
	}
}

func newBrokerJob(plgn *Plugin) enaml.InstanceJob {
	return enaml.InstanceJob{
		Name:       "cf-mysql-broker",
		Release:    "cf-mysql",
		Properties: nil,
	}
}
