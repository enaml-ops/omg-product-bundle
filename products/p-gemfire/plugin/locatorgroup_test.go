package gemfire_plugin_test

import (
	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-product-bundle/products/p-gemfire/enaml-gen/locator"
	. "github.com/enaml-ops/omg-product-bundle/products/p-gemfire/plugin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Locator Group", func() {
	var locatorGroup *LocatorGroup
	var controlVMMemory = 1024
	var controlPort = 55221
	var controlRestPort = 8080
	var controlVMType = "large"
	var controlStaticIPs = []string{"1.0.0.1", "1.0.0.2"}
	var controlNetworkName = "my-network"
	var controlJobName = "locator"
	var staticIPs []string
	var instanceGroup *enaml.InstanceGroup
	var controlArpCleanerJobName = "arp-cleaner"
	var controlReleaseName = "GemFire"

	BeforeEach(func() {
		locatorGroup = NewLocatorGroup(controlNetworkName, controlStaticIPs, controlPort, controlRestPort, controlVMMemory, controlVMType)
		instanceGroup = locatorGroup.GetInstanceGroup()
	})

	Context("when a locator ip list is set", func() {
		BeforeEach(func() {
			staticIPs = instanceGroup.GetNetworkByName(controlNetworkName).StaticIPs
		})

		It("should create an instance group with static IPs", func() {
			Expect(staticIPs).Should(Equal(controlStaticIPs))
		})

		It("should create the correct # of Locator instances", func() {
			Expect(len(staticIPs)).Should(Equal(instanceGroup.Instances), "ips should match number of instances to be created")
		})

		It("Should create map to properties.gemfire.locator.addresses", func() {
			jobProperties := instanceGroup.GetJobByName(controlJobName).Properties.(locator.LocatorJob)
			Ω(jobProperties.Gemfire.Locator.Addresses).Should(Equal(controlStaticIPs))
		})
	})

	Context("when a locator vmtype is set", func() {
		It("should set the instance groups vm type", func() {
			Ω(instanceGroup.VMType).Should(Equal(controlVMType))
		})
	})

	Context("when valid gemfire.properties are set", func() {
		var jobProperties locator.LocatorJob

		BeforeEach(func() {
			jobProperties = instanceGroup.GetJobByName("locator").Properties.(locator.LocatorJob)
		})

		It("should set a rest-port", func() {
			Ω(jobProperties.Gemfire.Locator.RestPort).Should(Equal(controlRestPort))
		})

		It("should set a port", func() {
			Ω(jobProperties.Gemfire.Locator.Port).Should(Equal(controlPort))
		})

		It("should set a vm-memory", func() {
			Ω(jobProperties.Gemfire.Locator.VmMemory).Should(Equal(controlVMMemory))
		})

		It("should set number of locators", func() {
			Ω(jobProperties.Gemfire.ClusterTopology.NumberOfLocators).Should(Equal(len(controlStaticIPs)), "this number should match the number of actual locators")
		})

		It("should set min number of locators", func() {
			Ω(jobProperties.Gemfire.ClusterTopology.MinNumberOfLocators).Should(Equal(len(controlStaticIPs)), "this number should match the number of actual locators")
		})

		It("should contain a job for arp-cleaner from the gemfire release", func() {
			Ω(instanceGroup.GetJobByName(controlArpCleanerJobName)).ShouldNot(BeNil())
			Ω(instanceGroup.GetJobByName(controlArpCleanerJobName).Release).Should(Equal(controlReleaseName))
			Ω(instanceGroup.GetJobByName(controlArpCleanerJobName).Properties).ShouldNot(BeNil())
		})
	})
})
