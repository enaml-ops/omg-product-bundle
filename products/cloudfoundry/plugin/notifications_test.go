package cloudfoundry_test

import (
	"fmt"

	"github.com/enaml-ops/enaml"
	dn "github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/deploy-notifications"
	dnui "github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/deploy-notifications-ui"
	tn "github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/test-notifications"
	tnui "github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/test-notifications-ui"
	. "github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/plugin"
	"github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/plugin/config"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("notifications", func() {
	var cfg *config.Config

	const (
		controlSystemDomain          = "sys.example.com"
		controlVMType                = "errand-vm-type"
		controlNetworkName           = "foundry-net"
		controlStemcellName          = "cool-ubuntu-animal"
		controlNotificationsDBUser   = "notificationsuser"
		controlNotificationsDBPass   = "notificationspass"
		controlCFAdminPass           = "cfadminpass"
		controlAdminSecret           = "adminsecret"
		controlEncryptionKey         = "dbencryptionkey"
		controlNotificationsUISecret = "notificationsuisecret"
		controlNotificationsSecret   = "notificationssecret"
		controlMySQLProxy            = "10.0.0.10"
	)

	BeforeEach(func() {
		cfg = &config.Config{
			VMType: config.VMType{
				ErrandVMType: controlVMType,
			},
			User: config.User{
				NotificationsDBUser: controlNotificationsDBUser,
			},
			Secret: config.Secret{
				NotificationsDBPassword:     controlNotificationsDBPass,
				AdminPassword:               controlCFAdminPass,
				AdminSecret:                 controlAdminSecret,
				DbEncryptionKey:             controlEncryptionKey,
				NotificationsClientSecret:   controlNotificationsSecret,
				NotificationsUIClientSecret: controlNotificationsUISecret,
			},
			IP: config.IP{
				MySQLProxyIPs: []string{controlMySQLProxy},
			},
			SystemDomain:      controlSystemDomain,
			NetworkName:       controlNetworkName,
			StemcellName:      controlStemcellName,
			AZs:               []string{"az1"},
			SkipSSLCertVerify: true,
		}
	})

	validateIG := func(ig *enaml.InstanceGroup) {
		Ω(ig.Instances).Should(Equal(1))
		Ω(ig.VMType).Should(Equal(controlVMType))
		Ω(ig.Lifecycle).Should(Equal("errand"))
		Ω(ig.Networks).Should(HaveLen(1))
		Ω(ig.Networks[0].Name).Should(Equal(controlNetworkName))
		Ω(ig.Stemcell).Should(Equal(controlStemcellName))
		Ω(ig.AZs).Should(ConsistOf("az1"))
		Ω(ig.Update.MaxInFlight).Should(Equal(1))
	}

	Context("notifications partition", func() {
		var (
			igc InstanceGroupCreator
			ig  *enaml.InstanceGroup
			dm  *enaml.DeploymentManifest
		)
		BeforeEach(func() {
			igc = NewNotifications(cfg)
			ig = igc.ToInstanceGroup()
			dm = new(enaml.DeploymentManifest)
			dm.AddInstanceGroup(ig)
			Ω(dm.InstanceGroups).Should(HaveLen(1))
		})

		It("configures the instance group", func() {
			Ω(ig.Name).Should(Equal("notifications"))
			validateIG(ig)
		})

		It("configures the deploy-notifications job", func() {
			job := ig.GetJobByName("deploy-notifications")
			Ω(job).ShouldNot(BeNil())
			Ω(job.Name).Should(Equal("deploy-notifications"))
			Ω(job.Release).Should(Equal(NotificationsReleaseName))

			props := job.Properties.(*dn.DeployNotificationsJob)

			Ω(props.Domain).Should(Equal(controlSystemDomain))

			Ω(props.Ssl).ShouldNot(BeNil())
			Ω(props.Ssl.SkipCertVerify).Should(BeTrue())

			Ω(props.Notifications.Network).Should(Equal("notifications"))
			Ω(props.Notifications.EncryptionKey).Should(Equal(controlEncryptionKey))
			Ω(props.Notifications.EnableDiego).Should(BeTrue())
			Ω(props.Notifications.InstanceCount).Should(Equal(3))

			// TODO	Ω(props.Notifications.SyslogUrl).Should(Equal())

			// TODO Ω(props.Notifications.Sender).Should(Equal())

			Ω(props.Notifications.Smtp).ShouldNot(BeNil())
			// TODO Ω(props.Notifications.Smtp.Host).Should(Equal())
			Ω(props.Notifications.Smtp.Port).Should(Equal(25))
			Ω(props.Notifications.Smtp.Tls).Should(BeFalse())
			Ω(props.Notifications.Smtp.AuthMechanism).Should(Equal("none"))

			Ω(props.Notifications.Database).ShouldNot(BeNil())
			Ω(props.Notifications.Database.Url).Should(Equal(fmt.Sprintf("mysql://%s:%s@%s:3306/notifications", controlNotificationsDBUser, controlNotificationsDBPass, controlMySQLProxy)))

			Ω(props.Notifications.Organization).Should(Equal("system"))
			Ω(props.Notifications.Space).Should(Equal("notifications-with-ui"))
			Ω(props.Notifications.ErrorOnMisconfiguration).Should(BeFalse())

			Ω(props.Notifications.Cf).ShouldNot(BeNil())
			Ω(props.Notifications.Cf.AdminUser).Should(Equal("admin"))
			Ω(props.Notifications.Cf.AdminPassword).Should(Equal(controlCFAdminPass))

			Ω(props.Notifications.Uaa).ShouldNot(BeNil())
			Ω(props.Notifications.Uaa.AdminClientId).Should(Equal("admin"))
			Ω(props.Notifications.Uaa.AdminClientSecret).Should(Equal(controlAdminSecret))
			Ω(props.Notifications.Uaa.ClientId).Should(Equal("notifications"))
			Ω(props.Notifications.Uaa.ClientSecret).Should(Equal(controlNotificationsSecret))
		})
	})

	Context("notifications-tests partition", func() {
		var (
			igc InstanceGroupCreator
			ig  *enaml.InstanceGroup
			dm  *enaml.DeploymentManifest
		)
		BeforeEach(func() {
			igc = NewNotificationsTest(cfg)
			ig = igc.ToInstanceGroup()
			dm = new(enaml.DeploymentManifest)
			dm.AddInstanceGroup(ig)
			Ω(dm.InstanceGroups).Should(HaveLen(1))
		})

		It("configures the instance group", func() {
			Ω(ig.Name).Should(Equal("notifications-tests"))
			validateIG(ig)
		})

		It("configures the test-notifications job", func() {
			job := ig.GetJobByName("test-notifications")
			Ω(job).ShouldNot(BeNil())
			Ω(job.Name).Should(Equal("test-notifications"))
			Ω(job.Release).Should(Equal(NotificationsReleaseName))

			props := job.Properties.(*tn.TestNotificationsJob)
			Ω(props.Domain).Should(Equal(controlSystemDomain))

			Ω(props.Notifications.Cf).ShouldNot(BeNil())
			Ω(props.Notifications.Cf.AdminUser).Should(Equal("admin"))
			Ω(props.Notifications.Cf.AdminPassword).Should(Equal(controlCFAdminPass))

			Ω(props.Notifications.AppDomain).Should(Equal(controlSystemDomain))

			Ω(props.Notifications.Uaa).ShouldNot(BeNil())
			Ω(props.Notifications.Uaa.AdminClientId).Should(Equal("admin"))
			Ω(props.Notifications.Uaa.AdminClientSecret).Should(Equal(controlAdminSecret))
		})
	})

	Context("notifications-ui partition", func() {
		var (
			igc InstanceGroupCreator
			ig  *enaml.InstanceGroup
			dm  *enaml.DeploymentManifest
		)
		BeforeEach(func() {
			igc = NewNotificationsUI(cfg)
			ig = igc.ToInstanceGroup()
			dm = new(enaml.DeploymentManifest)
			dm.AddInstanceGroup(ig)
			Ω(dm.InstanceGroups).Should(HaveLen(1))
		})

		It("configures the instance group", func() {
			Ω(ig.Name).Should(Equal("notifications-ui"))
			validateIG(ig)
		})

		It("configures the deploy-notifications-ui job", func() {
			job := ig.GetJobByName("deploy-notifications-ui")
			Ω(job).ShouldNot(BeNil())
			Ω(job.Name).Should(Equal("deploy-notifications-ui"))
			Ω(job.Release).Should(Equal(NotificationsUIReleaseName))

			props := job.Properties.(*dnui.DeployNotificationsUiJob)

			Ω(props.Domain).Should(Equal(controlSystemDomain))

			// TODO Ω(props.NotificationsUi.SyslogUrl).Should(Equal(""))
			Ω(props.NotificationsUi.Network).Should(Equal("notifications"))
			Ω(props.NotificationsUi.AppDomain).Should(Equal(controlSystemDomain))
			Ω(props.NotificationsUi.EncryptionKey).Should(Equal(controlEncryptionKey))
			Ω(props.NotificationsUi.EnableDiego).Should(BeTrue())
			Ω(props.NotificationsUi.InstanceCount).Should(Equal(1))
			Ω(props.NotificationsUi.Organization).Should(Equal("system"))
			Ω(props.NotificationsUi.Space).Should(Equal("notifications-with-ui"))

			Ω(props.NotificationsUi.Cf).ShouldNot(BeNil())
			Ω(props.NotificationsUi.Cf.AdminUser).Should(Equal("admin"))
			Ω(props.NotificationsUi.Cf.AdminPassword).Should(Equal(controlCFAdminPass))

			Ω(props.NotificationsUi.Uaa).ShouldNot(BeNil())
			Ω(props.NotificationsUi.Uaa.ClientId).Should(Equal("notifications_ui_client"))
			Ω(props.NotificationsUi.Uaa.ClientSecret).Should(Equal(controlNotificationsUISecret))

			Ω(props.Ssl).ShouldNot(BeNil())
			Ω(props.Ssl.SkipCertVerify).Should(BeTrue())
		})
	})

	Context("notifications-ui-tests partition", func() {
		var (
			igc InstanceGroupCreator
			ig  *enaml.InstanceGroup
			dm  *enaml.DeploymentManifest
		)
		BeforeEach(func() {
			igc = NewNotificationsUITest(cfg)
			ig = igc.ToInstanceGroup()
			dm = new(enaml.DeploymentManifest)
			dm.AddInstanceGroup(ig)
			Ω(dm.InstanceGroups).Should(HaveLen(1))
		})

		It("configures the instance group", func() {
			Ω(ig.Name).Should(Equal("notifications-ui-tests"))
			validateIG(ig)
		})

		It("configures the test-notifications-ui job", func() {
			job := ig.GetJobByName("test-notifications-ui")
			Ω(job).ShouldNot(BeNil())
			Ω(job.Name).Should(Equal("test-notifications-ui"))
			Ω(job.Release).Should(Equal(NotificationsUIReleaseName))

			props := job.Properties.(*tnui.TestNotificationsUiJob)
			Ω(props.Domain).Should(Equal(controlSystemDomain))

			Ω(props.NotificationsUi).ShouldNot(BeNil())
			Ω(props.NotificationsUi.Cf).ShouldNot(BeNil())
			Ω(props.NotificationsUi.Cf.AdminUser).Should(Equal("admin"))
			Ω(props.NotificationsUi.Cf.AdminPassword).Should(Equal(controlCFAdminPass))

			Ω(props.NotificationsUi.AppDomain).Should(Equal(controlSystemDomain))

			Ω(props.NotificationsUi.Uaa).ShouldNot(BeNil())
			Ω(props.NotificationsUi.Uaa.AdminClient).Should(Equal("admin"))
			Ω(props.NotificationsUi.Uaa.AdminSecret).Should(Equal(controlAdminSecret))

			Ω(props.NotificationsUi.Organization).Should(Equal("system"))
			Ω(props.NotificationsUi.Space).Should(Equal("notifications-with-ui"))
		})
	})
})
