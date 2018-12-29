package datastore

import (
	"testing"

	"github.com/google/uuid"
	"github.com/ricardocampos/goauth/oauth2"
	"github.com/stretchr/testify/assert"
)

func getPostgresTokenRepository() oauth2.ReferenceTokenRepository {
	repository, _ := NewPostgresTokenRepository("postgres://postgres:password@localhost/goauth?sslmode=disable")
	return repository
}

func TestPostgresAddTokenWithValidToken(t *testing.T) {
	// Act
	repository := getPostgresTokenRepository()
	token := getToken()
	// Assert
	err := repository.AddToken(token)
	assert.Nil(t, err, "Should not return an error response")
}

func TestPostgresAddTokenWithInvalidToken(t *testing.T) {
	// Act
	repository := getPostgresTokenRepository()
	// Assert
	err := repository.AddToken(nil)
	assert.NotNil(t, err, "Should return an error response")
}

func TestPostgresGetTokenWithValidTokenID(t *testing.T) {
	// Act
	repository := getPostgresTokenRepository()
	expiry := oauth2.GetNowInEpochTime() + 300000
	tokenID := uuid.New().String()
	token, _ := oauth2.NewReferenceToken(tokenID, "myclient", expiry, "an_access_token")

	repository.AddToken(token)

	// Assert
	result, ok, err := repository.GetToken(tokenID)
	assert.True(t, ok, "It should return ok")
	assert.Nil(t, err, "Should not return an error response")
	assert.NotNil(t, result, "Should return reference token")
	assert.Equal(t, tokenID, result.TokenID())
}
func TestPostgresGetTokenWNotInRepository(t *testing.T) {
	// Act
	repository := getPostgresTokenRepository()
	token := getToken()
	repository.AddToken(token)

	// Assert
	result, ok, err := repository.GetToken(uuid.New().String())
	assert.False(t, ok, "It should return not ok")
	assert.NotNil(t, err, "Should return an error response")
	assert.Nil(t, result, "Should not return an item")
}

func TestPostgresGetTokenWithInvalidTokenID(t *testing.T) {
	// Act
	repository := getPostgresTokenRepository()
	token := getToken()
	repository.AddToken(token)

	// Assert
	result, ok, err := repository.GetToken("")
	assert.False(t, ok, "It should return not ok")
	assert.NotNil(t, err, "Should return an error response")
	assert.Nil(t, result, "Should not return a token")
}
