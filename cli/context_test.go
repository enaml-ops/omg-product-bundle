package cli_test

import (
	"flag"
	"io/ioutil"

	cli "gopkg.in/urfave/cli.v2"

	. "github.com/enaml-ops/omg-product-bundle/cli"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("LoadResourceFromContext function", func() {
	Context("when called with a filename (prefixed by @)", func() {
		var ctx *cli.Context

		BeforeEach(func() {
			set := flag.NewFlagSet("test", 0)
			set.String("my-flag", "@fixtures/foo.txt", "")
			ctx = cli.NewContext(nil, set, nil)
		})

		It("reads from the specified file", func() {
			value, err := LoadResourceFromContext(ctx, "my-flag")
			Ω(err).Should(BeNil())

			exp, _ := ioutil.ReadFile("fixtures/foo.txt")
			Ω(value).Should(Equal(string(exp)))
		})
	})

	Context("when called with a standard string argument", func() {
		var ctx *cli.Context

		BeforeEach(func() {
			set := flag.NewFlagSet("test", 0)
			set.String("my-flag", "fixtures/deployment_task.json", "")
			ctx = cli.NewContext(nil, set, nil)
		})

		It("returns the argument value directly", func() {
			value, err := LoadResourceFromContext(ctx, "my-flag")
			Ω(err).Should(BeNil())
			Ω(value).Should(Equal("fixtures/deployment_task.json"))
		})
	})
})
