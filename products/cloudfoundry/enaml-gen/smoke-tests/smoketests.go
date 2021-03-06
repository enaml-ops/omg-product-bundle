package smoke_tests 
/*
* File Generated by enaml generator
* !!! Please do not edit this file !!!
*/
type SmokeTests struct {

	/*AppsDomain - Descr: The Elastic Runtime Application Domain Default: <nil>
*/
	AppsDomain interface{} `yaml:"apps_domain,omitempty"`

	/*RuntimeApp - Descr: The Elastic Runtime app name to use when running runtime tests Default: 
*/
	RuntimeApp interface{} `yaml:"runtime_app,omitempty"`

	/*Backend - Descr: Defines the backend to be used. ('dea', 'diego', '' (default)) Default: 
*/
	Backend interface{} `yaml:"backend,omitempty"`

	/*UseExistingOrg - Descr: Toggles setup and cleanup of the Elastic Runtime organization Default: false
*/
	UseExistingOrg interface{} `yaml:"use_existing_org,omitempty"`

	/*Api - Descr: The Elastic Runtime API endpoint URL Default: <nil>
*/
	Api interface{} `yaml:"api,omitempty"`

	/*EnableWindowsTests - Descr: Toggles a portion of the suite that exercises Windows platform support Default: false
*/
	EnableWindowsTests interface{} `yaml:"enable_windows_tests,omitempty"`

	/*SuiteName - Descr: A token used by the tests when creating Apps / Spaces Default: CF_SMOKE_TESTS
*/
	SuiteName interface{} `yaml:"suite_name,omitempty"`

	/*UseExistingSpace - Descr: Toggles setup and cleanup of the Elastic Runtime space Default: false
*/
	UseExistingSpace interface{} `yaml:"use_existing_space,omitempty"`

	/*User - Descr: The Elastic Runtime API user Default: <nil>
*/
	User interface{} `yaml:"user,omitempty"`

	/*Org - Descr: The Elastic Runtime organization name to use when running tests Default: <nil>
*/
	Org interface{} `yaml:"org,omitempty"`

	/*SkipSslValidation - Descr: Toggles cli verification of the Elastic Runtime API SSL certificate Default: false
*/
	SkipSslValidation interface{} `yaml:"skip_ssl_validation,omitempty"`

	/*LoggingApp - Descr: The Elastic Runtime app name to use when running logging tests Default: 
*/
	LoggingApp interface{} `yaml:"logging_app,omitempty"`

	/*GinkgoOpts - Descr: Ginkgo options for the smoke tests Default: 
*/
	GinkgoOpts interface{} `yaml:"ginkgo_opts,omitempty"`

	/*Password - Descr: The Elastic Runtime API user's password Default: <nil>
*/
	Password interface{} `yaml:"password,omitempty"`

	/*Space - Descr: The Elastic Runtime space name to use when running tests Default: <nil>
*/
	Space interface{} `yaml:"space,omitempty"`

}