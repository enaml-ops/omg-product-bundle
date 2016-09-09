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

// GetFlags returns the CLI flags accepted by the plugin.
func (p *Plugin) GetFlags() []pcli.Flag {
	return []pcli.Flag{
		pcli.CreateStringFlag("deployment-name", "the name bosh will use for the deployment", "p-rabbitmq"),
		pcli.CreateStringFlag("network", "the name of the network to use"),
		// pcli.CreateStringFlag("stemcell-url", "the url of the stemcell you wish to use"),
		// pcli.CreateStringFlag("stemcell-ver", "the version number of the stemcell you wish to use"),
		// pcli.CreateStringFlag("stemcell-sha", "the sha of the stemcell you will use"),
		// pcli.CreateStringFlag("stemcell-name", "the name of the stemcell you will use", "ubuntu-trusty"),
		pcli.CreateStringSliceFlag("server-ip", "rabbit-mq server IPs to use"),
		pcli.CreateStringFlag("syslog-address", "the address of your syslog drain"),
		pcli.CreateIntFlag("syslog-port", "the port for your syslog connection", "514"),
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
