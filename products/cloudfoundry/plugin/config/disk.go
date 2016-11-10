package config

// Disk contains the disk types for a Cloud Foundry deployment.
type Disk struct {
	EtcdPersistentDiskType       string `omg:"etcd-disk-type"`
	MySQLPersistentDiskType      string `omg:"mysql-disk-type"`
	NFSPersistentDiskType        string `omg:"nfs-disk-type"`
	DiegoDBPersistentDiskType    string `omg:"diego-db-disk-type"`
	DiegoCellPersistentDiskType  string `omg:"diego-cell-disk-type"`
	DiegoBrainPersistentDiskType string `omg:"diego-brain-disk-type"`
}
