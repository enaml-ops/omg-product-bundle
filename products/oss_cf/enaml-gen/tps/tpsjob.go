package tps 
/*
* File Generated by enaml generator
* !!! Please do not edit this file !!!
*/
type TpsJob struct {

	/*Diego - Descr: when connecting over https, ignore bad ssl certificates Default: false
*/
	Diego *Diego `yaml:"diego,omitempty"`

	/*Capi - Descr: Basic auth username for CC internal API Default: internal_user
*/
	Capi *Capi `yaml:"capi,omitempty"`

}