package cloudfoundry

import (
	"fmt"

	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-product-bundle/products/oss_cf/enaml-gen/route_registrar"
	"github.com/enaml-ops/omg-product-bundle/products/oss_cf/enaml-gen/uaa"
	"github.com/enaml-ops/omg-product-bundle/products/oss_cf/plugin/config"
)

//UAA -
type UAA struct {
	Config         *config.Config
	Metron         *Metron
	StatsdInjector *StatsdInjector
	ConsulAgent    *ConsulAgent
	Login          *uaa.Login
	UAA            *uaa.Uaa
}

//NewUAAPartition -
func NewUAAPartition(config *config.Config) InstanceGroupCreator {
	UAA := &UAA{
		Config:         config,
		Metron:         NewMetron(config),
		ConsulAgent:    NewConsulAgent([]string{"uaa"}, config),
		StatsdInjector: NewStatsdInjector(nil),
	}
	UAA.Login = UAA.CreateLogin()
	UAA.UAA = UAA.CreateUAA()
	return UAA
}

//CreateUAA - Helper method to create uaa structure
func (s *UAA) CreateUAA() (login *uaa.Uaa) {
	clientMap := make(map[string]UAAClient)
	clientMap["opentsdb-firehose-nozzle"] = UAAClient{
		AccessTokenValidity:  1209600,
		AuthorizedGrantTypes: "authorization_code,client_credentials,refresh_token",
		Override:             true,
		Secret:               s.Config.OpentsdbFirehoseNozzleClientSecret,
		Scope:                "openid,oauth.approvals,doppler.firehose",
		Authorities:          "oauth.login,doppler.firehose",
	}
	clientMap["identity"] = UAAClient{
		ID:                   "identity",
		Secret:               s.Config.IdentityClientSecret,
		Scope:                "cloud_controller.admin,cloud_controller.read,cloud_controller.write,openid,zones.*.*,zones.*.*.*,zones.read,zones.write,scim.read",
		ResourceIDs:          "none",
		AuthorizedGrantTypes: "authorization_code,client_credentials,refresh_token",
		AutoApprove:          true,
		Authorities:          "scim.zones,zones.read,uaa.resource,zones.write,cloud_controller.admin",
		RedirectURI:          fmt.Sprintf("%s://p-identity.%s/dashboard/,%s://p-identity.%s/dashboard/**", s.Config.UAALoginProtocol, s.Config.SystemDomain, s.Config.UAALoginProtocol, s.Config.SystemDomain),
	}
	clientMap["login"] = UAAClient{
		ID:                   "login",
		Secret:               s.Config.LoginClientSecret,
		AutoApprove:          true,
		Override:             true,
		Authorities:          "oauth.login,scim.write,clients.read,notifications.write,critical_notifications.write,emails.write,scim.userids,password.write",
		AuthorizedGrantTypes: "authorization_code,client_credentials,refresh_token",
		Scope:                "openid,oauth.approvals",
	}
	clientMap["portal"] = UAAClient{
		AuthorizedGrantTypes: "authorization_code,client_credentials,password,implicit",
		ID:                   "portal",
		Secret:               s.Config.PortalClientSecret,
		Override:             true,
		AutoApprove:          true,
		Authorities:          "scim.write,scim.read,cloud_controller.read,cloud_controller.write,password.write,uaa.admin,uaa.resource,cloud_controller.admin,emails.write,notifications.write",
		Scope:                "openid,cloud_controller.read,cloud_controller.write,password.write,console.admin,console.support,cloud_controller.admin",
		AccessTokenValidity:  1209600,
		RefreshTokenValidity: 1209600,
		Name:                 "Pivotal Apps Manager",
		AppLaunchURL:         fmt.Sprintf("%s://apps.%s", s.Config.UAALoginProtocol, s.Config.SystemDomain),
		ShowOnHomepage:       true,
		AppIcon:              "iVBORw0KGgoAAAANSUhEUgAAAGwAAABsCAYAAACPZlfNAAAAAXNSR0IArs4c6QAABYtJREFUeAHtnVtsFFUYx7/d3ruWotUKVIkNaCw02YgJGBRTMd4CokUejD4QH4gxQcIDeHnBmPjkhSghUYLGe3ywPtAHNCo0QgkWwi2tXG2V1kIpLXTbLt1tS9dzlmzSJssZhv32zDk7/2km2znn7Pd9+/vt2Z2dmW0D9Obat4gCiwiLBQQSLflSViAQeN6Can1fYiJBFPQ9BcsAQBiEWUbAsnIxwyDMMgKWlYsZBmGWEbCsXMwwCLOMgGXlYoZBmGUELCsXMwzCLCNgWbmYYRBmGQHLysUMgzDLCFhWLmYYhFlGwLJyMcMgzDIClpWLGQZhlhGwrFzMMAizjIBl5WKGQZhlBCwrV1xbb96y59V1VFJQmLawQNrWa43x8XEaHo1fW+Oj1H8lSqf6eulEbw+dvNhLvcNDinvb0WWksAdm3UWhwiJ2gt2RAWo80UY7jrdSU8cZGrt6lT1HtgMaKSxbD7qqfDq99tAjyTUSG6FP9v1BH+3dTUPxeLZSssf17U5HeXEJbXr8aerY+A6tf7iOxFeu2OFmI6BvhaVgVoRCtHl5PTW8/AoV5xekmo299b2wlJn6+WFqWrOWKkpDqSYjbyFskpZFs++hL1e9NKnFvF+t3OmQOwzdkcgUmnnBABXm5Ys1j8qKisVadFPvS8tramn1goX09eEDU+KbsmGlsMbjbbT6x++UDOVORGXoFppXOYMerLqbVsyrpcWzqykYdH5R+fjZlcnd/8sjV5Q5vOh0rt6LqhhyJsQ3uC+ID8ry89aHYtf90W1bKLzlffr19EnH6HIP8oXasOM4LwbkrLB0MP+6cJ6e+eoz+vTP5nTdU9peDC+Ysm3Khq+ESehy5r3e2ECHu7uUDuqq59Id4iXVtMV3wqSACSHt3V2/KF3I97qayjuVY7zo9KUwCfq3M6coNjamZD6zrFzZ70Wnb4XFxseoK3JZyXzWtGnKfi86fStMwu6LRpXMZ5RBmBKQ7k75XqZa8gLmPZ/Nq0hFkLnvttJSZUT5Oc60xbfC5CGs6lsrlD56hgaV/V50+lbYkuo5VFygPp3SMwxhXjwp0+bcsGRp2vZU48TEBB09153aNObWlzNMHo1/6r4apYTmsx10MTqsHONFp5VH6zMBtWbhYtq6YpVjiJ/ajjmO8WKAL4QFxamWZffPT1678dicex05D4jTKj8cO+Q4zosBOSXs7bonktci5ovjgPIUye3ieo3wzKrk+TC5faPLGz83On6ovtFY3ONySth7Ty67qbPMk6Hu+edv+vzg/slNRv3uy52O6xk40HWW6r/94nrdRrTn1AzLhOju9tP03DfbKTo6mkmYrN/X98L6xQHgTb/vpG0t+5LnybJOPMMEvhXWOXCJvj9yiD7Yu4sGRkYyxKjv7r4RJi+Na+05Rwf/66SG1qO0v/NffZQZM+WUsI07d1BC/MTE144GYzHxJYcYDYq1vb/f8WQlI9OshsopYZubm7IKy4Tg2K03wYKLGiDMBSwThkKYCRZc1ABhLmCZMBTCTLDgogYIcwHLhKEQZoIFFzVAmAtYJgyFMBMsuKgBwlzAMmEohJlgwUUNEOYClglDIcwECy5qgDAXsEwYCmEmWHBRA4S5gGXCUAgzwYKLGow84yyvuyhR/GW19kt9Lh5ibg01UtjS7VtzizLjo8FLIiNMHaEgTAdlxhwQxghTRygI00GZMQeEMcLUEQrCdFBmzAFhjDB1hIIwHZQZc0AYI0wdoSBMB2XGHBDGCFNHKAjTQZkxB4QxwtQRCsJ0UGbMAWGMMHWEgjAdlBlzQBgjTB2hIEwHZcYcEMYIU0coCNNBmTEHhDHC1BEKwnRQZswBYYwwdYSCMB2UGXNAGCNMHaEgTAdlxhziUu1Ei8M/+WFMh1CZEUi0/A+j7hNSB5Wo2wAAAABJRU5ErkJggg==",
	}
	clientMap["apps_manager_js"] = UAAClient{
		Override:             true,
		AutoApprove:          []string{"cloud_controller.read", "cloud_controller.write", "cloud_controller.admin"},
		Scope:                "cloud_controller.read,cloud_controller.write,cloud_controller.admin",
		AuthorizedGrantTypes: "implicit",
		AccessTokenValidity:  28800,
	}
	clientMap["cf"] = UAAClient{
		ID:                   "cf",
		Override:             true,
		Authorities:          "uaa.none",
		AuthorizedGrantTypes: "password,refresh_token",
		Scope:                "cloud_controller.read,cloud_controller.write,openid,password.write,cloud_controller.admin,scim.read,scim.write,doppler.firehose,uaa.user",
		AccessTokenValidity:  7200,
		RefreshTokenValidity: 1209600,
	}
	clientMap["autoscaling_service"] = UAAClient{
		ID:                   "autoscaling_service",
		Secret:               s.Config.AutoScalingServiceClientSecret,
		Override:             true,
		AutoApprove:          true,
		Authorities:          "cloud_controller.write,cloud_controller.read,cloud_controller.admin,notifications.write,critical_notifications.write,emails.write",
		AuthorizedGrantTypes: "client_credentials,authorization_code,refresh_token",
		Scope:                "openid,cloud_controller.permissions,cloud_controller.read,cloud_controller.write",
		AccessTokenValidity:  3600,
	}
	clientMap["system_passwords"] = UAAClient{
		ID:                   "system_passwords",
		Secret:               s.Config.SystemPasswordsClientSecret,
		Override:             true,
		AutoApprove:          true,
		Authorities:          "uaa.admin,scim.read,scim.write,password.write",
		AuthorizedGrantTypes: "client_credentials",
	}
	clientMap["cc-service-dashboards"] = UAAClient{
		ID:                   "cc-service-dashboards",
		Secret:               s.Config.CCServiceDashboardsClientSecret,
		Override:             true,
		Authorities:          "clients.read,clients.write,clients.admin",
		AuthorizedGrantTypes: "client_credentials",
		Scope:                "cloud_controller.write,openid,cloud_controller.read,cloud_controller_service_permissions.read",
	}
	clientMap["doppler"] = UAAClient{
		ID:          "doppler",
		Secret:      s.Config.DopplerSecret,
		Authorities: "uaa.resource",
	}
	clientMap["gorouter"] = UAAClient{
		ID:                   "gorouter",
		Secret:               s.Config.GoRouterClientSecret,
		Authorities:          "clients.read,clients.write,clients.admin,routing.routes.write,routing.routes.read",
		AuthorizedGrantTypes: "client_credentials,refresh_token",
		Scope:                "openid,cloud_controller_service_permissions.read",
	}
	clientMap["notifications"] = UAAClient{
		ID:                   "notifications",
		Secret:               s.Config.NotificationsClientSecret,
		Authorities:          "cloud_controller.admin,scim.read,notifications.write,critical_notifications.write,emails.write",
		AuthorizedGrantTypes: "client_credentials",
	}
	clientMap["notifications_template"] = UAAClient{
		ID:                   "notifications_template",
		Secret:               "bb6be96896c5ab64c897",
		Scope:                "openid,clients.read,clients.write,clients.secret",
		Authorities:          "clients.read,clients.write,clients.secret,notification_templates.write,notification_templates.read,notifications.manage",
		AuthorizedGrantTypes: "client_credentials",
	}
	clientMap["notifications_ui_client"] = UAAClient{
		ID:                   "notifications_ui_client",
		Secret:               s.Config.NotificationsUIClientSecret,
		Scope:                "notification_preferences.read,notification_preferences.write,openid",
		AuthorizedGrantTypes: "authorization_code,client_credentials,refresh_token",
		Authorities:          "notification_preferences.admin",
		AutoApprove:          true,
		Override:             true,
		RedirectURI:          fmt.Sprintf("%s://notifications-ui.%s/sessions/create", s.Config.UAALoginProtocol, s.Config.SystemDomain),
	}
	clientMap["cloud_controller_username_lookup"] = UAAClient{
		ID:                   "cloud_controller_username_lookup",
		Secret:               s.Config.CloudControllerUsernameLookupClientSecret,
		AuthorizedGrantTypes: "client_credentials",
		Authorities:          "scim.userids",
	}
	clientMap["cc_routing"] = UAAClient{
		Authorities:          "routing.router_groups.read",
		AuthorizedGrantTypes: "client_credentials",
		Secret:               s.Config.CCRoutingClientSecret,
	}
	clientMap["ssh-proxy"] = UAAClient{
		AuthorizedGrantTypes: "authorization_code",
		AutoApprove:          true,
		Override:             true,
		RedirectURI:          "/login",
		Scope:                "openid,cloud_controller.read,cloud_controller.write",
		Secret:               s.Config.SSHProxyClientSecret,
	}
	clientMap["apps_metrics"] = UAAClient{
		ID:                   "apps_metrics",
		Secret:               s.Config.AppsMetricsClientSecret,
		Override:             true,
		AuthorizedGrantTypes: "authorization_code,refresh_token",
		RedirectURI:          fmt.Sprintf("%s://apm.%s,%s://apm.%s/,%s://apm.%s/*,%s://metrics.%s,%s://metrics.%s/,%s://metrics.%s/*", s.Config.UAALoginProtocol, s.Config.SystemDomain, s.Config.UAALoginProtocol, s.Config.SystemDomain, s.Config.UAALoginProtocol, s.Config.SystemDomain, s.Config.UAALoginProtocol, s.Config.SystemDomain, s.Config.UAALoginProtocol, s.Config.SystemDomain, s.Config.UAALoginProtocol, s.Config.SystemDomain),
		Scope:                "cloud_controller.admin,cloud_controller.read,metrics.read",
		AccessTokenValidity:  900,
		RefreshTokenValidity: 2628000,
	}
	clientMap["apps_metrics_processing"] = UAAClient{
		ID:                   "apps_metrics_processing",
		Secret:               s.Config.AppsMetricsProcessingClientSecret,
		Override:             true,
		AuthorizedGrantTypes: "authorization_code,client_credentials,refresh_token",
		Authorities:          "oauth.login,doppler.firehose,cloud_controller.admin",
		Scope:                "openid,oauth.approvals,doppler.firehose,cloud_controller.admin",
		AccessTokenValidity:  1209600,
	}
	return &uaa.Uaa{
		RequireHttps: true,
		Ssl: &uaa.UaaSsl{
			Port: -1,
		},
		Admin: &uaa.Admin{
			ClientSecret: s.Config.AdminSecret,
		},
		Authentication: &uaa.Authentication{
			Policy: &uaa.AuthenticationPolicy{
				LockoutAfterFailures: 5,
			},
		},
		Password: &uaa.UaaPassword{
			Policy: &uaa.PasswordPolicy{
				MinLength:                 0,
				RequireLowerCaseCharacter: 0,
				RequireUpperCaseCharacter: 0,
				RequireDigit:              0,
				RequireSpecialCharacter:   0,
				ExpirePasswordInMonths:    0,
			},
		},

		Ldap: &uaa.Ldap{
			ProfileType:         "search-and-bind",
			Url:                 s.Config.LDAPUrl,
			UserDN:              s.Config.LDAPUserDN,
			UserPassword:        s.Config.LDAPUserPassword,
			SearchBase:          s.Config.LDAPSearchBase,
			SearchFilter:        s.Config.LDAPSearchFilter,
			SslCertificate:      "",
			SslCertificateAlias: "",
			MailAttributeName:   s.Config.LDAPMailAttributeName,
			Enabled:             s.Config.LDAPEnabled,
			Groups: &uaa.LdapGroups{
				ProfileType:       "no-groups",
				SearchBase:        "",
				GroupSearchFilter: "",
			},
		},
		CatalinaOpts: "-Xmx768m -XX:MaxPermSize=256m",
		Url:          fmt.Sprintf("%s://uaa.%s", s.Config.UAALoginProtocol, s.Config.SystemDomain),
		Jwt: &uaa.Jwt{
			SigningKey:      s.Config.JWTSigningKey,
			VerificationKey: s.Config.JWTVerificationKey,
		},
		Proxy: &uaa.Proxy{
			Servers: s.Config.RouterMachines,
		},
		Clients: clientMap,
		Scim: &uaa.Scim{
			User: &uaa.ScimUser{
				Override: true,
			},
			UseridsEnabled: true,
			Users: []UAAScimUser{
				UAAScimUser{
					Name:     "admin",
					Password: s.Config.AdminPassword,
					Groups: []string{
						"scim.write",
						"scim.read",
						"openid",
						"cloud_controller.admin",
						"dashboard.user",
						"console.admin",
						"console.support",
						"doppler.firehose",
						"notification_preferences.read",
						"notification_preferences.write",
						"notifications.manage",
						"notification_templates.read",
						"notification_templates.write",
						"emails.write",
						"notifications.write",
						"zones.read",
						"zones.write",
					},
				},
				UAAScimUser{
					Name:     "push_apps_manager",
					Password: s.Config.PushAppsManagerPassword,
					Groups:   []string{"cloud_controller.admin"},
				},
				UAAScimUser{
					Name:     "smoke_tests",
					Password: s.Config.SmokeTestsPassword,
					Groups:   []string{"cloud_controller.admin"},
				},
				UAAScimUser{
					Name:     "system_services",
					Password: s.Config.SystemServicesPassword,
					Groups:   []string{"cloud_controller.admin"},
				},
				UAAScimUser{
					Name:     "system_verification",
					Password: s.Config.SystemVerificationPassword,
					Groups: []string{
						"scim.write",
						"scim.read",
						"openid",
						"cloud_controller.admin",
						"dashboard.user",
						"console.admin",
						"console.support",
					},
				},
			},
		},
	}
}

