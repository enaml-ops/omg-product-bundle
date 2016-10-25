package cloudfoundry

import (
	"strings"

	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/plugin/config"
	"github.com/enaml-ops/pluginlib/pcli"
	"github.com/enaml-ops/pluginlib/pluginutil"
	"github.com/enaml-ops/pluginlib/product"
	"github.com/xchapter7x/lo"
	"gopkg.in/urfave/cli.v2"
)

func init() {
	RegisterInstanceGrouperFactory(NewConsulPartition)
	RegisterInstanceGrouperFactory(NewNatsPartition)
	RegisterInstanceGrouperFactory(NewEtcdPartition)
	RegisterInstanceGrouperFactory(NewDiegoDatabasePartition)
	RegisterInstanceGrouperFactory(NewNFSPartition)
	RegisterInstanceGrouperFactory(NewGoRouterPartition)
	RegisterInstanceGrouperFactory(NewMySQLProxyPartition)
	RegisterInstanceGrouperFactory(NewMySQLPartition)
	RegisterInstanceGrouperFactory(NewCloudControllerPartition)
	RegisterInstanceGrouperFactory(NewHaProxyPartition)
	RegisterInstanceGrouperFactory(NewClockGlobalPartition)
	RegisterInstanceGrouperFactory(NewCloudControllerWorkerPartition)
	RegisterInstanceGrouperFactory(NewUAAPartition)
	RegisterInstanceGrouperFactory(NewDiegoBrainPartition)
	RegisterInstanceGrouperFactory(NewDiegoCellPartition)
	RegisterInstanceGrouperFactory(NewDopplerPartition)
	RegisterInstanceGrouperFactory(NewLoggregatorTrafficController)

	//errands
	RegisterInstanceGrouperFactory(NewSmokeErrand)
	RegisterInstanceGrouperFactory(NewBootstrapPartition)
	acceptanceTests := func(config *config.Config) InstanceGroupCreator {
		return NewAcceptanceTestsPartition(true, config)
	}
	internetLessAcceptanceTests := func(config *config.Config) InstanceGroupCreator {
		return NewAcceptanceTestsPartition(false, config)
	}
	RegisterInstanceGrouperFactory(acceptanceTests)
	RegisterInstanceGrouperFactory(internetLessAcceptanceTests)
	RegisterInstanceGrouperFactory(NewPushAppsManager)
	RegisterInstanceGrouperFactory(NewDeployAutoscaling)
	RegisterInstanceGrouperFactory(NewAutoscaleRegisterBroker)
	RegisterInstanceGrouperFactory(NewAutoscaleDestroyBroker)
	RegisterInstanceGrouperFactory(NewAutoscalingTests)
	RegisterInstanceGrouperFactory(NewNotifications)
	RegisterInstanceGrouperFactory(NewNotificationsTest)
	RegisterInstanceGrouperFactory(NewNotificationsUI)
	RegisterInstanceGrouperFactory(NewNotificationsUITest)
}

