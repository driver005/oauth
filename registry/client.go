package registry

import (
	"github.com/ory/fosite"
	"github.com/ory/hydra/x"

	"github.com/driver005/oauth/client"

	"github.com/driver005/oauth/manager"
)

type InternalClient interface {
	x.RegistryWriter
	Client
}

type Client interface {
	ClientValidator() *client.Validator
	ClientManager() manager.Client
	ClientHasher() fosite.Hasher
}
