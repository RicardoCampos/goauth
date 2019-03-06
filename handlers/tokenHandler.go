package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ricardocampos/goauth/oauth2"
	"github.com/ricardocampos/goauth/services"
	log "github.com/sirupsen/logrus"
)

// TokenHandler handles tokens
type TokenHandler interface {
	HandleToken(http.ResponseWriter, *http.Request)
}

type tokenHandler struct {
	logger       *log.Logger
	tokenService services.TokenService
}

func NewTokenHandler(logger *log.Logger, svc services.TokenService) TokenHandler {
	return tokenHandler{logger, svc}
}

// TokenHandler handles incoming requests for new tokens.
func (handler tokenHandler) HandleToken(w http.ResponseWriter, r *http.Request) {
	handler.logger.Debug("Handling token request")
	request, err := decodeTokenRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	response, _ := handler.tokenService.CreateToken(request)
	encodeTokenResponse(w, response)
}

// decodeTokenRequest will extract the necessary information from a token request. It will not validate the content itself, however.
func decodeTokenRequest(r *http.Request) (oauth2.TokenRequest, error) {
	parseError := r.ParseForm()
	if parseError != nil {
		return oauth2.TokenRequest{}, parseError
	}
	var request oauth2.TokenRequest
	username, password, ok := r.BasicAuth()
	if ok {
		//we use header value
		request.ClientID = username
		request.ClientSecret = password
	} else {
		// we use form value
		request.ClientID = r.PostFormValue(oauth2.ClientID)
		request.ClientSecret = r.PostFormValue(oauth2.ClientSecret)
	}
	request.GrantType = r.PostFormValue(oauth2.GrantType)
	request.Scope = r.PostFormValue(oauth2.Scope)
	return request, nil
}

func encodeTokenResponse(w http.ResponseWriter, response oauth2.TokenResponse) error {
	token := response
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Pragma", "no-cache")
	if token.Err != nil {
		w.WriteHeader(oauth2.ErrorToHTTPCode(token.Err))
		var temp struct{}
		return json.NewEncoder(w).Encode(temp)
	}
	return json.NewEncoder(w).Encode(response)
}
