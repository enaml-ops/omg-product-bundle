package broker_registrar 
/*
* File Generated by enaml generator
* !!! Please do not edit this file !!!
*/
type Cf struct {

	/*AdminPassword - Descr: Password of the admin user Default: <nil>
*/
	AdminPassword interface{} `yaml:"admin_password,omitempty"`

	/*ApiUrl - Descr: Full URL of Cloud Foundry API Default: <nil>
*/
	ApiUrl interface{} `yaml:"api_url,omitempty"`

	/*SkipSslValidation - Descr: Skip SSL validation when connecting to Cloud Foundry API Default: false
*/
	SkipSslValidation interface{} `yaml:"skip_ssl_validation,omitempty"`

	/*AdminUsername - Descr: Username of the admin user Default: <nil>
*/
	AdminUsername interface{} `yaml:"admin_username,omitempty"`

}