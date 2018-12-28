package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/ricardocampos/goauth/oauth2"
)

// OAuth2Service has a token endpoint you can send client credentials to in exchange for a wonderful JWT (or... string)
type OAuth2Service interface {
	Token(tokenRequest) (tokenResponse, error)
}

type oAuth2Service struct {
	clientRepository oauth2.ClientRepository
	tokenRepository  oauth2.ReferenceTokenRepository
}

type tokenPayload struct {
	Issuer string `json:"iss"`
	Scope  string `json:"scope"`
}

func (svc oAuth2Service) Token(s tokenRequest) (tokenResponse, error) {
	if svc.clientRepository == nil {
		return tokenResponse{
			Err: "Cannot validate clients."}, oauth2.ErrInvalidGrant
	}
	if svc.tokenRepository == nil {
		return tokenResponse{
			Err: "Cannot store tokens."}, oauth2.ErrInvalidGrant
	}
	if s.grantType == "" || s.grantType != oauth2.ClientCredentialsGrantType {
		return tokenResponse{
			Err: "You have provided an invalid grant_type in your token request."}, oauth2.ErrInvalidGrant
	}
	if s.scope == "" {
		return tokenResponse{
			Err: "Invalid scope in token request."}, oauth2.ErrInvalidScope
	}
	if s.clientID == "" {
		return tokenResponse{
			Err: "You have provided an invalid grant in your token request."}, oauth2.ErrInvalidGrant
	}
	if s.clientSecret == "" {
		return tokenResponse{
			Err: "You have provided an invalid grant in your token request."}, oauth2.ErrInvalidGrant
	}

	// VALIDATE THE CLIENT

	client, ok := svc.clientRepository.GetClient(s.clientID)
	if client == nil || !ok {
		return tokenResponse{}, oauth2.ErrInvalidGrant
	}
	if client.ClientSecret() != s.clientSecret {
		return tokenResponse{}, oauth2.ErrInvalidGrant
	}

	validated, validationErr := client.ValidateScopes(strings.Fields(s.scope))
	if validationErr != nil || !validated {
		return tokenResponse{}, oauth2.ErrInvalidScope
	}

	nbfEpochTimeMs := oauth2.GetNowInEpochTime()
	expiry := nbfEpochTimeMs + client.AccessTokenLifetime()
	jti := uuid.New().String()

	claims := jwt.StandardClaims{
		Audience:  "youlot", // TODO: Really, this is optional, but will come from the client struct
		ExpiresAt: expiry,
		Id:        jti,
		IssuedAt:  nbfEpochTimeMs,
		Issuer:    "goauth.xom",   // TODO: This is a constant and should be the URI of the issuing server.
		NotBefore: nbfEpochTimeMs, // this should be the current time
	}

	clientTokenType := client.TokenType()
	accessToken := string(jti)
	token, _ := signToken(claims)
	if clientTokenType == oauth2.ReferenceTokenType {
		referenceToken, err := oauth2.NewReferenceToken(accessToken, client.ClientID(), expiry, token)
		if err != nil {
			return tokenResponse{}, oauth2.ErrInvalidGrant
		}
		addTokenError := svc.tokenRepository.AddToken(referenceToken)
		if addTokenError != nil {
			return tokenResponse{}, errors.New("unable to save token")
		}
	} else if clientTokenType == oauth2.BearerTokenType {
		accessToken = token
	}
	return tokenResponse{
		AccessToken: accessToken,
		TokenType:   clientTokenType,
		ExpiresIn:   client.AccessTokenLifetime(),
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
