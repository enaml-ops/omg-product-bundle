package loggregator_trafficcontroller 
/*
* File Generated by enaml generator
* !!! Please do not edit this file !!!
*/
type LoggregatorEtcd struct {

	/*RequireSsl - Descr: Enable ssl for all communication with etcd Default: false
*/
	RequireSsl interface{} `yaml:"require_ssl,omitempty"`

	/*Machines - Descr: IPs pointing to the ETCD cluster Default: <nil>
*/
	Machines interface{} `yaml:"machines,omitempty"`

	/*Maxconcurrentrequests - Descr: Number of concurrent requests to ETCD Default: 10
*/
	Maxconcurrentrequests interface{} `yaml:"maxconcurrentrequests,omitempty"`

	/*CaCert - Descr: PEM-encoded CA certificate Default: 
*/
	CaCert interface{} `yaml:"ca_cert,omitempty"`

}