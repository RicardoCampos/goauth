package services

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/ricardocampos/goauth/oauth2"
	log "github.com/sirupsen/logrus"
)

// TokenService has a token endpoint you can send client credentials to in exchange for a wonderful JWT (or... string)
type TokenService interface {
	CreateToken(oauth2.TokenRequest) (oauth2.TokenResponse, error)
	Validate(oauth2.ValidationRequest) (oauth2.ValidationResponse, error)
}

type tokenService struct {
	clientRepository *oauth2.ClientRepository
	tokenRepository  *oauth2.ReferenceTokenRepository
	rsaPrivateKey    *rsa.PrivateKey
	rsaPublicKey     *rsa.PublicKey
	logger           *log.Logger
}

func (svc tokenService) TokenRepository() oauth2.ReferenceTokenRepository {
	return *svc.tokenRepository
}
func (svc tokenService) ClientRepository() oauth2.ClientRepository {
	return *svc.clientRepository
}

func (svc tokenService) CreateToken(s oauth2.TokenRequest) (oauth2.TokenResponse, error) {
	svc.logger.Debug("Creating token")
	client, response, err := svc.ValidateTokenRequest(s)

	if err != nil {
		svc.logger.Warn("Error validating incoming token request", err)
		return response, oauth2.ErrInvalidGrant
	}
	if client == nil {
		panic("Client should not be empty if the error was nil.")
	}

	nbfEpochTimeMs := oauth2.GetNowInEpochTime()
	expiry := nbfEpochTimeMs + client.AccessTokenLifetime()
	jti := uuid.New().String()

	claims := oauth2.TokenPayload{
		Scope: s.Scope,
		StandardClaims: jwt.StandardClaims{
			Audience:  "youlot", // TODO: Really, this is optional, but will come from the client struct
			ExpiresAt: expiry,
			Id:        jti,
			IssuedAt:  nbfEpochTimeMs,
			Issuer:    "goauth.xom",      // TODO: This is a constant and should be the URI of the issuing server.
			NotBefore: nbfEpochTimeMs,    // this should be the current time
			Subject:   client.ClientID(), // Slightly off-piste here. We are not an authentication server so using the sub claim to identify the client.
		},
	}

	clientTokenType := client.TokenType()
	token, err := oauth2.SignToken(claims, svc.rsaPrivateKey)

	if err != nil {
		svc.logger.Error("Error signing token", err)
		return oauth2.TokenResponse{}, oauth2.ErrInvalidGrant
	}

	if clientTokenType == oauth2.ReferenceTokenType {
		referenceToken, err := oauth2.NewReferenceToken(jti, client.ClientID(), expiry, token)
		if err != nil {
			svc.logger.Error("Error creating reference token", err)
			return oauth2.TokenResponse{
				Err: oauth2.ErrInvalidGrant,
			}, oauth2.ErrInvalidGrant
		}
		addTokenError := svc.TokenRepository().AddToken(referenceToken)
		if addTokenError != nil {
			svc.logger.Error("Error adding reference token", addTokenError)
			return oauth2.TokenResponse{
				Err: oauth2.ErrInvalidGrant,
			}, errors.New("unable to save token")
		}
		return oauth2.TokenResponse{
			AccessToken: jti,
			TokenType:   clientTokenType,
			ExpiresIn:   client.AccessTokenLifetime(),
			Scope:       s.Scope,
		}, nil

	} else if clientTokenType == oauth2.BearerTokenType {

		return oauth2.TokenResponse{
			AccessToken: token,
			TokenType:   clientTokenType,
			ExpiresIn:   client.AccessTokenLifetime(),
			Scope:       s.Scope,
		}, nil
	}
	return oauth2.TokenResponse{}, oauth2.ErrInvalidTokenType
}

