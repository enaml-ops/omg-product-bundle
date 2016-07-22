package cloudfoundry_test

import (
	"io/ioutil"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"

	. "github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/plugin"
	"github.com/enaml-ops/pluginlib/pcli"
	"github.com/enaml-ops/pluginlib/util"
)

var _ = Describe("Cloud Foundry Plugin", func() {

	Describe("given InferFromCloudDecorate", func() {
		Context("when infer-from-cloud is set to true", func() {
			var flgs []pcli.Flag

			BeforeEach(func() {
				b, _ := ioutil.ReadFile("fixtures/cloudconfig.yml")
				var flagsToInferFromCloudConfig = map[string][]string{
					"disktype": []string{"mysql-disk-type"},
					"vmtype":   []string{"diego-brain-vm-type"},
					"az":       []string{"az"},
					"network":  []string{"network"},
				}
				flgs = []pcli.Flag{
					pcli.Flag{FlagType: pcli.BoolFlag, Name: "infer-from-cloud"},
					pcli.Flag{FlagType: pcli.StringFlag, Name: "diego-brain-vm-type"},
					pcli.Flag{FlagType: pcli.StringFlag, Name: "mysql-disk-type"},
					pcli.Flag{FlagType: pcli.StringFlag, Name: "az"},
					pcli.Flag{FlagType: pcli.StringFlag, Name: "network"},
				}
				args := []string{
					"mycoolness",
					"--infer-from-cloud",
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
	})

	Describe("given the GetProduct Method", func() {
		var plugin *Plugin

		BeforeEach(func() {
			plugin = new(Plugin)
		})

		Context("when called w/ a `vault-active` flag set to TRUE and an INCOMPLETE set of vault values", func() {
			It("then it should panic", func() {
				Ω(func() {
					plugin.GetProduct([]string{
						"my-app",
						"--vault-active",
					},
						[]byte(``),
					)
				}).Should(Panic())
			})
		})
	})
})
