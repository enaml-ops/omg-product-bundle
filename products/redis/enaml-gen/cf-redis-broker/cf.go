package cf_redis_broker 
/*
* File Generated by enaml generator
* !!! Please do not edit this file !!!
*/
type Cf struct {

	/*Nats - Descr: The password to use when authenticating with NATS Default: <nil>
*/
	Nats *Nats `yaml:"nats,omitempty"`

	/*AppsDomain - Descr: Domain shared by the UAA and CF API eg 'bosh-lite.com' Default: <nil>
*/
	AppsDomain interface{} `yaml:"apps_domain,omitempty"`

}
