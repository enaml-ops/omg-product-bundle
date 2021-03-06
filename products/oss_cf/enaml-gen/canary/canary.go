package canary 
/*
* File Generated by enaml generator
* !!! Please do not edit this file !!!
*/
type Canary struct {

	/*AppName - Descr: App name for the canary app Default: <nil>
*/
	AppName interface{} `yaml:"app_name,omitempty"`

	/*Space - Descr: The Elastic Runtime space name to use for the canary app Default: <nil>
*/
	Space interface{} `yaml:"space,omitempty"`

	/*AppDomain - Descr: The domain to use for the canary app Default: <nil>
*/
	AppDomain interface{} `yaml:"app_domain,omitempty"`

	/*DatadogApiKey - Descr: Datadog API key for the canary app Default: <nil>
*/
	DatadogApiKey interface{} `yaml:"datadog_api_key,omitempty"`

	/*Password - Descr: The Elastic Runtime API user's password Default: <nil>
*/
	Password interface{} `yaml:"password,omitempty"`

	/*Api - Descr: The Elastic Runtime API endpoint URL Default: <nil>
*/
	Api interface{} `yaml:"api,omitempty"`

	/*User - Descr: The Elastic Runtime API user Default: <nil>
*/
	User interface{} `yaml:"user,omitempty"`

	/*Org - Descr: The Elastic Runtime organization name to use for the canary app Default: <nil>
*/
	Org interface{} `yaml:"org,omitempty"`

	/*DeploymentName - Descr: Deployment name for the canary app Default: <nil>
*/
	DeploymentName interface{} `yaml:"deployment_name,omitempty"`

	/*InstanceCount - Descr: Number of instances of the canary app Default: <nil>
*/
	InstanceCount interface{} `yaml:"instance_count,omitempty"`

	/*CfStack - Descr: Stack for the canary app Default: cflinuxfs2
*/
	CfStack interface{} `yaml:"cf_stack,omitempty"`

}