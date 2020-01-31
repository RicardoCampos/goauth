package datastore

import "github.com/RicardoCampos/goauth/oauth2"

// InMemoryClientRepository is a hard coded, in-memory only placeholder
type inMemoryClientRepository struct {
	clients map[string]oauth2.Client
}

// NewInMemoryClientRepository is
func NewInMemoryClientRepository() oauth2.ClientRepository {
	return &inMemoryClientRepository{
		clients: make(map[string]oauth2.Client),
	}
}

// AddClient Adds a client to the in memory database
func (r *inMemoryClientRepository) AddClient(client oauth2.Client) {
	if len(r.clients) < 1 {
		r.clients = make(map[string]oauth2.Client)
	}
	r.clients[client.ClientID()] = client
}

// GetClients gets an in memory array of clients
func (r inMemoryClientRepository) GetClients() map[string]oauth2.Client {
	return r.clients
}

// GetClient gets a specified client
func (r inMemoryClientRepository) GetClient(clientID string) (oauth2.Client, bool) {
	v, ok := r.clients[clientID]
	return v, ok
}
