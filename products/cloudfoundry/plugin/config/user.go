package config

import "github.com/codegangsta/cli"

func NewUser(c *cli.Context) User {
	return User{
		NATSUser:              c.String("nats-user"),
		MySQLBootstrapUser:    c.String("mysql-bootstrap-username"),
		StagingUploadUser:     c.String("cc-staging-upload-user"),
		CCBulkAPIUser:         c.String("cc-bulk-api-user"),
		CCDBUsername:          c.String("db-ccdb-username"),
		UAADBUserName:         c.String("db-uaa-username"),
		MySQLProxyAPIUsername: c.String("mysql-proxy-api-username"),
		ConsoleDBUserName:     c.String("db-console-username"),
		RouterUser:            c.String("router-user"),
	}
}

type User struct {
	MySQLBootstrapUser    string
	NATSUser              string
	CCBulkAPIUser         string
	StagingUploadUser     string
	InternalAPIUser       string
	CCDBUsername          string
	ConsoleDBUserName     string
	MySQLProxyAPIUsername string
	UAADBUserName         string
	RouterUser            string
}
