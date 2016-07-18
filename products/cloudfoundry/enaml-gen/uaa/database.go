package uaa 
/*
* File Generated by enaml generator
* !!! Please do not edit this file !!!
*/
type Database struct {

	/*LogAbandoned - Descr: Should connections that are forcibly closed be logged. Default: true
*/
	LogAbandoned interface{} `yaml:"log_abandoned,omitempty"`

	/*CaseInsensitive - Descr: Set to true if you don't want to be using LOWER() SQL functions in search queries/filters, because you know that your DB is case insensitive. If this property is null, then it will be set to true if the UAA DB is MySQL and false otherwise, but even on MySQL you can override it by setting it explicitly to false Default: <nil>
*/
	CaseInsensitive interface{} `yaml:"case_insensitive,omitempty"`

	/*RemoveAbandoned - Descr: True if connections that are left open longer then abandoned_timeout seconds during a session(time between borrow and return from pool) should be forcibly closed Default: false
*/
	RemoveAbandoned interface{} `yaml:"remove_abandoned,omitempty"`

	/*MaxConnections - Descr: The max number of open connections to the DB from a running UAA instance Default: 100
*/
	MaxConnections interface{} `yaml:"max_connections,omitempty"`

	/*AbandonedTimeout - Descr: Timeout in seconds for the longest running queries. Take into DB migrations for this timeout as they may run during a long period of time. Default: 300
*/
	AbandonedTimeout interface{} `yaml:"abandoned_timeout,omitempty"`

	/*MaxIdleConnections - Descr: The max number of open idle connections to the DB from a running UAA instance Default: 10
*/
	MaxIdleConnections interface{} `yaml:"max_idle_connections,omitempty"`

}