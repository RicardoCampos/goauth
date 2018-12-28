package oauth2

import (
	"errors"
	"fmt"
)

// DefaultAccessTokenLifetimeInS is 5 minutes
const DefaultAccessTokenLifetimeInS = int64(300)

// MinimumAccesTokenLifetimeInS is 2 seconds, this is to stop overly optimistic token lifetimes!
const MinimumAccesTokenLifetimeInS = int64(2)

// MaximumAllowedScopes is an arbitrary number, designed to ensure requests for ridiculous numbers of scopes are not made
const MaximumAllowedScopes = 255

// Client adds behaviour to the client structure
type Client interface {
	ValidateScopes(scopes []string) (bool, error)
	ClientID() string           // The clients unique ID
	ClientSecret() string       // The clients secret
	AccessTokenLifetime() int64 // How long in ms the token should last
	AllowedScopes() []string    // Which scopes the client is allowed to request
	TokenType() string          // Reference or Bearer
}

// client represents a client application that is registered with the OAuth2 server
type client struct {
	clientID            string   // The clients unique ID
	clientSecret        string   // The clients secret
	accessTokenLifetime int64    // How long in ms the token should last
	allowedScopes       []string // Which scopes the client is allowed to request
	tokenType           string   // Reference or Bearer
}

// NewClient creates a new client. If you provide 0 as the access token lifetime, it will use the default.
// If you supply an access token lifetime smaller than MinimumAccesTokenLifetimeInMs it will error.
func NewClient(clientID string, clientSecret string, tokenType string, accessTokenLifetime int64, allowedScopes []string) (Client, error) {
	if len(clientID) < 1 || len(clientSecret) < 1 {
		return nil, errors.New("invalid input for clientID or clientSecret when creating client")
	}
	if len(tokenType) < 1 || (tokenType != ReferenceTokenType && tokenType != BearerTokenType) {
		return nil, ErrInvalidTokenType
	}
	lifetime := accessTokenLifetime
	if lifetime == 0 {
		lifetime = DefaultAccessTokenLifetimeInS
	} else if lifetime < MinimumAccesTokenLifetimeInS {
		return nil, errors.New("invalid input for allowedScopes when creating client - there must be at least one allowed scope")
	}
	if len(allowedScopes) < 1 {
		return nil, errors.New("invalid input for allowedScopes when creating client - there must be at least one allowed scope")
	}
	if len(allowedScopes) > MaximumAllowedScopes {
		return nil, fmt.Errorf("invalid input for allowedScopes when creating client - cannot exceed %d requested scopes", MaximumAllowedScopes)
	}
	result := client{
		clientID:            clientID,
		clientSecret:        clientSecret,
		accessTokenLifetime: lifetime,
		allowedScopes:       allowedScopes,
		tokenType:           tokenType,
	}
	return result, nil
}

// ValidateScopes returns true if the client has the requested scopes in its allowedScopes list, false if not.
func (c client) ValidateScopes(scopes []string) (bool, error) {
	if scopes == nil || len(scopes) < 1 {
		return false, errors.New("at least one requested scope must be supplied")
	}
	for _, requestedScope := range scopes {
		found := false
		for _, clientScope := range c.allowedScopes {
			if clientScope == requestedScope {
				found = true
			}
		}
		if !found {
			return false, fmt.Errorf("requested %s scope not allowed for this client", requestedScope)
		}
	}
	return true, nil
}

func (c client) ClientID() string {
	return c.clientID
}

func (c client) ClientSecret() string {
	return c.clientSecret
}

func (c client) AccessTokenLifetime() int64 {
	return c.accessTokenLifetime
}

func (c client) AllowedScopes() []string {
	return c.allowedScopes
}

func (c client) TokenType() string {
	return c.tokenType
}
