package cloudfoundry_test

import (
	"github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/haproxy"
	. "github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/plugin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("HaProxy Partition", func() {
	Context("when initialized WITHOUT a complete set of arguments", func() {
		It("then HasValidValues should return false", func() {
			plugin := new(Plugin)
			c := plugin.GetContext([]string{
				"cloudfoundry",
				"--skip-haproxy=false",
			})
			haproxyPartition := NewHaProxyPartition(c)
			Ω(haproxyPartition.HasValidValues()).Should(BeFalse())
		})
	})
	Context("when initialized WITH a complete set of arguments", func() {
		var haproxyPartition InstanceGrouper
		BeforeEach(func() {
			plugin := new(Plugin)
			c := plugin.GetContext([]string{
				"cloudfoundry",
				"--stemcell-name", "cool-ubuntu-animal",
				"--az", "eastprod-1",
				"--skip-haproxy=false",
				"--haproxy-ip", "1.0.11.1",
				"--haproxy-ip", "1.0.11.2",
				"--haproxy-ip", "1.0.11.3",
				"--network", "foundry-net",
				"--haproxy-vm-type", "blah",
				"--router-ip", "1.0.0.1",
				"--router-ip", "1.0.0.2",
				"--haproxy-sslpem", "blah",
			})
			haproxyPartition = NewHaProxyPartition(c)
		})
		It("then HasValidValues should be true", func() {
			Ω(haproxyPartition.HasValidValues()).Should(BeTrue())
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
