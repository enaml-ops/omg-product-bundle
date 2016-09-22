package pmysql

const (
	innodbBufferPoolSize       = 2147483648
	maxConnections             = 1500
	databaseStartupTimeout     = 600
	wsrepDebug                 = 1
	backupServerPort           = 8081
	seededDBUser               = "repcanary"
	seededDBName               = "canary_db"
	adminUsername              = "root"
	notificationClientUsername = "mysql-monitoring"
	natsPort                   = "4222"
	natsUser                   = "nats"
	switchboardCount           = 2
	pollFrequency              = 30
	writeReadDelay             = 20
	brokerQuotaPause           = 30
	brokerPersistentDisk       = 102400
	registrarUser              = "admin"
)
