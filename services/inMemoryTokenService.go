package services

import (
	"crypto/rsa"

	"github.com/ricardocampos/goauth/datastore"
	"github.com/ricardocampos/goauth/oauth2"
	log "github.com/sirupsen/logrus"
)

// NewInMemoryTokenService creates an TokenService that does all the things in memory
func NewInMemoryTokenService(logger *log.Logger, rsaPrivateKey *rsa.PrivateKey, rsaPublicKey *rsa.PublicKey) TokenService {
	if logger == nil || rsaPrivateKey == nil || rsaPublicKey == nil {
		panic("Required parameters not provided into TokenService")
	}
	logger.Debug("Running in in-memory mode")
	clientRepo := datastore.NewInMemoryClientRepository()
	c, _ := oauth2.NewClient("foo_bearer", "secret", oauth2.BearerTokenType, 0, []string{"read", "write"})
	c2, _ := oauth2.NewClient("foo_reference", "secret", oauth2.ReferenceTokenType, 0, []string{"read", "write"})
	clientRepo.AddClient(c)
	clientRepo.AddClient(c2)
	tokenRepo := datastore.NewInMemoryReferenceTokenRepository()

	return &tokenService{
		clientRepository: &clientRepo,
		tokenRepository:  &tokenRepo,
		rsaPrivateKey:    rsaPrivateKey,
		rsaPublicKey:     rsaPublicKey,
		logger:           logger,
	}
}
