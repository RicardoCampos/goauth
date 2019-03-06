package oauth2

import "errors"

// Errors

// ErrInvalidGrant returned when the input value for grant type is wrong
var ErrInvalidGrant = errors.New("Invalid Grant Type")

// ErrInvalidScope returned when the input value for scope is wrong
var ErrInvalidScope = errors.New("Invalid Scope Type")

// ErrInvalidTokenType returned when the input value for token is wrong
var ErrInvalidTokenType = errors.New("Invalid Token Type")

// ErrInvalidRequest Invalid input
var ErrInvalidRequest = errors.New("Invalid Request")

// ErrInvalidToken Token not found or invalid
var ErrInvalidToken = errors.New("not found or expired")

// ErrInsufficientScope returned when a token hasn't got the scope that has been requested
var ErrInsufficientScope = errors.New("insufficient_scope")

// Return value constants

// BearerTokenType will return the full JWT in an AccessToken field.
const BearerTokenType = "Bearer"

// ReferenceTokenType will return a UUID in the AccessToken field. This will need to be passed back to this service for validation.
const ReferenceTokenType = "Reference"

// Oauth2

// ClientCredentialsGrantType is the string representing the client credentials grant type
const ClientCredentialsGrantType = "client_credentials"

// Request input constants

// GrantType constant for form input
const GrantType = "grant_type"

// Scope constant for form input
const Scope = "scope"

// ClientID constant for form input
const ClientID = "client_id"

// ClientSecret constant for form input
const ClientSecret = "client_secret"
