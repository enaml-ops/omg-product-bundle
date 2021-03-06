package tsa 
/*
* File Generated by enaml generator
* !!! Please do not edit this file !!!
*/
type TsaJob struct {

	/*AuthorizedKeys - Descr: Public keys to authorize for SSH connections.
 Default: []
*/
	AuthorizedKeys interface{} `yaml:"authorized_keys,omitempty"`

	/*HeartbeatInterval - Descr: Interval on which to register workers with the ATC.
 Default: 30s
*/
	HeartbeatInterval interface{} `yaml:"heartbeat_interval,omitempty"`

	/*HostPublicKey - Descr: Public key component of the host's key. This property is exported via the `tsa` link so that workers can discover it.
 Default: 
*/
	HostPublicKey interface{} `yaml:"host_public_key,omitempty"`

	/*Yeller - Descr: If configured, errors emitted to the logs will also be emitted to Yeller.
This is only really useful for Concourse developers.
 Default: 
*/
	Yeller *Yeller `yaml:"yeller,omitempty"`

	/*TeamAuthorizedKeys - Descr: Public keys to authorize team workers for SSH connections.
 Default: []
*/
	TeamAuthorizedKeys interface{} `yaml:"team_authorized_keys,omitempty"`

	/*AuthorizeGeneratedWorkerKey - Descr: Permit access via generated worker key, local to the deployment. Set to
`false` if you plan on only ever using explicitly configured and
authorized worker keys.
 Default: true
*/
	AuthorizeGeneratedWorkerKey interface{} `yaml:"authorize_generated_worker_key,omitempty"`

	/*Atc - Descr: ATC API endpoint to which workers will be advertised.
If not specified, it will be autodiscovered via BOSH links.
 Default: http://127.0.0.1:8080
*/
	Atc *Atc `yaml:"atc,omitempty"`

	/*ForwardHost - Descr: Address to advertise forwarded worker connections to.

If not specified, the instance's address is used.
 Default: <nil>
*/
	ForwardHost interface{} `yaml:"forward_host,omitempty"`

	/*BindPort - Descr: Port on which to listen for SSH connections.
 Default: 2222
*/
	BindPort interface{} `yaml:"bind_port,omitempty"`

	/*HostKey - Descr: Private key to use for the SSH server.
If not specified, a deployment-scoped default is used.
 Default: 
*/
	HostKey interface{} `yaml:"host_key,omitempty"`

}