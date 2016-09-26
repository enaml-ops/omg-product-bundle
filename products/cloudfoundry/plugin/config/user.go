package config

import (
	"github.com/enaml-ops/pluginlib/pluginutil"
	"gopkg.in/urfave/cli.v2"
)

func RequiredUserFlags() []string {
	return []string{
		"nats-user",
		"mysql-bootstrap-username",
		"cc-staging-upload-user",
		"db-ccdb-username",
		"db-uaa-username",
		"mysql-proxy-api-username",
		"db-console-username",
		"router-user",
		"cc-internal-api-user",
		"db-autoscale-username",
		"db-notifications-username",
	}
}

func NewUser(c *cli.Context) User {
	return User{
		NATSUser:              c.String("nats-user"),
		MySQLBootstrapUser:    c.String("mysql-bootstrap-username"),
		StagingUploadUser:     c.String("cc-staging-upload-user"),
		CCDBUsername:          c.String("db-ccdb-username"),
		UAADBUserName:         c.String("db-uaa-username"),
		AutoscaleDBUser:       c.String("db-autoscale-username"),
		MySQLProxyAPIUsername: c.String("mysql-proxy-api-username"),
		ConsoleDBUserName:     c.String("db-console-username"),
		RouterUser:            c.String("router-user"),
		CCInternalAPIUser:     c.String("cc-internal-api-user"),
		AutoscaleBrokerUser:   pluginutil.NewPassword(16),
		NotificationsDBUser:   c.String("db-notifications-username"),
	}
}

type User struct {
	CCInternalAPIUser     string
	MySQLBootstrapUser    string
	NATSUser              string
	StagingUploadUser     string
	CCDBUsername          string
	ConsoleDBUserName     string
	MySQLProxyAPIUsername string
	UAADBUserName         string
	AutoscaleDBUser       string
	RouterUser            string
	AutoscaleBrokerUser   string
	NotificationsDBUser   string
}
