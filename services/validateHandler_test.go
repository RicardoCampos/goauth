package services

import (
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateDecodeRequestWithValidInput(t *testing.T) {
	// Arrange
	data := url.Values{}
	data.Set("token", "ba59e57d-bd63-459c-913d-82fff6f248f7")

	httpRequest, _ := http.NewRequest("POST", "/validate", strings.NewReader(data.Encode()))
	httpRequest.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Act
	result, err := decodeValidateRequest(nil, httpRequest)

	// Assert
	assert.Nil(t, err)
	assert.NotNil(t, "ba59e57d-bd63-459c-913d-82fff6f248f7", result.(validationRequest).tokenID)
}

func TestValidateDecodeRequestWithInValidInput(t *testing.T) {
	// Arrange
	data := url.Values{}

	httpRequest, _ := http.NewRequest("POST", "/validate", strings.NewReader(data.Encode()))
	httpRequest.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Act
	result, err := decodeValidateRequest(nil, httpRequest)

	// Assert
	assert.NotNil(t, err)
	assert.Nil(t, result)
}
