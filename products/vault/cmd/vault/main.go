package main

import (
	"github.com/enaml-ops/omg-cli/pluginlib/product"
	"github.com/enaml-ops/omg-product-bundle/products/vault/plugin"
)

func main() {
	product.Run(new(vault.Plugin))
}
