package nsync 
/*
* File Generated by enaml generator
* !!! Please do not edit this file !!!
*/
type Capi struct {

	/*Nsync - Descr: Whether or not to use privileged containers for  buildpack based LRPs and tasks. Containers with a docker-image-based rootfs will continue to always be unprivileged and cannot be changed. Default: false
*/
	Nsync *Nsync `yaml:"nsync,omitempty"`

}