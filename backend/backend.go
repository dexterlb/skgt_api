package backend

import (
	"database/sql"
	"fmt"
	"math/rand"

	"github.com/jmoiron/sqlx"

	// postgresql driver
	_ "github.com/lib/pq"
)

type Backend struct {
	db *sqlx.DB
}

func NewBackend(dbURN string) (*Backend, error) {
	var db *sqlx.DB
	db, err := sqlx.Connect("postgres", dbURN)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %s", err)
	}

	return &Backend{
		db: db,
	}, nil
}

func (b *Backend) Info() (string, error) {
	return "this is foo.", nil
}

func (b *Backend) CheckApiKey(apiKey string) error {
	err := b.db.Get(&apiKey, "select value from api_key where value = $1", apiKey)
	switch err {
	case nil:
		return nil
	case sql.ErrNoRows:
		return ErrWrongApiKey
	default:
		return fmt.Errorf("database error: %s", err)
	}
}

const apiKeySymbols = "abcdefghijklmnopqrstuvwxyz0123456789"

func (b *Backend) NewApiKey() (string, error) {
	keyBytes := make([]byte, 256)
	for i := range keyBytes {
		keyBytes[i] = apiKeySymbols[rand.Intn(len(apiKeySymbols))]
	}

	key := string(keyBytes)

	_, err := b.db.Exec("insert into api_key(value) values($1)", key)
	if err != nil {
		return "", fmt.Errorf("unable to create api key: %s", err)
	}
	return key, nil
}

func (b *Backend) InitDB() error {
	_, err := b.db.Exec(schema)
	if err != nil {
		return fmt.Errorf("unable to execute schema: %s", err)
	}
	return nil
}

func (b *Backend) DropDB() error {
	_, err := b.db.Exec(dropSchema)
	if err != nil {
		return fmt.Errorf("unable to drop schema: %s", err)
	}
	return nil
}
