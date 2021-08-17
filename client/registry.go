package client

import (
	"github.com/driver005/oauth/helpers"
	"github.com/ory/fosite"
)

type InternalRegistry interface {
	helpers.RegistryWriter
	Registry
}

type Registry interface {
	ClientValidator() *Validator
	ClientManager() Manager
	ClientHasher() fosite.Hasher
}
