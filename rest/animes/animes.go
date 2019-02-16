package animes

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gorilla/mux"
)

func CreateAnimeHandler(db *sql.DB) http.Handler {
	animeHandler := &AnimeHandler{Db: db}
	return animeHandler
}

type AnimeHandler struct {
	Db *sql.DB
}

func (a *AnimeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rows, queryErr := a.Db.Query("SELECT COUNT(*) FROM anime")
	if queryErr != nil {
		fmt.Println(queryErr)
	}
	defer rows.Close()
	var count sql.NullInt64
	if rows.Next() {
		err := rows.Scan(&count)
		if err != nil {
			fmt.Println(err)
		}
	}
	randowRowNumber := rand.Int63n(count.Int64) + 1
	animeRows, animeRowsErr := a.Db.Query("select russian, amine_url, poster_url from (select row_number() over(), russian, amine_url, poster_url from anime) as query where query.row_number = $1", randowRowNumber)
	if animeRowsErr != nil {
		fmt.Println(animeRowsErr)
	}
	defer animeRows.Close()
	animeRo := &AnimeRO{}
	if animeRows.Next() {
		var russianName sql.NullString
		var animeURL sql.NullString
		var posterURL sql.NullString
		animeRows.Scan(&russianName, &animeURL, &posterURL)
		animeRo.Name = russianName.String
		animeRo.URL = "https://shikimori.org" + animeURL.String
		animeRo.PosterURL = "https://shikimori.org" + posterURL.String
	}
	json.NewEncoder(w).Encode(animeRo)
}

func CreateSearchAnimeHandler(db *sql.DB, router *mux.Router) http.Handler {
	searchAnimeHandler := &SearchAnimeHandler{Db: db, Router: router}
	return searchAnimeHandler
}

type SearchAnimeHandler struct {
	Db     *sql.DB
	Router *mux.Router
}

func (as *SearchAnimeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars, parseErr := url.ParseQuery(r.URL.RawQuery)
	if parseErr != nil {
		fmt.Println(parseErr)
	}
	status, statusOk := vars["status"]
	limit, limitOk := vars["limit"]
	offset, offsetOk := vars["offset"]
	animes := []AnimeRO{}
	args := make([]interface{}, 0)
	sqlQueryString := "SELECT russian, amine_url, poster_url FROM anime WHERE 1=1"
	countOfParameter := 0
	if statusOk {
		countOfParameter++
		sqlQueryString += " AND status = $" + strconv.Itoa(countOfParameter)
		args = append(args, status[0])
	}
	if limitOk {
		countOfParameter++
		sqlQueryString += " LIMIT $" + strconv.Itoa(countOfParameter)
		value, err := strconv.ParseInt(limit[0], 10, 0)
		if err != nil {
			fmt.Println(err)
		}
		args = append(args, value)
	} else {
		sqlQueryString += " LIMIT 50"
	}
	if offsetOk {
		countOfParameter++
		sqlQueryString += " OFFSET $" + strconv.Itoa(countOfParameter)
		args = append(args, offset[0])
	}
	result, queryErr := as.Db.Query(sqlQueryString, args...)
	if queryErr != nil {
		fmt.Println(queryErr)
		panic(queryErr)
	}
	defer result.Close()
	for result.Next() {
		animeRo := AnimeRO{}
		var russianName sql.NullString
		var animeURL sql.NullString
		var posterURL sql.NullString
		result.Scan(&russianName, &animeURL, &posterURL)
		animeRo.Name = russianName.String
		animeRo.URL = "https://shikimori.org" + animeURL.String
		animeRo.PosterURL = "https://shikimori.org" + posterURL.String
		animes = append(animes, animeRo)
	}
	json.NewEncoder(w).Encode(animes)
}

//AnimeRO is rest object
type AnimeRO struct {
	Name      string `json:"name"`
	URL       string `json:"url"`
	PosterURL string `json:"poster_url"`
}
