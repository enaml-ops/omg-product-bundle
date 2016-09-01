package cloudfoundry

const (
	//CFReleaseName -
	CFReleaseName = "cf"
	//StemcellName -
	StemcellName = "ubuntu-trusty"
	//StemcellAlias -
	StemcellAlias = "trusty"

	//CFLinuxFSReleaseName -
	CFLinuxFSReleaseName = "cflinuxfs2-rootfs"

	//GardenReleaseName
	GardenReleaseName = "garden-linux"

	//DiegoReleaseName
	DiegoReleaseName = "diego"

	CFMysqlReleaseName         = "cf-mysql"
	CFLinuxReleaseName         = "cflinuxfs2-rootfs"
	EtcdReleaseName            = "etcd"
	PushAppsReleaseName        = "push-apps-manager-release"
	NotificationsReleaseName   = "notifications"
	NotificationsUIReleaseName = "notifications-ui"
	CFAutoscalingReleaseName   = "cf-autoscaling"

	defaultBBSAPILocation = "bbs.service.cf.internal:8889"
)

var (
	//DeploymentName -
	DeploymentName = "cf"
	//CFReleaseVersion -
	CFReleaseVersion = "235.5.62"
	//StemcellVersion -
	StemcellVersion = "3232.17"
	//DiegoReleaseVerion
	DiegoReleaseVersion = "0.1467.29"
	//CFMysqlReleaseVersion
	CFMysqlReleaseVersion = "25.2"

	GardenReleaseVersion          = "0.338.0"
	CFLinuxReleaseVersion         = "1.26.0"
	EtcdReleaseVersion            = "48"
	PushAppsReleaseVersion        = "629.7"
	NotificationsReleaseVersion   = "24"
	NotificationsUIReleaseVersion = "17"
	CFAutoscalingReleaseVersion   = "36"

	MySQLBackupReleaseVersion   = "1"
	ServiceBackupReleaseVersion = "1"
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
