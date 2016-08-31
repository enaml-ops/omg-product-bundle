package cloudfoundry_test

import (
	"encoding/json"
	"fmt"

	yaml "gopkg.in/yaml.v2"

	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/consul_agent"
	"github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/rep"
	. "github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/plugin"
	"github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/plugin/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("given a Diego Cell Partition", func() {
	Describe("given valid flags", func() {

		var instanceGroup *enaml.InstanceGroup
		var grouper InstanceGroupCreator

		Context("when ToInstanceGroup is called", func() {

			BeforeEach(func() {
				/*cf := new(Plugin)
				c := cf.GetContext([]string{
					"cloudfoundry",
					"--diego-cell-ip", "10.0.0.39",
					"--diego-cell-ip", "10.0.0.40",
					"--diego-cell-vm-type", "cellvmtype",
					"--diego-cell-disk-type", "celldisktype",
					"--bbs-server-ca-cert", "cacert",
					"--bbs-client-cert", "clientcert",
					"--bbs-client-key", "clientkey",
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
				})*/
				config := &config.Config{
					AZs:                         []string{"eastprod-1"},
					StemcellName:                "cool-ubuntu-animal",
					NetworkName:                 "foundry-net",
					AllowSSHAccess:              true,
					ConsulEncryptKeys:           []string{"encyption-key"},
					ConsulCaCert:                "ca-cert",
					ConsulAgentCert:             "agent-cert",
					ConsulAgentKey:              "agent-key",
					ConsulServerCert:            "server-cert",
					ConsulServerKey:             "server-key",
					DiegoCellIPs:                []string{"10.0.0.39", "10.0.0.40"},
					DiegoCellVMType:             "cellvmtype",
					DiegoCellPersistentDiskType: "celldisktype",
					BBSCACert:                   "cacert",
					BBSClientCert:               "clientcert",
					BBSClientKey:                "clientkey",
					MetronSecret:                "metronsecret",
					MetronZone:                  "metronzoneguid",
					SyslogAddress:               "syslog-server",
					SyslogPort:                  10601,
					SyslogTransport:             "tcp",
					EtcdMachines:                []string{"1.0.0.7", "1.0.0.8"},
					/*"--consul-ip", "1.0.0.1",
					"--consul-ip", "1.0.0.2",
					"--consul-vm-type", "blah",
					*/
				}
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
						Ω(job.Release).Should(Equal(CFLinuxFSReleaseName))
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
