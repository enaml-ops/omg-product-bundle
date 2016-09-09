package main

import (
	pmysql "github.com/enaml-ops/omg-product-bundle/products/p-mysql/plugin"
	"github.com/enaml-ops/pluginlib/product"
)

// Version is the version of the p-rabbitmq plugin.
var Version string = "v0.0.0" // overridden at link time

func main() {
	product.Run(&pmysql.Plugin{
		PluginVersion: Version,
	})
}
