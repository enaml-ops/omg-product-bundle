package cloudfoundry_test

import (
	"github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/consul_agent"
	. "github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/plugin"
	"github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/plugin/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Consul Partition", func() {
	Context("when initialized WITH a complete set of arguments", func() {
		var err error
		var consul InstanceGroupCreator
		BeforeEach(func() {

			config := &config.Config{
				StemcellName:    "cool-ubuntu-animal",
				AZs:             []string{"eastprod-1"},
				NetworkName:     "foundry-net",
				DopplerZone:      "DopplerZoneguid",
				SyslogAddress:   "syslog-server",
				SyslogPort:      10601,
				SyslogTransport: "tcp",
				Secret:          config.Secret{},
				User:            config.User{},
				Certs:           &config.Certs{},
				InstanceCount:   config.InstanceCount{},
				IP:              config.IP{},
			}
			config.EtcdMachines = []string{"1.0.0.7", "1.0.0.8"}
			config.ConsulEncryptKeys = []string{"encyption-key"}
			config.ConsulCaCert = "ca-cert"
			config.ConsulAgentCert = "agent-cert"
			config.ConsulAgentKey = "agent-key"
			config.ConsulServerCert = "server-cert"
			config.ConsulServerKey = "server-key"
			config.ConsulIPs = []string{"1.0.0.1", "1.0.0.2"}
			config.ConsulVMType = "blah"
			config.DopplerSharedSecret = "metronsecret"

			consul = NewConsulPartition(config)
		})
		It("then it should not return an error", func() {
			Ω(err).Should(BeNil())
		})
		It("then it should allow the user to configure the consul IPs", func() {
			ig := consul.ToInstanceGroup()
			network := ig.Networks[0]
			Ω(len(network.StaticIPs)).Should(Equal(2))
			Ω(network.StaticIPs).Should(ConsistOf("1.0.0.1", "1.0.0.2"))
		})
		It("then it should have 2 instances", func() {
			ig := consul.ToInstanceGroup()
			Ω(ig.Instances).Should(Equal(2))
		})
		It("then it should configure the correct number of instances automatically from the count of given IPs", func() {
			ig := consul.ToInstanceGroup()
			network := ig.Networks[0]
			Ω(len(network.StaticIPs)).Should(Equal(ig.Instances))
		})

		It("then it should allow the user to configure the AZs", func() {
			ig := consul.ToInstanceGroup()
			Ω(len(ig.AZs)).Should(Equal(1))
			Ω(ig.AZs[0]).Should(Equal("eastprod-1"))
		})

		It("then it should allow the user to configure vm-type", func() {
			ig := consul.ToInstanceGroup()
			Ω(ig.VMType).ShouldNot(BeEmpty())
			Ω(ig.VMType).Should(Equal("blah"))
		})

		It("then it should allow the user to configure network to use", func() {
			ig := consul.ToInstanceGroup()
			network := ig.GetNetworkByName("foundry-net")
			Ω(network).ShouldNot(BeNil())
		})

		It("then it should allow the user to configure the used stemcell", func() {
			ig := consul.ToInstanceGroup()
			Ω(ig.Stemcell).ShouldNot(BeEmpty())
			Ω(ig.Stemcell).Should(Equal("cool-ubuntu-animal"))
		})
		It("then it should have update max in-flight 1 and serial", func() {
			ig := consul.ToInstanceGroup()
			Ω(ig.Update.MaxInFlight).Should(Equal(1))
			Ω(ig.Update.Serial).Should(Equal(true))
		})

		It("then it should then have 3 jobs", func() {
			ig := consul.ToInstanceGroup()
			Ω(len(ig.Jobs)).Should(Equal(3))
		})
		It("then it should then have consul agent job", func() {
			ig := consul.ToInstanceGroup()
			job := ig.GetJobByName("consul_agent")
			Ω(job).ShouldNot(BeNil())
			props, _ := job.Properties.(*consul_agent.ConsulAgentJob)
			Ω(props.Consul.ServerKey).Should(Equal("server-key"))
			Ω(props.Consul.ServerCert).Should(Equal("server-cert"))
			Ω(props.Consul.AgentCert).Should(Equal("agent-cert"))
			Ω(props.Consul.AgentKey).Should(Equal("agent-key"))
			Ω(props.Consul.CaCert).Should(Equal("ca-cert"))
			Ω(props.Consul.EncryptKeys).Should(Equal([]string{"encyption-key"}))
			agent := props.Consul.Agent
			Ω(agent.Servers.Lan).Should(Equal([]string{"1.0.0.1", "1.0.0.2"}))
		})
	})
})
