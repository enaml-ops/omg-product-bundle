package converger 
/*
* File Generated by enaml generator
* !!! Please do not edit this file !!!
*/
type Converger struct {

	/*DropsondePort - Descr: local metron agent's port Default: 3457
*/
	DropsondePort interface{} `yaml:"dropsonde_port,omitempty"`

	/*RepeatIntervalInSeconds - Descr: the interval between runs of the converge process Default: 30
*/
	RepeatIntervalInSeconds interface{} `yaml:"repeat_interval_in_seconds,omitempty"`

	/*ExpireCompletedTaskDurationInSeconds - Descr: completed, unresolved tasks are deleted after this duration in seconds Default: 120
*/
	ExpireCompletedTaskDurationInSeconds interface{} `yaml:"expire_completed_task_duration_in_seconds,omitempty"`

	/*LogLevel - Descr: Log level Default: info
*/
	LogLevel interface{} `yaml:"log_level,omitempty"`

	/*KickTaskDurationInSeconds - Descr: the interval, in seconds, between kicks to tasks in seconds Default: 30
*/
	KickTaskDurationInSeconds interface{} `yaml:"kick_task_duration_in_seconds,omitempty"`

	/*ExpirePendingTaskDurationInSeconds - Descr: unclaimed tasks are marked as failed, after this duration in seconds Default: 1800
*/
	ExpirePendingTaskDurationInSeconds interface{} `yaml:"expire_pending_task_duration_in_seconds,omitempty"`

	/*Bbs - Descr: Address to the BBS Server Default: bbs.service.cf.internal:8889
*/
	Bbs *Bbs `yaml:"bbs,omitempty"`

	/*DebugAddr - Descr: address at which to serve debug info Default: 0.0.0.0:17002
*/
	DebugAddr interface{} `yaml:"debug_addr,omitempty"`

}