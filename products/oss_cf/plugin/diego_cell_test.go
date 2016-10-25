package cloudfoundry_test

import (
	"encoding/json"
	"fmt"

	yaml "gopkg.in/yaml.v2"

	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-product-bundle/products/oss_cf/enaml-gen/consul_agent"
	"github.com/enaml-ops/omg-product-bundle/products/oss_cf/enaml-gen/garden"
	"github.com/enaml-ops/omg-product-bundle/products/oss_cf/enaml-gen/rep"
	. "github.com/enaml-ops/omg-product-bundle/products/oss_cf/plugin"
	"github.com/enaml-ops/omg-product-bundle/products/oss_cf/plugin/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("given a Diego Cell Partition", func() {
	Describe("given valid flags", func() {

		var instanceGroup *enaml.InstanceGroup
		var grouper InstanceGroupCreator

		Context("when ToInstanceGroup is called", func() {

			BeforeEach(func() {
				config := &config.Config{
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
				config.SkipSSLCertVerify = true
				config.ConsulEncryptKeys = []string{"encyption-key"}
				config.ConsulAgentCert = "agent-cert"
				config.ConsulAgentKey = "agent-key"
				config.ConsulServerCert = "server-cert"
				config.ConsulServerKey = "server-key"
				config.DiegoCellIPs = []string{"10.0.0.39", "10.0.0.40"}
				config.DiegoCellVMType = "cellvmtype"
				config.DiegoCellPersistentDiskType = "celldisktype"
				config.BBSCACert = "cacert"
				config.BBSClientCert = "clientcert"
				config.BBSClientKey = "clientkey"
				config.DopplerSharedSecret = "metronsecret"
				config.EtcdMachines = []string{"1.0.0.7", "1.0.0.8"}
				grouper = NewDiegoCellPartition(config)
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
				Ω(instanceGroup.GetJobByName("rep")).ShouldNot(BeNil())
				Ω(instanceGroup.GetJobByName("consul_agent")).ShouldNot(BeNil())
				Ω(instanceGroup.GetJobByName("cflinuxfs2-rootfs-setup")).ShouldNot(BeNil())
				Ω(instanceGroup.GetJobByName("garden")).ShouldNot(BeNil())
				Ω(instanceGroup.GetJobByName("statsd-injector")).ShouldNot(BeNil())
				Ω(instanceGroup.GetJobByName("metron_agent")).ShouldNot(BeNil())
			})

			Describe("given a rep job", func() {
				Context("when defined", func() {
					var job *enaml.InstanceJob
					BeforeEach(func() {
						job = instanceGroup.GetJobByName("rep")
					})
					It("then it should use the correct release", func() {
						Ω(job.Release).Should(Equal(DiegoReleaseName))
					})
					It("then it should set the BBS API location", func() {
						props := job.Properties.(*rep.RepJob)
						Ω(props.Diego.Rep.Bbs.ApiLocation).Should(Equal("bbs.service.cf.internal:8889"))
					})
					It("then it should set SSL", func() {
						props := job.Properties.(*rep.RepJob)
						Ω(props.Diego.Ssl).ShouldNot(BeNil())
						Ω(props.Diego.Ssl.SkipCertVerify).Should(BeTrue())
					})
					It("should make rootfses an array", func() {
						props := job.Properties.(*rep.RepJob)
						Ω(props.Diego.Rep.PreloadedRootfses).Should(HaveLen(1))
						Ω(props.Diego.Rep.PreloadedRootfses).Should(ConsistOf("cflinuxfs2:/var/vcap/packages/cflinuxfs2/rootfs"))
					})
				})
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
						Ω(props.Consul.Agent.Services).Should(BeEmpty())
					})
				})
			})

			Describe("given a cflinuxfs2-rootfs-setup job", func() {
				Context("when defined", func() {
					var job *enaml.InstanceJob

					BeforeEach(func() {
						job = instanceGroup.GetJobByName("cflinuxfs2-rootfs-setup")
					})

					It("then it should use the correct release", func() {
						Ω(job.Release).Should(Equal(CFLinuxReleaseName))
					})

					It("then it should not generate null properties", func() {
						Ω(job.Properties).ShouldNot(BeNil())

						bytes, err := yaml.Marshal(job.Properties)
						Ω(err).ShouldNot(HaveOccurred())
						Ω(bytes).Should(MatchYAML("{}"))
					})
				})
			})

			Describe("given a garden job", func() {
				Context("when defined", func() {
					var job *enaml.InstanceJob
					BeforeEach(func() {
						job = instanceGroup.GetJobByName("garden")
					})
					It("then it should use the correct release", func() {
						Ω(job.Release).Should(Equal(GardenReleaseName))
					})

					It("then it should populate my properties", func() {
						Ω(job.Properties).ShouldNot(BeNil())
						gardenJob := job.Properties.(*garden.GardenJob)
						Ω(gardenJob.Garden.AllowHostAccess).Should(BeFalse())
						Ω(gardenJob.Garden.PersistentImageList).Should(ConsistOf("/var/vcap/packages/cflinuxfs2/rootfs"))
						Ω(gardenJob.Garden.DenyNetworks).Should(ConsistOf("0.0.0.0/0"))
						Ω(gardenJob.Garden.NetworkPool).Should(Equal("10.254.0.0/22"))
						Ω(gardenJob.Garden.NetworkMtu).Should(Equal(1454))
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

			Describe("given a rep job", func() {
				Context("when defined", func() {
					var job *enaml.InstanceJob
					BeforeEach(func() {
						job = instanceGroup.GetJobByName("rep")
					})
					It("then it should use the correct release", func() {
						Ω(job.Release).Should(Equal(DiegoReleaseName))
					})
					It("then it should populate my properties", func() {
						b, _ := json.Marshal(job.Properties)
						fmt.Println("job", string(b))
						Ω(job.Properties).ShouldNot(BeNil())
					})
				})
			})
		})
	})
})
