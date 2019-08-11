package models

import "database/sql"

func commitTransaction(tx *sql.Tx) error {
	if commitErr := tx.Commit(); commitErr != nil {
		return rollbackTransaction(tx, commitErr)
	}
	return nil
}

func rollbackTransaction(tx *sql.Tx, err error) error {
	if rollErr := tx.Rollback(); rollErr != nil {
		return rollErr
	}
	return err
}
