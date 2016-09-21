package pmysql

import (
	"encoding/json"
	"fmt"

	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-cli/utils"
	"github.com/enaml-ops/pluginlib/pcli"
	"github.com/enaml-ops/pluginlib/product"
	"github.com/enaml-ops/pluginlib/util"
	"github.com/xchapter7x/lo"
)

const (
	CFMysqlReleaseName            = "cf-mysql"
	CFMysqlReleaseVersion         = "24.6"
	MysqlBackupReleaseName        = "mysql-backup"
	MysqlBackupReleaseVersion     = "1.25.0"
	ServiceBackupReleaseName      = "service-backup"
	ServiceBackupReleaseVersion   = "14"
	MysqlMonitoringReleaseName    = "mysql-monitoring"
	MysqlMonitoringReleaseVersion = "3"

	defaultStemcellVersion = "3232.17"
)

type Plugin struct {
	PluginVersion               string
	DeploymentName              string
	NetworkName                 string
	IPs                         []string
	ProxyIPs                    []string
	VMTypeName                  string
	DiskTypeName                string
	AZs                         []string
	StemcellName                string
	StemcellURL                 string
	StemcellVersion             string
	StemcellSHA                 string
	SyslogAddress               string
	SyslogPort                  string
	SyslogTransport             string
	AdminPassword               string
	SeededDBPassword            string
	GaleraHealthcheckUsername   string
	GaleraHealthcheckPassword   string
	GaleraHealthcheckDBPassword string
	ClusterHealthPassword       string
	BaseDomain                  string
	NotificationClientSecret    string
	UaaAdminClientSecret        string
	NotificationRecipientEmail  string
	BackupEndpointUser          string
	BackupEndpointPassword      string
	NatsPassword                string
	NatsUser                    string
	NatsPort                    string
	ProxyAPIUser                string
	ProxyAPIPass                string
	MonitoringIPs               []string
	BrokerIPs                   []string
	BrokerQuotaEnforcerPassword string
	BrokerAuthUsername          string
	BrokerAuthPassword          string
	BrokerCookieSecret          string
	ServiceSecret               string
	CFAdminPassword             string
}

