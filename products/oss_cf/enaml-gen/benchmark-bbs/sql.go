package benchmark_bbs 
/*
* File Generated by enaml generator
* !!! Please do not edit this file !!!
*/
type Sql struct {

	/*DbConnectionString - Descr: Connection string to use for SQL backend [username:password@tcp(1.1.1.1:1234)/database] Default: <nil>
*/
	DbConnectionString interface{} `yaml:"db_connection_string,omitempty"`

	/*DbDriver - Descr: driver to use, e.g. postgres or mysql Default: <nil>
*/
	DbDriver interface{} `yaml:"db_driver,omitempty"`

}