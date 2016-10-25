package consul_agent 
/*
* File Generated by enaml generator
* !!! Please do not edit this file !!!
*/
type Consul struct {

	/*CaCert - Descr: PEM-encoded CA certificate Default: <nil>
*/
	CaCert interface{} `yaml:"ca_cert,omitempty"`

	/*EncryptKeys - Descr: A list of passphrases that will be converted into encryption keys, the first key in the list is the active one Default: <nil>
*/
	EncryptKeys interface{} `yaml:"encrypt_keys,omitempty"`

	/*ServerCert - Descr: PEM-encoded server certificate Default: <nil>
*/
	ServerCert interface{} `yaml:"server_cert,omitempty"`

	/*AgentKey - Descr: PEM-encoded client key Default: <nil>
*/
	AgentKey interface{} `yaml:"agent_key,omitempty"`

	/*ServerKey - Descr: PEM-encoded server key Default: <nil>
*/
	ServerKey interface{} `yaml:"server_key,omitempty"`

	/*AgentCert - Descr: PEM-encoded agent certificate Default: <nil>
*/
	AgentCert interface{} `yaml:"agent_cert,omitempty"`

	/*Agent - Descr: Agent log level. Default: info
*/
	Agent *Agent `yaml:"agent,omitempty"`

}