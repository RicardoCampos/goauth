package oauth2

// ClientRepository a Repository of OAuth2Clients
type ClientRepository interface {
	AddClient(client Client)
	GetClients() map[string]Client
	GetClient(clientID string) (Client, bool)
}
