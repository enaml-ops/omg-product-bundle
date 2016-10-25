package cloudfoundry_test

import (
	"io/ioutil"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
	"gopkg.in/urfave/cli.v2"

	. "github.com/enaml-ops/omg-product-bundle/products/oss_cf/plugin"
	"github.com/enaml-ops/omg-product-bundle/products/oss_cf/plugin/config"
	"github.com/enaml-ops/omg-product-bundle/products/oss_cf/plugin/pluginfakes"
	"github.com/enaml-ops/pluginlib/pluginutil"
)

var _ = Describe("Vault helpers", func() {
	Describe("given a RotatePasswordHash", func() {
		Context("when called with a vaultrotater and a valid hash", func() {
			var fakeVault *pluginfakes.FakeVaultRotater
			var err error
			BeforeEach(func() {
				fakeVault = new(pluginfakes.FakeVaultRotater)
				fakeVault.RotateSecretsReturns(nil)
				err = RotatePasswordHash(fakeVault, "secret/hash/of/stuff")
			})
			It("should set a valid set of secrets to vault", func() {
				_, givenSecrets := fakeVault.RotateSecretsArgsForCall(0)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(givenSecrets).ShouldNot(BeEmpty())
			})
		})
	})
	Describe("given a RotateCertHash", func() {
		Context("when called with a vaultrotater and a valid hash", func() {
			var fakeVault *pluginfakes.FakeVaultRotater
			var err error
			BeforeEach(func() {
				fakeVault = new(pluginfakes.FakeVaultRotater)
				fakeVault.RotateSecretsReturns(nil)
				err = RotateCertHash(fakeVault, "secret/hash/of/stuff", "sys.fake.domain.io", []string{"apps.fake.domain.io"})
			})
			It("should set a valid set of secrets to vault", func() {
				_, givenSecrets := fakeVault.RotateSecretsArgsForCall(0)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(givenSecrets).ShouldNot(BeEmpty())
			})
		})
	})
	XDescribe("given a VaultRotate", func() {
		Context("when called with the rotate flag and all required vault flags", func() {
			It("then it should rotate the password values in the vault", func() {
				Ω(true).Should(BeFalse())
			})

			It("then it should rotate the keycert values in the vault", func() {
				Ω(true).Should(BeFalse())
			})
		})
	})
	Describe("given vault unmarshal", func() {
		Context("when decorating with vault data", func() {
			var (
				server *ghttp.Server
				c      *cli.Context
			)

			BeforeEach(func() {
				b, _ := ioutil.ReadFile("fixtures/mysql_vault.json")
				server = ghttp.NewServer()
				server.AllowUnhandledRequests = true
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", "/v1/secret/pcf-np-1-password"),
						ghttp.RespondWith(http.StatusOK, b),
					),
				)

				p := new(Plugin)
				flags := p.GetFlags()
				args := []string{
					"app",
					"--system-domain", "sys.example.com",
					"--app-domain", "apps.example.com",
					"--app-domain", "apps2.example.com",
					"--vault-domain", server.URL(),
					"--vault-hash-password", "secret/pcf-np-1-password",
					"--vault-hash-keycert", "secret/pcf-np-1-keycert",
					"--vault-hash-host", "secret/pcf-np-1-hostname",
					"--vault-hash-ip", "secret/pcf-np-1-ips",
					"--vault-token", "dasdfasdf",
					"--vault-rotate",
					"--syslog-address", "10.113.82.164",
				}
				VaultDecorate(args, flags)
				c = pluginutil.NewContext(args, pluginutil.ToCliFlagArray(flags))
			})

			AfterEach(func() {
				server.Close()
			})

			It("generates seeded databases", func() {
				Ω(c.String("db-uaa-password")).ShouldNot(BeEmpty())
				ig := NewMySQLPartition(&config.Config{})
				mysql := ig.(*MySQL)
				Ω(mysql.MySQLSeededDatabases).ShouldNot(BeEmpty())
			})

			It("uses string slice values specified as args and ignores values in vault", func() {
				appDomains := c.StringSlice("app-domain")
				Ω(appDomains).Should(ConsistOf("apps.example.com", "apps2.example.com"))
			})
		})
	})
})
