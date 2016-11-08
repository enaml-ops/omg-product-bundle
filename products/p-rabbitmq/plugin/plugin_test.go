package prabbitmq_test

import (
	"io/ioutil"
	"net/http"

	"github.com/enaml-ops/enaml"
	prabbitmq "github.com/enaml-ops/omg-product-bundle/products/p-rabbitmq/plugin"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
	logging "github.com/op/go-logging"
)

var _ = Describe("prabbitmq plugin", func() {

	BeforeSuite(func() {
		logging.SetBackend(logging.NewLogBackend(ioutil.Discard, "", 0))
	})

	Context("when generating a manifest with incomplete input", func() {
		var (
			p *prabbitmq.Plugin
		)

		BeforeEach(func() {
			p = new(prabbitmq.Plugin)
		})
		It("then we should return error", func() {
			_, err := p.GetProduct([]string{"foo"}, []byte{}, nil)
			Ω(err).Should(HaveOccurred())
		})
	})

	Context("when we are testing defaults????", func() {

		var (
			p  *prabbitmq.Plugin
			dm *enaml.DeploymentManifest
		)

		BeforeEach(func() {
			p = new(prabbitmq.Plugin)
			manifestBytes, err := p.GetProduct([]string{"foo",
				"--az", "asdf",
				"--system-domain", "asdf",
				"--network", "asdf",
				"--rabbit-server-ip", "asdf",
				"--rabbit-broker-ip", "asdf",
				"--system-services-password", "asdf",
				"--doppler-zone", "asdf",
				"--doppler-shared-secret", "asdf",
				"--etcd-machine-ip", "asdf",
				"--rabbit-public-ip", "asdf",
				"--rabbit-broker-vm-type", "asdf",
				"--rabbit-server-vm-type", "asdf",
				"--rabbit-haproxy-vm-type", "asdf",
				"--syslog-address", "asdf",
				"--nats-machine-ip", "asdf",
			}, []byte{}, nil)
			Ω(err).ShouldNot(HaveOccurred())
			dm = enaml.NewDeploymentManifest(manifestBytes)
		})
		It("should have the correct releases", func() {
			hasRelease := func(name, version string) bool {
				for i := range dm.Releases {
					if dm.Releases[i].Name == name && dm.Releases[i].Version == version {
						return true
					}
				}
				return false
			}

			Ω(hasRelease(prabbitmq.CFRabbitMQReleaseName, prabbitmq.CFRabbitMQReleaseVersion)).Should(BeTrue())
			Ω(hasRelease(prabbitmq.ServiceMetricsReleaseName, prabbitmq.ServiceMetricsReleaseVersion)).Should(BeTrue())
			Ω(hasRelease(prabbitmq.LoggregatorReleaseName, prabbitmq.LoggregatorReleaseVersion)).Should(BeTrue())
			Ω(hasRelease(prabbitmq.RabbitMQMetricsReleaseName, prabbitmq.RabbitMQMetricsReleaseVersion)).Should(BeTrue())
		})

		It("should have the correct instance groups", func() {
			Ω(dm.GetInstanceGroupByName("rabbitmq-server-partition")).ShouldNot(BeNil())
			Ω(dm.GetInstanceGroupByName("rabbitmq-broker-partition")).ShouldNot(BeNil())
			Ω(dm.GetInstanceGroupByName("rabbitmq-haproxy-partition")).ShouldNot(BeNil())
			Ω(dm.GetInstanceGroupByName("broker-registrar")).ShouldNot(BeNil())
			Ω(dm.GetInstanceGroupByName("broker-deregistrar")).ShouldNot(BeNil())
		})

		It("should set the update", func() {
			Ω(dm.Update.Canaries).Should(Equal(1))
			Ω(dm.Update.CanaryWatchTime).Should(Equal("30000-300000"))
			Ω(dm.Update.UpdateWatchTime).Should(Equal("30000-300000"))
			Ω(dm.Update.MaxInFlight).Should(Equal(1))
			Ω(dm.Update.Serial).Should(BeTrue())
		})

		It("should configure the stemcell", func() {
			Ω(dm.Stemcells).Should(HaveLen(1))
			Ω(dm.Stemcells[0].OS).Should(Equal(prabbitmq.StemcellName))
			Ω(dm.Stemcells[0].Alias).Should(Equal(prabbitmq.StemcellAlias))
			Ω(dm.Stemcells[0].Version).Should(Equal(prabbitmq.StemcellVersion))
		})
	})

	Context("when inferring defaults from cloud config", func() {
		var (
			p  *prabbitmq.Plugin
			dm *enaml.DeploymentManifest
		)
		BeforeEach(func() {
			p = new(prabbitmq.Plugin)
			cc, err := ioutil.ReadFile("fixtures/cloudconfig.yml")
			Ω(err).ShouldNot(HaveOccurred())

			manifestBytes, err := p.GetProduct([]string{"foo",
				"--infer-from-cloud",
				"--rabbit-broker-vm-type", "large",
				"--system-domain", "asdf",
				"--rabbit-server-ip", "asdf",
				"--rabbit-broker-ip", "asdf",
				"--system-services-password", "asdf",
				"--doppler-zone", "asdf",
				"--doppler-shared-secret", "asdf",
				"--etcd-machine-ip", "asdf",
				"--rabbit-public-ip", "asdf",
				"--syslog-address", "asdf",
				"--nats-machine-ip", "asdf",
			}, cc, nil)
			Ω(err).ShouldNot(HaveOccurred())
			dm = enaml.NewDeploymentManifest(manifestBytes)
		})

		It("should not overwrite flags that were specified on the command line", func() {
			ig := dm.GetInstanceGroupByName("rabbitmq-broker-partition")
			Ω(ig).ShouldNot(BeNil())
			Ω(ig.VMType).Should(Equal("large"))
		})

		It("should infer VM types not specified on the command line", func() {
			ig := dm.GetInstanceGroupByName("rabbitmq-server-partition")
			Ω(ig).ShouldNot(BeNil())
			Ω(ig.VMType).Should(Equal("smallvm"))

			ig = dm.GetInstanceGroupByName("rabbitmq-haproxy-partition")
			Ω(ig).ShouldNot(BeNil())
			Ω(ig.VMType).Should(Equal("smallvm"))
		})

		It("should infer the network name", func() {
			for _, name := range []string{"rabbitmq-server-partition", "rabbitmq-haproxy-partition", "rabbitmq-broker-partition"} {
				ig := dm.GetInstanceGroupByName(name)
				Ω(ig).ShouldNot(BeNil(), "couldn't find instance group "+name)
				Ω(ig.Networks).Should(HaveLen(1))
				Ω(ig.Networks[0].Name).Should(Equal("privatenet"))
			}
		})

		It("should infer the AZs", func() {
			for _, name := range []string{"rabbitmq-server-partition", "rabbitmq-haproxy-partition", "rabbitmq-broker-partition"} {
				ig := dm.GetInstanceGroupByName(name)
				Ω(ig).ShouldNot(BeNil(), "couldn't find instance group "+name)
				Ω(ig.AZs).Should(ConsistOf("z1", "z2"))
			}
		})
	})

	Context("when using Vault integration", func() {
		var (
			p      *prabbitmq.Plugin
			dm     *enaml.DeploymentManifest
			server *ghttp.Server
		)

		BeforeEach(func() {
			hash1, err := ioutil.ReadFile("fixtures/hash1.json")
			Ω(err).ShouldNot(HaveOccurred())
			hash2, err := ioutil.ReadFile("fixtures/hash2.json")
			Ω(err).ShouldNot(HaveOccurred())

			server = ghttp.NewServer()
			server.AllowUnhandledRequests = true
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/v1/secret/hash1"),
					ghttp.RespondWith(http.StatusOK, hash1),
				),
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/v1/secret/hash2"),
					ghttp.RespondWith(http.StatusOK, hash2),
				),
			)

			p = new(prabbitmq.Plugin)
			manifestBytes, err := p.GetProduct([]string{
				"rabbitmq",
				"--vault-domain", server.URL(),
				"--vault-token", "asdfghjkl",
				"--vault-hash", "secret/hash1",
				"--vault-hash", "secret/hash2",
				"--az", "z1",
				"--system-domain", "asdf",
				"--network", "asdf",
				"--rabbit-server-ip", "asdf",
				"--rabbit-broker-ip", "asdf",
				"--system-services-password", "asdf",
				"--doppler-zone", "asdf",
				"--doppler-shared-secret", "asdf",
				"--etcd-machine-ip", "asdf",
				"--rabbit-public-ip", "asdf",
				"--rabbit-broker-vm-type", "asdf",
				"--rabbit-server-vm-type", "asdf",
				"--rabbit-haproxy-vm-type", "asdf",
			}, []byte{}, nil)
			Ω(err).ShouldNot(HaveOccurred())
			dm = enaml.NewDeploymentManifest(manifestBytes)
		})

		AfterEach(func() {
			server.Close()
		})

		It("then it should use the values from vault", func() {
			broker := dm.GetInstanceGroupByName("rabbitmq-broker-partition")
			brokerJob := broker.GetJobByName("rabbitmq-broker")
			props := brokerJob.Properties.(map[interface{}]interface{})

			Ω(props["syslog_aggregator"]).Should(HaveKeyWithValue("address", "1.0.0.5"))

			nats := props["cf"].(map[interface{}]interface{})["nats"]
			Ω(nats).Should(HaveKeyWithValue("password", "natspassword"))
			machines := nats.(map[interface{}]interface{})["machines"]
			Ω(machines).Should(ConsistOf("1.0.0.3", "1.0.0.4"))
		})

	})
})
