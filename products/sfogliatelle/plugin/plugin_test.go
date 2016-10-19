package sfogliatelle_test

import (
	"os"

	"github.com/enaml-ops/enaml"
	. "github.com/enaml-ops/omg-product-bundle/products/sfogliatelle/plugin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("sfogliatelle plugin", func() {
	var splugin *Plugin

	var f *os.File
	BeforeEach(func() {
		splugin = &Plugin{Version: "0.0"}
		f, _ = os.Open("fixtures/stdin.yml")
		splugin.Source = f
	})

	AfterEach(func() {
		f.Close()
	})

	Context("When a NOT given valid flags", func() {
		It("should pass validation of required flags", func() {
			_, err := splugin.GetProduct([]string{
				"sfogliatelle-command",
			}, []byte{})
			Expect(err).Should(HaveOccurred(), "we should error if the proper flags are not given")
		})
	})

	Context("When a given valid flags", func() {

		It("should pass validation of required flags", func() {
			_, err := splugin.GetProduct([]string{
				"sfogliatelle-command",
				"--layer-file", "fixtures/instance-layer.yml",
			}, []byte{})
			Expect(err).ShouldNot(HaveOccurred(), "should pass env var isset required value check")
		})

		Context("and a deployment manifest has been piped in for an instance group layer", func() {

			It("should read the yaml file and layer it ontop of any existing deployment manifest", func() {
				manifestBytes, err := splugin.GetProduct([]string{
					"sfogliatelle-command",
					"--layer-file", "fixtures/instance-layer.yml",
					"--instance-group-name", "rabbitmq-server-partition",
				}, []byte{})
				Expect(err).ShouldNot(HaveOccurred())
				manifest := enaml.NewDeploymentManifest(manifestBytes)
				ig := manifest.GetInstanceGroupByName("rabbitmq-server-partition")
				Ω(ig.Instances).ShouldNot(Equal(3))
				Ω(ig.Instances).Should(Equal(10))
			})
		})

		Context("and a deployment manifest has been piped in for a Job layer", func() {

			It("should read the yaml file and layer it ontop of any existing deployment manifest", func() {
				manifestBytes, err := splugin.GetProduct([]string{
					"sfogliatelle-command",
					"--layer-file", "fixtures/job-layer.yml",
					"--instance-group-name", "rabbitmq-server-partition",
					"--job-name", "metron_agent",
				}, []byte{})
				Expect(err).ShouldNot(HaveOccurred())
				manifest := enaml.NewDeploymentManifest(manifestBytes)
				ig := manifest.GetInstanceGroupByName("rabbitmq-server-partition")
				ma := ig.GetJobByName("metron_agent")
				maProperties := ma.Properties.(map[interface{}]interface{})
				secret := maProperties["metron_endpoint"].(map[interface{}]interface{})["shared_secret"]
				Ω(secret).ShouldNot(Equal("shhhhdonttell"))
				Ω(secret).Should(Equal("somethingelse"))
			})
		})

		Context("and a deployment manifest has been piped in for a deployment layer", func() {

			It("should read the yaml file and layer it ontop of any existing deployment manifest", func() {
				manifestBytes, err := splugin.GetProduct([]string{
					"sfogliatelle-command",
					"--layer-file", "fixtures/deploy-layer.yml",
				}, []byte{})
				Expect(err).ShouldNot(HaveOccurred())
				manifest := enaml.NewDeploymentManifest(manifestBytes)

				for _, release := range manifest.Releases {
					if release.Name == "cf-rabbitmq" {
						Ω(release.Version).ShouldNot(Equal("215.8.0"))
						Ω(release.Version).Should(Equal("300.8.0"))
					}

					if release.Name == "service-metrics" {
						Ω(release.Version).ShouldNot(Equal("1.4.3"))
						Ω(release.Version).Should(Equal("2.4.3"))
					}

					if release.Name == "loggregator" {
						Ω(release.Version).ShouldNot(Equal("9"))
						Ω(release.Version).Should(Equal("10"))
					}

					if release.Name == "rabbitmq-metrics" {
						Ω(release.Version).ShouldNot(Equal("1.29.0"))
						Ω(release.Version).Should(Equal("1.30.0"))
					}
				}
			})
		})
	})
})
