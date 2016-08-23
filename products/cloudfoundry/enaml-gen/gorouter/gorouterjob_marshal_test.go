package gorouter_test

import (
	"github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/gorouter"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	yaml "gopkg.in/yaml.v2"
)

const controlYAML = `routing-api:
  port: 1234
  auth_disabled: true
routing_api:
  enabled: true
`

var _ = Describe("given GorouterJob", func() {
	Context("when marshalling routing API", func() {
		It("produces the correct YAML", func() {
			j := gorouter.GorouterJob{
				RoutingApi: &gorouter.RoutingApi{
					Port:         1234,
					AuthDisabled: true,
					Enabled:      true,
				},
			}
			b, err := yaml.Marshal(&j)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(b).Should(MatchYAML(controlYAML))
		})
	})
})
