package cloudfoundry_test

import (
	"github.com/enaml-ops/omg-product-bundle/products/oss_cf/enaml-gen/doppler"
	"github.com/enaml-ops/omg-product-bundle/products/oss_cf/enaml-gen/syslog_drain_binder"
	. "github.com/enaml-ops/omg-product-bundle/products/oss_cf/plugin"
	"github.com/enaml-ops/omg-product-bundle/products/oss_cf/plugin/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Doppler Partition", func() {

	Context("when initialized WITH a complete set of arguments", func() {
		var dopplerPartition InstanceGroupCreator
		BeforeEach(func() {

			config := &config.Config{
				NetworkName:                   "foundry-net",
				StemcellName:                  "cool-ubuntu-animal",
				AZs:                           []string{"eastprod-1"},
				SystemDomain:                  "sys.test.com",
				SkipSSLCertVerify:             true,
				SyslogAddress:                 "syslog-server",
				SyslogPort:                    10601,
				SyslogTransport:               "tcp",
				DopplerZone:                   "dopplerzone",
				DopplerMessageDrainBufferSize: 100,
				Secret:        config.Secret{},
				User:          config.User{},
				Certs:         &config.Certs{},
				InstanceCount: config.InstanceCount{},
				IP:            config.IP{},
			}
			config.DopplerSharedSecret = "secret"
			config.CCBulkAPIPassword = "bulk-pwd"
			config.EtcdMachines = []string{"1.0.0.7", "1.0.0.8"}
			config.DopplerIPs = []string{"1.0.11.1", "1.0.11.2"}
			config.DopplerVMType = "blah"

			dopplerPartition = NewDopplerPartition(config)
		})
		It("then it should allow the user to configure the doppler IPs", func() {
			ig := dopplerPartition.ToInstanceGroup()
			network := ig.Networks[0]
			Ω(len(network.StaticIPs)).Should(Equal(2))
			Ω(network.StaticIPs).Should(ConsistOf("1.0.11.1", "1.0.11.2"))
		})
		It("then it should have 2 instances", func() {
			ig := dopplerPartition.ToInstanceGroup()
			Ω(ig.Instances).Should(Equal(2))
		})
		It("then it should configure the correct number of instances automatically from the count of given IPs", func() {
			ig := dopplerPartition.ToInstanceGroup()
			network := ig.Networks[0]
			Ω(len(network.StaticIPs)).Should(Equal(ig.Instances))
		})

		It("then it should allow the user to configure the AZs", func() {
			ig := dopplerPartition.ToInstanceGroup()
			Ω(len(ig.AZs)).Should(Equal(1))
			Ω(ig.AZs[0]).Should(Equal("eastprod-1"))
		})

		It("then it should allow the user to configure vm-type", func() {
			ig := dopplerPartition.ToInstanceGroup()
			Ω(ig.VMType).ShouldNot(BeEmpty())
			Ω(ig.VMType).Should(Equal("blah"))
		})

		It("then it should allow the user to configure network to use", func() {
			ig := dopplerPartition.ToInstanceGroup()
			network := ig.GetNetworkByName("foundry-net")
			Ω(network).ShouldNot(BeNil())
		})

		It("then it should allow the user to configure the used stemcell", func() {
			ig := dopplerPartition.ToInstanceGroup()
			Ω(ig.Stemcell).ShouldNot(BeEmpty())
			Ω(ig.Stemcell).Should(Equal("cool-ubuntu-animal"))
		})
		It("then it should have update max in-flight 1 and serial false", func() {
			ig := dopplerPartition.ToInstanceGroup()
			Ω(ig.Update.MaxInFlight).Should(Equal(1))
			Ω(ig.Update.Serial).Should(Equal(false))
		})

		It("then it should then have 4 jobs", func() {
			ig := dopplerPartition.ToInstanceGroup()
			Ω(len(ig.Jobs)).Should(Equal(4))
		})
		It("then it should then have doppler job", func() {
			ig := dopplerPartition.ToInstanceGroup()
			job := ig.GetJobByName("doppler")
			Ω(job).ShouldNot(BeNil())
			props, _ := job.Properties.(*doppler.DopplerJob)
			Ω(props.Doppler.Zone).Should(Equal("dopplerzone"))
			Ω(props.Doppler.MessageDrainBufferSize).Should(Equal(100))
			Ω(props.Loggregator.Etcd.Machines).Should(ConsistOf("1.0.0.7", "1.0.0.8"))
			Ω(props.DopplerEndpoint.SharedSecret).Should(Equal("secret"))
		})
		It("then it should then have syslog_drain_binder job", func() {
			ig := dopplerPartition.ToInstanceGroup()
			job := ig.GetJobByName("syslog_drain_binder")
			Ω(job).ShouldNot(BeNil())
			props, _ := job.Properties.(*syslog_drain_binder.SyslogDrainBinderJob)
			Ω(props.Ssl.SkipCertVerify).Should(Equal(true))
			Ω(props.SystemDomain).Should(Equal("sys.test.com"))
			Ω(props.Cc.BulkApiPassword).Should(Equal("bulk-pwd"))
			Ω(props.Cc.SrvApiUri).Should(Equal("https://api.sys.test.com"))
			Ω(props.Loggregator.Etcd.Machines).Should(ConsistOf("1.0.0.7", "1.0.0.8"))

		})
	})
})
