package file_server 
/*
* File Generated by enaml generator
* !!! Please do not edit this file !!!
*/
type FileServer struct {

	/*DropsondePort - Descr: local metron agent's port Default: 3457
*/
	DropsondePort interface{} `yaml:"dropsonde_port,omitempty"`

	/*ListenAddr - Descr: Address of interface on which to serve files Default: 0.0.0.0:8080
*/
	ListenAddr interface{} `yaml:"listen_addr,omitempty"`

	/*StaticDirectory - Descr: Fully-qualified path to the doc root for the file server's static files Default: /var/vcap/jobs/file_server/packages/
*/
	StaticDirectory interface{} `yaml:"static_directory,omitempty"`

	/*DebugAddr - Descr: address at which to serve debug info Default: 0.0.0.0:17005
*/
	DebugAddr interface{} `yaml:"debug_addr,omitempty"`

	/*LogLevel - Descr: Log level Default: info
*/
	LogLevel interface{} `yaml:"log_level,omitempty"`

}