package pmysql_test

import (
	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-product-bundle/products/p-mysql/enaml-gen/rejoin-unsafe"
	. "github.com/enaml-ops/omg-product-bundle/products/p-mysql/plugin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("given NewRejoinUnsafe func", func() {
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

		Context("when creating the rejoin-unsafe instancegroup", func() {
			var ig *enaml.InstanceGroup
			var jobProperties *rejoin_unsafe.RejoinUnsafeJob

			var controlUser = "lkaslkdfhlaksdf"
			var controlPass = "lkaslkdalksdklklnasdgkn"
			var controlIPs = []string{"kladsglkasdklgkhl"}

			BeforeEach(func() {
				plgn.GaleraHealthcheckUsername = controlUser
				plgn.GaleraHealthcheckPassword = controlPass
				plgn.IPs = controlIPs
				ig = NewRejoinUnsafe(plgn)
				jobProperties = ig.GetJobByName("rejoin-unsafe").Properties.(*rejoin_unsafe.RejoinUnsafeJob)
			})

			It("then it should contain instance jobs for all required releases", func() {
				Ω(len(ig.Jobs)).Should(Equal(1))
				Ω(checkJobExists(ig.Jobs, "rejoin-unsafe")).Should(BeTrue(), "there should be a job for replication-canary release")
			})
			It("then it should setup rejoin unsafe job information", func() {
				Ω(jobProperties.ClusterIps).Should(Equal(controlIPs), "the ips for the mysql cluser should be set")
				Ω(jobProperties.CfMysql.Mysql.GaleraHealthcheck.EndpointPassword).Should(Equal(controlPass), "the password for the galera healthcheck should be set")
				Ω(jobProperties.CfMysql.Mysql.GaleraHealthcheck.EndpointUsername).Should(Equal(controlUser), "the username for the galera healthcheck should be set")
			})
		})
	})
})
