package concourse_test

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"

	. "github.com/enaml-ops/omg-product-bundle/products/concourse"
	"github.com/enaml-ops/omg-product-bundle/products/concourse/enaml-gen/atc"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Concourse Deployment", func() {
	var deployment Deployment
	BeforeEach(func() {
		deployment = NewDeployment()
	})
	Describe("Given a CreateUpdate", func() {
		Context("when calling", func() {
			It("then we should return a valid enaml.Update", func() {
				update := deployment.CreateUpdate()
				Ω(update.Canaries).Should(Equal(1))
				Ω(update.MaxInFlight).Should(Equal(3))
				Ω(update.Serial).Should(Equal(false))
				Ω(update.UpdateWatchTime).Should(Equal("1000-60000"))
				Ω(update.CanaryWatchTime).Should(Equal("1000-60000"))
			})
		})
	})

	Describe("given CreateAtcJob", func() {
		Context("when called without TLS cert/key flags", func() {
			It("should not emit YAML for TLS settings", func() {
				job := deployment.CreateAtcJob()
				Ω(job).ShouldNot(BeNil())
				b, _ := yaml.Marshal(job.Properties)

				var m map[string]interface{}
				yaml.Unmarshal(b, &m)
				Ω(m).ShouldNot(HaveKey("tls_key"))
				Ω(m).ShouldNot(HaveKey("tls_cert"))
			})
		})
		Context("when called with an external URL", func() {
			const controlURL = "https://myconcourse.com"
			It("sets the URL correctly", func() {
				deployment.ConcourseURL = controlURL
				job := deployment.CreateAtcJob()
				props := job.Properties.(atc.AtcJob)
				Ω(props.ExternalUrl).Should(Equal(controlURL))
			})
		})
	})

	Describe("Given a CreateWebInstanceGroup", func() {
		Context("when calling with WebAzs and StemcellAlias on deployment", func() {
			It("then we should return a valid *enaml.InstanceGroup", func() {
				deployment.WebIPs = []string{"10.0.0.10"}
				deployment.AZs = []string{"z1"}
				deployment.StemcellAlias = "trusty"
				deployment.WebVMType = "small"
				web, err := deployment.CreateWebInstanceGroup()
				Ω(err).Should(BeNil())
				Ω(web.Name).Should(Equal("web"))
				Ω(web.Instances).Should(Equal(1))
				Ω(web.ResourcePool).Should(Equal(""))
				Ω(web.AZs).Should(Equal([]string{"z1"}))
				Ω(web.PersistentDisk).Should(Equal(0))
				Ω(web.Stemcell).Should(Equal("trusty"))
				Ω(web.VMType).Should(Equal("small"))
				Ω(len(web.Networks)).Should(Equal(1))
				Ω(len(web.Jobs)).Should(Equal(2))
			})
		})
	})
	Describe("Given a CreateDatabaseInstanceGroup", func() {
		Context("when calling with Azs and Stemcell on deployment", func() {
			It("then we should return a valid *enaml.InstanceGroup", func() {
				deployment.AZs = []string{"z1"}
				deployment.StemcellAlias = "trusty"
				deployment.DatabaseVMType = "medium"
				deployment.DatabaseStorageType = "medium"
				database, err := deployment.CreateDatabaseInstanceGroup()
				Ω(err).Should(BeNil())
				Ω(database.Name).Should(Equal("db"))
				Ω(database.Instances).Should(Equal(1))
				Ω(database.ResourcePool).Should(Equal(""))
				Ω(database.AZs).Should(Equal([]string{"z1"}))
				Ω(database.PersistentDisk).Should(Equal(0))
				Ω(database.PersistentDiskType).Should(Equal("medium"))
				Ω(database.Stemcell).Should(Equal("trusty"))
				Ω(database.VMType).Should(Equal("medium"))
				Ω(len(database.Networks)).Should(Equal(1))
				Ω(len(database.Jobs)).Should(Equal(1))
			})
		})
	})
	Describe("Given a CreateWorkerInstanceGroup", func() {
		Context("when calling with Azs and Stemcell on deployment", func() {
			It("then we should return a valid *enaml.InstanceGroup", func() {
				deployment.AZs = []string{"z1"}
				deployment.StemcellAlias = "trusty"
				deployment.WorkerVMType = "medium"
				deployment.WorkerInstances = 1
				worker, err := deployment.CreateWorkerInstanceGroup()
				Ω(err).Should(BeNil())
				Ω(worker.Name).Should(Equal("worker"))
				Ω(worker.Instances).Should(Equal(1))
				Ω(worker.ResourcePool).Should(Equal(""))
				Ω(worker.AZs).Should(Equal([]string{"z1"}))
				Ω(worker.PersistentDisk).Should(Equal(0))
				Ω(worker.Stemcell).Should(Equal("trusty"))
				Ω(worker.VMType).Should(Equal("medium"))
				Ω(len(worker.Networks)).Should(Equal(1))
				Ω(len(worker.Jobs)).Should(Equal(3))
			})
		})
	})
	Describe("Given a new deployment and a valid cloud config", func() {
		var cc []byte

		BeforeEach(func() {
			cc, _ = ioutil.ReadFile("fixtures/cloudconfig.yml")
			azs := []string{"z1"}
			deployment.WebVMType = "small"
			deployment.WorkerVMType = "small"
			deployment.AZs = azs
			deployment.DatabaseVMType = "small"
			deployment.DatabaseStorageType = "large"
		})

		XContext("when calling Initialize without a strong password", func() {
			deployment.ConcoursePassword = "test"
			It("then we should error and prompt the user for a better pass", func() {
				err := deployment.Initialize([]byte(""))
				Ω(err).ShouldNot(BeNil())
			})
		})

		Context("when initializing without remote stemcell flags", func() {
			BeforeEach(func() {
				deployment.StemcellAlias = "trusty"
				deployment.StemcellVersion = "3262.2"
				deployment.ConcoursePassword = "password"
			})
			It("initializes without error", func() {
				err := deployment.Initialize(cc)
				Ω(err).ShouldNot(HaveOccurred())
			})
		})

		Context("when initializing with all required remote stemcell flags", func() {
			BeforeEach(func() {
				deployment.StemcellAlias = "trusty"
				deployment.StemcellVersion = "3262.2"

				It("initializes without error", func() {
					err := deployment.Initialize(cc)
					Ω(err).ShouldNot(HaveOccurred())
				})

			})
		})
	})

	Describe("Given a cloud config", func() {
		Context("when validating", func() {
			var cc []byte
			BeforeEach(func() {
				cc, _ = ioutil.ReadFile("fixtures/cloudconfig.yml")
				azs := []string{"z1"}
				deployment.WebVMType = "small"
				deployment.WorkerVMType = "small"
				deployment.AZs = azs
				deployment.DatabaseVMType = "small"
				deployment.DatabaseStorageType = "large"
				deployment.StemcellAlias = "ubuntu"
				deployment.ConcoursePassword = "password"
			})

			It("should be valid", func() {
				Ω(len(cc)).ShouldNot(BeZero())
				err := deployment.Initialize(cc)
				Ω(err).Should(BeNil())
			})
		})
	})
})
