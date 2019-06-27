package animes

import (
	"database/sql"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"

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
		log.Println(queryErr)
	}
	defer rows.Close()
	var count sql.NullInt64
	if rows.Next() {
		err := rows.Scan(&count)
		if err != nil {
			log.Println(err)
		}
	}
	randowRowNumber := rand.Int63n(count.Int64) + 1
	animeRows, animeRowsErr := a.Db.Query("select russian, amine_url, poster_url from (select row_number() over(), russian, amine_url, poster_url from anime) as query where query.row_number = $1", randowRowNumber)
	if animeRowsErr != nil {
		log.Println(animeRowsErr)
	}
	defer animeRows.Close()
	animeRo := &AnimeRO{}
	if animeRows.Next() {
		var russianName sql.NullString
		var animeURL sql.NullString
		var posterURL sql.NullString
		animeRows.Scan(&russianName, &animeURL, &posterURL)
		animeRo.Name = russianName.String
		animeRo.URL = "https://shikimori.one" + animeURL.String
		animeRo.PosterURL = "https://shikimori.one" + posterURL.String
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
		log.Println(parseErr)
	}
	status, statusOk := vars["status"]
	kind, kindOk := vars["kind"]
	phrase, phraseOk := vars["phrase"]
	order, orderOK := vars["order"]
	score, scoreOk := vars["score"]
	genre, genreOk := vars["genre"]
	studio, studioOk := vars["studio"]
	duration, durationOk := vars["duration"]
	rating, ratingOk := vars["rating"]
	franchise, franchiseOk := vars["franchise"]
	ids, idsOk := vars["ids"]
	excludeIds, excludeIdsOk := vars["exclude_ids"]

	limit, limitOk := vars["limit"]
	offset, offsetOk := vars["offset"]
	animes := []AnimeRO{}
	args := make([]interface{}, 0)

	sqlQueryString := "SELECT animes.anime_internal_id," +
		"animes.name," +
		"animes.anime_external_id," +
		"animes.russian," +
		"animes.amine_url," +
		"animes.kind," +
		"animes.anime_status," +
		"animes.epizodes," +
		"animes.epizodes_aired," +
		"animes.aired_on," +
		"animes.released_on," +
		"animes.poster_url," +
		"animes.score," +
		"animes.duration," +
		"animes.rating," +
		"animes.franchase " +
		"FROM (SELECT " +
		"anime.id AS anime_internal_id," +
		"anime.name," +
		"anime.external_id as anime_external_id," +
		"anime.russian," +
		"anime.amine_url," +
		"anime.kind," +
		"anime.anime_status," +
		"anime.epizodes," +
		"anime.epizodes_aired," +
		"anime.aired_on," +
		"anime.released_on," +
		"anime.poster_url," +
		"anime.score," +
		"anime.duration," +
		"anime.rating," +
		"anime.franchase "
	countOfParameter := 0
	if genreOk {
		sqlQueryString += ", genre.external_id as genre_external_id"
	}
	if studioOk {
		sqlQueryString += ", studio.external_id as studio_external_id"
	}
	if phraseOk {
		countOfParameter++
		sqlQueryString += ", to_tsvector(anime.russian) as russian_tsvector, to_tsvector(anime.name) as english_tsvector, phraseto_tsquery($" + strconv.Itoa(countOfParameter) + ") as ts_query"
		args = append(args, phrase[0])
	}
	sqlQueryString += " FROM anime"
	if genreOk {
		sqlQueryString += " JOIN anime_genre ON anime.id = anime_genre.anime_id" +
			" JOIN genre ON genre.id = anime_genre.genre_id"
	}
	if studioOk {
		sqlQueryString += " JOIN anime_studio ON anime.id = anime_studio.anime_id" +
			" JOIN studio ON studio.id = anime_studio.studio_id"
	}
	sqlQueryString += ") as animes"
	sqlQueryString += " WHERE 1=1"
	if phraseOk {
		sqlQueryString += " AND (animes.russian_tsvector @@ animes.ts_query OR animes.english_tsvector @@ animes.ts_query)"
	}
	if genreOk {
		countOfParameter++
		sqlQueryString += " AND genre_external_id IN ($" + strconv.Itoa(countOfParameter)
		var params = strings.Split(genre[0], ",")
		args = append(args, params[0])
		for ind, genreExternalID := range params {
			if ind == 0 {
				continue
			} else if ind == len(params)-1 {
				countOfParameter++
				sqlQueryString += ", $" + strconv.Itoa(countOfParameter)
				args = append(args, genreExternalID)
			} else {
				countOfParameter++
				sqlQueryString += ", $" + strconv.Itoa(countOfParameter) + ", "
				args = append(args, genreExternalID)
			}
		}
		sqlQueryString += ")"
	}
	if studioOk {
		countOfParameter++
		sqlQueryString += " AND studio_external_id IN ($" + strconv.Itoa(countOfParameter)
		var params = strings.Split(studio[0], ",")
		args = append(args, params[0])
		for ind, studioExternalID := range params {
			if ind == 0 {
				continue
			} else if ind == len(params)-1 {
				countOfParameter++
				sqlQueryString += ", $" + strconv.Itoa(countOfParameter)
				args = append(args, studioExternalID)
			} else {
				countOfParameter++
				sqlQueryString += ", $" + strconv.Itoa(countOfParameter) + ", "
				args = append(args, studioExternalID)
			}
		}
		sqlQueryString += ")"
	}
	if statusOk {
		countOfParameter++
		sqlQueryString += " AND animes.anime_status = $" + strconv.Itoa(countOfParameter)
		args = append(args, status[0])
	}
	if kindOk {
		var kinds = [...]string{"tv", "movie", "ova", "ona", "special", "music", "tv_13", "tv_24", "tv_48"}
		for _, s := range kinds {
			if s == kind[0] {
				countOfParameter++
				sqlQueryString += " AND animes.kind = $" + strconv.Itoa(countOfParameter)
				args = append(args, kind[0])
				break
			}
		}
	}
	if idsOk {
		countOfParameter++
		sqlQueryString += " AND anime_external_id IN ($" + strconv.Itoa(countOfParameter)
		var params = strings.Split(ids[0], ",")
		args = append(args, params[0])
		for ind, id := range params {
			if ind == 0 {
				continue
			} else if ind == len(params)-1 {
				countOfParameter++
				sqlQueryString += ", $" + strconv.Itoa(countOfParameter)
				args = append(args, id)
			} else {
				countOfParameter++
				sqlQueryString += ", $" + strconv.Itoa(countOfParameter) + ", "
				args = append(args, id)
			}
		}
		sqlQueryString += ")"
	}
	//log.Panicln("query = " + sqlQueryString)
	if excludeIdsOk {
		countOfParameter++
		sqlQueryString += " AND anime_external_id NOT IN ($" + strconv.Itoa(countOfParameter)
		var params = strings.Split(excludeIds[0], ",")
		args = append(args, params[0])
		for ind, excludeID := range params {
			if ind == 0 {
				continue
			} else if ind == len(params)-1 {
				countOfParameter++
				sqlQueryString += ", $" + strconv.Itoa(countOfParameter)
				args = append(args, excludeID)
			} else {
				countOfParameter++
				sqlQueryString += ", $" + strconv.Itoa(countOfParameter) + ", "
				args = append(args, excludeID)
			}
		}
		sqlQueryString += ")"
	}
	if durationOk {
		switch duration[0] {
		case "S":
			{
				sqlQueryString += " AND animes.duration < 10"
			}
		case "D":
			{
				sqlQueryString += " AND animes.duration < 30"
			}
		case "F":
			{
				sqlQueryString += " AND animes.duration >= 30"
			}
		}
	}
	if franchiseOk {
		countOfParameter++
		sqlQueryString += " AND animes.franchase = $" + strconv.Itoa(countOfParameter)
		args = append(args, franchise[0])
	}
	if ratingOk {
		var ratings = [...]string{"none", "g", "pg", "pg_13", "r", "r_plus", "rx"}
		for _, r := range ratings {
			if r == rating[0] {
				countOfParameter++
				sqlQueryString += " AND animes.rating = $" + strconv.Itoa(countOfParameter)
				args = append(args, rating[0])
				break
			}
		}
	}
	if scoreOk {
		//need to validate score
		countOfParameter++
		sqlQueryString += " AND animes.score >= $" + strconv.Itoa(countOfParameter)
		args = append(args, score[0])
	}
	if orderOK {
		sqlQueryString += " ORDER BY "
		switch order[0] {
		case "id":
			{
				countOfParameter++
				sqlQueryString += "$" + strconv.Itoa(countOfParameter)
				args = append(args, "anime_external_id")
			}
		case "kind":
			{
				countOfParameter++
				sqlQueryString += "$" + strconv.Itoa(countOfParameter)
				args = append(args, "animes.kind")
			}
		case "name":
			{
				countOfParameter++
				sqlQueryString += "$" + strconv.Itoa(countOfParameter)
				args = append(args, "animes.name")
			}
		case "aired_on":
			{
				countOfParameter++
				sqlQueryString += "$" + strconv.Itoa(countOfParameter)
				args = append(args, "animes.aired_on")
			}
		case "episodes":
			{
				countOfParameter++
				sqlQueryString += "$" + strconv.Itoa(countOfParameter)
				args = append(args, "animes.epizodes")
			}
		case "status":
			{
				countOfParameter++
				sqlQueryString += "$" + strconv.Itoa(countOfParameter)
				args = append(args, "animes.status")
			}
		case "relevance":
			{
				if phraseOk {
					countOfParameter++
					sqlQueryString += "$" + strconv.Itoa(countOfParameter)
					args = append(args, "get_rank(animes.russian_tsvector, animes.english_tsvector, animes.ts_query) DESC")
				}
			}
		}
	}
	if limitOk {
		countOfParameter++
		sqlQueryString += " LIMIT $" + strconv.Itoa(countOfParameter)
		value, err := strconv.ParseInt(limit[0], 10, 0)
		if err != nil {
			log.Println(err)
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
	log.Println(sqlQueryString)
	result, queryErr := as.Db.Query(sqlQueryString, args...)
	if queryErr != nil {
		log.Println(queryErr)
		panic(queryErr)
	}
	defer result.Close()
	for result.Next() {
		animeRo := AnimeRO{}
		var id sql.NullInt64
		var name sql.NullString
		var externalID sql.NullString
		var russianName sql.NullString
		var animeURL sql.NullString
		var kind sql.NullString
		var animeStatus sql.NullString
		var epizodes sql.NullInt64
		var epizodesAired sql.NullInt64
		var airedOn sql.NullString
		var releasedOn sql.NullString
		var posterURL sql.NullString
		var score sql.NullFloat64
		var duration sql.NullFloat64
		var rating sql.NullString
		var franchase sql.NullString
		result.Scan(&id,
			&name, &externalID,
			&russianName,
			&animeURL,
			&kind,
			&animeStatus,
			&epizodes,
			&epizodesAired,
			&airedOn,
			&releasedOn,
			&posterURL,
			&score,
			&duration,
			&rating,
			&franchase)
		animeRo.Name = name.String
		animeRo.RussuanName = russianName.String
		animeRo.URL = "https://shikimori.org" + animeURL.String
		animeRo.PosterURL = "https://shikimori.org" + posterURL.String
		animes = append(animes, animeRo)
	}
	json.NewEncoder(w).Encode(animes)
}

//LowensteinDistance copypasting from wikipedia
func LowensteinDistance(s1, s2 string) int {
	min := func(values ...int) int {
		m := values[0]
		for _, v := range values {
			if v < m {
				m = v
			}
		}
		return m
	}
	r1, r2 := []rune(s1), []rune(s2)
	n, m := len(r1), len(r2)
	if n > m {
		r1, r2 = r2, r1
		n, m = m, n
	}
	currentRow := make([]int, n+1)
	previousRow := make([]int, n+1)
	for i := range currentRow {
		currentRow[i] = i
	}
	for i := 1; i <= m; i++ {
		for j := range currentRow {
			previousRow[j] = currentRow[j]
			if j == 0 {
				currentRow[j] = i
				continue
			} else {
				currentRow[j] = 0
			}
			add, del, change := previousRow[j]+1, currentRow[j-1]+1, previousRow[j-1]
			if r1[j-1] != r2[i-1] {
				change++
			}
			currentRow[j] = min(add, del, change)
		}
	}
	return currentRow[n]
}

//AnimeRO is rest object
type AnimeRO struct {
	Name        string `json:"name"`
	RussuanName string `json:"russian_name"`
	URL         string `json:"url"`
	PosterURL   string `json:"poster_url"`
	Ld          int
}
