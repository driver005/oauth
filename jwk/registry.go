package jwk

import (
	"github.com/driver005/oauth/helpers"
)

type InternalRegistry interface {
	helpers.RegistryWriter
	helpers.RegistryLogger
	Registry
}

type Registry interface {
	KeyManager() Manager
	KeyGenerators() map[string]KeyGenerator
	KeyCipher() *AEAD
}
