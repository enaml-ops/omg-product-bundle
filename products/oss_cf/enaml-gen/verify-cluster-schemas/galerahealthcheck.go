package verify_cluster_schemas 
/*
* File Generated by enaml generator
* !!! Please do not edit this file !!!
*/
type GaleraHealthcheck struct {

	/*EndpointUsername - Descr: Username used to authenticate with galera healthcheck Default: <nil>
*/
	EndpointUsername interface{} `yaml:"endpoint_username,omitempty"`

	/*EndpointPassword - Descr: Password used to authenticate with galera healthcheck Default: <nil>
*/
	EndpointPassword interface{} `yaml:"endpoint_password,omitempty"`

}