package services

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/RicardoCampos/goauth/oauth2"
)

// NewTokenHandler creates a token endpoint handler
func NewTokenHandler(svc OAuth2Service) *httptransport.Server {
	return httptransport.NewServer(
		makeTokenEndpoint(svc),
		decodeTokenRequest,
		encodeTokenResponse,
	)
}

// the request should be x-www-form-urlencoded
// it MUST include grant_type
// it MUST include scope
// We should also validate the authorisation section, though
type tokenRequest struct {
	clientID     string
	clientSecret string
	grantType    string
	scope        string
}

// In reality we will want to respond with a JWT. ugh.
type tokenResponse struct {
	AccessToken string `json:"access_token,v"`
	TokenType   string `json:"token_type,omitempty"`
	ExpiresIn   int64  `json:"expires_in,omitempty"`
	Scope       string `json:"scope,omitempty"`
	ErrMsg      string `json:"errMsg,omitempty"` // errors don't JSON-marshal, so we use a string
	Err         error  `json:"-"`                // actual raw error, always omitted
}

func makeTokenEndpoint(svc OAuth2Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(tokenRequest)
		token, err := svc.Token(req)
		if err != nil {
			return token, nil
		}
		return token, nil
	}
}

// decodeTokenRequest will extract the necessary information from a token request. It will not validate the content itself, however.
func decodeTokenRequest(_ context.Context, r *http.Request) (interface{}, error) {
	parseError := r.ParseForm()
	if parseError != nil {
		return nil, parseError
	}
	var request tokenRequest
	username, password, ok := r.BasicAuth()
	if ok {
		//we use header value
		request.clientID = username
		request.clientSecret = password
	} else {
		// we use form value
		request.clientID = r.PostFormValue(oauth2.ClientID)
		request.clientSecret = r.PostFormValue(oauth2.ClientSecret)
	}
	request.grantType = r.PostFormValue(oauth2.GrantType)
	request.scope = r.PostFormValue(oauth2.Scope)
	return request, nil
}

func encodeTokenResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	token := response.(tokenResponse)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if token.Err != nil {
		w.WriteHeader(errorToHTTPCode(token.Err))
		var temp struct{}
		return json.NewEncoder(w).Encode(temp)
	}
	return json.NewEncoder(w).Encode(response)
}
