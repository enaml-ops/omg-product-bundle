package cloudfoundry

import (
	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/acceptance-tests"
	"github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/plugin/config"
)

type acceptanceTests struct {
	Config                   *config.Config
	IncludeInternetDependent bool
}

func NewAcceptanceTestsPartition(internet bool, config *config.Config) InstanceGroupCreator {
	return &acceptanceTests{
		Config: config,
		IncludeInternetDependent: internet,
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
		VMType:    a.Config.ErrandVMType,
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
			AdminPassword:            a.Config.AdminPassword,
			IncludeLogging:           true,
			IncludeInternetDependent: internet,
			IncludeOperator:          true,
			IncludeServices:          true,
			IncludeSecurityGroups:    true,
			SkipSslValidation:        a.Config.SkipSSLCertVerify,
			SkipRegex:                "lucid64",
			JavaBuildpackName:        javaBuildpackName,
		},
	}
}
