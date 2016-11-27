package backend

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

// Wrap wraps a function with queries to the database in transaction
func (b *Backend) Wrap(f func(tx *sqlx.Tx) (interface{}, error)) (data interface{}, err error) {
	tx, err := b.db.Beginx()
	if err != nil {
		return nil, fmt.Errorf("cannot open transaction: %s", err)
	}

	defer func() {
		// if something has failed, and tx hasn't been commited, roll it back
		if tx != nil {
			txErr := tx.Rollback()
			if txErr != nil {
				err = fmt.Errorf("trying to rollback transaction because of error [%s] failed: %s", err, txErr)
			}
		}
	}()

	// do the actual work
	data, err = f(tx)
	if err != nil {
		return nil, fmt.Errorf("cannot exec function in transaction: %s", err)
	}

	// commit transaction
	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("cannot commit transaction: %s", err)
	}

	// set tx to nil so that the deferred error check knows not to rollback
	tx = nil

	return data, nil
}
