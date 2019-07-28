package rest

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func CreateGenreHandler(db *sql.DB) http.Handler {
	genreHandler := &GenreHandler{Db: db}
	return genreHandler
}

type GenreHandler struct {
	Db *sql.DB
}

func (g *GenreHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars, parseErr := url.ParseQuery(r.URL.RawQuery)
	if parseErr != nil {
		log.Println(parseErr)
	}
	genreSQLBuilder := GenreQueryBuilder{}
	if limit, limitOk := vars["limit"]; limitOk {
		if limitInt64, parseErr := strconv.ParseInt(limit[0], 10, 32); parseErr != nil {
			//TODO error processing
		} else {
			genreSQLBuilder.SetOffset(int32(limitInt64))
		}
	}
	if offset, offsetOk := vars["offset"]; offsetOk {
		if offsetInt64, parseErr := strconv.ParseInt(offset[0], 10, 32); parseErr != nil {
			//TODO error processing
		} else {
			genreSQLBuilder.SetOffset(int32(offsetInt64))
		}
	}
	sqlQuery, args := genreSQLBuilder.Build()
	rows, rowsErr := g.Db.Query(sqlQuery, args...)
	if rowsErr != nil {
		log.Println(rowsErr)
	}
	defer rows.Close()
	genres := []GenreRo{}
	for rows.Next() {
		genreRo := GenreRo{}
		var id sql.NullString
		var name sql.NullString
		var russian sql.NullString
		var kind sql.NullString
		rows.Scan(&id, &name, &russian, &kind)
		genreRo.ID = &id.String
		genreRo.Name = &name.String
		genreRo.Russian = &russian.String
		genreRo.Kind = &kind.String
		genres = append(genres, genreRo)
	}
	json.NewEncoder(w).Encode(genres)
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

type GenreRo struct {
	ID      *string `json:"id"`
	Name    *string `json:"name"`
	Russian *string `json:"russian"`
	Kind    *string `json:"kind"`
}
