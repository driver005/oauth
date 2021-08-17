package persistence

import (
	"context"

	"github.com/driver005/oauth/client"
	"github.com/driver005/oauth/consent"
	"github.com/driver005/oauth/helpers"
	"github.com/driver005/oauth/jwk"
	"github.com/ory/x/popx"

	"github.com/gobuffalo/pop/v5"
)

type (
	Persister interface {
		consent.Manager
		client.Manager
		helpers.FositeStorer
		jwk.Manager

		MigrationStatus(ctx context.Context) (popx.MigrationStatuses, error)
		MigrateDown(context.Context, int) error
		MigrateUp(context.Context) error
		PrepareMigration(context.Context) error
		Connection(context.Context) *pop.Connection
	}
	Provider interface {
		Persister() Persister
	}
)
