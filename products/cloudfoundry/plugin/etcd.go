package cloudfoundry

import (
	"github.com/codegangsta/cli"
	"github.com/enaml-ops/enaml"
	etcdlib "github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/etcd"
	etcdmetricslib "github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/etcd_metrics_server"
	"github.com/xchapter7x/lo"
)

//NewEtcdPartition -
func NewEtcdPartition(c *cli.Context) (igf InstanceGrouper) {
	igf = &Etcd{
		AZs:                c.StringSlice("az"),
		StemcellName:       c.String("stemcell-name"),
		NetworkIPs:         c.StringSlice("etcd-machine-ip"),
		NetworkName:        c.String("network"),
		VMTypeName:         c.String("etcd-vm-type"),
		PersistentDiskType: c.String("etcd-disk-type"),
		Metron:             NewMetron(c),
		StatsdInjector:     NewStatsdInjector(c),
		Nats: &etcdmetricslib.Nats{
			Username: c.String("nats-user"),
			Password: c.String("nats-pass"),
			Machines: c.StringSlice("nats-machine-ip"),
		},
	}
	return
}

//ToInstanceGroup -
func (s *Etcd) ToInstanceGroup() (ig *enaml.InstanceGroup) {
	ig = &enaml.InstanceGroup{
		Name:      "etcd_server-partition",
		Instances: len(s.NetworkIPs),
		VMType:    s.VMTypeName,
		AZs:       s.AZs,
		Stemcell:  s.StemcellName,
		Jobs: []enaml.InstanceJob{
			s.newEtcdJob(),
			s.newEtcdMetricsServerJob(),
			s.Metron.CreateJob(),
			s.StatsdInjector.CreateJob(),
		},
		Networks: []enaml.Network{
			enaml.Network{Name: s.NetworkName, StaticIPs: s.NetworkIPs},
		},
		PersistentDiskType: s.PersistentDiskType,
		Update: enaml.Update{
			MaxInFlight: 1,
			Serial:      false,
		},
	}
	return
}

func (s *Etcd) newEtcdJob() enaml.InstanceJob {
	return enaml.InstanceJob{
		Name:    "etcd",
		Release: "cf",
		Properties: &etcdlib.EtcdJob{
			Etcd: &etcdlib.Etcd{
				PeerRequireSsl: false,
				RequireSsl:     false,
				Machines:       s.NetworkIPs,
			},
		},
	}
}

func (s *Etcd) newEtcdMetricsServerJob() enaml.InstanceJob {
	return enaml.InstanceJob{
		Name:    "etcd_metrics_server",
		Release: "cf",
		Properties: &etcdmetricslib.EtcdMetricsServerJob{
			EtcdMetricsServer: &etcdmetricslib.EtcdMetricsServer{
				Nats: s.Nats,
			},
		},
	}
}

//HasValidValues -
func (s *Etcd) HasValidValues() bool {

	lo.G.Debugf("checking '%s' for valid flags", "etcd")

	if len(s.AZs) <= 0 {
		lo.G.Debugf("could not find the correct number of AZs configured '%v' : '%v'", len(s.AZs), s.AZs)
	}
	if len(s.NetworkIPs) <= 0 {
		lo.G.Debugf("could not find the correct number of network ips configured '%v' : '%v'", len(s.NetworkIPs), s.NetworkIPs)
	}
	if s.StemcellName == "" {
		lo.G.Debugf("could not find a valid stemcellname '%v'", s.StemcellName)
	}
	if s.VMTypeName == "" {
		lo.G.Debugf("could not find a valid vmtypename '%v'", s.VMTypeName)
	}
	if s.PersistentDiskType == "" {
		lo.G.Debugf("could not find a valid PersistentDiskType '%v'", s.PersistentDiskType)
	}
	if s.NetworkName == "" {
		lo.G.Debugf("could not find a valid networkname '%v'", s.NetworkName)
	}
	return (len(s.AZs) > 0 &&
		s.StemcellName != "" &&
		s.VMTypeName != "" &&
		s.NetworkName != "" &&
		len(s.NetworkIPs) > 0 &&
		s.PersistentDiskType != "")
}
