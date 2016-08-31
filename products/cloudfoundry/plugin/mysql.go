package cloudfoundry

import (
	"github.com/enaml-ops/enaml"
	mysqllib "github.com/enaml-ops/omg-product-bundle/products/cf-mysql/enaml-gen/mysql"
)

//MySQL -
type MySQL struct {
	Config                 *Config
	DatabaseStartupTimeout int
	InnodbBufferPoolSize   int
	MaxConnections         int
	MySQLSeededDatabases   []MySQLSeededDatabase
}

//MySQLSeededDatabase -
type MySQLSeededDatabase struct {
	Name     string `yaml:"name"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

//NewMySQLPartition -
func NewMySQLPartition(config *Config) (igf InstanceGroupCreator) {
	igf = &MySQL{
		Config:                 config,
		DatabaseStartupTimeout: 1200,
		InnodbBufferPoolSize:   2147483648,
		MaxConnections:         1500,
		MySQLSeededDatabases:   parseSeededDBs(config),
	}
	return
}

func parseSeededDBs(config *Config) (dbs []MySQLSeededDatabase) {
	dbs = append(dbs, MySQLSeededDatabase{
		Name:     "uaa",
		Username: config.UAADBUserName,
		Password: config.UAADBPassword,
	})
	dbs = append(dbs, MySQLSeededDatabase{
		Name:     "ccdb",
		Username: config.CCDBUsername,
		Password: config.CCDBPassword,
	})
	dbs = append(dbs, MySQLSeededDatabase{
		Name:     "console",
		Username: config.ConsoleDBUserName,
		Password: config.ConsoleDBPassword,
	})

	return
}

// GetSeededDBByName returns a pointer to the seeded database with a particular
// name.  It returns null if no matching database is found.
func (s *MySQL) GetSeededDBByName(name string) *MySQLSeededDatabase {
	for i := range s.MySQLSeededDatabases {
		if s.MySQLSeededDatabases[i].Name == name {
			return &s.MySQLSeededDatabases[i]
		}
	}
	return nil
}

//ToInstanceGroup -
func (s *MySQL) ToInstanceGroup() (ig *enaml.InstanceGroup) {
	ig = &enaml.InstanceGroup{
		Name:               "mysql-partition",
		Instances:          len(s.Config.MySQLIPs),
		VMType:             s.Config.MySQLVMType,
		AZs:                s.Config.AZs,
		Stemcell:           s.Config.StemcellName,
		PersistentDiskType: s.Config.MySQLPersistentDiskType,
		Jobs: []enaml.InstanceJob{
			s.newMySQLJob(),
		},
		Networks: []enaml.Network{
			enaml.Network{Name: s.Config.NetworkName, StaticIPs: s.Config.MySQLIPs},
		},
		Update: enaml.Update{
			MaxInFlight: 1,
		},
	}
	return
}

func (s *MySQL) newMySQLJob() enaml.InstanceJob {
	return enaml.InstanceJob{
		Name:    "mysql",
		Release: "cf-mysql",
		Properties: &mysqllib.MysqlJob{
			AdminPassword:          s.Config.MySQLAdminPassword,
			ClusterIps:             s.Config.MySQLIPs,
			DatabaseStartupTimeout: s.DatabaseStartupTimeout,
			InnodbBufferPoolSize:   s.InnodbBufferPoolSize,
			MaxConnections:         s.MaxConnections,
			BootstrapEndpoint: &mysqllib.BootstrapEndpoint{
				Username: s.Config.MySQLBootstrapUser,
				Password: s.Config.MySQLBootstrapPassword,
			},
			SeededDatabases: s.MySQLSeededDatabases,
			SyslogAggregator: &mysqllib.SyslogAggregator{
				Address:   s.Config.SyslogAddress,
				Port:      s.Config.SyslogPort,
				Transport: s.Config.SyslogTransport,
			},
		},
	}
}
