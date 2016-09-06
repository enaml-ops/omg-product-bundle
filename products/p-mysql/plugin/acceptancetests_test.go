package pmysql_test

import (
	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-product-bundle/products/p-mysql/enaml-gen/acceptance-tests"
	. "github.com/enaml-ops/omg-product-bundle/products/p-mysql/plugin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("given NewAcceptanceTests func", func() {
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

		Context("when creating the acceptance-tests instancegroup", func() {
			var ig *enaml.InstanceGroup
			var jobProperties *acceptance_tests.AcceptanceTestsJob

			BeforeEach(func() {
				ig = NewAcceptanceTests(plgn)
				jobProperties = ig.GetJobByName("acceptance-tests").Properties.(*acceptance_tests.AcceptanceTestsJob)
			})

			It("then it should contain instance jobs for all required releases", func() {
				Ω(len(ig.Jobs)).Should(Equal(1))
				Ω(checkJobExists(ig.Jobs, "acceptance-tests")).Should(BeTrue(), "there should be a job for replication-canary release")
			})
			XIt("then it should setup acceptance tests job information", func() {
				Ω(jobProperties.TimeoutScale).Should(Equal(1))
				Ω(jobProperties.Cf).ShouldNot(BeNil(), "should have a properly intitialized cf")
				Ω(jobProperties.Proxy).ShouldNot(BeNil(), "should have a properly intitialized proxy")
				Ω(jobProperties.Broker).ShouldNot(BeNil(), "should have a properly intitialized broker")
				Ω(jobProperties.Service).ShouldNot(BeNil(), "should have a properly intitialized service")
			})
		})
	})
})