// ValidateTokenRequest Validates a request for a token.
func (svc tokenService) ValidateTokenRequest(s oauth2.TokenRequest) (oauth2.Client, oauth2.TokenResponse, error) {
	if svc.clientRepository == nil {
		svc.logger.Error("Client repository should not be nil")
		return nil, oauth2.TokenResponse{
			ErrMsg: "Cannot validate clients.",
			Err:    oauth2.ErrInvalidGrant,
		}, oauth2.ErrInvalidGrant
	}
	if svc.tokenRepository == nil {
		svc.logger.Error("Token repository should not be nil")
		return nil, oauth2.TokenResponse{
			ErrMsg: "Cannot store tokens.",
			Err:    oauth2.ErrInvalidGrant,
		}, oauth2.ErrInvalidGrant
	}
	if s.GrantType == "" || s.GrantType != oauth2.ClientCredentialsGrantType {
		ErrMsg := "You have provided an invalid grant_type in your token request."
		svc.logger.Debug(ErrMsg)
		return nil, oauth2.TokenResponse{
			ErrMsg: ErrMsg,
			Err:    oauth2.ErrInvalidGrant,
		}, oauth2.ErrInvalidGrant
	}
	if s.Scope == "" {
		ErrMsg := "Invalid scope in token request."
		svc.logger.Debug(ErrMsg)
		return nil, oauth2.TokenResponse{
			ErrMsg: ErrMsg,
			Err:    oauth2.ErrInvalidScope,
		}, oauth2.ErrInvalidScope
	}
	if s.ClientID == "" {
		svc.logger.Debug("Client ID is missing.")
		return nil, oauth2.TokenResponse{
			ErrMsg: "You have provided an invalid grant in your token request.",
			Err:    oauth2.ErrInvalidGrant,
		}, oauth2.ErrInvalidGrant
	}
	if s.ClientSecret == "" {
		svc.logger.Debug("Client secret is missing")
		return nil, oauth2.TokenResponse{
			ErrMsg: "You have provided an invalid grant in your token request.",
			Err:    oauth2.ErrInvalidGrant,
		}, oauth2.ErrInvalidGrant
	}

	// VALIDATE THE CLIENT

	client, ok := svc.ClientRepository().GetClient(s.ClientID)
	if client == nil || !ok {
		svc.logger.Debug("Cannot find client.")
		return nil, oauth2.TokenResponse{
			Err: oauth2.ErrInvalidGrant,
		}, oauth2.ErrInvalidGrant
	}
	if client.ClientSecret() != s.ClientSecret {
		svc.logger.Debug("Client secret does not match.")
		return nil, oauth2.TokenResponse{
			Err: oauth2.ErrInvalidGrant,
		}, oauth2.ErrInvalidGrant
	}

	validated, validationErr := client.ValidateScopes(strings.Fields(s.Scope))
	if validationErr != nil || !validated {
		svc.logger.Debug("Scopes requested does not match the clients allowed scopes.")
		return nil, oauth2.TokenResponse{
			Err: oauth2.ErrInvalidScope,
		}, oauth2.ErrInvalidScope
	}

	return client, oauth2.TokenResponse{}, nil
}

func (svc tokenService) Validate(r oauth2.ValidationRequest) (oauth2.ValidationResponse, error) {
	if svc.tokenRepository == nil {
		err := errors.New("missing required dependencies")
		return oauth2.ValidationResponse{
			ErrMsg: "Cannot retrieve tokens.",
			Err:    err,
		}, err
	}
	token, ok, _ := svc.TokenRepository().GetToken(r.TokenID)
	if !ok {
		return oauth2.ValidationResponse{
			Err: oauth2.ErrInvalidToken,
		}, oauth2.ErrInvalidToken
	}
	if token != nil {
		result, err := oauth2.ValidateToken(token, r.Scopes, svc.rsaPublicKey)
		// Now we validate the token
		if err != nil || !result {
			return oauth2.ValidationResponse{Err: err, ErrMsg: "Invalid Token"}, err
		}
		return oauth2.ValidationResponse{}, nil
	}
	return oauth2.ValidationResponse{
		Err: oauth2.ErrInvalidToken,
	}, oauth2.ErrInvalidToken
}

func validateScopes(allowedScopes []string, requestedScopes []string) (bool, error) {
	if allowedScopes == nil || len(allowedScopes) < 1 {
		return false, errors.New("at least one scope must be allowed")
	}
	if requestedScopes == nil || len(requestedScopes) < 1 {
		return false, errors.New("at least one requested scope must be supplied")
	}
	for _, requestedScope := range requestedScopes {
		found := false
		for _, tokenScope := range allowedScopes {
			if tokenScope == requestedScope {
				found = true
			}
		}
		if !found {
			return false, fmt.Errorf("requested %s scope not allowed for this client", requestedScope)
		}
	}
	return true, nil
}
