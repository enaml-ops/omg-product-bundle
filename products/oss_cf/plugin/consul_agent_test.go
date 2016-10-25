package cloudfoundry_test

import (
	"github.com/enaml-ops/omg-product-bundle/products/oss_cf/enaml-gen/consul_agent"
	. "github.com/enaml-ops/omg-product-bundle/products/oss_cf/plugin"
	"github.com/enaml-ops/omg-product-bundle/products/oss_cf/plugin/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Consul Agent", func() {
	Context("when initialized WITH a complete set of arguments", func() {

		var cfg *config.Config
		BeforeEach(func() {
			cfg = &config.Config{
				Secret:        config.Secret{},
				User:          config.User{},
				Certs:         &config.Certs{},
				InstanceCount: config.InstanceCount{},
				IP:            config.IP{},
			}
			cfg.ConsulAgentCert = "agent-cert"
			cfg.ConsulAgentKey = "agent-key"
			cfg.ConsulServerCert = "server-cert"
			cfg.ConsulEncryptKeys = []string{"encyption-key"}
			cfg.ConsulServerKey = "server-key"
			cfg.ConsulIPs = []string{"1.0.0.1", "1.0.0.2"}
		})
		It("then consul with server false", func() {
			consulAgent := NewConsulAgent([]string{}, cfg)
			Ω(consulAgent.Mode).Should(Equal(""))
		})
		It("then consul with server true", func() {
			consulAgent := NewConsulAgentServer(cfg)
			Ω(consulAgent.Mode).Should(Equal("server"))
		})
		It("then job properties are set properly for server false", func() {
			consulAgent := NewConsulAgent([]string{}, cfg)
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
			consulAgent := NewConsulAgent([]string{"etcd"}, cfg)
			job := consulAgent.CreateJob()
			Ω(job).ShouldNot(BeNil())
			props := job.Properties.(*consul_agent.ConsulAgentJob)
			etcdMap := make(map[string]map[string]string)
			etcdMap["etcd"] = make(map[string]string)
			Ω(props.Consul.Agent.Services).Should(Equal(etcdMap))
		})
		It("then job properties are set properly etcd and uaa service", func() {
			consulAgent := NewConsulAgent([]string{"etcd", "uaa"}, cfg)
			job := consulAgent.CreateJob()
			Ω(job).ShouldNot(BeNil())
			props := job.Properties.(*consul_agent.ConsulAgentJob)
			servicesMap := make(map[string]map[string]string)
			servicesMap["etcd"] = make(map[string]string)
			servicesMap["uaa"] = make(map[string]string)
			Ω(props.Consul.Agent.Services).Should(Equal(servicesMap))
		})
		It("then job properties are set properly for server true", func() {
			consulAgent := NewConsulAgentServer(cfg)
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
