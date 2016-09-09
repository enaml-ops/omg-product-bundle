package config

import (
	"gopkg.in/urfave/cli.v2"
	"github.com/enaml-ops/pluginlib/util"
)

func RequiredCertFlags() []string {
	return []string{"consul-server-ca-cert", "consul-agent-cert", "consul-agent-key", "consul-server-cert", "consul-server-key",
		"bbs-server-ca-cert", "bbs-client-cert", "bbs-client-key", "bbs-server-cert", "bbs-server-key",
		"etcd-server-cert", "etcd-server-key", "etcd-client-cert", "etcd-client-key", "etcd-peer-cert", "etcd-peer-key",
		"uaa-saml-service-provider-key", "uaa-saml-service-provider-cert", "uaa-jwt-signing-key", "uaa-jwt-verification-key",
		"router-ssl-cert", "router-ssl-key",
	}
}

func NewCerts(c *cli.Context) (*Certs, error) {
	certs := &Certs{}

	caCert, err := pluginutil.LoadResourceFromContext(c, "consul-server-ca-cert")
	if err != nil {
		return nil, err
	}
	agentCert, err := pluginutil.LoadResourceFromContext(c, "consul-agent-cert")
	if err != nil {
		return nil, err
	}
	agentKey, err := pluginutil.LoadResourceFromContext(c, "consul-agent-key")
	if err != nil {
		return nil, err
	}
	serverCert, err := pluginutil.LoadResourceFromContext(c, "consul-server-cert")
	if err != nil {
		return nil, err
	}
	serverKey, err := pluginutil.LoadResourceFromContext(c, "consul-server-key")
	if err != nil {
		return nil, err
	}
	bbsCaCert, err := pluginutil.LoadResourceFromContext(c, "bbs-server-ca-cert")
	if err != nil {
		return nil, err
	}

	bbsClientCert, err := pluginutil.LoadResourceFromContext(c, "bbs-client-cert")
	if err != nil {
		return nil, err
	}

	bbsClientKey, err := pluginutil.LoadResourceFromContext(c, "bbs-client-key")
	if err != nil {
		return nil, err
	}

	bbsServerCert, err := pluginutil.LoadResourceFromContext(c, "bbs-server-cert")
	if err != nil {
		return nil, err
	}

	bbsServerKey, err := pluginutil.LoadResourceFromContext(c, "bbs-server-key")
	if err != nil {
		return nil, err
	}

	etcdServerCert, err := pluginutil.LoadResourceFromContext(c, "etcd-server-cert")
	if err != nil {
		return nil, err
	}

	etcdServerKey, err := pluginutil.LoadResourceFromContext(c, "etcd-server-key")
	if err != nil {
		return nil, err
	}

	etcdClientCert, err := pluginutil.LoadResourceFromContext(c, "etcd-client-cert")
	if err != nil {
		return nil, err
	}

	etcdClientKey, err := pluginutil.LoadResourceFromContext(c, "etcd-client-key")
	if err != nil {
		return nil, err
	}

	etcdPeerCert, err := pluginutil.LoadResourceFromContext(c, "etcd-peer-cert")
	if err != nil {
		return nil, err
	}

	etcdPeerKey, err := pluginutil.LoadResourceFromContext(c, "etcd-peer-key")
	if err != nil {
		return nil, err
	}

	sslpem, err := pluginutil.LoadResourceFromContext(c, "haproxy-sslpem")
	if err != nil {
		return nil, err
	}

	routerCert, err := pluginutil.LoadResourceFromContext(c, "router-ssl-cert")
	if err != nil {
		return nil, err
	}
	routerKey, err := pluginutil.LoadResourceFromContext(c, "router-ssl-key")
	if err != nil {
		return nil, err
	}

	samlKey, err := pluginutil.LoadResourceFromContext(c, "uaa-saml-service-provider-key")
	if err != nil {
		return nil, err
	}

	samlCert, err := pluginutil.LoadResourceFromContext(c, "uaa-saml-service-provider-cert")
	if err != nil {
		return nil, err
	}

	jwtSigningKey, err := pluginutil.LoadResourceFromContext(c, "uaa-jwt-signing-key")
	if err != nil {
		return nil, err
	}

	jwtVerificationKey, err := pluginutil.LoadResourceFromContext(c, "uaa-jwt-verification-key")
	if err != nil {
		return nil, err
	}

	certs.SAMLServiceProviderCertificate = samlCert
	certs.JWTSigningKey = jwtSigningKey
	certs.JWTVerificationKey = jwtVerificationKey
	certs.SAMLServiceProviderKey = samlKey
	certs.RouterSSLCert = routerCert
	certs.RouterSSLKey = routerKey
	certs.ConsulCaCert = caCert
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
	ConsulCaCert                   string
	ConsulAgentCert                string
	ConsulAgentKey                 string
	ConsulServerCert               string
	ConsulServerKey                string
	JWTVerificationKey             string
	SAMLServiceProviderKey         string
	SAMLServiceProviderCertificate string
	JWTSigningKey                  string
}
