package cloud_controller_worker 
/*
* File Generated by enaml generator
* !!! Please do not edit this file !!!
*/
type Newrelic struct {

	/*LicenseKey - Descr: The api key for NewRelic Default: <nil>
*/
	LicenseKey interface{} `yaml:"license_key,omitempty"`

	/*DeveloperMode - Descr: Activate NewRelic developer mode Default: false
*/
	DeveloperMode interface{} `yaml:"developer_mode,omitempty"`

	/*CaptureParams - Descr: Capture and send query params to NewRelic Default: false
*/
	CaptureParams interface{} `yaml:"capture_params,omitempty"`

	/*EnvironmentName - Descr: The environment name used by NewRelic Default: development
*/
	EnvironmentName interface{} `yaml:"environment_name,omitempty"`

	/*TransactionTracer - Descr: Enable transaction tracing in NewRelic Default: false
*/
	TransactionTracer *TransactionTracer `yaml:"transaction_tracer,omitempty"`

	/*MonitorMode - Descr: Activate NewRelic monitor mode Default: false
*/
	MonitorMode interface{} `yaml:"monitor_mode,omitempty"`

	/*LogFilePath - Descr: The location for NewRelic to log to Default: /var/vcap/sys/log/cloud_controller_ng/newrelic
*/
	LogFilePath interface{} `yaml:"log_file_path,omitempty"`

}