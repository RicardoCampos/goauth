package datastore

import (
	"database/sql"
	"errors"
	"log"

	// we mask the actual driver for now
	_ "github.com/lib/pq"
	"github.com/ricardocampos/goauth/oauth2"
)

type pgTokenRepository struct {
	db     *sql.DB
	tokens map[string]oauth2.ReferenceToken
}

//NewPostgresTokenRepository creates a new repository backed by Postgres
func NewPostgresTokenRepository(dataSourceName string) (oauth2.ReferenceTokenRepository, error) {
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	repository := pgTokenRepository{
		db:     db,
		tokens: make(map[string]oauth2.ReferenceToken),
	}
	return repository, nil
}

func (r pgTokenRepository) AddToken(token oauth2.ReferenceToken) error {
	if token == nil {
		return errors.New("requires a valid token to store")
	}
	db := r.db
	stmt, err := db.Prepare("INSERT INTO public.tokens(\"tokenID\", \"clientID\", \"expiry\", \"accessToken\") VALUES ($1, $2, $3, $4);")
	if err != nil {
		log.Fatal(err)
		return err
	}
	_, err = stmt.Exec(token.TokenID(), token.ClientID(), token.Expiry(), token.AccessToken())
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

// GetToken Gets a token by ID
func (r pgTokenRepository) GetToken(tokenID string) (oauth2.ReferenceToken, bool, error) {
	if len(tokenID) < 1 {
		return nil, false, errors.New("please provide a valid tokenID")
	}
	db := r.db
	var (
		dbTokenID   string
		clientID    string
		expiry      int64
		accessToken string
	)
	stmt, err := db.Prepare("SELECT \"tokenID\", \"clientID\", \"expiry\", \"accessToken\" FROM public.tokens WHERE \"tokenID\" = $1 LIMIT 1;")

	if err != nil {
		return nil, false, err
	}
	defer stmt.Close()
	err = stmt.QueryRow(tokenID).Scan(&dbTokenID, &clientID, &expiry, &accessToken)
	if err != nil {
		return nil, false, err
	}
	token, err := oauth2.NewReferenceToken(dbTokenID, clientID, expiry, accessToken)
	if err != nil {
		return nil, false, err
	}
	return token, true, nil
}
