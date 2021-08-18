package main

import (
	"context"
	"fmt"
	"os"

	"github.com/driver005/oauth/config"
	"github.com/driver005/oauth/driver"
	"github.com/driver005/oauth/registry"
	"github.com/ory/x/configx"
	"github.com/ory/x/errorsx"
)

type MigrateHandler struct{}

func newMigrateHandler() *MigrateHandler {
	return &MigrateHandler{}
}

func (h *MigrateHandler) MigrateSQL() {
	var d registry.Registry

	d = driver.New(
		context.Background(),
		driver.WithOptions(
			configx.SkipValidation(),
			configx.WithValue(config.KeyDSN),
		),
		driver.DisableValidation(),
		driver.DisablePreloading(),
	)

	p := d.Persister()
	conn := p.Connection(context.Background())
	if conn == nil {
		fmt.Println("")
		fmt.Printf("Migrations can only be executed against a SQL-compatible driver but DSN is not a SQL source.\n")
		os.Exit(1)
		return
	}

	if err := conn.Open(); err != nil {
		fmt.Printf("Could not open the database connection:\n%+v\n", err)
		os.Exit(1)
		return
	}

	// convert migration tables
	if err := p.PrepareMigration(context.Background()); err != nil {
		fmt.Printf("Could not convert the migration table:\n%+v\n", err)
		os.Exit(1)
		return
	}

	// print migration status
	fmt.Println("The following migration is planned:")
	fmt.Println("")

	status, err := p.MigrationStatus(context.Background())
	if err != nil {
		fmt.Printf("Could not get the migration status:\n%+v\n", errorsx.WithStack(err))
		os.Exit(1)
		return
	}
	_ = status.Write(os.Stdout)

	// apply migrations
	if err := p.MigrateUp(context.Background()); err != nil {
		fmt.Printf("Could not apply migrations:\n%+v\n", errorsx.WithStack(err))
	}

	fmt.Println("Successfully applied migrations!")
}
