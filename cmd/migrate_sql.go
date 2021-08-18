package cmd

import (
	"github.com/spf13/cobra"

	"github.com/driver005/oauth/cmd/cli"
)

// migrateSqlCmd represents the sql command
func NewMigrateSqlCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sql <database-url>",
		Short: "Create SQL schemas and apply migration plans",
		Long: `Run this command on a fresh SQL installation and when you upgrade Hydra to a new minor version. For example,
upgrading Hydra 0.7.0 to 0.8.0 requires running this command.

It is recommended to run this command close to the SQL instance (e.g. same subnet) instead of over the public internet.
This decreases risk of failure and decreases time required.

You can read in the database URL using the -e flag, for example:
	export DSN=...
	hydra migrate sql -e

### WARNING ###

Before running this command on an existing database, create a back up!
`,
		Run: cli.NewHandler().Migration.MigrateSQL,
	}

	cmd.Flags().BoolP("read-from-env", "e", false, "If set, reads the database connection string from the environment variable DSN or config file key dsn.")
	cmd.Flags().BoolP("yes", "y", false, "If set all confirmation requests are accepted without user interaction.")

	return cmd
}
