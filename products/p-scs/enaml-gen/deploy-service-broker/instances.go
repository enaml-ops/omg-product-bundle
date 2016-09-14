package deploy_service_broker 
/*
* File Generated by enaml generator
* !!! Please do not edit this file !!!
*/
type Instances struct {

	/*InstancesUser - Descr: Username with access to instances space Default: <nil>
*/
	InstancesUser interface{} `yaml:"instances_user,omitempty"`

	/*OrgName - Descr: Org that will host Instances Default: p-spring-cloud-services
*/
	OrgName interface{} `yaml:"org_name,omitempty"`

	/*InstancesPassword - Descr: Password for the username that has access to the instances space Default: <nil>
*/
	InstancesPassword interface{} `yaml:"instances_password,omitempty"`

	/*SpaceName - Descr: Space that will host Instances Default: instances
*/
	SpaceName interface{} `yaml:"space_name,omitempty"`

}