package oauth2

type ValidationRequest struct {
	TokenID string
	Scopes  []string
	Err     error
}

type ValidationResponse struct {
	ErrMsg string `json:"errMsg,omitempty"` // errors don't JSON-marshal, so we use a string
	Err    error  `json:"-"`                // actual raw error
}
