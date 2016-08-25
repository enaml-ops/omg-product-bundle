package mysql 
/*
* File Generated by enaml generator
* !!! Please do not edit this file !!!
*/
type MysqlJob struct {

	/*InnodbBufferPoolSize - Descr: The size in bytes of the memory buffer InnoDB uses to cache data and indexes of its tables Default: <nil>
*/
	InnodbBufferPoolSize interface{} `yaml:"innodb_buffer_pool_size,omitempty"`

	/*WsrepMaxWsRows - Descr: Maximum permitted number of rows per writeset. Default: 131072
*/
	WsrepMaxWsRows interface{} `yaml:"wsrep_max_ws_rows,omitempty"`

	/*IbLogFileSize - Descr: Size of the ib_log_file used by innodb, in MB Default: 1024
*/
	IbLogFileSize interface{} `yaml:"ib_log_file_size,omitempty"`

	/*MaxConnections - Descr: Maximum total number of database connections for the node Default: 1500
*/
	MaxConnections interface{} `yaml:"max_connections,omitempty"`

	/*DatabaseStartupTimeout - Descr: How long the startup script waits for the database to come online (in seconds) Default: <nil>
*/
	DatabaseStartupTimeout interface{} `yaml:"database_startup_timeout,omitempty"`

	/*WsrepMaxWsSize - Descr: Maximum permitted size in bytes per writeset. Default: 1073741824
*/
	WsrepMaxWsSize interface{} `yaml:"wsrep_max_ws_size,omitempty"`

	/*BootstrapEndpoint - Descr: Password used by the bootstrap endpoints for Basic Auth Default: <nil>
*/
	BootstrapEndpoint *BootstrapEndpoint `yaml:"bootstrap_endpoint,omitempty"`

	/*AdminUsername - Descr: Username for the MySQL server admin user Default: root
*/
	AdminUsername interface{} `yaml:"admin_username,omitempty"`

	/*AdminPassword - Descr: Password for the MySQL server admin user Default: <nil>
*/
	AdminPassword interface{} `yaml:"admin_password,omitempty"`

	/*ClusterIps - Descr: List of nodes.  Must have the same number of ips as there are nodes in the cluster Default: <nil>
*/
	ClusterIps interface{} `yaml:"cluster_ips,omitempty"`

	/*RoadminPassword - Descr: Password for the MySQL server read-only admin user Default: <nil>
*/
	RoadminPassword interface{} `yaml:"roadmin_password,omitempty"`

	/*MaxHeapTableSize - Descr: The maximum size (in rows) to which user-created MEMORY tables are permitted to grow Default: 16777216
*/
	MaxHeapTableSize interface{} `yaml:"max_heap_table_size,omitempty"`

	/*Port - Descr: Port the mysql server should bind to Default: 3306
*/
	Port interface{} `yaml:"port,omitempty"`

	/*GcacheSize - Descr: Cache size used by galera (maximum amount of data possible in an IST), in MB Default: 512
*/
	GcacheSize interface{} `yaml:"gcache_size,omitempty"`

	/*SeededDatabases - Descr: Set of databases to seed Default: map[]
*/
	SeededDatabases interface{} `yaml:"seeded_databases,omitempty"`

	/*RoadminEnabled - Descr: Whether read only user is enabled Default: false
*/
	RoadminEnabled interface{} `yaml:"roadmin_enabled,omitempty"`

	/*NetworkName - Descr: The name of the network (needed for the syslog aggregator) Default: <nil>
*/
	NetworkName interface{} `yaml:"network_name,omitempty"`

	/*SyslogAggregator - Descr: TCP port of syslog aggregator Default: <nil>
*/
	SyslogAggregator *SyslogAggregator `yaml:"syslog_aggregator,omitempty"`

	/*TmpTableSize - Descr: The maximum size (in bytes) of internal in-memory temporary tables Default: 33554432
*/
	TmpTableSize interface{} `yaml:"tmp_table_size,omitempty"`

	/*SkipNameResolve - Descr: Do not restrict connections to database based on hostname Default: false
*/
	SkipNameResolve interface{} `yaml:"skip_name_resolve,omitempty"`

	/*HealthcheckPort - Descr: Port used by healthcheck process to listen on Default: 9200
*/
	HealthcheckPort interface{} `yaml:"healthcheck_port,omitempty"`

}