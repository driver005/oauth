package consent

import (
	"github.com/driver005/oauth/client"
	"github.com/driver005/oauth/helpers"
	"github.com/ory/fosite"
	"github.com/ory/fosite/handler/openid"
	"github.com/ory/hydra/jwk"
)

type InternalRegistry interface {
	helpers.RegistryWriter
	helpers.RegistryCookieStore
	helpers.RegistryLogger
	Registry
	client.Registry

	OAuth2Storage() helpers.FositeStorer
	OpenIDJWTStrategy() jwk.JWTStrategy
	OpenIDConnectRequestValidator() *openid.OpenIDConnectRequestValidator
	ScopeStrategy() fosite.ScopeStrategy
}

type Registry interface {
	ConsentManager() Manager
	ConsentStrategy() Strategy
	SubjectIdentifierAlgorithm() map[string]SubjectIdentifierAlgorithm
}
