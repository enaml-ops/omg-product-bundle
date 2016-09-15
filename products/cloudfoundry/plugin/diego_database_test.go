package cloudfoundry_test

import (
	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/bbs"
	"github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/consul_agent"
	"github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/etcd"
	. "github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/plugin"
	"github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/plugin/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("given a Diego Database Partition", func() {
	Describe("given valid flags", func() {

		var instanceGroup *enaml.InstanceGroup
		var grouper InstanceGroupCreator

		Context("when ToInstanceGroup is called", func() {

			BeforeEach(func() {

				config := &config.Config{
					SystemDomain:    "service.cf.domain.com",
					AZs:             []string{"eastprod-1"},
					StemcellName:    "cool-ubuntu-animal",
					NetworkName:     "foundry-net",
					AllowSSHAccess:  true,
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
				config.EtcdMachines = []string{"1.0.0.7", "1.0.0.8"}
				config.ConsulEncryptKeys = []string{"encyption-key"}
				config.ConsulAgentCert = "agent-cert"
				config.ConsulAgentKey = "agent-key"
				config.ConsulServerCert = "server-cert"
				config.ConsulServerKey = "server-key"
				config.BBSCACert = "cacert"
				config.BBSServerCert = "clientcert"
				config.BBSServerKey = "clientkey"
				config.DiegoDBIPs = []string{"10.0.0.39", "10.0.0.40"}
				config.DiegoDBVMType = "dbvmtype"
				config.DiegoDBPersistentDiskType = "dbdisktype"
				config.DiegoDBPassphrase = "random-db-encrytionkey"
				config.EtcdServerCert = "blah-cert"
				config.EtcdServerKey = "blah-key"
				config.EtcdClientCert = "bleh-cert"
				config.EtcdClientKey = "bleh-key"
				config.EtcdPeerCert = "blee-cert"
				config.EtcdPeerKey = "blee-key"
				config.BBSClientCert = "clientcert"
				config.BBSClientKey = "clientkey"
				config.ConsulIPs = []string{"1.0.0.1", "1.0.0.2"}
				config.DopplerSharedSecret = "metronsecret"

				grouper = NewDiegoDatabasePartition(config)
				instanceGroup = grouper.ToInstanceGroup()
			})

			It("then it should be populated with valid network configs", func() {
				ignet := instanceGroup.GetNetworkByName("foundry-net")
				Ω(ignet).ShouldNot(BeNil())
				Ω(ignet.StaticIPs).Should(ConsistOf("10.0.0.39", "10.0.0.40"))
			})

			It("then it should have an instance count in line with given IPs", func() {
				ignet := instanceGroup.GetNetworkByName("foundry-net")
				Ω(len(ignet.StaticIPs)).Should(Equal(instanceGroup.Instances))
			})

			It("then it should be populated the required jobs", func() {
				Ω(instanceGroup.GetJobByName("etcd")).ShouldNot(BeNil())
				Ω(instanceGroup.GetJobByName("bbs")).ShouldNot(BeNil())
				Ω(instanceGroup.GetJobByName("consul_agent")).ShouldNot(BeNil())
				Ω(instanceGroup.GetJobByName("metron_agent")).ShouldNot(BeNil())
				Ω(instanceGroup.GetJobByName("statsd-injector")).ShouldNot(BeNil())
			})
			Describe("given a consul_agent job", func() {
				Context("when defined", func() {
					var job *enaml.InstanceJob
					BeforeEach(func() {
						job = instanceGroup.GetJobByName("consul_agent")
					})
					It("then it should use the correct release", func() {
						Ω(job.Release).Should(Equal(CFReleaseName))
					})

					It("then it should populate my properties", func() {
						Ω(job.Properties).ShouldNot(BeNil())
						props := job.Properties.(*consul_agent.ConsulAgentJob)
						Ω(props.Consul.Agent.Mode).Should(BeNil())
						services := props.Consul.Agent.Services.(map[string]map[string]string)
						Ω(services).Should(HaveKey("etcd"))
					})
				})
			})

			Describe("given a statsd-injector job", func() {
				Context("when defined", func() {
					var job *enaml.InstanceJob
					BeforeEach(func() {
						job = instanceGroup.GetJobByName("statsd-injector")
					})
					It("then it should use the correct release", func() {
						Ω(job.Release).Should(Equal(CFReleaseName))
					})

					It("then it should populate my properties", func() {
						Ω(job.Properties).ShouldNot(BeNil())
					})
				})
			})

			Describe("given a metron_agent job", func() {
				Context("when defined", func() {
					var job *enaml.InstanceJob
					BeforeEach(func() {
						job = instanceGroup.GetJobByName("metron_agent")
					})
					It("then it should use the correct release", func() {
						Ω(job.Release).Should(Equal(CFReleaseName))
					})
					It("then it should populate my properties", func() {
						Ω(job.Properties).ShouldNot(BeNil())
					})
				})
			})

			Describe("given a etcd job", func() {
				Context("when defined", func() {
					var job *enaml.InstanceJob
					BeforeEach(func() {
						job = instanceGroup.GetJobByName("etcd")
					})
					It("then it should use the correct release", func() {
						Ω(job.Release).Should(Equal(EtcdReleaseName))
					})
					It("then it should populate my properties under etcd", func() {
						Ω(job.Properties).ShouldNot(BeNil())
						props := job.Properties.(*etcd.EtcdJob)
						Ω(props.Etcd).ShouldNot(BeNil())
					})
					It("then it should have the cluster as an array", func() {
						props := job.Properties.(*etcd.EtcdJob)
						arr := props.Etcd.Cluster.([]map[string]interface{})
						Ω(arr).Should(HaveLen(1))
					})
					It("then it should use internal hostnames for etcd", func() {
						props := job.Properties.(*etcd.EtcdJob)
						Ω(props.Etcd.AdvertiseUrlsDnsSuffix).Should(Equal("etcd.service.cf.internal"))
						Ω(props.Etcd.Machines).Should(ConsistOf("etcd.service.cf.internal"))
					})
				})
			})

			Describe("given a bbs job", func() {
				Context("when defined", func() {
					var job *enaml.InstanceJob

					BeforeEach(func() {
						job = instanceGroup.GetJobByName("bbs")
					})

					It("then it should use the correct release", func() {
						Ω(job.Release).Should(Equal(DiegoReleaseName))
					})

					It("then it should populate my properties", func() {
						Ω(job.Properties).ShouldNot(BeNil())
					})

					It("should properly set my server key/cert", func() {
						props := job.Properties.(*bbs.BbsJob)
						Ω(props.Diego.Bbs.ServerCert).Should(Equal("clientcert"))
						Ω(props.Diego.Bbs.ServerKey).Should(Equal("clientkey"))
					})

					It("should properly set my db passphrase", func() {
						props := job.Properties.(*bbs.BbsJob)
						arr := props.Diego.Bbs.EncryptionKeys.([]map[string]string)
						Ω(arr).Should(HaveLen(1))
						Ω(arr[0]["passphrase"]).Should(Equal("random-db-encrytionkey"))
					})

					It("should properly set my auctioneer API URL", func() {
						props := job.Properties.(*bbs.BbsJob)
						Ω(props.Diego.Bbs.Auctioneer.ApiUrl).Should(Equal("http://auctioneer.service.cf.internal:9016"))
					})

					It("should properly set my bbs.etcd", func() {
						props := job.Properties.(*bbs.BbsJob)
						Ω(props.Diego.Bbs.Etcd).ShouldNot(BeNil())
						Ω(props.Diego.Bbs.Etcd.Machines).Should(ConsistOf("etcd.service.cf.internal"))
					})

					It("should have encryption keys as an array", func() {
						props := job.Properties.(*bbs.BbsJob)
						arr := props.Diego.Bbs.EncryptionKeys.([]map[string]string)
						Ω(arr).Should(HaveLen(1))
					})
				})
			})
		})
	})
})
