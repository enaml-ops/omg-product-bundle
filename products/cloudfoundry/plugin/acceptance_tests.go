package cloudfoundry

import (
	"github.com/codegangsta/cli"
	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/acceptance-tests"
	"github.com/xchapter7x/lo"
)

type acceptanceTests struct {
	Config                   *Config
	AdminPassword            string
	VMType                   string
	IncludeInternetDependent bool
}

func NewAcceptanceTestsPartition(c *cli.Context, internet bool, config *Config) InstanceGrouper {
	return &acceptanceTests{
		Config:                   config,
		AdminPassword:            c.String("admin-password"),
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
		AZs:       a.Config.AZs,
		Stemcell:  a.Config.StemcellName,
		Networks: []enaml.Network{
			{Name: a.Config.NetworkName},
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
	if len(a.Config.AppDomains) > 0 {
		ad = a.Config.AppDomains[0]
	}
	return &acceptance_tests.AcceptanceTestsJob{
		AcceptanceTests: &acceptance_tests.AcceptanceTests{
			Api:                      prefixSystemDomain(a.Config.SystemDomain, "api"),
			AppsDomain:               ad,
			AdminUser:                "admin",
			AdminPassword:            a.AdminPassword,
			IncludeLogging:           true,
			IncludeInternetDependent: internet,
			IncludeOperator:          true,
			IncludeServices:          true,
			IncludeSecurityGroups:    true,
			SkipSslValidation:        a.Config.SkipSSLCertVerify,
			SkipRegex:                "lucid64",
			JavaBuildpackName:        "java_buildpack_offline",
		},
	}
}

func (a *acceptanceTests) HasValidValues() bool {
	lo.G.Debugf("checking '%s' for valid flags", "acceptanceTests")

	if a.AdminPassword == "" {
		lo.G.Debugf("could not find a valid admin password '%v'", a.AdminPassword)
	}

	return a.AdminPassword != ""
}