//GetFlags -
func (s *Plugin) GetFlags() (flags []pcli.Flag) {
	return []pcli.Flag{
		pcli.CreateStringFlag("cf-release-version", "version for cf bosh release", CFReleaseVersion),
		pcli.CreateStringFlag("garden-release-version", "version for garden bosh release", GardenReleaseVersion),
		pcli.CreateStringFlag("diego-release-version", "version for diego bosh release", DiegoReleaseVersion),
		pcli.CreateStringFlag("etcd-release-version", "version for etcd bosh release", EtcdReleaseVersion),
		pcli.CreateStringFlag("cf-mysql-release-version", "version for cf-mysql bosh release", CFMysqlReleaseVersion),
		pcli.CreateStringFlag("cflinuxfs2-release-version", "version for cflinuxfs2 bosh release", CFLinuxReleaseVersion),
		pcli.CreateStringFlag("stemcell-version", "version of stemcell", s.GetMeta().Stemcell.Version),

		// shared for all instance groups:
		pcli.CreateBoolFlag("infer-from-cloud", "setting this flag will attempt to pull as many defaults from your targetted bosh's cloud config as it can (vmtype, network, disk, etc)."),
		pcli.CreateStringFlag("stemcell-name", "the alias of your desired stemcell", s.GetMeta().Stemcell.Alias),
		pcli.CreateStringSliceFlag("az", "list of AZ names to use"),
		pcli.CreateStringFlag("network", "the name of the network to use"),
		pcli.CreateStringFlag("system-domain", "System Domain"),
		pcli.CreateStringSliceFlag("app-domain", "Applications Domains"),
		pcli.CreateBoolFlag("allow-app-ssh-access", "Allow SSH access for CF applications"),

		pcli.CreateStringSliceFlag("router-ip", "a list of the router ips you wish to use"),
		pcli.CreateStringFlag("router-vm-type", "the name of your desired vm size"),
		pcli.CreateStringFlag("router-ssl-cert", "the go router ssl cert, or a filename preceded by '@'"),
		pcli.CreateStringFlag("router-ssl-key", "the go router ssl key, or a filename preceded by '@'"),
		pcli.CreateStringFlag("router-user", "the username of the go-routers", "router_status"),
		pcli.CreateStringFlag("router-pass", "the password of the go-routers"),
		pcli.CreateBoolFlag("router-enable-ssl", "enable or disable ssl on your routers"),

		pcli.CreateBoolTFlag("skip-haproxy", "this flag is on by default and it will skip installing haproxy"),
		pcli.CreateStringSliceFlag("haproxy-ip", "a list of the haproxy ips you wish to use"),
		pcli.CreateStringFlag("haproxy-vm-type", "the name of your desired vm size"),

		pcli.CreateStringFlag("nats-vm-type", "the name of your desired vm size for NATS"),
		pcli.CreateStringFlag("nats-user", "username for your nats pool", "nats"),
		pcli.CreateStringFlag("nats-pass", "password for your nats pool", "nats-password"),
		pcli.CreateIntFlag("nats-port", "the port for the NATS server to listen on", "4222"),
		pcli.CreateStringSliceFlag("nats-machine-ip", "ip of a nats node vm"),

		pcli.CreateStringSliceFlag("consul-ip", "a list of the consul ips you wish to use"),
		pcli.CreateStringFlag("consul-vm-type", "the name of your desired vm size for consul"),
		pcli.CreateStringSliceFlag("consul-encryption-key", "encryption key for consul"),
		pcli.CreateStringFlag("consul-agent-cert", "agent cert for consul, or a filename preceded by '@'"),
		pcli.CreateStringFlag("consul-agent-key", "agent key for consul, or a filename preceded by '@'"),
		pcli.CreateStringFlag("consul-server-cert", "server cert for consul, or a filename preceded by '@'"),
		pcli.CreateStringFlag("consul-server-key", "server key for consul, or a filename preceded by '@'"),

		pcli.CreateStringFlag("syslog-address", "address of syslog server"),
		pcli.CreateIntFlag("syslog-port", "port of syslog server", "514"),
		pcli.CreateStringFlag("syslog-transport", "transport to syslog server", "tcp"),

		pcli.CreateStringSliceFlag("etcd-machine-ip", "ip of a etcd node vm"),
		pcli.CreateStringFlag("etcd-vm-type", "the name of your desired vm size for etcd"),
		pcli.CreateStringFlag("etcd-disk-type", "the name of your desired persistent disk type for etcd"),

		pcli.CreateStringFlag("nfs-ip", "a list of the nfs ips you wish to use"),
		pcli.CreateStringFlag("nfs-vm-type", "the name of your desired vm size for nfs"),
		pcli.CreateStringFlag("nfs-disk-type", "the name of your desired persistent disk type for nfs"),
		pcli.CreateStringFlag("nfs-share-path", "NFS Share Path", "/var/vcap/store"),
		pcli.CreateStringSliceFlag("nfs-allow-from-network-cidr", "the network cidr you wish to allow connections to nfs from"),

		//Mysql Flags
		pcli.CreateStringSliceFlag("mysql-ip", "a list of the mysql ips you wish to use"),
		pcli.CreateStringFlag("mysql-vm-type", "the name of your desired vm size for mysql"),
		pcli.CreateStringFlag("mysql-disk-type", "the name of your desired persistent disk type for mysql"),
		pcli.CreateStringFlag("mysql-admin-password", "admin password for mysql"),
		pcli.CreateStringFlag("mysql-bootstrap-username", "bootstrap username for mysql", "enamlmbu"),
		pcli.CreateStringFlag("mysql-bootstrap-password", "bootstrap password for mysql"),

		//MySQL proxy flags
		pcli.CreateStringSliceFlag("mysql-proxy-ip", "a list of -mysql proxy ips you wish to use"),
		pcli.CreateStringFlag("mysql-proxy-vm-type", "the name of your desired vm size for mysql proxy"),
		pcli.CreateStringFlag("mysql-proxy-external-host", "Host name of MySQL proxy"),
		pcli.CreateStringFlag("mysql-proxy-api-username", "Proxy API user name", "enamlmpa"),
		pcli.CreateStringFlag("mysql-proxy-api-password", "Proxy API password"),

		//CC Worker Partition Flags
		pcli.CreateIntFlag("cc-worker-instances", "the number of vms for cc workers", "1"),
		pcli.CreateStringFlag("cc-worker-vm-type", "the name of the desired vm type for cc worker"),
		pcli.CreateStringFlag("cc-staging-upload-user", "user name for staging upload", "staging_upload_user"),
		pcli.CreateStringFlag("cc-staging-upload-password", "password for staging upload"),
		pcli.CreateStringFlag("cc-bulk-api-password", "password for bulk api calls"),
		pcli.CreateStringFlag("cc-db-encryption-key", "Cloud Controller DB encryption key"),
		pcli.CreateStringFlag("cc-internal-api-user", "user name for Internal API calls"),
		pcli.CreateStringFlag("cc-internal-api-password", "password for Internal API calls"),
		pcli.CreateIntFlag("cc-uploader-poll-interval", "CC uploader job polling interval, in seconds", "25"),
		pcli.CreateStringFlag("cc-vm-type", "Cloud Controller VM Type"),
		pcli.CreateIntFlag("cc-instances", "the number of vms for cc", "1"),
		pcli.CreateStringFlag("host-key-fingerprint", "Host Key Fingerprint"),
		pcli.CreateStringFlag("support-address", "Support URL", "https://support.pivotal.io"),
		pcli.CreateStringFlag("min-cli-version", "Min CF CLI Version supported", "6.7.0"),

		pcli.CreateStringFlag("db-uaa-username", "uaa db username", "enamluaa"),
		pcli.CreateStringFlag("db-uaa-password", "uaa db password"),
		pcli.CreateStringFlag("db-ccdb-username", "ccdb db username", "enamlccdb"),
		pcli.CreateStringFlag("db-ccdb-password", "ccdb db password"),
		pcli.CreateStringFlag("db-console-username", "console db username", "enamlconsole"),
		pcli.CreateStringFlag("db-console-password", "console db password"),
		pcli.CreateStringFlag("db-app_usage-password", "app usage db password"),
		pcli.CreateStringFlag("db-autoscale-username", "autoscale db user", "enamlautoscale"),
		pcli.CreateStringFlag("db-autoscale-password", "autoscale db password"),
		pcli.CreateStringFlag("db-notifications-username", "notifications db user", "enamlnotifications"),
		pcli.CreateStringFlag("db-notifications-password", "notifications db password"),

		//Diego Database
		pcli.CreateStringSliceFlag("diego-db-ip", "a list of static IPs for the diego database partitions"),
		pcli.CreateStringFlag("diego-db-vm-type", "the name of the desired vm type for the diego db"),
		pcli.CreateStringFlag("diego-db-disk-type", "the name of your desired persistent disk type for the diego db"),
		pcli.CreateStringFlag("diego-db-passphrase", "the passphrase for your database"),
		pcli.CreateStringFlag("bbs-server-cert", "BBS server SSL cert (or a file containing it: file format `@filepath`)"),
		pcli.CreateStringFlag("bbs-server-key", "BBS server SSL key (or a file containing it: file format `@filepath`)"),
		pcli.CreateStringFlag("etcd-server-key", "etcd server SSL key (or a file containing it: file format `@filepath`)"),
		pcli.CreateStringFlag("etcd-server-cert", "etcd server cert  (or a file containing it: file format `@filepath`)"),
		pcli.CreateStringFlag("etcd-client-key", "etcd client SSL key (or a file containing it: file format `@filepath`)"),
		pcli.CreateStringFlag("etcd-client-cert", "etcd client SSL cert (or a file containing it: file format `@filepath`)"),
		pcli.CreateStringFlag("etcd-peer-key", "etcd peer SSL key (or a file containing it: file format `@filepath`)"),
		pcli.CreateStringFlag("etcd-peer-cert", "etcd peer SSL cert (or a file containing it: file format `@filepath`)"),

		// Diego Cell
		pcli.CreateStringSliceFlag("diego-cell-ip", "a list of static IPs for the diego cell"),
		pcli.CreateStringFlag("diego-cell-vm-type", "the name of the desired vm type for the diego cell"),
		pcli.CreateStringFlag("diego-cell-disk-type", "the name of your desired persistent disk type for the diego cell"),

		// Diego Brain
		pcli.CreateStringSliceFlag("diego-brain-ip", "a list of static IPs for the diego brain"),
		pcli.CreateStringFlag("diego-brain-vm-type", "the name of the desired vm type for the diego brain"),
		pcli.CreateStringFlag("diego-brain-disk-type", "the name of your desired persistent disk type for the diego brain"),

		pcli.CreateStringFlag("bbs-server-ca-cert", "BBS CA SSL cert (or a file containing it: file format `@filepath`)"),
		pcli.CreateStringFlag("bbs-client-cert", "BBS client SSL cert (or a file containing it: file format `@filepath`)"),
		pcli.CreateStringFlag("bbs-client-key", "BBS client SSL key (or a file containing it: file format `@filepath`)"),
		pcli.CreateBoolTFlag("bbs-require-ssl", "enable SSL for all communications with the BBS"),

		pcli.CreateBoolTFlag("skip-cert-verify", "ignore bad SSL certificates when connecting over HTTPS"),

		pcli.CreateIntFlag("cc-external-port", "external port of the Cloud Controller API", "9022"),
		pcli.CreateStringFlag("ssh-proxy-uaa-secret", "the OAuth client secret used to authenticate the SSH proxy"),
		pcli.CreateIntFlag("loggregator-port", "port for loggregator", "443"),
		pcli.CreateStringFlag("clock-global-vm-type", "the name of the desired vm type for the clock global partition"),

		//Doppler
		pcli.CreateStringSliceFlag("doppler-ip", "a list of the doppler ips you wish to use"),
		pcli.CreateStringFlag("doppler-vm-type", "the name of your desired vm size for doppler"),
		pcli.CreateStringFlag("doppler-zone", "the name zone for doppler"),
		pcli.CreateIntFlag("doppler-drain-buffer-size", "message drain buffer size", "10000"),
		pcli.CreateStringFlag("doppler-shared-secret", "doppler shared secret"),

		//Loggregator Traffic Controller
		pcli.CreateStringSliceFlag("loggregator-traffic-controller-ip", "a list of loggregator traffic controller IPs"),
		pcli.CreateStringFlag("loggregator-traffic-controller-vmtype", "the name of your desired vm size for the loggregator traffic controller"),

		//UAA
		pcli.CreateStringFlag("uaa-vm-type", "the name of your desired vm size for uaa"),
		pcli.CreateIntFlag("uaa-instances", "the number of your desired vms for uaa", "1"),

		pcli.CreateStringFlag("uaa-company-name", "name of company for UAA branding", "Pivotal"),
		pcli.CreateStringFlag("uaa-product-logo", "product logo for UAA branding", "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"),
		pcli.CreateStringFlag("uaa-square-logo", "square logo for UAA branding", "iVBORw0KGgoAAAANSUhEUgAAAGwAAABsCAYAAACPZlfNAAAAAXNSR0IArs4c6QAABYtJREFUeAHtnVtsFFUYx7/d3ruWotUKVIkNaCw02YgJGBRTMd4CokUejD4QH4gxQcIDeHnBmPjkhSghUYLGe3ywPtAHNCo0QgkWwi2tXG2V1kIpLXTbLt1tS9dzlmzSJssZhv32zDk7/2km2znn7Pd9+/vt2Z2dmW0D9Obat4gCiwiLBQQSLflSViAQeN6Can1fYiJBFPQ9BcsAQBiEWUbAsnIxwyDMMgKWlYsZBmGWEbCsXMwwCLOMgGXlYoZBmGUELCsXMwzCLCNgWbmYYRBmGQHLysUMgzDLCFhWLmYYhFlGwLJyMcMgzDIClpWLGQZhlhGwrFzMMAizjIBl5WKGQZhlBCwrV1xbb96y59V1VFJQmLawQNrWa43x8XEaHo1fW+Oj1H8lSqf6eulEbw+dvNhLvcNDinvb0WWksAdm3UWhwiJ2gt2RAWo80UY7jrdSU8cZGrt6lT1HtgMaKSxbD7qqfDq99tAjyTUSG6FP9v1BH+3dTUPxeLZSssf17U5HeXEJbXr8aerY+A6tf7iOxFeu2OFmI6BvhaVgVoRCtHl5PTW8/AoV5xekmo299b2wlJn6+WFqWrOWKkpDqSYjbyFskpZFs++hL1e9NKnFvF+t3OmQOwzdkcgUmnnBABXm5Ys1j8qKisVadFPvS8tramn1goX09eEDU+KbsmGlsMbjbbT6x++UDOVORGXoFppXOYMerLqbVsyrpcWzqykYdH5R+fjZlcnd/8sjV5Q5vOh0rt6LqhhyJsQ3uC+ID8ry89aHYtf90W1bKLzlffr19EnH6HIP8oXasOM4LwbkrLB0MP+6cJ6e+eoz+vTP5nTdU9peDC+Ysm3Khq+ESehy5r3e2ECHu7uUDuqq59Id4iXVtMV3wqSACSHt3V2/KF3I97qayjuVY7zo9KUwCfq3M6coNjamZD6zrFzZ70Wnb4XFxseoK3JZyXzWtGnKfi86fStMwu6LRpXMZ5RBmBKQ7k75XqZa8gLmPZ/Nq0hFkLnvttJSZUT5Oc60xbfC5CGs6lsrlD56hgaV/V50+lbYkuo5VFygPp3SMwxhXjwp0+bcsGRp2vZU48TEBB09153aNObWlzNMHo1/6r4apYTmsx10MTqsHONFp5VH6zMBtWbhYtq6YpVjiJ/ajjmO8WKAL4QFxamWZffPT1678dicex05D4jTKj8cO+Q4zosBOSXs7bonktci5ovjgPIUye3ieo3wzKrk+TC5faPLGz83On6ovtFY3ONySth7Ty67qbPMk6Hu+edv+vzg/slNRv3uy52O6xk40HWW6r/94nrdRrTn1AzLhOju9tP03DfbKTo6mkmYrN/X98L6xQHgTb/vpG0t+5LnybJOPMMEvhXWOXCJvj9yiD7Yu4sGRkYyxKjv7r4RJi+Na+05Rwf/66SG1qO0v/NffZQZM+WUsI07d1BC/MTE144GYzHxJYcYDYq1vb/f8WQlI9OshsopYZubm7IKy4Tg2K03wYKLGiDMBSwThkKYCRZc1ABhLmCZMBTCTLDgogYIcwHLhKEQZoIFFzVAmAtYJgyFMBMsuKgBwlzAMmEohJlgwUUNEOYClglDIcwECy5qgDAXsEwYCmEmWHBRA4S5gGXCUAgzwYKLGow84yyvuyhR/GW19kt9Lh5ibg01UtjS7VtzizLjo8FLIiNMHaEgTAdlxhwQxghTRygI00GZMQeEMcLUEQrCdFBmzAFhjDB1hIIwHZQZc0AYI0wdoSBMB2XGHBDGCFNHKAjTQZkxB4QxwtQRCsJ0UGbMAWGMMHWEgjAdlBlzQBgjTB2hIEwHZcYcEMYIU0coCNNBmTEHhDHC1BEKwnRQZswBYYwwdYSCMB2UGXNAGCNMHaEgTAdlxhziUu1Ei8M/+WFMh1CZEUi0/A+j7hNSB5Wo2wAAAABJRU5ErkJggg=="),
		pcli.CreateStringFlag("uaa-footer-legal-txt", "legal text for UAA branding", "Legal Text"),
		pcli.CreateBoolTFlag("uaa-enable-selfservice-links", "enable self service links"),
		pcli.CreateBoolTFlag("uaa-signups-enabled", "enable signups"),
		pcli.CreateStringFlag("uaa-login-protocol", "uaa login protocol, default https", "https"),
		pcli.CreateStringFlag("uaa-saml-service-provider-key", "saml service provider key for uaa"),
		pcli.CreateStringFlag("uaa-saml-service-provider-cert", "saml service provider certificate for uaa"),
		pcli.CreateStringFlag("uaa-jwt-signing-key", "signing key for jwt used by UAA"),
		pcli.CreateStringFlag("uaa-jwt-verification-key", "verification key for jwt used by UAA"),
		pcli.CreateBoolFlag("uaa-ldap-enabled", "is ldap enabled for UAA"),
		pcli.CreateStringFlag("uaa-ldap-url", "url for ldap server"),
		pcli.CreateStringFlag("uaa-ldap-user-dn", "userDN to bind to ldap with"),
		pcli.CreateStringFlag("uaa-ldap-user-password", "bind password for ldap user"),
		pcli.CreateStringFlag("uaa-ldap-search-filter", "search filter for users"),
		pcli.CreateStringFlag("uaa-ldap-search-base", "search base for users"),
		pcli.CreateStringFlag("uaa-ldap-mail-attributename", "attribute name for mail"),
		pcli.CreateStringFlag("uaa-admin-secret", "admin account client secret"),

		//User accounts
		pcli.CreateStringFlag("admin-password", "password for admin account"),
		pcli.CreateStringFlag("push-apps-manager-password", "password for push_apps_manager account"),
		pcli.CreateStringFlag("smoke-tests-password", "password for smoke_tests account"),
		pcli.CreateStringFlag("system-services-password", "password for system_services account"),
		pcli.CreateStringFlag("system-verification-password", "password for system_verification account"),

		//Client secrets
		pcli.CreateStringFlag("opentsdb-firehose-nozzle-client-secret", "client-secret for opentsdb firehose nozzle"),
		pcli.CreateStringFlag("identity-client-secret", "client-secret for identity"),
		pcli.CreateStringFlag("login-client-secret", "client-secret for login"),
		pcli.CreateStringFlag("portal-client-secret", "client-secret for portal"),
		pcli.CreateStringFlag("autoscaling-service-client-secret", "client-secret for autoscaling service"),
		pcli.CreateStringFlag("system-passwords-client-secret", "client-secret for system-passwords"),
		pcli.CreateStringFlag("cc-service-dashboards-client-secret", "client-secret for cc-service-dashboards"),
		pcli.CreateStringFlag("doppler-client-secret", "client-secret for doppler"),
		pcli.CreateStringFlag("gorouter-client-secret", "client-secret for gorouter"),
		pcli.CreateStringFlag("notifications-client-secret", "client-secret for notifications"),
		pcli.CreateStringFlag("notifications-ui-client-secret", "client-secret for notification-ui"),
		pcli.CreateStringFlag("cloud-controller-username-lookup-client-secret", "client-secret for cloud controller username lookup"),
		pcli.CreateStringFlag("cc-routing-client-secret", "client-secret for cc routing"),
		pcli.CreateStringFlag("ssh-proxy-client-secret", "client-secret for ssh proxy"),
		pcli.CreateStringFlag("apps-metrics-client-secret", "client-secret for apps metrics "),
		pcli.CreateStringFlag("apps-metrics-processing-client-secret", "client-secret for apps metrics processing"),

		pcli.CreateStringFlag("errand-vm-type", "vm type to be used for running errands"),
		pcli.CreateStringFlag("haproxy-sslpem", "SSL pem for HAProxy"),
		pcli.CreateStringFlag("apps-manager-secret-token", "apps manager secret token for signing cookies"),

		//Vault stuff
		pcli.CreateStringFlag("vault-domain", "the location of your vault server (ie. http://10.0.0.1:8200)"),
		pcli.CreateStringFlag("vault-hash-misc", "the hashname for misc CLI flags"),
		pcli.CreateStringFlag("vault-hash-password", "the hashname of your secret (ie. secret/pcf-1-passwords"),
		pcli.CreateStringFlag("vault-hash-keycert", "the hashname of your secret (ie. secret/pcf-1-keycert"),
		pcli.CreateStringFlag("vault-hash-ip", "the hashname of your secret (ie. secret/pcf-1-ips"),
		pcli.CreateStringFlag("vault-hash-host", "the hashname of your secret (ie. secret/pcf-1-hosts"),
		pcli.CreateStringFlag("vault-token", "the token to make connections to your vault"),
		pcli.CreateBoolFlag("vault-rotate", "set this flag to true if you would like re/set the values in vault. this will rotate internal certs and passwords"),
		pcli.CreateBoolTFlag("vault-active", "use the data which is stored in vault for the flag values it contains"),
	}
}

