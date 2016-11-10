package config

// InstanceCount contains the configurable instance counts for a cloud foundry deployment.
type InstanceCount struct {
	CloudControllerWorkerInstances int `omg:"cc-worker-instances"`
	CloudControllerInstances       int `omg:"cc-instances"`
	UAAInstances                   int `omg:"uaa-instances"`
}
