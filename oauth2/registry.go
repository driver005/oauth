package oauth2

import (
	"github.com/driver005/oauth/client"
	"github.com/driver005/oauth/consent"
	"github.com/driver005/oauth/helpers"
	"github.com/driver005/oauth/jwt"
	"github.com/ory/fosite"
	"github.com/ory/fosite/handler/openid"
)

type InternalRegistry interface {
	client.Registry
	helpers.RegistryWriter
	helpers.RegistryLogger
	consent.Registry
	Registry
}

type Registry interface {
	OAuth2Storage() helpers.FositeStorer
	OAuth2Provider() fosite.OAuth2Provider
	AudienceStrategy() fosite.AudienceMatchingStrategy
	ScopeStrategy() fosite.ScopeStrategy

	AccessTokenJWTStrategy() jwt.JWTStrategy
	OpenIDJWTStrategy() jwt.JWTStrategy

	OpenIDConnectRequestValidator() *openid.OpenIDConnectRequestValidator
}
