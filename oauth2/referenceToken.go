package oauth2

import (
	"errors"
)

type ReferenceToken interface {
	TokenID() string
	ClientID() string
	Expiry() int64
	AccessToken() string
}

// referenceToken stores all relevant information for reference tokens
type referenceToken struct {
	tokenID     string
	clientID    string
	expiry      int64
	accessToken string
}

func (r referenceToken) TokenID() string {
	return r.tokenID
}

func (r referenceToken) ClientID() string {
	return r.clientID
}

func (r referenceToken) Expiry() int64 {
	return r.expiry
}

func (r referenceToken) AccessToken() string {
	return r.accessToken
}

// NewReferenceToken creates a reference token
func NewReferenceToken(tokenID string, clientID string, expiry int64, accessToken string) (ReferenceToken, error) {
	if len(tokenID) < 1 {
		return nil, errors.New("tokenID is required")
	}
	if len(clientID) < 1 {
		return nil, errors.New("clientID is required")
	}
	if len(accessToken) < 1 {
		return nil, errors.New("accessToken is required")
	}
	return referenceToken{
		tokenID:     tokenID,
		clientID:    clientID,
		expiry:      expiry,
		accessToken: accessToken,
	}, nil
}
