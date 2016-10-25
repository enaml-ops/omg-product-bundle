package cloudfoundry

import (
	"encoding/json"

	"github.com/enaml-ops/pluginlib/pcli"
	"github.com/enaml-ops/pluginlib/pluginutil"
	"github.com/xchapter7x/lo"
)

func VaultRotate(args []string, flgs []pcli.Flag) error {
	var err error
	c := pluginutil.NewContext(args, pluginutil.ToCliFlagArray(flgs))

	if c.Bool("vault-rotate") && hasValidVaultFlags(c) && c.String("system-domain") != "" {
		lo.G.Debug("rotating your vault values")
		vault := pluginutil.NewVaultUnmarshal(c.String("vault-domain"), c.String("vault-token"))

		lo.G.Debug("rotating password values")
		if err = RotatePasswordHash(vault, c.String("vault-hash-password")); err == nil {
			lo.G.Debug("rotating keycert values")
			err = RotateCertHash(vault, c.String("vault-hash-keycert"), c.String("system-domain"), c.StringSlice("app-domain"))
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

	if err = vault.RotateSecrets(hash, secrets); err != nil {
		lo.G.Errorf("error updating hash: %v", err.Error())
	}
	return err
}

func RotateCertHash(vault VaultRotater, hash, systemDomain string, appsDomain []string) error {
	secrets, err := getKeyCertObject(systemDomain, appsDomain)
	if err != nil {
		return err
	}

	if err = vault.RotateSecrets(hash, secrets); err != nil {
		lo.G.Errorf("error updating hash: %v", err.Error())
	}
	return err
}

const passLength = 20

func getPasswordObject() []byte {
	fieldnames := []string{
		"cc-staging-upload-user",
		"cc-bulk-api-user",
		"cc-internal-api-user",
		"router-pass",
		"nats-pass",
		"mysql-admin-password",
		"mysql-bootstrap-password",
		"mysql-proxy-api-password",
		"cc-staging-upload-password",
		"cc-bulk-api-password",
		"cc-internal-api-password",
		"db-autoscale-password",
		"db-uaa-password",
		"db-ccdb-password",
		"db-console-password",
		"db-app_usage-password",
		"db-notifications-password",
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
		"doppler-zone",
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
		"apps-manager-secret-token",
	}

	passVault := make(map[string]string)
	for _, fn := range fieldnames {
		passVault[fn] = pluginutil.NewPassword(passLength)
	}
	b, _ := json.Marshal(passVault)
	return b
}

func getKeyCertObject(systemDomain string, appDomain []string) ([]byte, error) {
	const (
		keysuffix    = "-key"
		certsuffix   = "-cert"
		caCertSuffix = "-ca-cert"
	)

	type certGenerator struct{ flag, host string }

	fieldnames := []certGenerator{
		{"router-ssl", systemDomain},
		{"consul-agent", "consul_agent_cert"},
		{"consul-server", "server.dc1.cf.internal"},
		{"bbs-client", "bbs_client_cert"},
		{"bbs-server", "bbs.service.cf.internal"},
		{"etcd-server", "etcd.service.cf.internal"},
		{"etcd-client", "etcd_client_cert"},
		{"etcd-peer", "etcd.service.cf.internal"},
		{"uaa-saml-service-provider", "service_provider_key_credentials"},
	}

	certVault := make(map[string]string)
	caKey, caCert, err := pluginutil.Initialize()
	if err != nil {
		return nil, err
	}
	for _, fn := range fieldnames {
		ca, cert, key, err := pluginutil.GenerateCertWithCA([]string{fn.host, "*." + fn.host}, caCert, caKey)
		if err != nil {
			lo.G.Errorf("couldn't create cert for flag %s", fn.flag)
			return nil, err
		}
		certVault[fn.flag+certsuffix] = cert
		certVault[fn.flag+keysuffix] = key
		certVault[fn.flag+caCertSuffix] = ca
	}

	hosts := []string{
		systemDomain,
		"*." + systemDomain,
		"*.uaa." + systemDomain,
		"*.login." + systemDomain,
	}
	for _, ad := range appDomain {
		hosts = append(hosts, "*."+ad)
	}
	_, cert, key, err := pluginutil.GenerateCertWithCA(hosts, caCert, caKey)
	if err != nil {
		lo.G.Error("coudln't generate haproxy cert")
		return nil, err
	}
	certVault["haproxy-sslpem"] = cert + key

	jwtPublicKey, jwtPrivateKey, err := pluginutil.GenerateKeys()
	if err != nil {
		lo.G.Error("couldn't generate UAA JWT keys")
		return nil, err
	}
	certVault["uaa-jwt-signing-key"] = jwtPrivateKey
	certVault["uaa-jwt-verification-key"] = jwtPublicKey

	b, err := json.Marshal(certVault)
	return b, err
}
