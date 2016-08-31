package cloudfoundry_test

import (
	"github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/haproxy"
	. "github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/plugin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("HaProxy Partition", func() {
	Context("when initialized WITH a complete set of arguments", func() {
		var haproxyPartition InstanceGroupCreator
		BeforeEach(func() {
			config := &Config{
				StemcellName:   "cool-ubuntu-animal",
				AZs:            []string{"eastprod-1"},
				NetworkName:    "foundry-net",
				HAProxySkip:    false,
				HAProxyIPs:     []string{"1.0.11.1", "1.0.11.2", "1.0.11.3"},
				HAProxyVMType:  "blah",
				RouterMachines: []string{"1.0.0.1", "1.0.0.2"},
				HAProxySSLPem:  "blah",
			}
			haproxyPartition = NewHaProxyPartition(config)
		})
		It("then it should allow the user to configure the haproxy IPs", func() {
			ig := haproxyPartition.ToInstanceGroup()
			network := ig.Networks[0]
			Ω(len(network.StaticIPs)).Should(Equal(3))
			Ω(network.StaticIPs).Should(ConsistOf("1.0.11.1", "1.0.11.2", "1.0.11.3"))
		})
		It("then it should have 3 instances", func() {
			ig := haproxyPartition.ToInstanceGroup()
			Ω(ig.Instances).Should(Equal(3))
		})
		It("then it should configure the correct number of instances automatically from the count of given IPs", func() {
			ig := haproxyPartition.ToInstanceGroup()
			network := ig.Networks[0]
			Ω(len(network.StaticIPs)).Should(Equal(ig.Instances))
		})

		It("then it should allow the user to configure the AZs", func() {
			ig := haproxyPartition.ToInstanceGroup()
			Ω(len(ig.AZs)).Should(Equal(1))
			Ω(ig.AZs[0]).Should(Equal("eastprod-1"))
		})

		It("then it should allow the user to configure vm-type", func() {
			ig := haproxyPartition.ToInstanceGroup()
			Ω(ig.VMType).ShouldNot(BeEmpty())
			Ω(ig.VMType).Should(Equal("blah"))
		})

		It("then it should allow the user to configure network to use", func() {
			ig := haproxyPartition.ToInstanceGroup()
			network := ig.GetNetworkByName("foundry-net")
			Ω(network).ShouldNot(BeNil())
		})

		It("then it should allow the user to configure the used stemcell", func() {
			ig := haproxyPartition.ToInstanceGroup()
			Ω(ig.Stemcell).ShouldNot(BeEmpty())
			Ω(ig.Stemcell).Should(Equal("cool-ubuntu-animal"))
		})
		It("then it should have update max in-flight 1 and serial false", func() {
			ig := haproxyPartition.ToInstanceGroup()
			Ω(ig.Update.MaxInFlight).Should(Equal(1))
			Ω(ig.Update.Serial).Should(Equal(false))
		})

		It("then it should then have 4 job", func() {
			ig := haproxyPartition.ToInstanceGroup()
			Ω(len(ig.Jobs)).Should(Equal(4))
		})
		It("then it should then have haproxy job", func() {
			ig := haproxyPartition.ToInstanceGroup()
			job := ig.GetJobByName("haproxy")
			Ω(job).ShouldNot(BeNil())
			props, _ := job.Properties.(*haproxy.HaproxyJob)
			Ω(props).ShouldNot(BeNil())
			Ω(props.RequestTimeoutInSeconds).Should(Equal(180))
			Ω(props.HaProxy).ShouldNot(BeNil())
			Ω(props.HaProxy.SslPem).Should(Equal("blah"))
			Ω(props.HaProxy.DisableHttp).Should(BeTrue())
			Ω(props.Cc).ShouldNot(BeNil())
			Ω(props.Cc.AllowAppSshAccess).Should(BeTrue())
			Ω(props.Router).ShouldNot(BeNil())
			Ω(props.Router.Servers).ShouldNot(BeNil())
			Ω(props.Router.Servers.Z1).Should(ConsistOf("1.0.0.1", "1.0.0.2"))
		})
	})
})
