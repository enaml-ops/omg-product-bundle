package smoke_tests 
/*
* File Generated by enaml generator
* !!! Please do not edit this file !!!
*/
type CfSmokeTests struct {

	/*UseExistingOrg - Descr: Runs smoke test errand as an existing org. Creates a new org if false Default: false
*/
	UseExistingOrg interface{} `yaml:"use_existing_org,omitempty"`

	/*Org - Descr: The name of the Org to run acceptance tests against Default: 
*/
	Org interface{} `yaml:"org,omitempty"`

}