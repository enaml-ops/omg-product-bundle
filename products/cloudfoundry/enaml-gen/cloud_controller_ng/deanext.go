package cloud_controller_ng 
/*
* File Generated by enaml generator
* !!! Please do not edit this file !!!
*/
type DeaNext struct {

	/*ClientKey - Descr: PEM-encoded server key Default: <nil>
*/
	ClientKey interface{} `yaml:"client_key,omitempty"`

	/*StagingDiskLimitMb - Descr: Disk limit in mb for staging tasks Default: 6144
*/
	StagingDiskLimitMb interface{} `yaml:"staging_disk_limit_mb,omitempty"`

	/*ClientCert - Descr: PEM-encoded server certificate Default: <nil>
*/
	ClientCert interface{} `yaml:"client_cert,omitempty"`

	/*CaCert - Descr: PEM-encoded CA certificate Default: <nil>
*/
	CaCert interface{} `yaml:"ca_cert,omitempty"`

	/*StagingMemoryLimitMb - Descr: Memory limit in mb for staging tasks Default: 1024
*/
	StagingMemoryLimitMb interface{} `yaml:"staging_memory_limit_mb,omitempty"`

	/*AdvertiseIntervalInSeconds - Descr: Advertise interval for DEAs Default: 5
*/
	AdvertiseIntervalInSeconds interface{} `yaml:"advertise_interval_in_seconds,omitempty"`

}