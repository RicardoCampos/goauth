package services

import (
	"crypto/rsa"
	"io/ioutil"
	"testing"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-kit/kit/log"
	"github.com/google/uuid"
	"github.com/ricardocampos/goauth/oauth2"
	"github.com/stretchr/testify/assert"
)

func loadTestPublicKey() *rsa.PublicKey {
	k, err := ioutil.ReadFile("test_key.pub")
	if err != nil {
		panic(err)
	}
	key, err := jwt.ParseRSAPublicKeyFromPEM(k)
	if err != nil {
		panic(err)
	}
	return key
}

func TestValidateReturnsNoErrorsOnSuccess(t *testing.T) {
	// Arrange
	svc := NewInMemoryOAuth2Service(log.NewNopLogger(), loadTestKey())
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
	svc := NewInMemoryOAuth2Service(log.NewNopLogger(), loadTestKey())

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

func TestValidateTokenChecksExpiryInPast(t *testing.T) {
	// Arrange
	expiry := oauth2.GetNowInEpochTime() + 1
	referenceToken, err := oauth2.NewReferenceToken(uuid.New().String(), "client", expiry, "asdasdsad")
	if err != nil {
		t.Fatal()
	}

	// Act
	result := validateToken(referenceToken, loadTestPublicKey())

	// Assert
	assert.False(t, result)
}

func TestValidateTokenChecksExpiryInFuture(t *testing.T) {
	// Arrange
	now := oauth2.GetNowInEpochTime()
	expiry := oauth2.GetNowInEpochTime() + 10000
	jti := uuid.New().String()
	claims := jwt.StandardClaims{
		Audience:  "youlot", // TODO: Really, this is optional, but will come from the client struct
		ExpiresAt: expiry,
		Id:        jti,
		IssuedAt:  now,
		Issuer:    "goauth.xom", // TODO: This is a constant and should be the URI of the issuing server.
		NotBefore: now,          // this should be the current time
		Subject:   "client",
	}
	token, _ := signToken(claims, loadTestKey())
	referenceToken, err := oauth2.NewReferenceToken(jti, "client", expiry, token)
	if err != nil {
		t.Fatal()
	}

	// Act
	result := validateToken(referenceToken, loadTestPublicKey())

	// Assert
	assert.True(t, result)
}

func TestValidateTokenChecksInput(t *testing.T) {
	// Act
	result := validateToken(nil, loadTestPublicKey())

	// Assert
	assert.False(t, result)
}
