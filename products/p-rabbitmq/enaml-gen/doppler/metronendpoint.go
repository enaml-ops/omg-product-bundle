package doppler 
/*
* File Generated by enaml generator
* !!! Please do not edit this file !!!
*/
type MetronEndpoint struct {

	/*Host - Descr: The host used to emit messages to the Metron agent Default: 127.0.0.1
*/
	Host interface{} `yaml:"host,omitempty"`

	/*DropsondePort - Descr: The port used to emit dropsonde messages to the Metron agent Default: 3457
*/
	DropsondePort interface{} `yaml:"dropsonde_port,omitempty"`

}