type UAAScimUser struct {
	Name     string   `yaml:"name,omitempty"`
	Password string   `yaml:"password,omitempty"`
	Groups   []string `yaml:"groups,omitempty"`
}

//CreateLogin - Helper method to create login structure
func (s *UAA) CreateLogin() (login *uaa.Login) {
	return &uaa.Login{
		Branding:                s.CreateBranding(),
		SelfServiceLinksEnabled: s.Config.SelfServiceLinksEnabled,
		Protocol:                s.Config.UAALoginProtocol,
		Links: &uaa.Links{
			Signup: fmt.Sprintf("%s://login.%s/create_account", s.Config.UAALoginProtocol, s.Config.SystemDomain),
			Passwd: fmt.Sprintf("%s://login.%s/forgot_password", s.Config.UAALoginProtocol, s.Config.SystemDomain),
		},
		Url: fmt.Sprintf("%s://uaa.%s", s.Config.UAALoginProtocol, s.Config.SystemDomain),
		Notifications: &uaa.Notifications{
			Url: fmt.Sprintf("%s://notifications.%s", s.Config.UAALoginProtocol, s.Config.SystemDomain),
		},
		Saml: &uaa.Saml{
			Entityid:                   fmt.Sprintf("%s://login.%s", s.Config.UAALoginProtocol, s.Config.SystemDomain),
			ServiceProviderKey:         s.Config.SAMLServiceProviderKey,
			ServiceProviderCertificate: s.Config.SAMLServiceProviderCertificate,
			SignRequest:                true,
			WantAssertionSigned:        false,
		},
		Logout: &uaa.Logout{
			Redirect: &uaa.Redirect{
				Parameter: &uaa.Parameter{
					Disable:   false,
					Whitelist: []string{fmt.Sprintf("%s://console.%s", s.Config.UAALoginProtocol, s.Config.SystemDomain), fmt.Sprintf("%s://apps.%s", s.Config.UAALoginProtocol, s.Config.SystemDomain)},
				},
				Url: "/login",
			},
		},
	}

}

