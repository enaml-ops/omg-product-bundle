package cloudfoundry

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/enaml-ops/enaml"
	ccnglib "github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/cloud_controller_ng"
	"github.com/xchapter7x/lo"
)

func NewCloudControllerPartition(c *cli.Context) InstanceGrouper {
	return &CloudControllerPartition{
		AZs:                   c.StringSlice("az"),
		VMTypeName:            c.String("cc-vm-type"),
		StemcellName:          c.String("stemcell-name"),
		NetworkName:           c.String("network"),
		SystemDomain:          c.String("system-domain"),
		AppDomains:            c.StringSlice("app-domain"),
		AllowAppSSHAccess:     c.Bool("allow-app-ssh-access"),
		Metron:                NewMetron(c),
		ConsulAgent:           NewConsulAgent(c, []string{}),
		NFSMounter:            NewNFSMounter(c),
		StatsdInjector:        NewStatsdInjector(c),
		StagingUploadUser:     c.String("cc-staging-upload-user"),
		StagingUploadPassword: c.String("cc-staging-upload-password"),
		BulkAPIUser:           c.String("cc-bulk-api-user"),
		BulkAPIPassword:       c.String("cc-bulk-api-password"),
		InternalAPIUser:       c.String("cc-internal-api-user"),
		InternalAPIPassword:   c.String("cc-internal-api-password"),
		DbEncryptionKey:       c.String("cc-db-encryption-key"),
		HostKeyFingerprint:    c.String("host-key-fingerprint"),
		SupportAddress:        c.String("support-address"),
		MinCliVersion:         c.String("min-cli-version"),
	}
}

//ToInstanceGroup - Convert CLoud Controller Partition to an Instance Group
func (s *CloudControllerPartition) ToInstanceGroup() (ig *enaml.InstanceGroup) {
	ig = &enaml.InstanceGroup{
		Name:      "cloud_controller-partition",
		AZs:       s.AZs,
		Instances: 2, //Not sure where this number should be coming from!
		VMType:    s.VMTypeName,
		Stemcell:  s.StemcellName,
		Networks: []enaml.Network{
			enaml.Network{Name: s.NetworkName},
		},
		Jobs: []enaml.InstanceJob{
			newCloudControllerNgWorkerJob(s),
			s.ConsulAgent.CreateJob(),
			s.NFSMounter.CreateJob(),
			s.Metron.CreateJob(),
			s.StatsdInjector.CreateJob(),
		},
		Update: enaml.Update{
			MaxInFlight: 1,
			Serial:      true,
		},
	}
	return
}

func newCloudControllerNgWorkerJob(c *CloudControllerPartition) enaml.InstanceJob {
	return enaml.InstanceJob{
		Name:    "cloud_controller_worker",
		Release: "cf",
		Properties: &ccnglib.CloudControllerNgJob{
			AppSsh: &ccnglib.AppSsh{
				HostKeyFingerprint: c.HostKeyFingerprint,
			},
			Domain:                   c.SystemDomain,
			SystemDomain:             c.SystemDomain,
			AppDomains:               c.AppDomains,
			SystemDomainOrganization: "system",
			SupportAddress:           c.SupportAddress,
			Login: &ccnglib.Login{
				Url: fmt.Sprintf("https://login.%s", c.SystemDomain),
			},
			Cc: &ccnglib.Cc{
				AllowAppSshAccess:     c.AllowAppSSHAccess,
				DefaultToDiegoBackend: true,
				Buildpacks: &ccnglib.Buildpacks{
					BlobstoreType: "fog",
					FogConnection: &ccnglib.DefaultFogConnection{
						Provider:  "Local",
						LocalRoot: "/var/vcap/nfs/shared",
					},
				},
				Droplets: &ccnglib.Droplets{
					BlobstoreType: "fog",
					FogConnection: &ccnglib.DefaultFogConnection{
						Provider:  "Local",
						LocalRoot: "/var/vcap/nfs/shared",
					},
				},
				Packages: &ccnglib.Packages{
					BlobstoreType: "fog",
					FogConnection: &ccnglib.DefaultFogConnection{
						Provider:  "Local",
						LocalRoot: "/var/vcap/nfs/shared",
					},
				},
				ResourcePool: &ccnglib.ResourcePool{
					BlobstoreType: "fog",
					FogConnection: &ccnglib.DefaultFogConnection{
						Provider:  "Local",
						LocalRoot: "/var/vcap/nfs/shared",
					},
				},
				ExternalProtocol:          "https",
				LoggingLevel:              "debug",
				MaximumHealthCheckTimeout: 600,
				StagingUploadUser:         c.StagingUploadUser,
				StagingUploadPassword:     c.StagingUploadPassword,
				BulkApiUser:               c.BulkAPIUser,
				BulkApiPassword:           c.BulkAPIPassword,
				InternalApiUser:           c.InternalAPIUser,
				InternalApiPassword:       c.InternalAPIPassword,
				DbEncryptionKey:           c.DbEncryptionKey,
				DefaultRunningSecurityGroups: []string{
					"all_open",
				},
				DefaultStagingSecurityGroups: []string{
					"all_open",
				},
				DisableCustomBuildpacks: false,
				ExternalHost:            "api",
				InstallBuildpacks:       []string{},

				QuotaDefinitions: []string{},

				SecurityGroupDefinitions: []string{},

				Stacks: []string{},

				MinCliVersion:            c.MinCliVersion,
				MinRecommendedCliVersion: c.MinCliVersion,
			},
		},
	}
}

//HasValidValues - Check if valid values has been populated
func (s *CloudControllerPartition) HasValidValues() bool {

	lo.G.Debugf("checking '%s' for valid flags", "acceptanceTests")

	if len(s.AZs) <= 0 {
		lo.G.Debugf("could not find the correct number of AZs configured '%v' : '%v'", len(s.AZs), s.AZs)
	}

	if s.StemcellName == "" {
		lo.G.Debugf("could not find a valid stemcellname '%v'", s.StemcellName)
	}

	if s.NetworkName == "" {
		lo.G.Debugf("could not find a valid networkname '%v'", s.NetworkName)
	}

	if len(s.AppDomains) <= 0 {
		lo.G.Debugf("could not find the correct number of app domains configured '%v' : '%v'", len(s.AppDomains), s.AppDomains)
	}

	if s.SystemDomain == "" {
		lo.G.Debugf("could not find a valid system domain '%v'", s.SystemDomain)
	}

	if s.VMTypeName == "" {
		lo.G.Debugf("could not find a valid vmtypename '%v'", s.VMTypeName)
	}

	if s.Metron.Zone == "" {
		lo.G.Debugf("could not find a valid metron zone '%v'", s.Metron.Zone)
	}

	if s.Metron.Secret == "" {
		lo.G.Debugf("could not find a valid metron secret '%v'", s.Metron.Secret)
	}

	return (len(s.AZs) > 0 &&
		s.StemcellName != "" &&
		s.VMTypeName != "" &&
		s.Metron.Zone != "" &&
		s.Metron.Secret != "" &&
		s.NetworkName != "" &&
		s.SystemDomain != "" &&
		len(s.AppDomains) > 0 &&
		s.NFSMounter.hasValidValues() &&
		s.ConsulAgent.HasValidValues())
}
