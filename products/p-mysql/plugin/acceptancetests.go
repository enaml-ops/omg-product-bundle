package pmysql

import (
	"fmt"

	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-product-bundle/products/p-mysql/enaml-gen/acceptance-tests"
)

func NewAcceptanceTests(plgn *Plugin) *enaml.InstanceGroup {
	return &enaml.InstanceGroup{
		Name:      "acceptance-tests",
		Lifecycle: "errand",
		Instances: 1,
		VMType:    plgn.VMTypeName,
		AZs:       plgn.AZs,
		Stemcell:  StemcellAlias,
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
		Release: CFMysqlReleaseName,
		Properties: &acceptance_tests.AcceptanceTestsJob{
			TimeoutScale: 1,
			Cf: &acceptance_tests.Cf{
				ApiUrl:            fmt.Sprintf("https://api.sys.%s", plgn.BaseDomain),
				AdminUsername:     "admin",
				AdminPassword:     plgn.CFAdminPassword,
				AppsDomain:        fmt.Sprintf("apps.%s", plgn.BaseDomain),
				SkipSslValidation: true,
			},
			Proxy: &acceptance_tests.Proxy{
				ExternalHost: fmt.Sprintf("p-mysql.sys.%s", plgn.BaseDomain),
				ApiUsername:  plgn.ProxyAPIUser,
				ApiPassword:  plgn.ProxyAPIPass,
				ProxyCount:   len(plgn.ProxyIPs),
			},
			Broker: &acceptance_tests.Broker{
				Host: fmt.Sprintf("p-mysql.sys.%s", plgn.BaseDomain),
			},
			Service: &acceptance_tests.Service{
				Name: "p-mysql",
				Plans: []map[string]interface{}{
					{
						"name":           "100mb-dev",
						"max_storage_mb": 100,
					},
				},
			},
		},
	}
}
