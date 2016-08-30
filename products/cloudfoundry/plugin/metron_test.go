package cloudfoundry_test

import (
	"github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/metron_agent"
	. "github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/plugin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Metron", func() {
	Context("when initialized WITH a complete set of arguments", func() {
		var metron *Metron
		BeforeEach(func() {
			metron = NewMetron(BuildConfig())
		})
		It("then it should allow the user to configure the metron agent", func() {
			job := metron.CreateJob()
			Ω(job).ShouldNot(BeNil())
			props, _ := job.Properties.(*metron_agent.MetronAgentJob)
			Ω(props.MetronAgent.Zone).Should(Equal("metronzoneguid"))
			Ω(props.SyslogDaemonConfig.Address).Should(Equal("syslog-server"))
			Ω(props.SyslogDaemonConfig.Port).Should(Equal(10601))
			Ω(props.SyslogDaemonConfig.Transport).Should(Equal("tcp"))
			Ω(props.MetronEndpoint.SharedSecret).Should(Equal("metronsecret"))
			Ω(props.Loggregator.Etcd.Machines).Should(Equal([]string{"10.0.1.2", "10.0.1.3", "10.0.1.4"}))
		})
	})
})
