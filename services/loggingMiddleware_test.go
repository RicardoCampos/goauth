package services

import (
	"bytes"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/stretchr/testify/assert"
)

type fakeOAuth2Service struct{}

func (svc fakeOAuth2Service) Token(input tokenRequest) (output tokenResponse, err error) {
	return tokenResponse{}, nil
}
func (svc fakeOAuth2Service) Validate(input validationRequest) (output validationResponse, err error) {
	return validationResponse{}, nil
}

func TestLoggingMiddlewareActuallyLogs(t *testing.T) {
	// Arrange
	svc := fakeOAuth2Service{}
	var b bytes.Buffer
	logger := log.NewJSONLogger(&b)
	sut := NewLoggingMiddleware(logger, svc)

	// Act
	response, err := sut.Token(tokenRequest{})

	loggedEvent := b.String()
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.NotEmpty(t, loggedEvent)

}
