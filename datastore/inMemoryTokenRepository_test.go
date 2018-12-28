package datastore

import (
	"testing"

	"github.com/ricardocampos/goauth/oauth2"
	"github.com/stretchr/testify/assert"
)

func getToken() oauth2.ReferenceToken {
	expiry := oauth2.GetNowInEpochTime() + 300000
	token, _ := oauth2.NewReferenceToken("abc123", "myclient", expiry, "an_access_token")
	return token
}

func TestAddTokenWithValidToken(t *testing.T) {
	// Act
	repository := NewInMemoryReferenceTokenRepository()
	token := getToken()
	// Assert
	err := repository.AddToken(token)
	assert.Nil(t, err, "Should not return an error response")
}

func TestAddTokenWithInvalidToken(t *testing.T) {
	// Act
	repository := NewInMemoryReferenceTokenRepository()
	// Assert
	err := repository.AddToken(nil)
	assert.NotNil(t, err, "Should return an error response")
}

func TestGetTokenWithValidTokenID(t *testing.T) {
	// Act
	repository := NewInMemoryReferenceTokenRepository()
	token := getToken()
	tokenID := token.TokenID()
	repository.AddToken(token)

	// Assert
	result, ok, err := repository.GetToken(tokenID)
	assert.True(t, ok, "It should return ok")
	assert.Nil(t, err, "Should not return an error response")
	assert.NotNil(t, result, "Should return reference token")
	assert.Equal(t, tokenID, result.TokenID())
}
func TestGetTokenWNotInRepository(t *testing.T) {
	// Act
	repository := NewInMemoryReferenceTokenRepository()
	token := getToken()
	repository.AddToken(token)

	// Assert
	result, ok, err := repository.GetToken("foobar")
	assert.False(t, ok, "It should return not ok")
	assert.Nil(t, err, "Should not return an error response")
	assert.Nil(t, result, "Should not return an item")
}

func TestGetTokenWithInvalidTokenID(t *testing.T) {
	// Act
	repository := NewInMemoryReferenceTokenRepository()
	token := getToken()
	repository.AddToken(token)

	// Assert
	result, ok, err := repository.GetToken("")
	assert.False(t, ok, "It should return not ok")
	assert.NotNil(t, err, "Should return an error response")
	assert.Nil(t, result, "Should not return a token")
}
