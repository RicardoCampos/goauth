package oauth2client

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewShouldNotAllowEmptyClientID(t *testing.T) {
	// Act
	c, err := New("", "bar", 0, []string{"read", "write"})

	// Assert

	assert.Nil(t, c, "Should not return a client")
	assert.NotNil(t, err, "Should return an error response")
}

func TestNewShouldNotAllowEmptyClientSecret(t *testing.T) {
	// Act
	c, err := New("foo", "", 0, []string{"read", "write"})

	// Assert

	assert.Nil(t, c, "Should not return a client")
	assert.NotNil(t, err, "Should return an error response")
}

func TestNewShouldNotAllowEmptyAllowedScopes(t *testing.T) {
	// Act
	c, err := New("foo", "bar", 0, []string{})

	// Assert

	assert.Nil(t, c, "Should not return a client")
	assert.NotNil(t, err, "Should return an error response")
}
func TestNewShouldNotAllowLifetimeLessThanMinimum(t *testing.T) {
	// Act
	c, err := New("foo", "bar", MinimumAccesTokenLifetimeInMs-1, []string{"read", "write"})

	// Assert

	assert.Nil(t, c, "Should not return a client")
	assert.NotNil(t, err, "Should return an error response")
}

func TestValidateScopesCatchesInvalidRequest(t *testing.T) {
	// Arrange
	c, _ := New("foo", "bar", 0, []string{"read", "write"})

	// Act
	ok, err := c.ValidateScopes([]string{"delete"})

	// Assert
	assert.NotNil(t, err, "Should return an error with an invalid request.")
	assert.False(t, ok, "Validation should fail when asking for a scope not in allowed list.")
}

func TestValidateScopesAllowsValidRequest(t *testing.T) {
	// Arrange
	c, _ := New("foo", "bar", 0, []string{"read", "write"})

	// Act
	ok, err := c.ValidateScopes([]string{"write"})

	// Assert
	assert.Nil(t, err, "Should not return an error with a valid request.")
	assert.True(t, ok, "Validation should pass when asking for a scope in the allowed list.")
}
