package atc 
/*
* File Generated by enaml generator
* !!! Please do not edit this file !!!
*/
type GenericOauth struct {

	/*TokenUrl - Descr: Generic OAuth provider token endpoint URL. Default: 
*/
	TokenUrl interface{} `yaml:"token_url,omitempty"`

	/*ClientId - Descr: Application client ID for enabling generic OAuth. Default: 
*/
	ClientId interface{} `yaml:"client_id,omitempty"`

	/*AuthUrlParams - Descr: List Parameter to pass to the authentication server authorization url.
 Default: map[]
*/
	AuthUrlParams interface{} `yaml:"auth_url_params,omitempty"`

	/*DisplayName - Descr: Name of the authentication method to be displayed on the Web UI Default: 
*/
	DisplayName interface{} `yaml:"display_name,omitempty"`

	/*ClientSecret - Descr: Application client secret for enabling generic OAuth. Default: 
*/
	ClientSecret interface{} `yaml:"client_secret,omitempty"`

	/*AuthUrl - Descr: Generic OAuth provider authorization endpoint url. Default: 
*/
	AuthUrl interface{} `yaml:"auth_url,omitempty"`

	/*Scope - Descr: OAuth scope required for users who will have access. Default: 
*/
	Scope interface{} `yaml:"scope,omitempty"`

}