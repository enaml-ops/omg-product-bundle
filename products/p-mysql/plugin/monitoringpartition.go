package pmysql

import (
	"strings"

	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-product-bundle/products/p-mysql/enaml-gen/replication-canary"
)

func NewMonitoringPartition(plgn *Plugin) *enaml.InstanceGroup {
	return &enaml.InstanceGroup{
		Name:               "monitoring-partition",
		Lifecycle:          "service",
		Instances:          len(plgn.MonitoringIPs),
		VMType:             plgn.VMTypeName,
		AZs:                plgn.AZs,
		Stemcell:           plgn.StemcellName,
		PersistentDiskType: plgn.DiskTypeName,
		Jobs: []enaml.InstanceJob{
			newReplicationCanaryJob(plgn),
		},
		Networks: []enaml.Network{
			enaml.Network{
				Name:      plgn.NetworkName,
				StaticIPs: plgn.MonitoringIPs,
				Default:   []interface{}{"dns", "gateway"},
			},
		},
		Update: enaml.Update{
			MaxInFlight: 1,
		},
	}
}

func newReplicationCanaryJob(plgn *Plugin) enaml.InstanceJob {
	return enaml.InstanceJob{
		Name:    "replication-canary",
		Release: "mysql-monitoring",
		Properties: &replication_canary.ReplicationCanaryJob{
			Domain: strings.Join([]string{"sys", plgn.BaseDomain}, "."),
			SyslogAggregator: &replication_canary.SyslogAggregator{
				Address:   plgn.SyslogAddress,
				Port:      plgn.SyslogPort,
				Transport: plgn.SyslogTransport,
			},
			MysqlMonitoring: &replication_canary.MysqlMonitoring{
				RecipientEmail: plgn.NotificationRecipientEmail,
				NotifyOnly:     true,
				ReplicationCanary: &replication_canary.ReplicationCanary{
					UaaAdminClientSecret:        plgn.UaaAdminClientSecret,
					ClusterIps:                  plgn.IPs,
					CanaryUsername:              seededDBUser,
					CanaryPassword:              plgn.SeededDBPassword,
					NotificationsClientUsername: notificationClientUsername,
					NotificationsClientSecret:   plgn.NotificationClientSecret,
					SwitchboardCount:            switchboardCount,
					SwitchboardUsername:         plgn.ProxyAPIUser,
					SwitchboardPassword:         plgn.ProxyAPIPass,
					PollFrequency:               pollFrequency,
					WriteReadDelay:              writeReadDelay,
				},
			},
		},
	}
}
