package cloudfoundry

import (
	"fmt"

	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/smoke-tests"
	"github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/plugin/config"
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
			SmokeTests: &smoke_tests.SmokeTests{
				UseExistingOrg:   false,
				UseExistingSpace: false,
				Space:            "CF_SMOKE_TEST_SPACE",
				Org:              "CF_SMOKE_TEST_ORG",
				Password:         s.Config.SmokeTestsPassword,
				User:             "smoke_tests",
				Api:              fmt.Sprintf("%s://api.%s", s.Config.UAALoginProtocol, s.Config.SystemDomain),
				AppsDomain:       s.Config.AppDomains[0],
			},
		},
	}
}
