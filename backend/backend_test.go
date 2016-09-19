package backend

import (
	"flag"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var dbURN string

func init() {
	var (
		dbName string
		dbUser string
	)

	flag.StringVar(&dbName, "db.name", "", "name of database to connect to")
	flag.StringVar(&dbUser, "db.user", "", "username to connect to database with")

	flag.Parse()

	dbURN = fmt.Sprintf(
		"user=%s dbname=%s sslmode=disable", dbUser, dbName,
	)
}

func openBackend(t *testing.T) (*Backend, string) {
	backend, err := NewBackend(dbURN)
	if err != nil {
		t.Fatalf("cannot create backend: %s", err)
	}

	err = backend.InitDB()
	if err != nil {
		t.Errorf("cannot create database: %s", err)
	}

	apiKey, err := backend.NewApiKey()
	if err != nil {
		t.Fatalf("cannot create api key: %s", err)
	}

	return backend, apiKey
}

func closeBackend(t *testing.T, backend *Backend) {
	err := backend.DropDB()
	if err != nil {
		t.Errorf("cannot drop database: %s", err)
	}
}

func TestBackend_Info(t *testing.T) {
	assert := assert.New(t)
	backend, apiKey := openBackend(t)
	defer closeBackend(t, backend)

	info, err := backend.Info(apiKey)
	if err != nil {
		t.Error(err)
	}

	assert.Equal("this is foo.", info)
}

func TestBackend_GetAge(t *testing.T) {
	// assert := assert.New(t)

	backend, err := NewBackend(dbURN)
	if err != nil {
		t.Fatalf("cannot create backend: %s", err)
	}

	backend.InitDB()
	defer backend.DropDB()
}
