package cloudfoundry

import (
	"fmt"

	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-product-bundle/products/oss_cf/enaml-gen/smoke-tests"
	"github.com/enaml-ops/omg-product-bundle/products/oss_cf/plugin/config"
)

type smokeErrand struct {
	Config *config.Config
}

//NewSmokeErrand - errand definition for smoke tests
func NewSmokeErrand(config *config.Config) InstanceGroupCreator {
	return &smokeErrand{
		Config: config,
	}
}

//ToInstanceGroup -
func (s *smokeErrand) ToInstanceGroup() (ig *enaml.InstanceGroup) {
	ig = &enaml.InstanceGroup{
		Name:      "smoke-tests",
		Instances: 1,
		VMType:    s.Config.ErrandVMType,
		AZs:       s.Config.AZs,
		Stemcell:  s.Config.StemcellName,
		Jobs: []enaml.InstanceJob{
			s.createSmokeJob(),
		},
		Lifecycle: "errand",
		Networks: []enaml.Network{
			enaml.Network{Name: s.Config.NetworkName},
		},
		Update: enaml.Update{
			MaxInFlight: 1,
		},
	}
	return
}

func (s *smokeErrand) createSmokeJob() enaml.InstanceJob {
	return enaml.InstanceJob{
		Name:    "smoke-tests",
		Release: "cf",
		Properties: &smoke_tests.SmokeTestsJob{
			Cf: &smoke_tests.Cf{
				AdminUsername:     "smoke_tests",
				AdminPassword:     s.Config.SmokeTestsPassword,
				ApiUrl:            fmt.Sprintf("%s://api.%s", s.Config.UAALoginProtocol, s.Config.SystemDomain),
				AppDomains:        s.Config.AppDomains[0],
				SkipSslValidation: s.Config.SkipSSLCertVerify,
				SmokeTests: &smoke_tests.CfSmokeTests{
					UseExistingOrg: false,
					Org:            "CF_SMOKE_TEST_ORG",
				},
			},
			CfMysql: &smoke_tests.CfMysql{
				SmokeTests: &smoke_tests.CfMysqlSmokeTests{
					Password: s.Config.SmokeTestsPassword,
				},
				Mysql: &smoke_tests.Mysql{
					AdminPassword: s.Config.SmokeTestsPassword,
					AdminUsername: "smoke_tests",
				},
				ExternalHost: s.Config.AppDomains[0],
			},
		},
	}
}
