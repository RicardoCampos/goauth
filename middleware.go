package main

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
