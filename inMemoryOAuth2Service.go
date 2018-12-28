package main

import (
	"github.com/ricardocampos/goauth/datastore"
	"github.com/ricardocampos/goauth/oauth2"
)

// NewInMemoryOAuth2Service creates an OAuth2Service that does all the things in memory
func NewInMemoryOAuth2Service() OAuth2Service {
	clientRepo := datastore.NewInMemoryClientRepository()
	c, _ := oauth2.NewClient("foo", "secret", oauth2.BearerTokenType, 0, []string{"read", "write"})
	clientRepo.AddClient(c)

	tokenRepo := datastore.NewInMemoryReferenceTokenRepository()

	return &oAuth2Service{
		clientRepository: clientRepo,
		tokenRepository:  tokenRepo,
	}
}
