package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/endpoint"
)

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
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
	Scope       string `json:"scope"`
	Err         string `json:"err,omitempty"` // errors don't JSON-marshal, so we use a string
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
		request.clientID = r.PostFormValue(ClientID)
		request.clientSecret = r.PostFormValue(ClientSecret)
	}
	request.grantType = r.PostFormValue(GrantType)
	request.scope = r.PostFormValue(Scope)
	return request, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
