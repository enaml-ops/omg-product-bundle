package cloudfoundry_test

import (
	"github.com/enaml-ops/omg-product-bundle/products/cf-mysql/enaml-gen/proxy"
	. "github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/plugin"
	"github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/plugin/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("MySQL Proxy Partition", func() {
	Context("when initialized WITH a complete set of arguments", func() {
		var mysqlProxyPartition InstanceGroupCreator
		BeforeEach(func() {
			config := &config.Config{
				StemcellName:           "cool-ubuntu-animal",
				AZs:                    []string{"eastprod-1"},
				NetworkName:            "foundry-net",
				NATSPort:               4222,
				MySQLProxyExternalHost: "mysqlhostname",
				SyslogAddress:          "syslog-server",
				SyslogPort:             10601,
				SyslogTransport:        "tcp",
				Secret:                 config.Secret{},
				User:                   config.User{},
				Certs:                  &config.Certs{},
				InstanceCount:          config.InstanceCount{},
				IP:                     config.IP{},
			}
			config.NATSUser = "nats"
			config.NATSPassword = "pass"
			config.NATSMachines = []string{"1.0.0.5", "1.0.0.6"}
			config.MySQLIPs = []string{"1.0.10.1", "1.0.10.2"}
			config.MySQLProxyIPs = []string{"1.0.10.3", "1.0.10.4"}
			config.MySQLProxyVMType = "blah"
			config.MySQLProxyAPIUsername = "apiuser"
			config.MySQLProxyAPIPassword = "apipassword"

			mysqlProxyPartition = NewMySQLProxyPartition(config)
		})
		It("then it should allow the user to configure the mysql proxy IPs", func() {
			ig := mysqlProxyPartition.ToInstanceGroup()
			Ω(len(ig.Networks)).Should(Equal(1))
			network := ig.Networks[0]
			Ω(len(network.StaticIPs)).Should(Equal(2))
			Ω(network.StaticIPs).Should(ConsistOf("1.0.10.3", "1.0.10.4"))
		})
		It("then it should have 2 instances", func() {
			ig := mysqlProxyPartition.ToInstanceGroup()
			Ω(ig.Instances).Should(Equal(2))
		})
		It("then it should configure the correct number of instances automatically from the count of given IPs", func() {
			ig := mysqlProxyPartition.ToInstanceGroup()
			network := ig.Networks[0]
			Ω(len(network.StaticIPs)).Should(Equal(ig.Instances))
		})

		It("then it should allow the user to configure the AZs", func() {
			ig := mysqlProxyPartition.ToInstanceGroup()
			Ω(len(ig.AZs)).Should(Equal(1))
			Ω(ig.AZs[0]).Should(Equal("eastprod-1"))
		})

		It("then it should allow the user to configure vm-type", func() {
			ig := mysqlProxyPartition.ToInstanceGroup()
			Ω(ig.VMType).ShouldNot(BeEmpty())
			Ω(ig.VMType).Should(Equal("blah"))
		})

		It("then it should allow the user to configure network to use", func() {
			ig := mysqlProxyPartition.ToInstanceGroup()
			network := ig.GetNetworkByName("foundry-net")
			Ω(network).ShouldNot(BeNil())
		})

		It("then it should allow the user to configure the used stemcell", func() {
			ig := mysqlProxyPartition.ToInstanceGroup()
			Ω(ig.Stemcell).ShouldNot(BeEmpty())
			Ω(ig.Stemcell).Should(Equal("cool-ubuntu-animal"))
		})
		It("then it should have update max in-flight 1 and serial", func() {
			ig := mysqlProxyPartition.ToInstanceGroup()
			Ω(ig.Update.MaxInFlight).Should(Equal(1))
			Ω(ig.Update.Serial).Should(Equal(false))
		})

		It("then it should then have 1 job", func() {
			ig := mysqlProxyPartition.ToInstanceGroup()
			Ω(len(ig.Jobs)).Should(Equal(1))
		})
		It("then it should then have proxy job", func() {
			ig := mysqlProxyPartition.ToInstanceGroup()
			job := ig.GetJobByName("proxy")
			Ω(job).ShouldNot(BeNil())
			Ω(job.Release).Should(Equal("cf-mysql"))
			props, _ := job.Properties.(*proxy.ProxyJob)
			Ω(props.Proxy.ApiUsername).Should(Equal("apiuser"))
			Ω(props.Proxy.ApiPassword).Should(Equal("apipassword"))
			Ω(props.Proxy.ProxyIps).Should(ConsistOf("1.0.10.3", "1.0.10.4"))
			Ω(props.ExternalHost).Should(Equal("mysqlhostname"))
			Ω(props.ClusterIps).Should(ConsistOf("1.0.10.1", "1.0.10.2"))
			Ω(props.SyslogAggregator.Address).Should(Equal("syslog-server"))
			Ω(props.SyslogAggregator.Port).Should(Equal(10601))
			Ω(props.SyslogAggregator.Transport).Should(Equal("tcp"))
			Ω(props.Nats).ShouldNot(BeNil())
			Ω(props.Nats.Port).Should(Equal(4222))
			Ω(props.Nats.User).Should(Equal("nats"))
			Ω(props.Nats.Password).Should(Equal("pass"))
			Ω(props.Nats.Machines).Should(ConsistOf("1.0.0.5", "1.0.0.6"))
		})

	})
})