func (s *Plugin) GetFlags() (flags []pcli.Flag) {
	return []pcli.Flag{
		pcli.Flag{FlagType: pcli.StringFlag, Name: "deployment-name", Value: "p-mysql", Usage: "the name bosh will use for this deployment"},
		pcli.Flag{FlagType: pcli.StringFlag, Name: "base-domain", Usage: "the base domain you wish to associate your mysql routes with"},
		pcli.Flag{FlagType: pcli.StringFlag, Name: "notification-recipient-email", Usage: "email to send monitoring notifications to"},

		pcli.Flag{FlagType: pcli.StringFlag, Name: "stemcell-url", Usage: "the url of the stemcell you wish to use"},
		pcli.Flag{FlagType: pcli.StringFlag, Name: "stemcell-ver", Usage: "the version number of the stemcell you wish to use", Value: defaultStemcellVersion},
		pcli.Flag{FlagType: pcli.StringFlag, Name: "stemcell-sha", Usage: "the sha of the stemcell you will use"},
		pcli.Flag{FlagType: pcli.StringFlag, Name: "stemcell-name", Value: "ubuntu-trusty", Usage: "the OS of the stemcell you will use"},

		pcli.Flag{FlagType: pcli.StringFlag, Name: "cf-mysql-release-version", Value: CFMysqlReleaseVersion, Usage: "the cf-mysql release version to use"},
		pcli.Flag{FlagType: pcli.StringFlag, Name: "mysql-backup-release-version", Value: MysqlBackupReleaseVersion, Usage: "the mysql-backup release version to use"},
		pcli.Flag{FlagType: pcli.StringFlag, Name: "service-backup-release-version", Value: ServiceBackupReleaseVersion, Usage: "the service-backup release version to use"},
		pcli.Flag{FlagType: pcli.StringFlag, Name: "mysql-monitoring-release-version", Value: MysqlMonitoringReleaseVersion, Usage: "the mysql-monitoring release version to user"},

		pcli.Flag{FlagType: pcli.BoolFlag, Name: "infer-from-cloud", Usage: "setting this flag will attempt to pull as many defaults from your targetted bosh's cloud config as it can (vmtype, network, disk, etc)."},
		pcli.Flag{FlagType: pcli.StringSliceFlag, Name: "az", Usage: "list of AZ names to use"},
		pcli.Flag{FlagType: pcli.StringFlag, Name: "network", Usage: "the name of the network to use"},
		pcli.Flag{FlagType: pcli.StringFlag, Name: "vm-type", Usage: "name of your desired vm type"},
		pcli.Flag{FlagType: pcli.StringFlag, Name: "disk-type", Usage: "name of your desired disk type"},

		pcli.Flag{FlagType: pcli.StringSliceFlag, Name: "ip", Usage: "multiple static ips for each mysql vm"},
		pcli.Flag{FlagType: pcli.StringSliceFlag, Name: "proxy-ip", Usage: "multiple static ips for each mysql-proxy vm"},
		pcli.Flag{FlagType: pcli.StringSliceFlag, Name: "monitoring-ip", Usage: "multiple static ips for each monitoring job vm"},
		pcli.Flag{FlagType: pcli.StringSliceFlag, Name: "broker-ip", Usage: "multiple static ips for each broker job vm"},

		// this set of values comes from the ERT deployment
		pcli.Flag{FlagType: pcli.StringFlag, Name: "admin-password", Usage: "the CF admin user's password"},
		pcli.Flag{FlagType: pcli.StringFlag, Name: "notifications-client-secret", Usage: "client secret for monitoring notifications"},
		pcli.Flag{FlagType: pcli.StringFlag, Name: "uaa-admin-secret", Usage: "uaa client secret for monitoring notifications"},
		pcli.Flag{FlagType: pcli.StringFlag, Name: "nats-user", Value: natsUser, Usage: "the user to access the nats instance"},
		pcli.Flag{FlagType: pcli.StringFlag, Name: "nats-pass", Usage: "the password to access the nats instance"},
		pcli.Flag{FlagType: pcli.StringFlag, Name: "nats-port", Value: natsPort, Usage: "the port to access the nats instance"},
		pcli.Flag{FlagType: pcli.StringFlag, Name: "syslog-address", Usage: "the address of your syslog drain"},
		pcli.Flag{FlagType: pcli.StringFlag, Name: "syslog-port", Value: "514", Usage: "the port for your syslog connection"},
		pcli.Flag{FlagType: pcli.StringFlag, Name: "syslog-transport", Value: "tcp", Usage: "the proto for your syslog connection"},

		pcli.Flag{FlagType: pcli.StringFlag, Name: "mysql-admin-password", Usage: "the admin password for your mysql"},
		pcli.Flag{FlagType: pcli.StringFlag, Name: "seeded-db-password", Usage: "canary seeded db password"},
		pcli.Flag{FlagType: pcli.StringFlag, Name: "galera-healthcheck-username", Usage: "galera healthcheck endpoint user"},
		pcli.Flag{FlagType: pcli.StringFlag, Name: "galera-healthcheck-password", Usage: "galera healthcheck endpoint user's password"},
		pcli.Flag{FlagType: pcli.StringFlag, Name: "galera-healthcheck-db-password", Usage: "galera healthcheck db password"},
		pcli.Flag{FlagType: pcli.StringFlag, Name: "cluster-health-password", Usage: "clusterhealth password"},
		pcli.Flag{FlagType: pcli.StringFlag, Name: "backup-endpoint-user", Usage: "the user to access the backup rest endpoint"},
		pcli.Flag{FlagType: pcli.StringFlag, Name: "backup-endpoint-password", Usage: "the password to access the backup rest endpoint"},
		pcli.Flag{FlagType: pcli.StringFlag, Name: "broker-quota-enforcer-password", Usage: "the password to the broker quota enforcer"},
		pcli.Flag{FlagType: pcli.StringFlag, Name: "proxy-api-username", Usage: "the api username for the proxy", Value: "admin"},
		pcli.Flag{FlagType: pcli.StringFlag, Name: "proxy-api-password", Usage: "the api password for the proxy"},
		pcli.Flag{FlagType: pcli.StringFlag, Name: "broker-auth-username", Usage: "a basic auth user for mysql broker"},
		pcli.Flag{FlagType: pcli.StringFlag, Name: "broker-auth-password", Usage: "a basic auth password for mysql broker"},
		pcli.Flag{FlagType: pcli.StringFlag, Name: "broker-cookie-secret", Usage: "the broker cookie secret"},
		pcli.Flag{FlagType: pcli.StringFlag, Name: "service-secret", Usage: "the broker service secret"},

		pcli.Flag{FlagType: pcli.StringFlag, Name: "vault-domain", Usage: "the location of your vault server (ie. http://10.0.0.1:8200)"},
		pcli.Flag{FlagType: pcli.StringFlag, Name: "vault-token", Usage: "the token to make connections to your vault"},
		pcli.Flag{FlagType: pcli.StringSliceFlag, Name: "vault-hash-ert", Usage: "hashes containing ERT secrets.  these hashes are only read, never written (ie. secret/pcf-np-1-passwords"},
		pcli.Flag{FlagType: pcli.StringFlag, Name: "vault-hash-mysql-secret", Usage: "the hash of your secret (ie. secret/p-mysql-1-passwords"},
		pcli.Flag{FlagType: pcli.StringFlag, Name: "vault-hash-mysql-ip", Usage: "the hash of your secret (ie. secret/p-mysql-1-ips"},
		pcli.CreateBoolFlag("vault-rotate", "set this flag to reset the values in vault. this will rotate internal passwords in the 'vault-hash-mysql-secret' hash"),
	}
}

