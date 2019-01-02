package services

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
)

// ErrInvalidRequest Invalid input
var ErrInvalidRequest = errors.New("Invalid Request")

// ErrInvalidToken Token not found or invalid
var ErrInvalidToken = errors.New("not found or expired")

type validationRequest struct {
	tokenID string
}

type validationResponse struct {
	ErrMsg string `json:"errMsg,omitempty"` // errors don't JSON-marshal, so we use a string
	Err    error  `json:"-"`                // actual raw error
}

// NewValidateHandler creates a token endpoint handler
func NewValidateHandler(svc OAuth2Service) *httptransport.Server {
	return httptransport.NewServer(
		makeValidateEndpoint(svc),
		decodeValidateRequest,
		encodeValidateResponse,
	)
}

func makeValidateEndpoint(svc OAuth2Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(validationRequest)
		result, err := svc.Validate(req)
		if err != nil {
			return result, nil
		}
		return result, nil
	}
}

func decodeValidateRequest(_ context.Context, r *http.Request) (interface{}, error) {
	parseError := r.ParseForm()
	if parseError != nil {
		return nil, parseError
	}
	tokenID := r.PostFormValue("token")
	if len(tokenID) < 1 {
		return nil, ErrInvalidRequest
	}
	req := validationRequest{
		tokenID,
	}
	return req, nil
}

func encodeValidateResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	valid := response.(validationResponse)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if valid.Err != nil {
		w.WriteHeader(errorToHTTPCode(valid.Err))
		var temp struct{}
		return json.NewEncoder(w).Encode(temp)
	}
	return json.NewEncoder(w).Encode(response)
}
