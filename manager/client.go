package manager

import (
	"context"

	"github.com/ory/fosite"

	"github.com/driver005/oauth/models"
)

type Client interface {
	ClientStorage

	Authenticate(ctx context.Context, id string, secret []byte) (*models.Client, error)
}

type ClientStorage interface {
	GetClient(ctx context.Context, id string) (fosite.Client, error)

	CreateClient(ctx context.Context, c *models.Client) error

	UpdateClient(ctx context.Context, c *models.Client) error

	DeleteClient(ctx context.Context, id string) error

	GetClients(ctx context.Context, filters Filter) ([]models.Client, error)

	CountClients(ctx context.Context) (int, error)

	GetConcreteClient(ctx context.Context, id string) (*models.Client, error)
}
