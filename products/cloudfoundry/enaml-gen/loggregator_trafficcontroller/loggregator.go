package loggregator_trafficcontroller 
/*
* File Generated by enaml generator
* !!! Please do not edit this file !!!
*/
type Loggregator struct {

	/*OutgoingDropsondePort - Descr: Port for outgoing dropsonde messages Default: 8081
*/
	OutgoingDropsondePort interface{} `yaml:"outgoing_dropsonde_port,omitempty"`

	/*DopplerPort - Descr: Port for outgoing doppler messages Default: 8081
*/
	DopplerPort interface{} `yaml:"doppler_port,omitempty"`

	/*Etcd - Descr: IPs pointing to the ETCD cluster Default: <nil>
*/
	Etcd *Etcd `yaml:"etcd,omitempty"`

}