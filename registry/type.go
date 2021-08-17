package registry

import (
	"context"
	"fmt"

	"github.com/driver005/oauth/client"
	"github.com/driver005/oauth/config"
	"github.com/driver005/oauth/consent"
	"github.com/driver005/oauth/helpers"
	"github.com/driver005/oauth/jwk"
	"github.com/driver005/oauth/oauth2"
	"github.com/driver005/oauth/persistence"
	"github.com/ory/fosite"
	foauth2 "github.com/ory/fosite/handler/oauth2"
	"github.com/ory/x/dbal"
	"github.com/ory/x/errorsx"
	"github.com/ory/x/healthx"
	"github.com/ory/x/logrusx"
	prometheus "github.com/ory/x/prometheusx"
	"github.com/pkg/errors"
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

func NewRegistryFromDSN(ctx context.Context, c *config.Provider, l *logrusx.Logger) (Registry, error) {
	driver, err := dbal.GetDriverFor(c.DSN())
	if err != nil {
		return nil, errorsx.WithStack(err)
	}

	registry, ok := driver.(Registry)
	if !ok {
		return nil, errors.Errorf("driver of type %T does not implement interface Registry", driver)
	}

	registry = registry.WithLogger(l).WithConfig(c)

	if err := registry.Init(ctx); err != nil {
		return nil, err
	}

	return registry, nil
}

func CallRegistry(ctx context.Context, r Registry) {
	fmt.Println("test 1")
	r.ClientValidator()
	fmt.Println("test 2")
	r.ClientManager()
	fmt.Println("test 3")
	r.ClientHasher()
	fmt.Println("test 4")
	r.ConsentManager()
	fmt.Println("test 5")
	r.ConsentStrategy()
	fmt.Println("test 6")
	r.SubjectIdentifierAlgorithm()
	fmt.Println("test 7")
	r.KeyManager()
	fmt.Println("test 8")
	r.KeyGenerators()
	fmt.Println("test 9")
	r.KeyCipher()
	fmt.Println("test 10")
	r.OAuth2Storage()
	fmt.Println("test 11")
	r.OAuth2Provider()
	fmt.Println("test 12")
	r.AudienceStrategy()
	fmt.Println("test 14")
	r.ScopeStrategy()
	fmt.Println("test 15")
	r.AccessTokenJWTStrategy()
	fmt.Println("test 16")
	r.OpenIDJWTStrategy()
	fmt.Println("test 17")
	r.OpenIDConnectRequestValidator()
	fmt.Println("test 18")
	r.PrometheusManager()
	fmt.Println("test 19")
	r.Tracer(ctx)
	fmt.Println("finish")
}
