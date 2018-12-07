package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecodeTokenRequestWithCredentialsInHeader(t *testing.T) {
	// Arrange
	data := url.Values{}
	data.Set(GrantType, ClientCredentialsGrantType)
	data.Set(Scope, "read")
	httpRequest, err := http.NewRequest("POST", "/connect/token", strings.NewReader(data.Encode()))
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
	assert.Equal(t, ClientCredentialsGrantType, result.(tokenRequest).grantType, "It should extract grant type correctly.")
	assert.Equal(t, "read", result.(tokenRequest).scope, "It should extract scopes correctly.")
	assert.Equal(t, "myAwesomeClient", result.(tokenRequest).clientID, "It should extract clientID correctly.")
	assert.Equal(t, "supersecretpassword", result.(tokenRequest).clientSecret, "It should extract clientSecret correctly.")
}

func TestDecodeTokenRequestWithCredentialsInBody(t *testing.T) {
	// Arrange
	data := url.Values{}
	data.Set(GrantType, ClientCredentialsGrantType)
	data.Set(Scope, "read")
	data.Set(ClientID, "myAwesomeClient")
	data.Set(ClientSecret, "supersecretpassword")
	httpRequest, err := http.NewRequest("POST", "/connect/token", strings.NewReader(data.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	httpRequest.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Act
	result, err := decodeTokenRequest(nil, httpRequest)
	fmt.Println(result)

	// Assert
	assert.Nil(t, err, "We should not have an error returned.")
	assert.NotNil(t, result, "tokenRequest should not be null.")
	assert.Equal(t, ClientCredentialsGrantType, result.(tokenRequest).grantType, "It should extract grant type correctly.")
	assert.Equal(t, "read", result.(tokenRequest).scope, "It should extract scopes correctly.")
	assert.Equal(t, "myAwesomeClient", result.(tokenRequest).clientID, "It should extract clientID correctly.")
	assert.Equal(t, "supersecretpassword", result.(tokenRequest).clientSecret, "It should extract clientSecret correctly.")
}
