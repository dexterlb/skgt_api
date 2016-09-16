package backend

import (
	"flag"
	"fmt"
	"testing"

	"github.com/jmoiron/sqlx"
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

func TestBackend_Foo(t *testing.T) {
	assert := assert.New(t)

	backend, err := NewBackend(dbURN)
	if err != nil {
		t.Fatalf("cannot create backend: %s", err)
	}

	foo, err := backend.Foo()
	if err != nil {
		t.Error(err)
	}

	assert.Equal("this is foo.", foo)
}

func TestBackend_GetAge(t *testing.T) {
	assert := assert.New(t)

	backend, err := NewBackend(dbURN)
	if err != nil {
		t.Fatalf("cannot create backend: %s", err)
	}

	backend.InitDB()
	defer backend.DropDB()

	db := sqlx.MustConnect("postgres", dbURN)
	db.MustExec(`insert into person(name, age) values('pesho', 42)`)

	age, err := backend.GetAge("pesho")
	if err != nil {
		t.Error(err)
	}

	assert.Equal(42, age)
}
