package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/ricardocampos/goauth/oauth2"
	"github.com/ricardocampos/goauth/services"
	log "github.com/sirupsen/logrus"
)

type ValidateHandler interface {
	HandleValidation(http.ResponseWriter, *http.Request)
}

type validateHandler struct {
	logger       *log.Logger
	tokenService services.TokenService
}

func NewValidateHandler(logger *log.Logger, svc services.TokenService) ValidateHandler {
	return validateHandler{logger, svc}
}

// ValidateHandler handles incoming requests to validate reference tokens.
func (handler validateHandler) HandleValidation(w http.ResponseWriter, r *http.Request) {
	handler.logger.Debug("Handling validate request")
	request, err := decodeValidateRequest(r, handler.logger)
	if err != nil {
		handler.logger.Debug("Unable to decode incoming request")
		encodeValidateResponse(w, oauth2.ValidationResponse{Err: err})
		return
	}
	response, _ := handler.tokenService.Validate(request)
	encodeValidateResponse(w, response)
}

func decodeValidateRequest(r *http.Request, logger *log.Logger) (oauth2.ValidationRequest, error) {
	parseError := r.ParseForm()
	if parseError != nil {
		logger.Debug("Unable to parse incoming form request")
		return oauth2.ValidationRequest{Err: parseError}, parseError
	}
	tokenID := r.PostFormValue("token")
	if len(tokenID) < 1 {
		logger.Debug("token value was missing")
		return oauth2.ValidationRequest{Err: oauth2.ErrInvalidRequest}, oauth2.ErrInvalidRequest
	}
	scopes := strings.Fields(r.PostFormValue("expectedScope"))

	if len(scopes) < 1 {
		logger.Debug("expectedScope was missing.")
		return oauth2.ValidationRequest{Err: oauth2.ErrInvalidRequest}, oauth2.ErrInvalidRequest
	}
	req := oauth2.ValidationRequest{
		TokenID: tokenID,
		Scopes:  scopes,
	}
	return req, nil
}

func encodeValidateResponse(w http.ResponseWriter, response oauth2.ValidationResponse) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if response.Err != nil {
		w.WriteHeader(oauth2.ErrorToHTTPCode(response.Err))
		var temp struct{}
		return json.NewEncoder(w).Encode(temp)
	}
	return json.NewEncoder(w).Encode(response)
}
