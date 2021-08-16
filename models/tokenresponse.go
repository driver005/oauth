package models

// The token response
// swagger:model oauthTokenResponse
type swaggerOAuthTokenResponse struct {
	// The lifetime in seconds of the access token.  For
	//  example, the value "3600" denotes that the access token will
	// expire in one hour from the time the response was generated.
	ExpiresIn int `json:"expires_in"`

	// The scope of the access token
	Scope int `json:"scope"`

	// To retrieve a refresh token request the id_token scope.
	IDToken int `json:"id_token"`

	// The access token issued by the authorization server.
	AccessToken string `json:"access_token"`

	// The refresh token, which can be used to obtain new
	// access tokens. To retrieve it add the scope "offline" to your access token request.
	RefreshToken string `json:"refresh_token"`

	// The type of the token issued
	TokenType string `json:"token_type"`
}
