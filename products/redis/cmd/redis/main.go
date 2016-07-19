package main

import (
	"github.com/enaml-ops/omg-cli/pluginlib/product"
	"github.com/enaml-ops/omg-product-bundle/products/redis/plugin"
)

func main() {
	product.Run(new(redis.Plugin))
}
