package cloudfoundry

const (
	CFReleaseName    = "cf"
	CFReleaseVersion = "235.5.62"

	StemcellName    = "ubuntu-trusty"
	StemcellAlias   = "trusty"
	StemcellVersion = "3232.17"

	CFLinuxReleaseName    = "cflinuxfs2-rootfs"
	CFLinuxReleaseVersion = "1.26.0"

	GardenReleaseName    = "garden-linux"
	GardenReleaseVersion = "0.338.0"

	DiegoReleaseName    = "diego"
	DiegoReleaseVersion = "0.1467.29"

	CFMysqlReleaseName    = "cf-mysql"
	CFMysqlReleaseVersion = "25.2"

	EtcdReleaseName    = "etcd"
	EtcdReleaseVersion = "48"

	PushAppsReleaseName    = "push-apps-manager-release"
	PushAppsReleaseVersion = "629.7"

	NotificationsReleaseName    = "notifications"
	NotificationsReleaseVersion = "24"

	NotificationsUIReleaseName    = "notifications-ui"
	NotificationsUIReleaseVersion = "17"

	CFAutoscalingReleaseName    = "cf-autoscaling"
	CFAutoscalingReleaseVersion = "36"

	MySQLBackupReleaseVersion   = "1"
	ServiceBackupReleaseVersion = "1"
)

const (
	DeploymentName = "cf"

	defaultBBSAPILocation = "bbs.service.cf.internal:8889"

	javaBuildpackName    = "java-offline-buildpack"
	javaBuildpackPackage = "buildpack_java_offline"
)

var flagsToInferFromCloudConfig = map[string][]string{
	"disktype": []string{
		"mysql-disk-type",
		"diego-db-disk-type",
		"diego-cell-disk-type",
		"diego-brain-disk-type",
		"etcd-disk-type",
		"nfs-disk-type",
	},
	"vmtype": []string{
		"diego-brain-vm-type",
		"errand-vm-type",
		"clock-global-vm-type",
		"doppler-vm-type",
		"uaa-vm-type",
		"diego-cell-vm-type",
		"diego-db-vm-type",
		"router-vm-type",
		"haproxy-vm-type",
		"nats-vm-type",
		"consul-vm-type",
		"etcd-vm-type",
		"nfs-vm-type",
		"mysql-vm-type",
		"mysql-proxy-vm-type",
		"cc-worker-vm-type",
		"cc-vm-type",
		"loggregator-traffic-controller-vmtype",
	},
	"az":      []string{"az"},
	"network": []string{"network"},
}