//GetMeta -
func (s *Plugin) GetMeta() product.Meta {
	return product.Meta{
		Name: "cloudfoundry",

		Stemcell: enaml.Stemcell{
			Name:    StemcellName,
			Alias:   StemcellAlias,
			Version: StemcellVersion,
		},
		Releases: []enaml.Release{
			enaml.Release{
				Name:    CFReleaseName,
				Version: CFReleaseVersion,
			},
			enaml.Release{
				Name:    CFMysqlReleaseName,
				Version: CFMysqlReleaseVersion,
			},
			enaml.Release{
				Name:    DiegoReleaseName,
				Version: DiegoReleaseVersion,
			},
			enaml.Release{
				Name:    GardenReleaseName,
				Version: GardenReleaseVersion,
			},
			enaml.Release{
				Name:    CFLinuxReleaseName,
				Version: CFLinuxReleaseVersion,
			},
			enaml.Release{
				Name:    EtcdReleaseName,
				Version: EtcdReleaseVersion,
			},
			enaml.Release{
				Name:    PushAppsReleaseName,
				Version: PushAppsReleaseVersion,
			},
			enaml.Release{
				Name:    NotificationsReleaseName,
				Version: NotificationsReleaseVersion,
			},
			enaml.Release{
				Name:    NotificationsUIReleaseName,
				Version: NotificationsUIReleaseVersion,
			},
			enaml.Release{
				Name:    CFAutoscalingReleaseName,
				Version: CFAutoscalingReleaseVersion,
			},
		},
		Properties: map[string]interface{}{
			"version":                   s.PluginVersion,
			"pivotal-elastic-runtime":   strings.Join([]string{"pivotal-elastic-runtime", PivotalERTVersion}, " / "),
			"cf-release":                strings.Join([]string{CFReleaseName, CFReleaseVersion}, " / "),
			"cf-mysql-release":          strings.Join([]string{CFMysqlReleaseName, CFMysqlReleaseVersion}, " / "),
			"diego-release":             strings.Join([]string{DiegoReleaseName, DiegoReleaseVersion}, " / "),
			"garden-linux-release":      strings.Join([]string{GardenReleaseName, GardenReleaseVersion}, " / "),
			"cflinuxfs2-rootfs-release": strings.Join([]string{CFLinuxReleaseName, CFLinuxReleaseVersion}, " / "),
			"etcd-release":              strings.Join([]string{EtcdReleaseName, EtcdReleaseVersion}, " / "),
			"pushapp-release":           strings.Join([]string{PushAppsReleaseName, PushAppsReleaseVersion}, " / "),
			"stemcell":                  StemcellVersion,
			"notifications-release":     strings.Join([]string{NotificationsReleaseName, NotificationsReleaseVersion}, " / "),
			"notifications-ui-release":  strings.Join([]string{NotificationsUIReleaseName, NotificationsUIReleaseVersion}, " / "),
			"cf-autoscaling-release":    strings.Join([]string{CFAutoscalingReleaseName, CFAutoscalingReleaseVersion}, " / "),
		},
	}
}

