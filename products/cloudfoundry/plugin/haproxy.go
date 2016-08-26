package cloudfoundry

import (
	"github.com/codegangsta/cli"
	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/haproxy"
	"github.com/enaml-ops/pluginlib/util"
	"github.com/xchapter7x/lo"
)

//NewHaProxyPartition -
func NewHaProxyPartition(c *cli.Context) InstanceGrouper {
	sslpem, err := pluginutil.LoadResourceFromContext(c, "haproxy-sslpem")
	if err != nil {
		lo.G.Error("couldn't load haproxy-sslpem:" + err.Error())
		return nil
	}

	return &HAProxy{
		Skip:           c.BoolT("skip-haproxy"),
		AZs:            c.StringSlice("az"),
		StemcellName:   c.String("stemcell-name"),
		NetworkIPs:     c.StringSlice("haproxy-ip"),
		NetworkName:    c.String("network"),
		VMTypeName:     c.String("haproxy-vm-type"),
		ConsulAgent:    NewConsulAgent(c, []string{}),
		Metron:         NewMetron(c),
		StatsdInjector: NewStatsdInjector(c),
		RouterMachines: c.StringSlice("router-ip"),
		SSLPem:         sslpem,
	}
}

//ToInstanceGroup -
func (s *HAProxy) ToInstanceGroup() (ig *enaml.InstanceGroup) {
	if !s.Skip {
		ig = &enaml.InstanceGroup{
			Name:      "ha_proxy-partition",
			Instances: len(s.NetworkIPs),
			VMType:    s.VMTypeName,
			AZs:       s.AZs,
			Stemcell:  s.StemcellName,
			Jobs: []enaml.InstanceJob{
				s.createHAProxyJob(),
				s.ConsulAgent.CreateJob(),
				s.Metron.CreateJob(),
				s.StatsdInjector.CreateJob(),
			},
			Networks: []enaml.Network{
				enaml.Network{Name: s.NetworkName, StaticIPs: s.NetworkIPs},
			},
			Update: enaml.Update{
				MaxInFlight: 1,
			},
		}
	}
	return
}

func (s *HAProxy) createHAProxyJob() enaml.InstanceJob {
	return enaml.InstanceJob{
		Name:    "haproxy",
		Release: "cf",
		Properties: &haproxy.HaproxyJob{
			RequestTimeoutInSeconds: 180,
			HaProxy: &haproxy.HaProxy{
				DisableHttp: true,
				SslPem:      s.SSLPem,
			},
			Router: &haproxy.Router{
				Servers: &haproxy.Servers{
					Z1: s.RouterMachines,
				},
			},
			Cc: &haproxy.Cc{
				AllowAppSshAccess: true,
			},
		},
	}
}

//HasValidValues - Check if the datastructure has valid fields
func (s *HAProxy) HasValidValues() bool {

	if s.Skip {
		lo.G.Debug("we are not using haproxy")
		return true
	}

	lo.G.Debugf("checking '%s' for valid flags", "haproxy")

	if len(s.AZs) <= 0 {
		lo.G.Debugf("could not find the correct number of AZs configured '%v' : '%v'", len(s.AZs), s.AZs)
	}
	if len(s.NetworkIPs) <= 0 {
		lo.G.Debugf("could not find the correct number of network ips configured '%v' : '%v'", len(s.NetworkIPs), s.NetworkIPs)
	}
	if len(s.RouterMachines) <= 0 {
		lo.G.Debugf("could not find the correct number of RouterMachines configured '%v' : '%v'", len(s.RouterMachines), s.RouterMachines)
	}
	if s.StemcellName == "" {
		lo.G.Debugf("could not find a valid stemcellname '%v'", s.StemcellName)
	}
	if s.VMTypeName == "" {
		lo.G.Debugf("could not find a valid vmtypename '%v'", s.VMTypeName)
	}
	if s.NetworkName == "" {
		lo.G.Debugf("could not find a valid NetworkName '%v'", s.NetworkName)
	}
	return (len(s.AZs) > 0 &&
		s.StemcellName != "" &&
		s.VMTypeName != "" &&
		s.NetworkName != "" &&
		len(s.NetworkIPs) > 0 &&
		len(s.RouterMachines) > 0 &&
		s.SSLPem != "")
}
