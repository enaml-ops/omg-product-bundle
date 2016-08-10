package main

import (
	"github.com/enaml-ops/omg-product-bundle/products/dockerregistry/plugin"
	"github.com/enaml-ops/pluginlib/product"
)

func main() {
	product.Run(new(plugin.Plugin))
}
