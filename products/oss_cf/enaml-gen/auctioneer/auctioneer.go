package auctioneer 
/*
* File Generated by enaml generator
* !!! Please do not edit this file !!!
*/
type Auctioneer struct {

	/*Bbs - Descr: maximum number of idle http connections Default: <nil>
*/
	Bbs *Bbs `yaml:"bbs,omitempty"`

	/*ListenAddr - Descr: address where auctioneer listens for LRP and task start auction requests Default: 0.0.0.0:9016
*/
	ListenAddr interface{} `yaml:"listen_addr,omitempty"`

	/*StartingContainerWeight - Descr: Factor to bias against cells with starting containers (0.0 - 1.0) Default: 0.25
*/
	StartingContainerWeight interface{} `yaml:"starting_container_weight,omitempty"`

	/*DebugAddr - Descr: address at which to serve debug info Default: 127.0.0.1:17001
*/
	DebugAddr interface{} `yaml:"debug_addr,omitempty"`

	/*LogLevel - Descr: Log level Default: info
*/
	LogLevel interface{} `yaml:"log_level,omitempty"`

	/*CellStateTimeout - Descr: Timeout applied to HTTP requests to the Cell State endpoint. Default: 1s
*/
	CellStateTimeout interface{} `yaml:"cell_state_timeout,omitempty"`

	/*DropsondePort - Descr: local metron agent's port Default: 3457
*/
	DropsondePort interface{} `yaml:"dropsonde_port,omitempty"`

}