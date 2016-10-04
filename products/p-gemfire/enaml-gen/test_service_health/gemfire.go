package test_service_health 
/*
* File Generated by enaml generator
* !!! Please do not edit this file !!!
*/
type Gemfire struct {

	/*ServicePlanName - Descr: Name of service plan registered with cloud controller Default: <nil>
*/
	ServicePlanName interface{} `yaml:"service_plan_name,omitempty"`

	/*Testing - Descr: Name to give test GemFire CF service Default: gemfire-smoke-test-service-70552f90-7785-4347-9516
*/
	Testing *Testing `yaml:"testing,omitempty"`

	/*Locator - Descr: IP address of a Gemfire cache server to test data validity Default: <nil>
*/
	Locator *Locator `yaml:"locator,omitempty"`

	/*ServiceName - Descr: Name of service registered with cloud controller Default: <nil>
*/
	ServiceName interface{} `yaml:"service_name,omitempty"`

	/*Capabilities - Descr: A List of functionality a cluster will provide Default: []
*/
	Capabilities interface{} `yaml:"capabilities,omitempty"`

}