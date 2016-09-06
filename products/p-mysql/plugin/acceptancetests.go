package pmysql

import (
	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-product-bundle/products/p-mysql/enaml-gen/acceptance-tests"
)

func NewAcceptanceTests(plgn *Plugin) *enaml.InstanceGroup {
	return &enaml.InstanceGroup{
		Name:               "acceptance-tests",
		Lifecycle:          "errand",
		Instances:          1,
		VMType:             plgn.VMTypeName,
		AZs:                plgn.AZs,
		Stemcell:           plgn.StemcellName,
		PersistentDiskType: plgn.DiskTypeName,
		Jobs: []enaml.InstanceJob{
			newAcceptanceTestsJob(plgn),
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

func newAcceptanceTestsJob(plgn *Plugin) enaml.InstanceJob {
	return enaml.InstanceJob{
		Name:    "acceptance-tests",
		Release: "cf-mysql",
		Properties: &acceptance_tests.AcceptanceTestsJob{
			TimeoutScale: 1,
			Cf:           nil,
			Proxy:        nil,
			Broker:       nil,
			Service:      nil,
		},
	}
}
