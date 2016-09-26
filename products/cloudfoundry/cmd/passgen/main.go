package main

import (
	"encoding/json"
	"os"

	"github.com/enaml-ops/pluginlib/pluginutil"
)

const passLength = 20

func main() {
	fieldnames := []string{
		"router-pass",
		"nats-pass",
		"mysql-admin-password",
		"mysql-bootstrap-password",
		"mysql-proxy-api-password",
		"cc-staging-upload-password",
		"cc-bulk-api-password",
		"cc-internal-api-password",
		"db-uaa-password",
		"db-ccdb-password",
		"db-console-password",
		"diego-db-passphrase",
		"uaa-ldap-user-password",
		"admin-password",
		"push-apps-manager-password",
		"smoke-tests-password",
		"system-services-password",
		"system-verification-password",
		"system-passwords-client-secret",
		"doppler-shared-secret",
		"ssh-proxy-uaa-secret",
		"metron-secret",
		"uaa-admin-secret",
		"opentsdb-firehose-nozzle-client-secret",
		"identity-client-secret",
		"login-client-secret",
		"portal-client-secret",
		"autoscaling-service-client-secret",
		"cc-service-dashboards-client-secret",
		"doppler-client-secret",
		"gorouter-client-secret",
		"notifications-client-secret",
		"notifications-ui-client-secret",
		"cloud-controller-username-lookup-client-secret",
		"cc-routing-client-secret",
		"ssh-proxy-client-secret",
		"apps-metrics-client-secret",
		"apps-metrics-processing-client-secret",
		"consul-encryption-key",
		"cc-db-encryption-key",
		"host-key-fingerprint",
		"uaa-jwt-signing-key",
		"uaa-jwt-verification-key",
		"uaa-saml-service-provider-key",
	}

	passVault := make(map[string]string)
	for _, fn := range fieldnames {
		passVault[fn] = pluginutil.NewPassword(passLength)
	}
	json.NewEncoder(os.Stdout).Encode(passVault)
}
