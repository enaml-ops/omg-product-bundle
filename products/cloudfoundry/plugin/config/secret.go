package config

import (
	"github.com/enaml-ops/omg-cli/utils"
	"gopkg.in/urfave/cli.v2"
)

func RequiredSecretFlags() []string {
	return []string{
		"db-uaa-password",
		"push-apps-manager-password",
		"system-services-password",
		"system-verification-password",
		"opentsdb-firehose-nozzle-client-secret",
		"identity-client-secret",
		"login-client-secret",
		"portal-client-secret",
		"autoscaling-service-client-secret",
		"system-passwords-client-secret",
		"cc-service-dashboards-client-secret",
		"gorouter-client-secret",
		"notifications-client-secret",
		"notifications-ui-client-secret",
		"cloud-controller-username-lookup-client-secret",
		"cc-routing-client-secret",
		"apps-metrics-client-secret",
		"apps-metrics-processing-client-secret",
		"admin-password",
		"nats-pass",
		"mysql-bootstrap-password",
		"consul-encryption-key",
		"smoke-tests-password",
		"doppler-shared-secret",
		"doppler-client-secret",
		"cc-bulk-api-password",
		"cc-internal-api-password",
		"ssh-proxy-uaa-secret",
		"cc-db-encryption-key",
		"db-ccdb-password",
		"diego-db-passphrase",
		"uaa-admin-secret",
		"router-pass",
		"mysql-proxy-api-password",
		"mysql-admin-password",
		"db-console-password",
		"cc-staging-upload-password",
		"db-app_usage-password",
		"apps-manager-secret-token",
		"db-autoscale-password",
		"db-notifications-password",
	}
}

func NewSecret(c *cli.Context) Secret {
	return Secret{
		UAADBPassword:                             c.String("db-uaa-password"),
		PushAppsManagerPassword:                   c.String("push-apps-manager-password"),
		SystemServicesPassword:                    c.String("system-services-password"),
		SystemVerificationPassword:                c.String("system-verification-password"),
		OpentsdbFirehoseNozzleClientSecret:        c.String("opentsdb-firehose-nozzle-client-secret"),
		IdentityClientSecret:                      c.String("identity-client-secret"),
		LoginClientSecret:                         c.String("login-client-secret"),
		PortalClientSecret:                        c.String("portal-client-secret"),
		AutoScalingServiceClientSecret:            c.String("autoscaling-service-client-secret"),
		SystemPasswordsClientSecret:               c.String("system-passwords-client-secret"),
		CCServiceDashboardsClientSecret:           c.String("cc-service-dashboards-client-secret"),
		GoRouterClientSecret:                      c.String("gorouter-client-secret"),
		NotificationsClientSecret:                 c.String("notifications-client-secret"),
		NotificationsUIClientSecret:               c.String("notifications-ui-client-secret"),
		CloudControllerUsernameLookupClientSecret: c.String("cloud-controller-username-lookup-client-secret"),
		CCRoutingClientSecret:                     c.String("cc-routing-client-secret"),
		AppsMetricsClientSecret:                   c.String("apps-metrics-client-secret"),
		AppsMetricsProcessingClientSecret:         c.String("apps-metrics-processing-client-secret"),
		AdminPassword:                             c.String("admin-password"),
		NATSPassword:                              c.String("nats-pass"),
		MySQLBootstrapPassword:                    c.String("mysql-bootstrap-password"),
		ConsulEncryptKeys:                         c.StringSlice("consul-encryption-key"),
		SmokeTestsPassword:                        c.String("smoke-tests-password"),
		DopplerSharedSecret:                       c.String("doppler-shared-secret"),
		DopplerSecret:                             c.String("doppler-client-secret"),
		CCBulkAPIPassword:                         c.String("cc-bulk-api-password"),
		CCInternalAPIPassword:                     c.String("cc-internal-api-password"),
		SSHProxyClientSecret:                      c.String("ssh-proxy-uaa-secret"),
		DbEncryptionKey:                           c.String("cc-db-encryption-key"),
		CCDBPassword:                              c.String("db-ccdb-password"),
		DiegoDBPassphrase:                         c.String("diego-db-passphrase"),
		AdminSecret:                               c.String("uaa-admin-secret"),
		LDAPUserPassword:                          c.String("uaa-ldap-user-password"),
		RouterPass:                                c.String("router-pass"),
		MySQLProxyAPIPassword:                     c.String("mysql-proxy-api-password"),
		MySQLAdminPassword:                        c.String("mysql-admin-password"),
		ConsoleDBPassword:                         c.String("db-console-password"),
		StagingUploadPassword:                     c.String("cc-staging-upload-password"),
		AppUsageDBPassword:                        c.String("db-app_usage-password"),
		AppsManagerSecretToken:                    c.String("apps-manager-secret-token"),
		AutoscaleBrokerPassword:                   utils.NewPassword(16),
		AutoscaleDBPassword:                       c.String("db-autoscale-password"),
	}
}

type Secret struct {
	AdminPassword                             string
	MySQLBootstrapPassword                    string
	NATSPassword                              string
	SmokeTestsPassword                        string
	DopplerSecret                             string
	DopplerSharedSecret                       string
	CCBulkAPIPassword                         string
	CCInternalAPIPassword                     string
	SSHProxyClientSecret                      string
	DiegoDBPassphrase                         string
	AdminSecret                               string
	UAADBPassword                             string
	PushAppsManagerPassword                   string
	SystemServicesPassword                    string
	SystemVerificationPassword                string
	OpentsdbFirehoseNozzleClientSecret        string
	IdentityClientSecret                      string
	LoginClientSecret                         string
	PortalClientSecret                        string
	AutoScalingServiceClientSecret            string
	SystemPasswordsClientSecret               string
	CCServiceDashboardsClientSecret           string
	GoRouterClientSecret                      string
	NotificationsClientSecret                 string
	NotificationsUIClientSecret               string
	CloudControllerUsernameLookupClientSecret string
	CCRoutingClientSecret                     string
	AppsMetricsClientSecret                   string
	AppsMetricsProcessingClientSecret         string
	LDAPUserPassword                          string
	DbEncryptionKey                           string
	CCDBPassword                              string
	StagingUploadPassword                     string
	MySQLProxyAPIPassword                     string
	RouterPass                                string
	MySQLAdminPassword                        string
	ConsoleDBPassword                         string
	ConsulEncryptKeys                         []string
	AppUsageDBPassword                        string
	AppsManagerSecretToken                    string
	AutoscaleBrokerPassword                   string
	AutoscaleDBPassword                       string
	NotificationsDBPassword                   string
}
