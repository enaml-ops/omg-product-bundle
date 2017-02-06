package groundcrew 
/*
* File Generated by enaml generator
* !!! Please do not edit this file !!!
*/
type GroundcrewJob struct {

	/*DrainTimeout - Descr: Maximum wait time in Go duration format (1m = 1 minute) for worker drain
to be finished. Only applies when worker is getting shutdown.

If not specified, it will be indefinite.
 Default: <nil>
*/
	DrainTimeout interface{} `yaml:"drain_timeout,omitempty"`

	/*AdditionalResourceTypes - Descr: Additional resource types supported by the workers in `[{type: string, image: string}]` format.
 Default: []
*/
	AdditionalResourceTypes interface{} `yaml:"additional_resource_types,omitempty"`

	/*NoProxy - Descr: A list domains and IPs with optional port for which the proxy should be bypassed, e.g. [localhost, 127.0.0.1, example.com, domain.com:8080]
 Default: []
*/
	NoProxy interface{} `yaml:"no_proxy,omitempty"`

	/*Tsa - Descr: Port of the TSA server to register with.

Only used when `tsa.host` is also specified. Otherwise the port is
autodiscovered via the `tsa` link.
 Default: 2222
*/
	Tsa *Tsa `yaml:"tsa,omitempty"`

	/*Tags - Descr: An array of tags to advertise for each worker.
 Default: []
*/
	Tags interface{} `yaml:"tags,omitempty"`

	/*Garden - Descr: Garden server connection address to forward through SSH to the TSA.

If not specified, the Garden server address is registered directly.
 Default: <nil>
*/
	Garden *Garden `yaml:"garden,omitempty"`

	/*HttpsProxyUrl - Descr: Proxy to use for outgoing https requests from containers.
 Default: <nil>
*/
	HttpsProxyUrl interface{} `yaml:"https_proxy_url,omitempty"`

	/*Platform - Descr: Platform to advertise for each worker.
 Default: linux
*/
	Platform interface{} `yaml:"platform,omitempty"`

	/*Team - Descr: Register the worker for a single team.

If not specified, the worker will be shared across all teams.
 Default: 
*/
	Team interface{} `yaml:"team,omitempty"`

	/*HttpProxyUrl - Descr: Proxy to use for outgoing http requests from containers.
 Default: <nil>
*/
	HttpProxyUrl interface{} `yaml:"http_proxy_url,omitempty"`

	/*Yeller - Descr: API key to output errors from Concourse to Yeller.
 Default: 
*/
	Yeller *Yeller `yaml:"yeller,omitempty"`

	/*Baggageclaim - Descr: Baggageclaim server URL to advertise directly to the
TSA.

If not specified, either the `baggageclaim` link is
used.
 Default: <nil>
*/
	Baggageclaim *Baggageclaim `yaml:"baggageclaim,omitempty"`

}