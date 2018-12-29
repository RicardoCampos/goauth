package services

import (
	"crypto/rsa"
	"fmt"

	jwt "github.com/dgrijalva/jwt-go"
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
