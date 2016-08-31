package config

import "github.com/codegangsta/cli"

func RequiredInstanceCountFlags() []string {
	return []string{"cc-instances", "uaa-instances", "cc-worker-instances"}
}

func NewInstanceCount(c *cli.Context) InstanceCount {
	return InstanceCount{
		CloudControllerInstances:       c.Int("cc-instances"),
		UAAInstances:                   c.Int("uaa-instances"),
		CloudControllerWorkerInstances: c.Int("cc-worker-instances"),
	}
}

type InstanceCount struct {
	CloudControllerWorkerInstances int
	CloudControllerInstances       int
	UAAInstances                   int
}
