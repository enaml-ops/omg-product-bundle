package cloud_controller_worker 
/*
* File Generated by enaml generator
* !!! Please do not edit this file !!!
*/
type Thresholds struct {

	/*Worker - Descr: The cc will alert if memory remains above this threshold for 3 monit cycles Default: 384
*/
	Worker *Worker `yaml:"worker,omitempty"`

}