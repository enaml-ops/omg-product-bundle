package test_service_health 
/*
* File Generated by enaml generator
* !!! Please do not edit this file !!!
*/
type Api struct {

	/*Org - Descr: CF Org to use for test app Default: system
*/
	Org interface{} `yaml:"org,omitempty"`

	/*Url - Descr: Cloud Controller API address Default: <nil>
*/
	Url interface{} `yaml:"url,omitempty"`

	/*Space - Descr: CF Space with in Org to use for test app Default: gemfire-smoke-test-space-57818572-4437-45d8-a25b-d71e5f5eae7d
*/
	Space interface{} `yaml:"space,omitempty"`

	/*Password - Descr: Password for authentication Default: <nil>
*/
	Password interface{} `yaml:"password,omitempty"`

	/*Username - Descr: Username to authenticate with Default: <nil>
*/
	Username interface{} `yaml:"username,omitempty"`

}