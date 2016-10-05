package gemfire_plugin_test

import (
	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-product-bundle/products/p-gemfire/enaml-gen/server"
	. "github.com/enaml-ops/omg-product-bundle/products/p-gemfire/plugin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("server Group", func() {
	var serverGroup *ServerGroup

	Context("when a server ip list is set", func() {
		var controlStaticIPs = []string{"1.0.0.1", "1.0.0.2"}
		var controlNetworkName = "my-network"
		var controlJobName = "server"
		var staticIPs []string
		var instanceGroup *enaml.InstanceGroup
		BeforeEach(func() {
			serverGroup = NewServerGroup(controlNetworkName, controlStaticIPs)
			instanceGroup = serverGroup.GetInstanceGroup()
			staticIPs = instanceGroup.GetNetworkByName(controlNetworkName).StaticIPs
		})

		It("should create an instance group with static IPs", func() {
			Expect(staticIPs).Should(Equal(controlStaticIPs))
		})

		It("should create the correct # of server instances", func() {
			Expect(len(staticIPs)).Should(Equal(instanceGroup.Instances), "ips should match number of instances to be created")
		})

		It("Should create map to properties.gemfire.server.addresses", func() {
			jobProperties := instanceGroup.GetJobByName(controlJobName).Properties.(server.ServerJob)
			Ω(jobProperties.Gemfire.Locator.Addresses).Should(Equal(controlStaticIPs))
		})
	})
})
