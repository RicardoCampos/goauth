package main

import (
	"testing"

	"github.com/ricardocampos/goauth/oauth2"
	"github.com/stretchr/testify/assert"
)

func TestTokenChecksClientID(t *testing.T) {
	// Arrange
	svc := NewInMemoryOAuth2Service()
	request := tokenRequest{
		clientID:     "",
		clientSecret: "secret",
		grantType:    oauth2.ClientCredentialsGrantType,
		scope:        "read",
	}

	// Act
	token, err := svc.Token(request)

	// Assert
	assert.Empty(t, token.AccessToken, "We should not be provided with a token.")
	assert.NotNil(t, err, "We should have an error returned if we do not provide a known client.")
	assert.Equal(t, oauth2.ErrInvalidGrant, err)
}

func TestTokenChecksClientSecret(t *testing.T) {
	// Arrange
	svc := NewInMemoryOAuth2Service()
	request := tokenRequest{
		clientID:     "foo",
		clientSecret: "wrongpassword",
		grantType:    oauth2.ClientCredentialsGrantType,
		scope:        "read",
	}

	// Act
	token, err := svc.Token(request)

	// Assert
	assert.Empty(t, token.AccessToken, "We should not be provided with a token.")
	assert.NotNil(t, err, "We should have an error returned if we provide an incorrect secret.")
	assert.Equal(t, oauth2.ErrInvalidGrant, err)

}

func TestTokenChecksClientScopes(t *testing.T) {
	// Arrange
	svc := NewInMemoryOAuth2Service()
	request := tokenRequest{
		clientID:     "foo",
		clientSecret: "secret",
		grantType:    oauth2.ClientCredentialsGrantType,
		scope:        "thisisnotreal",
	}

	// Act
	token, err := svc.Token(request)

	// Assert
	assert.Empty(t, token.AccessToken, "We should not be provided with a token.")
	assert.NotNil(t, err, "We should have an error returned if we ask for a scope we are not allowed.")
	assert.Equal(t, oauth2.ErrInvalidScope, err)
}

func TestTokenWillNotWorkIfRepositoryNotInitialised(t *testing.T) {
	// Arrange
	svc := oAuth2Service{}
	request := tokenRequest{
		clientID:     "foo",
		clientSecret: "secret",
		grantType:    oauth2.ClientCredentialsGrantType,
		scope:        "read",
	}

	// Act
	token, err := svc.Token(request)

	// Assert
	assert.NotNil(t, err, "We should not have an error returned.")
	assert.NotNil(t, token, "token should not be null.")
	assert.Equal(t, "Cannot validate clients.", token.Err)
	assert.Equal(t, err, oauth2.ErrInvalidGrant)
}

func TestTokenWillNotErrorWhenAllInputIsOk(t *testing.T) {
	// Arrange
	svc := NewInMemoryOAuth2Service()
	request := tokenRequest{
		clientID:     "foo",
		clientSecret: "secret",
		grantType:    oauth2.ClientCredentialsGrantType,
		scope:        "read",
	}

	// Act
	token, err := svc.Token(request)

	// Assert
	assert.Nil(t, err, "We should not have an error returned.")
	assert.NotNil(t, token, "token should not be null.")
}

func TestTokenRequiresScope(t *testing.T) {
	// Arrange
	svc := NewInMemoryOAuth2Service()
	request := tokenRequest{
		clientID:     "foo",
		clientSecret: "secret",
		grantType:    oauth2.ClientCredentialsGrantType,
	}

	// Act
	token, err := svc.Token(request)

	// Assert
	assert.Empty(t, token.AccessToken, "We should not be provided with a token.")
	assert.NotNil(t, err, "We should have an error returned if we do not ask for a scope.")
	assert.Equal(t, oauth2.ErrInvalidScope, err)
}

func TestTokenRejectsInvalidGrant(t *testing.T) {
	// Arrange
	svc := NewInMemoryOAuth2Service()
	request := tokenRequest{
		clientID:     "foo",
		clientSecret: "secret",
		grantType:    "numberwang!",
		scope:        "thisisnotreal",
	}

	// Act
	token, err := svc.Token(request)

	// Assert
	assert.Empty(t, token.AccessToken, "We should not be provided with a token.")
	assert.NotNil(t, err, "We should have an error returned if we ask for a grant type we do not support.")
	assert.Equal(t, oauth2.ErrInvalidGrant, err)
}

func TestTokenRejectsEmptyGrant(t *testing.T) {
	// Arrange
	svc := NewInMemoryOAuth2Service()
	request := tokenRequest{
		clientID:     "foo",
		clientSecret: "secret",
		grantType:    "",
		scope:        "thisisnotreal",
	}

	// Act
	token, err := svc.Token(request)

	// Assert
	assert.Empty(t, token.AccessToken, "We should not be provided with a token.")
	assert.NotNil(t, err, "We should have an error returned if we ask for an empty grant type.")
	assert.Equal(t, oauth2.ErrInvalidGrant, err)
}
