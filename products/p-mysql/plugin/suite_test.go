package pmysql_test

import (
	"github.com/enaml-ops/enaml"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Test Suite")
}

func checkJobExists(jobs []enaml.InstanceJob, name string) bool {
	for _, j := range jobs {
		if j.Name == name {
			return true
		}
	}
	return false
}

func checkGroupExists(groups []*enaml.InstanceGroup, name string) bool {
	for _, g := range groups {

		if g.Name == name {
			return true
		}
	}
	return false
}
