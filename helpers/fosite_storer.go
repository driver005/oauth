package helpers

import (
	"context"
	"time"

	"github.com/ory/fosite"
	"github.com/ory/fosite/handler/oauth2"
	"github.com/ory/fosite/handler/openid"
	"github.com/ory/fosite/handler/pkce"
)

type FositeStorer interface {
	fosite.Storage
	oauth2.CoreStorage
	openid.OpenIDConnectRequestStorage
	pkce.PKCERequestStorage

	RevokeRefreshToken(ctx context.Context, requestID string) error

	RevokeAccessToken(ctx context.Context, requestID string) error

	// flush the access token requests from the database.
	// no data will be deleted after the 'notAfter' timeframe.
	FlushInactiveAccessTokens(ctx context.Context, notAfter time.Time, limit int, batchSize int) error

	// flush the login requests from the database.
	// this will address the database long-term growth issues discussed in https://github.com/ory/hydra/issues/1574.
	// no data will be deleted after the 'notAfter' timeframe.
	FlushInactiveLoginConsentRequests(ctx context.Context, notAfter time.Time, limit int, batchSize int) error

	DeleteAccessTokens(ctx context.Context, clientID string) error

	FlushInactiveRefreshTokens(ctx context.Context, notAfter time.Time, limit int, batchSize int) error
}
