package prabbitmq

import (
	"fmt"

	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/pluginlib/pcli"
	"github.com/enaml-ops/pluginlib/product"
	"github.com/enaml-ops/pluginlib/util"
	"github.com/xchapter7x/lo"
)

// Plugin is an omg product plugin for deploying p-rabbitmq.
type Plugin struct{}

func (p *Plugin) GetFlags() []pcli.Flag {
	return []pcli.Flag{
		pcli.Flag{FlagType: pcli.StringFlag, Name: "deployment-name", Value: "p-rabbitmq", Usage: "the name bosh will use for the deployment"},
		pcli.Flag{FlagType: pcli.StringFlag, Name: "network", Usage: "the name of the network to use"},
		// pcli.Flag{FlagType: pcli.StringFlag, Name: "stemcell-url", Usage: "the url of the stemcell you wish to use"},
		// pcli.Flag{FlagType: pcli.StringFlag, Name: "stemcell-ver", Usage: "the version number of the stemcell you wish to use"},
		// pcli.Flag{FlagType: pcli.StringFlag, Name: "stemcell-sha", Usage: "the sha of the stemcell you will use"},
		// pcli.Flag{FlagType: pcli.StringFlag, Name: "stemcell-name", Value: "ubuntu-trusty", Usage: "the name of the stemcell you will use"},
		pcli.Flag{FlagType: pcli.StringSliceFlag, Name: "server-ip", Usage: "rabbit-mq server IPs to use"},
		pcli.Flag{FlagType: pcli.StringFlag, Name: "syslog-address", Usage: "the address of your syslog drain"},
		pcli.Flag{FlagType: pcli.IntFlag, Name: "syslog-port", Value: "514", Usage: "the port for your syslog connection"},
	}
}

// GetMeta returns metadata about the p-rabbitmq product.
func (p *Plugin) GetMeta() product.Meta {
	return product.Meta{
		Name: "p-rabbitmq",
		Properties: map[string]interface{}{
			"version":                  "",                                                          // TODO GET FROM PLUGIN MAIN FILE?
			"pivotal-rabbit-mq":        fmt.Sprintf("%s / %s", "pivotal-rabbit-mq", ProductVersion), // TODO match pivnet on name
			"cf-rabbitmq-release":      fmt.Sprintf("%s / %s", CFRabbitMQReleaseName, CFRabbitMQReleaseVersion),
			"service-metrics-release":  fmt.Sprintf("%s / %s", ServiceMetricsReleaseName, ServiceMetricsReleaseVersion),
			"loggregator-release":      fmt.Sprintf("%s / %s", LoggregatorReleaseName, LoggregatorReleaseVersion),
			"rabbitmq-metrics-release": fmt.Sprintf("%s / %s", RabbitMQMetricsReleaseName, RabbitMQMetricsReleaseVersion),
		},
	}
}

// GetProduct generates a BOSH deployment manifest for p-rabbitmq.
func (p *Plugin) GetProduct(args []string, cloudConfig []byte) []byte {
	c := pluginutil.NewContext(args, pluginutil.ToCliFlagArray(p.GetFlags()))
	cfg, err := configFromContext(c)
	if err != nil {
		lo.G.Error(err.Error())
	}

	dm := new(enaml.DeploymentManifest)
	dm.SetName(cfg.DeploymentName)

	dm.AddRelease(enaml.Release{Name: CFRabbitMQReleaseName, Version: CFRabbitMQReleaseVersion})
	dm.AddRelease(enaml.Release{Name: ServiceMetricsReleaseName, Version: ServiceMetricsReleaseVersion})
	dm.AddRelease(enaml.Release{Name: LoggregatorReleaseName, Version: LoggregatorReleaseVersion})
	dm.AddRelease(enaml.Release{Name: RabbitMQMetricsReleaseName, Version: RabbitMQMetricsReleaseVersion})

	// TODO add stemcell
	// ubuntu-trusty, 3232.17
	//dm.AddRemoteStemcell(os, alias, ver, url, sha1)

	// add instance groups
	dm.AddInstanceGroup(p.NewRabbitMQServerPartition(cfg))

	dm.Update = enaml.Update{
		Canaries:        1,
		CanaryWatchTime: "30000-300000",
		UpdateWatchTime: "30000-300000",
		MaxInFlight:     1,
		Serial:          true,
	}

	return dm.Bytes()
}