func (s *Plugin) GetMeta() product.Meta {
	return product.Meta{
		Name: "p-mysql",
		Properties: map[string]interface{}{
			"version":        defaultStemcellVersion,
			"stemcell":       s.StemcellVersion,
			"pivotal-mysql":  "1.7.12",
			"cf-mysql":       fmt.Sprintf("%s / %s", CFMysqlReleaseName, CFMysqlReleaseVersion),
			"mysql-backup":   fmt.Sprintf("%s / %s", MysqlBackupReleaseName, MysqlBackupReleaseVersion),
			"service-backup": fmt.Sprintf("%s / %s", ServiceBackupReleaseName, ServiceBackupReleaseVersion),
			"mysql-monitor":  fmt.Sprintf("%s / %s", MysqlMonitoringReleaseName, MysqlMonitoringReleaseVersion),
		},
	}
}

func (s *Plugin) GetProduct(args []string, cloudConfig []byte) (b []byte) {
	flgs := s.GetFlags()
	InferFromCloudDecorate(flagsToInferFromCloudConfig, cloudConfig, args, flgs)
	c := pluginutil.NewContext(args, pluginutil.ToCliFlagArray(flgs))

	domain := c.String("vault-domain")
	token := c.String("vault-token")
	if domain != "" && token != "" {
		v := pluginutil.NewVaultUnmarshal(domain, token)

		for _, hash := range c.StringSlice("vault-hash-ert") {
			ertVaultDecorate(flgs, hash, v)
		}

		hash := c.String("vault-hash-mysql-ip")
		if hash != "" {
			if err := v.UnmarshalFlags(hash, flgs); err != nil {
				lo.G.Error("error unmarshalling vault hash", hash, err)
			}
		}

		hash = c.String("vault-hash-mysql-secret")
		if hash != "" {
			if c.Bool("vault-rotate") {
				if err := vaultRotateMySQL(hash, v); err != nil {
					lo.G.Error("error rotating mysql secrets:", err)
				}
			}

			if err := v.UnmarshalFlags(hash, flgs); err != nil {
				lo.G.Error("error unmarshalling vault hash", hash, err)
			}
		}

		c = pluginutil.NewContext(args, pluginutil.ToCliFlagArray(flgs))
	}

	s.IPs = c.StringSlice("ip")
	s.ProxyIPs = c.StringSlice("proxy-ip")
	s.AZs = c.StringSlice("az")
	s.DeploymentName = c.String("deployment-name")
	s.NetworkName = c.String("network")
	s.StemcellName = c.String("stemcell-name")
	s.StemcellVersion = c.String("stemcell-ver")
	s.StemcellSHA = c.String("stemcell-sha")
	s.StemcellURL = c.String("stemcell-url")
	s.VMTypeName = c.String("vm-type")
	s.DiskTypeName = c.String("disk-type")
	s.AdminPassword = c.String("mysql-admin-password")
	s.SeededDBPassword = c.String("seeded-db-password")
	s.SyslogAddress = c.String("syslog-address")
	s.SyslogPort = c.String("syslog-port")
	s.SyslogTransport = c.String("syslog-transport")
	s.GaleraHealthcheckUsername = c.String("galera-healthcheck-username")
	s.GaleraHealthcheckPassword = c.String("galera-healthcheck-password")
	s.GaleraHealthcheckDBPassword = c.String("galera-healthcheck-db-password")
	s.ClusterHealthPassword = c.String("cluster-health-password")
	s.BaseDomain = c.String("base-domain")
	s.NotificationRecipientEmail = c.String("notification-recipient-email")
	s.NotificationClientSecret = c.String("notifications-client-secret")
	s.UaaAdminClientSecret = c.String("uaa-admin-secret")
	s.BackupEndpointUser = c.String("backup-endpoint-user")
	s.BackupEndpointPassword = c.String("backup-endpoint-password")
	s.NatsUser = c.String("nats-user")
	s.NatsPassword = c.String("nats-pass")
	s.NatsPort = c.String("nats-port")
	s.ProxyAPIUser = c.String("proxy-api-username")
	s.ProxyAPIPass = c.String("proxy-api-password")
	s.MonitoringIPs = c.StringSlice("monitoring-ip")
	s.BrokerIPs = c.StringSlice("broker-ip")
	s.BrokerQuotaEnforcerPassword = c.String("broker-quota-enforcer-password")
	s.BrokerAuthUsername = c.String("broker-auth-username")
	s.BrokerAuthPassword = c.String("broker-auth-password")
	s.BrokerCookieSecret = c.String("broker-cookie-secret")
	s.ServiceSecret = c.String("service-secret")
	s.CFAdminPassword = c.String("admin-password")

	var err error
	if err = s.flagValidation(); err != nil {
		lo.G.Error("invalid arguments: ", err)
		lo.G.Fatal("exiting due to invalid args")
	}

	if err = s.cloudconfigValidation(enaml.NewCloudConfigManifest(cloudConfig)); err != nil {
		lo.G.Error("invalid settings for cloud config on target bosh: ", err)
		lo.G.Fatal("your deployment is not compatible with your cloud config, exiting")
	}

	lo.G.Debug("context", c)
	var dm = new(enaml.DeploymentManifest)
	dm.SetName(s.DeploymentName)
	dm.AddRelease(enaml.Release{Name: CFMysqlReleaseName, Version: c.String("cf-mysql-release-version")})
	dm.AddRelease(enaml.Release{Name: MysqlBackupReleaseName, Version: c.String("mysql-backup-release-version")})
	dm.AddRelease(enaml.Release{Name: ServiceBackupReleaseName, Version: c.String("service-backup-release-version")})
	dm.AddRelease(enaml.Release{Name: MysqlMonitoringReleaseName, Version: c.String("mysql-monitoring-release-version")})
	dm.AddRemoteStemcell(s.StemcellName, s.StemcellName, s.StemcellVersion, s.StemcellURL, s.StemcellSHA)
	dm.AddInstanceGroup(NewMysqlPartition(s))
	dm.AddInstanceGroup(NewProxyPartition(s))
	dm.AddInstanceGroup(NewMonitoringPartition(s))
	dm.AddInstanceGroup(NewCfMysqlBrokerPartition(s))
	//dm.AddInstanceGroup(NewBackupPreparePartition())
	dm.AddInstanceGroup(NewBrokerRegistrar(s))
	dm.AddInstanceGroup(NewBrokerDeRegistrar(s))
	dm.AddInstanceGroup(NewRejoinUnsafe(s))
	dm.AddInstanceGroup(NewAcceptanceTests(s))
	dm.Update = enaml.Update{
		MaxInFlight:     1,
		UpdateWatchTime: "30000-300000",
		CanaryWatchTime: "30000-300000",
		Serial:          false,
		Canaries:        1,
	}
	return dm.Bytes()
}

