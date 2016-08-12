package main

import (
	"github.com/enaml-ops/omg-product-bundle/products/docker/plugin"
	"github.com/enaml-ops/pluginlib/product"
)

func main() {
	product.Run(new(docker.Plugin))
}
