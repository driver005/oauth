package oauth2cors

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/driver005/oauth/client"
	"github.com/driver005/oauth/config"
	"github.com/driver005/oauth/helpers"

	"github.com/driver005/oauth/oauth2"

	"github.com/gobwas/glob"
	"github.com/rs/cors"

	"github.com/ory/fosite"
)

func Middleware(reg interface {
	Config() *config.Provider
	helpers.RegistryLogger
	oauth2.Registry
	client.Registry
}) func(h http.Handler) http.Handler {
	opts, enabled := reg.Config().CORS(config.PublicInterface)
	if !enabled {
		return func(h http.Handler) http.Handler {
			return h
		}
	}

	var alwaysAllow = len(opts.AllowedOrigins) == 0
	var patterns []glob.Glob
	for _, o := range opts.AllowedOrigins {
		if o == "*" {
			alwaysAllow = true
		}
		// if the protocol (http or https) is specified, but the url is wildcard, use special ** glob, which ignore the '.' separator.
		// This way g := glob.Compile("http://**") g.Match("http://google.com") returns true.
		if splittedO := strings.Split(o, "://"); len(splittedO) != 1 && splittedO[1] == "*" {
			o = fmt.Sprintf("%s://**", splittedO[0])
		}
		g, err := glob.Compile(strings.ToLower(o), '.')
		if err != nil {
			reg.Logger().WithError(err).Fatalf("Unable to parse cors origin: %s", o)
		}

		patterns = append(patterns, g)
	}

	options := cors.Options{
		AllowedOrigins:     opts.AllowedOrigins,
		AllowedMethods:     opts.AllowedMethods,
		AllowedHeaders:     opts.AllowedHeaders,
		ExposedHeaders:     opts.ExposedHeaders,
		MaxAge:             opts.MaxAge,
		AllowCredentials:   opts.AllowCredentials,
		OptionsPassthrough: opts.OptionsPassthrough,
		Debug:              opts.Debug,
		AllowOriginRequestFunc: func(r *http.Request, origin string) bool {
			if alwaysAllow {
				return true
			}

			origin = strings.ToLower(origin)
			for _, p := range patterns {
				if p.Match(origin) {
					return true
				}
			}

			username, _, ok := r.BasicAuth()
			if !ok || username == "" {
				token := fosite.AccessTokenFromRequest(r)
				if token == "" {
					return false
				}

				session := oauth2.NewSessionWithCustomClaims("", reg.Config().AllowedTopLevelClaims())
				_, ar, err := reg.OAuth2Provider().IntrospectToken(context.Background(), token, fosite.AccessToken, session)
				if err != nil {
					return false
				}

				username = ar.GetClient().GetID()
			}

			cl, err := reg.ClientManager().GetConcreteClient(r.Context(), username)
			if err != nil {
				return false
			}

			for _, o := range cl.AllowedCORSOrigins {
				if o == "*" {
					return true
				}
				g, err := glob.Compile(strings.ToLower(o), '.')
				if err != nil {
					return false
				}
				if g.Match(origin) {
					return true
				}
			}

			return false
		},
	}

	return cors.New(options).Handler
}
