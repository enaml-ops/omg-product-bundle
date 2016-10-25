package cloudfoundry

const (
	CFReleaseName    = "cf"
	CFReleaseVersion = "245"
	CFReleaseURI     = "https://bosh.io/d/github.com/cloudfoundry/cf-release?v=245"
	CFReleaseSHA     = "0c9f485f640c2b9e3136fcc89047b3d76dd6863c"

	StemcellName    = "ubuntu-trusty"
	StemcellAlias   = "trusty"
	StemcellVersion = "3232.17"

	CFLinuxReleaseName    = "cflinuxfs2-rootfs"
	CFLinuxReleaseVersion = "1.37.0"
	CFLinuxReleaseURI     = "https://bosh.io/d/github.com/cloudfoundry/cflinuxfs2-rootfs-release?v=1.37.0"
	CFLinuxReleaseSHA     = "fc2a1f3cfbb953a22f9e0bff93f964477bd921b9"

	GardenReleaseName    = "garden-linux"
	GardenReleaseVersion = "0.342.0"
	GardenReleaseURI     = "https://bosh.io/d/github.com/cloudfoundry-incubator/garden-linux-release?v=0.342.0"
	GardenReleaseSHA     = "dbfd8e7e3560286b6d8c02ba9065a50289e8e0f3"

	DiegoReleaseName    = "diego"
	DiegoReleaseVersion = "0.1487.0"
	DiegoReleaseURI     = "https://bosh.io/d/github.com/cloudfoundry/diego-release?v=0.1487.0"
	DiegoReleaseSHA     = "f173af7117baa34cff97cf4355f958ea536a45d0"

	CFMysqlReleaseName    = "cf-mysql"
	CFMysqlReleaseVersion = "32"
	CFMysqlReleaseURI     = "https://bosh.io/d/github.com/cloudfoundry/cf-mysql-release?v=32"
	CFMysqlReleaseSHA     = "a41bb2cadd4311bc9977ccc2c1fca07ba41ccef2"

	EtcdReleaseName    = "etcd"
	EtcdReleaseVersion = "78"
	EtcdReleaseURI     = "https://bosh.io/d/github.com/cloudfoundry-incubator/etcd-release?v=78"
	EtcdReleaseSHA     = "56f7cbcbf72ff2c8f9be76f27d267aef621d3175"
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
