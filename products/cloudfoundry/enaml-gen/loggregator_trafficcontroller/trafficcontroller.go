package loggregator_trafficcontroller 
/*
* File Generated by enaml generator
* !!! Please do not edit this file !!!
*/
type TrafficController struct {

	/*SecurityEventLogging - Descr: Enable logging of all requests made to the Traffic Controller in CEF format Default: false
*/
	SecurityEventLogging *SecurityEventLogging `yaml:"security_event_logging,omitempty"`

	/*DisableAccessControl - Descr: Traffic controller bypasses authentication with the UAA and CC Default: false
*/
	DisableAccessControl interface{} `yaml:"disable_access_control,omitempty"`

	/*OutgoingPort - Descr: Port on which the traffic controller listens to for requests Default: 8080
*/
	OutgoingPort interface{} `yaml:"outgoing_port,omitempty"`

	/*Zone - Descr: Zone of the loggregator_trafficcontroller Default: <nil>
*/
	Zone interface{} `yaml:"zone,omitempty"`

	/*Debug - Descr: boolean value to turn on verbose logging for loggregator system (dea agent & loggregator server) Default: false
*/
	Debug interface{} `yaml:"debug,omitempty"`

}