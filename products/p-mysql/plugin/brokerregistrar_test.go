package pmysql_test

import (
	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-product-bundle/products/p-mysql/enaml-gen/broker-registrar"
	. "github.com/enaml-ops/omg-product-bundle/products/p-mysql/plugin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("given NewBrokerRegistrar func", func() {
	Describe("given a valid plugin object argument", func() {
		var plgn *Plugin
		BeforeEach(func() {
			plgn = new(Plugin)
		})

		Context("when the plugin sets monitoring ip values", func() {
			var ig *enaml.InstanceGroup
			BeforeEach(func() {
				ig = NewBrokerRegistrar(plgn)
			})

			It("then it should create vm instances for each given ip", func() {
				Ω(ig.Instances).Should(Equal(1), "does not create the correct number of instances")
			})

			It("then it should create a static ip block for the given ips", func() {
				Ω(len(ig.Networks)).Should(Equal(1))
			})
		})

		Context("when creating the broker-registrar instancegroup", func() {
			var ig *enaml.InstanceGroup
			var jobProperties *broker_registrar.BrokerRegistrarJob

			var controlAuthUser = "lkaslkdfhlaksdf"
			var controlAuthPass = "lkaslkdalksdklklnasdgkn"
			var controlCFPass = "kladsglkasdklgkhl"

			BeforeEach(func() {
				plgn.BrokerAuthUsername = controlAuthUser
				plgn.BrokerAuthPassword = controlAuthPass
				plgn.CFAdminPassword = controlCFPass
				ig = NewBrokerRegistrar(plgn)
				jobProperties = ig.GetJobByName("broker-registrar").Properties.(*broker_registrar.BrokerRegistrarJob)
			})

			It("then it should contain instance jobs for all required releases", func() {
				Ω(len(ig.Jobs)).Should(Equal(1))
				Ω(checkJobExists(ig.Jobs, "broker-registrar")).Should(BeTrue(), "there should be a job for replication-canary release")
			})
			It("then it should setup registrar user information", func() {
				Ω(jobProperties.Cf.AdminPassword).Should(Equal(controlCFPass))
				Ω(jobProperties.Broker.Password).Should(Equal(controlAuthPass))
				Ω(jobProperties.Broker.Username).Should(Equal(controlAuthUser))
			})
		})
	})
})
