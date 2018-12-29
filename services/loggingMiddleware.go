package services

import (
	"encoding/json"
	"time"

	"github.com/go-kit/kit/log"
)

type loggingMiddleware struct {
	logger log.Logger
	next   OAuth2Service
}

func (mw loggingMiddleware) Token(input tokenRequest) (output tokenResponse, err error) {
	defer func(begin time.Time) {
		inputJSON, _ := json.Marshal(input)
		outputJSON, _ := json.Marshal(output)
		mw.logger.Log(
			"method", "token",
			"input", inputJSON,
			"output", outputJSON,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	output, err = mw.next.Token(input)
	return
}

// Validate Implements OAuth2Service
func (mw loggingMiddleware) Validate(input validationRequest) (output validationResponse, err error) {
	defer func(begin time.Time) {
		inputJSON, _ := json.Marshal(input)
		outputJSON, _ := json.Marshal(output)
		mw.logger.Log(
			"method", "validate",
			"input", inputJSON,
			"output", outputJSON,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	output, err = mw.next.Validate(input)
	return
}

// NewLoggingMiddleware returns new NewLoggingMiddleware instance
func NewLoggingMiddleware(logger log.Logger, svc OAuth2Service) OAuth2Service {
	return loggingMiddleware{logger: logger, next: svc}
}