func (s *Plugin) cloudconfigValidation(cloudConfig *enaml.CloudConfigManifest) (err error) {
	lo.G.Debug("running cloud config validation")

	for _, vmtype := range cloudConfig.VMTypes {
		err = fmt.Errorf("vm size %s does not exist in cloud config. options are: %v", s.VMTypeName, cloudConfig.VMTypes)
		if vmtype.Name == s.VMTypeName {
			err = nil
			break
		}
	}

	for _, disktype := range cloudConfig.DiskTypes {
		err = fmt.Errorf("disk size %s does not exist in cloud config. options are: %v", s.DiskTypeName, cloudConfig.DiskTypes)
		if disktype.Name == s.DiskTypeName {
			err = nil
			break
		}
	}

	for _, net := range cloudConfig.Networks {
		err = fmt.Errorf("network %s does not exist in cloud config. options are: %v", s.NetworkName, cloudConfig.Networks)
		if net.(map[interface{}]interface{})["name"] == s.NetworkName {
			err = nil
			break
		}
	}

	if len(cloudConfig.VMTypes) == 0 {
		err = fmt.Errorf("no vm sizes found in cloud config")
	}

	if len(cloudConfig.DiskTypes) == 0 {
		err = fmt.Errorf("no disk sizes found in cloud config")
	}

	if len(cloudConfig.Networks) == 0 {
		err = fmt.Errorf("no networks found in cloud config")
	}
	return
}

