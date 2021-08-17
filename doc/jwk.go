package doc

import (
	"github.com/driver005/oauth/models"
)

// swagger:parameters getJsonWebKey deleteJsonWebKey
type swaggerJsonWebKeyQuery struct {
	// The kid of the desired key
	// in: path
	// required: true
	KID string `json:"kid"`

	// The set
	// in: path
	// required: true
	Set string `json:"set"`
}

// swagger:parameters updateJsonWebKeySet
type swaggerJwkUpdateSet struct {
	// The set
	// in: path
	// required: true
	Set string `json:"set"`

	// in: body
	Body swaggerJSONWebKeySet
}

// swagger:parameters updateJsonWebKey
type swaggerJwkUpdateSetKey struct {
	// The kid of the desired key
	// in: path
	// required: true
	KID string `json:"kid"`

	// The set
	// in: path
	// required: true
	Set string `json:"set"`

	// in: body
	Body models.JSONWebKey
}

// swagger:parameters createJsonWebKeySet
type swaggerJwkCreateSet struct {
	// The set
	// in: path
	// required: true
	Set string `json:"set"`

	// in: body
	Body models.CreateRequest
}

// swagger:parameters getJsonWebKeySet deleteJsonWebKeySet
type swaggerJwkSetQuery struct {
	// The set
	// in: path
	// required: true
	Set string `json:"set"`
}

// It is important that this model object is named JSONWebKeySet for
// "swagger generate spec" to generate only on definition of a
// JSONWebKeySet. Since one with the same name is previously defined as
// client.Client.JSONWebKeys and this one is last, this one will be
// effectively written in the swagger spec.
//
// swagger:model JSONWebKeySet
type swaggerJSONWebKeySet struct {
	// The value of the "keys" parameter is an array of JWK values.  By
	// default, the order of the JWK values within the array does not imply
	// an order of preference among them, although applications of JWK Sets
	// can choose to assign a meaning to the order for their purposes, if
	// desired.
	Keys []models.JSONWebKey `json:"keys"`
}
