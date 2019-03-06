package services

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/ricardocampos/goauth/oauth2"
	"github.com/stretchr/testify/assert"
)

func TestDecodeTokenRequestWithCredentialsInHeader(t *testing.T) {
	// Arrange
	data := url.Values{}
	data.Set(oauth2.GrantType, oauth2.ClientCredentialsGrantType)
	data.Set(oauth2.Scope, "read")
	httpRequest, err := http.NewRequest("POST", "/token", strings.NewReader(data.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	httpRequest.Header.Set("Authorization", "Basic bXlBd2Vzb21lQ2xpZW50OnN1cGVyc2VjcmV0cGFzc3dvcmQ=")
	httpRequest.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Act
	result, err := decodeTokenRequest(nil, httpRequest)

	// Assert
	assert.Nil(t, err, "We should not have an error returned.")
	assert.NotNil(t, result, "tokenRequest should not be null.")
	assert.Equal(t, oauth2.ClientCredentialsGrantType, result.(tokenRequest).grantType, "It should extract grant type correctly.")
	assert.Equal(t, "read", result.(tokenRequest).scope, "It should extract scopes correctly.")
	assert.Equal(t, "myAwesomeClient", result.(tokenRequest).clientID, "It should extract clientID correctly.")
	assert.Equal(t, "supersecretpassword", result.(tokenRequest).clientSecret, "It should extract clientSecret correctly.")
}

func TestDecodeTokenRequestWithCredentialsInBody(t *testing.T) {
	// Arrange
	data := url.Values{}
	data.Set(oauth2.GrantType, oauth2.ClientCredentialsGrantType)
	data.Set(oauth2.Scope, "read")
	data.Set(oauth2.ClientID, "myAwesomeClient")
	data.Set(oauth2.ClientSecret, "supersecretpassword")
	httpRequest, err := http.NewRequest("POST", "/connect/token", strings.NewReader(data.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	httpRequest.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Act
	result, err := decodeTokenRequest(nil, httpRequest)

	// Assert
	assert.Nil(t, err, "We should not have an error returned.")
	assert.NotNil(t, result, "tokenRequest should not be null.")
	assert.Equal(t, oauth2.ClientCredentialsGrantType, result.(tokenRequest).grantType, "It should extract grant type correctly.")
	assert.Equal(t, "read", result.(tokenRequest).scope, "It should extract scopes correctly.")
	assert.Equal(t, "myAwesomeClient", result.(tokenRequest).clientID, "It should extract clientID correctly.")
	assert.Equal(t, "supersecretpassword", result.(tokenRequest).clientSecret, "It should extract clientSecret correctly.")
}

func TestEncodeTokenResponse(t *testing.T) {
	recorder := httptest.NewRecorder()
	// create a token request object
	response := tokenResponse{
		AccessToken: "abc",
		TokenType:   "reference",
		ExpiresIn:   10000,
		Scope:       "none",
	}
	encodeTokenResponse(nil, recorder, response)
	result := recorder.Result()
	assert.Equal(t, "no-store", result.Header["Cache-Control"][0])
	assert.Equal(t, "no-cache", result.Header["Pragma"][0])
}
