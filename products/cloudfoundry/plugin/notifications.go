package cloudfoundry

import (
	"fmt"

	"github.com/enaml-ops/enaml"
	dn "github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/deploy-notifications"
	dnui "github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/deploy-notifications-ui"
	tn "github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/test-notifications"
	tnui "github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/test-notifications-ui"
	"github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/plugin/config"
)

type notifications struct{ *config.Config }

func NewNotifications(c *config.Config) InstanceGroupCreator {
	return notifications{c}
}

func (n notifications) ToInstanceGroup() *enaml.InstanceGroup {
	return &enaml.InstanceGroup{
		Name:      "notifications",
		Instances: 1,
		VMType:    n.ErrandVMType,
		Lifecycle: "errand",
		AZs:       n.AZs,
		Stemcell:  n.StemcellName,
		Networks: []enaml.Network{
			{Name: n.NetworkName},
		},
		Update: enaml.Update{
			MaxInFlight: 1,
		},
		Jobs: []enaml.InstanceJob{
			{
				Name:    "deploy-notifications",
				Release: NotificationsReleaseName,
				Properties: &dn.DeployNotificationsJob{
					Domain: n.SystemDomain,
					Ssl: &dn.Ssl{
						SkipCertVerify: n.SkipSSLCertVerify,
					},
					Notifications: &dn.Notifications{
						AppDomain:               n.SystemDomain, // yes, the system domain
						Network:                 "notifications",
						EncryptionKey:           n.DbEncryptionKey,
						EnableDiego:             true,
						InstanceCount:           3,
						SyslogUrl:               n.SyslogPort, // really??
						Organization:            "system",
						Space:                   "notifications-with-ui",
						Sender:                  "", // TODO new flag?
						ErrorOnMisconfiguration: false,
						Cf: &dn.Cf{
							AdminUser:     "admin",
							AdminPassword: n.AdminPassword,
						},
						Smtp: &dn.Smtp{
							Host:          "", // TODO need new flag?
							Port:          25,
							Tls:           false,
							AuthMechanism: "none",
						},
						Uaa: &dn.Uaa{
							AdminClientId:     "admin",
							AdminClientSecret: n.AdminSecret,
							ClientId:          "notifications",
							ClientSecret:      n.NotificationsClientSecret,
						},
						Database: &dn.Database{
							Url:                fmt.Sprintf("mysql://%s:%s@%s:3306/notifications", n.NotificationsDBUser, n.NotificationsDBPassword, n.MySQLProxyHost()),
							MaxOpenConnections: 10,
						},
						DefaultTemplate: `{
              "name": "Default Template",
              "subject": "CF Notification: {{.Subject}}",
              "html": "\u003ctable width=\"100%\" cellpadding=\"0\" cellspacing=\"0\" border=\"0\" style=\"border-collapse:collapse;font-family: Helvetica, Arial,sans-serif\"\u003e\n    \u003ctbody\u003e\n        \u003ctr\u003e\n            \u003ctd width=\"100%\" align=\"center\" bgcolor=\"#F8f8f8\" style=\"padding-right:10px;padding-left:10px\"\u003e\n                \u003ctable width=\"600\" cellpadding=\"0\" cellspacing=\"0\" border=\"0\" style=\"border-collapse:collapse\"\u003e\n                    \u003ctbody\u003e\n                        \u003ctr\u003e\n                            \u003ctd align=\"left\" valign=\"top\" style=\"padding-bottom:30px;padding-top:30px\"\u003e\n                              \u003cimg src=\"http://notifications-ui.{{.Domain}}/assets/pivotal_logo.png\"\n                                alt=\"Pivotal CF\" border=\"0\" style=\"display:block\" class=\"CToWUd\"\n                                width=\"166px\" height=\"35px\"\u003e\n                            \u003c/td\u003e\n                        \u003c/tr\u003e\n                    \u003c/tbody\u003e\n                \u003c/table\u003e\n                \u003ctable width=\"600\" cellpadding=\"0\" cellspacing=\"0\" border=\"0\" style=\"border-collapse:collapse\"\u003e\n                    \u003ctbody\u003e\n                        \u003ctr\u003e\n                            \u003ctd align=\"left\" valign=\"top\" bgcolor=\"#FFFFFF\"\u003e\n                                \u003ctable width=\"600\" cellpadding=\"40\" cellspacing=\"0\" border=\"0\" style=\"border-collapse:collapse;border:1px solid #e0e4e5\"\u003e\n                                    \u003ctbody\u003e\n                                        \u003ctr\u003e\n                                            \u003ctd style=\"color:#666666;font-size:16px;line-height:22px\"\u003e\n                                                \u003cdiv\u003e\n                                                    {{.HTML}}\n                                                \u003c/div\u003e\n                                            \u003c/td\u003e\n                                        \u003c/tr\u003e\n                                    \u003c/tbody\u003e\n                                \u003c/table\u003e\n                            \u003c/td\u003e\n                        \u003c/tr\u003e\n                    \u003c/tbody\u003e\n                \u003c/table\u003e\n                \u003ctable width=\"600\" cellpadding=\"0\" cellspacing=\"0\" border=\"0\" style=\"border-collapse:collapse;font-size:12px;color:#b4b4b4\"\u003e\n                    \u003ctbody\u003e\n                        \u003ctr\u003e\n                          \u003ctd align=\"left\" valign=\"middle\" style=\"padding-top:20px;padding-bottom:20px\"\u003eManage your \u003ca href=\"https://notifications-ui.{{.Domain}}/preferences\"\n                                target=\"_blank\"\u003enotification preferences\u003c/a\u003e or\n                              \u003ca href=\"https://notifications-ui.{{.Domain}}/unsubscribe/{{.UnsubscribeID}}\"\n                                target=\"_blank\"\u003eunsubscribe\u003c/a\u003e\n                            \u003c/td\u003e\n                        \u003c/tr\u003e\n                        \u003ctr\u003e\n                            \u003ctd align=\"left\" valign=\"middle\" style=\"padding-bottom:20px\"\u003ePivotal Cloud Foundry, and the Pivotal Cloud Foundry logo are registered\n                                trademarks or trademarks of Pivotal Software, Inc.\n                                in the United States and other countries. All other\n                                trademarks used herein are the property of their\n                                respective owners.\n                                \u003ca\u003e\n                            \u003c/td\u003e\n                        \u003c/tr\u003e\n                        \u003ctr\u003e\n                            \u003ctd align=\"left\" valign=\"middle\" style=\"padding-bottom:20px\"\u003e\u00A9 2015 Pivotal Software, Inc. All rights reserved.\n                                Published in the USA.\n                                \u003ca\u003e\n                            \u003c/td\u003e\n                        \u003c/tr\u003e\n                    \u003c/tbody\u003e\n                \u003c/table\u003e\n            \u003c/td\u003e\n        \u003c/tr\u003e\n    \u003c/tbody\u003e\n\u003c/table\u003e\n",
              "text": "_____\n\n\nPivotal CF\n\n{{.Text}}\n\nManage your notification preferences:\nhttps://notifications-ui.{{.Domain}}/preferences\n\nUnsubscribe from Pivotal emails:\nhttps://notifications-ui.{{.Domain}}/unsubscribe/{{.UnsubscribeID}}\n\nPivotal Cloud Foundry, and the Pivotal Cloud Foundry logo are registered trademarks or trademarks of Pivotal Software, Inc. in the United States and other countries. All other trademarks used herein are the property of their respective owners.\n\n\u00A9 2015 Pivotal Software, Inc. All rights reserved. Published in the USA.\n\n______\n",
              "metadata": {}
              }`,
					},
				},
			},
		},
	}
}

