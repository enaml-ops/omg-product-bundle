package cloudfoundry_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/plugin"
	"github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/plugin/pluginfakes"
)

var _ = Describe("Vault helpers", func() {
	Describe("given a RotatePasswordHash", func() {
		Context("when called with a vaultrotater and a valid hash", func() {
			var fakeVault *pluginfakes.FakeVaultRotater
			var err error
			BeforeEach(func() {
				fakeVault = new(pluginfakes.FakeVaultRotater)
				fakeVault.RotateSecretsReturns(nil)
				err = RotatePasswordHash(fakeVault, "secret/hash/of/stuff")
			})
			It("should set a valid set of secrets to vault", func() {
				_, givenSecrets := fakeVault.RotateSecretsArgsForCall(0)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(givenSecrets).ShouldNot(BeEmpty())
			})
		})
	})
	Describe("given a RotateCertHash", func() {
		Context("when called with a vaultrotater and a valid hash", func() {
			var fakeVault *pluginfakes.FakeVaultRotater
			var err error
			BeforeEach(func() {
				fakeVault = new(pluginfakes.FakeVaultRotater)
				fakeVault.RotateSecretsReturns(nil)
				err = RotateCertHash(fakeVault, "secret/hash/of/stuff", "fake.domain.io")
			})
			It("should set a valid set of secrets to vault", func() {
				_, givenSecrets := fakeVault.RotateSecretsArgsForCall(0)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(givenSecrets).ShouldNot(BeEmpty())
			})
		})
	})
	XDescribe("given a VaultRotate", func() {
		Context("when called with the rotate flag and all required vault flags", func() {
			It("then it should rotate the password values in the vault", func() {
				Ω(true).Should(BeFalse())
			})

			It("then it should rotate the keycert values in the vault", func() {
				Ω(true).Should(BeFalse())
			})
		})
	})
})
