package loggregator_trafficcontroller 
/*
* File Generated by enaml generator
* !!! Please do not edit this file !!!
*/
type TrafficController struct {

	/*LockedMemoryLimit - Descr: Size (KB) of shell's locked memory limit. Set to 'kernel' to use the kernel's default. Non-numeric values other than 'kernel', 'soft', 'hard', and 'unlimited' will result in an error. Default: unlimited
*/
	LockedMemoryLimit interface{} `yaml:"locked_memory_limit,omitempty"`

	/*Debug - Descr: boolean value to turn on verbose logging for loggregator system (dea agent & loggregator server) Default: false
*/
	Debug interface{} `yaml:"debug,omitempty"`

	/*DisableAccessControl - Descr: Traffic controller bypasses authentication with the UAA and CC Default: false
*/
	DisableAccessControl interface{} `yaml:"disable_access_control,omitempty"`

	/*Etcd - Descr: PEM-encoded client key Default: 
*/
	Etcd *TrafficControllerEtcd `yaml:"etcd,omitempty"`

	/*SecurityEventLogging - Descr: Enable logging of all requests made to the Traffic Controller in CEF format Default: false
*/
	SecurityEventLogging *SecurityEventLogging `yaml:"security_event_logging,omitempty"`

	/*OutgoingPort - Descr: DEPRECATED Default: 8080
*/
	OutgoingPort interface{} `yaml:"outgoing_port,omitempty"`

}