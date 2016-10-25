package redis

import (
	"fmt"
	"strings"

	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/pluginlib/pcli"
	"github.com/enaml-ops/pluginlib/pluginutil"
	"github.com/enaml-ops/pluginlib/product"
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
	PluginVersion string
}

func (s *Plugin) GetFlags() (flags []pcli.Flag) {
	return []pcli.Flag{
		pcli.Flag{FlagType: pcli.StringSliceFlag, Name: "leader-ip", Usage: "multiple static ips for each redis leader vm"},
		pcli.Flag{FlagType: pcli.IntFlag, Name: "leader-instances", Value: "1", Usage: "the number of leader instances to provision"},
		pcli.Flag{FlagType: pcli.StringFlag, Name: "redis-pass", Value: "red1s", Usage: "the password to use for connecting redis nodes"},
		pcli.Flag{FlagType: pcli.IntFlag, Name: "pool-instances", Value: "2", Usage: "number of instances in the redis cluster"},
		pcli.Flag{FlagType: pcli.StringFlag, Name: "disk-size", Value: "4096", Usage: "size of disk on VMs"},
		pcli.Flag{FlagType: pcli.IntFlag, Name: "slave-instances", Value: "1", Usage: "number of slave VMs"},
		pcli.Flag{FlagType: pcli.IntFlag, Name: "errand-instances", Value: "1", Usage: "number of errand VMs"},
		pcli.Flag{FlagType: pcli.StringSliceFlag, Name: "slave-ip", Usage: "list of slave VM Ips"},
		pcli.Flag{FlagType: pcli.StringFlag, Name: "network-name", Usage: "name of your target network"},
		pcli.Flag{FlagType: pcli.StringFlag, Name: "vm-size", Usage: "name of your desired vm size"},
		pcli.Flag{FlagType: pcli.StringFlag, Name: "stemcell-url", Usage: "the url of the stemcell you wish to use"},
		pcli.Flag{FlagType: pcli.StringFlag, Name: "stemcell-ver", Usage: "the version number of the stemcell you wish to use"},
		pcli.Flag{FlagType: pcli.StringFlag, Name: "stemcell-sha", Usage: "the sha of the stemcell you will use"},
		pcli.Flag{FlagType: pcli.StringFlag, Name: "stemcell-name", Value: s.GetMeta().Stemcell.Name, Usage: "the name of the stemcell you will use"},
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

func (s *Plugin) GetProduct(args []string, cloudConfig []byte) (b []byte, err error) {
	c := pluginutil.NewContext(args, pluginutil.ToCliFlagArray(s.GetFlags()))

	if err := s.flagValidation(c); err != nil {
		lo.G.Error("invalid arguments: ", err)
		return nil, err
	}

	if err := s.cloudconfigValidation(c, enaml.NewCloudConfigManifest(cloudConfig)); err != nil {
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
	dm.AddRemoteStemcell(c.String("stemcell-name"), c.String("stemcell-name"), c.String("stemcell-ver"), c.String("stemcell-url"), c.String("stemcell-sha"))

	for _, bkt := range []jobBucket{
		jobBucket{JobName: "redis_leader_z1", JobType: Master, Instances: c.Int("leader-instances")},
		jobBucket{JobName: "redis_z1", JobType: Pool, Instances: c.Int("pool-instances")},
		jobBucket{JobName: "redis_test_slave_z1", JobType: Slave, Instances: c.Int("slave-instances")},
		jobBucket{JobName: "acceptance-tests", JobType: Errand, Instances: c.Int("errand-instances")},
	} {
		dm.AddJob(NewRedisJob(
			bkt.JobName,
			c.String("network-name"),
			c.String("redis-pass"),
			c.String("disk-size"),
			c.String("vm-size"),
			c.StringSlice("leader-ip"),
			c.StringSlice("slave-ip"),
			bkt.Instances,
			bkt.JobType,
		))
	}
	return dm.Bytes(), err
}

func (s *Plugin) cloudconfigValidation(c *cli.Context, cloudConfig *enaml.CloudConfigManifest) (err error) {
	lo.G.Debug("running cloud config validation")
	var vmsize = c.String("vm-size")
	var disksize = c.String("disk-size")
	var netname = c.String("network-name")

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

func (s *Plugin) flagValidation(c *cli.Context) (err error) {
	lo.G.Debug("validating given flags")

	if len(c.StringSlice("leader-ip")) <= 0 {
		err = fmt.Errorf("no `leader-ip` given")
	}

	if len(c.StringSlice("slave-ip")) <= 0 {
		err = fmt.Errorf("no `slave-ip` given")
	}

	if len(c.String("network-name")) <= 0 {
		err = fmt.Errorf("no `network-name` given")
	}

	if len(c.String("vm-size")) <= 0 {
		err = fmt.Errorf("no `vm-size` given")
	}

	if len(c.String("stemcell-ver")) <= 0 {
		err = fmt.Errorf("no `stemcell-ver` given")
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
