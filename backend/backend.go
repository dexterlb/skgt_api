package backend

import (
	"database/sql"
	"fmt"
	"math/rand"

	"github.com/jmoiron/sqlx"

	// postgresql driver
	_ "github.com/lib/pq"
)

// Backend manages the API database and any requests made towards it
type Backend struct {
	db *sqlx.DB
}

// New returns a Backend instance, initialising a connection to the
// specified postgresql database.
func New(dbURN string) (*Backend, error) {
	var db *sqlx.DB
	db, err := sqlx.Connect("postgres", dbURN)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %s", err)
	}

	return &Backend{
		db: db,
	}, nil
}

// Info returns human-readable API version information
func (b *Backend) Info() (string, error) {
	return "skgt-api, pre-release", nil
}

// CheckAPIKey checks if the given string is a correct API key
// (refer to NewAPIKey() for generating API keys)
func (b *Backend) CheckAPIKey(apiKey string) error {
	err := b.db.Get(&apiKey, "select value from api_key where value = $1", apiKey)
	switch err {
	case nil:
		return nil
	case sql.ErrNoRows:
		return ErrWrongAPIKey
	default:
		return fmt.Errorf("database error: %s", err)
	}
}

const apiKeySymbols = "abcdefghijklmnopqrstuvwxyz0123456789"

// NewAPIKey generates a valid API key and returns it
func (b *Backend) NewAPIKey() (string, error) {
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

// InitDB initialises the database (creates tables, constraints etc)
func (b *Backend) InitDB() error {
	_, err := b.db.Exec(schema)
	if err != nil {
		return fmt.Errorf("unable to execute schema: %s", err)
	}
	return nil
}

// DropDB drops the database, performing the reverse operations of those
// InitDB() does
func (b *Backend) DropDB() error {
	_, err := b.db.Exec(dropSchema)
	if err != nil {
		return fmt.Errorf("unable to drop schema: %s", err)
	}
	return nil
}
