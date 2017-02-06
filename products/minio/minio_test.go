package minio_test

import (
	"io/ioutil"

	. "github.com/enaml-ops/omg-product-bundle/products/minio"
	yaml "gopkg.in/yaml.v2"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Minio", func() {
	Context("NewDeployment", func() {
		It("Returns non nil deployment", func() {
			deployment := NewDeployment(Config{})
			Ω(deployment).Should(Not(BeNil()))
		})
	})
	Context("CloudConfigValidation", func() {
		var cloudConfig []byte
		var err error
		BeforeEach(func() {
			cloudConfig, err = ioutil.ReadFile("./fixtures/cloudconfig.yml")
			Ω(err).Should(Not(HaveOccurred()))
			Ω(cloudConfig).Should(Not(BeNil()))
		})
		It("No errors returned", func() {
			deployment := NewDeployment(Config{
				AZ:       "z1",
				VMType:   "small",
				DiskType: "small",
			})
			err := deployment.CloudConfigValidation(cloudConfig)
			Ω(err).Should(Not(HaveOccurred()))
		})
		It("Returns error for az name not found", func() {
			deployment := NewDeployment(Config{
				AZ: "z",
			})
			err := deployment.CloudConfigValidation(cloudConfig)
			Ω(err).Should(HaveOccurred())
			Ω(err.Error()).Should(BeEquivalentTo("AZ [z] is not defined as a AZ in cloud config"))
		})
		It("Returns error for vmtype name not found", func() {
			deployment := NewDeployment(Config{
				AZ:     "z1",
				VMType: "blah",
			})
			err := deployment.CloudConfigValidation(cloudConfig)
			Ω(err).Should(HaveOccurred())
			Ω(err.Error()).Should(BeEquivalentTo("VMType[blah] is not defined as a VMType in cloud config"))
		})
		It("Returns error for disktype name not found", func() {
			deployment := NewDeployment(Config{
				AZ:       "z1",
				VMType:   "small",
				DiskType: "blah",
			})
			err := deployment.CloudConfigValidation(cloudConfig)
			Ω(err).Should(HaveOccurred())
			Ω(err.Error()).Should(BeEquivalentTo("DiskType[blah] is not defined as a DiskType in cloud config"))
		})
		It("Returns error when no yaml for cloud config", func() {
			deployment := NewDeployment(Config{})
			err := deployment.CloudConfigValidation([]byte("asdfasd"))
			Ω(err).Should(HaveOccurred())
		})
	})
	Context("CreateDeploymentManifest", func() {
		var cloudConfig []byte
		var err error
		BeforeEach(func() {
			cloudConfig, err = ioutil.ReadFile("./fixtures/cloudconfig.yml")
			Ω(err).Should(Not(HaveOccurred()))
			Ω(cloudConfig).Should(Not(BeNil()))
		})
		It("No errors returned", func() {
			config := NewConfig()
			config.AZ = "z1"
			config.VMType = "small"
			config.DiskType = "small"
			config.NetworkName = "theNetwork"
			config.IP = "10.244.0.2"
			config.Region = DefaultRegion
			config.AccessKey = "sample-access-key"
			config.SecretKey = "sample-secret-key"
			deployment := NewDeployment(config)
			dm, err := deployment.CreateDeploymentManifest(cloudConfig)
			Ω(err).Should(Not(HaveOccurred()))
			Ω(dm).Should(Not(BeNil()))

			controlYaml, err := ioutil.ReadFile("./fixtures/manifest.yml")
			Ω(err).ShouldNot(HaveOccurred())
			b, err := yaml.Marshal(dm)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(b).Should(MatchYAML(string(controlYaml)))
		})
	})
})