type notificationsTest struct{ *config.Config }

func NewNotificationsTest(c *config.Config) InstanceGroupCreator {
	return notificationsTest{c}
}

func (n notificationsTest) ToInstanceGroup() *enaml.InstanceGroup {
	return &enaml.InstanceGroup{
		Name:      "notifications-tests",
		Instances: 1,
		VMType:    n.ErrandVMType,
		Lifecycle: "errand",
		AZs:       n.AZs,
		Stemcell:  n.StemcellName,
		Networks: []enaml.Network{
			{Name: n.NetworkName},
		},
		Update: enaml.Update{
			MaxInFlight: 1,
		},
		Jobs: []enaml.InstanceJob{
			{
				Name:    "test-notifications",
				Release: NotificationsReleaseName,
				Properties: &tn.TestNotificationsJob{
					Domain: n.SystemDomain,
					Notifications: &tn.Notifications{
						Cf: &tn.Cf{
							AdminUser:     "admin",
							AdminPassword: n.AdminPassword,
						},
						AppDomain: n.SystemDomain, // yes, the system domain
						Uaa: &tn.Uaa{
							AdminClientId:     "admin",
							AdminClientSecret: n.AdminSecret,
						},
						Organization: "system",
						Space:        "notifications-with-ui",
					},
				},
			},
		},
	}
}

