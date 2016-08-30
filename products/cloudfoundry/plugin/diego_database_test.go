package cloudfoundry_test

import (
	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/bbs"
	"github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/consul_agent"
	"github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/etcd"
	. "github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/plugin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("given a Diego Database Partition", func() {
	Describe("given valid flags", func() {

		var instanceGroup *enaml.InstanceGroup
		var grouper InstanceGrouper

		Context("when ToInstanceGroup is called", func() {

			BeforeEach(func() {
				cf := new(Plugin)
				c := cf.GetContext([]string{
					"cloudfoundry",
					"--diego-db-ip", "10.0.0.39",
					"--diego-db-ip", "10.0.0.40",
					"--diego-db-vm-type", "dbvmtype",
					"--diego-db-disk-type", "dbdisktype",
					"--diego-db-passphrase", "random-db-encrytionkey",
					"--bbs-server-ca-cert", "cacert",
					"--etcd-server-cert", "blah-cert",
					"--etcd-server-key", "blah-key",
					"--etcd-client-cert", "bleh-cert",
					"--etcd-client-key", "bleh-key",
					"--etcd-peer-cert", "blee-cert",
					"--etcd-peer-key", "blee-key",
					"--bbs-client-cert", "clientcert",
					"--bbs-client-key", "clientkey",
					"--bbs-server-cert", "clientcert",
					"--bbs-server-key", "clientkey",
					"--consul-ip", "1.0.0.1",
					"--consul-ip", "1.0.0.2",
					"--consul-vm-type", "blah",
					"--metron-secret", "metronsecret",
					"--metron-zone", "metronzoneguid",
					"--syslog-address", "syslog-server",
					"--syslog-port", "10601",
					"--syslog-transport", "tcp",
					"--etcd-machine-ip", "1.0.0.7",
					"--etcd-machine-ip", "1.0.0.8",
				})
				config := &Config{
					SystemDomain:      "service.cf.domain.com",
					AZs:               []string{"eastprod-1"},
					StemcellName:      "cool-ubuntu-animal",
					NetworkName:       "foundry-net",
					AllowSSHAccess:    true,
					ConsulEncryptKeys: []string{"encyption-key"},
					ConsulCaCert:      "ca-cert",
					ConsulAgentCert:   "agent-cert",
					ConsulAgentKey:    "agent-key",
					ConsulServerCert:  "server-cert",
					ConsulServerKey:   "server-key",
					BBSCACert:         "cacert",
					BBSServerCert:     "clientcert",
					BBSServerKey:      "clientkey",
				}
				grouper = NewDiegoDatabasePartition(c, config)
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
						Ω(services).Should(HaveKey("bbs"))
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
