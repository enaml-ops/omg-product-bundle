package cloudfoundry_test

import (
	"fmt"

	"github.com/enaml-ops/enaml"
	das "github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/deploy-autoscaling"
	db "github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/destroy-broker"
	rb "github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/register-broker"
	ta "github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/test-autoscaling"

	. "github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/plugin"
	"github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/plugin/config"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("autoscaling", func() {

	Describe("deploy autoscaling errand", func() {
		Context("when initialized with a complete set of arguments", func() {
			var (
				igc InstanceGroupCreator
				ig  *enaml.InstanceGroup
				dm  *enaml.DeploymentManifest
			)
			const (
				controlProxyIP           = "10.1.2.3"
				controlVMType            = "small"
				controlBrokerUser        = "brokeruser"
				controlBrokerPassword    = "brokerpassword"
				controlAdminPassword     = "adminPassword"
				controlAutoscaleDBUser   = "autoscale"
				controlAutoscaleDBPass   = "autopass"
				controlAutoscalingSecret = "autoscalesecret"
			)
			BeforeEach(func() {
				c := &config.Config{
					AZs:          []string{"z1"},
					StemcellName: "cool-ubuntu-animal",
					NetworkName:  "foundry-net",
					AppDomains:   []string{"apps.example.com"},
					SystemDomain: "sys.example.com",
					Secret: config.Secret{
						AutoscaleBrokerPassword:        controlBrokerPassword,
						AdminPassword:                  controlAdminPassword,
						AutoScalingServiceClientSecret: controlAutoscalingSecret,
						AutoscaleDBPassword:            controlAutoscaleDBPass,
					},
					Certs:         &config.Certs{},
					InstanceCount: config.InstanceCount{},
					IP: config.IP{
						MySQLProxyIPs: []string{controlProxyIP},
					},
					VMType: config.VMType{
						ErrandVMType: controlVMType,
					},
					User: config.User{
						AutoscaleDBUser:     controlAutoscaleDBUser,
						AutoscaleBrokerUser: controlBrokerUser,
					},
					SkipSSLCertVerify: true,
				}
				igc = NewDeployAutoscaling(c)
				dm = new(enaml.DeploymentManifest)
				ig = igc.ToInstanceGroup()
				dm.AddInstanceGroup(ig)
				Ω(dm.GetInstanceGroupByName("autoscaling")).ShouldNot(BeNil())
			})

			It("should configure the instance group", func() {
				Ω(ig).ShouldNot(BeNil())
				Ω(ig.Name).Should(Equal("autoscaling"))
				Ω(ig.Instances).Should(Equal(1))
				Ω(ig.VMType).Should(Equal(controlVMType))
				Ω(ig.AZs).Should(ConsistOf("z1"))
				Ω(ig.Stemcell).Should(Equal("cool-ubuntu-animal"))
				Ω(ig.Networks).Should(HaveLen(1))
				Ω(ig.Networks[0].Name).Should(Equal("foundry-net"))
				Ω(ig.Update.MaxInFlight).Should(Equal(1))
			})

			It("should configure the deploy-autoscaling job", func() {
				Ω(ig.Jobs).Should(HaveLen(1))
				Ω(ig.Jobs[0].Name).Should(Equal("deploy-autoscaling"))
				Ω(ig.Jobs[0].Release).Should(Equal(CFAutoscalingReleaseName))
				props := ig.Jobs[0].Properties
				Ω(props).ShouldNot(BeNil())

				job := props.(*das.DeployAutoscalingJob)

				Ω(job.AppDomains).Should(ConsistOf("sys.example.com"))

				Ω(job.Autoscale).ShouldNot(BeNil())
				Ω(job.Autoscale.Broker).ShouldNot(BeNil())
				Ω(job.Autoscale.Broker.User).Should(Equal(controlBrokerUser))
				Ω(job.Autoscale.Broker.Password).Should(Equal(controlBrokerPassword))

				Ω(job.Autoscale.Cf).ShouldNot(BeNil())
				Ω(job.Autoscale.Cf.AdminUser).Should(Equal("admin"))
				Ω(job.Autoscale.Cf.AdminPassword).Should(Equal(controlAdminPassword))

				Ω(job.Autoscale.InstanceCount).Should(Equal(1))

				Ω(job.Autoscale.Database).ShouldNot(BeNil())
				Ω(job.Autoscale.Database.Url).Should(Equal(fmt.Sprintf("mysql://%s:%s@%s:3306/autoscale", controlAutoscaleDBUser, controlAutoscaleDBPass, controlProxyIP)))

				Ω(job.Autoscale.EncryptionKey).ShouldNot(BeEmpty())
				Ω(job.Autoscale.EnableDiego).Should(BeTrue())
				Ω(job.Autoscale.NotificationsHost).Should(Equal("https://notifications.sys.example.com"))
				Ω(job.Autoscale.Organization).Should(Equal("system"))
				Ω(job.Autoscale.Space).Should(Equal("autoscaling"))

				Ω(job.Autoscale.MarketplaceCompanyName).Should(Equal("Pivotal"))
				Ω(job.Autoscale.MarketplaceImageUrl).ShouldNot(BeEmpty())
				Ω(job.Autoscale.MarketplaceDocumentationUrl).Should(Equal("http://docs.gopivotal.com/pivotalcf/"))

				Ω(job.Domain).Should(Equal("sys.example.com"))

				Ω(job.Ssl).ShouldNot(BeNil())
				Ω(job.Ssl.SkipCertVerify).Should(BeTrue())

				Ω(job.Uaa).ShouldNot(BeNil())
				Ω(job.Uaa.Clients).ShouldNot(BeNil())
				Ω(job.Uaa.Clients.AutoscalingService).ShouldNot(BeNil())
				Ω(job.Uaa.Clients.AutoscalingService.Secret).Should(Equal(controlAutoscalingSecret))

			})
		})
	})

	Describe("register broker errand", func() {
		Context("when initialized with a complete set of arguments", func() {
			var (
				igc InstanceGroupCreator
				ig  *enaml.InstanceGroup
				dm  *enaml.DeploymentManifest
			)
			const (
				controlProxyIP           = "10.1.2.3"
				controlVMType            = "small"
				controlBrokerUser        = "brokeruser"
				controlBrokerPassword    = "brokerpassword"
				controlAdminPassword     = "adminPassword"
				controlAutoscaleDBUser   = "autoscale"
				controlAutoscaleDBPass   = "autopass"
				controlAutoscalingSecret = "autoscalesecret"
			)
			BeforeEach(func() {
				c := &config.Config{
					AZs:          []string{"z1"},
					StemcellName: "cool-ubuntu-animal",
					NetworkName:  "foundry-net",
					SystemDomain: "sys.example.com",
					Secret: config.Secret{
						AutoscaleBrokerPassword:        controlBrokerPassword,
						AdminPassword:                  controlAdminPassword,
						AutoScalingServiceClientSecret: controlAutoscalingSecret,
						AutoscaleDBPassword:            controlAutoscaleDBPass,
					},
					Certs:         &config.Certs{},
					InstanceCount: config.InstanceCount{},
					VMType: config.VMType{
						ErrandVMType: controlVMType,
					},
					User: config.User{
						AutoscaleDBUser:     controlAutoscaleDBUser,
						AutoscaleBrokerUser: controlBrokerUser,
					},
					SkipSSLCertVerify: true,
				}
				igc = NewAutoscaleRegisterBroker(c)
				dm = new(enaml.DeploymentManifest)
				ig = igc.ToInstanceGroup()
				dm.AddInstanceGroup(ig)
				Ω(dm.GetInstanceGroupByName("autoscaling-register-broker")).ShouldNot(BeNil())
			})

			It("should configure the instance group", func() {
				Ω(ig).ShouldNot(BeNil())
				Ω(ig.Name).Should(Equal("autoscaling-register-broker"))
				Ω(ig.Instances).Should(Equal(1))
				Ω(ig.VMType).Should(Equal(controlVMType))
				Ω(ig.AZs).Should(ConsistOf("z1"))
				Ω(ig.Stemcell).Should(Equal("cool-ubuntu-animal"))
				Ω(ig.Networks).Should(HaveLen(1))
				Ω(ig.Networks[0].Name).Should(Equal("foundry-net"))
				Ω(ig.Update.MaxInFlight).Should(Equal(1))
			})

			It("should configure the register-broker job", func() {
				Ω(ig.Jobs).Should(HaveLen(1))
				Ω(ig.Jobs[0].Name).Should(Equal("register-broker"))
				Ω(ig.Jobs[0].Release).Should(Equal(CFAutoscalingReleaseName))
				props := ig.Jobs[0].Properties
				Ω(props).ShouldNot(BeNil())

				job := props.(*rb.RegisterBrokerJob)
				Ω(job.AppDomains).Should(ConsistOf("sys.example.com"))
				Ω(job.Autoscale).ShouldNot(BeNil())
				Ω(job.Autoscale.Broker).ShouldNot(BeNil())
				Ω(job.Autoscale.Broker.User).Should(Equal(controlBrokerUser))
				Ω(job.Autoscale.Broker.Password).Should(Equal(controlBrokerPassword))

				Ω(job.Autoscale.Cf).ShouldNot(BeNil())
				Ω(job.Autoscale.Cf.AdminUser).Should(Equal("admin"))
				Ω(job.Autoscale.Cf.AdminPassword).Should(Equal(controlAdminPassword))

				Ω(job.Autoscale.Organization).Should(Equal("system"))
				Ω(job.Autoscale.Space).Should(Equal("autoscaling"))
				Ω(job.Domain).Should(Equal("sys.example.com"))
				Ω(job.Ssl).ShouldNot(BeNil())
				Ω(job.Ssl.SkipCertVerify).Should(BeTrue())
			})
		})
	})

	Describe("destroy broker errand", func() {
		Context("when initialized with a complete set of arguments", func() {
			var (
				igc InstanceGroupCreator
				ig  *enaml.InstanceGroup
				dm  *enaml.DeploymentManifest
			)
			const (
				controlProxyIP           = "10.1.2.3"
				controlVMType            = "small"
				controlBrokerUser        = "brokeruser"
				controlBrokerPassword    = "brokerpassword"
				controlAdminPassword     = "adminPassword"
				controlAutoscaleDBUser   = "autoscale"
				controlAutoscaleDBPass   = "autopass"
				controlAutoscalingSecret = "autoscalesecret"
			)
			BeforeEach(func() {
				c := &config.Config{
					AZs:          []string{"z1"},
					StemcellName: "cool-ubuntu-animal",
					NetworkName:  "foundry-net",
					SystemDomain: "sys.example.com",
					Secret: config.Secret{
						AutoscaleBrokerPassword:        controlBrokerPassword,
						AdminPassword:                  controlAdminPassword,
						AutoScalingServiceClientSecret: controlAutoscalingSecret,
						AutoscaleDBPassword:            controlAutoscaleDBPass,
					},
					Certs:         &config.Certs{},
					InstanceCount: config.InstanceCount{},
					VMType: config.VMType{
						ErrandVMType: controlVMType,
					},
					User: config.User{
						AutoscaleDBUser:     controlAutoscaleDBUser,
						AutoscaleBrokerUser: controlBrokerUser,
					},
					SkipSSLCertVerify: true,
				}
				igc = NewAutoscaleDestroyBroker(c)
				dm = new(enaml.DeploymentManifest)
				ig = igc.ToInstanceGroup()
				dm.AddInstanceGroup(ig)
				Ω(dm.GetInstanceGroupByName("autoscaling-destroy-broker")).ShouldNot(BeNil())
			})

			It("should configure the instance group", func() {
				Ω(ig).ShouldNot(BeNil())
				Ω(ig.Name).Should(Equal("autoscaling-destroy-broker"))
				Ω(ig.Instances).Should(Equal(1))
				Ω(ig.VMType).Should(Equal(controlVMType))
				Ω(ig.AZs).Should(ConsistOf("z1"))
				Ω(ig.Stemcell).Should(Equal("cool-ubuntu-animal"))
				Ω(ig.Networks).Should(HaveLen(1))
				Ω(ig.Networks[0].Name).Should(Equal("foundry-net"))
				Ω(ig.Update.MaxInFlight).Should(Equal(1))
			})

			It("should configure the destroy-broker job", func() {
				Ω(ig.Jobs).Should(HaveLen(1))
				Ω(ig.Jobs[0].Name).Should(Equal("destroy-broker"))
				Ω(ig.Jobs[0].Release).Should(Equal(CFAutoscalingReleaseName))
				props := ig.Jobs[0].Properties
				Ω(props).ShouldNot(BeNil())

				job := props.(*db.DestroyBrokerJob)
				Ω(job.Autoscale).ShouldNot(BeNil())
				Ω(job.Autoscale.Broker).ShouldNot(BeNil())
				Ω(job.Autoscale.Broker.User).Should(Equal(controlBrokerUser))
				Ω(job.Autoscale.Broker.Password).Should(Equal(controlBrokerPassword))

				Ω(job.Autoscale.Cf).ShouldNot(BeNil())
				Ω(job.Autoscale.Cf.AdminUser).Should(Equal("admin"))
				Ω(job.Autoscale.Cf.AdminPassword).Should(Equal(controlAdminPassword))

				Ω(job.Autoscale.Organization).Should(Equal("system"))
				Ω(job.Autoscale.Space).Should(Equal("autoscaling"))
				Ω(job.Domain).Should(Equal("sys.example.com"))
				Ω(job.Ssl).ShouldNot(BeNil())
				Ω(job.Ssl.SkipCertVerify).Should(BeTrue())
			})
		})
	})

	Describe("autoscaling tests errand", func() {
		Context("when initialized with a complete set of arguments", func() {
			var (
				igc InstanceGroupCreator
				ig  *enaml.InstanceGroup
				dm  *enaml.DeploymentManifest
			)
			const (
				controlProxyIP           = "10.1.2.3"
				controlVMType            = "small"
				controlBrokerUser        = "brokeruser"
				controlBrokerPassword    = "brokerpassword"
				controlAdminPassword     = "adminPassword"
				controlAutoscaleDBUser   = "autoscale"
				controlAutoscaleDBPass   = "autopass"
				controlAutoscalingSecret = "autoscalesecret"
			)
			BeforeEach(func() {
				c := &config.Config{
					AZs:          []string{"z1"},
					StemcellName: "cool-ubuntu-animal",
					NetworkName:  "foundry-net",
					SystemDomain: "sys.example.com",
					Secret: config.Secret{
						AutoscaleBrokerPassword:        controlBrokerPassword,
						AdminPassword:                  controlAdminPassword,
						AutoScalingServiceClientSecret: controlAutoscalingSecret,
						AutoscaleDBPassword:            controlAutoscaleDBPass,
					},
					Certs:         &config.Certs{},
					InstanceCount: config.InstanceCount{},
					VMType: config.VMType{
						ErrandVMType: controlVMType,
					},
					User: config.User{
						AutoscaleDBUser:     controlAutoscaleDBUser,
						AutoscaleBrokerUser: controlBrokerUser,
					},
					SkipSSLCertVerify: true,
				}
				igc = NewAutoscalingTests(c)
				dm = new(enaml.DeploymentManifest)
				ig = igc.ToInstanceGroup()
				dm.AddInstanceGroup(ig)
				Ω(dm.GetInstanceGroupByName("autoscaling-tests")).ShouldNot(BeNil())
			})

			It("should configure the instance group", func() {
				Ω(ig).ShouldNot(BeNil())
				Ω(ig.Name).Should(Equal("autoscaling-tests"))
				Ω(ig.Instances).Should(Equal(1))
				Ω(ig.VMType).Should(Equal(controlVMType))
				Ω(ig.AZs).Should(ConsistOf("z1"))
				Ω(ig.Stemcell).Should(Equal("cool-ubuntu-animal"))
				Ω(ig.Networks).Should(HaveLen(1))
				Ω(ig.Networks[0].Name).Should(Equal("foundry-net"))
				Ω(ig.Update.MaxInFlight).Should(Equal(1))
			})

			It("should configure the test-autoscaling job", func() {
				Ω(ig.Jobs).Should(HaveLen(1))
				Ω(ig.Jobs[0].Name).Should(Equal("test-autoscaling"))
				Ω(ig.Jobs[0].Release).Should(Equal(CFAutoscalingReleaseName))
				props := ig.Jobs[0].Properties
				Ω(props).ShouldNot(BeNil())

				job := props.(*ta.TestAutoscalingJob)
				Ω(job.Autoscale).ShouldNot(BeNil())
				Ω(job.Autoscale.Cf).ShouldNot(BeNil())
				Ω(job.Autoscale.Cf.AdminUser).Should(Equal("admin"))
				Ω(job.Autoscale.Cf.AdminPassword).Should(Equal(controlAdminPassword))

				Ω(job.Domain).Should(Equal("sys.example.com"))
				Ω(job.Ssl).ShouldNot(BeNil())
				Ω(job.Ssl.SkipCertVerify).Should(BeTrue())
			})
		})
	})

})
