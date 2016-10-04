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

	Context("when a locator ip list is set", func() {
		var controlStaticIPs = []string{"1.0.0.1", "1.0.0.2"}
		var controlNetworkName = "my-network"
		var controlJobName = "locator"
		var staticIPs []string
		var instanceGroup *enaml.InstanceGroup
		BeforeEach(func() {
			locatorGroup = NewLocatorGroup(controlNetworkName, controlStaticIPs)
			instanceGroup = locatorGroup.GetInstanceGroup()
			staticIPs = instanceGroup.GetNetworkByName(controlNetworkName).StaticIPs
		})

		It("should create an instance group with static IPs", func() {
			Expect(staticIPs).Should(Equal(controlStaticIPs))
		})

		It("should create the correct # of Locator instances", func() {
			Expect(len(staticIPs)).Should(Equal(instanceGroup.Instances), "ips should match number of instances to be created")
		})

		It("Should create map to properties.gemfire.locator.addresses", func() {
			fakeProperties := getFakeGemfirePropertiesWithAddresses(controlStaticIPs)
			jobProperties := instanceGroup.GetJobByName(controlJobName).Properties
			Ω(jobProperties).Should(Equal(fakeProperties))
		})
	})
})

func getFakeGemfirePropertiesWithAddresses(addresses []string) locator.LocatorJob {
	return locator.LocatorJob{
		Gemfire: &locator.Gemfire{
			Locator: &locator.Locator{
				Addresses: addresses,
			},
		},
	}
}
