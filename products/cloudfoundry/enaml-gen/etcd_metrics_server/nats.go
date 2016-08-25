package etcd_metrics_server 
/*
* File Generated by enaml generator
* !!! Please do not edit this file !!!
*/
type Nats struct {

	/*Password - Descr: NATS server password Default: <nil>
*/
	Password interface{} `yaml:"password,omitempty"`

	/*Machines - Descr: array of NATS addresses Default: <nil>
*/
	Machines interface{} `yaml:"machines,omitempty"`

	/*Port - Descr: NATS server port Default: 4222
*/
	Port interface{} `yaml:"port,omitempty"`

	/*Username - Descr: NATS server username Default: <nil>
*/
	Username interface{} `yaml:"username,omitempty"`

}