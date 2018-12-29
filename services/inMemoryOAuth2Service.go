package services

import (
	"crypto/rsa"

	"github.com/ricardocampos/goauth/datastore"
	"github.com/ricardocampos/goauth/oauth2"
)

// NewInMemoryOAuth2Service creates an OAuth2Service that does all the things in memory
func NewInMemoryOAuth2Service(rsaKey *rsa.PrivateKey) OAuth2Service {
	clientRepo := datastore.NewInMemoryClientRepository()
	c, _ := oauth2.NewClient("foo", "secret", oauth2.BearerTokenType, 0, []string{"read", "write"})
	c2, _ := oauth2.NewClient("foo_reference", "secret", oauth2.ReferenceTokenType, 0, []string{"read", "write"})
	clientRepo.AddClient(c)
	clientRepo.AddClient(c2)
	tokenRepo := datastore.NewInMemoryReferenceTokenRepository()

	return &oAuth2Service{
		clientRepository: clientRepo,
		tokenRepository:  tokenRepo,
		rsaKey:           rsaKey,
	}
}
