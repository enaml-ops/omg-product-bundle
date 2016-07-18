package debian_nfs_server 
/*
* File Generated by enaml generator
* !!! Please do not edit this file !!!
*/
type NfsServer struct {

	/*PipefsDirectory - Descr: Pipefs directory for NFS idmapd Default: /var/lib/nfs/rpc_pipef
*/
	PipefsDirectory interface{} `yaml:"pipefs_directory,omitempty"`

	/*NoRootSquash - Descr: Exports /var/vcap/store with no_root_squash when set to true Default: false
*/
	NoRootSquash interface{} `yaml:"no_root_squash,omitempty"`

	/*AllowFromEntries - Descr: An array of Hosts, Domains, Wildcard Domains, CIDR Networks and/or IPs from which /var/vcap/store is accessible Default: <nil>
*/
	AllowFromEntries interface{} `yaml:"allow_from_entries,omitempty"`

	/*IdmapdDomain - Descr: Domain name for NFS idmapd Default: localdomain
*/
	IdmapdDomain interface{} `yaml:"idmapd_domain,omitempty"`

}