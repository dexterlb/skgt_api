package backend

import (
	"fmt"

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
