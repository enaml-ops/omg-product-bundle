package etcd_metrics_server 
/*
* File Generated by enaml generator
* !!! Please do not edit this file !!!
*/
type EtcdMetricsServer struct {

	/*Etcd - Descr: PEM-encoded client key Default: 
*/
	Etcd *Etcd `yaml:"etcd,omitempty"`

	/*Status - Descr: basic auth password for metrics server (leave empty for generated) Default: 
*/
	Status *Status `yaml:"status,omitempty"`

}