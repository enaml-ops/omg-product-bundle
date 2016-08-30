package cloudfoundry_test

import (
	"github.com/enaml-ops/enaml"
	ltc "github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/loggregator_trafficcontroller"
	"github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/metron_agent"
	"github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/route_registrar"
	. "github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/plugin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("given the loggregator traffic controller partition", func() {

	Context("when initialized with a complete set of arguments", func() {
		var grouper InstanceGroupCreator
		var dm *enaml.DeploymentManifest
		BeforeEach(func() {

			config := &Config{
				AZs:                []string{"eastprod-1"},
				StemcellName:       "cool-ubuntu-animal",
				NetworkName:        "foundry-net",
				SystemDomain:       "sys.yourdomain.com",
				SkipSSLCertVerify:  false,
				NATSUser:           "natsuser",
				NATSPassword:       "natspass",
				NATSPort:           4222,
				NATSMachines:       []string{"10.0.0.10", "10.0.0.11"},
				LoggregratorIPs:    []string{"10.0.0.39", "10.0.0.40"},
				LoggregratorVMType: "vmtype",
				EtcdMachines:       []string{"10.0.1.2", "10.0.1.3", "10.0.1.4"},
				DopplerSecret:      "dopplersecret",
				MetronSecret:       "metronsecret",
				MetronZone:         "metronzoneguid",
			}
			grouper = NewLoggregatorTrafficController(config)

			dm = new(enaml.DeploymentManifest)
			dm.AddInstanceGroup(grouper.ToInstanceGroup())
		})

		It("then it should allow the user to configure the AZs", func() {
			ig := dm.GetInstanceGroupByName("loggregator_trafficcontroller-partition")
			Ω(len(ig.AZs)).Should(Equal(1))
			Ω(ig.AZs[0]).Should(Equal("eastprod-1"))
		})

		It("then it should allow the user to configure the used stemcell", func() {
			ig := dm.GetInstanceGroupByName("loggregator_trafficcontroller-partition")
			Ω(ig.Stemcell).Should(Equal("cool-ubuntu-animal"))
		})

		It("then it should allow the user to configure the network to use", func() {
			ig := dm.GetInstanceGroupByName("loggregator_trafficcontroller-partition")
			network := ig.GetNetworkByName("foundry-net")
			Ω(network).ShouldNot(BeNil())
			Ω(len(network.StaticIPs)).Should(Equal(2))
			Ω(network.StaticIPs[0]).Should(Equal("10.0.0.39"))
			Ω(network.StaticIPs[1]).Should(Equal("10.0.0.40"))
		})

		It("should use the correct number of instances based on the network IPs", func() {
			ig := dm.GetInstanceGroupByName("loggregator_trafficcontroller-partition")
			network := ig.Networks[0]
			Ω(ig.Instances).Should(Equal(len(network.StaticIPs)))
		})

		It("then it should allow the user to configure the VM type", func() {
			ig := dm.GetInstanceGroupByName("loggregator_trafficcontroller-partition")
			Ω(ig.VMType).Should(Equal("vmtype"))
		})

		It("then it should have update max in flight 1", func() {
			ig := dm.GetInstanceGroupByName("loggregator_trafficcontroller-partition")
			Ω(ig.Update.MaxInFlight).Should(Equal(1))
		})

		It("then it should have correctly configured the loggregator traffic controller job", func() {
			ig := dm.GetInstanceGroupByName("loggregator_trafficcontroller-partition")
			job := ig.GetJobByName("loggregator_trafficcontroller")
			Ω(job).ShouldNot(BeNil())
			Ω(job.Release).Should(Equal(CFReleaseName))

			props := job.Properties.(*ltc.LoggregatorTrafficcontrollerJob)
			Ω(props.SystemDomain).Should(Equal("sys.yourdomain.com"))
			Ω(props.Cc.SrvApiUri).Should(Equal("https://api.sys.yourdomain.com"))
			Ω(props.Ssl.SkipCertVerify).Should(BeFalse())
			Ω(props.TrafficController.Zone).Should(Equal("metronzoneguid"))
			Ω(props.Doppler.UaaClientId).Should(Equal("doppler"))
			Ω(props.Uaa.Clients.Doppler.Secret).Should(Equal("dopplersecret"))
			Ω(props.Loggregator.Etcd.Machines).Should(ConsistOf("10.0.1.2", "10.0.1.3", "10.0.1.4"))
		})

		It("then it should have the metron_agent job", func() {
			ig := dm.GetInstanceGroupByName("loggregator_trafficcontroller-partition")
			job := ig.GetJobByName("metron_agent")
			Ω(job).ShouldNot(BeNil())
			Ω(job.Release).Should(Equal(CFReleaseName))

			props := job.Properties.(*metron_agent.MetronAgentJob)
			Ω(props.MetronAgent.Zone).Should(Equal("metronzoneguid"))
			Ω(props.MetronAgent.Deployment).Should(Equal(CFReleaseName))
			Ω(props.MetronEndpoint.SharedSecret).Should(Equal("metronsecret"))
			Ω(props.Loggregator.Etcd.Machines).Should(ConsistOf("10.0.1.2", "10.0.1.3", "10.0.1.4"))
		})

		It("then it should have the route_registrar job", func() {
			ig := dm.GetInstanceGroupByName("loggregator_trafficcontroller-partition")
			job := ig.GetJobByName("route_registrar")
			Ω(job).ShouldNot(BeNil())

			props := job.Properties.(*route_registrar.RouteRegistrarJob)
			Ω(props.RouteRegistrar.Routes).ShouldNot(BeNil())
			routes := props.RouteRegistrar.Routes.([]map[string]interface{})
			Ω(len(routes)).Should(Equal(2))
			Ω(routes[0]).Should(HaveKeyWithValue("name", "doppler"))
			Ω(routes[0]).Should(HaveKeyWithValue("port", 8081))
			Ω(routes[0]).Should(HaveKeyWithValue("registration_interval", "20s"))
			Ω(routes[0]).Should(HaveKey("uris"))
			Ω(routes[0]["uris"]).Should(ConsistOf("doppler.sys.yourdomain.com"))
			Ω(routes[1]).Should(HaveKeyWithValue("name", "loggregator"))
			Ω(routes[1]).Should(HaveKeyWithValue("port", 8080))
			Ω(routes[1]).Should(HaveKeyWithValue("registration_interval", "20s"))
			Ω(routes[1]).Should(HaveKey("uris"))
			Ω(routes[1]["uris"]).Should(ConsistOf("loggregator.sys.yourdomain.com"))

			Ω(props.Nats).ShouldNot(BeNil())
			Ω(props.Nats.User).Should(Equal("natsuser"))
			Ω(props.Nats.Password).Should(Equal("natspass"))
			Ω(props.Nats.Port).Should(Equal(4222))
			Ω(props.Nats.Machines).Should(ConsistOf("10.0.0.10", "10.0.0.11"))
		})

		It("then it should have the statsd-injector job", func() {
			ig := dm.GetInstanceGroupByName("loggregator_trafficcontroller-partition")
			job := ig.GetJobByName("statsd-injector")
			Ω(job).ShouldNot(BeNil())
			Ω(job.Release).Should(Equal(CFReleaseName))
		})
	})
})
