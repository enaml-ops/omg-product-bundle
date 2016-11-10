package redis

import (
	"fmt"
	"strings"

	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/pluginlib/cred"
	"github.com/enaml-ops/pluginlib/pcli"
	"github.com/enaml-ops/pluginlib/pluginutil"
	"github.com/enaml-ops/pluginlib/productv1"
	"github.com/xchapter7x/lo"
	"gopkg.in/urfave/cli.v2"
)

const (
	StemcellName    = "trusty"
	StemcellAlias   = "trusty"
	StemcellVersion = "3263.8"
	BoshReleaseURL  = "https://bosh.io/d/github.com/cloudfoundry-community/redis-boshrelease?v=12"
	BoshReleaseVer  = "12"
	BoshReleaseSHA  = "324910eaf68e8803ad2317d5a2f5f6a06edc0a40"
	Master          = iota
	Slave
	Errand
	Pool
)

type jobBucket struct {
	JobName   string
	JobType   int
	Instances int
}

type Plugin struct {
	PluginVersion   string   `omg:"-"`
	LeaderIP        []string `omg:"leader-ip"`
	LeaderInstances int
	RedisPassword   string `omg:"redis-pass"`
	PoolInstances   int
	DiskSize        string
	SlaveInstances  int
	ErrandInstances int
	SlaveIP         []string `omg:"slave-ip"`
	NetworkName     string
	VMSize          string `omg:"vm-size"`
	StemcellURL     string `omg:"stemcell-url"`
	StemcellVersion string `omg:"stemcell-ver"`
	StemcellSHA     string `omg:"stemcell-sha"`
	StemcellName    string
}

func (s *Plugin) GetFlags() (flags []pcli.Flag) {
	return []pcli.Flag{
		pcli.CreateStringSliceFlag("leader-ip", "multiple static ips for each redis leader vm"),
		pcli.CreateIntFlag("leader-instances", "the number of leader instances to provision", "1"),
		pcli.CreateStringFlag("redis-pass", "the password to use for connecting redis nodes", "red1s"),
		pcli.CreateIntFlag("pool-instances", "number of instances in the redis cluster", "2"),
		pcli.CreateStringFlag("disk-size", "size of disk on VMs", "4096"),
		pcli.CreateIntFlag("slave-instances", "number of slave VMs", "1"),
		pcli.CreateIntFlag("errand-instances", "number of errand VMs", "1"),
		pcli.CreateStringSliceFlag("slave-ip", "list of slave VM Ips"),
		pcli.CreateStringFlag("network-name", "name of your target network"),
		pcli.CreateStringFlag("vm-size", "name of your desired vm size"),
		pcli.CreateStringFlag("stemcell-url", "the url of the stemcell you wish to use"),
		pcli.CreateStringFlag("stemcell-ver", "the version number of the stemcell you wish to use"),
		pcli.CreateStringFlag("stemcell-sha", "the sha of the stemcell you will use"),
		pcli.CreateStringFlag("stemcell-name", "the name of the stemcell you will use", s.GetMeta().Stemcell.Name),
	}
}

func (s *Plugin) GetMeta() product.Meta {
	return product.Meta{
		Name: "redis",
		Stemcell: enaml.Stemcell{
			Name:    StemcellName,
			Alias:   StemcellAlias,
			Version: StemcellVersion,
		},
		Releases: []enaml.Release{
			enaml.Release{
				Name:    "redis",
				Version: BoshReleaseVer,
				URL:     BoshReleaseURL,
				SHA1:    BoshReleaseSHA,
			},
		},
		Properties: map[string]interface{}{
			"version":       s.PluginVersion,
			"redis-release": strings.Join([]string{BoshReleaseURL, BoshReleaseVer, BoshReleaseSHA}, " / "),
		},
	}
}

