package config

import (
	"github.com/enaml-ops/pluginlib/pcli"
	"github.com/enaml-ops/pluginlib/pluginutil"
	"gopkg.in/urfave/cli.v2"
)

func NewUser(c *cli.Context) User {
	var u User
	pcli.UnmarshalFlags(&u, c)
	u.AutoscaleBrokerUser = pluginutil.NewPassword(16)
	return u
}

type User struct {
	CCInternalAPIUser     string `omg:"cc-internal-api-user"`
	MySQLBootstrapUser    string `omg:"mysql-bootstrap-username"`
	NATSUser              string `omg:"nats-user"`
	StagingUploadUser     string `omg:"cc-staging-upload-user"`
	CCDBUsername          string `omg:"db-ccdb-username"`
	ConsoleDBUserName     string `omg:"db-console-username"`
	MySQLProxyAPIUsername string `omg:"mysql-proxy-api-username"`
	UAADBUserName         string `omg:"db-uaa-username"`
	AutoscaleDBUser       string `omg:"db-autoscale-username"`
	RouterUser            string
	AutoscaleBrokerUser   string `omg:"-"`
	NotificationsDBUser   string `omg:"db-notifications-username"`
}
