package etcd_metrics_server 
/*
* File Generated by enaml generator
* !!! Please do not edit this file !!!
*/
type Status struct {

	/*Username - Descr: basic auth username for metrics server (leave empty for generated) Default: 
*/
	Username interface{} `yaml:"username,omitempty"`

	/*Password - Descr: basic auth password for metrics server (leave empty for generated) Default: 
*/
	Password interface{} `yaml:"password,omitempty"`

	/*Port - Descr: listening port for metrics server Default: 5678
*/
	Port interface{} `yaml:"port,omitempty"`

}