func (s *Plugin) GetProduct(args []string, cloudConfig []byte, cs cred.Store) ([]byte, error) {
	c := pluginutil.NewContext(args, pluginutil.ToCliFlagArray(s.GetFlags()))

	err := pcli.UnmarshalFlags(s, c)
	if err != nil {
		return nil, err
	}

	err = s.cloudconfigValidation(c, enaml.NewCloudConfigManifest(cloudConfig))

	if err != nil {
		lo.G.Error("invalid settings for cloud config on target bosh: ", err)
		return nil, err
	}

	lo.G.Debug("context", c)
	var dm = new(enaml.DeploymentManifest)
	dm.SetName("enaml-redis")
	dm.Update = enaml.Update{
		Canaries:        1,
		CanaryWatchTime: "1000-100000",
		MaxInFlight:     50,
		UpdateWatchTime: "1000-100000",
	}
	dm.Properties = enaml.Properties{
		"redis": struct{}{},
	}
	dm.AddRemoteRelease("redis", BoshReleaseVer, BoshReleaseURL, BoshReleaseSHA)
	dm.AddRemoteStemcell(s.StemcellName, s.StemcellName, s.StemcellVersion, s.StemcellURL, s.StemcellSHA)

	for _, bkt := range []jobBucket{
		jobBucket{JobName: "redis_leader_z1", JobType: Master, Instances: s.LeaderInstances},
		jobBucket{JobName: "redis_z1", JobType: Pool, Instances: s.PoolInstances},
		jobBucket{JobName: "redis_test_slave_z1", JobType: Slave, Instances: s.SlaveInstances},
		jobBucket{JobName: "acceptance-tests", JobType: Errand, Instances: s.ErrandInstances},
	} {
		dm.AddJob(NewRedisJob(
			bkt.JobName,
			s.NetworkName,
			s.RedisPassword,
			s.DiskSize,
			s.VMSize,
			s.LeaderIP,
			s.SlaveIP,
			bkt.Instances,
			bkt.JobType,
		))
	}
	return dm.Bytes(), err
}

func (s *Plugin) cloudconfigValidation(c *cli.Context, cloudConfig *enaml.CloudConfigManifest) (err error) {
	lo.G.Debug("running cloud config validation")
	var vmsize = s.VMSize
	var disksize = s.DiskSize
	var netname = s.NetworkName

	for _, vmtype := range cloudConfig.VMTypes {
		err = fmt.Errorf("vm size %s does not exist in cloud config. options are: %v", vmsize, cloudConfig.VMTypes)
		if vmtype.Name == vmsize {
			err = nil
			break
		}
	}

	for _, disktype := range cloudConfig.DiskTypes {
		err = fmt.Errorf("disk size %s does not exist in cloud config. options are: %v", disksize, cloudConfig.DiskTypes)
		if disktype.Name == disksize {
			err = nil
			break
		}
	}

	for _, net := range cloudConfig.Networks {
		err = fmt.Errorf("network %s does not exist in cloud config. options are: %v", netname, cloudConfig.Networks)
		if net.(map[interface{}]interface{})["name"] == netname {
			err = nil
			break
		}
	}

	if len(cloudConfig.VMTypes) == 0 {
		err = fmt.Errorf("no vm sizes found in cloud config")
	}

	if len(cloudConfig.DiskTypes) == 0 {
		err = fmt.Errorf("no disk sizes found in cloud config")
	}

	if len(cloudConfig.Networks) == 0 {
		err = fmt.Errorf("no networks found in cloud config")
	}
	return
}

func NewRedisJob(name, networkName, pass, disk, vmSize string, masterIPs, slaveIPs []string, instances int, jobType int) (job enaml.Job) {
	var lifecycle string
	network := enaml.Network{Name: networkName}
	properties := enaml.Properties{
		"network": networkName,
		"redis": map[string]interface{}{
			"password": pass,
		},
	}
	template := enaml.Template{Name: "redis", Release: "redis"}

	switch jobType {
	case Master:
		network.StaticIPs = masterIPs

	case Slave:
		network.StaticIPs = slaveIPs
		properties["redis"].(map[string]interface{})["master"] = masterIPs[0]
		properties["redis"].(map[string]interface{})["slave"] = slaveIPs[0]

	case Errand:
		lifecycle = "errand"
		properties["redis"].(map[string]interface{})["master"] = masterIPs[0]
		properties["redis"].(map[string]interface{})["slave"] = slaveIPs[0]
		template = enaml.Template{Name: "acceptance-tests", Release: "redis"}

	default:
		properties["redis"].(map[string]interface{})["master"] = masterIPs[0]
	}

	job = enaml.Job{
		Name:      name,
		Lifecycle: lifecycle,
		Instances: instances,
		Networks: []enaml.Network{
			network,
		},
		Templates:      []enaml.Template{template},
		PersistentDisk: disk,
		ResourcePool:   vmSize,
		Update: enaml.Update{
			Canaries: 10,
		},
		Properties: make(map[string]interface{}),
	}
	job.AddProperty("redis", properties["redis"])
	job.AddProperty("network", properties["network"])
	lo.G.Debug("job", job)
	return
}
