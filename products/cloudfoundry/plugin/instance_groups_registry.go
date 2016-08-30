package cloudfoundry

import (
	"github.com/codegangsta/cli"
	"github.com/enaml-ops/enaml"
)

var (
	factories []InstanceGrouperFactory

	configFactories []InstanceGrouperConfigFactory
)

//InstanceGroupCreator creates and validates InstanceGroups.
type InstanceGroupCreator interface {
	ToInstanceGroup() *enaml.InstanceGroup
}

// InstanceGrouper creates and validates InstanceGroups.
type InstanceGrouper interface {
	ToInstanceGroup() (ig *enaml.InstanceGroup)
	HasValidValues() bool
}

// InstanceGrouperFactory is a function that creates InstanceGroupCreator from Config.
type InstanceGrouperFactory func(*Config) InstanceGroupCreator

type InstanceGrouperConfigFactory func(*cli.Context, *Config) InstanceGrouper

// RegisterInstanceGrouperFactory registers an InstanceGrouperFactory.
// InstanceGrouperFactories should generally be registered in their package's
// init() function.
func RegisterInstanceGrouperFactory(igf InstanceGrouperFactory) {
	factories = append(factories, igf)
}

func RegisterInstanceGrouperConfigFactory(igf InstanceGrouperConfigFactory) {
	configFactories = append(configFactories, igf)
}
