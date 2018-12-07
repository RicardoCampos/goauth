package oauth2client

import (
	"errors"
	"fmt"
)

// DefaultAccessTokenLifetimeInMs is 5 minutes
const DefaultAccessTokenLifetimeInMs = int64(300000)

// MinimumAccesTokenLifetimeInMs is 2 seconds, this is to stop overly optimistic token lifetimes!
const MinimumAccesTokenLifetimeInMs = int64(2000)

// MaximumAllowedScopes is an arbitrary number, designed to ensure requests for ridiculous numbers of scopes are not made
const MaximumAllowedScopes = 255

// Client adds behaviour to the client structure
type Client interface {
	ValidateScopes(scopes []string) (bool, error)
}

// client represents a client application that is registered with the OAuth2 server
type client struct {
	clientID            string   // The clients unique ID
	clientSecret        string   // The clients secret
	accessTokenLifetime int64    // How long in ms the token should last
	allowedScopes       []string // Which scopes the client is allowed to request
}

// New creates a new client. If you provide 0 as the access token lifetime, it will use the default.
// If you supply an access token lifetimw smaller than MinimumAccesTokenLifetimeInMs it will error.
func New(clientID string, clientSecret string, accessTokenLifetime int64, allowedScopes []string) (Client, error) {
	if len(clientID) < 1 || len(clientSecret) < 1 {
		return nil, errors.New("invalid input for clientID or clientSecret when creating client")
	}
	lifetime := accessTokenLifetime
	if lifetime == 0 {
		lifetime = DefaultAccessTokenLifetimeInMs
	} else if lifetime < MinimumAccesTokenLifetimeInMs {
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
