package oauth2

import (
	"net/http"
)

func ErrorToHTTPCode(err error) int {
	switch err {
	case ErrInvalidGrant, ErrInvalidScope, ErrInvalidToken:
		return http.StatusUnauthorized
	case ErrInvalidRequest:
		return http.StatusBadRequest
	case ErrInsufficientScope:
		return http.StatusForbidden
	}
	return http.StatusInternalServerError
}
