package gemfire_plugin

import "fmt"

const (
	arpCleanerJobName                  = "arp-cleaner"
	serverJobName                      = "server"
	serverGroup                        = "server-group"
	locatorJobName                     = "locator"
	locatorGroup                       = "locator-group"
	releaseName                        = "GemFire"
	releaseVersion                     = "latest"
	defaultServerPort                  = "55001"
	defaultLocatorRestPort             = "8080"
	defaultLocatorPort                 = "55221"
	defaultLocatorVMMemory             = "1024"
	defaultDeploymentName              = "p-gemfire"
	defaultStemcellName                = "ubuntu-trusty"
	defaultStemcellAlias               = "trusty"
	defaultStemcellVersion             = "3232.17"
	defaultServerInstanceCount         = "2"
	defaultDevRestPort                 = "7070"
	defaultDevRestActive               = "true"
	SecurityClientAuthenticatorDefault = "DummyAuthenticator.create"
	SecurityClientAccessorDefault      = "templates.security.SimpleAuthorization.create"
	KeystoreRemotePathDefault          = "/usr/local/share/ca-certificates/gemfire.cer"
)

var (
	ActiveAuthNErr = fmt.Errorf("When you activate authn you must set (public-key-pass, keystore-local-path, security-jar-local-path )")
)
