package cloudfoundry_test

import (
	"github.com/codegangsta/cli"
	"github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/consul_agent"
	. "github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/plugin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Consul Agent", func() {
	Context("when initialized WITHOUT a complete set of arguments", func() {
		It("then hasValidValues should return false", func() {
			plugin := new(Plugin)
			c := plugin.GetContext([]string{
				"cloudfoundry",
			})
			consulAgent := NewConsulAgent(c, []string{}, &Config{})
			Ω(consulAgent.HasValidValues()).Should(BeFalse())
		})
	})
	Context("when initialized WITH a complete set of arguments", func() {

		var c *cli.Context
		var config *Config
		BeforeEach(func() {
			plugin := new(Plugin)
			c = plugin.GetContext([]string{
				"cloudfoundry",
				"--consul-ip", "1.0.0.1",
				"--consul-ip", "1.0.0.2",
			})
			config = &Config{
				ConsulAgentCert:   "agent-cert",
				ConsulAgentKey:    "agent-key",
				ConsulServerCert:  "server-cert",
				ConsulEncryptKeys: []string{"encyption-key"},
				ConsulServerKey:   "server-key",
				ConsulCaCert:      "ca-cert",
			}
		})
		It("then hasValidValues should return true for consul with server false", func() {
			consulAgent := NewConsulAgent(c, []string{}, config)
			Ω(consulAgent.HasValidValues()).Should(BeTrue())
		})
		It("then hasValidValues should return true for consul with server true", func() {
			consulAgent := NewConsulAgentServer(c, config)
			Ω(consulAgent.HasValidValues()).Should(BeTrue())
		})
		It("then job properties are set properly for server false", func() {
			consulAgent := NewConsulAgent(c, []string{}, config)
			job := consulAgent.CreateJob()
			Ω(job).ShouldNot(BeNil())
			props := job.Properties.(*consul_agent.ConsulAgentJob)
			Ω(props.Consul.Agent.Servers.Lan).Should(ConsistOf("1.0.0.1", "1.0.0.2"))
			Ω(props.Consul.AgentCert).Should(Equal("agent-cert"))
			Ω(props.Consul.AgentKey).Should(Equal("agent-key"))
			Ω(props.Consul.ServerCert).Should(Equal("server-cert"))
			Ω(props.Consul.ServerKey).Should(Equal("server-key"))
			Ω(props.Consul.EncryptKeys).Should(ConsistOf("encyption-key"))
			Ω(props.Consul.Agent.Domain).Should(Equal("cf.internal"))
			Ω(props.Consul.Agent.Mode).Should(BeNil())
		})
		It("then job properties are set properly etcd service", func() {
			consulAgent := NewConsulAgent(c, []string{"etcd"}, config)
			job := consulAgent.CreateJob()
			Ω(job).ShouldNot(BeNil())
			props := job.Properties.(*consul_agent.ConsulAgentJob)
			etcdMap := make(map[string]map[string]string)
			etcdMap["etcd"] = make(map[string]string)
			Ω(props.Consul.Agent.Services).Should(Equal(etcdMap))
		})
		It("then job properties are set properly etcd and uaa service", func() {
			consulAgent := NewConsulAgent(c, []string{"etcd", "uaa"}, config)
			job := consulAgent.CreateJob()
			Ω(job).ShouldNot(BeNil())
			props := job.Properties.(*consul_agent.ConsulAgentJob)
			servicesMap := make(map[string]map[string]string)
			servicesMap["etcd"] = make(map[string]string)
			servicesMap["uaa"] = make(map[string]string)
			Ω(props.Consul.Agent.Services).Should(Equal(servicesMap))
		})
		It("then job properties are set properly for server true", func() {
			consulAgent := NewConsulAgentServer(c, config)
			job := consulAgent.CreateJob()
			Ω(job).ShouldNot(BeNil())
			props := job.Properties.(*consul_agent.ConsulAgentJob)
			Ω(props.Consul.Agent.Servers.Lan).Should(ConsistOf("1.0.0.1", "1.0.0.2"))
			Ω(props.Consul.AgentCert).Should(Equal("agent-cert"))
			Ω(props.Consul.AgentKey).Should(Equal("agent-key"))
			Ω(props.Consul.ServerCert).Should(Equal("server-cert"))
			Ω(props.Consul.ServerKey).Should(Equal("server-key"))
			Ω(props.Consul.EncryptKeys).Should(ConsistOf("encyption-key"))
			Ω(props.Consul.Agent.Domain).Should(Equal("cf.internal"))
			Ω(props.Consul.Agent.Mode).Should(Equal("server"))
		})

	})
})
