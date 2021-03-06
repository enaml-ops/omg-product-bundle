package test_notifications 
/*
* File Generated by enaml generator
* !!! Please do not edit this file !!!
*/
type Notifications struct {

	/*Organization - Descr: Organization that contains the app Default: <nil>
*/
	Organization interface{} `yaml:"organization,omitempty"`

	/*Space - Descr: Space that contains the app Default: <nil>
*/
	Space interface{} `yaml:"space,omitempty"`

	/*Tests - Descr: Toggle for running the performance tests Default: <nil>
*/
	Tests *Tests `yaml:"tests,omitempty"`

	/*Uaa - Descr: UAA Admin client ID Default: <nil>
*/
	Uaa *Uaa `yaml:"uaa,omitempty"`

	/*Cf - Descr: Username of the CF admin user Default: <nil>
*/
	Cf *Cf `yaml:"cf,omitempty"`

	/*AppDomain - Descr: Domain used to host application Default: <nil>
*/
	AppDomain interface{} `yaml:"app_domain,omitempty"`

}