package main

import (
	"github.com/enaml-ops/omg-product-bundle/products/concourse/plugin"
	"github.com/enaml-ops/omg-cli/pluginlib/product"
)

func main() {
	product.Run(new(concourseplugin.ConcoursePlugin))
}
