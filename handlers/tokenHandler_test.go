package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/ricardocampos/goauth/oauth2"
	"github.com/ricardocampos/goauth/services"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func createRequest(t *testing.T, inBody bool) *http.Request {
	if inBody {
		data := url.Values{}
		data.Set(oauth2.GrantType, oauth2.ClientCredentialsGrantType)
		data.Set(oauth2.Scope, "read")
		data.Set(oauth2.ClientID, "foo_bearer")
		data.Set(oauth2.ClientSecret, "secret")
		httpRequest, err := http.NewRequest("POST", "/connect/token", strings.NewReader(data.Encode()))
		if err != nil {
			t.Fatal(err)
		}
		httpRequest.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		return httpRequest
	}
	data := url.Values{}
	data.Set(oauth2.GrantType, oauth2.ClientCredentialsGrantType)
	data.Set(oauth2.Scope, "read")
	httpRequest, err := http.NewRequest("POST", "/token", strings.NewReader(data.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	httpRequest.Header.Set("Authorization", "Basic Zm9vX2JlYXJlcjpzZWNyZXQ=")
	httpRequest.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	return httpRequest
}

func TestDecodeTokenRequestWithCredentialsInHeader(t *testing.T) {
	// Arrange
	httpRequest := createRequest(t, false)

	// Act
	result, err := decodeTokenRequest(httpRequest)

	// Assert
	assert.Nil(t, err, "We should not have an error returned.")
	assert.NotNil(t, result, "tokenRequest should not be null.")
	assert.Equal(t, oauth2.ClientCredentialsGrantType, result.GrantType, "It should extract grant type correctly.")
	assert.Equal(t, "read", result.Scope, "It should extract scopes correctly.")
	assert.Equal(t, "foo_bearer", result.ClientID, "It should extract clientID correctly.")
	assert.Equal(t, "secret", result.ClientSecret, "It should extract clientSecret correctly.")
}

func TestDecodeTokenRequestWithCredentialsInBody(t *testing.T) {
	// Arrange
	httpRequest := createRequest(t, true)

	// Act
	result, err := decodeTokenRequest(httpRequest)

	// Assert
	assert.Nil(t, err, "We should not have an error returned.")
	assert.NotNil(t, result, "tokenRequest should not be null.")
	assert.Equal(t, oauth2.ClientCredentialsGrantType, result.GrantType, "It should extract grant type correctly.")
	assert.Equal(t, "read", result.Scope, "It should extract scopes correctly.")
	assert.Equal(t, "foo_bearer", result.ClientID, "It should extract clientID correctly.")
	assert.Equal(t, "secret", result.ClientSecret, "It should extract clientSecret correctly.")
}

func TestEncodeTokenResponse(t *testing.T) {
	recorder := httptest.NewRecorder()
	// create a token request object
	response := oauth2.TokenResponse{
		AccessToken: "abc",
		TokenType:   "reference",
		ExpiresIn:   10000,
		Scope:       "none",
	}
	encodeTokenResponse(recorder, response)
	result := recorder.Result()
	assert.Equal(t, "no-store", result.Header["Cache-Control"][0])
	assert.Equal(t, "no-cache", result.Header["Pragma"][0])
}

func TestHandleTokens(t *testing.T) {
	logger := log.New()
	recorder := httptest.NewRecorder()
	svc := services.NewInMemoryTokenService(logger, oauth2.LoadTestPrivateRsaKey("../oauth2/test_key"), oauth2.LoadTestPublicRsaKey("../oauth2/test_key.pub"))
	handler := NewTokenHandler(logger, svc)

	httpRequest := createRequest(t, false)

	handler.HandleToken(recorder, httpRequest)
	result := recorder.Result()
	assert.Equal(t, http.StatusOK, result.StatusCode)
	assert.Equal(t, "no-store", result.Header["Cache-Control"][0])
	assert.Equal(t, "no-cache", result.Header["Pragma"][0])
	assert.Equal(t, "application/json; charset=utf-8", result.Header["Content-Type"][0])
}
