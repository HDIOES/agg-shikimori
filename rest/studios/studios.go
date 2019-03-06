package studios

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

func CreateStudioHandler(db *sql.DB) http.Handler {
	studioHandler := &StudioHandler{Db: db}
	return studioHandler
}

type StudioHandler struct {
	Db *sql.DB
}

func (g *StudioHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars, parseErr := url.ParseQuery(r.URL.RawQuery)
	if parseErr != nil {
		log.Println(parseErr)
	}
	limit, limitOk := vars["limit"]
	offset, offsetOk := vars["offset"]
	sqlQueryString := "SELECT external_id, studio_name, filtered_studio_name FROM studio WHERE 1=1"
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
	studios := []StudioRo{}
	for rows.Next() {
		studioRo := StudioRo{}
		var id sql.NullString
		var name sql.NullString
		var filteredName sql.NullString
		rows.Scan(&id, &name, &filteredName)
		studioRo.ID = &id.String
		studioRo.Name = &name.String
		studioRo.FilteredName = &filteredName.String
		studios = append(studios, studioRo)
	}
	json.NewEncoder(w).Encode(studios)
}

type StudioRo struct {
	ID           *string `json:"id"`
	Name         *string `json:"name"`
	FilteredName *string `json:"filtered_name"`
}
