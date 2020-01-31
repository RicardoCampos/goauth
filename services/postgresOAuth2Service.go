package services

import (
	"crypto/rsa"

	"github.com/go-kit/kit/log"
	"github.com/RicardoCampos/goauth/datastore"
)

// NewPostgresOAuth2Service creates an OAuth2Service that does all the things in memory
func NewPostgresOAuth2Service(logger log.Logger, connectionString string, rsaKey *rsa.PrivateKey) OAuth2Service {
	logger.Log("msg", "Running in PostgreSQL backed mode")
	clientRepo, err := datastore.NewPostgresClientRepository(connectionString, logger)
	if err != nil {
		logger.Log("msg", err)
		panic("unable to continue without a client repository")
	}
	tokenRepo, err := datastore.NewPostgresTokenRepository(connectionString, logger)
	if err != nil {
		logger.Log("msg", err)
		panic("unable to continue without a token repository")
	}
	return &oAuth2Service{
		clientRepository: clientRepo,
		tokenRepository:  tokenRepo,
		rsaKey:           rsaKey,
		logger:           logger,
	}
}
