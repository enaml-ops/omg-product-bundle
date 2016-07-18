package director 
/*
* File Generated by enaml generator
* !!! Please do not edit this file !!!
*/
type Nginx struct {

	/*SslProtocols - Descr: SSL/TLS protocols to allow Default: TLSv1 TLSv1.1 TLSv1.2
*/
	SslProtocols interface{} `yaml:"ssl_protocols,omitempty"`

	/*SslCiphers - Descr: List of SSL ciphers to allow (format: https://www.openssl.org/docs/apps/ciphers.html#CIPHER_LIST_FORMAT) Default: ECDHE-RSA-AES128-GCM-SHA256:ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES256-GCM-SHA384:ECDHE-ECDSA-AES256-GCM-SHA384:DHE-RSA-AES128-GCM-SHA256:DHE-DSS-AES128-GCM-SHA256:kEDH+AESGCM:ECDHE-RSA-AES128-SHA256:ECDHE-ECDSA-AES128-SHA256:ECDHE-RSA-AES128-SHA:ECDHE-ECDSA-AES128-SHA:ECDHE-RSA-AES256-SHA384:ECDHE-ECDSA-AES256-SHA384:ECDHE-RSA-AES256-SHA:ECDHE-ECDSA-AES256-SHA:DHE-RSA-AES128-SHA256:DHE-RSA-AES128-SHA:DHE-DSS-AES128-SHA256:DHE-RSA-AES256-SHA256:DHE-DSS-AES256-SHA:DHE-RSA-AES256-SHA:!aNULL:!eNULL:!EXPORT:!DES:!RC4:!3DES:!MD5:!PSK
*/
	SslCiphers interface{} `yaml:"ssl_ciphers,omitempty"`

	/*SslPreferServerCiphers - Descr: Prefer server's cipher priority instead of client's (true for On, false for Off) Default: true
*/
	SslPreferServerCiphers interface{} `yaml:"ssl_prefer_server_ciphers,omitempty"`

	/*Workers - Descr: Number of nginx workers for director Default: 2
*/
	Workers interface{} `yaml:"workers,omitempty"`

}