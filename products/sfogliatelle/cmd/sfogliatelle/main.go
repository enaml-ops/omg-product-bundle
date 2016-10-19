package main

import (
	"os"

	"github.com/enaml-ops/omg-product-bundle/products/sfogliatelle/plugin"
	"github.com/enaml-ops/pluginlib/product"
)

// Version is the version of the sfogliatelle plugin.
var Version string = "v0.0.0" // overridden at link time

func main() {
	product.Run(&sfogliatelle.Plugin{
		Version: Version,
		Source:  os.Stdin,
	})
}
