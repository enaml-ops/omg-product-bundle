package cloudfoundry

import (
	"fmt"
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

	Describe("given a call to the plugin", func() {

		Context("when 'stemcell-name' flag is given by the user", func() {
			var ertPlugin *Plugin

			BeforeEach(func() {
				ertPlugin = new(Plugin)
			})

			It("then it should overwrite the default and use the value given", func() {
				controlStemcellAlias := "ubuntu-magic"
				manifestBytes, err := ertPlugin.GetProduct(append(
					[]string{"ert-command"},
					ertRequiredFlags(controlStemcellAlias)...,
				), []byte{})
				Expect(err).ShouldNot(HaveOccurred())
				manifest := enaml.NewDeploymentManifest(manifestBytes)
				Ω(manifest.Stemcells).ShouldNot(BeNil())
				Ω(manifest.Stemcells[0].Alias).Should(Equal(controlStemcellAlias))
				for _, instanceGroup := range manifest.InstanceGroups {
					Expect(instanceGroup.Stemcell).Should(Equal(controlStemcellAlias), fmt.Sprintf("stemcell for instance group %v was not set properly", instanceGroup.Name))
				}
			})
		})
	})

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

func ertRequiredFlags(stemcellAlias string) []string {
	return []string{
		"--stemcell-name", stemcellAlias,
		"--vault-active=false",
		"--network", "1",
		"--system-domain", "1",
		"--host-key-fingerprint", "1",
		"--support-address", "1",
		"--min-cli-version", "1",
		"--nfs-share-path", "1",
		"--uaa-login-protocol", "1",
		"--doppler-zone", "1",
		"--uaa-company-name", "1",
		"--uaa-product-logo", "1",
		"--uaa-square-logo", "1",
		"--uaa-footer-legal-txt", "1",
		"--az", "1",
		"--app-domain", "1",
		"--nfs-allow-from-network-cidr", "1",
		"--nats-port", "1",
		"--doppler-drain-buffer-size", "1",
		"--cc-uploader-poll-interval", "1",
		"--cc-external-port", "1",
		"--loggregator-port", "1",
		"--consul-agent-cert", "1",
		"--consul-agent-key", "1",
		"--consul-server-cert", "1",
		"--consul-server-key", "1",
		"--bbs-server-ca-cert", "1",
		"--bbs-client-cert", "1",
		"--bbs-client-key", "1",
		"--bbs-server-cert", "1",
		"--bbs-server-key", "1",
		"--etcd-server-cert", "1",
		"--etcd-server-key", "1",
		"--etcd-client-cert", "1",
		"--etcd-client-key", "1",
		"--etcd-peer-cert", "1",
		"--etcd-peer-key", "1",
		"--uaa-saml-service-provider-key", "1",
		"--uaa-saml-service-provider-cert", "1",
		"--uaa-jwt-signing-key", "1",
		"--uaa-jwt-verification-key", "1",
		"--router-ssl-cert", "1",
		"--router-ssl-key", "1",
		"--diego-cell-disk-type", "1",
		"--diego-brain-disk-type", "1",
		"--diego-db-disk-type", "1",
		"--nfs-disk-type", "1",
		"--etcd-disk-type", "1",
		"--mysql-disk-type", "1",
		"--cc-instances", "1",
		"--uaa-instances", "1",
		"--cc-worker-instances", "1",
		"--db-uaa-password", "1",
		"--push-apps-manager-password", "1",
		"--system-services-password", "1",
		"--system-verification-password", "1",
		"--opentsdb-firehose-nozzle-client-secret", "1",
		"--identity-client-secret", "1",
		"--login-client-secret", "1",
		"--portal-client-secret", "1",
		"--autoscaling-service-client-secret", "1",
		"--system-passwords-client-secret", "1",
		"--cc-service-dashboards-client-secret", "1",
		"--gorouter-client-secret", "1",
		"--notifications-client-secret", "1",
		"--notifications-ui-client-secret", "1",
		"--cloud-controller-username-lookup-client-secret", "1",
		"--cc-routing-client-secret", "1",
		"--apps-metrics-client-secret", "1",
		"--apps-metrics-processing-client-secret", "1",
		"--admin-password", "1",
		"--nats-pass", "1",
		"--mysql-bootstrap-password", "1",
		"--consul-encryption-key", "1",
		"--smoke-tests-password", "1",
		"--doppler-shared-secret", "1",
		"--doppler-client-secret", "1",
		"--cc-bulk-api-password", "1",
		"--cc-internal-api-password", "1",
		"--ssh-proxy-uaa-secret", "1",
		"--cc-db-encryption-key", "1",
		"--db-ccdb-password", "1",
		"--diego-db-passphrase", "1",
		"--uaa-admin-secret", "1",
		"--router-pass", "1",
		"--mysql-proxy-api-password", "1",
		"--mysql-admin-password", "1",
		"--db-console-password", "1",
		"--cc-staging-upload-password", "1",
		"--db-app_usage-password", "1",
		"--apps-manager-secret-token", "1",
		"--db-autoscale-password", "1",
		"--db-notifications-password", "1",
		"--nats-user", "1",
		"--mysql-bootstrap-username", "1",
		"--cc-staging-upload-user", "1",
		"--db-ccdb-username", "1",
		"--db-uaa-username", "1",
		"--mysql-proxy-api-username", "1",
		"--db-console-username", "1",
		"--router-user", "1",
		"--cc-internal-api-user", "1",
		"--db-autoscale-username", "1",
		"--db-notifications-username", "1",
		"--mysql-proxy-vm-type", "1",
		"--clock-global-vm-type", "1",
		"--cc-vm-type", "1",
		"--diego-brain-vm-type", "1",
		"--diego-cell-vm-type", "1",
		"--doppler-vm-type", "1",
		"--loggregator-traffic-controller-vmtype", "1",
		"--cc-worker-vm-type", "1",
		"--errand-vm-type", "1",
		"--etcd-vm-type", "1",
		"--nats-vm-type", "1",
		"--consul-vm-type", "1",
		"--mysql-vm-type", "1",
		"--diego-db-vm-type", "1",
		"--uaa-vm-type", "1",
		"--router-vm-type", "1",
		"--nfs-vm-type", "1",
		"--nfs-ip", "1",
		"--mysql-ip", "1",
		"--diego-cell-ip", "1",
		"--consul-ip", "1",
		"--doppler-ip", "1",
		"--mysql-proxy-ip", "1",
		"--loggregator-traffic-controller-ip", "1",
		"--nats-machine-ip", "1",
		"--etcd-machine-ip", "1",
		"--diego-brain-ip", "1",
		"--diego-db-ip", "1",
		"--router-ip", "1",
	}
}
