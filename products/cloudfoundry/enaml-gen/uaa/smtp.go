package uaa 
/*
* File Generated by enaml generator
* !!! Please do not edit this file !!!
*/
type Smtp struct {

	/*Password - Descr: SMTP server password Default: <nil>
*/
	Password interface{} `yaml:"password,omitempty"`

	/*Port - Descr: SMTP server port Default: 2525
*/
	Port interface{} `yaml:"port,omitempty"`

	/*Host - Descr: SMTP server host address Default: localhost
*/
	Host interface{} `yaml:"host,omitempty"`

	/*User - Descr: SMTP server username Default: <nil>
*/
	User interface{} `yaml:"user,omitempty"`

}