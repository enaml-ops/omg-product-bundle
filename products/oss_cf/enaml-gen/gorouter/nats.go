package gorouter 
/*
* File Generated by enaml generator
* !!! Please do not edit this file !!!
*/
type Nats struct {

	/*User - Descr:  Default: <nil>
*/
	User interface{} `yaml:"user,omitempty"`

	/*Port - Descr:  Default: <nil>
*/
	Port interface{} `yaml:"port,omitempty"`

	/*Password - Descr:  Default: <nil>
*/
	Password interface{} `yaml:"password,omitempty"`

	/*Machines - Descr: IP of each NATS cluster member. Default: <nil>
*/
	Machines interface{} `yaml:"machines,omitempty"`

}