package config

// VMType contains the VM types for a Cloud Foundry deployment.
type VMType struct {
	NatsVMType                  string `omg:"nats-vm-type"`
	ConsulVMType                string `omg:"consul-vm-type"`
	MySQLVMType                 string `omg:"mysql-vm-type"`
	EtcdVMType                  string `omg:"etcd-vm-type"`
	ClockGlobalVMType           string `omg:"clock-global-vm-type"`
	MySQLProxyVMType            string `omg:"mysql-proxy-vm-type"`
	HAProxyVMType               string `omg:"haproxy-vm-type,optional"`
	RouterVMType                string `omg:"router-vm-type"`
	NFSVMType                   string `omg:"nfs-vm-type"`
	CloudControllerVMType       string `omg:"cc-vm-type"`
	CloudControllerWorkerVMType string `omg:"cc-worker-vm-type"`
	DiegoDBVMType               string `omg:"diego-db-vm-type"`
	UAAVMType                   string `omg:"uaa-vm-type"`
	DiegoCellVMType             string `omg:"diego-cell-vm-type"`
	DiegoBrainVMType            string `omg:"diego-brain-vm-type"`
	DopplerVMType               string `omg:"doppler-vm-type"`
	ErrandVMType                string `omg:"errand-vm-type"`
	LoggregratorVMType          string `omg:"loggregator-traffic-controller-vmtype"`
}
