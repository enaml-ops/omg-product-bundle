package cloudfoundry

import (
	"github.com/codegangsta/cli"
	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/acceptance-tests"
	"github.com/xchapter7x/lo"
)

func NewAcceptanceTestsPartition(c *cli.Context, internet bool) InstanceGrouper {
	return &acceptanceTests{
		AZs:                      c.StringSlice("az"),
		StemcellName:             c.String("stemcell-name"),
		NetworkName:              c.String("network"),
		AppsDomain:               c.StringSlice("app-domain"),
		SystemDomain:             c.String("system-domain"),
		AdminPassword:            c.String("admin-password"),
		SkipCertVerify:           c.BoolT("skip-cert-verify"),
		IncludeInternetDependent: internet,
		VMType: c.String("acceptance-tests-vm-type"),
	}
}

func (a *acceptanceTests) ToInstanceGroup() *enaml.InstanceGroup {
	instanceGroupName := "acceptance-tests"
	if !a.IncludeInternetDependent {
		instanceGroupName += "-internetless"
	}
	return &enaml.InstanceGroup{
		Name:      instanceGroupName,
		Instances: 1,
		VMType:    a.VMType,
		Lifecycle: "errand",
		AZs:       a.AZs,
		Stemcell:  a.StemcellName,
		Networks: []enaml.Network{
			{Name: a.NetworkName},
		},
		Update: enaml.Update{
			MaxInFlight: 1,
		},
		Jobs: []enaml.InstanceJob{
			{
				Name:       "acceptance-tests",
				Release:    CFReleaseName,
				Properties: a.newAcceptanceTestsProperties(a.IncludeInternetDependent),
			},
		},
	}
}

func (a *acceptanceTests) newAcceptanceTestsProperties(internet bool) *acceptance_tests.AcceptanceTestsJob {
	var ad string
	if len(a.AppsDomain) > 0 {
		ad = a.AppsDomain[0]
	}
	return &acceptance_tests.AcceptanceTestsJob{
		AcceptanceTests: &acceptance_tests.AcceptanceTests{
			Api:                      prefixSystemDomain(a.SystemDomain, "api"),
			AppsDomain:               ad,
			AdminUser:                "admin",
			AdminPassword:            a.AdminPassword,
			IncludeLogging:           true,
			IncludeInternetDependent: internet,
			IncludeOperator:          true,
			IncludeServices:          true,
			IncludeSecurityGroups:    true,
			SkipSslValidation:        a.SkipCertVerify,
			SkipRegex:                "lucid64",
			JavaBuildpackName:        "java_buildpack_offline",
		},
	}
}

func (a *acceptanceTests) HasValidValues() bool {
	lo.G.Debugf("checking '%s' for valid flags", "acceptanceTests")

	if len(a.AZs) <= 0 {
		lo.G.Debugf("could not find the correct number of AZs configured '%v' : '%v'", len(a.AZs), a.AZs)
	}

	if a.StemcellName == "" {
		lo.G.Debugf("could not find a valid stemcellname '%v'", a.StemcellName)
	}

	if a.NetworkName == "" {
		lo.G.Debugf("could not find a valid networkname '%v'", a.NetworkName)
	}

	if len(a.AppsDomain) <= 0 {
		lo.G.Debugf("could not find the correct number of app domains configured '%v' : '%v'", len(a.AppsDomain), a.AppsDomain)
	}

	if a.SystemDomain == "" {
		lo.G.Debugf("could not find a valid system domain '%v'", a.SystemDomain)
	}

	if a.AdminPassword == "" {
		lo.G.Debugf("could not find a valid admin password '%v'", a.AdminPassword)
	}

	return len(a.AZs) > 0 &&
		a.StemcellName != "" &&
		a.NetworkName != "" &&
		len(a.AppsDomain) > 0 &&
		a.SystemDomain != "" &&
		a.AdminPassword != ""
}
