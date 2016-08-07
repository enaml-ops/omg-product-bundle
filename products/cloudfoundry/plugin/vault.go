package cloudfoundry

import (
	"encoding/json"
	"math/rand"
	"time"

	"github.com/enaml-ops/pluginlib/pcli"
	"github.com/enaml-ops/pluginlib/util"
	"github.com/xchapter7x/lo"
)

func VaultRotate(args []string, flgs []pcli.Flag) error {
	var err error
	c := pluginutil.NewContext(args, pluginutil.ToCliFlagArray(flgs))

	if c.Bool("vault-rotate") && hasValidVaultFlags(c) {
		lo.G.Debug("rotating your vault values")
		vault := pluginutil.NewVaultUnmarshal(c.String("vault-domain"), c.String("vault-token"), pluginutil.DefaultClient())

		lo.G.Debug("rotating password values")
		if err = RotatePasswordHash(vault, c.String("vault-hash-password")); err == nil {
			lo.G.Debug("rotating keycert values")
			err = RotateCertHash(vault, c.String("vault-hash-keycert"))
		}
		lo.G.Debugf("checking respone from rotate: %v", err)

	} else {
		lo.G.Debug("we are not rotating vault values at this time")
	}
	return err
}

func RotatePasswordHash(vault VaultRotater, hash string) error {
	var err error
	secrets := getPasswordObject()
	lo.G.Debugf("secrets: %v", string(secrets))

	if err = vault.RotateSecrets(hash, secrets); err != nil {
		lo.G.Errorf("error updating hash: %v", err.Error())
	}
	return err
}

func RotateCertHash(vault VaultRotater, hash string) error {
	var secrets interface{}
	var err error

	if err = vault.RotateSecrets(hash, secrets); err != nil {
		lo.G.Errorf("error updating hash: %v", err.Error())
	}
	return err
}

func randomString(strlen int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	const chars = "abcdefghipqrstuvwxyz0123456789"
	result := make([]byte, strlen)
	for i := 0; i < strlen; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}

const passLength = 20

func getPasswordObject() []byte {
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
	var passVault map[string]string = make(map[string]string)

	for _, fn := range fieldnames {
		passVault[fn] = randomString(passLength)
	}
	b, _ := json.Marshal(passVault)
	return b
}
