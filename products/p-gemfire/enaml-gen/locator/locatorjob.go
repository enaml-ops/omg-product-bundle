package locator 
/*
* File Generated by enaml generator
* !!! Please do not edit this file !!!
*/
type LocatorJob struct {

	/*ExternalDependencies - Descr: System domain Default: <nil>
*/
	ExternalDependencies *ExternalDependencies `yaml:"external_dependencies,omitempty"`

	/*Gemfire - Descr: min number of locators which should be present Default: 2
*/
	Gemfire *Gemfire `yaml:"gemfire,omitempty"`

}