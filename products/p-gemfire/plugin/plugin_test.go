package gemfire_plugin_test

import (
	"fmt"
	"os"
	"strconv"

	"gopkg.in/yaml.v2"

	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-product-bundle/products/p-gemfire/enaml-gen/server"
	. "github.com/enaml-ops/omg-product-bundle/products/p-gemfire/plugin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("pgemfire plugin", func() {
	var gPlugin *Plugin

	Context("When a commnd line args are populate as ENV Vars", func() {

		var controlAZName = "blah"
		BeforeEach(func() {
			gPlugin = &Plugin{Version: "0.0"}
			os.Setenv("OMG_AZ", controlAZName)
		})
		AfterEach(func() {
			os.Setenv("OMG_AZ", "")
		})

		It("should pass validation of required flags", func() {
			_, err := gPlugin.GetProduct([]string{
				"pgemfire-command",
				"--network-name", "net1",
				"--locator-static-ip", "1.0.0.2",
				"--server-instance-count", "1",
				"--gemfire-locator-vm-size", "asdf",
				"--gemfire-server-vm-size", "asdf",
			}, []byte{})
			Expect(err).ShouldNot(HaveOccurred(), "should pass env var isset required value check")
		})

		It("should properly set up the Availability Zones", func() {
			manifestBytes, err := gPlugin.GetProduct([]string{
				"pgemfire-command",
				"--network-name", "net1",
				"--locator-static-ip", "1.0.0.2",
				"--server-instance-count", "1",
				"--gemfire-locator-vm-size", "asdf",
				"--gemfire-server-vm-size", "asdf",
			}, []byte{})
			Expect(err).ShouldNot(HaveOccurred())
			manifest := enaml.NewDeploymentManifest(manifestBytes)
			for _, instanceGroup := range manifest.InstanceGroups {
				Expect(instanceGroup.AZs).Should(Equal([]string{controlAZName}), fmt.Sprintf("Availability ZOnes for instance group %v was not set properly", instanceGroup.Name))
			}
		})
	})

	Context("When a commnd line args are passed", func() {
		BeforeEach(func() {
			gPlugin = &Plugin{Version: "0.0"}
		})

		It("should contain a server and locator instance group", func() {
			manifestBytes, err := gPlugin.GetProduct([]string{
				"pgemfire-command",
				"--az", "z1",
				"--network-name", "net1",
				"--locator-static-ip", "1.0.0.2",
				"--server-instance-count", "1",
				"--gemfire-locator-vm-size", "asdf",
				"--gemfire-server-vm-size", "asdf",
			}, []byte{})
			Expect(err).ShouldNot(HaveOccurred())
			manifest := enaml.NewDeploymentManifest(manifestBytes)
			locator := manifest.GetInstanceGroupByName("locator-group")
			Expect(locator).ShouldNot(BeNil())
			server := manifest.GetInstanceGroupByName("server-group")
			Expect(server).ShouldNot(BeNil())
		})

		It("should return error when AZ/s are not provided", func() {
			_, err := gPlugin.GetProduct([]string{
				"pgemfire-command",
				"--network-name", "asdf",
				"--locator-static-ip", "asdf",
				"--server-instance-count", "1",
				"--gemfire-locator-vm-size", "asdf",
				"--gemfire-server-vm-size", "asdf",
			}, []byte{})
			Expect(err).Should(HaveOccurred())
		})

		It("should return error when network name is not provided", func() {
			_, err := gPlugin.GetProduct([]string{
				"pgemfire-command",
				"--az", "asdf",
				"--locator-static-ip", "asdf",
				"--server-instance-count", "1",
				"--gemfire-locator-vm-size", "asdf",
				"--gemfire-server-vm-size", "asdf",
			}, []byte{})
			Expect(err).Should(HaveOccurred())
		})

		It("should return error when locator IPs are not given", func() {
			_, err := gPlugin.GetProduct([]string{
				"pgemfire-command",
				"--az", "asdf",
				"--network-name", "asdf",
				"--server-instance-count", "1",
				"--gemfire-locator-vm-size", "asdf",
				"--gemfire-server-vm-size", "asdf",
			}, []byte{})
			Expect(err).Should(HaveOccurred())
		})

		It("should not require server-instance-count field", func() {
			_, err := gPlugin.GetProduct([]string{
				"pgemfire-command",
				"--az", "asdf",
				"--network-name", "asdf",
				"--locator-static-ip", "asdf",
				"--gemfire-locator-vm-size", "asdf",
				"--gemfire-server-vm-size", "asdf",
			}, []byte{})
			Expect(err).ShouldNot(HaveOccurred(), "server-instance-count should not be required")
		})

		It("should return error when --gemfire-locator-vm-size is not provided", func() {
			_, err := gPlugin.GetProduct([]string{
				"pgemfire-command",
				"--az", "asdf",
				"--network-name", "asdf",
				"--locator-static-ip", "asdf",
				"--server-instance-count", "1",
				"--gemfire-server-vm-size", "asdf",
			}, []byte{})
			Expect(err).Should(HaveOccurred())
		})

		It("should return error when --gemfire-server-vm-size is not provided", func() {
			_, err := gPlugin.GetProduct([]string{
				"pgemfire-command",
				"--az", "asdf",
				"--network-name", "asdf",
				"--locator-static-ip", "asdf",
				"--server-instance-count", "1",
				"--gemfire-locator-vm-size", "asdf",
			}, []byte{})
			Expect(err).Should(HaveOccurred())
		})

		It("should properly set up the Update segment", func() {
			controlStemcellAlias := "ubuntu-magic"
			manifestBytes, err := gPlugin.GetProduct([]string{
				"pgemfire-command",
				"--az", "z1",
				"--network-name", "net1",
				"--locator-static-ip", "1.0.0.2",
				"--server-instance-count", "1",
				"--gemfire-locator-vm-size", "asdf",
				"--gemfire-server-vm-size", "asdf",
				"--stemcell-alias", controlStemcellAlias,
			}, []byte{})
			Expect(err).ShouldNot(HaveOccurred())
			manifest := enaml.NewDeploymentManifest(manifestBytes)
			Ω(manifest.Update.MaxInFlight).ShouldNot(Equal(0), "max in flight")
			Ω(manifest.Update.UpdateWatchTime).ShouldNot(BeEmpty(), "update watch time")
			Ω(manifest.Update.CanaryWatchTime).ShouldNot(BeEmpty(), "canary watch time")
			Ω(manifest.Update.Canaries).ShouldNot(Equal(0), "number of canaries")
		})

		It("should properly set up the gemfire release", func() {
			controlStemcellAlias := "ubuntu-magic"
			manifestBytes, err := gPlugin.GetProduct([]string{
				"pgemfire-command",
				"--az", "z1",
				"--network-name", "net1",
				"--locator-static-ip", "1.0.0.2",
				"--server-instance-count", "1",
				"--gemfire-locator-vm-size", "asdf",
				"--gemfire-server-vm-size", "asdf",
				"--stemcell-alias", controlStemcellAlias,
			}, []byte{})
			Expect(err).ShouldNot(HaveOccurred())
			manifest := enaml.NewDeploymentManifest(manifestBytes)
			Ω(manifest.Releases).ShouldNot(BeEmpty())
			Ω(manifest.Releases[0]).ShouldNot(BeNil())
			Ω(manifest.Releases[0].Name).Should(Equal("GemFire"))
			Ω(manifest.Releases[0].Version).Should(Equal("latest"))
		})

		It("should properly set up the stemcells", func() {
			controlStemcellAlias := "ubuntu-magic"
			manifestBytes, err := gPlugin.GetProduct([]string{
				"pgemfire-command",
				"--az", "z1",
				"--network-name", "net1",
				"--locator-static-ip", "1.0.0.2",
				"--server-instance-count", "1",
				"--gemfire-locator-vm-size", "asdf",
				"--gemfire-server-vm-size", "asdf",
				"--stemcell-alias", controlStemcellAlias,
			}, []byte{})
			Expect(err).ShouldNot(HaveOccurred())
			manifest := enaml.NewDeploymentManifest(manifestBytes)
			Ω(manifest.Stemcells).ShouldNot(BeNil())
			Ω(manifest.Stemcells[0].Alias).Should(Equal(controlStemcellAlias))
			for _, instanceGroup := range manifest.InstanceGroups {
				Expect(instanceGroup.Stemcell).Should(Equal(controlStemcellAlias), fmt.Sprintf("stemcell for instance group %v was not set properly", instanceGroup.Name))
			}
		})

		It("should properly set up the deployment name", func() {
			var controlName = "p-gemfire-name"
			manifestBytes, err := gPlugin.GetProduct([]string{
				"pgemfire-command",
				"--az", "z1",
				"--deployment-name", controlName,
				"--network-name", "net1",
				"--locator-static-ip", "1.0.0.2",
				"--server-instance-count", "1",
				"--gemfire-locator-vm-size", "asdf",
				"--gemfire-server-vm-size", "asdf",
			}, []byte{})
			Expect(err).ShouldNot(HaveOccurred())
			manifest := enaml.NewDeploymentManifest(manifestBytes)
			Ω(manifest.Name).Should(Equal(controlName))
		})

		It("should properly set up the Availability Zones", func() {
			var controlAZ = "z1"
			manifestBytes, err := gPlugin.GetProduct([]string{
				"pgemfire-command",
				"--az", controlAZ,
				"--network-name", "net1",
				"--locator-static-ip", "1.0.0.2",
				"--server-instance-count", "1",
				"--gemfire-locator-vm-size", "asdf",
				"--gemfire-server-vm-size", "asdf",
			}, []byte{})
			Expect(err).ShouldNot(HaveOccurred())
			manifest := enaml.NewDeploymentManifest(manifestBytes)
			for _, instanceGroup := range manifest.InstanceGroups {
				Expect(instanceGroup.AZs).Should(Equal([]string{controlAZ}), fmt.Sprintf("Availability ZOnes for instance group %v was not set properly", instanceGroup.Name))
			}
		})
		Context("and command line args are setting server static ips", func() {
			It("should properly set up the number of servers in cluster topology", func() {
				var controlAZ = "z1"
				var controlIPCount = "6"
				var givenStaticCount = 1
				manifestBytes, err := gPlugin.GetProduct([]string{
					"pgemfire-command",
					"--az", controlAZ,
					"--network-name", "net1",
					"--locator-static-ip", "1.0.0.2",
					"--server-instance-count", controlIPCount,
					"--server-static-ip", "1.0.0.3",
					"--gemfire-locator-vm-size", "asdf",
					"--gemfire-server-vm-size", "asdf",
				}, []byte{})
				Expect(err).ShouldNot(HaveOccurred())
				manifest := enaml.NewDeploymentManifest(manifestBytes)
				instanceGroup := manifest.GetInstanceGroupByName("server-group")
				var properties = new(server.ServerJob)
				propertiesBytes, _ := yaml.Marshal(instanceGroup.GetJobByName("server").Properties)
				yaml.Unmarshal(propertiesBytes, properties)
				Expect(properties.Gemfire.ClusterTopology.NumberOfServers).Should(Equal(givenStaticCount), fmt.Sprintf("we should match ips given not instance count given"))
			})
		})

		Context("and command line args are setting dev rest apiu config", func() {
			It("should properly set up the settings for the dev rest api", func() {
				var controlPort = 4567
				var controlActive = "true"
				manifestBytes, err := gPlugin.GetProduct([]string{
					"pgemfire-command",
					"--az", "blah",
					"--network-name", "net1",
					"--locator-static-ip", "1.0.0.2",
					"--server-instance-count", "3",
					"--server-static-ip", "1.0.0.3",
					"--gemfire-locator-vm-size", "asdf",
					"--gemfire-server-vm-size", "asdf",
					"--gemfire-dev-rest-api-port", strconv.Itoa(controlPort),
					"--gemfire-dev-rest-api-active", controlActive,
				}, []byte{})
				Expect(err).ShouldNot(HaveOccurred())
				manifest := enaml.NewDeploymentManifest(manifestBytes)
				instanceGroup := manifest.GetInstanceGroupByName("server-group")
				var properties = new(server.ServerJob)
				propertiesBytes, _ := yaml.Marshal(instanceGroup.GetJobByName("server").Properties)
				yaml.Unmarshal(propertiesBytes, properties)
				Expect(properties.Gemfire.Server.DevRestApi.Port).Should(Equal(controlPort), fmt.Sprintf("should overwrite the port default"))
				Expect(properties.Gemfire.Server.DevRestApi.Active).Should(BeTrue(), fmt.Sprintf("should overwrite the active default"))
			})
		})
	})
})
