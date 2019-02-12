package animes

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
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

//AnimeRO is rest object
type AnimeRO struct {
	Name      string `json:"name"`
	URL       string `json:"url"`
	PosterURL string `json:"poster_url"`
}