//GetProduct -
func (s *Plugin) GetProduct(args []string, cloudConfig []byte) (b []byte, err error) {
	flgs := s.GetFlags()
	InferFromCloudDecorate(flagsToInferFromCloudConfig, cloudConfig, args, flgs)

	if err := VaultRotate(args, flgs); err != nil {
		lo.G.Errorf("unable to rotate vault values: %v", err.Error())
		return nil, err
	}
	VaultDecorate(args, flgs)
	c := pluginutil.NewContext(args, pluginutil.ToCliFlagArray(flgs))
	var dm *enaml.DeploymentManifest
	var cfg *config.Config

	if cfg, err = config.NewConfig(c); err == nil {
		dm, err = s.getDeploymentManifest(c, cfg)

		if err != nil {
			lo.G.Errorf("error creating manifest: %v", err.Error())
		}
	} else {
		lo.G.Errorf("error getting config: %v", err.Error())
	}
	return dm.Bytes(), err
}

func (s *Plugin) getDeploymentManifest(c *cli.Context, config *config.Config) (*enaml.DeploymentManifest, error) {
	dm := enaml.NewDeploymentManifest([]byte(``))
	dm.SetName(DeploymentName)

	dm.AddRelease(enaml.Release{
		Name:    CFReleaseName,
		Version: c.String("cf-release-version"),
		URL:     c.String("cf-release-url"),
		SHA1:    c.String("cf-release-sha"),
	})
	dm.AddRelease(enaml.Release{
		Name:    CFMysqlReleaseName,
		Version: c.String("cf-mysql-release-version"),
		URL:     c.String("cf-mysql-release-url"),
		SHA1:    c.String("cf-mysql-release-sha"),
	})
	dm.AddRelease(enaml.Release{
		Name:    DiegoReleaseName,
		Version: c.String("diego-release-version"),
		URL:     c.String("diego-release-url"),
		SHA1:    c.String("diego-release-sha"),
	})
	dm.AddRelease(enaml.Release{
		Name:    GardenReleaseName,
		Version: c.String("garden-release-version"),
		URL:     c.String("garden-release-url"),
		SHA1:    c.String("garden-release-sha"),
	})
	dm.AddRelease(enaml.Release{
		Name:    CFLinuxReleaseName,
		Version: c.String("cflinuxfs2-release-version"),
		URL:     c.String("cflinuxfs2-release-url"),
		SHA1:    c.String("cflinuxfs2-release-sha"),
	})
	dm.AddRelease(enaml.Release{
		Name:    EtcdReleaseName,
		Version: c.String("etcd-release-version"),
		URL:     c.String("etcd-release-url"),
		SHA1:    c.String("etcd-release-sha"),
	})

	dm.AddRelease(enaml.Release{Name: MySQLBackupReleaseName, Version: MySQLBackupReleaseVersion})
	dm.AddRelease(enaml.Release{Name: PushAppsReleaseName, Version: PushAppsReleaseVersion})
	dm.AddRelease(enaml.Release{Name: CFAutoscalingReleaseName, Version: CFAutoscalingReleaseVersion})
	dm.AddRelease(enaml.Release{Name: NotificationsReleaseName, Version: NotificationsReleaseVersion})
	dm.AddRelease(enaml.Release{Name: NotificationsUIReleaseName, Version: NotificationsUIReleaseVersion})

	dm.AddStemcell(enaml.Stemcell{OS: StemcellName, Version: c.String("stemcell-version"), Alias: c.String("stemcell-name")})

	dm.Update.MaxInFlight = 1
	dm.Update.Canaries = 1
	dm.Update.Serial = false
	dm.Update.CanaryWatchTime = "30000-300000"
	dm.Update.UpdateWatchTime = "30000-300000"

	for _, factory := range factories {
		grouper := factory(config)
		if ig := grouper.ToInstanceGroup(); ig != nil {
			lo.G.Debug("instance-group: ", ig)
			dm.AddInstanceGroup(ig)
		}
	}
	return dm, nil
}

