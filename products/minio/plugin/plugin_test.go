package plugin_test

import (
	"io/ioutil"

	. "github.com/enaml-ops/omg-product-bundle/products/minio/plugin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("given MinioPlugin Plugin", func() {
	Context("when plugin is not properly initialized", func() {

		Context("when GetProduct is called with an empty cloud config", func() {
			It("should return an error", func() {
				p := new(Plugin)
				_, err := p.GetProduct([]string{
					"test",
					"--network-name", "private",
					"--az", "z1",
				}, []byte{}, nil)
				Ω(err).Should(HaveOccurred())
			})

			Context("when GetProduct is called with missing flag values", func() {
				It("should return an error", func() {
					p := new(Plugin)
					_, err := p.GetProduct([]string{
						"test",
						"--network-name", "private",
						"--az", "z1",
					}, []byte{0, 1, 2}, nil)
					Ω(err).Should(HaveOccurred())
				})
			})
		})
	})

	Context("when plugin is properly initialized", func() {
		var myplugin *Plugin
		BeforeEach(func() {
			myplugin = new(Plugin)
		})

		Context("when GetProduct is called with valid args", func() {
			var myconcourse []byte
			var err error
			BeforeEach(func() {
				cloudBytes, _ := ioutil.ReadFile("../fixtures/cloudconfig.yml")
				myconcourse, err = myplugin.GetProduct([]string{
					"test",
					"--network-name", "private",
					"--az", "z1",
					"--ip", "10.0.1.2",
					"--vm-type", "small",
					"--disk-type", "small",
					"--access-key", "sample-access-key",
					"--secret-key", "sample-secret-key",
				}, cloudBytes, nil)
			})
			It("then it should return the bytes representation of the object", func() {
				Ω(err).ShouldNot(HaveOccurred())
				Ω(myconcourse).ShouldNot(BeEmpty())
			})
		})
	})
})
