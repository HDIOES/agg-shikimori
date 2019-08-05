package models

import (
	"database/sql"
	"errors"
	"strconv"
	"strings"
)

//StudioDAO struct
type StudioDAO struct {
	Db *sql.DB
}

//FindByID function
func (dao *StudioDAO) FindByID(id int64) (*StudioDTO, error) {
	stmt, stmtErr := dao.Db.Prepare("SELECT ")
	if stmtErr != nil {
		return nil, stmtErr
	}
	defer stmt.Close()
	result, resultErr := stmt.Query(id)
	if resultErr != nil {
		return nil, resultErr
	}
	defer result.Close()
	if result.Next() {
		dto := StudioDTO{}
		var ID sql.NullInt64
		var externalID sql.NullString
		var name sql.NullString
		var filteredName sql.NullString
		var isReal sql.NullBool
		var imageURL sql.NullString
		result.Scan(&ID, &externalID, &name, &filteredName, &isReal, &imageURL)
		dto.ID = ID.Int64
		dto.ExternalID = externalID.String
		dto.Name = name.String
		dto.FilteredStudioName = filteredName.String
		dto.IsReal = isReal.Bool
		dto.ImageURL = imageURL.String
		return &dto, nil
	}
	return nil, errors.New("Studio not found")
}

//FindByExternalID function
func (dao *StudioDAO) FindByExternalID(externalID string) (*StudioDTO, error) {
	sqlBuilder := StudioQueryBuilder{}
	sqlBuilder.SetExternalID(externalID)
	studios, err := dao.FindByFilter(sqlBuilder)
	if err != nil {
		return nil, err
	}
	if len(studios) > 0 {
		return &studios[0], nil
	}
	return nil, errors.New("Studio not found")
}

//FindByFilter function
func (dao *StudioDAO) FindByFilter(sqlBuilder StudioQueryBuilder) ([]StudioDTO, error) {
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
	dtos := []StudioDTO{}
	for result.Next() {
		var ID sql.NullInt64
		var externalID sql.NullString
		var name sql.NullString
		var filteredName sql.NullString
		var isReal sql.NullBool
		var imageURL sql.NullString
		result.Scan(&ID, &externalID, &name, &filteredName, &isReal, &imageURL)
		dto := StudioDTO{ID: ID.Int64, ExternalID: externalID.String, Name: name.String, FilteredStudioName: filteredName.String, IsReal: isReal.Bool, ImageURL: imageURL.String}
		dtos = append(dtos, dto)
	}
	return dtos, nil
}

//Create function
func (dao *StudioDAO) Create(dto StudioDTO) (int64, error) {
	tx, txErr := dao.Db.Begin()
	if txErr != nil {
		return 0, txErr
	}
	stmt, stmtErr := tx.Prepare("INSERT INTO studio (external_id, studio_name, filtered_studio_name, is_real, image_url) VALUES($1, $2, $3, $4, $5)")
	if stmtErr != nil {
		return 0, rollbackTransaction(tx, stmtErr)
	}
	defer stmt.Close()
	result, resultErr := stmt.Query(dto.ExternalID, dto.Name, dto.FilteredStudioName, dto.IsReal, dto.ImageURL)
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
func (dao *StudioDAO) Update(dto StudioDTO) error {
	tx, txErr := dao.Db.Begin()
	if txErr != nil {
		return txErr
	}
	stmt, stmtErr := tx.Prepare("UPDATE studio SET external_id = $1, SET studio_name = $2, SET filtered_studio_name = $3, SET is_real = $4, SET image_url = $5 WHERE id = $6")
	if stmtErr != nil {
		return rollbackTransaction(tx, stmtErr)
	}
	defer stmt.Close()
	_, resultErr := stmt.Exec(dto.ExternalID, dto.Name, dto.FilteredStudioName, dto.IsReal, dto.ImageURL, dto.ID)
	if resultErr != nil {
		return rollbackTransaction(tx, resultErr)
	}
	if commitErr := tx.Commit(); commitErr != nil {
		return rollbackTransaction(tx, commitErr)
	}
	return nil
}

//StudioDTO struct
type StudioDTO struct {
	ID                 int64
	ExternalID         string
	Name               string
	FilteredStudioName string
	IsReal             bool
	ImageURL           string
}

//StudioQueryBuilder struct
type StudioQueryBuilder struct {
	Limit      int32
	Offset     int32
	ExternalID string
}

//Build func
func (sqb *StudioQueryBuilder) Build() (string, []interface{}) {
	query := strings.Builder{}
	args := make([]interface{}, 0)
	query.WriteString("SELECT external_id, studio_name, filtered_studio_name, kind FROM studio WHERE 1=1")
	countOfParameter := 0
	if len(sqb.ExternalID) > 0 {
		countOfParameter++
		args = append(args, sqb.ExternalID)
		query.WriteString(" AND $")
		query.WriteString(strconv.Itoa(countOfParameter))
	}
	if sqb.Limit > 0 {
		countOfParameter++
		args = append(args, sqb.Limit)
		query.WriteString(" LIMIT $")
		query.WriteString(strconv.Itoa(countOfParameter))
	} else {
		query.WriteString(" LIMIT 50")
	}
	if sqb.Offset > 0 {
		countOfParameter++
		args = append(args, sqb.Offset)
		query.WriteString(" OFFSET $")
		query.WriteString(strconv.Itoa(countOfParameter))
	}
	return query.String(), args
}

//SetLimit func
func (sqb *StudioQueryBuilder) SetLimit(limit int32) {
	sqb.Limit = limit
}

//SetOffset func
func (sqb *StudioQueryBuilder) SetOffset(offset int32) {
	sqb.Offset = offset
}

//SetExternalID func
func (sqb *StudioQueryBuilder) SetExternalID(externalID string) {
	sqb.ExternalID = externalID
}
