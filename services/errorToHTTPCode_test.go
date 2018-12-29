package services

import (
	"errors"
	"net/http"
	"testing"

	"github.com/ricardocampos/goauth/oauth2"
	"github.com/stretchr/testify/assert"
)

func TestErrorToHTTPCodeReturnsUnauthorizedForInvalidGrants(t *testing.T) {
	status := errorToHTTPCode(oauth2.ErrInvalidGrant)
	assert.Equal(t, http.StatusUnauthorized, status)
}

func TestErrorToHTTPCodeReturnsUnauthorizedForInvalidScopes(t *testing.T) {
	status := errorToHTTPCode(oauth2.ErrInvalidScope)
	assert.Equal(t, http.StatusUnauthorized, status)
}

func TestErrorToHTTPCodeReturnsInternalServerErrorForGenericErrors(t *testing.T) {
	status := errorToHTTPCode(oauth2.ErrInvalidTokenType)
	assert.Equal(t, http.StatusInternalServerError, status)
	status = errorToHTTPCode(errors.New("foo"))
	assert.Equal(t, http.StatusInternalServerError, status)
}

func TestErrorToHTTPCodeReturnsBadRequestWhenInvalidToken(t *testing.T) {
	status := errorToHTTPCode(ErrInvalidToken)
	assert.Equal(t, http.StatusBadRequest, status)
}
