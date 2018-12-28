package oauth2

// ReferenceTokenRepository Stores reference tokens
type ReferenceTokenRepository interface {
	AddToken(token ReferenceToken) error
	GetToken(tokenID string) (ReferenceToken, bool, error)
	// RemoveToken(tokenID string) error
	// RemoveAllTokensForClient(clientID string) error
}
