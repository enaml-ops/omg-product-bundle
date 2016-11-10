package proxy 
/*
* File Generated by enaml generator
* !!! Please do not edit this file !!!
*/
type SyslogAggregator struct {

	/*Port - Descr: TCP port of syslog aggregator Default: <nil>
*/
	Port interface{} `yaml:"port,omitempty"`

	/*Transport - Descr: Transport to be used when forwarding logs (tcp|udp|relp). Default: tcp
*/
	Transport interface{} `yaml:"transport,omitempty"`

	/*Address - Descr: IP address for syslog aggregator Default: <nil>
*/
	Address interface{} `yaml:"address,omitempty"`

}