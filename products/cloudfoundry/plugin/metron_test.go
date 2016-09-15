package cloudfoundry_test

import (
	"github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/metron_agent"
	. "github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/plugin"
	"github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/plugin/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Metron", func() {
	Context("when initialized WITH a complete set of arguments", func() {
		var metron *Metron
		BeforeEach(func() {
			config := &config.Config{
				StemcellName:    "cool-ubuntu-animal",
				AZs:             []string{"eastprod-1"},
				NetworkName:     "foundry-net",
				DopplerZone:     "DopplerZoneguid",
				SyslogAddress:   "syslog-server",
				SyslogPort:      10601,
				SyslogTransport: "tcp",
				Secret:          config.Secret{},
				User:            config.User{},
				Certs:           &config.Certs{},
				InstanceCount:   config.InstanceCount{},
				IP:              config.IP{},
			}
			config.EtcdMachines = []string{"10.0.1.2", "10.0.1.3", "10.0.1.4"}
			config.ConsulEncryptKeys = []string{"encyption-key"}
			config.ConsulAgentCert = "agent-cert"
			config.ConsulAgentKey = "agent-key"
			config.ConsulServerCert = "server-cert"
			config.ConsulServerKey = "server-key"
			config.ConsulIPs = []string{"1.0.0.1", "1.0.0.2"}
			config.ConsulVMType = "blah"
			config.DopplerSharedSecret = "metronsecret"

			metron = NewMetron(config)
		})
		It("then it should allow the user to configure the metron agent", func() {
			job := metron.CreateJob()
			Ω(job).ShouldNot(BeNil())
			props, _ := job.Properties.(*metron_agent.MetronAgentJob)
			Ω(props.MetronAgent.Zone).Should(Equal("DopplerZoneguid"))
			Ω(props.SyslogDaemonConfig.Address).Should(Equal("syslog-server"))
			Ω(props.SyslogDaemonConfig.Port).Should(Equal(10601))
			Ω(props.SyslogDaemonConfig.Transport).Should(Equal("tcp"))
			Ω(props.MetronEndpoint.SharedSecret).Should(Equal("metronsecret"))
			Ω(props.Loggregator.Etcd.Machines).Should(Equal([]string{"10.0.1.2", "10.0.1.3", "10.0.1.4"}))
		})
	})
})