type notificationsUI struct{ *config.Config }

func NewNotificationsUI(c *config.Config) InstanceGroupCreator {
	return notificationsUI{c}
}

func (n notificationsUI) ToInstanceGroup() *enaml.InstanceGroup {
	return &enaml.InstanceGroup{
		Name:      "notifications-ui",
		Instances: 1,
		VMType:    n.ErrandVMType,
		Lifecycle: "errand",
		AZs:       n.AZs,
		Stemcell:  n.StemcellName,
		Networks: []enaml.Network{
			{Name: n.NetworkName},
		},
		Update: enaml.Update{
			MaxInFlight: 1,
		},
		Jobs: []enaml.InstanceJob{
			{
				Name:    "deploy-notifications-ui",
				Release: NotificationsUIReleaseName,
				Properties: &dnui.DeployNotificationsUiJob{
					Domain: n.SystemDomain,
					Ssl: &dnui.Ssl{
						SkipCertVerify: n.SkipSSLCertVerify,
					},
					NotificationsUi: &dnui.NotificationsUi{
						Network:       "notifications",
						SyslogUrl:     n.SyslogPort, // TODO ??
						EncryptionKey: n.DbEncryptionKey,
						EnableDiego:   true,
						InstanceCount: 1,
						Cf: &dnui.Cf{
							AdminUser:     "admin",
							AdminPassword: n.AdminPassword,
						},
						AppDomain: n.SystemDomain, // yes, the system domain
						Uaa: &dnui.Uaa{
							ClientId:     "notifications_ui_client",
							ClientSecret: n.NotificationsUIClientSecret,
						},
						Organization: "system",
						Space:        "notifications-with-ui",
					},
				},
			},
		},
	}
}

type notificationsUITest struct{ *config.Config }

func NewNotificationsUITest(c *config.Config) InstanceGroupCreator {
	return notificationsUITest{c}
}

func (n notificationsUITest) ToInstanceGroup() *enaml.InstanceGroup {
	return &enaml.InstanceGroup{
		Name:      "notifications-ui-tests",
		Instances: 1,
		VMType:    n.ErrandVMType,
		Lifecycle: "errand",
		AZs:       n.AZs,
		Stemcell:  n.StemcellName,
		Networks: []enaml.Network{
			{Name: n.NetworkName},
		},
		Update: enaml.Update{
			MaxInFlight: 1,
		},
		Jobs: []enaml.InstanceJob{
			{
				Name:    "test-notifications-ui",
				Release: NotificationsUIReleaseName,
				Properties: &tnui.TestNotificationsUiJob{
					Domain: n.SystemDomain,
					NotificationsUi: &tnui.NotificationsUi{
						Cf: &tnui.Cf{
							AdminUser:     "admin",
							AdminPassword: n.AdminPassword,
						},
						AppDomain: n.SystemDomain, // yes, the system domain
						Uaa: &tnui.Uaa{
							AdminClient: "admin",
							AdminSecret: n.AdminSecret,
						},
						Organization: "system",
						Space:        "notifications-with-ui",
					},
				},
			},
		},
	}
}
