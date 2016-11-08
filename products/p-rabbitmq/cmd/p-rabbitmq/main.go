package main

import (
	prabbitmq "github.com/enaml-ops/omg-product-bundle/products/p-rabbitmq/plugin"
	"github.com/enaml-ops/pluginlib/productv1"
)

// Version is the version of the p-rabbitmq plugin.
var Version string = "v0.0.0" // overridden at link time

func main() {
	product.Run(&prabbitmq.Plugin{
		Version: Version,
	})
}
