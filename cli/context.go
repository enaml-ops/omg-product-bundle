package cli

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/urfave/cli.v2"
)

// LoadResourceFromContext loads data from a CLI flag.
//
// If the value of the specified flag starts with '@', the flag is interpreted
// as a filename and the contents of the file are returned.
//
// In all other cases, the flag value is returned directly.
func LoadResourceFromContext(c *cli.Context, flag string) (string, error) {
	value := c.String(flag)
	if len(value) > 0 && value[0] == '@' {
		b, err := ioutil.ReadFile(value[1:])
		if err != nil {
			return "", fmt.Errorf("couldn't read %s: %s\n", value[1:], err.Error())
		}
		value = string(b)
	}
	return value, nil
}
