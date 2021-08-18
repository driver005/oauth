package cli

import (
	"net/url"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ory/x/cmdx"
	"github.com/ory/x/flagx"
)

type Handler struct {
	Migration *MigrateHandler
}

func Remote(cmd *cobra.Command) string {
	if endpoint := flagx.MustGetString(cmd, "endpoint"); endpoint != "" {
		return strings.TrimRight(endpoint, "/")
	} else if endpoint := os.Getenv("HYDRA_URL"); endpoint != "" {
		return strings.TrimRight(endpoint, "/")
	}

	cmdx.Fatalf("To execute this command, the endpoint URL must point to the URL where ORY Hydra is located. To set the endpoint URL, use flag --endpoint or environment variable HYDRA_URL if an administrative command was used.")
	return ""
}

func RemoteURI(cmd *cobra.Command) *url.URL {
	endpoint, err := url.ParseRequestURI(Remote(cmd))
	cmdx.Must(err, "Unable to parse remote url: %s", err)
	return endpoint
}

func NewHandler() *Handler {
	return &Handler{
		Migration: newMigrateHandler(),
	}
}
