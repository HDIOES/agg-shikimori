package models

import (
	"database/sql"

	"github.com/pkg/errors"
)

func commitTransaction(tx *sql.Tx) error {
	if commitErr := tx.Commit(); commitErr != nil {
		return rollbackTransaction(tx, errors.Wrap(commitErr, ""))
	}
	return nil
}

func rollbackTransaction(tx *sql.Tx, err error) error {
	if rollErr := tx.Rollback(); rollErr != nil {
		return errors.Wrap(rollErr, "")
	}
	return errors.Wrap(err, "")
}
