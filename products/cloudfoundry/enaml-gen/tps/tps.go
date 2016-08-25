package tps 
/*
* File Generated by enaml generator
* !!! Please do not edit this file !!!
*/
type Tps struct {

	/*MaxInFlightRequests - Descr: Maximum number of requests to handle at once. Default: 200
*/
	MaxInFlightRequests interface{} `yaml:"max_in_flight_requests,omitempty"`

	/*Bbs - Descr: PEM-encoded client key Default: <nil>
*/
	Bbs *Bbs `yaml:"bbs,omitempty"`

	/*Cc - Descr: External port to access the Cloud Controller Default: 9022
*/
	Cc *Cc `yaml:"cc,omitempty"`

	/*DropsondePort - Descr: local metron agent's port Default: 3457
*/
	DropsondePort interface{} `yaml:"dropsonde_port,omitempty"`

	/*Watcher - Descr: address at which to serve debug info Default: 0.0.0.0:17015
*/
	Watcher *Watcher `yaml:"watcher,omitempty"`

	/*LogLevel - Descr: Log level Default: info
*/
	LogLevel interface{} `yaml:"log_level,omitempty"`

	/*Listener - Descr: address at which to serve API requests Default: 0.0.0.0:1518
*/
	Listener *Listener `yaml:"listener,omitempty"`

	/*ConsulAgentPort - Descr: local consul agent's port Default: 8500
*/
	ConsulAgentPort interface{} `yaml:"consul_agent_port,omitempty"`

	/*TrafficControllerUrl - Descr: URL of Traffic controller Default: ws://loggregator-trafficcontroller.service.cf.internal:8081
*/
	TrafficControllerUrl interface{} `yaml:"traffic_controller_url,omitempty"`

}