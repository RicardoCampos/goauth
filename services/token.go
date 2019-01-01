package services

import (
	"crypto/rsa"
	"fmt"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/ricardocampos/goauth/oauth2"
)

type tokenPayload struct {
	Issuer string `json:"iss"`
	Scope  string `json:"scope"`
}

func signToken(input jwt.StandardClaims, rsaKey *rsa.PrivateKey) (string, error) {
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

func validateToken(r oauth2.ReferenceToken, key *rsa.PublicKey) bool {
	if r == nil {
		return false
	}

	// expiry check
	if r.Expiry() < oauth2.GetNowInEpochTime() {
		return false
	}

	token, err := jwt.Parse(r.AccessToken(), func(token *jwt.Token) (interface{}, error) {
		// Validate algorithm
		// TODO: we load by `kid` to support multiple certs
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			// TODO: return typed error
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return key, nil
	})
	if err != nil {
		return false
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Because we are paranoid, check that the data in the token matches the convenience fields in the reference token (in case database has compromises)
		if claims["jti"] != r.TokenID() {
			return false
		}
		if claims["exp"] != float64(r.Expiry()) { //TODO: raise PR as this is a bug. All inputs are int64 but it gets mapped to float64
			return false
		}
		if claims["sub"] != r.ClientID() {
			return false
		}
		return true
	}
	return false
}
