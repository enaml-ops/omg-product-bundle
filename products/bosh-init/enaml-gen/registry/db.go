package registry 
/*
* File Generated by enaml generator
* !!! Please do not edit this file !!!
*/
type Db struct {

	/*ConnectionOptions - Descr: Additional options for the database Default: map[max_connections:32 pool_timeout:10]
*/
	ConnectionOptions interface{} `yaml:"connection_options,omitempty"`

	/*Port - Descr: Port of the registry database Default: 5432
*/
	Port interface{} `yaml:"port,omitempty"`

	/*Password - Descr: Password used for the registry database Default: <nil>
*/
	Password interface{} `yaml:"password,omitempty"`

	/*Database - Descr: Name of the registry database Default: bosh_registry
*/
	Database interface{} `yaml:"database,omitempty"`

	/*Adapter - Descr: The type of database used Default: postgres
*/
	Adapter interface{} `yaml:"adapter,omitempty"`

	/*User - Descr: Username used for the registry database Default: bosh
*/
	User interface{} `yaml:"user,omitempty"`

	/*Host - Descr: Address of the registry database Default: 127.0.0.1
*/
	Host interface{} `yaml:"host,omitempty"`

}