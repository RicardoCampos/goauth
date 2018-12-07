package main

import (
	"fmt"
	"io/ioutil"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

// OAuth2Service has a token endpoint you can send client credentials to in exchange for a wonderful JWT (or... string)
type OAuth2Service interface {
	Token(tokenRequest) (tokenResponse, error)
}

type oAuth2Service struct{}

type tokenPayload struct {
	Issuer string `json:"iss"`
	Scope  string `json:"scope"`
}

func (oAuth2Service) Token(s tokenRequest) (tokenResponse, error) {
	if s.grantType == "" || s.grantType != ClientCredentialsGrantType {
		return tokenResponse{
			Err: "You have provided an invalid grant_type in your token request."}, ErrInvalidGrant
	}
	if s.scope == "" {
		return tokenResponse{
			Err: "Invalid scope in token request."}, ErrInvalidScope
	}
	// TODO: this actually needs to authorise the client. If the username/password matches then we should check to see if it has that scope.
	// This way an invalid client could be detected.
	// This will also allow us to validate the scopes, and the bearer type.
	// We should also have the token lifetime from the client.
	// TODO: With reference tokens we only issue a UUID in the AccessToken field and we store the actual token in the data store.

	tokenValidTimeSpanMs := int64(300000) // 5 minutes
	nbfEpochTimeMs := time.Now().UnixNano() / 1000000
	expiry := nbfEpochTimeMs + tokenValidTimeSpanMs
	jti := uuid.New().String() // this is the reference UUID we would return for ref tokens.

	claims := jwt.StandardClaims{
		Audience:  "youlot", // Really, this is optional, but will come from the client struct
		ExpiresAt: expiry,
		Id:        jti,
		IssuedAt:  1544214970,
		Issuer:    "goauth.xom",   // This is a constant and should be the URI of the issuing server.
		NotBefore: nbfEpochTimeMs, // this should be the current time
	}
	accessToken, _ := signToken(claims)

	// tokenType must be retrieved from the client's details.
	return tokenResponse{
		AccessToken: accessToken,
		TokenType:   BearerTokenType,
		ExpiresIn:   expiry,
		Scope:       s.scope,
	}, nil
}

// signToken Create, sign, and output a token.
func signToken(input jwt.StandardClaims) (string, error) {

	k, _ := ioutil.ReadFile("sample_key")

	key, err := jwt.ParseRSAPrivateKeyFromPEM(k)
	if err != nil {
		return fmt.Sprintf("Failed to parse valid private key: %v", err), err
	}

	// get the signing alg
	alg := jwt.GetSigningMethod("RS512")
	if alg == nil {
		return "Couldn't find signing method", fmt.Errorf("Couldn't find signing method: %v", "RS512")
	}

	// create a new token
	token := jwt.NewWithClaims(alg, input)

	if out, err := token.SignedString(key); err == nil {
		// fmt.Println(out)
		return out, nil
	}
	return fmt.Sprintf("Error signing token: %v", err), fmt.Errorf("Error signing token: %v", err)
}
