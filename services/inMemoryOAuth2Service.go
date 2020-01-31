package services

import (
	"crypto/rsa"

	"github.com/go-kit/kit/log"
	"github.com/RicardoCampos/goauth/datastore"
	"github.com/RicardoCampos/goauth/oauth2"
)

// NewInMemoryOAuth2Service creates an OAuth2Service that does all the things in memory
func NewInMemoryOAuth2Service(logger log.Logger, rsaKey *rsa.PrivateKey) OAuth2Service {
	logger.Log("msg", "Running in in-memory mode")
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
		logger:           logger,
	}
}
