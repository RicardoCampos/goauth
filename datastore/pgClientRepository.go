package datastore

import (
	"database/sql"
	"strings"

	"github.com/go-kit/kit/log"
	// we mask the actual driver for now
	_ "github.com/lib/pq"
	"github.com/ricardocampos/goauth/oauth2"
)

type pgClientRepository struct {
	db     *sql.DB
	logger log.Logger
}

//NewPostgresClientRepository creates a new repository backed by Postgres
func NewPostgresClientRepository(dataSourceName string, logger log.Logger) (oauth2.ClientRepository, error) {
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	repository := pgClientRepository{
		db,
		logger,
	}
	return repository, nil
}

// AddClient Adds a client to the in memory databsae
func (r pgClientRepository) AddClient(client oauth2.Client) {
	db := r.db
	stmt, err := db.Prepare("INSERT INTO public.clients(\"clientId\", \"clientSecret\", \"accessTokenLifetime\", \"tokenType\", \"allowedScopes\" ) VALUES($1,$2,$3,$4,$5)")
	if err != nil {
		r.logger.Log("msg", "Unable to prepare client insert statement", err)
		return
	}
	_, er := stmt.Exec(client.ClientID(), client.ClientSecret(), client.AccessTokenLifetime(), client.TokenType(), strings.Join(client.AllowedScopes(), " "))
	if er != nil {
		r.logger.Log("msg", "Unable to execute client insert statement", err)
	}
}

// GetClients gets an in memory array of clients
func (r pgClientRepository) GetClients() map[string]oauth2.Client {
	db := r.db
	stmt, err := db.Prepare("SELECT \"clientId\", \"clientSecret\", \"accessTokenLifetime\", \"tokenType\", \"allowedScopes\" FROM public.clients")
	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		r.logger.Log("msg", "Unable to prepare client get statement", err)
		return nil
	}
	defer rows.Close()
	clients := make(map[string]oauth2.Client)
	for rows.Next() {
		var (
			dbClientID          string
			clientSecret        string
			accessTokenLifetime int64
			tokenType           string
			allowedScopes       string
		)
		err := rows.Scan(&dbClientID, &clientSecret, &accessTokenLifetime, &tokenType, &allowedScopes)
		if err != nil {
			r.logger.Log("msg", "Unable to scan client get statement", err)
		}
		client, err := oauth2.NewClient(dbClientID, clientSecret, tokenType, accessTokenLifetime, strings.Fields(allowedScopes))
		if err != nil {
			r.logger.Log("msg", "Unable to create client from rows", err)
		}
		clients[dbClientID] = client
	}
	if err = rows.Err(); err != nil {
		r.logger.Log("msg", "Unknown error getting clients", err)
	}
	return clients
}

// GetClient gets a specified client
func (r pgClientRepository) GetClient(clientID string) (oauth2.Client, bool) {
	db := r.db
	var (
		dbClientID          string
		clientSecret        string
		accessTokenLifetime int64
		tokenType           string
		allowedScopes       string
	)
	stmt, err := db.Prepare("SELECT \"clientId\", \"clientSecret\", \"accessTokenLifetime\", \"tokenType\", \"allowedScopes\" FROM public.clients WHERE \"clientId\" = $1 LIMIT 1;")

	if err != nil {
		r.logger.Log("msg", "Could not compile the client query.", err)
		return nil, false
	}
	defer stmt.Close()
	err = stmt.QueryRow(clientID).Scan(&dbClientID, &clientSecret, &accessTokenLifetime, &tokenType, &allowedScopes)
	if err != nil {
		r.logger.Log("msg", "Unable to find a matching client.", err)
		return nil, false
	}
	client, err := oauth2.NewClient(dbClientID, clientSecret, tokenType, accessTokenLifetime, strings.Fields(allowedScopes))
	if err != nil {
		r.logger.Log("msg", "Failed to create a client from the row retrieved", err)
		return nil, false
	}
	return client, true
}
