package debian_nfs_server 
/*
* File Generated by enaml generator
* !!! Please do not edit this file !!!
*/
type DebianNfsServerJob struct {

	/*NfsServer - Descr: Exports /var/vcap/store with no_root_squash when set to true Default: false
*/
	NfsServer *NfsServer `yaml:"nfs_server,omitempty"`

}