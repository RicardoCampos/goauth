package datastore

import (
	"testing"

	"github.com/RicardoCampos/goauth/oauth2"
	"github.com/stretchr/testify/assert"
)

func TestInMemoryClientRepositoryAddClient(t *testing.T) {
	// Arrange/Act
	repository := NewInMemoryClientRepository()

	// Act
	c, _ := oauth2.NewClient("foo", "bar", oauth2.ReferenceTokenType, 0, []string{"read", "write"})
	repository.AddClient(c)

	// Assert

	assert.Equal(t, "foo", c.ClientID())
}

func TestInMemoryClientRepositoryGetClients(t *testing.T) {
	// Arrange
	repository := NewInMemoryClientRepository()
	c, _ := oauth2.NewClient("foo", "bar", oauth2.ReferenceTokenType, 0, []string{"read", "write"})
	repository.AddClient(c)

	//Act
	clients := repository.GetClients()

	// Assert
	assert.Equal(t, 1, len(clients))
}

func TestInMemoryClientRepositoryGetClientsWhenNoneAdded(t *testing.T) {
	// Arrange
	repository := NewInMemoryClientRepository()

	//Act
	clients := repository.GetClients()

	// Assert
	assert.Equal(t, 0, len(clients))
}

func TestInMemoryClientRepositoryGetClient(t *testing.T) {
	// Arrange
	repository := NewInMemoryClientRepository()
	c, _ := oauth2.NewClient("foo", "bar", oauth2.ReferenceTokenType, 0, []string{"read", "write"})
	repository.AddClient(c)

	// Act
	sut, _ := repository.GetClient("foo")

	// Assert
	assert.Equal(t, "foo", sut.ClientID())
}

func TestInMemoryClientRepositoryGetClientThatDoesNotExist(t *testing.T) {
	// Arrange/Act
	repository := NewInMemoryClientRepository()

	// Act
	c, _ := repository.GetClient("bar")

	// Assert
	assert.Nil(t, c)
}
