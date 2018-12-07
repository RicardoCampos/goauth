package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenWillNotError(t *testing.T) {
	// Arrange
	svc := oAuth2Service{}
	var request tokenRequest
	request.grantType = ClientCredentialsGrantType
	request.scope = "read"

	// Act
	token, err := svc.Token(request)

	// Assert
	assert.Nil(t, err, "We should not have an error returned.")
	assert.NotNil(t, token, "total should not be null.")
}

func TestTokenRequiresScope(t *testing.T) {
	// Arrange
	svc := oAuth2Service{}
	var request tokenRequest
	request.grantType = ClientCredentialsGrantType

	// Act
	token, err := svc.Token(request)

	// Assert
	assert.Empty(t, token.AccessToken, "We should not be provided with a token.")
	assert.NotNil(t, err, "We should have an error returned if we do not ask for a scope.")
	assert.Equal(t, ErrInvalidScope, err)
}

func TestTokenRejectsInvalidGrant(t *testing.T) {
	// Arrange
	svc := oAuth2Service{}
	var request tokenRequest
	request.grantType = "code"
	request.scope = "read"

	// Act
	token, err := svc.Token(request)

	// Assert
	assert.Empty(t, token.AccessToken, "We should not be provided with a token.")
	assert.NotNil(t, err, "We should have an error returned if we ask for a grant type we do not support.")
	assert.Equal(t, ErrInvalidGrant, err)
}

func TestTokenRejectsEmptyGrant(t *testing.T) {
	// Arrange
	svc := oAuth2Service{}
	var request tokenRequest
	request.grantType = ""
	request.scope = "read"

	// Act
	token, err := svc.Token(request)

	// Assert
	assert.Empty(t, token.AccessToken, "We should not be provided with a token.")
	assert.NotNil(t, err, "We should have an error returned if we ask for an empty grant type.")
	assert.Equal(t, ErrInvalidGrant, err)
}
