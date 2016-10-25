package uaa 
/*
* File Generated by enaml generator
* !!! Please do not edit this file !!!
*/
type LdapSsl struct {

	/*Tls - Descr: If using StartTLS, what mode to enable. Default is none, not enabled. Possible values are none, simple, external Default: none
*/
	Tls interface{} `yaml:"tls,omitempty"`

	/*Skipverification - Descr: Set to true, and LDAPS connection will not validate the server certificate. Default: false
*/
	Skipverification interface{} `yaml:"skipverification,omitempty"`

}