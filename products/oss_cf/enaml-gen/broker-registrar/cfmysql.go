package broker_registrar 
/*
* File Generated by enaml generator
* !!! Please do not edit this file !!!
*/
type CfMysql struct {

	/*Broker - Descr: Basic Auth password for the service broker Default: <nil>
*/
	Broker *Broker `yaml:"broker,omitempty"`

	/*ExternalHost - Descr: Host address of the service broker Default: <nil>
*/
	ExternalHost interface{} `yaml:"external_host,omitempty"`

}