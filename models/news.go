package models

import (
	"database/sql"
	"errors"
)

//NewDAO struct
type NewDAO struct {
	Db *sql.DB
}

//Find function
func (ndao *NewDAO) Find(id int64) (*NewDTO, error) {
	stmt, prepareStmtErr := ndao.Db.Prepare("SELECT id, name, body FROM new WHERE id = $1")
	if prepareStmtErr != nil {
		return nil, prepareStmtErr
	}
	defer stmt.Close()
	result, stmtErr := stmt.Query(id)
	if stmtErr != nil {
		return nil, stmtErr
	}
	defer result.Close()
	newDTO := &NewDTO{}
	if result.Next() {
		var sqlID sql.NullInt64
		var sqlName sql.NullString
		var sqlBody sql.NullString
		result.Scan(&sqlID, &sqlName, &sqlBody)
		newDTO.ID = sqlID.Int64
		newDTO.Name = sqlName.String
		newDTO.Body = sqlBody.String
	} else {
		return nil, errors.New("New not found")
	}
	return newDTO, nil
}

//Create function
func (ndao *NewDAO) Create(dto NewDTO) (int64, error) {
	tx, beginErr := ndao.Db.Begin()
	if beginErr != nil {
		return 0, rollbackTransaction(tx, beginErr)
	}
	stmt, prepareStmtErr := tx.Prepare("INSERT INTO new (name, body) VALUES($1, $2) RETURNING id")
	if prepareStmtErr != nil {
		return 0, rollbackTransaction(tx, prepareStmtErr)
	}
	defer stmt.Close()
	result, stmtErr := stmt.Query(dto.Name, dto.Body)
	if stmtErr != nil {
		return 0, rollbackTransaction(tx, stmtErr)
	}
	var lastID sql.NullInt64
	if result.Next() {
		result.Scan(&lastID)
	}
	result.Close()
	if commitErr := tx.Commit(); commitErr != nil {
		return 0, rollbackTransaction(tx, commitErr)
	}
	return lastID.Int64, nil
}

//Update function
func (ndao *NewDAO) Update(dto NewDTO) error {
	tx, beginErr := ndao.Db.Begin()
	if beginErr != nil {
		return rollbackTransaction(tx, beginErr)
	}
	stmt, prepareStmtErr := tx.Prepare("UPDATE new (name, body) VALUES($2, $3) WHERE id = $1")
	if prepareStmtErr != nil {
		return rollbackTransaction(tx, prepareStmtErr)
	}
	defer stmt.Close()
	_, stmtErr := stmt.Exec(dto.ID, dto.Name, dto.Body)
	if stmtErr != nil {
		return rollbackTransaction(tx, stmtErr)
	}
	if commitErr := tx.Commit(); commitErr != nil {
		return rollbackTransaction(tx, commitErr)
	}
	return nil
}

//Delete function
func (ndao *NewDAO) Delete(id int64) error {
	tx, beginErr := ndao.Db.Begin()
	if beginErr != nil {
		return rollbackTransaction(tx, beginErr)
	}
	stmt, prepareStmtErr := tx.Prepare("DELETE FROM new WHERE id = $1")
	if prepareStmtErr != nil {
		return rollbackTransaction(tx, prepareStmtErr)
	}
	_, stmtErr := stmt.Exec(id)
	if stmtErr != nil {
		stmt.Close()
		return rollbackTransaction(tx, stmtErr)
	}
	stmt.Close()
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

//NewDTO struct
type NewDTO struct {
	ID   int64
	Name string
	Body string
}
