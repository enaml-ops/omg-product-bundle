package cloudfoundry

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/smoke-tests"
	"github.com/xchapter7x/lo"
)

type SmokeErrand struct {
	Config     *Config
	VMTypeName string
	Instances  int
	Protocol   string
	Password   string
}

//NewSmokeErrand - errand definition for smoke tests
func NewSmokeErrand(c *cli.Context, config *Config) InstanceGrouper {
	return &SmokeErrand{
		Config:     config,
		VMTypeName: c.String("errand-vm-type"),
		Protocol:   c.String("uaa-login-protocol"),
		Password:   c.String("smoke-tests-password"),
	}
}

//ToInstanceGroup -
func (s *SmokeErrand) ToInstanceGroup() (ig *enaml.InstanceGroup) {
	ig = &enaml.InstanceGroup{
		Name:      "smoke-tests",
		Instances: 1,
		VMType:    s.VMTypeName,
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

func (s *SmokeErrand) createSmokeJob() enaml.InstanceJob {
	return enaml.InstanceJob{
		Name:    "smoke-tests",
		Release: "cf",
		Properties: &smoke_tests.SmokeTestsJob{
			SmokeTests: &smoke_tests.SmokeTests{
				UseExistingOrg:   false,
				UseExistingSpace: false,
				Space:            "CF_SMOKE_TEST_SPACE",
				Org:              "CF_SMOKE_TEST_ORG",
				Password:         s.Password,
				User:             "smoke_tests",
				Api:              fmt.Sprintf("%s://api.%s", s.Protocol, s.Config.SystemDomain),
				AppsDomain:       s.Config.AppDomains[0],
			},
		},
	}
}

//HasValidValues - Check if the datastructure has valid fields
func (s *SmokeErrand) HasValidValues() bool {
	lo.G.Debugf("checking '%s' for valid flags", "smoke")

	if s.VMTypeName == "" {
		lo.G.Debugf("could not find a valid vmtypename '%v'", s.VMTypeName)
	}
	if s.Protocol == "" {
		lo.G.Debugf("could not find a valid Protocol '%v'", s.Protocol)
	}
	if s.Password == "" {
		lo.G.Debugf("could not find a valid Password '%v'", s.Password)
	}
	return (s.VMTypeName != "" &&
		s.Protocol != "" &&
		s.Password != "")
}
