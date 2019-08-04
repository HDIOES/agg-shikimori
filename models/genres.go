package models

import (
	"database/sql"
	"errors"
	"strconv"
	"strings"
)

//GenreDAO struct
type GenreDAO struct {
	Db sql.DB
}

//FindByFilter function
func (dao *GenreDAO) FindByFilter(sqlBuilder GenreQueryBuilder) ([]GenreDTO, error) {
	query, args := sqlBuilder.Build()
	stmt, stmtErr := dao.Db.Prepare(query)
	if stmtErr != nil {
		return nil, stmtErr
	}
	defer stmt.Close()
	result, resultErr := stmt.Query(args)
	if resultErr != nil {
		return nil, resultErr
	}
	defer result.Close()
	dtos := []GenreDTO{}
	for result.Next() {
		var ID sql.NullInt64
		var externalID sql.NullString
		var name sql.NullString
		var russian sql.NullString
		var kind sql.NullString
		result.Scan(&ID, &externalID, &name, &russian, &kind)
		dto := GenreDTO{ID: ID.Int64, ExternalID: externalID.String, Name: name.String, Russian: russian.String, Kind: kind.String}
		dtos = append(dtos, dto)
	}
	return dtos, nil
}

//Create function
func (dao *GenreDAO) Create(dto GenreDTO) (int64, error) {
	tx, txErr := dao.Db.Begin()
	if txErr != nil {
		return 0, txErr
	}
	stmt, stmtErr := tx.Prepare("INSERT INTO genre (external_id, genre_name, russian, kind) VALUES($1, $2, $3, $4)")
	if stmtErr != nil {
		return 0, rollbackTransaction(tx, stmtErr)
	}
	defer stmt.Close()
	result, resultErr := stmt.Query(dto.ExternalID, dto.Name, dto.Russian, dto.Kind)
	if resultErr != nil {
		return 0, rollbackTransaction(tx, resultErr)
	}
	var ID sql.NullInt64
	if result.Next() {
		result.Scan(&ID)
	} else {
		result.Close()
		return 0, rollbackTransaction(tx, errors.New("Genre not found"))
	}
	result.Close()
	if commitErr := tx.Commit(); commitErr != nil {
		return 0, rollbackTransaction(tx, commitErr)
	}
	return ID.Int64, nil
}

//Update function
func (dao *GenreDAO) Update(dto GenreDTO) error {
	tx, txErr := dao.Db.Begin()
	if txErr != nil {
		return txErr
	}
	stmt, stmtErr := tx.Prepare("UPDATE genre SET external_id = $1, SET genre_name = $2, SET russian = $3, SET kind = $4 WHERE id = $5")
	if stmtErr != nil {
		return rollbackTransaction(tx, stmtErr)
	}
	defer stmt.Close()
	_, resultErr := stmt.Exec(dto.ExternalID, dto.Name, dto.Russian, dto.Kind, dto.ID)
	if resultErr != nil {
		return rollbackTransaction(tx, resultErr)
	}
	if commitErr := tx.Commit(); commitErr != nil {
		return rollbackTransaction(tx, commitErr)
	}
	return nil
}

//GenreDTO struct
type GenreDTO struct {
	ID         int64
	ExternalID string
	Name       string
	Russian    string
	Kind       string
}

//GenreQueryBuilder struct
type GenreQueryBuilder struct {
	Limit  int32
	Offset int32
}

//Build func
func (gqb *GenreQueryBuilder) Build() (string, []interface{}) {
	query := strings.Builder{}
	args := make([]interface{}, 0)
	query.WriteString("SELECT external_id, genre_name, russian, kind FROM genre WHERE 1=1")
	countOfParameter := 0
	if gqb.Limit > 0 {
		countOfParameter++
		args = append(args, gqb.Limit)
		query.WriteString(" LIMIT $")
		query.WriteString(strconv.Itoa(countOfParameter))
	} else {
		query.WriteString(" LIMIT 50")
	}
	if gqb.Offset > 0 {
		countOfParameter++
		args = append(args, gqb.Offset)
		query.WriteString(" OFFSET $")
		query.WriteString(strconv.Itoa(countOfParameter))
	}
	return query.String(), args
}

//SetLimit func
func (gqb *GenreQueryBuilder) SetLimit(limit int32) {
	gqb.Limit = limit
}

//SetOffset func
func (gqb *GenreQueryBuilder) SetOffset(offset int32) {
	gqb.Offset = offset
}
