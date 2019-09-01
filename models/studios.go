package models

import (
	"database/sql"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

//StudioDAO struct
type StudioDAO struct {
	Db *sql.DB
}

//DeleteAll function
func (dao *StudioDAO) DeleteAll() error {
	tx, beginErr := dao.Db.Begin()
	if beginErr != nil {
		return rollbackTransaction(tx, errors.Wrap(beginErr, ""))
	}
	stmt, prepareStmtErr := tx.Prepare("TRUNCATE studio CASCADE")
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

//FindByID function
func (dao *StudioDAO) FindByID(id int64) (*StudioDTO, error) {
	stmt, stmtErr := dao.Db.Prepare("SELECT ")
	if stmtErr != nil {
		return nil, errors.Wrap(stmtErr, "")
	}
	defer stmt.Close()
	result, resultErr := stmt.Query(id)
	if resultErr != nil {
		return nil, errors.Wrap(resultErr, "")
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
		if externalID.Valid {
			dto.ExternalID = externalID.String
		}
		if name.Valid {
			dto.Name = &name.String
		}
		if filteredName.Valid {
			dto.FilteredStudioName = &filteredName.String
		}
		if isReal.Valid {
			dto.IsReal = &isReal.Bool
		}
		if imageURL.Valid {
			dto.ImageURL = &imageURL.String
		}
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
		return nil, errors.Wrap(err, "")
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
		return nil, errors.Wrap(stmtErr, "")
	}
	defer stmt.Close()
	result, resultErr := stmt.Query(args...)
	if resultErr != nil {
		return nil, errors.Wrap(resultErr, "")
	}
	defer result.Close()
	var capacity int32 = 50
	if sqlBuilder.Limit > 0 && sqlBuilder.Limit <= 50 {
		capacity = sqlBuilder.Limit
	}
	dtos := make([]StudioDTO, 0, capacity)
	for result.Next() {
		var ID sql.NullInt64
		var externalID sql.NullString
		var name sql.NullString
		var filteredName sql.NullString
		var isReal sql.NullBool
		var imageURL sql.NullString
		result.Scan(&ID, &externalID, &name, &filteredName, &isReal, &imageURL)
		dto := StudioDTO{}
		if ID.Valid {
			dto.ID = ID.Int64
		}
		if externalID.Valid {
			dto.ExternalID = externalID.String
		}
		if name.Valid {
			dto.Name = &name.String
		}
		if filteredName.Valid {
			dto.FilteredStudioName = &filteredName.String
		}
		if isReal.Valid {
			dto.IsReal = &isReal.Bool
		}
		if imageURL.Valid {
			dto.ImageURL = &imageURL.String
		}
		dtos = append(dtos, dto)
	}
	return dtos, nil
}

//Create function
func (dao *StudioDAO) Create(dto StudioDTO) (int64, error) {
	tx, txErr := dao.Db.Begin()
	if txErr != nil {
		return 0, errors.Wrap(txErr, "")
	}
	stmt, stmtErr := tx.Prepare("INSERT INTO studio (external_id, studio_name, filtered_studio_name, is_real, image_url) VALUES($1, $2, $3, $4, $5) RETURNING id")
	if stmtErr != nil {
		time.Now()
		return 0, rollbackTransaction(tx, errors.Wrap(stmtErr, ""))
	}
	defer stmt.Close()
	result, resultErr := stmt.Query(dto.ExternalID, dto.Name, dto.FilteredStudioName, dto.IsReal, dto.ImageURL)
	if resultErr != nil {
		return 0, rollbackTransaction(tx, errors.Wrap(resultErr, ""))
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
		return 0, rollbackTransaction(tx, errors.Wrap(commitErr, ""))
	}
	return ID.Int64, nil
}

//Update function
func (dao *StudioDAO) Update(dto StudioDTO) error {
	tx, txErr := dao.Db.Begin()
	if txErr != nil {
		return errors.Wrap(txErr, "")
	}
	stmt, stmtErr := tx.Prepare("UPDATE studio SET external_id = $1, studio_name = $2, filtered_studio_name = $3, is_real = $4, image_url = $5 WHERE id = $6")
	if stmtErr != nil {
		return rollbackTransaction(tx, errors.Wrap(stmtErr, ""))
	}
	defer stmt.Close()
	_, resultErr := stmt.Exec(dto.ExternalID, dto.Name, dto.FilteredStudioName, dto.IsReal, dto.ImageURL, dto.ID)
	if resultErr != nil {
		return rollbackTransaction(tx, errors.Wrap(resultErr, ""))
	}
	if commitErr := tx.Commit(); commitErr != nil {
		return rollbackTransaction(tx, errors.Wrap(commitErr, ""))
	}
	return nil
}

//StudioDTO struct
type StudioDTO struct {
	ID                 int64
	ExternalID         string
	Name               *string
	FilteredStudioName *string
	IsReal             *bool
	ImageURL           *string
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
	query.WriteString("SELECT id, external_id, studio_name, filtered_studio_name, is_real, image_url FROM studio WHERE 1=1")
	countOfParameter := 0
	if len(sqb.ExternalID) > 0 {
		countOfParameter++
		args = append(args, sqb.ExternalID)
		query.WriteString(" AND external_id = $")
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
