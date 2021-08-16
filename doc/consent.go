package doc

import "time"

// swagger:parameters getLoginRequest
type swaggerGetLoginRequestByChallenge struct {
	// in: query
	// required: true
	Challenge string `json:"login_challenge"`
}

// swagger:parameters getConsentRequest
type swaggerGetConsentRequestByChallenge struct {
	// in: query
	// required: true
	Challenge string `json:"consent_challenge"`
}

// swagger:parameters getLogoutRequest
type swaggerGetLogoutRequestByChallenge struct {
	// in: query
	// required: true
	Challenge string `json:"logout_challenge"`
}

// swagger:parameters revokeConsentSessions
type swaggerRevokeConsentSessions struct {
	// The subject (Subject) who's consent sessions should be deleted.
	//
	// in: query
	// required: true
	Subject string `json:"subject"`

	// If set, deletes only those consent sessions by the Subject that have been granted to the specified OAuth 2.0 Client ID
	//
	// in: query
	Client string `json:"client"`

	// If set to `?all=true`, deletes all consent sessions by the Subject that have been granted.
	//
	// in: query
	All bool `json:"all"`
}

// swagger:parameters listSubjectConsentSessions
type swaggerListSubjectConsentSessionsPayload struct {
	// in: query
	// required: true
	Subject string `json:"subject"`
}

// swagger:parameters revokeAuthenticationSession
type swaggerRevokeAuthenticationSessionPayload struct {
	// in: query
	// required: true
	Subject string `json:"subject"`
}

// swagger:parameters acceptLoginRequest
type swaggerAcceptLoginRequest struct {
	// in: query
	// required: true
	Challenge string `json:"login_challenge"`

	// in: body
	Body HandledLoginRequest
}

// swagger:parameters acceptConsentRequest
type swaggerAcceptConsentRequest struct {
	// in: query
	// required: true
	Challenge string `json:"consent_challenge"`

	// in: body
	Body HandledConsentRequest
}

// swagger:parameters acceptLogoutRequest
type swaggerAcceptLogoutRequest struct {
	// in: query
	// required: true
	Challenge string `json:"logout_challenge"`
}

// swagger:parameters rejectConsentRequest
type swaggerRejectConsentRequest struct {
	// in: query
	// required: true
	Challenge string `json:"consent_challenge"`

	// in: body
	Body RequestDeniedError
}

// swagger:parameters rejectLoginRequest
type swaggerRejectLoginRequest struct {
	// in: query
	// required: true
	Challenge string `json:"login_challenge"`

	// in: body
	Body RequestDeniedError
}

// swagger:parameters rejectLogoutRequest
type swaggerRejectLogoutRequest struct {
	// in: query
	// required: true
	Challenge string `json:"logout_challenge"`

	// in: body
	Body RequestDeniedError
}

// A list of used consent requests.
// swagger:response handledConsentRequestList
type swaggerListHandledConsentRequestsResult struct {
	// in: body
	// type: array
	Body []PreviousConsentSession
}

// swagger:model flushLoginConsentRequest
type FlushLoginConsentRequest struct {
	// NotAfter sets after which point tokens should not be flushed. This is useful when you want to keep a history
	// of recent login and consent database entries for auditing.
	NotAfter time.Time `json:"notAfter"`
}
