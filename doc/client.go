package doc

import "github.com/driver005/oauth/models"

// Package client implements OAuth 2.0 client management capabilities
//
// OAuth 2.0 clients are used to perform OAuth 2.0 and OpenID Connect flows. Usually, OAuth 2.0 clients are granted
// to applications that want to use OAuth 2.0 access and refresh tokens.
//
// In ORY Hydra, OAuth 2.0 clients are used to manage ORY Hydra itself. These clients may gain highly privileged access
// if configured that way. This endpoint should be well protected and only called by code you trust.
//

// swagger:parameters createOAuth2Client
type swaggerCreateClientPayload struct {
	// in: body
	// required: true
	Body models.Client
}

// swagger:parameters updateOAuth2Client
type swaggerUpdateClientPayload struct {
	// in: path
	// required: true
	ID string `json:"id"`

	// in: body
	// required: true
	Body models.Client
}

// swagger:parameters patchOAuth2Client
type swaggerPatchClientPayload struct {
	// in: path
	// required: true
	ID string `json:"id"`

	// in: body
	// required: true
	Body patchRequest
}

// A JSONPatch request
//
// swagger:model patchRequest
type patchRequest []patchDocument

// A JSONPatch document as defined by RFC 6902
//
// swagger:model patchDocument
type patchDocument struct {
	// The operation to be performed
	//
	// required: true
	// example: "replace"
	Op string `json:"op"`

	// A JSON-pointer
	//
	// required: true
	// example: "/name"
	Path string `json:"path"`

	// The value to be used within the operations
	Value interface{} `json:"value"`

	// A JSON-pointer
	From string `json:"from"`
}

// A list of clients.
// swagger:response oAuth2ClientList
type swaggerListClientsResult struct {
	// in: body
	// type: array
	Body []models.Client
}

// swagger:parameters getOAuth2Client
type swaggerGetOAuth2Client struct {
	// The id of the OAuth 2.0 Client.
	//
	// in: path
	ID string `json:"id"`
}

// swagger:parameters deleteOAuth2Client
type swaggerDeleteOAuth2Client struct {
	// The id of the OAuth 2.0 Client.
	//
	// in: path
	ID string `json:"id"`
}
