package jwk

import (
	"context"
	"time"

	jose "gopkg.in/square/go-jose.v2"
)

type (
	Manager interface {
		AddKey(ctx context.Context, set string, key *jose.JSONWebKey) error

		AddKeySet(ctx context.Context, set string, keys *jose.JSONWebKeySet) error

		GetKey(ctx context.Context, set, kid string) (*jose.JSONWebKeySet, error)

		GetKeySet(ctx context.Context, set string) (*jose.JSONWebKeySet, error)

		DeleteKey(ctx context.Context, set, kid string) error

		DeleteKeySet(ctx context.Context, set string) error
	}

	SQLData struct {
		ID        int       `db:"pk"`
		Set       string    `db:"sid"`
		KID       string    `db:"kid"`
		Version   int       `db:"version"`
		CreatedAt time.Time `db:"created_at"`
		Key       string    `db:"keydata"`
	}
)

func (d SQLData) TableName() string {
	return "hydra_jwk"
}
