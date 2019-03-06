package genres

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strconv"
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
	limit, limitOk := vars["limit"]
	offset, offsetOk := vars["offset"]
	sqlQueryString := "SELECT external_id, genre_name, russian, kind FROM genre WHERE 1=1"
	countOfParameter := 0
	args := make([]interface{}, 0)
	if offsetOk {
		args = append(args, offset[0])
		countOfParameter++
		sqlQueryString += " OFFSET $" + strconv.Itoa(countOfParameter)
	}
	if limitOk {
		countOfParameter++
		args = append(args, limit[0])
		sqlQueryString += " LIMIT $" + strconv.Itoa(countOfParameter)
	} else {
		sqlQueryString += " LIMIT 50"
	}
	rows, rowsErr := g.Db.Query(sqlQueryString, args...)
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

type GenreRo struct {
	ID      *string `json:"id"`
	Name    *string `json:"name"`
	Russian *string `json:"russian"`
	Kind    *string `json:"kind"`
}
