package acceptance_tests 
/*
* File Generated by enaml generator
* !!! Please do not edit this file !!!
*/
type Cf struct {

	/*AdminPassword - Descr: Password of the admin user Default: <nil>
*/
	AdminPassword interface{} `yaml:"admin_password,omitempty"`

	/*SkipSslValidation - Descr: Whether to add --skip-ssl-validation for cf cli Default: false
*/
	SkipSslValidation interface{} `yaml:"skip_ssl_validation,omitempty"`

	/*ApiUrl - Descr: Full URL of Cloud Foundry API Default: <nil>
*/
	ApiUrl interface{} `yaml:"api_url,omitempty"`

	/*AppsDomain - Descr: Shared domain for pushed apps Default: <nil>
*/
	AppsDomain interface{} `yaml:"apps_domain,omitempty"`

	/*AdminUsername - Descr: Username of the admin user Default: <nil>
*/
	AdminUsername interface{} `yaml:"admin_username,omitempty"`

}