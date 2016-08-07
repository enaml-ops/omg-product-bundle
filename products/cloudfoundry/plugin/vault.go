package cloudfoundry

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"math/big"
	mrand "math/rand"
	"net"
	"os"
	"strings"
	"time"

	"github.com/enaml-ops/pluginlib/pcli"
	"github.com/enaml-ops/pluginlib/util"
	"github.com/xchapter7x/lo"
)

func VaultRotate(args []string, flgs []pcli.Flag) error {
	var err error
	c := pluginutil.NewContext(args, pluginutil.ToCliFlagArray(flgs))

	if c.Bool("vault-rotate") && hasValidVaultFlags(c) && c.String("system-domain") != "" {
		lo.G.Debug("rotating your vault values")
		vault := pluginutil.NewVaultUnmarshal(c.String("vault-domain"), c.String("vault-token"), pluginutil.DefaultClient())

		lo.G.Debug("rotating password values")
		if err = RotatePasswordHash(vault, c.String("vault-hash-password")); err == nil {
			lo.G.Debug("rotating keycert values")
			err = RotateCertHash(vault, c.String("vault-hash-keycert"), c.String("system-domain"))
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

func RotateCertHash(vault VaultRotater, hash string, host string) error {
	var err error
	secrets := getKeyCertObject(host)

	if err = vault.RotateSecrets(hash, secrets); err != nil {
		lo.G.Errorf("error updating hash: %v", err.Error())
	}
	return err
}

func randomString(strlen int) string {
	mrand.Seed(time.Now().UTC().UnixNano())
	const chars = "abcdefghipqrstuvwxyz0123456789"
	result := make([]byte, strlen)
	for i := 0; i < strlen; i++ {
		result[i] = chars[mrand.Intn(len(chars))]
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

func getKeyCertObject(host string) []byte {
	keysuffix := "-key"
	certsuffix := "-cert"
	fieldnames := []string{
		"router-ssl",
		"consul-agent",
		"consul-server",
		"bbs-server",
		"etcd-server",
		"etcd-client",
		"etcd-peer",
		"bbs-client",
	}
	cafieldnames := []string{
		"consul-ca-cert",
		"bbs-ca-cert",
	}
	var certVault map[string]string = make(map[string]string)

	for _, fn := range fieldnames {
		ca := false
		curve := ""
		cert, key := certgen(&host, &ca, &curve)
		certVault[fn+certsuffix] = cert
		certVault[fn+keysuffix] = key
	}

	for _, fn := range cafieldnames {
		ca := true
		curve := ""
		certVault[fn+certsuffix], _ = certgen(&host, &ca, &curve)
	}
	b, _ := json.Marshal(certVault)
	return b
}

func publicKey(priv interface{}) interface{} {
	switch k := priv.(type) {
	case *rsa.PrivateKey:
		return &k.PublicKey
	case *ecdsa.PrivateKey:
		return &k.PublicKey
	default:
		return nil
	}
}

func pemBlockForKey(priv interface{}) *pem.Block {
	switch k := priv.(type) {
	case *rsa.PrivateKey:
		return &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(k)}
	case *ecdsa.PrivateKey:
		b, err := x509.MarshalECPrivateKey(k)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to marshal ECDSA private key: %v", err)
			os.Exit(2)
		}
		return &pem.Block{Type: "EC PRIVATE KEY", Bytes: b}
	default:
		return nil
	}
}

func certgen(host *string, isCA *bool, ecdsaCurve *string) (cert string, key string) {
	validFor := 365 * 24 * time.Hour
	rsaBits := func() *int {
		v := 2048
		return &v
	}()

	if len(*host) == 0 {
		lo.G.Fatalf("Missing required --host parameter")
	}
	var priv interface{}
	var err error
	switch *ecdsaCurve {
	case "":
		priv, err = rsa.GenerateKey(rand.Reader, *rsaBits)
	case "P224":
		priv, err = ecdsa.GenerateKey(elliptic.P224(), rand.Reader)
	case "P256":
		priv, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	case "P384":
		priv, err = ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	case "P521":
		priv, err = ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	default:
		fmt.Fprintf(os.Stderr, "Unrecognized elliptic curve: %q", *ecdsaCurve)
		os.Exit(1)
	}
	if err != nil {
		lo.G.Fatalf("failed to generate private key: %s", err)
	}

	var notBefore time.Time
	notBefore = time.Now()
	notAfter := notBefore.Add(validFor)

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		lo.G.Fatalf("failed to generate serial number: %s", err)
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Acme Co"},
		},
		NotBefore: notBefore,
		NotAfter:  notAfter,

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	hosts := strings.Split(*host, ",")
	for _, h := range hosts {
		if ip := net.ParseIP(h); ip != nil {
			template.IPAddresses = append(template.IPAddresses, ip)
		} else {
			template.DNSNames = append(template.DNSNames, h)
		}
	}

	if *isCA {
		template.IsCA = true
		template.KeyUsage |= x509.KeyUsageCertSign
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, publicKey(priv), priv)
	if err != nil {
		lo.G.Fatalf("Failed to create certificate: %s", err)
	}
	certOut := bytes.NewBufferString("")
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	lo.G.Debug("written cert")
	keyOut := bytes.NewBufferString("")
	pem.Encode(keyOut, pemBlockForKey(priv))
	lo.G.Debug("written key")
	return certOut.String(), keyOut.String()
}
