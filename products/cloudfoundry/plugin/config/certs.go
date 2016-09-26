package config

import (
	clictx "github.com/enaml-ops/omg-product-bundle/cli"
	"gopkg.in/urfave/cli.v2"
)

func RequiredCertFlags() []string {
	return []string{
		"consul-agent-cert", "consul-agent-key", "consul-server-cert", "consul-server-key",
		"bbs-server-ca-cert", "bbs-client-cert", "bbs-client-key", "bbs-server-cert", "bbs-server-key",
		"etcd-server-cert", "etcd-server-key", "etcd-client-cert", "etcd-client-key", "etcd-peer-cert", "etcd-peer-key",
		"uaa-saml-service-provider-key", "uaa-saml-service-provider-cert", "uaa-jwt-signing-key", "uaa-jwt-verification-key",
		"router-ssl-cert", "router-ssl-key",
	}
}

func NewCerts(c *cli.Context) (*Certs, error) {
	certs := &Certs{}

	agentCert, err := clictx.LoadResourceFromContext(c, "consul-agent-cert")
	if err != nil {
		return nil, err
	}
	agentKey, err := clictx.LoadResourceFromContext(c, "consul-agent-key")
	if err != nil {
		return nil, err
	}
	serverCert, err := clictx.LoadResourceFromContext(c, "consul-server-cert")
	if err != nil {
		return nil, err
	}
	serverKey, err := clictx.LoadResourceFromContext(c, "consul-server-key")
	if err != nil {
		return nil, err
	}
	bbsCaCert, err := clictx.LoadResourceFromContext(c, "bbs-server-ca-cert")
	if err != nil {
		return nil, err
	}

	bbsClientCert, err := clictx.LoadResourceFromContext(c, "bbs-client-cert")
	if err != nil {
		return nil, err
	}

	bbsClientKey, err := clictx.LoadResourceFromContext(c, "bbs-client-key")
	if err != nil {
		return nil, err
	}

	bbsServerCert, err := clictx.LoadResourceFromContext(c, "bbs-server-cert")
	if err != nil {
		return nil, err
	}

	bbsServerKey, err := clictx.LoadResourceFromContext(c, "bbs-server-key")
	if err != nil {
		return nil, err
	}

	etcdServerCert, err := clictx.LoadResourceFromContext(c, "etcd-server-cert")
	if err != nil {
		return nil, err
	}

	etcdServerKey, err := clictx.LoadResourceFromContext(c, "etcd-server-key")
	if err != nil {
		return nil, err
	}

	etcdClientCert, err := clictx.LoadResourceFromContext(c, "etcd-client-cert")
	if err != nil {
		return nil, err
	}

	etcdClientKey, err := clictx.LoadResourceFromContext(c, "etcd-client-key")
	if err != nil {
		return nil, err
	}

	etcdPeerCert, err := clictx.LoadResourceFromContext(c, "etcd-peer-cert")
	if err != nil {
		return nil, err
	}

	etcdPeerKey, err := clictx.LoadResourceFromContext(c, "etcd-peer-key")
	if err != nil {
		return nil, err
	}

	sslpem, err := clictx.LoadResourceFromContext(c, "haproxy-sslpem")
	if err != nil {
		return nil, err
	}

	routerCert, err := clictx.LoadResourceFromContext(c, "router-ssl-cert")
	if err != nil {
		return nil, err
	}
	routerKey, err := clictx.LoadResourceFromContext(c, "router-ssl-key")
	if err != nil {
		return nil, err
	}

	samlKey, err := clictx.LoadResourceFromContext(c, "uaa-saml-service-provider-key")
	if err != nil {
		return nil, err
	}

	samlCert, err := clictx.LoadResourceFromContext(c, "uaa-saml-service-provider-cert")
	if err != nil {
		return nil, err
	}

	jwtSigningKey, err := clictx.LoadResourceFromContext(c, "uaa-jwt-signing-key")
	if err != nil {
		return nil, err
	}

	jwtVerificationKey, err := clictx.LoadResourceFromContext(c, "uaa-jwt-verification-key")
	if err != nil {
		return nil, err
	}

	certs.SAMLServiceProviderCertificate = samlCert
	certs.JWTSigningKey = jwtSigningKey
	certs.JWTVerificationKey = jwtVerificationKey
	certs.SAMLServiceProviderKey = samlKey
	certs.RouterSSLCert = routerCert
	certs.RouterSSLKey = routerKey
	certs.ConsulAgentCert = agentCert
	certs.ConsulServerCert = serverCert
	certs.ConsulAgentKey = agentKey
	certs.ConsulServerKey = serverKey
	certs.BBSCACert = bbsCaCert
	certs.BBSClientCert = bbsClientCert
	certs.BBSClientKey = bbsClientKey
	certs.BBSServerCert = bbsServerCert
	certs.BBSServerKey = bbsServerKey
	certs.EtcdClientCert = etcdClientCert
	certs.EtcdClientKey = etcdClientKey
	certs.EtcdPeerCert = etcdPeerCert
	certs.EtcdPeerKey = etcdPeerKey
	certs.EtcdServerKey = etcdServerKey
	certs.EtcdServerCert = etcdServerCert
	certs.HAProxySSLPem = sslpem

	return certs, nil
}

type Certs struct {
	EtcdServerCert                 string
	EtcdServerKey                  string
	EtcdClientCert                 string
	EtcdClientKey                  string
	EtcdPeerCert                   string
	EtcdPeerKey                    string
	BBSCACert                      string
	BBSClientCert                  string
	BBSClientKey                   string
	BBSServerCert                  string
	BBSServerKey                   string
	RouterSSLCert                  string
	RouterSSLKey                   string
	HAProxySSLPem                  string
	ConsulAgentCert                string
	ConsulAgentKey                 string
	ConsulServerCert               string
	ConsulServerKey                string
	JWTVerificationKey             string
	SAMLServiceProviderKey         string
	SAMLServiceProviderCertificate string
	JWTSigningKey                  string
}