func (s *UAA) CreateBranding() (branding *uaa.Branding) {
	branding = &uaa.Branding{
		CompanyName:     s.Config.CompanyName,
		ProductLogo:     s.Config.ProductLogo,
		SquareLogo:      s.Config.SquareLogo,
		FooterLegalText: s.Config.FooterLegalText,
	}
	return
}

//ToInstanceGroup -
func (s *UAA) ToInstanceGroup() (ig *enaml.InstanceGroup) {
	ig = &enaml.InstanceGroup{
		Name:      "uaa-partition",
		Instances: s.Config.UAAInstances,
		VMType:    s.Config.UAAVMType,
		AZs:       s.Config.AZs,
		Stemcell:  s.Config.StemcellName,
		Jobs: []enaml.InstanceJob{
			s.createUAAJob(),
			s.Metron.CreateJob(),
			s.ConsulAgent.CreateJob(),
			s.StatsdInjector.CreateJob(),
			s.createRouteRegistrarJob(),
		},
		Networks: []enaml.Network{
			enaml.Network{Name: s.Config.NetworkName},
		},
		Update: enaml.Update{
			MaxInFlight: 1,
		},
	}
	return
}

func (s *UAA) createRouteRegistrarJob() enaml.InstanceJob {
	return enaml.InstanceJob{
		Name:    "route_registrar",
		Release: "cf",
		Properties: &route_registrar.RouteRegistrarJob{
			RouteRegistrar: &route_registrar.RouteRegistrar{
				Routes: []map[string]interface{}{
					map[string]interface{}{
						"name":                  "uaa",
						"port":                  8080,
						"registration_interval": "40s",
						"uris":                  []string{fmt.Sprintf("uaa.%s", s.Config.SystemDomain), fmt.Sprintf("*.uaa.%s", s.Config.SystemDomain), fmt.Sprintf("login.%s", s.Config.SystemDomain), fmt.Sprintf("*.login.%s", s.Config.SystemDomain)},
					},
				},
			},
			Nats: &route_registrar.Nats{
				User:     s.Config.NATSUser,
				Password: s.Config.NATSPassword,
				Machines: s.Config.NATSMachines,
				Port:     s.Config.NATSPort,
			},
		},
	}
}

func (s *UAA) createUAAJob() enaml.InstanceJob {
	return enaml.InstanceJob{
		Name:    "uaa",
		Release: "cf",
		Properties: &uaa.UaaJob{
			Login: s.Login,
			Uaa:   s.UAA,
			Uaadb: s.createUAADB(),
		},
	}
}

func (s *UAA) createUAADB() (uaadb *uaa.Uaadb) {
	const uaaVal = "uaa"

	return &uaa.Uaadb{
		Address:  s.Config.MySQLProxyHost(),
		Port:     3306,
		DbScheme: "mysql",
		Roles: []map[string]interface{}{
			map[string]interface{}{
				"tag":      "admin",
				"name":     s.Config.UAADBUserName,
				"password": s.Config.UAADBPassword,
			},
		},
		Databases: []map[string]interface{}{
			map[string]interface{}{
				"tag":  uaaVal,
				"name": uaaVal,
			},
		},
	}
}
