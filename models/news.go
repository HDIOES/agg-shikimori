package models

import (
	"database/sql"

	"github.com/pkg/errors"
)

//NewDAO struct
type NewDAO struct {
	Db *sql.DB
}

//DeleteAll function
func (dao *NewDAO) DeleteAll() error {
	tx, beginErr := dao.Db.Begin()
	if beginErr != nil {
		return rollbackTransaction(tx, errors.Wrap(beginErr, ""))
	}
	stmt, prepareStmtErr := tx.Prepare("TRUNCATE new CASCADE")
	if prepareStmtErr != nil {
		return rollbackTransaction(tx, errors.Wrap(prepareStmtErr, ""))
	}
	defer stmt.Close()
	if _, stmtErr := stmt.Exec(); stmtErr != nil {
		return rollbackTransaction(tx, errors.Wrap(stmtErr, ""))
	}
	if cErr := commitTransaction(tx); cErr != nil {
		return errors.Wrap(cErr, "")
	}
	return nil
}

//Find function
func (dao *NewDAO) Find(id int64) (*NewDTO, error) {
	stmt, prepareStmtErr := dao.Db.Prepare("SELECT id, name, body FROM new WHERE id = $1")
	if prepareStmtErr != nil {
		return nil, errors.Wrap(prepareStmtErr, "")
	}
	defer stmt.Close()
	result, stmtErr := stmt.Query(id)
	if stmtErr != nil {
		return nil, errors.Wrap(stmtErr, "")
	}
	defer result.Close()
	newDTO := &NewDTO{}
	if result.Next() {
		var sqlID sql.NullInt64
		var sqlName sql.NullString
		var sqlBody sql.NullString
		result.Scan(&sqlID, &sqlName, &sqlBody)
		newDTO.ID = &sqlID.Int64
		newDTO.Name = &sqlName.String
		newDTO.Body = &sqlBody.String
	} else {
		return nil, errors.New("New not found")
	}
	return newDTO, nil
}

//Create function
func (dao *NewDAO) Create(dto NewDTO) (int64, error) {
	tx, beginErr := dao.Db.Begin()
	if beginErr != nil {
		return 0, rollbackTransaction(tx, errors.Wrap(beginErr, ""))
	}
	stmt, prepareStmtErr := tx.Prepare("INSERT INTO new (name, body) VALUES($1, $2) RETURNING id")
	if prepareStmtErr != nil {
		return 0, rollbackTransaction(tx, errors.Wrap(prepareStmtErr, ""))
	}
	defer stmt.Close()
	result, stmtErr := stmt.Query(dto.Name, dto.Body)
	if stmtErr != nil {
		return 0, rollbackTransaction(tx, errors.Wrap(stmtErr, ""))
	}
	var lastID sql.NullInt64
	if result.Next() {
		result.Scan(&lastID)
	}
	result.Close()
	if commitErr := tx.Commit(); commitErr != nil {
		return 0, rollbackTransaction(tx, errors.Wrap(commitErr, ""))
	}
	return lastID.Int64, nil
}

//Update function
func (dao *NewDAO) Update(dto NewDTO) error {
	tx, beginErr := dao.Db.Begin()
	if beginErr != nil {
		return rollbackTransaction(tx, errors.Wrap(beginErr, ""))
	}
	stmt, prepareStmtErr := tx.Prepare("UPDATE new (name, body) VALUES($2, $3) WHERE id = $1")
	if prepareStmtErr != nil {
		return rollbackTransaction(tx, errors.Wrap(prepareStmtErr, ""))
	}
	defer stmt.Close()
	_, stmtErr := stmt.Exec(dto.ID, dto.Name, dto.Body)
	if stmtErr != nil {
		return rollbackTransaction(tx, errors.Wrap(stmtErr, ""))
	}
	if commitErr := tx.Commit(); commitErr != nil {
		return rollbackTransaction(tx, errors.Wrap(commitErr, ""))
	}
	return nil
}

//Delete function
func (dao *NewDAO) Delete(id int64) error {
	tx, beginErr := dao.Db.Begin()
	if beginErr != nil {
		return rollbackTransaction(tx, errors.Wrap(beginErr, ""))
	}
	stmt, prepareStmtErr := tx.Prepare("DELETE FROM new WHERE id = $1")
	if prepareStmtErr != nil {
		return rollbackTransaction(tx, errors.Wrap(prepareStmtErr, ""))
	}
	defer stmt.Close()
	_, stmtErr := stmt.Exec(id)
	if stmtErr != nil {
		return rollbackTransaction(tx, errors.Wrap(stmtErr, ""))
	}
	if commitErr := tx.Commit(); commitErr != nil {
		return rollbackTransaction(tx, errors.Wrap(commitErr, ""))
	}
	return nil
}

//NewDTO struct
type NewDTO struct {
	ID   *int64
	Name *string
	Body *string
}
