package services

import (
	"testing"

	"github.com/ricardocampos/goauth/oauth2"
	"github.com/stretchr/testify/assert"
)

func TestValidateReturnsNoErrorsOnSuccess(t *testing.T) {
	// Arrange
	svc := NewInMemoryOAuth2Service(loadTestKey())
	tokenRequest := tokenRequest{
		clientID:     "foo_reference",
		clientSecret: "secret",
		grantType:    oauth2.ClientCredentialsGrantType,
		scope:        "read",
	}
	token, _ := svc.Token(tokenRequest)

	// Act

	request := validationRequest{
		tokenID: token.AccessToken,
	}

	// Act
	response, err := svc.Validate(request)

	// Assert
	assert.Empty(t, response.ErrMsg, "We should not have an error message")
	assert.Nil(t, response.Err, "We should not have an error returned.")
	assert.Nil(t, err, "We should not have an error returned.")
}

func TestValidateReturnsErrorsOnFailure(t *testing.T) {
	// Arrange
	svc := NewInMemoryOAuth2Service(loadTestKey())

	// Act

	request := validationRequest{
		tokenID: "this will not exist",
	}

	// Act
	response, err := svc.Validate(request)

	// Assert
	assert.NotNil(t, err)
	assert.NotNil(t, response.Err, "We should have an error returned.")
}
