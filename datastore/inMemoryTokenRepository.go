package datastore

import (
	"errors"

	"github.com/ricardocampos/goauth/oauth2"
)

type inMemoryReferenceTokenRepository struct {
	tokens map[string]oauth2.ReferenceToken
}

func (r inMemoryReferenceTokenRepository) AddToken(token oauth2.ReferenceToken) error {
	if token == nil {
		return errors.New("you cannot add an empty token")
	}
	if r.tokens != nil {
		r.tokens[token.TokenID()] = token
		return nil
	}
	return errors.New("not implemented")
}

// GetToken Gets a token by ID
func (r inMemoryReferenceTokenRepository) GetToken(tokenID string) (oauth2.ReferenceToken, bool, error) {
	if len(tokenID) < 1 {
		return nil, false, errors.New("please provide a valid tokenID")
	}
	v, ok := r.tokens[tokenID]
	if ok {
		return v, true, nil
	}
	return nil, false, errors.New("could not find token")
}

// NewInMemoryReferenceTokenRepository returns an inMemoryReferenceTokenRepository
func NewInMemoryReferenceTokenRepository() oauth2.ReferenceTokenRepository {
	return inMemoryReferenceTokenRepository{
		tokens: make(map[string]oauth2.ReferenceToken),
	}
}
