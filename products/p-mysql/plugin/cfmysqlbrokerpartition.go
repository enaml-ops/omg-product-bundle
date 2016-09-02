package pmysql

import (
	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-product-bundle/products/p-mysql/enaml-gen/cf-mysql-broker"
)

func NewCfMysqlBrokerPartition(plgn *Plugin) *enaml.InstanceGroup {
	return &enaml.InstanceGroup{
		Name:               "cf-mysql-broker-partition",
		Lifecycle:          "service",
		Instances:          len(plgn.BrokerIPs),
		VMType:             plgn.VMTypeName,
		AZs:                plgn.AZs,
		Stemcell:           plgn.StemcellName,
		PersistentDiskType: plgn.DiskTypeName,
		Jobs: []enaml.InstanceJob{
			newBrokerJob(plgn),
		},
		Networks: []enaml.Network{
			enaml.Network{
				Name:      plgn.NetworkName,
				StaticIPs: plgn.BrokerIPs,
				Default:   []interface{}{"dns", "gateway"},
			},
		},
		Update: enaml.Update{
			MaxInFlight: 1,
		},
	}
}

func newBrokerJob(plgn *Plugin) enaml.InstanceJob {
	return enaml.InstanceJob{
		Name:    "cf-mysql-broker",
		Release: "cf-mysql",
		Properties: &cf_mysql_broker.CfMysqlBrokerJob{
			Broker: &cf_mysql_broker.Broker{
				QuotaEnforcer: &cf_mysql_broker.QuotaEnforcer{
					Password: plgn.BrokerQuotaEnforcerPassword,
					Pause:    brokerQuotaPause,
				},
			},
			Networks: &cf_mysql_broker.Networks{
				BrokerNetwork: plgn.NetworkName,
			},
			SslEnabled:        true,
			SkipSslValidation: true,
		},
	}
}
