package push_apps_manager 
/*
* File Generated by enaml generator
* !!! Please do not edit this file !!!
*/
type PushAppsManagerJob struct {

	/*Services - Descr: Cloud Foundry login server URL Default: <nil>
*/
	Services *Services `yaml:"services,omitempty"`

	/*Ssl - Descr: When rendering app route urls, prefix all of them with https Default: false
*/
	Ssl *Ssl `yaml:"ssl,omitempty"`

	/*Env - Descr: Default 'from' address for e-mails generated by Console application Default: <nil>
*/
	Env *Env `yaml:"env,omitempty"`

	/*AppUsageService - Descr: Whether or not to allow manual creation of AppEvents Default: false
*/
	AppUsageService *AppUsageService `yaml:"app_usage_service,omitempty"`

	/*Cf - Descr: Cloud Foundry system domain, used for the Console application's URL Default: <nil>
*/
	Cf *Cf `yaml:"cf,omitempty"`

	/*Databases - Descr: IP of database server for the app usage service Default: <nil>
*/
	Databases *Databases `yaml:"databases,omitempty"`

}