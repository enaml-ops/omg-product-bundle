package cloudfoundry_test

import (
	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/auctioneer"
	"github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/cc_uploader"
	"github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/consul_agent"
	"github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/converger"
	"github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/file_server"
	"github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/metron_agent"
	"github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/nsync"
	"github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/route_emitter"
	"github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/ssh_proxy"
	"github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/stager"
	"github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/tps"
	. "github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/plugin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("given a Diego Brain Partition", func() {
	Context("when initialized WITHOUT a complete set of arguments", func() {
		var ig InstanceGrouper
		BeforeEach(func() {
			cf := new(Plugin)
			c := cf.GetContext([]string{
				"cloudfoundry",
				"--metron-secret", "metronsecret",
				"--metron-zone", "metronzoneguid",
			})
			ig = NewDiegoBrainPartition(c)
		})

		It("then it should not validate", func() {
			Ω(ig.HasValidValues()).Should(BeFalse())
		})
	})

	Context("when initialized with a complete set of arguments", func() {
		var deploymentManifest *enaml.DeploymentManifest
		var grouper InstanceGrouper

		BeforeEach(func() {
			cf := new(Plugin)
			c := cf.GetContext([]string{
				"cloudfoundry",
				"--az", "eastprod-1",
				"--stemcell-name", "cool-ubuntu-animal",
				"--network", "foundry-net",
				"--allow-app-ssh-access",
				"--diego-brain-ip", "10.0.0.39",
				"--diego-brain-ip", "10.0.0.40",
				"--diego-brain-vm-type", "brainvmtype",
				"--diego-brain-disk-type", "braindisktype",
				"--bbs-server-ca-cert", "cacert",
				"--bbs-client-cert", "clientcert",
				"--bbs-client-key", "clientkey",
				"--bbs-require-ssl=false",
				"--skip-cert-verify=false",
				"--cc-uploader-poll-interval", "25",
				"--cc-external-port", "9023",
				"--system-domain", "sys.test.com",
				"--cc-internal-api-user", "internaluser",
				"--cc-internal-api-password", "internalpassword",
				"--cc-bulk-batch-size", "5",
				"--cc-fetch-timeout", "30",
				"--fs-listen-addr", "0.0.0.0:12345",
				"--fs-static-dir", "/foo/bar/baz",
				"--fs-debug-addr", "10.0.1.2:22222",
				"--fs-log-level", "debug",
				"--metron-port", "3458",
				"--nats-user", "nats",
				"--nats-port", "1234",
				"--nats-pass", "natspass",
				"--nats-machine-ip", "10.0.0.11",
				"--nats-machine-ip", "10.0.0.12",
				"--ssh-proxy-uaa-secret", "secret",
				"--traffic-controller-url", "wss://doppler.sys.yourdomain.com:443",
				"--consul-vm-type", "blah",
				"--consul-encryption-key", "encyption-key",
				"--consul-server-ca-cert", "ca-cert",
				"--consul-agent-cert", "agent-cert",
				"--consul-agent-key", "agent-key",
				"--consul-server-cert", "server-cert",
				"--consul-server-key", "server-key",
				"--consul-ip", "1.0.0.1",
				"--consul-ip", "1.0.0.2",
				"--metron-secret", "metronsecret",
				"--metron-zone", "metronzoneguid",
				"--etcd-machine-ip", "1.0.0.7",
				"--etcd-machine-ip", "1.0.0.8",
			})
			grouper = NewDiegoBrainPartition(c)
			deploymentManifest = new(enaml.DeploymentManifest)
			deploymentManifest.AddInstanceGroup(grouper.ToInstanceGroup())
		})

		It("then it should configure the instance group correctly", func() {

			By("having valid values")
			Ω(grouper.HasValidValues()).Should(BeTrue())

			ig := deploymentManifest.GetInstanceGroupByName("diego_brain-partition")

			By("configuring the AZs")
			Ω(len(ig.AZs)).Should(Equal(1))
			Ω(ig.AZs[0]).Should(Equal("eastprod-1"))

			By("configuring the stemcell")
			Ω(ig.Stemcell).Should(Equal("cool-ubuntu-animal"))

			By("configuring the network")
			network := ig.GetNetworkByName("foundry-net")
			Ω(network).ShouldNot(BeNil())
			Ω(len(network.StaticIPs)).Should(Equal(2))
			Ω(network.StaticIPs[0]).Should(Equal("10.0.0.39"))
			Ω(network.StaticIPs[1]).Should(Equal("10.0.0.40"))

			By("configuring the VM type")
			Ω(ig.VMType).Should(Equal("brainvmtype"))

			By("configuring the disk type")
			Ω(ig.PersistentDiskType).Should(Equal("braindisktype"))

			By("setting the correct number of instances")
			Ω(len(ig.Networks)).Should(Equal(1))
			Ω(len(ig.Networks[0].StaticIPs)).Should(Equal(ig.Instances))

			By("configuring update")
			Ω(ig.Update.MaxInFlight).Should(Equal(1))

			By("configuring the auctioneer job")
			job := ig.GetJobByName("auctioneer")
			Ω(job.Release).Should(Equal(DiegoReleaseName))
			a := job.Properties.(*auctioneer.AuctioneerJob)
			Ω(a.Diego.Auctioneer.Bbs.CaCert).Should(Equal("cacert"))
			Ω(a.Diego.Auctioneer.Bbs.ClientCert).Should(Equal("clientcert"))
			Ω(a.Diego.Auctioneer.Bbs.ClientKey).Should(Equal("clientkey"))
			Ω(a.Diego.Auctioneer.Bbs.ApiLocation).Should(Equal("bbs.service.cf.internal:8889"))

			By("configuring the CC uploader")
			job = ig.GetJobByName("cc_uploader")
			Ω(job.Release).Should(Equal(DiegoReleaseName))
			cc := job.Properties.(*cc_uploader.CcUploaderJob)
			Ω(cc.Diego.Ssl.SkipCertVerify).Should(BeFalse())
			Ω(cc.Diego.CcUploader.Cc.JobPollingIntervalInSeconds).Should(Equal(25))

			By("configuring the converger")
			job = ig.GetJobByName("converger")
			Ω(job.Release).Should(Equal(DiegoReleaseName))
			c := job.Properties.(*converger.Converger)
			Ω(c.Bbs.ApiLocation).Should(Equal("bbs.service.cf.internal:8889"))
			Ω(c.Bbs.CaCert).Should(Equal("cacert"))
			Ω(c.Bbs.ClientCert).Should(Equal("clientcert"))
			Ω(c.Bbs.ClientKey).Should(Equal("clientkey"))

			By("configuring the file server")
			job = ig.GetJobByName("file_server")
			Ω(job.Release).Should(Equal(DiegoReleaseName))
			fs := job.Properties.(*file_server.FileServerJob)
			Ω(fs.Diego.Ssl.SkipCertVerify).Should(BeFalse())

			By("configuring nsync")
			job = ig.GetJobByName("nsync")
			Ω(job.Release).Should(Equal(DiegoReleaseName))
			n := job.Properties.(*nsync.NsyncJob)
			Ω(n.Diego.Ssl.SkipCertVerify).Should(BeFalse())
			Ω(n.Diego.Nsync.Bbs.ApiLocation).Should(Equal("bbs.service.cf.internal:8889"))
			Ω(n.Diego.Nsync.Bbs.CaCert).Should(Equal("cacert"))
			Ω(n.Diego.Nsync.Bbs.ClientCert).Should(Equal("clientcert"))
			Ω(n.Diego.Nsync.Bbs.ClientKey).Should(Equal("clientkey"))
			Ω(n.Diego.Nsync.Cc.BaseUrl).Should(Equal("https://api.sys.test.com"))
			Ω(n.Diego.Nsync.Cc.BasicAuthUsername).Should(Equal("internaluser"))
			Ω(n.Diego.Nsync.Cc.BasicAuthPassword).Should(Equal("internalpassword"))
			Ω(n.Diego.Nsync.Cc.BulkBatchSize).Should(Equal(5))
			Ω(n.Diego.Nsync.Cc.FetchTimeoutInSeconds).Should(Equal(30))
			Ω(n.Diego.Nsync.Cc.PollingIntervalInSeconds).Should(Equal(25))

			By("configuring the route emitter")
			job = ig.GetJobByName("route_emitter")
			Ω(job.Release).Should(Equal(DiegoReleaseName))
			r := job.Properties.(*route_emitter.RouteEmitterJob)
			Ω(r.Diego.RouteEmitter.Bbs.ApiLocation).Should(Equal("bbs.service.cf.internal:8889"))
			Ω(r.Diego.RouteEmitter.Bbs.CaCert).Should(Equal("cacert"))
			Ω(r.Diego.RouteEmitter.Bbs.ClientCert).Should(Equal("clientcert"))
			Ω(r.Diego.RouteEmitter.Bbs.ClientKey).Should(Equal("clientkey"))
			Ω(r.Diego.RouteEmitter.Bbs.RequireSsl).Should(BeFalse())
			Ω(r.Diego.RouteEmitter.Nats.User).Should(Equal("nats"))
			Ω(r.Diego.RouteEmitter.Nats.Password).Should(Equal("natspass"))
			Ω(r.Diego.RouteEmitter.Nats.Port).Should(Equal(1234))
			Ω(r.Diego.RouteEmitter.Nats.Machines).Should(ContainElement("10.0.0.11"))
			Ω(r.Diego.RouteEmitter.Nats.Machines).Should(ContainElement("10.0.0.12"))

			By("configuring the SSH proxy")
			job = ig.GetJobByName("ssh_proxy")
			s := job.Properties.(*ssh_proxy.SshProxyJob)
			Ω(s.Diego.Ssl.SkipCertVerify).Should(BeFalse())
			Ω(s.Diego.SshProxy.Bbs.ApiLocation).Should(Equal("bbs.service.cf.internal:8889"))
			Ω(s.Diego.SshProxy.Bbs.CaCert).Should(Equal("cacert"))
			Ω(s.Diego.SshProxy.Bbs.ClientCert).Should(Equal("clientcert"))
			Ω(s.Diego.SshProxy.Bbs.ClientKey).Should(Equal("clientkey"))
			Ω(s.Diego.SshProxy.Bbs.RequireSsl).Should(BeFalse())
			Ω(s.Diego.SshProxy.EnableCfAuth).Should(BeTrue())    // tied to allow-app-ssh-access
			Ω(s.Diego.SshProxy.EnableDiegoAuth).Should(BeTrue()) // tied to allow-app-ssh-access
			Ω(s.Diego.SshProxy.Cc.ExternalPort).Should(Equal(9023))
			Ω(s.Diego.SshProxy.UaaTokenUrl).Should(Equal("https://uaa.sys.test.com/oauth/token"))
			Ω(s.Diego.SshProxy.UaaSecret).Should(Equal("secret"))
			Ω(s.Diego.SshProxy.HostKey).ShouldNot(BeEmpty())

			By("configuring the stager")
			job = ig.GetJobByName("stager")
			Ω(job.Release).Should(Equal(DiegoReleaseName))
			stager := job.Properties.(*stager.StagerJob)
			Ω(s.Diego.Ssl.SkipCertVerify).Should(BeFalse())
			Ω(stager.Diego.Stager.Bbs.ApiLocation).Should(Equal("bbs.service.cf.internal:8889"))
			Ω(stager.Diego.Stager.Bbs.CaCert).Should(Equal("cacert"))
			Ω(stager.Diego.Stager.Bbs.ClientCert).Should(Equal("clientcert"))
			Ω(stager.Diego.Stager.Bbs.ClientKey).Should(Equal("clientkey"))
			Ω(stager.Diego.Stager.Bbs.RequireSsl).Should(BeFalse())
			Ω(stager.Diego.Stager.Cc.ExternalPort).Should(Equal(9023))
			Ω(stager.Diego.Stager.Cc.BasicAuthUsername).Should(Equal("internaluser"))
			Ω(stager.Diego.Stager.Cc.BasicAuthPassword).Should(Equal("internalpassword"))

			By("configuring the tps")
			job = ig.GetJobByName("tps")
			Ω(job.Release).Should(Equal(DiegoReleaseName))
			t := job.Properties.(*tps.TpsJob)
			Ω(t.Diego.Ssl.SkipCertVerify).Should(BeFalse())
			Ω(t.Diego.Tps.TrafficControllerUrl).Should(Equal("wss://doppler.sys.yourdomain.com:443"))
			Ω(t.Diego.Tps.Bbs.ApiLocation).Should(Equal("bbs.service.cf.internal:8889"))
			Ω(t.Diego.Tps.Bbs.CaCert).Should(Equal("cacert"))
			Ω(t.Diego.Tps.Bbs.ClientCert).Should(Equal("clientcert"))
			Ω(t.Diego.Tps.Bbs.ClientKey).Should(Equal("clientkey"))
			Ω(t.Diego.Tps.Bbs.RequireSsl).Should(BeFalse())
			Ω(t.Diego.Tps.Cc.ExternalPort).Should(Equal(9023))
			Ω(t.Diego.Tps.Cc.BasicAuthUsername).Should(Equal("internaluser"))
			Ω(t.Diego.Tps.Cc.BasicAuthPassword).Should(Equal("internalpassword"))

			By("configuring the consul agent")
			job = ig.GetJobByName("consul_agent")
			consul := job.Properties.(*consul_agent.ConsulAgentJob)
			Ω(consul.Consul.ServerKey).Should(Equal("server-key"))
			Ω(consul.Consul.ServerCert).Should(Equal("server-cert"))
			Ω(consul.Consul.AgentCert).Should(Equal("agent-cert"))
			Ω(consul.Consul.AgentKey).Should(Equal("agent-key"))
			Ω(consul.Consul.CaCert).Should(Equal("ca-cert"))
			Ω(consul.Consul.EncryptKeys).Should(Equal([]string{"encyption-key"}))
			Ω(consul.Consul.Agent.Servers.Lan).Should(Equal([]string{"1.0.0.1", "1.0.0.2"}))

			By("configuring the metron agent")
			job = ig.GetJobByName("metron_agent")
			m := job.Properties.(*metron_agent.MetronAgentJob)
			Ω(m.MetronAgent.Zone).Should(Equal("metronzoneguid"))
			Ω(m.MetronAgent.Deployment).Should(Equal("cf"))
			Ω(m.MetronEndpoint.SharedSecret).Should(Equal("metronsecret"))
			Ω(m.Loggregator.Etcd.Machines).Should(Equal([]string{"1.0.0.7", "1.0.0.8"}))

			By("configuring the statsd injector")
			job = ig.GetJobByName("statsd-injector")
			Ω(job.Properties).Should(BeEmpty())
		})
	})
})
