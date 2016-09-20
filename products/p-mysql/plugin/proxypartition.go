package pmysql

import (
	"strings"

	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-product-bundle/products/p-mysql/enaml-gen/proxy"
)

func NewProxyPartition(plgn *Plugin) *enaml.InstanceGroup {
	return &enaml.InstanceGroup{
		Name:               "proxy-partition",
		Lifecycle:          "service",
		Instances:          len(plgn.ProxyIPs),
		VMType:             plgn.VMTypeName,
		AZs:                plgn.AZs,
		Stemcell:           plgn.StemcellName,
		PersistentDiskType: plgn.DiskTypeName,
		Jobs: []enaml.InstanceJob{
			newProxyJob(plgn),
		},
		Networks: []enaml.Network{
			enaml.Network{
				Name:      plgn.NetworkName,
				StaticIPs: plgn.ProxyIPs,
				Default:   []interface{}{"dns", "gateway"},
			},
		},
		Update: enaml.Update{
			MaxInFlight: 1,
		},
	}
}

func newProxyJob(plgn *Plugin) enaml.InstanceJob {
	return enaml.InstanceJob{
		Name:    "proxy",
		Release: CFMysqlReleaseName,
		Properties: &proxy.ProxyJob{
			ExternalHost: strings.Join([]string{"p-mysql", "sys", plgn.BaseDomain}, "."),
			ClusterIps:   plgn.IPs,
			Nats: &proxy.Nats{
				Machines: plgn.ProxyIPs,
				Password: plgn.NatsPassword,
				User:     plgn.NatsUser,
				Port:     plgn.NatsPort,
			},
			SyslogAggregator: &proxy.SyslogAggregator{
				Address:   plgn.SyslogAddress,
				Port:      plgn.SyslogPort,
				Transport: plgn.SyslogTransport,
			},
			Proxy: &proxy.Proxy{
				ProxyIps:    plgn.ProxyIPs,
				ApiUsername: plgn.ProxyAPIUser,
				ApiPassword: plgn.ProxyAPIPass,
			},
		},
	}
}
