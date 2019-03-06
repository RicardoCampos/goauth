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

func createValidateRequest() *http.Request {
	data := url.Values{}
	data.Set("token", "ba59e57d-bd63-459c-913d-82fff6f248f7")
	data.Set("expectedScope", "read")
	httpRequest, _ := http.NewRequest("POST", "/validate", strings.NewReader(data.Encode()))
	httpRequest.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	return httpRequest
}

func TestValidateDecodeRequestWithValidInput(t *testing.T) {
	// Arrange
	logger := log.New()
	httpRequest := createValidateRequest()

	// Act
	result, err := decodeValidateRequest(httpRequest, logger)

	// Assert
	assert.Nil(t, err)
	assert.NotNil(t, "ba59e57d-bd63-459c-913d-82fff6f248f7", result.TokenID)
}

func TestValidateDecodeRequestWithInValidInput(t *testing.T) {
	// Arrange
	logger := log.New()

	data := url.Values{}
	httpRequest, _ := http.NewRequest("POST", "/validate", strings.NewReader(data.Encode()))
	httpRequest.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Act

	_, err := decodeValidateRequest(httpRequest, logger)

	// Assert
	assert.NotNil(t, err)
}

func TestHandleValidate(t *testing.T) {
	logger := log.New()
	recorder := httptest.NewRecorder()
	svc := services.NewInMemoryTokenService(logger, oauth2.LoadTestPrivateRsaKey("../oauth2/test_key"), oauth2.LoadTestPublicRsaKey("../oauth2/test_key.pub"))
	handler := NewValidateHandler(logger, svc)

	httpRequest := createValidateRequest()

	handler.HandleValidation(recorder, httpRequest)
	result := recorder.Result()
	assert.Equal(t, http.StatusUnauthorized, result.StatusCode)
}
