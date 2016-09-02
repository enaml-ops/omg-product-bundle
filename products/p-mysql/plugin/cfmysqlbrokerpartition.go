package pmysql

import (
	"strings"

	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-product-bundle/products/p-mysql/enaml-gen/cf-mysql-broker"
	"github.com/xchapter7x/lo"
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
	var hostValue string
	if len(plgn.ProxyIPs) >= 1 {
		hostValue = plgn.ProxyIPs[0]
	} else {
		lo.G.Error("could not find any proxy hosts defined")
	}
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
			ExternalHost:      strings.Join([]string{"p-mysql", "sys", plgn.BaseDomain}, "."),
			CcApiUri:          strings.Join([]string{"https://api", "sys", plgn.BaseDomain}, "."),
			CookieSecret:      plgn.BrokerCookieSecret,
			AuthUsername:      plgn.BrokerAuthUsername,
			AuthPassword:      plgn.BrokerAuthPassword,
			Nats: &cf_mysql_broker.Nats{
				Machines: plgn.ProxyIPs,
				Password: plgn.NatsPassword,
				User:     plgn.NatsUser,
				Port:     plgn.NatsPort,
			},
			SyslogAggregator: &cf_mysql_broker.SyslogAggregator{
				Address:   plgn.SyslogAddress,
				Port:      plgn.SyslogPort,
				Transport: plgn.SyslogTransport,
			},
			MysqlNode: &cf_mysql_broker.MysqlNode{
				Host:           hostValue,
				AdminPassword:  plgn.AdminPassword,
				PersistentDisk: brokerPersistentDisk,
			},
		},
	}
}
