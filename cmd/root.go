package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// This represents the base command when called without any subcommands
func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "hydra",
		Short: "Run and manage ORY Hydra",
	}
	RegisterCommandRecursive(cmd)
	return cmd
}

func RegisterCommandRecursive(parent *cobra.Command) {
	// Clients

	migrateCmd := NewMigrateCmd()
	parent.AddCommand(migrateCmd)
	migrateCmd.AddCommand(NewMigrateSqlCmd())

	serveCmd := NewServeCmd()
	parent.AddCommand(serveCmd)
	serveCmd.AddCommand(NewServeAdminCmd())
	serveCmd.AddCommand(NewServePublicCmd())
	serveCmd.AddCommand(NewServeAllCmd())

	parent.AddCommand(NewVersionCmd())
}

// Execute adds all child commands to the root command sets flags appropriately.
func Execute() {
	if err := NewRootCmd().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
