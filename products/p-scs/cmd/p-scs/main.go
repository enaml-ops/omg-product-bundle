package main

import (
	pscs "github.com/enaml-ops/omg-product-bundle/products/p-scs/plugin"
	"github.com/enaml-ops/pluginlib/productv1"
)

// Version is the version of the p-rabbitmq plugin.
var Version string = "v0.0.0" // overridden at link time

func main() {
	product.Run(&pscs.Plugin{
		Version: Version,
	})
}