func (s *Plugin) flagValidation() (err error) {
	lo.G.Debug("validating given flags")

	if len(s.IPs) <= 0 {
		err = fmt.Errorf("no `ip` given")
	}

	if len(s.AZs) <= 0 {
		err = fmt.Errorf("no `az` given")
	}

	if s.NetworkName == "" {
		err = fmt.Errorf("no `network-name` given")
	}

	if s.VMTypeName == "" {
		err = fmt.Errorf("no `vm-type` given")
	}
	if s.DiskTypeName == "" {
		err = fmt.Errorf("no `disk-type` given")
	}

	if s.StemcellVersion == "" {
		err = fmt.Errorf("no `stemcell-ver` given")
	}
	return
}

func ertVaultDecorate(flags []pcli.Flag, hash string, v *pluginutil.VaultUnmarshal) {
	err := v.UnmarshalSomeFlags(hash, flags,
		"syslog-address",
		"uaa-admin-secret",
		"notifications-client-secret",
		"nats-pass",
		"admin-password")
	if err != nil {
		lo.G.Errorf("Error unmarshalling ERT flags: %s", err.Error())
	}
}

func vaultRotateMySQL(hash string, v pluginutil.VaultRotater) error {
	secrets := map[string]string{
		"mysql-admin-password":           utils.NewPassword(16),
		"seeded-db-password":             utils.NewPassword(16),
		"galera-healthcheck-username":    utils.NewPassword(16),
		"galera-healthcheck-password":    utils.NewPassword(16),
		"galera-healthcheck-db-password": utils.NewPassword(16),
		"cluster-health-password":        utils.NewPassword(16),
		"backup-endpoint-user":           utils.NewPassword(16),
		"backup-endpoint-password":       utils.NewPassword(16),
		"broker-quota-enforcer-password": utils.NewPassword(16),
		"proxy-api-password":             utils.NewPassword(16),
		"broker-auth-username":           utils.NewPassword(16),
		"broker-auth-password":           utils.NewPassword(16),
		"broker-cookie-secret":           utils.NewPassword(16),
		"service-secret":                 utils.NewPassword(16),
	}

	b, err := json.Marshal(secrets)
	if err != nil {
		return err
	}
	lo.G.Debug("rotating secrets for hash", hash)
	return v.RotateSecrets(hash, b)
}

func InferFromCloudDecorate(inferFlagMap map[string][]string, cloudConfig []byte, args []string, flgs []pcli.Flag) {
	c := pluginutil.NewContext(args, pluginutil.ToCliFlagArray(flgs))

	if c.Bool("infer-from-cloud") {
		ccinf := pluginutil.NewCloudConfigInferFromBytes(cloudConfig)
		setAllInferredFlagDefaults(inferFlagMap["disktype"], ccinf.InferDefaultDiskType(), flgs)
		setAllInferredFlagDefaults(inferFlagMap["vmtype"], ccinf.InferDefaultVMType(), flgs)
		setAllInferredFlagDefaults(inferFlagMap["az"], ccinf.InferDefaultAZ(), flgs)
		setAllInferredFlagDefaults(inferFlagMap["network"], ccinf.InferDefaultNetwork(), flgs)
	}
}

func setAllInferredFlagDefaults(matchlist []string, defaultvalue string, flgs []pcli.Flag) {
	for _, match := range matchlist {
		setFlagDefault(match, defaultvalue, flgs)
	}
}

func setFlagDefault(flagname, defaultvalue string, flgs []pcli.Flag) {
	for idx, flg := range flgs {
		if flg.Name == flagname {
			flgs[idx].Value = defaultvalue
		}
	}
}

var flagsToInferFromCloudConfig = map[string][]string{
	"disktype": []string{
		"disk-type",
	},
	"vmtype": []string{
		"vm-type",
	},
	"az":      []string{"az"},
	"network": []string{"network"},
}
