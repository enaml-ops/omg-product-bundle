package pmysql

import (
	"strings"

	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-product-bundle/products/p-mysql/enaml-gen/mysql"
	"github.com/enaml-ops/omg-product-bundle/products/p-mysql/enaml-gen/send-email"
	"github.com/enaml-ops/omg-product-bundle/products/p-mysql/enaml-gen/streaming-mysql-backup-tool"
)

func NewMysqlPartition(plgn *Plugin) *enaml.InstanceGroup {
	return &enaml.InstanceGroup{
		Name:               "mysql-partition",
		Lifecycle:          "service",
		Instances:          len(plgn.IPs),
		VMType:             plgn.VMTypeName,
		AZs:                plgn.AZs,
		Stemcell:           plgn.StemcellName,
		PersistentDiskType: plgn.DiskTypeName,
		Jobs: []enaml.InstanceJob{
			newCFMySQLJob(plgn),
			newSendEmailJob(plgn),
			newStreamingMysqlBackupToolJob(plgn),
		},
		Networks: []enaml.Network{
			enaml.Network{
				Name:      plgn.NetworkName,
				StaticIPs: plgn.IPs,
				Default:   []interface{}{"dns", "gateway"},
			},
		},
		Update: enaml.Update{
			MaxInFlight: 1,
		},
	}
}

func newSendEmailJob(plgn *Plugin) enaml.InstanceJob {
	return enaml.InstanceJob{
		Name:    "send-email",
		Release: MysqlMonitoringReleaseName,
		Properties: &send_email.SendEmailJob{
			Ssl: &send_email.Ssl{
				SkipCertVerify: true,
			},
			Domain: strings.Join([]string{"sys", plgn.BaseDomain}, "."),
			MysqlMonitoring: &send_email.MysqlMonitoring{
				RecipientEmail: plgn.NotificationRecipientEmail,
				AdminClient: &send_email.AdminClient{
					Secret: plgn.UaaAdminClientSecret,
				},
				Client: &send_email.Client{
					Username: notificationClientUsername,
					Secret:   plgn.NotificationClientSecret,
				},
			},
		},
	}
}

func newStreamingMysqlBackupToolJob(plgn *Plugin) enaml.InstanceJob {
	return enaml.InstanceJob{
		Name:    "streaming-mysql-backup-tool",
		Release: MysqlBackupReleaseName,
		Properties: &streaming_mysql_backup_tool.StreamingMysqlBackupToolJob{
			CfMysqlBackup: &streaming_mysql_backup_tool.CfMysqlBackup{
				BackupServer: &streaming_mysql_backup_tool.BackupServer{
					Port: backupServerPort,
				},
				EndpointCredentials: &streaming_mysql_backup_tool.EndpointCredentials{
					Username: plgn.BackupEndpointUser,
					Password: plgn.BackupEndpointPassword,
				},
			},
			CfMysql: &streaming_mysql_backup_tool.CfMysql{
				Mysql: &streaming_mysql_backup_tool.Mysql{
					AdminUsername: adminUsername,
					AdminPassword: plgn.AdminPassword,
				},
			},
		},
	}
}

func newCFMySQLJob(plgn *Plugin) enaml.InstanceJob {
	return enaml.InstanceJob{
		Name:    "mysql",
		Release: CFMysqlReleaseName,
		Properties: &mysql.MysqlJob{
			AdminUsername: adminUsername,
			AdminPassword: plgn.AdminPassword,
			CfMysql: &mysql.CfMysql{
				Mysql: &mysql.Mysql{
					DisableAutoSst:     true,
					InterruptNotifyCmd: "/var/vcap/jobs/send-email/bin/run",
					ClusterHealth: &mysql.ClusterHealth{
						Password: plgn.ClusterHealthPassword,
					},
					GaleraHealthcheck: &mysql.GaleraHealthcheck{
						DbPassword:       plgn.GaleraHealthcheckDBPassword,
						EndpointPassword: plgn.GaleraHealthcheckPassword,
						EndpointUsername: plgn.GaleraHealthcheckUsername,
					},
				},
			},
			ClusterIps:             plgn.IPs,
			DatabaseStartupTimeout: databaseStartupTimeout,
			InnodbBufferPoolSize:   innodbBufferPoolSize,
			MaxConnections:         maxConnections,
			WsrepDebug:             wsrepDebug,
			SeededDatabases: []map[string]string{
				map[string]string{
					"name":     seededDBName,
					"username": seededDBUser,
					"password": plgn.SeededDBPassword,
				},
			},
			SyslogAggregator: &mysql.SyslogAggregator{
				Address:   plgn.SyslogAddress,
				Port:      plgn.SyslogPort,
				Transport: plgn.SyslogTransport,
			},
		},
	}
}
