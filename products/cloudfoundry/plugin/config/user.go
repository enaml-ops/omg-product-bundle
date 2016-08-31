package config

import "github.com/codegangsta/cli"

func RequiredUserFlags() []string {
	return []string{
		"nats-user",
		"mysql-bootstrap-username",
		"cc-staging-upload-user",
		"cc-bulk-api-user",
		"db-ccdb-username",
		"db-uaa-username",
		"mysql-proxy-api-username",
		"db-console-username",
		"router-user",
		"cc-internal-api-user",
	}
}

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
		CCInternalAPIUser:     c.String("cc-internal-api-user"),
	}
}

type User struct {
	CCInternalAPIUser     string
	MySQLBootstrapUser    string
	NATSUser              string
	CCBulkAPIUser         string
	StagingUploadUser     string
	CCDBUsername          string
	ConsoleDBUserName     string
	MySQLProxyAPIUsername string
	UAADBUserName         string
	RouterUser            string
}
