package consent

import (
	"net/http"

	"github.com/ory/fosite"
)

var _ Strategy = new(DefaultStrategy)

type Strategy interface {
	HandleOAuth2AuthorizationRequest(w http.ResponseWriter, r *http.Request, req fosite.AuthorizeRequester) (*HandledConsentRequest, error)
	HandleOpenIDConnectLogout(w http.ResponseWriter, r *http.Request) (*LogoutResult, error)
}
