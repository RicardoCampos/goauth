package oauth2

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewShouldNotAllowEmptyClientID(t *testing.T) {
	// Act
	c, err := NewClient("", "bar", ReferenceTokenType, 0, []string{"read", "write"})

	// Assert

	assert.Nil(t, c, "Should not return a client")
	assert.NotNil(t, err, "Should return an error response")
}

func TestNewShouldNotAllowEmptyClientSecret(t *testing.T) {
	// Act
	c, err := NewClient("foo", "", ReferenceTokenType, 0, []string{"read", "write"})

	// Assert

	assert.Nil(t, c, "Should not return a client")
	assert.NotNil(t, err, "Should return an error response")
}

func TestNewShouldAllowBearerTokenTypes(t *testing.T) {
	// Act
	c, err := NewClient("foo", "bar", BearerTokenType, 0, []string{"read", "write"})

	// Assert

	assert.NotNil(t, c, "Should return a client")
	assert.Nil(t, err, "Should not return an error response")
}

func TestNewShouldAllowReferenceTokenTypes(t *testing.T) {
	// Act
	c, err := NewClient("foo", "bar", ReferenceTokenType, 0, []string{"read", "write"})

	// Assert

	assert.NotNil(t, c, "Should return a client")
	assert.Nil(t, err, "Should not return an error response")
}

func TestNewShouldNotAllowInvalidTokenTypes(t *testing.T) {
	// Act
	c, err := NewClient("foo", "bar", "blarg", 0, []string{"read", "write"})

	// Assert

	assert.Nil(t, c, "Should not return a client")
	assert.NotNil(t, err, "Should return an error response")
	assert.Equal(t, ErrInvalidTokenType, err)
}

func TestNewShouldNotAllowEmptyAllowedScopes(t *testing.T) {
	// Act
	c, err := NewClient("foo", "bar", ReferenceTokenType, 0, []string{})

	// Assert

	assert.Nil(t, c, "Should not return a client")
	assert.NotNil(t, err, "Should return an error response")
}
func TestNewShouldNotAllowLifetimeLessThanMinimum(t *testing.T) {
	// Act
	c, err := NewClient("foo", "bar", ReferenceTokenType, MinimumAccesTokenLifetimeInS-1, []string{"read", "write"})

	// Assert

	assert.Nil(t, c, "Should not return a client")
	assert.NotNil(t, err, "Should return an error response")
}

func TestValidateScopesCatchesInvalidRequest(t *testing.T) {
	// Arrange
	c, _ := NewClient("foo", "bar", ReferenceTokenType, 0, []string{"read", "write"})

	// Act
	ok, err := c.ValidateScopes([]string{"delete"})

	// Assert
	assert.NotNil(t, err, "Should return an error with an invalid request.")
	assert.False(t, ok, "Validation should fail when asking for a scope not in allowed list.")
}

func TestValidateScopesAllowsValidRequest(t *testing.T) {
	// Arrange
	c, _ := NewClient("foo", "bar", ReferenceTokenType, 0, []string{"read", "write"})

	// Act
	ok, err := c.ValidateScopes([]string{"write"})

	// Assert
	assert.Nil(t, err, "Should not return an error with a valid request.")
	assert.True(t, ok, "Validation should pass when asking for a scope in the allowed list.")
}
