package concourseplugin_test

import (
	"io/ioutil"

	. "github.com/enaml-ops/omg-product-bundle/products/concourse/plugin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("given ConcoursePlugin Plugin", func() {
	Context("when plugin is not properly initialized", func() {

		Context("when GetProduct is called with an empty cloud config", func() {
			It("should return an error", func() {
				p := new(ConcoursePlugin)
				_, err := p.GetProduct([]string{
					"test",
					"--network-name", "private",
					"--external-url", "http://concourse.caleb-washburn.com",
					"--concourse-username", "concourse",
					"--concourse-password", "concourse",
					"--az", "z1",
				}, []byte{}, nil)
				Ω(err).Should(HaveOccurred())
			})

			Context("when GetProduct is called with missing flag values", func() {
				It("should return an error", func() {
					p := new(ConcoursePlugin)
					_, err := p.GetProduct([]string{
						"test",
						"--network-name", "private",
						"--external-url", "http://concourse.caleb-washburn.com",
						"--concourse-username", "concourse",
						"--concourse-password", "concourse",
						"--az", "z1",
					}, []byte{0, 1, 2}, nil)
					Ω(err).Should(HaveOccurred())
				})
			})
		})

		It("returns an error if the external URL is missing the scheme", func() {
			p := new(ConcoursePlugin)
			cloudBytes, _ := ioutil.ReadFile("../fixtures/cloudconfig.yml")
			_, err := p.GetProduct([]string{
				"test",
				"--network-name", "private",
				"--external-url", "concourse.caleb-washburn.com", // <-- MISSING http://
				"--concourse-username", "concourse",
				"--concourse-password", "concourse",
				"--az", "z1",
				"--web-ip", "10.0.1.2",
				"--web-vm-type", "small",
				"--worker-vm-type", "medium",
				"--database-vm-type", "medium",
				"--database-storage-type", "large",
			}, cloudBytes, nil)
			Ω(err).Should(HaveOccurred())
		})
	})

	Context("when plugin is properly initialized", func() {
		var myplugin *ConcoursePlugin
		BeforeEach(func() {
			myplugin = new(ConcoursePlugin)
		})

		Context("when GetProduct is called with valid args", func() {
			var myconcourse []byte
			var err error
			BeforeEach(func() {
				cloudBytes, _ := ioutil.ReadFile("../fixtures/cloudconfig.yml")
				myconcourse, err = myplugin.GetProduct([]string{
					"test",
					"--network-name", "private",
					"--external-url", "http://concourse.caleb-washburn.com",
					"--concourse-username", "concourse",
					"--concourse-password", "concourse",
					"--az", "z1",
					"--web-ip", "10.0.1.2",
					"--web-vm-type", "small",
					"--worker-vm-type", "medium",
					"--database-vm-type", "medium",
					"--database-storage-type", "large",
				}, cloudBytes, nil)
			})
			It("then it should return the bytes representation of the object", func() {
				Ω(err).ShouldNot(HaveOccurred())
				Ω(myconcourse).ShouldNot(BeEmpty())
			})
		})
	})
})
