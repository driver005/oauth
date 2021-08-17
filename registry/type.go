package registry

import (
	"context"

	"github.com/driver005/oauth/client"
	"github.com/driver005/oauth/config"
	"github.com/driver005/oauth/consent"
	"github.com/driver005/oauth/helpers"
	"github.com/driver005/oauth/jwk"
	"github.com/driver005/oauth/oauth2"
	"github.com/driver005/oauth/persistence"
	"github.com/ory/fosite"
	foauth2 "github.com/ory/fosite/handler/oauth2"
	"github.com/ory/x/healthx"
	"github.com/ory/x/logrusx"
	prometheus "github.com/ory/x/prometheusx"
)

type Driver interface {
	// CanHandle returns true if the driver is capable of handling the given DSN or false otherwise.
	CanHandle(dsn string) bool

	// Ping returns nil if the driver has connectivity and is healthy or an error otherwise.
	Ping() error
}

type Registry interface {
	Driver

	Init(ctx context.Context) error

	WithConfig(c *config.Provider) Registry
	WithLogger(l *logrusx.Logger) Registry

	Config() *config.Provider
	persistence.Provider
	helpers.RegistryLogger
	helpers.RegistryWriter
	helpers.RegistryCookieStore
	client.Registry
	consent.Registry
	jwk.Registry
	oauth2.Registry
	PrometheusManager() *prometheus.MetricsManager
	helpers.TracingProvider

	RegisterRoutes(admin *helpers.RouterAdmin, public *helpers.RouterPublic)
	ClientHandler() *client.Handler
	KeyHandler() *jwk.Handler
	ConsentHandler() *consent.Handler
	OAuth2Handler() *oauth2.Handler
	HealthHandler() *healthx.Handler

	OAuth2HMACStrategy() *foauth2.HMACSHAStrategy
	WithOAuth2Provider(f fosite.OAuth2Provider)
	WithConsentStrategy(c consent.Strategy)
}
