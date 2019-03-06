package services

import (
	"crypto/rsa"

	"github.com/ricardocampos/goauth/datastore"
	log "github.com/sirupsen/logrus"
)

// NewPostgresTokenService creates an TokenService that does all the things in memory
func NewPostgresTokenService(logger *log.Logger, connectionString string, rsaPrivateKey *rsa.PrivateKey, rsaPublicKey *rsa.PublicKey) TokenService {
	if logger == nil || rsaPrivateKey == nil || rsaPublicKey == nil {
		panic("Required parameters not provided into TokenService")
	}
	logger.Debug("msg", "Running in PostgreSQL backed mode")
	clientRepo, err := datastore.NewPostgresClientRepository(connectionString, logger)
	if err != nil {
		logger.Error("msg", err)
		panic("unable to continue without a client repository")
	}
	tokenRepo, err := datastore.NewPostgresTokenRepository(connectionString, logger)
	if err != nil {
		logger.Panic("msg", err)
		panic("unable to continue without a token repository")
	}
	return &tokenService{
		clientRepository: &clientRepo,
		tokenRepository:  &tokenRepo,
		rsaPrivateKey:    rsaPrivateKey,
		rsaPublicKey:     rsaPublicKey,
		logger:           logger,
	}
}
