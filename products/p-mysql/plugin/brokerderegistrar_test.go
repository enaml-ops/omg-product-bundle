package pmysql_test

import (
	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-product-bundle/products/p-mysql/enaml-gen/broker-deregistrar"
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
				ig = NewBrokerDeRegistrar(plgn)
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
			var jobProperties *broker_deregistrar.BrokerDeregistrarJob

			var controlAuthUser = "lkaslkdfhlaksdf"
			var controlAuthPass = "lkaslkdalksdklklnasdgkn"
			var controlCFPass = "kladsglkasdklgkhl"

			BeforeEach(func() {
				plgn.BrokerAuthUsername = controlAuthUser
				plgn.BrokerAuthPassword = controlAuthPass
				plgn.CFAdminPassword = controlCFPass
				ig = NewBrokerDeRegistrar(plgn)
				jobProperties = ig.GetJobByName("broker-deregistrar").Properties.(*broker_deregistrar.BrokerDeregistrarJob)
			})

			It("then it should contain instance jobs for all required releases", func() {
				Ω(len(ig.Jobs)).Should(Equal(1))
				Ω(checkJobExists(ig.Jobs, "broker-deregistrar")).Should(BeTrue(), "there should be a job for broker-deregistrar release release")
			})
			It("then it should setup deregistrar user information", func() {
				Ω(jobProperties.Cf.AdminPassword).Should(Equal(controlCFPass))
			})
		})
	})
})
