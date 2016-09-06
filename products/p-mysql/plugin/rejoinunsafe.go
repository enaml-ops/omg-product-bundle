package pmysql

import (
	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-product-bundle/products/p-mysql/enaml-gen/rejoin-unsafe"
)

func NewRejoinUnsafe(plgn *Plugin) *enaml.InstanceGroup {
	return &enaml.InstanceGroup{
		Name:               "rejoin-unsafe",
		Lifecycle:          "errand",
		Instances:          1,
		VMType:             plgn.VMTypeName,
		AZs:                plgn.AZs,
		Stemcell:           plgn.StemcellName,
		PersistentDiskType: plgn.DiskTypeName,
		Jobs: []enaml.InstanceJob{
			newRejoinUnsafeJob(plgn),
		},
		Networks: []enaml.Network{
			enaml.Network{
				Name:    plgn.NetworkName,
				Default: []interface{}{"dns", "gateway"},
			},
		},
		Update: enaml.Update{
			MaxInFlight: 1,
		},
	}
}

func newRejoinUnsafeJob(plgn *Plugin) enaml.InstanceJob {
	return enaml.InstanceJob{
		Name:    "rejoin-unsafe",
		Release: "cf-mysql",
		Properties: &rejoin_unsafe.RejoinUnsafeJob{
			ClusterIps: plgn.IPs,
			CfMysql: &rejoin_unsafe.CfMysql{
				Mysql: &rejoin_unsafe.Mysql{
					GaleraHealthcheck: &rejoin_unsafe.GaleraHealthcheck{
						EndpointUsername: plgn.GaleraHealthcheckUsername,
						EndpointPassword: plgn.GaleraHealthcheckPassword,
					},
				},
			},
		},
	}
}
