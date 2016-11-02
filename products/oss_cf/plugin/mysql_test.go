package cloudfoundry_test

import (
	"github.com/enaml-ops/omg-product-bundle/products/oss_cf/enaml-gen/mysql"
	. "github.com/enaml-ops/omg-product-bundle/products/oss_cf/plugin"
	"github.com/enaml-ops/omg-product-bundle/products/oss_cf/plugin/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("MySQL Partition", func() {

	Context("when initialized WITH a complete set of arguments", func() {
		var mysqlPartition InstanceGroupCreator

		BeforeEach(func() {

			config := &config.Config{
				StemcellName:    "cool-ubuntu-animal",
				AZs:             []string{"eastprod-1"},
				NetworkName:     "foundry-net",
				SyslogAddress:   "syslog-server",
				SyslogPort:      10601,
				SyslogTransport: "tcp",
				Secret:          config.Secret{},
				User:            config.User{},
				Certs:           &config.Certs{},
				InstanceCount:   config.InstanceCount{},
				IP:              config.IP{},
			}
			config.UAADBPassword = "uaapassword"
			config.MySQLBootstrapUser = "mysqlbootstrap"
			config.MySQLBootstrapPassword = "mysqlbootstrappwd"
			config.MySQLIPs = []string{"1.0.10.1", "1.0.10.2"}
			config.MySQLVMType = "blah1"
			config.MySQLPersistentDiskType = "blah2"
			config.MySQLAdminPassword = "mysqladmin"

			mysqlPartition = NewMySQLPartition(config)
		})
		It("then it should allow the user to configure the mysql IPs", func() {
			ig := mysqlPartition.ToInstanceGroup()
			network := ig.Networks[0]
			Ω(len(network.StaticIPs)).Should(Equal(2))
			Ω(network.StaticIPs).Should(ConsistOf("1.0.10.1", "1.0.10.2"))
		})

		It("then it should have 2 instances", func() {
			ig := mysqlPartition.ToInstanceGroup()
			Ω(ig.Instances).Should(Equal(2))
		})

		It("then it should configure the correct number of instances automatically from the count of given IPs", func() {
			ig := mysqlPartition.ToInstanceGroup()
			network := ig.Networks[0]
			Ω(len(network.StaticIPs)).Should(Equal(ig.Instances))
		})

		It("then it should allow the user to configure the AZs", func() {
			ig := mysqlPartition.ToInstanceGroup()
			Ω(len(ig.AZs)).Should(Equal(1))
			Ω(ig.AZs[0]).Should(Equal("eastprod-1"))
		})

		It("then it should allow the user to configure vm-type", func() {
			ig := mysqlPartition.ToInstanceGroup()
			Ω(ig.VMType).ShouldNot(BeEmpty())
			Ω(ig.VMType).Should(Equal("blah1"))
		})

		It("then it should allow the user to configure the persistent disk", func() {
			ig := mysqlPartition.ToInstanceGroup()
			Ω(ig.PersistentDiskType).Should(Equal("blah2"))
		})

		It("then it should allow the user to configure network to use", func() {
			ig := mysqlPartition.ToInstanceGroup()
			network := ig.GetNetworkByName("foundry-net")
			Ω(network).ShouldNot(BeNil())
		})

		It("then it should allow the user to configure the used stemcell", func() {
			ig := mysqlPartition.ToInstanceGroup()
			Ω(ig.Stemcell).ShouldNot(BeEmpty())
			Ω(ig.Stemcell).Should(Equal("cool-ubuntu-animal"))
		})

		It("then it should have update max in-flight 1 and serial", func() {
			ig := mysqlPartition.ToInstanceGroup()
			Ω(ig.Update.MaxInFlight).Should(Equal(1))
			Ω(ig.Update.Serial).Should(Equal(false))
		})

		It("then it should then have 1 job", func() {
			ig := mysqlPartition.ToInstanceGroup()
			Ω(len(ig.Jobs)).Should(Equal(1))
		})

		It("then it should then have mysql job", func() {
			ig := mysqlPartition.ToInstanceGroup()
			job := ig.GetJobByName("mysql")
			Ω(job).ShouldNot(BeNil())
			props, _ := job.Properties.(*mysql.MysqlJob)
			Ω(props.CfMysql.Mysql.DatabaseStartupTimeout).Should(Equal(1200))
			Ω(props.CfMysql.Mysql.MaxConnections).Should(Equal(1500))
			Ω(props.CfMysql.Mysql.InnodbBufferPoolSize).Should(Equal(2147483648))
			Ω(props.CfMysql.Mysql.AdminUsername).Should(Equal("mysqlbootstrap"))
			Ω(props.CfMysql.Mysql.AdminPassword).Should(Equal("mysqlbootstrappwd"))
			Ω(props.SyslogAggregator.Address).Should(Equal("syslog-server"))
			Ω(props.SyslogAggregator.Port).Should(Equal(10601))
			Ω(props.SyslogAggregator.Transport).Should(Equal("tcp"))
			Ω(props.CfMysql.Mysql.ClusterIps).Should(ConsistOf("1.0.10.1", "1.0.10.2"))
			Ω(props.CfMysql.Mysql.SeededDatabases).ShouldNot(BeEmpty())
		})

		It("then the mysql job should have seeded databases", func() {
			mysql := mysqlPartition.(*MySQL)
			Ω(mysql.MySQLSeededDatabases).ShouldNot(BeEmpty())
			for _, db := range mysql.MySQLSeededDatabases {
				if db.Name == "uaa" {
					Ω(db.Password).Should(Equal("uaapassword"))
				}
			}
		})
	})
})
