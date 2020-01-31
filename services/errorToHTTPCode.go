package services

import (
	"net/http"

	"github.com/RicardoCampos/goauth/oauth2"
)

func errorToHTTPCode(err error) int {
	switch err {
	case oauth2.ErrInvalidGrant, oauth2.ErrInvalidScope:
		return http.StatusUnauthorized
	case ErrInvalidToken:
		return http.StatusBadRequest
	}
	return http.StatusInternalServerError
}
