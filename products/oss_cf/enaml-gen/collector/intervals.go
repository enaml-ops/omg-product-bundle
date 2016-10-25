package collector 
/*
* File Generated by enaml generator
* !!! Please do not edit this file !!!
*/
type Intervals struct {

	/*Varz - Descr: the interval in seconds that varz is checked Default: 30
*/
	Varz interface{} `yaml:"varz,omitempty"`

	/*LocalMetrics - Descr: the interval in seconds that local_metrics are checked Default: 30
*/
	LocalMetrics interface{} `yaml:"local_metrics,omitempty"`

	/*Healthz - Descr: the interval in seconds that healthz is checked Default: 30
*/
	Healthz interface{} `yaml:"healthz,omitempty"`

	/*NatsPing - Descr: the interval in seconds that the collector pings nats to record latency Default: 30
*/
	NatsPing interface{} `yaml:"nats_ping,omitempty"`

	/*Discover - Descr: the interval in seconds that the collector attempts to discover components Default: 60
*/
	Discover interface{} `yaml:"discover,omitempty"`

	/*Prune - Descr: the interval in seconds that the collector attempts to prune unresponsive components Default: 300
*/
	Prune interface{} `yaml:"prune,omitempty"`

}