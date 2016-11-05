package backend

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
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

func openBackend(t *testing.T) *Backend {
	backend, err := New(dbURN)
	if err != nil {
		t.Fatalf("cannot create backend: %s", err)
	}

	err = backend.InitDB()
	if err != nil {
		t.Errorf("cannot create database: %s", err)
	}

	return backend
}

func closeBackend(t *testing.T, backend *Backend) {
	err := backend.DropDB()
	if err != nil {
		t.Errorf("cannot drop database: %s", err)
	}
}

func pause() {
	sigs := make(chan os.Signal, 1)
	fmt.Printf("Press CTRL-C to continue...\n")
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
}

func TestBackend_CheckAPIKey(t *testing.T) {
	backend := openBackend(t)
	defer closeBackend(t, backend)

	apiKey, err := backend.NewAPIKey()
	if err != nil {
		t.Fatalf("cannot create api key: %s", err)
	}

	err = backend.CheckAPIKey(apiKey)
	if err != nil {
		t.Fatalf("generated api key is wrong: %s", err)
	}

	err = backend.CheckAPIKey("42")
	if err != ErrWrongAPIKey {
		t.Fatalf("wrong error for wrong api key: %s", err)
	}
}

func TestBackend_Info(t *testing.T) {
	assert := assert.New(t)
	backend := openBackend(t)
	defer closeBackend(t, backend)

	info, err := backend.Info()
	if err != nil {
		t.Error(err)
	}

	assert.Equal("skgt-api, pre-release", info)
}

func TestBackend_GetAge(t *testing.T) {
	// assert := assert.New(t)

	backend, err := New(dbURN)
	if err != nil {
		t.Fatalf("cannot create backend: %s", err)
	}

	backend.InitDB()
	defer backend.DropDB()
}
