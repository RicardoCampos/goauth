package datastore

import (
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/google/uuid"
	"github.com/RicardoCampos/goauth/oauth2"
	"github.com/stretchr/testify/assert"
)

func getRepository() oauth2.ClientRepository {
	db, _ := NewPostgresClientRepository("postgres://postgres:password@localhost/goauth?sslmode=disable", log.NewNopLogger())
	return db
}
func TestPgRepositoryGetClient(t *testing.T) {
	// Arrange
	repository := getRepository()
	clientID := "foo_bearer"
	//  Act

	c, ok := repository.GetClient(clientID)

	// Assert

	assert.True(t, ok)
	assert.NotNil(t, c, "Should return the client if it is in the database")
	assert.Equal(t, clientID, c.ClientID(), "The clientID's should match")
}

func TestPgRepositoryGetClients(t *testing.T) {
	// Arrange
	repository := getRepository()

	//  Act

	clients := repository.GetClients()

	// Assert

	assert.NotNil(t, clients, "Should return the clients in the database")
}

func TestPgRepositoryAddClient(t *testing.T) {
	// Arrange
	repository := getRepository()
	//Create random client ID as the repository does not currently allow deleting
	clientID := uuid.New().String()

	//  Act
	c, _ := oauth2.NewClient(clientID, "bar", oauth2.ReferenceTokenType, 0, []string{"read", "write"})
	repository.AddClient(c)

	// Assert
	retrieved, ok := repository.GetClient(clientID)
	assert.True(t, ok)
	assert.NotNil(t, retrieved)
	assert.Equal(t, clientID, retrieved.ClientID())
}
