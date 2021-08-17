package consent

import (
	"net/http"

	"github.com/ory/x/errorsx"

	"github.com/gorilla/sessions"

	"github.com/ory/fosite"
	"github.com/ory/x/mapx"

	"github.com/driver005/oauth/models"
)

func sanitizeClientFromRequest(ar fosite.AuthorizeRequester) *models.Client {
	return sanitizeClient(ar.GetClient().(*models.Client))
}

func sanitizeClient(c *models.Client) *models.Client {
	cc := new(models.Client)
	// Remove the hashed secret here
	*cc = *c
	cc.Secret = ""
	return cc
}

func matchScopes(scopeStrategy fosite.ScopeStrategy, previousConsent []HandledConsentRequest, requestedScope []string) *HandledConsentRequest {
	for _, cs := range previousConsent {
		var found = true
		for _, scope := range requestedScope {
			if !scopeStrategy(cs.GrantedScope, scope) {
				found = false
				break
			}
		}

		if found {
			return &cs
		}
	}

	return nil
}

func createCsrfSession(w http.ResponseWriter, r *http.Request, store sessions.Store, name, csrf string, secure bool, sameSiteMode http.SameSite, sameSiteLegacyWorkaround bool) error {
	// Errors can be ignored here, because we always get a session session back. Error typically means that the
	// session doesn't exist yet.
	session, _ := store.Get(r, CookieName(secure, name))
	session.Values["csrf"] = csrf
	session.Options.HttpOnly = true
	session.Options.Secure = secure
	session.Options.SameSite = sameSiteMode
	if err := session.Save(r, w); err != nil {
		return errorsx.WithStack(err)
	}
	if sameSiteMode == http.SameSiteNoneMode && sameSiteLegacyWorkaround {
		return createCsrfSession(w, r, store, legacyCsrfSessionName(name), csrf, secure, 0, false)
	}
	return nil
}

func validateCsrfSession(r *http.Request, store sessions.Store, name, expectedCSRF string, sameSiteLegacyWorkaround, secure bool) error {
	if cookie, err := getCsrfSession(r, store, name, sameSiteLegacyWorkaround, secure); err != nil {
		return errorsx.WithStack(fosite.ErrRequestForbidden.WithHint("CSRF session cookie could not be decoded."))
	} else if csrf, err := mapx.GetString(cookie.Values, "csrf"); err != nil {
		return errorsx.WithStack(fosite.ErrRequestForbidden.WithHint("No CSRF value available in the session cookie."))
	} else if csrf != expectedCSRF {
		return errorsx.WithStack(fosite.ErrRequestForbidden.WithHint("The CSRF value from the token does not match the CSRF value from the data store."))
	}

	return nil
}

func getCsrfSession(r *http.Request, store sessions.Store, name string, sameSiteLegacyWorkaround, secure bool) (*sessions.Session, error) {
	cookie, err := store.Get(r, CookieName(secure, name))
	if sameSiteLegacyWorkaround && (err != nil || len(cookie.Values) == 0) {
		return store.Get(r, legacyCsrfSessionName(name))
	}
	return cookie, err
}

func legacyCsrfSessionName(name string) string {
	return name + "_legacy"
}

func CookieName(secure bool, name string) string {
	if !secure {
		return name + "_insecure"
	}
	return name
}
