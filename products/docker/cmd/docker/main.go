package main

import (
	"github.com/enaml-ops/omg-product-bundle/products/docker/plugin"
	"github.com/enaml-ops/pluginlib/product"
)

var Version string = "v0.0.0"

func main() {
	product.Run(&docker.Plugin{
		PluginVersion: Version,
	})
}
