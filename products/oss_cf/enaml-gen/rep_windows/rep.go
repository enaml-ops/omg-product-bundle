package rep_windows 
/*
* File Generated by enaml generator
* !!! Please do not edit this file !!!
*/
type Rep struct {

	/*DebugAddr - Descr: address at which to serve debug info Default: 127.0.0.1:17008
*/
	DebugAddr interface{} `yaml:"debug_addr,omitempty"`

	/*PreloadedRootfses - Descr: Array of name:absolute_path pairs representing root filesystems preloaded onto the underlying garden Default: [windows2012R2:/tmp/windows2012R2]
*/
	PreloadedRootfses interface{} `yaml:"preloaded_rootfses,omitempty"`

	/*EvacuationTimeoutInSeconds - Descr: The time to wait for evacuation to complete in seconds Default: 600
*/
	EvacuationTimeoutInSeconds interface{} `yaml:"evacuation_timeout_in_seconds,omitempty"`

	/*Stack - Descr: The stack for which to handle requests Default: windows2012R2
*/
	Stack interface{} `yaml:"stack,omitempty"`

	/*RootfsProviders - Descr: Array of schemes for which the underlying garden can support arbitrary root filesystems Default: [docker]
*/
	RootfsProviders interface{} `yaml:"rootfs_providers,omitempty"`

	/*Bbs - Descr: maximum number of idle http connections Default: <nil>
*/
	Bbs *Bbs `yaml:"bbs,omitempty"`

	/*EvacuationPollingIntervalInSeconds - Descr: The interval to look for completed tasks and LRPs during evacuation in seconds Default: 10
*/
	EvacuationPollingIntervalInSeconds interface{} `yaml:"evacuation_polling_interval_in_seconds,omitempty"`

	/*PollingIntervalInSeconds - Descr: The interval to look for completed tasks and LRPs in seconds Default: 30
*/
	PollingIntervalInSeconds interface{} `yaml:"polling_interval_in_seconds,omitempty"`

	/*Zone - Descr: The zone associated with the rep Default: <nil>
*/
	Zone interface{} `yaml:"zone,omitempty"`

	/*TrustedCerts - Descr: Concatenation of trusted CA certificates to be made available on the cell. Default: <nil>
*/
	TrustedCerts interface{} `yaml:"trusted_certs,omitempty"`

	/*ListenAddr - Descr: address to serve auction and LRP stop requests on Default: 0.0.0.0:1800
*/
	ListenAddr interface{} `yaml:"listen_addr,omitempty"`

	/*LogLevel - Descr: Log level Default: info
*/
	LogLevel interface{} `yaml:"log_level,omitempty"`

	/*ConsulAgentPort - Descr: local consul agent's port Default: 8500
*/
	ConsulAgentPort interface{} `yaml:"consul_agent_port,omitempty"`

	/*DropsondePort - Descr: local metron agent's port Default: 3457
*/
	DropsondePort interface{} `yaml:"dropsonde_port,omitempty"`

}