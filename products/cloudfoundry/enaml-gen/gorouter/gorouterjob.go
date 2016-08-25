package gorouter 
/*
* File Generated by enaml generator
* !!! Please do not edit this file !!!
*/
type GorouterJob struct {

	/*Dropsonde - Descr: Enable the dropsonde emitter library Default: false
*/
	Dropsonde *Dropsonde `yaml:"dropsonde,omitempty"`

	/*RoutingApi - Descr: Enable the GoRouter to receive routes from the Routing API Default: false
*/
	RoutingApi *RoutingApi `yaml:"-"`

	/*RequestTimeoutInSeconds - Descr: Timeout in seconds for Router -> Endpoint roundtrip. Default: 900
*/
	RequestTimeoutInSeconds interface{} `yaml:"request_timeout_in_seconds,omitempty"`

	/*Uaa - Descr: Port on which UAA is running. Default: 8080
*/
	Uaa *Uaa `yaml:"uaa,omitempty"`

	/*Router - Descr: Enables streaming of access log to syslog. Warning: this comes with a performance cost; due to higher I/O, max request rate is reduced. Default: false
*/
	Router *Router `yaml:"router,omitempty"`

	/*Nats - Descr:  Default: <nil>
*/
	Nats *Nats `yaml:"nats,omitempty"`

	/*MetronEndpoint - Descr: The port used to emit legacy messages to the Metron agent. Default: 3456
*/
	MetronEndpoint *MetronEndpoint `yaml:"metron_endpoint,omitempty"`

}