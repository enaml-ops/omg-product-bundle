package cf_mysql_broker 
/*
* File Generated by enaml generator
* !!! Please do not edit this file !!!
*/
type CfMysqlBrokerJob struct {

	/*Cf - Descr: Determines whether dashboard verifies SSL certificates when communicating with Cloud Controller and UAA Default: false
*/
	Cf *Cf `yaml:"cf,omitempty"`

	/*Nats - Descr: Username for broker to register a route with NATS Default: <nil>
*/
	Nats *Nats `yaml:"nats,omitempty"`

	/*CfMysql - Descr: Optional, The ip to be registered with the cf router for the broker. Defaults to the ip of the vm Default: <nil>
*/
	CfMysql *CfMysql `yaml:"cf_mysql,omitempty"`

	/*SyslogAggregator - Descr: Transport to be used when forwarding logs (tcp|udp|relp). Default: tcp
*/
	SyslogAggregator *SyslogAggregator `yaml:"syslog_aggregator,omitempty"`

}