func InferFromCloudDecorate(inferFlagMap map[string][]string, cloudConfig []byte, args []string, flgs []pcli.Flag) {
	c := pluginutil.NewContext(args, pluginutil.ToCliFlagArray(flgs))

	if c.Bool("infer-from-cloud") {
		ccinf := pluginutil.NewCloudConfigInferFromBytes(cloudConfig)
		setAllInferredFlagDefaults(inferFlagMap["disktype"], ccinf.InferDefaultDiskType(), flgs, c)
		setAllInferredFlagDefaults(inferFlagMap["vmtype"], ccinf.InferDefaultVMType(), flgs, c)
		setAllInferredFlagDefaults(inferFlagMap["az"], ccinf.InferDefaultAZ(), flgs, c)
		setAllInferredFlagDefaults(inferFlagMap["network"], ccinf.InferDefaultNetwork(), flgs, c)
	}
}

func setAllInferredFlagDefaults(matchlist []string, defaultvalue string, flgs []pcli.Flag, c *cli.Context) {
	for _, match := range matchlist {
		// only infer flags that weren't manually set
		if !c.IsSet(match) {
			setFlagDefault(match, defaultvalue, flgs)
		}
	}
}

func setFlagDefault(flagname, defaultvalue string, flgs []pcli.Flag) {
	for idx, flg := range flgs {
		if flg.Name == flagname {
			flgs[idx].Value = defaultvalue
		}
	}
}

