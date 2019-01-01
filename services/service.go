package services

import (
	"crypto/rsa"
	"errors"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/ricardocampos/goauth/oauth2"
)

// OAuth2Service has a token endpoint you can send client credentials to in exchange for a wonderful JWT (or... string)
type OAuth2Service interface {
	Token(tokenRequest) (tokenResponse, error)
	Validate(validationRequest) (validationResponse, error)
}

type oAuth2Service struct {
	clientRepository oauth2.ClientRepository
	tokenRepository  oauth2.ReferenceTokenRepository
	rsaKey           *rsa.PrivateKey
}

func (svc oAuth2Service) Token(s tokenRequest) (tokenResponse, error) {
	if svc.clientRepository == nil {
		return tokenResponse{
			ErrMsg: "Cannot validate clients.",
			Err:    oauth2.ErrInvalidGrant,
		}, oauth2.ErrInvalidGrant
	}
	if svc.tokenRepository == nil {
		return tokenResponse{
			ErrMsg: "Cannot store tokens.",
			Err:    oauth2.ErrInvalidGrant,
		}, oauth2.ErrInvalidGrant
	}
	if s.grantType == "" || s.grantType != oauth2.ClientCredentialsGrantType {
		return tokenResponse{
			ErrMsg: "You have provided an invalid grant_type in your token request.",
			Err:    oauth2.ErrInvalidGrant,
		}, oauth2.ErrInvalidGrant
	}
	if s.scope == "" {
		return tokenResponse{
			ErrMsg: "Invalid scope in token request.",
			Err:    oauth2.ErrInvalidScope,
		}, oauth2.ErrInvalidScope
	}
	if s.clientID == "" {
		return tokenResponse{
			ErrMsg: "You have provided an invalid grant in your token request.",
			Err:    oauth2.ErrInvalidGrant,
		}, oauth2.ErrInvalidGrant
	}
	if s.clientSecret == "" {
		return tokenResponse{
			ErrMsg: "You have provided an invalid grant in your token request.",
			Err:    oauth2.ErrInvalidGrant,
		}, oauth2.ErrInvalidGrant
	}

	// VALIDATE THE CLIENT

	client, ok := svc.clientRepository.GetClient(s.clientID)
	if client == nil || !ok {
		return tokenResponse{
			Err: oauth2.ErrInvalidGrant,
		}, oauth2.ErrInvalidGrant
	}
	if client.ClientSecret() != s.clientSecret {
		return tokenResponse{
			Err: oauth2.ErrInvalidGrant,
		}, oauth2.ErrInvalidGrant
	}

	validated, validationErr := client.ValidateScopes(strings.Fields(s.scope))
	if validationErr != nil || !validated {
		return tokenResponse{
			Err: oauth2.ErrInvalidScope,
		}, oauth2.ErrInvalidScope
	}

	nbfEpochTimeMs := oauth2.GetNowInEpochTime()
	expiry := nbfEpochTimeMs + client.AccessTokenLifetime()
	jti := uuid.New().String()

	claims := jwt.StandardClaims{
		Audience:  "youlot", // TODO: Really, this is optional, but will come from the client struct
		ExpiresAt: expiry,
		Id:        jti,
		IssuedAt:  nbfEpochTimeMs,
		Issuer:    "goauth.xom",      // TODO: This is a constant and should be the URI of the issuing server.
		NotBefore: nbfEpochTimeMs,    // this should be the current time
		Subject:   client.ClientID(), // Slightly off-piste here. We are not an authentication server so using the sub claim to identify the client.
	}

	clientTokenType := client.TokenType()
	token, err := signToken(claims, svc.rsaKey)

	if err != nil {
		return tokenResponse{}, oauth2.ErrInvalidGrant
	}

	if clientTokenType == oauth2.ReferenceTokenType {
		referenceToken, err := oauth2.NewReferenceToken(jti, client.ClientID(), expiry, token)
		if err != nil {
			return tokenResponse{
				Err: oauth2.ErrInvalidGrant,
			}, oauth2.ErrInvalidGrant
		}
		addTokenError := svc.tokenRepository.AddToken(referenceToken)
		if addTokenError != nil {
			return tokenResponse{
				Err: oauth2.ErrInvalidGrant,
			}, errors.New("unable to save token")
		}
		return tokenResponse{
			AccessToken: jti,
			TokenType:   clientTokenType,
			ExpiresIn:   client.AccessTokenLifetime(),
			Scope:       s.scope,
		}, nil

	} else if clientTokenType == oauth2.BearerTokenType {

		return tokenResponse{
			AccessToken: token,
			TokenType:   clientTokenType,
			ExpiresIn:   client.AccessTokenLifetime(),
			Scope:       s.scope,
		}, nil
	}
	return tokenResponse{}, oauth2.ErrInvalidTokenType
}

func (svc oAuth2Service) Validate(r validationRequest) (validationResponse, error) {
	if svc.tokenRepository == nil {
		err := errors.New("missing required dependencies")
		return validationResponse{
			ErrMsg: "Cannot retrieve tokens.",
			Err:    err,
		}, err
	}
	token, ok, err := svc.tokenRepository.GetToken(r.tokenID)
	if err != nil {
		return validationResponse{
			Err: err,
		}, err
	}
	if !ok {
		return validationResponse{
			Err: ErrInvalidToken,
		}, ErrInvalidToken
	}
	if token != nil {
		// Now we validate the token

		return validationResponse{}, nil
	}
	return validationResponse{
		Err: ErrInvalidToken,
	}, ErrInvalidToken
}
