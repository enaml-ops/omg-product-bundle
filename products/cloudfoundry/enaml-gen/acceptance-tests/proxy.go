package acceptance_tests 
/*
* File Generated by enaml generator
* !!! Please do not edit this file !!!
*/
type Proxy struct {

	/*ExternalHost - Descr: Proxy external host (e.g. p-mysql.example.com => proxy-0-p-mysql.example.com) Default: <nil>
*/
	ExternalHost interface{} `yaml:"external_host,omitempty"`

	/*ApiPassword - Descr: Proxy API password Default: <nil>
*/
	ApiPassword interface{} `yaml:"api_password,omitempty"`

	/*ApiForceHttps - Descr: Expect proxy to force redirect to HTTPS Default: true
*/
	ApiForceHttps interface{} `yaml:"api_force_https,omitempty"`

	/*SkipSslValidation - Descr: Tests will skip validation of SSL certificates Default: true
*/
	SkipSslValidation interface{} `yaml:"skip_ssl_validation,omitempty"`

	/*ProxyCount - Descr: Number of proxy instances. Use to construct an array of proxy dashboard url (e.g. https://proxy-INDEX-EXTERNAL_HOST) Default: <nil>
*/
	ProxyCount interface{} `yaml:"proxy_count,omitempty"`

	/*ApiUsername - Descr: Proxy API username Default: <nil>
*/
	ApiUsername interface{} `yaml:"api_username,omitempty"`

}