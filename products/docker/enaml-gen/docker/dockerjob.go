package docker 
/*
* File Generated by enaml generator
* !!! Please do not edit this file !!!
*/
type DockerJob struct {

	/*Docker - Descr: Array of log driver options Default: <nil>
*/
	Docker *Docker `yaml:"docker,omitempty"`

	/*Env - Descr: HTTP proxy that Docker should use Default: <nil>
*/
	Env *Env `yaml:"env,omitempty"`

}