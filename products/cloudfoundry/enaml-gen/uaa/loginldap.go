package uaa 
/*
* File Generated by enaml generator
* !!! Please do not edit this file !!!
*/
type LoginLdap struct {

	/*SearchFilter - Descr: Deprecated. Use uaa.ldap.searchFilter - login.ldap prefix is used for backwards compatibility to enable ldap from login config Default: cn={0}
*/
	SearchFilter interface{} `yaml:"searchFilter,omitempty"`

	/*PasswordEncoder - Descr: Deprecated. Use uaa.ldap.passwordEncoder - login.ldap prefix is used for backwards compatibility to enable ldap from login config Default: org.cloudfoundry.identity.uaa.login.ldap.DynamicPasswordComparator
*/
	PasswordEncoder interface{} `yaml:"passwordEncoder,omitempty"`

	/*UserDNPattern - Descr: Deprecated. Use uaa.ldap.userDNPattern - login.ldap prefix is used for backwards compatibility to enable ldap from login config Default: <nil>
*/
	UserDNPattern interface{} `yaml:"userDNPattern,omitempty"`

	/*PasswordAttributeName - Descr: Deprecated. Use uaa.ldap.passwordAttributeName - login.ldap prefix is used for backwards compatibility to enable ldap from login config Default: userPassword
*/
	PasswordAttributeName interface{} `yaml:"passwordAttributeName,omitempty"`

	/*SearchBase - Descr: Deprecated. Use uaa.ldap.searchBase - login.ldap prefix is used for backwards compatibility to enable ldap from login config Default: 
*/
	SearchBase interface{} `yaml:"searchBase,omitempty"`

	/*LocalPasswordCompare - Descr: Deprecated. Use uaa.ldap.localPasswordCompare - login.ldap prefix is used for backwards compatibility to enable ldap from login config Default: true
*/
	LocalPasswordCompare interface{} `yaml:"localPasswordCompare,omitempty"`

	/*UserPassword - Descr: Deprecated. Use uaa.ldap.userPassword - login.ldap prefix is used for backwards compatibility to enable ldap from login config Default: <nil>
*/
	UserPassword interface{} `yaml:"userPassword,omitempty"`

	/*Url - Descr: Deprecated. Use uaa.ldap.url - login.ldap prefix is used for backwards compatibility to enable ldap from login config Default: <nil>
*/
	Url interface{} `yaml:"url,omitempty"`

	/*SslCertificateAlias - Descr: Deprecated. Use uaa.ldap.sslCertificateAlias - login.ldap prefix is used for backwards compatibility to enable ldap from login config Default: <nil>
*/
	SslCertificateAlias interface{} `yaml:"sslCertificateAlias,omitempty"`

	/*ProfileType - Descr: Deprecated. Use uaa.ldap.profile_type - login.ldap prefix is used for backwards compatibility to enable ldap from login config Default: <nil>
*/
	ProfileType interface{} `yaml:"profile_type,omitempty"`

	/*UserDN - Descr: Deprecated. Use uaa.ldap.userDN - login.ldap prefix is used for backwards compatibility to enable ldap from login config Default: <nil>
*/
	UserDN interface{} `yaml:"userDN,omitempty"`

	/*SslCertificate - Descr: Deprecated. Use uaa.ldap.sslCertificate - login.ldap prefix is used for backwards compatibility to enable ldap from login config Default: <nil>
*/
	SslCertificate interface{} `yaml:"sslCertificate,omitempty"`

}