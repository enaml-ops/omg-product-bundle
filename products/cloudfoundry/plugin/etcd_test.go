package cloudfoundry_test

import (
	"github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/etcd"
	"github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/etcd_metrics_server"
	. "github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/plugin"
	"github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/plugin/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Etcd Partition", func() {
	Context("when initialized WITH a complete set of arguments", func() {
		var etcdPartition InstanceGroupCreator
		BeforeEach(func() {
			config := &config.Config{
				StemcellName:    "cool-ubuntu-animal",
				AZs:             []string{"eastprod-1"},
				NetworkName:     "foundry-net",
				NATSPort:        4222,
				MetronZone:      "metronzoneguid",
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
			config.EtcdVMType = "blah"
			config.EtcdPersistentDiskType = "blah-disk"
			config.MetronSecret = "metronsecret"
			config.NATSUser = "nats"
			config.NATSPassword = "pass"
			config.NATSMachines = []string{"1.0.0.5", "1.0.0.6"}

			etcdPartition = NewEtcdPartition(config)
		})
		It("then it should allow the user to configure the etcd IPs", func() {
			ig := etcdPartition.ToInstanceGroup()
			network := ig.Networks[0]
			Ω(len(network.StaticIPs)).Should(Equal(2))
			Ω(network.StaticIPs).Should(ConsistOf("1.0.0.7", "1.0.0.8"))
		})
		It("then it should have 2 instances", func() {
			ig := etcdPartition.ToInstanceGroup()
			Ω(ig.Instances).Should(Equal(2))
		})
		It("then it should configure the correct number of instances automatically from the count of given IPs", func() {
			ig := etcdPartition.ToInstanceGroup()
			network := ig.Networks[0]
			Ω(len(network.StaticIPs)).Should(Equal(ig.Instances))
		})

		It("then it should allow the user to configure the AZs", func() {
			ig := etcdPartition.ToInstanceGroup()
			Ω(len(ig.AZs)).Should(Equal(1))
			Ω(ig.AZs[0]).Should(Equal("eastprod-1"))
		})

		It("then it should allow the user to configure vm-type", func() {
			ig := etcdPartition.ToInstanceGroup()
			Ω(ig.VMType).ShouldNot(BeEmpty())
			Ω(ig.VMType).Should(Equal("blah"))
		})

		It("then it should allow the user to configure network to use", func() {
			ig := etcdPartition.ToInstanceGroup()
			network := ig.GetNetworkByName("foundry-net")
			Ω(network).ShouldNot(BeNil())
		})

		It("then it should allow the user to configure the used stemcell", func() {
			ig := etcdPartition.ToInstanceGroup()
			Ω(ig.Stemcell).ShouldNot(BeEmpty())
			Ω(ig.Stemcell).Should(Equal("cool-ubuntu-animal"))
		})

		It("then it should allow the user to configure disk type to use", func() {
			ig := etcdPartition.ToInstanceGroup()
			Ω(ig.PersistentDiskType).Should(Equal("blah-disk"))
		})
		It("then it should have update max in-flight 1 and serial", func() {
			ig := etcdPartition.ToInstanceGroup()
			Ω(ig.Update.MaxInFlight).Should(Equal(1))
			Ω(ig.Update.Serial).Should(Equal(false))
		})
		It("then it should then have 4 jobs", func() {
			ig := etcdPartition.ToInstanceGroup()
			Ω(len(ig.Jobs)).Should(Equal(4))
		})
		It("then it should then have etcd job", func() {
			ig := etcdPartition.ToInstanceGroup()
			job := ig.GetJobByName("etcd")
			Ω(job).ShouldNot(BeNil())
			props, _ := job.Properties.(*etcd.EtcdJob)
			Ω(props.Etcd.RequireSsl).Should(Equal(false))
			Ω(props.Etcd.PeerRequireSsl).Should(Equal(false))
			Ω(props.Etcd.Machines).Should(ConsistOf("1.0.0.7", "1.0.0.8"))
		})
		It("then it should then have etcd_metrics_server job", func() {
			ig := etcdPartition.ToInstanceGroup()
			job := ig.GetJobByName("etcd_metrics_server")
			Ω(job).ShouldNot(BeNil())
			props, _ := job.Properties.(*etcd_metrics_server.EtcdMetricsServerJob)
			Ω(props.EtcdMetricsServer.Nats.Machines).Should(ConsistOf("1.0.0.5", "1.0.0.6"))
			Ω(props.EtcdMetricsServer.Nats.Username).Should(Equal("nats"))
			Ω(props.EtcdMetricsServer.Nats.Password).Should(Equal("pass"))
		})
	})
})
