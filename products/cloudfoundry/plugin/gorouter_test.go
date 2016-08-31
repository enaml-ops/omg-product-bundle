package cloudfoundry_test

import (
	"github.com/enaml-ops/enaml"

	grtrlib "github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/gorouter"
	"github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/metron_agent"
	. "github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/plugin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Go-Router Partition", func() {
	Context("when the plugin is called by a operator with a complete set of arguments", func() {
		var deploymentManifest *enaml.DeploymentManifest
		const controlSecret = "goroutersecret"
		BeforeEach(func() {
			config := &Config{
				StemcellName:         "cool-ubuntu-animal",
				AZs:                  []string{"eastprod-1"},
				NetworkName:          "foundry-net",
				NATSUser:             "nats",
				NATSPassword:         "pass",
				NATSMachines:         []string{"1.0.0.5", "1.0.0.6"},
				NATSPort:             4222,
				RouterMachines:       []string{"1.0.0.1", "1.0.0.2"},
				GoRouterClientSecret: controlSecret,
				RouterVMType:         "blah",
				RouterSSLCert:        "@fixtures/sample.cert",
				RouterSSLKey:         "@fixtures/sample.key",
				RouterPass:           "blabadebleblahblah",
				MetronSecret:         "metronsecret",
				MetronZone:           "metronzoneguid",
				EtcdMachines:         []string{"1.0.0.7", "1.0.0.8"},
				RouterEnableSSL:      true,
				RouterUser:           "router_status",
			}
			gr := NewGoRouterPartition(config)
			deploymentManifest = new(enaml.DeploymentManifest)
			deploymentManifest.AddInstanceGroup(gr.ToInstanceGroup())
		})

		It("then it should allow the user to configure the router IPs", func() {
			ig := deploymentManifest.GetInstanceGroupByName("router-partition")
			network := ig.Networks[0]
			Ω(len(network.StaticIPs)).Should(Equal(2))
			Ω(network.StaticIPs).Should(ConsistOf("1.0.0.1", "1.0.0.2"))
		})

		It("then it should configure the correct number of instances automatically from the count of given IPs", func() {
			ig := deploymentManifest.GetInstanceGroupByName("router-partition")
			network := ig.Networks[0]
			Ω(len(network.StaticIPs)).Should(Equal(ig.Instances))
		})

		It("then it should allow the user to configure the AZs", func() {
			ig := deploymentManifest.GetInstanceGroupByName("router-partition")
			Ω(len(ig.AZs)).Should(Equal(1))
			Ω(ig.AZs[0]).Should(Equal("eastprod-1"))
		})

		It("then it should allow the user to configure vm-type", func() {
			ig := deploymentManifest.GetInstanceGroupByName("router-partition")
			Ω(ig.VMType).ShouldNot(BeEmpty())
			Ω(ig.VMType).Should(Equal("blah"))
		})

		It("then it should allow the user to configure network to use", func() {
			ig := deploymentManifest.GetInstanceGroupByName("router-partition")
			network := ig.GetNetworkByName("foundry-net")
			Ω(network).ShouldNot(BeNil())
		})

		It("then it should allow the user to configure the used stemcell", func() {
			ig := deploymentManifest.GetInstanceGroupByName("router-partition")
			Ω(ig.Stemcell).ShouldNot(BeEmpty())
			Ω(ig.Stemcell).Should(Equal("cool-ubuntu-animal"))
		})

		It("then it should allow the user to configure if we enable ssl", func() {
			ig := deploymentManifest.GetInstanceGroupByName("router-partition")
			job := ig.GetJobByName("gorouter")
			properties := job.Properties.(*grtrlib.GorouterJob)
			Ω(properties.Router.EnableSsl).Should(BeTrue())
		})

		It("then it should allow the user to configure the nats pool to use", func() {
			ig := deploymentManifest.GetInstanceGroupByName("router-partition")
			job := ig.GetJobByName("gorouter")
			properties := job.Properties.(*grtrlib.GorouterJob)
			Ω(properties.Nats.Machines).Should(ConsistOf("1.0.0.5", "1.0.0.6"))
			Ω(properties.Nats.User).Should(Equal("nats"))
			Ω(properties.Nats.Password).Should(Equal("pass"))
			Ω(properties.Nats.Port).Should(Equal(4222))
		})

		It("then it should allow the user to configure the loggregator pool to use", func() {
			ig := deploymentManifest.GetInstanceGroupByName("router-partition")
			job := ig.GetJobByName("metron_agent")
			properties := job.Properties.(*metron_agent.MetronAgentJob)
			Ω(properties.Loggregator.Etcd.Machines).Should(ConsistOf("1.0.0.7", "1.0.0.8"))
		})

		It("then it should allow the user to configure the metron agent", func() {
			ig := deploymentManifest.GetInstanceGroupByName("router-partition")
			job := ig.GetJobByName("metron_agent")
			properties := job.Properties.(*metron_agent.MetronAgentJob)
			Ω(properties.MetronAgent.Zone).Should(Equal("metronzoneguid"))
			Ω(properties.MetronAgent.Deployment).Should(Equal(DeploymentName))
			Ω(properties.MetronEndpoint.SharedSecret).Should(Equal("metronsecret"))
		})

		It("then it should allow the user to configure the router user/pass", func() {
			ig := deploymentManifest.GetInstanceGroupByName("router-partition")
			job := ig.GetJobByName("gorouter")
			properties := job.Properties.(*grtrlib.GorouterJob)
			Ω(properties.Router.Status.User).Should(Equal("router_status"))
			Ω(properties.Router.Status.Password).Should(Equal("blabadebleblahblah"))
		})

		It("then it should configure UAA", func() {
			ig := deploymentManifest.GetInstanceGroupByName("router-partition")
			job := ig.GetJobByName("gorouter")
			properties := job.Properties.(*grtrlib.GorouterJob)
			Ω(properties.Uaa).ShouldNot(BeNil())
			Ω(properties.Uaa.Ssl.Port).Should(Equal(-1))
			Ω(properties.Uaa.Clients.Gorouter.Secret).Should(Equal(controlSecret))
		})

		Context("when the plugin is called by a operator with arguments for ssl cert/key strings", func() {
			var deploymentManifest *enaml.DeploymentManifest
			BeforeEach(func() {
				gr := NewGoRouterPartition(&Config{
					RouterSSLCert: "blah",
					RouterSSLKey:  "blahblah",
				})
				deploymentManifest = new(enaml.DeploymentManifest)
				deploymentManifest.AddInstanceGroup(gr.ToInstanceGroup())
			})

			It("then it should use the provided string values directly", func() {
				ig := deploymentManifest.GetInstanceGroupByName("router-partition")
				job := ig.GetJobByName("gorouter")
				properties := job.Properties.(*grtrlib.GorouterJob)
				Ω(properties.Router.SslCert).Should(Equal("blah"))
				Ω(properties.Router.SslKey).Should(Equal("blahblah"))
			})
		})

		Context("when the plugin is called by a operator with arguments for just ssl cert/key strings", func() {
			var deploymentManifest *enaml.DeploymentManifest
			BeforeEach(func() {
				gr := NewGoRouterPartition(&Config{
					RouterSSLCert: "blah",
					RouterSSLKey:  "blahblah",
				})
				deploymentManifest = new(enaml.DeploymentManifest)
				deploymentManifest.AddInstanceGroup(gr.ToInstanceGroup())
			})

			It("then it should allow the user to configure the cert & key used from a string flag", func() {
				ig := deploymentManifest.GetInstanceGroupByName("router-partition")
				job := ig.GetJobByName("gorouter")
				properties := job.Properties.(*grtrlib.GorouterJob)
				Ω(properties.Router.SslCert).Should(Equal("blah"))
				Ω(properties.Router.SslKey).Should(Equal("blahblah"))
			})
		})
	})
})
