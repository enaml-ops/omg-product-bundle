package cloudfoundry

import (
	"gopkg.in/urfave/cli.v2"
	"github.com/enaml-ops/enaml"
	"github.com/xchapter7x/lo"
)

//NewStatsdInjector -
func NewStatsdInjector(c *cli.Context) (statsdInjector *StatsdInjector) {
	statsdInjector = &StatsdInjector{}
	return
}

//CreateJob -
func (s *StatsdInjector) CreateJob() enaml.InstanceJob {
	return enaml.InstanceJob{
		Name:       "statsd-injector",
		Release:    "cf",
		Properties: make(map[interface{}]interface{}),
	}
}

//HasValidValues -
func (s *StatsdInjector) HasValidValues() bool {
	lo.G.Debug("checking statsdinjector for valid values")
	return true
}
