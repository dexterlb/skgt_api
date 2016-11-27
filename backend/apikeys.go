package backend

import (
	"database/sql"
	"fmt"
	"math/rand"
)

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

// DeleteAPIKey deletes an API key from the database
func (b *Backend) DeleteAPIKey(apiKey string) error {
	_, err := b.db.Exec("delete from api_key where value = $1", apiKey)
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
	keyBytes := make([]byte, 64)
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
