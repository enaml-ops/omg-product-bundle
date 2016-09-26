package cloudfoundry

import (
	"io/ioutil"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
	"github.com/xchapter7x/lo"
	"github.com/xchapter7x/lo/lofakes"

	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/pluginlib/pcli"
	"github.com/enaml-ops/pluginlib/pluginutil"
)

var _ = Describe("Cloud Foundry Plugin", func() {

	Describe("given InferFromCloudDecorate", func() {
		Context("when infer-from-cloud is set to true", func() {
			var flgs []pcli.Flag
			var args []string

			BeforeEach(func() {
				b, _ := ioutil.ReadFile("fixtures/cloudconfig.yml")
				var flagsToInferFromCloudConfig = map[string][]string{
					"disktype": []string{"mysql-disk-type", "nfs-disk-type"},
					"vmtype":   []string{"diego-brain-vm-type"},
					"az":       []string{"az"},
					"network":  []string{"network"},
				}
				flgs = []pcli.Flag{
					pcli.Flag{FlagType: pcli.BoolFlag, Name: "infer-from-cloud"},
					pcli.Flag{FlagType: pcli.StringFlag, Name: "diego-brain-vm-type"},
					pcli.Flag{FlagType: pcli.StringFlag, Name: "mysql-disk-type"},
					pcli.Flag{FlagType: pcli.StringFlag, Name: "nfs-disk-type"},
					pcli.Flag{FlagType: pcli.StringFlag, Name: "az"},
					pcli.Flag{FlagType: pcli.StringFlag, Name: "network"},
				}
				args = []string{
					"mycoolness",
					"--infer-from-cloud",
					"--nfs-disk-type", "large",
				}
				InferFromCloudDecorate(flagsToInferFromCloudConfig, b, args, flgs)
			})

			It("then it should decorate the given flag array with cloudconfig values as defaults", func() {
				ctx := pluginutil.NewContext([]string{"mycoolapp"}, pluginutil.ToCliFlagArray(flgs))
				Ω(ctx.String("diego-brain-vm-type")).Should(Equal("smallvm"))
				Ω(ctx.String("mysql-disk-type")).Should(Equal("smalldisk"))
				Ω(ctx.String("az")).Should(Equal("z1,z2"))
				Ω(ctx.String("network")).Should(Equal("private"))
			})

			It("then it should not override flags that were manually provided", func() {
				ctx := pluginutil.NewContext(args, pluginutil.ToCliFlagArray(flgs))
				Ω(ctx.String("nfs-disk-type")).Should(Equal("large"))
			})
		})
	})

	Describe("given VaultDecorate", func() {
		Context("when called with a set of args and flags that can be overwritten from a vault", func() {
			var server *ghttp.Server
			var flgs []pcli.Flag

			BeforeEach(func() {
				b, _ := ioutil.ReadFile("fixtures/vault.json")
				server = ghttp.NewServer()
				server.AllowUnhandledRequests = true
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", "/v1/secret/move-along-nothing-to-see-here", ""),
						ghttp.RespondWith(http.StatusOK, string(b)),
					),
				)
				flgs = []pcli.Flag{
					pcli.Flag{FlagType: pcli.StringFlag, Name: "knock"},
					pcli.Flag{FlagType: pcli.BoolTFlag, Name: "vault-active"},
					pcli.Flag{FlagType: pcli.StringFlag, Name: "vault-domain"},
					pcli.Flag{FlagType: pcli.StringFlag, Name: "vault-token"},
					pcli.Flag{FlagType: pcli.StringFlag, Name: "vault-hash-password"},
					pcli.Flag{FlagType: pcli.StringFlag, Name: "vault-hash-keycert"},
					pcli.Flag{FlagType: pcli.StringFlag, Name: "vault-hash-ip"},
					pcli.Flag{FlagType: pcli.StringFlag, Name: "vault-hash-host"},
				}
				args := []string{
					"mycoolness",
					"--vault-token", "lshdglkahsdlgkhaskldghalsdhgk",
					"--vault-domain", server.URL(),
					"--vault-hash-password", "secret/move-along-nothing-to-see-here",
					"--vault-hash-keycert", "secret/move-along-nothing-to-see-here",
					"--vault-hash-ip", "secret/move-along-nothing-to-see-here",
					"--vault-hash-host", "secret/move-along-nothing-to-see-here",
				}
				VaultDecorate(args, flgs)
			})

			AfterEach(func() {
				server.Close()
			})

			It("then it should decorate the given flag array with vault values as defaults", func() {
				ctx := pluginutil.NewContext([]string{"mycoolapp"}, pluginutil.ToCliFlagArray(flgs))
				Ω(ctx.String("knock")).Should(Equal("knocks"))
			})
		})

		Context("when called w/ a `vault-active` flag set to TRUE and an INCOMPLETE set of vault values", func() {
			var logHolder = lo.G
			var logfake = new(lofakes.FakeLogger)

			BeforeEach(func() {
				logfake = new(lofakes.FakeLogger)
				logHolder = lo.G
				lo.G = logfake
				flgs := []pcli.Flag{
					pcli.Flag{FlagType: pcli.StringFlag, Name: "knock"},
					pcli.Flag{FlagType: pcli.BoolTFlag, Name: "vault-active"},
					pcli.Flag{FlagType: pcli.StringFlag, Name: "vault-domain"},
					pcli.Flag{FlagType: pcli.StringFlag, Name: "vault-token"},
					pcli.Flag{FlagType: pcli.StringFlag, Name: "vault-hash-password"},
					pcli.Flag{FlagType: pcli.StringFlag, Name: "vault-hash-keycert"},
					pcli.Flag{FlagType: pcli.StringFlag, Name: "vault-hash-ip"},
					pcli.Flag{FlagType: pcli.StringFlag, Name: "vault-hash-host"},
				}
				args := []string{
					"mycoolness",
					"--vault-token", "lshdglkahsdlgkhaskldghalsdhgk",
					"--vault-hash-ip", "secret/move-along-nothing-to-see-here",
					"--vault-hash-host", "secret/move-along-nothing-to-see-here",
				}
				VaultDecorate(args, flgs)
			})

			AfterEach(func() {
				lo.G = logHolder
			})

			It("then it exit and print an error", func() {
				Ω(logfake.FatalCallCount()).Should(Equal(1))
			})
		})
	})

	Describe("given getDeploymentManifest", func() {
		Context("when called without any instance groupers", func() {
			var p *Plugin
			var dm *enaml.DeploymentManifest
			var err error
			var oldfactories []InstanceGrouperFactory

			BeforeEach(func() {
				oldfactories = factories
				factories = factories[:0]
				p = new(Plugin)
				ctx := pluginutil.NewContext([]string{"cloudfoundry"}, pluginutil.ToCliFlagArray(p.GetFlags()))
				dm, err = p.getDeploymentManifest(ctx, nil)
				Ω(err).ShouldNot(HaveOccurred())
			})

			AfterEach(func() {
				factories = oldfactories
			})

			It("includes the correct releases", func() {
				hasRelease := func(name string) bool {
					for i := range dm.Releases {
						if dm.Releases[i].Name == name {
							return true
						}
					}
					return false
				}
				Ω(hasRelease(CFReleaseName)).Should(BeTrue())
				Ω(hasRelease(CFMysqlReleaseName)).Should(BeTrue())
				Ω(hasRelease(DiegoReleaseName)).Should(BeTrue())
				Ω(hasRelease(GardenReleaseName)).Should(BeTrue())
				Ω(hasRelease(CFLinuxReleaseName)).Should(BeTrue())
				Ω(hasRelease(EtcdReleaseName)).Should(BeTrue())
				Ω(hasRelease(MySQLBackupReleaseName)).Should(BeTrue())
				Ω(hasRelease(PushAppsReleaseName)).Should(BeTrue())
				Ω(hasRelease(CFAutoscalingReleaseName)).Should(BeTrue())
				// Ω(hasRelease(NotificationsReleaseName)).Should(BeTrue())
				// Ω(hasRelease(NotificationsUIReleaseName)).Should(BeTrue())
			})
		})
	})
})