func VaultDecorate(args []string, flgs []pcli.Flag) {
	c := pluginutil.NewContext(args, pluginutil.ToCliFlagArray(flgs))

	if hasValidVaultFlags(c) {
		lo.G.Debug("connecting to vault at: ", c.String("vault-domain"))
		vault := pluginutil.NewVaultUnmarshal(c.String("vault-domain"), c.String("vault-token"))
		hashes := []string{
			c.String("vault-hash-misc"),
			c.String("vault-hash-password"),
			c.String("vault-hash-keycert"),
			c.String("vault-hash-ip"),
			c.String("vault-hash-host"),
		}

		for _, hash := range hashes {
			if hash != "" {
				vault.UnmarshalFlags(hash, flgs)
			}
		}

	} else {
		lo.G.Debug("complete vault flagset not found:",
			"active: ", c.Bool("vault-active"),
			"domain: ", c.String("vault-domain"),
			"passhash: ", c.String("vault-hash-password"),
			"keycerthash: ", c.String("vault-hash-keycert"),
			"iphash: ", c.String("vault-hash-ip"),
			"hosthash: ", c.String("vault-hash-host"),
			"vaulttoken: ", c.String("vault-token"),
		)

		if c.Bool("vault-active") {
			lo.G.Fatal("you've activated vault, but have not provided a complete set of values... exiting program now")
		}
	}
}

func hasValidVaultFlags(c *cli.Context) bool {
	return c.Bool("vault-active") &&
		c.String("vault-domain") != "" &&
		c.String("vault-token") != ""
}

//GetContext -
func (s *Plugin) GetContext(args []string) (c *cli.Context) {
	c = pluginutil.NewContext(args, pluginutil.ToCliFlagArray(s.GetFlags()))
	return
}
