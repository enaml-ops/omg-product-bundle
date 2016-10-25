package cloud_controller_worker 
/*
* File Generated by enaml generator
* !!! Please do not edit this file !!!
*/
type Diego struct {

	/*TpsUrl - Descr: URL of the Diego tps service Default: http://tps.service.cf.internal:1518
*/
	TpsUrl interface{} `yaml:"tps_url,omitempty"`

	/*StagerUrl - Descr: URL of the Diego stager service Default: http://stager.service.cf.internal:8888
*/
	StagerUrl interface{} `yaml:"stager_url,omitempty"`

	/*NsyncUrl - Descr: URL of the Diego nsync service Default: http://nsync.service.cf.internal:8787
*/
	NsyncUrl interface{} `yaml:"nsync_url,omitempty"`

}