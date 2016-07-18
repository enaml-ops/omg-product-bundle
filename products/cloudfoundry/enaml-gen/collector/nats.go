package collector 
/*
* File Generated by enaml generator
* !!! Please do not edit this file !!!
*/
type Nats struct {

	/*Password - Descr: NATS password Default: <nil>
*/
	Password interface{} `yaml:"password,omitempty"`

	/*User - Descr: NATS user Default: <nil>
*/
	User interface{} `yaml:"user,omitempty"`

	/*Machines - Descr: IP of each NATS cluster member. Default: <nil>
*/
	Machines interface{} `yaml:"machines,omitempty"`

	/*Port - Descr: NATS TCP port Default: <nil>
*/
	Port interface{} `yaml:"port,omitempty"`

}