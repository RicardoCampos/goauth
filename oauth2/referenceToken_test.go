package oauth2

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewReferenceTokenShouldAllowCorrectInputs(t *testing.T) {
	// Act
	expiry := GetNowInEpochTime() + 10000
	token, err := NewReferenceToken("qweqwewq", "myclient", expiry, "asdsasd")

	// Assert

	assert.NotNil(t, token, "Should return a token")
	assert.Nil(t, err, "Should not return an error response")
}

func TestNewReferenceTokenShouldNotAllowEmptyTokenID(t *testing.T) {
	// Act
	token, err := NewReferenceToken("", "client", 1234, "1231232")

	// Assert

	assert.Nil(t, token, "Should not return a token")
	assert.NotNil(t, err, "Should return an error response")
}

func TestNewReferenceTokenShouldNotAllowEmptyClientID(t *testing.T) {
	// Act
	token, err := NewReferenceToken("qweqwewq", "", 1234, "1231232")

	// Assert

	assert.Nil(t, token, "Should not return a token")
	assert.NotNil(t, err, "Should return an error response")
}

func TestNewReferenceTokenShouldNotAllowExpiryBeforeNow(t *testing.T) {
	// Act
	thePast := GetNowInEpochTime() - 1
	token, err := NewReferenceToken("qweqwewq", "myclient", thePast, "1231232")

	// Assert

	assert.Nil(t, token, "Should not return a token")
	assert.NotNil(t, err, "Should return an error response")
}

func TestNewReferenceTokenShouldNotAllowAccessToken(t *testing.T) {
	// Act
	token, err := NewReferenceToken("qweqwewq", "myclient", 1234, "")

	// Assert

	assert.Nil(t, token, "Should not return a token")
	assert.NotNil(t, err, "Should return an error response")
}
