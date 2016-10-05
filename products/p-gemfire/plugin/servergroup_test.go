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
		var controlNetworkName = "my-network"
		var controlJobName = "server"
		var staticIPs []string
		var instanceGroup *enaml.InstanceGroup
		var locatorGroup *LocatorGroup
		var controlVMMemory = 1024
		var controlServerVMMemory = 2024
		var controlPort = 55221
		var controlServerPort = 55001
		var controlRestPort = 8080
		var controlVMType = "large"
		var controlStaticIPs = []string{"1.0.0.1", "1.0.0.2"}
		var controlInstanceCount = 6

		BeforeEach(func() {
			locatorGroup = NewLocatorGroup(controlNetworkName, controlStaticIPs, controlPort, controlRestPort, controlVMMemory, controlVMType)
			serverGroup = NewServerGroup(controlNetworkName, controlServerPort, controlInstanceCount, controlVMType, controlServerVMMemory, locatorGroup)
			instanceGroup = serverGroup.GetInstanceGroup()
			staticIPs = instanceGroup.GetNetworkByName(controlNetworkName).StaticIPs
		})

		It("should create an instance group with static IPs", func() {
			Expect(staticIPs).Should(BeNil())
		})

		It("should create the correct # of server instances", func() {
			Expect(controlInstanceCount).Should(Equal(instanceGroup.Instances), "should match number of server instances requested")
		})

		It("should create the correct vmtype for our servers", func() {
			Expect(controlVMType).Should(Equal(instanceGroup.VMType))
		})

		It("Should create map to properties.gemfire.server.addresses", func() {
			jobProperties := instanceGroup.GetJobByName(controlJobName).Properties.(server.ServerJob)
			Ω(jobProperties.Gemfire.Locator.Addresses).Should(Equal(controlStaticIPs))
		})

		It("Should create valid job properties for cluster topology", func() {
			jobProperties := instanceGroup.GetJobByName(controlJobName).Properties.(server.ServerJob)
			Ω(jobProperties.Gemfire.ClusterTopology.NumberOfLocators).Should(Equal(len(controlStaticIPs)), "number of locators should be derived from the number of StaticIPs")
			Ω(jobProperties.Gemfire.ClusterTopology.NumberOfServers).Should(Equal(controlInstanceCount), "number of locators should be derived from the given instance count value")
		})

		It("Should create valid job properties for server configuration", func() {
			jobProperties := instanceGroup.GetJobByName(controlJobName).Properties.(server.ServerJob)
			Ω(jobProperties.Gemfire.Server.Port).Should(Equal(controlServerPort), "server port should match the user given value")
			Ω(jobProperties.Gemfire.Server.VmMemory).Should(Equal(controlServerVMMemory), "server vm memory should mathc the user given value")
		})
	})
})
