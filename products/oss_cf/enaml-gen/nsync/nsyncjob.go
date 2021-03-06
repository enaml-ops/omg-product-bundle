package nsync 
/*
* File Generated by enaml generator
* !!! Please do not edit this file !!!
*/
type NsyncJob struct {

	/*Capi - Descr: Whether or not to use privileged containers for  buildpack based LRPs and tasks. Containers with a docker-image-based rootfs will continue to always be unprivileged and cannot be changed. Default: false
*/
	Capi *Capi `yaml:"capi,omitempty"`

	/*Diego - Descr: when connecting over https, ignore bad ssl certificates Default: false
*/
	Diego *Diego `yaml:"diego,omitempty"`

}