package cloudfoundry

import "github.com/enaml-ops/enaml"

var (
	factories []InstanceGrouperFactory
)

//InstanceGroupCreator creates and validates InstanceGroups.
type InstanceGroupCreator interface {
	ToInstanceGroup() *enaml.InstanceGroup
}

// InstanceGrouperFactory is a function that creates InstanceGroupCreator from Config.
type InstanceGrouperFactory func(*Config) InstanceGroupCreator

// RegisterInstanceGrouperFactory registers an InstanceGrouperFactory.
// InstanceGrouperFactories should generally be registered in their package's
// init() function.
func RegisterInstanceGrouperFactory(igf InstanceGrouperFactory) {
	factories = append(factories, igf)
}
