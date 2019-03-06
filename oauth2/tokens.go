package oauth2

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
)

type TokenRequest struct {
	ClientID     string
	ClientSecret string
	GrantType    string
	Scope        string
}

type TokenResponse struct {
	AccessToken string `json:"access_token,v"`
	TokenType   string `json:"token_type,omitempty"`
	ExpiresIn   int64  `json:"expires_in,omitempty"`
	Scope       string `json:"scope,omitempty"`
	ErrMsg      string `json:"errMsg,omitempty"` // errors don't JSON-marshal, so we use a string
	Err         error  `json:"-"`                // actual raw error, always omitted
}

type TokenPayload struct {
	Scope string `json:"scope"`
	jwt.StandardClaims
}

// Valid returns the validation ffrom the underlying standard claims object
func (t TokenPayload) Valid() error {
	return t.StandardClaims.Valid()
}

// SignToken cryptographically signs a JWT with a private key
func SignToken(input jwt.Claims, rsaKey *rsa.PrivateKey) (string, error) {
	// get the signing alg
	alg := jwt.GetSigningMethod("RS512")
	if alg == nil {
		return "Couldn't find signing method", fmt.Errorf("Couldn't find signing method: %v", "RS512")
	}
	// create a new token
	token := jwt.NewWithClaims(alg, input)
	// sign it
	out, err := token.SignedString(rsaKey)
	if err == nil {
		return out, nil
	}
	return fmt.Sprintf("Error signing token: %v", err), fmt.Errorf("Error signing token: %v", err)
}

func ValidateToken(r ReferenceToken, expectedScopes []string, key *rsa.PublicKey) (bool, error) {
	if r == nil {
		return false, ErrInvalidToken
	}

	// expiry check
	if r.Expiry() < GetNowInEpochTime() {
		return false, ErrInvalidToken
	}
	at := r.AccessToken()
	token, err := jwt.Parse(at, func(token *jwt.Token) (interface{}, error) {
		// Validate algorithm
		// TODO: we load by `kid` to support multiple certs
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			// TODO: return typed error
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return key, nil
	})
	if err != nil {
		return false, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Because we are paranoid, check that the data in the token matches the convenience fields in the reference token (in case database has compromises)
		if claims["jti"] != r.TokenID() {
			return false, ErrInvalidToken
		}
		if claims["exp"] != float64(r.Expiry()) { //TODO: raise PR as this is a bug. All inputs are int64 but it gets mapped to float64
			return false, ErrInvalidToken
		}
		if claims["sub"] != r.ClientID() {
			return false, ErrInvalidToken
		}

		// Check to see if scope is included
		if len(expectedScopes) > 0 {
			includedScopes := strings.Fields(claims["scope"].(string))
			return validateScopes(includedScopes, expectedScopes)
		}

		return true, nil
	}
	return false, ErrInvalidToken
}

func validateScopes(includedScopes []string, requestedScopes []string) (bool, error) {
	if includedScopes == nil || len(includedScopes) < 1 {
		return false, errors.New("at least one scope must be allowed")
	}
	if requestedScopes == nil || len(requestedScopes) < 1 {
		return false, errors.New("at least one requested scope must be supplied")
	}
	for _, requestedScope := range requestedScopes {
		found := false
		for _, tokenScope := range includedScopes {
			if tokenScope == requestedScope {
				found = true
			}
		}
		if !found {
			return false, ErrInsufficientScope
		}
	}
	return true, nil
}
