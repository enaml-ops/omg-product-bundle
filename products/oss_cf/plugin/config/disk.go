package config

import "gopkg.in/urfave/cli.v2"

func RequiredDiskFlags() []string {
	return []string{"diego-cell-disk-type", "diego-brain-disk-type", "diego-db-disk-type", "nfs-disk-type", "etcd-disk-type", "mysql-disk-type"}
}

func NewDisk(c *cli.Context) Disk {
	return Disk{
		DiegoCellPersistentDiskType:  c.String("diego-cell-disk-type"),
		DiegoBrainPersistentDiskType: c.String("diego-brain-disk-type"),
		DiegoDBPersistentDiskType:    c.String("diego-db-disk-type"),
		NFSPersistentDiskType:        c.String("nfs-disk-type"),
		EtcdPersistentDiskType:       c.String("etcd-disk-type"),
		MySQLPersistentDiskType:      c.String("mysql-disk-type"),
	}
}

type Disk struct {
	EtcdPersistentDiskType       string
	MySQLPersistentDiskType      string
	NFSPersistentDiskType        string
	DiegoDBPersistentDiskType    string
	DiegoCellPersistentDiskType  string
	DiegoBrainPersistentDiskType string
}
