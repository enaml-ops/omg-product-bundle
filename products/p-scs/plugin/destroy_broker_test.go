package pscs_test

import (
	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-product-bundle/products/p-scs/enaml-gen/destroy-service-broker"
	pscs "github.com/enaml-ops/omg-product-bundle/products/p-scs/plugin"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("destroy broker partition", func() {
	Context("when initialized with a complete configuration", func() {
		var ig *enaml.InstanceGroup

		const (
			controlNetworkName       = "foundry-net"
			controlSystemDomain      = "sys.example.com"
			controlBrokerUser        = "broker"
			controlBrokerPass        = "brokerpass"
			controlWorkerSecret      = "workersecret"
			controlWorkerPassword    = "workerpassword"
			controlInstancesPassword = "instancespassword"
			controlDashboardSecret   = "dashbaordsecret"
			controlEncryptionKey     = "encryptionkey"
			controlCFAdminPass       = "cfadmin"
			controlUAAAdminSecret    = "uaaadminsecret"
		)

		BeforeEach(func() {
			cfg := &pscs.Config{
				Network:               controlNetworkName,
				SystemDomain:          controlSystemDomain,
				AppDomains:            []string{"apps1.example.com", "apps2.example.com"},
				BrokerUsername:        controlBrokerUser,
				BrokerPassword:        controlBrokerPass,
				SkipSSLVerify:         true,
				WorkerClientSecret:    controlWorkerSecret,
				WorkerPassword:        controlWorkerPassword,
				InstancesPassword:     controlInstancesPassword,
				BrokerDashboardSecret: controlDashboardSecret,
				EncryptionKey:         controlEncryptionKey,
				CFAdminPassword:       controlCFAdminPass,
				UAAAdminClientSecret:  controlUAAAdminSecret,
			}
			ig = pscs.NewDestroyServiceBroker(cfg)
		})

		It("should configure the instance group", func() {
			Ω(ig.Name).Should(Equal("destroy-service-broker"))
			Ω(ig.Lifecycle).Should(Equal("errand"))
			Ω(ig.Instances).Should(Equal(1))
			Ω(ig.Networks).Should(HaveLen(1))
			Ω(ig.Networks[0].Name).Should(Equal(controlNetworkName))
			Ω(ig.Networks[0].Default).Should(ConsistOf("dns", "gateway"))
			Ω(ig.Networks[0].StaticIPs).Should(BeNil())
		})

		It("should configure the destroy-service-broker job", func() {
			job := ig.GetJobByName("destroy-service-broker")
			Ω(job).ShouldNot(BeNil())

			props := job.Properties.(*destroy_service_broker.DestroyServiceBrokerJob)
			Ω(props.Domain).Should(Equal(controlSystemDomain))
			Ω(props.Ssl).ShouldNot(BeNil())
			Ω(props.Ssl.SkipCertVerify).Should(BeTrue())

			Ω(props.SpringCloudBroker).ShouldNot(BeNil())
			Ω(props.SpringCloudBroker.Broker).ShouldNot(BeNil())
			Ω(props.SpringCloudBroker.Broker.User).Should(Equal(controlBrokerUser))
			Ω(props.SpringCloudBroker.Broker.Password).Should(Equal(controlBrokerPass))

			Ω(props.SpringCloudBroker.Worker).ShouldNot(BeNil())
			Ω(props.SpringCloudBroker.Worker.ClientSecret).Should(Equal(controlWorkerSecret))

			Ω(props.SpringCloudBroker.Instances).ShouldNot(BeNil())
			Ω(props.SpringCloudBroker.Instances.InstancesUser).Should(Equal("p-spring-cloud-services"))

			Ω(props.SpringCloudBroker.Cf).ShouldNot(BeNil())
			Ω(props.SpringCloudBroker.Cf.AdminUser).Should(Equal("admin"))
			Ω(props.SpringCloudBroker.Cf.AdminPassword).Should(Equal(controlCFAdminPass))

			Ω(props.SpringCloudBroker.Uaa).ShouldNot(BeNil())
			Ω(props.SpringCloudBroker.Uaa.AdminClientId).Should(Equal("admin"))
			Ω(props.SpringCloudBroker.Uaa.AdminClientSecret).Should(Equal(controlUAAAdminSecret))
		})
	})
})
