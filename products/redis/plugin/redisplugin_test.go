package redis_test

import (
	"fmt"
	"io/ioutil"

	"github.com/enaml-ops/enaml"
	. "github.com/enaml-ops/omg-product-bundle/products/redis/plugin"
	"github.com/enaml-ops/pluginlib/pluginutil"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/urfave/cli.v2"
)

var _ = Describe("given redis Plugin", func() {
	var plgn *Plugin

	BeforeEach(func() {
		plgn = new(Plugin)
	})

	Context("when calling GetProduct while targeting an un-compatible cloud config'd bosh", func() {
		var cloudConfigBytes []byte
		var controlInstances = "1"
		var controlNetName = "hello"
		var controlPass = "pss"
		var controlDisk = "4033"
		var controlVM = "large"
		var controlIP = "1.2.3.4"

		BeforeEach(func() {
			cloudConfigBytes, _ = ioutil.ReadFile("./fixtures/sample-aws.yml")
		})

		It("then we should fail fast and give the user guidance on what is wrong", func() {
			_, err := plgn.GetProduct([]string{
				"appname",
				"--disk-size", controlDisk,
				"--leader-instances", controlInstances,
				"--network-name", controlNetName,
				"--redis-pass", controlPass,
				"--vm-size", controlVM,
				"--leader-ip", controlIP,
				"--slave-ip", controlIP,
				"--stemcell-ver", "12.3.44",
			}, cloudConfigBytes, nil)
			Ω(err).Should(HaveOccurred())
		})
	})

	Context("when calling plugin without all required flags", func() {
		It("then it should fail fast and give the user guidance on what is wrong", func() {
			_, err := plgn.GetProduct([]string{"appname"}, []byte(``), nil)
			Ω(err).Should(HaveOccurred())
		})
	})

	Context("when calling GetProduct w/ valid flags and matching cloud config", func() {
		var deployment *enaml.DeploymentManifest
		var controlInstances = "1"
		var controlNetName = "private"
		var controlPass = "pss"
		var controlDisk = "4033"
		var controlVM = "medium"
		var controlIP = "1.2.3.4"

		BeforeEach(func() {
			cloudConfigBytes, _ := ioutil.ReadFile("./fixtures/sample-aws.yml")
			dmBytes, err := plgn.GetProduct([]string{
				"appname",
				"--disk-size", controlDisk,
				"--leader-instances", controlInstances,
				"--network-name", controlNetName,
				"--redis-pass", controlPass,
				"--vm-size", controlVM,
				"--leader-ip", controlIP,
				"--slave-ip", controlIP,
				"--stemcell-ver", "12.3.44",
			}, cloudConfigBytes, nil)
			Ω(err).ShouldNot(HaveOccurred())
			deployment = enaml.NewDeploymentManifest(dmBytes)
		})
		It("then we should have a properly initialized deployment set", func() {
			Ω(deployment.Update).ShouldNot(BeNil())
			Ω(len(deployment.Releases)).Should(Equal(1))
			Ω(len(deployment.Stemcells)).Should(Equal(1))
			Ω(len(deployment.Jobs)).Should(Equal(4))
		})
	})

	Context("when calling the plugin", func() {
		var flags []cli.Flag

		BeforeEach(func() {
			flags = pluginutil.ToCliFlagArray(plgn.GetFlags())
		})
		It("then there should be valid flags available", func() {
			for _, flagname := range []string{
				"leader-ip",
				"leader-instances",
				"redis-pass",
				"pool-instances",
				"disk-size",
				"slave-instances",
				"slave-ip",
				"network-name",
				"vm-size",
				"errand-instances",
				"stemcell-url",
				"stemcell-ver",
				"stemcell-sha",
				"stemcell-name",
			} {
				Ω(checkFlags(flags, flagname)).ShouldNot(HaveOccurred())
			}
		})
	})
})

func checkFlags(flags []cli.Flag, flagName string) error {
	var err = fmt.Errorf("could not find an ip flag %s in plugin", flagName)
	for _, f := range flags {
		if f.Names()[0] == flagName {
			err = nil
		}
	}
	return err
}
