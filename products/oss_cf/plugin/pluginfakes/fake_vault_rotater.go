// This file was generated by counterfeiter
package pluginfakes

import (
	"sync"

	"github.com/enaml-ops/omg-product-bundle/products/oss_cf/plugin"
)

type FakeVaultRotater struct {
	RotateSecretsStub        func(hash string, secrets interface{}) error
	rotateSecretsMutex       sync.RWMutex
	rotateSecretsArgsForCall []struct {
		hash    string
		secrets interface{}
	}
	rotateSecretsReturns struct {
		result1 error
	}
}

func (fake *FakeVaultRotater) RotateSecrets(hash string, secrets interface{}) error {
	fake.rotateSecretsMutex.Lock()
	fake.rotateSecretsArgsForCall = append(fake.rotateSecretsArgsForCall, struct {
		hash    string
		secrets interface{}
	}{hash, secrets})
	fake.rotateSecretsMutex.Unlock()
	if fake.RotateSecretsStub != nil {
		return fake.RotateSecretsStub(hash, secrets)
	} else {
		return fake.rotateSecretsReturns.result1
	}
}

func (fake *FakeVaultRotater) RotateSecretsCallCount() int {
	fake.rotateSecretsMutex.RLock()
	defer fake.rotateSecretsMutex.RUnlock()
	return len(fake.rotateSecretsArgsForCall)
}

func (fake *FakeVaultRotater) RotateSecretsArgsForCall(i int) (string, interface{}) {
	fake.rotateSecretsMutex.RLock()
	defer fake.rotateSecretsMutex.RUnlock()
	return fake.rotateSecretsArgsForCall[i].hash, fake.rotateSecretsArgsForCall[i].secrets
}

func (fake *FakeVaultRotater) RotateSecretsReturns(result1 error) {
	fake.RotateSecretsStub = nil
	fake.rotateSecretsReturns = struct {
		result1 error
	}{result1}
}

var _ cloudfoundry.VaultRotater = new(FakeVaultRotater)
