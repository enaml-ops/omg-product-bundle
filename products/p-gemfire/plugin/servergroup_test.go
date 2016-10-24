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

	Context("when valid server argument values are given w/o static server IPs", func() {
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
		var controlLocatorStaticIPs = []string{"1.0.0.1", "1.0.0.2"}
		var controlInstanceCount = 6
		var controlArpCleanerJobName = "arp-cleaner"
		var controlReleaseName = "GemFire"
		var controlDevRestPort = 7070
		var controlDevRestActive = false

		BeforeEach(func() {
			locatorGroup = NewLocatorGroup(controlNetworkName, controlLocatorStaticIPs, controlPort, controlRestPort, controlVMMemory, controlVMType)
			serverGroup = NewServerGroup(controlNetworkName, controlServerPort, controlInstanceCount, []string{}, controlVMType, controlServerVMMemory, controlDevRestPort, controlDevRestActive, locatorGroup)
			instanceGroup = serverGroup.GetInstanceGroup()
			staticIPs = instanceGroup.GetNetworkByName(controlNetworkName).StaticIPs
		})

		Context("and static server IPs are given", func() {
			var controlServerStaticIPs = []string{"1.0.0.2", "1.0.0.3"}

			BeforeEach(func() {

				serverGroup = NewServerGroup(controlNetworkName, controlServerPort, controlInstanceCount, controlServerStaticIPs, controlVMType, controlServerVMMemory, controlDevRestPort, controlDevRestActive, locatorGroup)
				instanceGroup = serverGroup.GetInstanceGroup()
				staticIPs = instanceGroup.GetNetworkByName(controlNetworkName).StaticIPs
			})

			It("should ignore the given instance count", func() {
				Ω(instanceGroup.Instances).ShouldNot(Equal(controlInstanceCount))
			})

			It("should set the instance count to the number of static server ips given", func() {
				Ω(instanceGroup.Instances).Should(Equal(len(controlServerStaticIPs)), "number of staticIPs compared to instance count")
			})

			It("should set static server IPs from the values given", func() {
				Expect(staticIPs).Should(Equal(controlServerStaticIPs))
			})
		})

		It("should contain a job for arp-cleaner from the gemfire release", func() {
			Ω(instanceGroup.GetJobByName(controlArpCleanerJobName)).ShouldNot(BeNil())
			Ω(instanceGroup.GetJobByName(controlArpCleanerJobName).Release).Should(Equal(controlReleaseName))
			Ω(instanceGroup.GetJobByName(controlArpCleanerJobName).Properties).ShouldNot(BeNil())
		})

		It("should create an instance group with static IPs for locators", func() {
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
			Ω(jobProperties.Gemfire.Locator.Addresses).Should(Equal(controlLocatorStaticIPs))
		})

		It("Should create valid job properties for cluster topology", func() {
			jobProperties := instanceGroup.GetJobByName(controlJobName).Properties.(server.ServerJob)
			Ω(jobProperties.Gemfire.ClusterTopology.NumberOfLocators).Should(Equal(len(controlLocatorStaticIPs)), "number of locators should be derived from the number of StaticIPs")
			Ω(jobProperties.Gemfire.ClusterTopology.NumberOfServers).Should(Equal(controlInstanceCount), "number of locators should be derived from the given instance count value")
		})

		It("Should create valid job properties for server configuration", func() {
			jobProperties := instanceGroup.GetJobByName(controlJobName).Properties.(server.ServerJob)
			Ω(jobProperties.Gemfire.Server.Port).Should(Equal(controlServerPort), "server port should match the user given value")
			Ω(jobProperties.Gemfire.Server.VmMemory).Should(Equal(controlServerVMMemory), "server vm memory should mathc the user given value")
		})
		Context("when given a dev rest api config value", func() {
			It("should set the values in the server", func() {
				jobProperties := instanceGroup.GetJobByName(controlJobName).Properties.(server.ServerJob)
				Ω(jobProperties.Gemfire.Server.DevRestApi.Port).Should(Equal(controlDevRestPort))
				Ω(jobProperties.Gemfire.Server.DevRestApi.Active).Should(Equal(controlDevRestActive))
			})
		})
	})
})
