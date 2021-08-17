package consent

import (
	"context"
	"time"

	"github.com/driver005/oauth/models"
)

type ForcedObfuscatedLoginSession struct {
	ClientID          string `db:"client_id"`
	Subject           string `db:"subject"`
	SubjectObfuscated string `db:"subject_obfuscated"`
}

func (ForcedObfuscatedLoginSession) TableName() string {
	return "hydra_oauth2_obfuscated_authentication_session"
}

type Manager interface {
	CreateConsentRequest(ctx context.Context, req *ConsentRequest) error
	GetConsentRequest(ctx context.Context, challenge string) (*ConsentRequest, error)
	HandleConsentRequest(ctx context.Context, challenge string, r *HandledConsentRequest) (*ConsentRequest, error)
	RevokeSubjectConsentSession(ctx context.Context, user string) error
	RevokeSubjectClientConsentSession(ctx context.Context, user, client string) error

	VerifyAndInvalidateConsentRequest(ctx context.Context, verifier string) (*HandledConsentRequest, error)
	FindGrantedAndRememberedConsentRequests(ctx context.Context, client, user string) ([]HandledConsentRequest, error)
	FindSubjectsGrantedConsentRequests(ctx context.Context, user string, limit, offset int) ([]HandledConsentRequest, error)
	CountSubjectsGrantedConsentRequests(ctx context.Context, user string) (int, error)

	// Cookie management
	GetRememberedLoginSession(ctx context.Context, id string) (*LoginSession, error)
	CreateLoginSession(ctx context.Context, session *LoginSession) error
	DeleteLoginSession(ctx context.Context, id string) error
	RevokeSubjectLoginSession(ctx context.Context, user string) error
	ConfirmLoginSession(ctx context.Context, id string, authTime time.Time, subject string, remember bool) error

	CreateLoginRequest(ctx context.Context, req *LoginRequest) error
	GetLoginRequest(ctx context.Context, challenge string) (*LoginRequest, error)
	HandleLoginRequest(ctx context.Context, challenge string, r *HandledLoginRequest) (*LoginRequest, error)
	VerifyAndInvalidateLoginRequest(ctx context.Context, verifier string) (*HandledLoginRequest, error)

	CreateForcedObfuscatedLoginSession(ctx context.Context, session *ForcedObfuscatedLoginSession) error
	GetForcedObfuscatedLoginSession(ctx context.Context, client, obfuscated string) (*ForcedObfuscatedLoginSession, error)

	ListUserAuthenticatedClientsWithFrontChannelLogout(ctx context.Context, subject, sid string) ([]models.Client, error)
	ListUserAuthenticatedClientsWithBackChannelLogout(ctx context.Context, subject, sid string) ([]models.Client, error)

	CreateLogoutRequest(ctx context.Context, request *LogoutRequest) error
	GetLogoutRequest(ctx context.Context, challenge string) (*LogoutRequest, error)
	AcceptLogoutRequest(ctx context.Context, challenge string) (*LogoutRequest, error)
	RejectLogoutRequest(ctx context.Context, challenge string) error
	VerifyAndInvalidateLogoutRequest(ctx context.Context, verifier string) (*LogoutRequest, error)
}
