package gemfire_plugin_test

import (
	"github.com/enaml-ops/enaml"
	. "github.com/enaml-ops/omg-product-bundle/products/p-gemfire/plugin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("pgemfire plugin", func() {
	var gPlugin *Plugin
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

		It("should return error when server instance count is not provided", func() {
			_, err := gPlugin.GetProduct([]string{
				"pgemfire-command",
				"--az", "asdf",
				"--network-name", "asdf",
				"--locator-static-ip", "asdf",
				"--gemfire-locator-vm-size", "asdf",
				"--gemfire-server-vm-size", "asdf",
			}, []byte{})
			Expect(err).Should(HaveOccurred())
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
	})
})
