package docker_test

import (
	"fmt"
	"io/ioutil"

	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-product-bundle/products/docker/enaml-gen/docker"
	. "github.com/enaml-ops/omg-product-bundle/products/docker/plugin"
	"github.com/enaml-ops/pluginlib/pluginutil"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/xchapter7x/lo"
	"github.com/xchapter7x/lo/lofakes"
	"gopkg.in/urfave/cli.v2"
)

var _ = Describe("given docker Plugin", func() {
	var plgn *Plugin

	BeforeEach(func() {
		plgn = new(Plugin)
	})

	Context("when called with a `--insecure-registry` stringslice flag value/s given", func() {
		var deployment *enaml.DeploymentManifest
		var controlRegistry1 = "blah"
		var controlRegistry2 = "bleh"

		BeforeEach(func() {
			cloudConfigBytes, _ := ioutil.ReadFile("./fixtures/sample-aws.yml")
			dmBytes, err := plgn.GetProduct([]string{
				"appname",
				"--network", "private",
				"--vm-type", "medium",
				"--disk-type", "medium",
				"--ip", "1.2.3.4",
				"--az", "z1",
				"--registry-mirror", "my.mirror.com",
				"--stemcell-ver", "12.3.44",
				"--stemcell-sha", "abcdef",
				"--stemcell-url", "https://stemcells.com/foo",
				"--container-definition", "./fixtures/sample-docker.yml",
				"--insecure-registry", controlRegistry1,
				"--insecure-registry", controlRegistry2,
			}, cloudConfigBytes, nil)
			deployment = enaml.NewDeploymentManifest(dmBytes)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(len(deployment.InstanceGroups)).Should(BeNumerically(">", 0), "we expect there to be some instance groups defined")
		})

		It("then it should properly pass the flag value to the plugin", func() {
			Ω(plgn.InsecureRegistries).Should(ConsistOf(controlRegistry1, controlRegistry2), "there should be insecure registries in the job properties")
		})
	})

	Context("when called with a `--docker-release-ver` `--docker-release-url` `--docker-release-sha` flag", func() {

		It("then it should have those registered as valid flags", func() {
			cloudConfigBytes, _ := ioutil.ReadFile("./fixtures/sample-aws.yml")
			Ω(func() {
				plgn.GetProduct([]string{
					"appname",
					"--network", "private",
					"--vm-type", "medium",
					"--disk-type", "medium",
					"--ip", "1.2.3.4",
					"--az", "z1",
					"--stemcell-ver", "12.3.44",
					"--container-definition", "./fixtures/sample-docker.yml",
					"--docker-release-ver", "skjdf",
					"--docker-release-url", "asdfasdf",
					"--docker-release-sha", "asdfasdf",
				}, cloudConfigBytes, nil)
			}).ShouldNot(Panic(), "these flags should not cause a panic, b/c they should exist")
		})

		It("then it should set the give values as the release values", func() {
			cloudConfigBytes, _ := ioutil.ReadFile("./fixtures/sample-aws.yml")
			controlver := "asdfasdf"
			controlurl := "fasdfasdf"
			controlsha := "akjhasdkghasdg"
			dmBytes, err := plgn.GetProduct([]string{
				"appname",
				"--network", "private",
				"--vm-type", "medium",
				"--disk-type", "medium",
				"--ip", "1.2.3.4",
				"--az", "z1",
				"--stemcell-ver", "12.3.44",
				"--stemcell-url", "https://stemcells.com/foo",
				"--stemcell-sha", "abcdef",
				"--container-definition", "./fixtures/sample-docker.yml",
				"--docker-release-ver", controlver,
				"--docker-release-url", controlurl,
				"--docker-release-sha", controlsha,
			}, cloudConfigBytes, nil)
			deployment := enaml.NewDeploymentManifest(dmBytes)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(deployment.Releases[0].Version).Should(Equal(controlver))
			Ω(deployment.Releases[0].URL).Should(Equal(controlurl))
			Ω(deployment.Releases[0].SHA1).Should(Equal(controlsha))
		})
	})

	Context("when called with a `--registry-mirror` stringslice flag value/s given", func() {
		var deployment *enaml.DeploymentManifest
		var controlMirror1 = "blah"
		var controlMirror2 = "bleh"

		BeforeEach(func() {
			cloudConfigBytes, _ := ioutil.ReadFile("./fixtures/sample-aws.yml")
			dmBytes, err := plgn.GetProduct([]string{
				"appname",
				"--network", "private",
				"--vm-type", "medium",
				"--disk-type", "medium",
				"--ip", "1.2.3.4",
				"--az", "z1",
				"--stemcell-ver", "12.3.44",
				"--stemcell-url", "https://stemcells.com/foo",
				"--stemcell-sha", "abcdef",
				"--container-definition", "./fixtures/sample-docker.yml",
				"--registry-mirror", controlMirror1,
				"--registry-mirror", controlMirror2,
			}, cloudConfigBytes, nil)
			deployment = enaml.NewDeploymentManifest(dmBytes)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(len(deployment.InstanceGroups)).Should(BeNumerically(">", 0), "we expect there to be some instance groups defined")
		})

		It("then it should properly pass the flag value to the plugin", func() {
			Ω(plgn.RegistryMirrors).Should(ConsistOf(controlMirror1, controlMirror2), "there should be registry mirrors in the job properties")
		})
	})

	Context("when the plugin has a InsecureRegistries value set", func() {
		var plgn *Plugin
		var ig *enaml.InstanceGroup
		var controlRegistry1 = "blah"
		var controlRegistry2 = "bleh"
		BeforeEach(func() {
			plgn = new(Plugin)
			plgn.InsecureRegistries = []string{controlRegistry1, controlRegistry2}
			ig = plgn.NewDockerInstanceGroup()
		})
		It("then it should set the insecure-registries array in the bosh deployment manifest the plugin generates", func() {
			var dockerJobProperties *docker.DockerJob = ig.GetJobByName("docker").Properties.(*docker.DockerJob)
			Ω(dockerJobProperties.Docker.InsecureRegistries).Should(ConsistOf(controlRegistry1, controlRegistry2), "there should be insecure registries in the job properties")
		})
	})

	Context("when the plugin has a RegistryMirrors value set", func() {
		var plgn *Plugin
		var ig *enaml.InstanceGroup
		var controlMirror1 = "blah"
		var controlMirror2 = "bleh"
		BeforeEach(func() {
			plgn = new(Plugin)
			plgn.RegistryMirrors = []string{controlMirror1, controlMirror2}
			ig = plgn.NewDockerInstanceGroup()
		})
		It("then it should set the insecure-registries array in the bosh deployment manifest the plugin generates", func() {
			var dockerJobProperties *docker.DockerJob = ig.GetJobByName("docker").Properties.(*docker.DockerJob)
			Ω(dockerJobProperties.Docker.RegistryMirrors).Should(ConsistOf(controlMirror1, controlMirror2), "there should be insecure registries in the job properties")
		})
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
			_, err := plgn.GetProduct([]string{
				"appname",
				"--disk-type", controlDisk,
				"--network", controlNetName,
				"--vm-type", controlVM,
				"--ip", controlIP,
				"--az", "z1-nothere",
				"--stemcell-url", "something",
				"--stemcell-ver", "12.3.44",
				"--stemcell-sha", "ilkjag09dhsg90ahsd09gsadg9",
				"--container-definition", "./fixtures/sample-docker.yml",
			}, cloudConfigBytes, nil)
			Ω(err).Should(HaveOccurred())
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
			_, err := plgn.GetProduct([]string{"appname"}, []byte(``), nil)
			Ω(err).Should(HaveOccurred())
		})
	})

	Context("when calling GetProduct without a valid docker def file", func() {
		var controlNetName = "private"
		var controlDisk = "medium"
		var controlVM = "medium"
		var controlIP = "1.2.3.4"
		var realLog = lo.G
		var logfake = new(lofakes.FakeLogger)

		BeforeEach(func() {

			realLog = lo.G
			logfake = new(lofakes.FakeLogger)
			lo.G = logfake
			cloudConfigBytes, _ := ioutil.ReadFile("./fixtures/sample-aws.yml")
			plgn.GetProduct([]string{
				"appname",
				"--network", controlNetName,
				"--vm-type", controlVM,
				"--disk-type", controlDisk,
				"--ip", controlIP,
				"--az", "z1",
				"--stemcell-ver", "12.3.44",
				"--stemcell-url", "https://stemcells.com/foo",
				"--stemcell-sha", "abcdef",
				"--container-definition", "this-file-does-not-exist",
			}, cloudConfigBytes, nil)
		})

		AfterEach(func() {
			lo.G = realLog
		})

		It("then we should have a properly initialized deployment set", func() {
			Ω(logfake.FatalfCallCount()).Should(Equal(1))
		})
	})

	Context("when calling GetProduct w/ valid flags and matching cloud config", func() {
		var deployment *enaml.DeploymentManifest
		var controlNetName = "private"
		var controlDisk = "medium"
		var controlVM = "medium"
		var controlIP = "1.2.3.4"

		BeforeEach(func() {
			cloudConfigBytes, _ := ioutil.ReadFile("./fixtures/sample-aws.yml")
			dmBytes, err := plgn.GetProduct([]string{
				"appname",
				"--network", controlNetName,
				"--vm-type", controlVM,
				"--disk-type", controlDisk,
				"--ip", controlIP,
				"--az", "z1",
				"--registry-mirror", "my.registry.com",
				"--stemcell-url", "http://stemcells.com/foo",
				"--stemcell-sha", "abcdef",
				"--stemcell-ver", "12.3.44",
				"--container-definition", "./fixtures/sample-docker.yml",
			}, cloudConfigBytes, nil)
			Ω(err).ShouldNot(HaveOccurred())
			deployment = enaml.NewDeploymentManifest(dmBytes)
		})
		It("then we should have a properly initialized deployment set", func() {
			Ω(deployment.Update).ShouldNot(BeNil())
			Ω(len(deployment.Releases)).Should(Equal(1))
			Ω(len(deployment.Stemcells)).Should(Equal(1))
			Ω(len(deployment.InstanceGroups)).Should(Equal(1))
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
			dmBytes, err := plgn.GetProduct([]string{
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
				"--container-definition", "./fixtures/sample-docker.yml",
			}, cloudConfigBytes, nil)
			Ω(err).ShouldNot(HaveOccurred())
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
				"container-definition",
			} {
				Ω(checkFlags(flags, flagname)).ShouldNot(HaveOccurred())
			}
		})
	})
})

func checkFlags(flags []cli.Flag, flagName string) error {
	var err = fmt.Errorf("could not find an flag %s in plugin", flagName)
	for _, f := range flags {
		if len(f.Names()) > 0 && f.Names()[0] == flagName {
			err = nil
		}
	}
	return err
}
