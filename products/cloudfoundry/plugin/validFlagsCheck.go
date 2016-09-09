package cloudfoundry

import (
	"gopkg.in/urfave/cli.v2"
	"github.com/xchapter7x/lo"
)

func hasValidStringFlags(c *cli.Context, flaglist []string) bool {
	for _, v := range flaglist {

		if c.String(v) == "" {
			lo.G.Errorf("empty flag value for required field: %v", v)
			return false
		}
	}
	return true
}

func hasValidStringSliceFlags(c *cli.Context, flaglist []string) bool {

	for _, v := range flaglist {
		lo.G.Debug(c.StringSlice(v))
		if len(c.StringSlice(v)) > 0 {
			lo.G.Errorf("empty flag value for required field: %v", v)
			return false
		}
	}
	return true
}
