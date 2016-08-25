package bbs 
/*
* File Generated by enaml generator
* !!! Please do not edit this file !!!
*/
type Bbs struct {

	/*ListenAddr - Descr: address at which to serve API requests Default: 0.0.0.0:8889
*/
	ListenAddr interface{} `yaml:"listen_addr,omitempty"`

	/*ServerCert - Descr: PEM-encoded client certificate Default: <nil>
*/
	ServerCert interface{} `yaml:"server_cert,omitempty"`

	/*RequireSsl - Descr: require ssl for all communication the bbs Default: true
*/
	RequireSsl interface{} `yaml:"require_ssl,omitempty"`

	/*CaCert - Descr: PEM-encoded CA certificate Default: <nil>
*/
	CaCert interface{} `yaml:"ca_cert,omitempty"`

	/*Auctioneer - Descr: Address of the auctioneer API Default: http://auctioneer.service.cf.internal:9016
*/
	Auctioneer *Auctioneer `yaml:"auctioneer,omitempty"`

	/*ActiveKeyLabel - Descr: Label of the encryption key to be used when writing to the database Default: <nil>
*/
	ActiveKeyLabel interface{} `yaml:"active_key_label,omitempty"`

	/*AdvertisementBaseHostname - Descr: Suffix for the BBS advertised hostname Default: bbs.service.cf.internal
*/
	AdvertisementBaseHostname interface{} `yaml:"advertisement_base_hostname,omitempty"`

	/*DesiredLrpCreationTimeout - Descr: expected maximum time to create all components of a desired LRP Default: <nil>
*/
	DesiredLrpCreationTimeout interface{} `yaml:"desired_lrp_creation_timeout,omitempty"`

	/*DebugAddr - Descr: address at which to serve debug info Default: 0.0.0.0:17017
*/
	DebugAddr interface{} `yaml:"debug_addr,omitempty"`

	/*Sql - Descr: EXPERIMENTAL: connection string to use for SQL backend [username:password@tcp(1.1.1.1:1234)/database] Default: <nil>
*/
	Sql *Sql `yaml:"sql,omitempty"`

	/*EncryptionKeys - Descr: List of encryption keys to be used Default: []
*/
	EncryptionKeys interface{} `yaml:"encryption_keys,omitempty"`

	/*ServerKey - Descr: PEM-encoded client key Default: <nil>
*/
	ServerKey interface{} `yaml:"server_key,omitempty"`

	/*DropsondePort - Descr: local metron agent's port Default: 3457
*/
	DropsondePort interface{} `yaml:"dropsonde_port,omitempty"`

	/*LogLevel - Descr: Log level Default: info
*/
	LogLevel interface{} `yaml:"log_level,omitempty"`

	/*Etcd - Descr: PEM-encoded CA certificate Default: <nil>
*/
	Etcd *Etcd `yaml:"etcd,omitempty"`

}