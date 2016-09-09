package pmysql_test

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/urfave/cli.v2"
	"github.com/enaml-ops/enaml"
	. "github.com/enaml-ops/omg-product-bundle/products/p-mysql/plugin"
	"github.com/enaml-ops/pluginlib/util"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/xchapter7x/lo"
	"github.com/xchapter7x/lo/lofakes"
)

var _ = Describe("given p-mysql Plugin", func() {
	var plgn *Plugin

	BeforeEach(func() {
		plgn = new(Plugin)
	})

	Context("when calling GetProduct while targeting an un-compatible cloud config'd bosh", func() {
		var logHolder = lo.G
		var logfake = new(lofakes.FakeLogger)
		var cloudConfigBytes []byte
		var controlNetName = "hello"
		var controlDisk = "medium"
		var controlVM = "large"
		var controlIP = "1.2.3.4"

		BeforeEach(func() {
			logfake = new(lofakes.FakeLogger)
			logHolder = lo.G
			lo.G = logfake
			cloudConfigBytes, _ = ioutil.ReadFile("./fixtures/sample-aws.yml")
		})

		AfterEach(func() {
			lo.G = logHolder
		})

		It("then we should fail fast and give the user guidance on what is wrong", func() {
			plgn.GetProduct([]string{
				"appname",
				"--disk-type", controlDisk,
				"--network", controlNetName,
				"--vm-type", controlVM,
				"--ip", controlIP,
				"--az", "z1-nothere",
				"--stemcell-url", "something",
				"--stemcell-ver", "12.3.44",
				"--stemcell-sha", "ilkjag09dhsg90ahsd09gsadg9",
			}, cloudConfigBytes)
			Ω(logfake.FatalCallCount()).Should(Equal(1))
		})
	})

	Context("when calling plugin without all required flags", func() {

		var logHolder = lo.G
		var logfake = new(lofakes.FakeLogger)

		BeforeEach(func() {
			logfake = new(lofakes.FakeLogger)
			logHolder = lo.G
			lo.G = logfake
		})

		AfterEach(func() {
			lo.G = logHolder
		})

		It("then it should fail fast and give the user guidance on what is wrong", func() {
			plgn.GetProduct([]string{"appname"}, []byte(``))
			Ω(logfake.FatalCallCount()).Should(BeNumerically(">=", 1))
		})
	})

	Context("when calling GetProduct w/ valid flags and matching cloud config", func() {
		var deployment *enaml.DeploymentManifest
		var controlNetName = "private"
		var controlDisk = "medium"
		var controlVM = "medium"
		var controlIP = "1.2.3.4"
		var controlInstanceGroups = []string{
			"mysql-partition",
			"acceptance-tests",
			"rejoin-unsafe",
			"broker-deregistrar",
			"broker-registrar",
			"cf-mysql-broker-partition",
			"monitoring-partition",
			"proxy-partition",
			"broker-deregistrar",
			"broker-registrar",
			"cf-mysql-broker-partition",
			"monitoring-partition",
			"proxy-partition",
			"backup-prepare-partition",
		}

		BeforeEach(func() {
			cloudConfigBytes, _ := ioutil.ReadFile("./fixtures/sample-aws.yml")
			dmBytes := plgn.GetProduct([]string{
				"appname",
				"--network", controlNetName,
				"--vm-type", controlVM,
				"--disk-type", controlDisk,
				"--ip", controlIP,
				"--az", "z1",
				"--stemcell-ver", "12.3.44",
			}, cloudConfigBytes)
			deployment = enaml.NewDeploymentManifest(dmBytes)
		})
		XIt("then we should have a properly initialized deployment set", func() {
			Ω(deployment.Update).ShouldNot(BeNil())
			Ω(len(deployment.Releases)).Should(Equal(4), "invalid number of releases")
			Ω(len(deployment.Stemcells)).Should(Equal(1), "invalid number of stemcells")

			for _, groupname := range controlInstanceGroups {
				Ω(checkGroupExists(deployment.InstanceGroups, groupname)).Should(BeTrue(), fmt.Sprintf("we expect instance group for %s", groupname))
			}
		})
	})

	Context("when calling GetProduct w/ a stemcell name flag ", func() {
		var deployment *enaml.DeploymentManifest
		var controlNetName = "private"
		var controlDisk = "medium"
		var controlVM = "medium"
		var controlIP = "1.2.3.4"

		BeforeEach(func() {
			cloudConfigBytes, _ := ioutil.ReadFile("./fixtures/sample-aws.yml")
			dmBytes := plgn.GetProduct([]string{
				"appname",
				"--network", controlNetName,
				"--vm-type", controlVM,
				"--disk-type", controlDisk,
				"--ip", controlIP,
				"--az", "z1",
				"--stemcell-name", "blahname",
				"--stemcell-url", "something",
				"--stemcell-ver", "12.3.44",
				"--stemcell-sha", "ilkjag09dhsg90ahsd09gsadg9",
			}, cloudConfigBytes)
			deployment = enaml.NewDeploymentManifest(dmBytes)
		})
		It("then we should have a properly configured stemcell definition in our deployment (os & alias from flag value)", func() {
			Ω(len(deployment.Stemcells)).Should(Equal(1))
			Ω(deployment.Stemcells[0].OS).Should(Equal("blahname"))
			Ω(deployment.Stemcells[0].Alias).Should(Equal("blahname"))
		})
	})

	Context("when calling the plugin", func() {
		var flags []cli.Flag

		BeforeEach(func() {
			flags = pluginutil.ToCliFlagArray(plgn.GetFlags())
		})
		It("then there should be valid flags available", func() {
			for _, flagname := range []string{
				"ip",
				"az",
				"network",
				"vm-type",
				"disk-type",
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
	var err = fmt.Errorf("could not find an flag %s in plugin", flagName)
	for _, f := range flags {
		if f.Names()[0] == flagName {
			err = nil
		}
	}
	return err
}
