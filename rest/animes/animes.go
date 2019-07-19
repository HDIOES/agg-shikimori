package animes

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/HDIOES/cpa-backend/util"
	"github.com/gorilla/mux"
)

func CreateSearchAnimeHandler(db *sql.DB, router *mux.Router, config *util.Configuration) http.Handler {
	animeDao := AnimeDao{Db: db, Config: config}
	searchAnimeHandler := &SearchAnimeHandler{Dao: &animeDao}
	return searchAnimeHandler
}

type SearchAnimeHandler struct {
	Dao *AnimeDao
}

func (as *SearchAnimeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars, parseErr := url.ParseQuery(r.URL.RawQuery)
	if parseErr != nil {
		log.Println(parseErr)
	}
	animeSQLBuilder := AnimeQueryBuilder{}
	if status, statusOk := vars["status"]; statusOk {
		animeSQLBuilder.SetStatus(status[0])
	}
	if kind, kindOk := vars["kind"]; kindOk {
		animeSQLBuilder.SetKind(kind[0])
	}
	if phrase, phraseOk := vars["phrase"]; phraseOk {
		animeSQLBuilder.SetPhrase(phrase[0])
	}
	if order, orderOK := vars["order"]; orderOK {
		animeSQLBuilder.SetOrder(order[0])
	}
	if score, scoreOk := vars["score"]; scoreOk {
		if scoreInt64, parseErr := strconv.ParseInt(score[0], 10, 32); parseErr != nil {
			//TODO error processing
		} else {
			animeSQLBuilder.SetScore(int32(scoreInt64))
		}
	}
	if genre, genreOk := vars["genre"]; genreOk {
		if scoreInt64, parseErr := strconv.ParseInt(genre[0], 10, 32); parseErr != nil {
			//TODO error processing
		} else {
			animeSQLBuilder.SetScore(int32(scoreInt64))
		}
	}
	if studio, studioOk := vars["studio"]; studioOk {
		if studioInt64, parseErr := strconv.ParseInt(studio[0], 10, 64); parseErr != nil {
			//TODO error processing
		} else {
			animeSQLBuilder.AddStudioID(studioInt64)
		}
	}
	if duration, durationOk := vars["duration"]; durationOk {
		animeSQLBuilder.SetDuration(duration[0])
	}
	if rating, ratingOk := vars["rating"]; ratingOk {
		animeSQLBuilder.SetRating(rating[0])
	}
	if franchise, franchiseOk := vars["franchise"]; franchiseOk {
		animeSQLBuilder.SetFranchise(franchise[0])
	}
	if ids, idsOk := vars["ids"]; idsOk {
		for _, id := range strings.Split(ids[0], " ") {
			if idInt64, parseErr := strconv.ParseInt(id, 10, 64); parseErr != nil {
				//TODO error processing
			} else {
				animeSQLBuilder.AddId(idInt64)
			}
		}
	}
	if excludeIds, excludeIdsOk := vars["exclude_ids"]; excludeIdsOk {
		for _, id := range strings.Split(excludeIds[0], " ") {
			if excludeIDInt64, parseErr := strconv.ParseInt(id, 10, 64); parseErr != nil {
				//TODO error processing
			} else {
				animeSQLBuilder.AddExcludeId(excludeIDInt64)
			}
		}
	}

	if limit, limitOk := vars["limit"]; limitOk {
		if limitInt64, parseErr := strconv.ParseInt(limit[0], 10, 32); parseErr != nil {
			//TODO error processing
		} else {
			animeSQLBuilder.SetLimit(int32(limitInt64))
		}
	}
	if offset, offsetOk := vars["offset"]; offsetOk {
		if offsetInt64, parseErr := strconv.ParseInt(offset[0], 10, 32); parseErr != nil {
			//TODO error processing
		} else {
			animeSQLBuilder.SetLimit(int32(offsetInt64))
		}
	}
	if animes, err := as.Dao.SearchAnimes(animeSQLBuilder); err != nil {
		//TODO error processing
	} else {
		json.NewEncoder(w).Encode(animes)
	}
}

type AnimeDao struct {
	Db     *sql.DB
	Config *util.Configuration
}

//SearchAnimes function returns anime array by predefined filter
func (a *AnimeDao) SearchAnimes(sqlBuilder AnimeQueryBuilder) (animes *[]AnimeRO, err error) {
	animes = &[]AnimeRO{}
	sqlQuery, args := sqlBuilder.Build()
	result, queryErr := a.Db.Query(sqlQuery, args...)
	if queryErr != nil {
		return nil, queryErr
	}
	defer result.Close()
	for result.Next() {
		animeRo := AnimeRO{}
		var rowNumber sql.NullInt64
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
		scanErr := result.Scan(
			&rowNumber,
			&id,
			&name,
			&externalID,
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
		if scanErr != nil {
			return nil, scanErr
		}
		animeRo.Name = name.String
		animeRo.RussuanName = russianName.String
		animeRo.URL = a.Config.ShikimoriURL + animeURL.String
		animeRo.PosterURL = a.Config.ShikimoriURL + posterURL.String
		*animes = append(*animes, animeRo)
	}
	return animes, err
}

func (a *AnimeDao) GetCount(sqlBuilder AnimeQueryBuilder) (int64, error) {
	sqlQuery, args := sqlBuilder.Build()
	result, queryErr := a.Db.Query(sqlQuery, args...)
	if queryErr != nil {
		return 0, queryErr
	}
	defer result.Close()
	if result.Next() {
		var count sql.NullInt64
		result.Scan(&count)
		return count.Int64, nil
	} else {
		return 0, nil
	}
}

//function GetRandomAnime returns random anime by predefined filter
func (a *AnimeDao) GetRandomAnime(sqlBuilder AnimeQueryBuilder) (*AnimeRO, error) {
	sqlQuery, args := sqlBuilder.Build()
	result, queryErr := a.Db.Query(sqlQuery, args...)
	if queryErr != nil {
		return nil, queryErr
	}
	defer result.Close()
	if result.Next() {
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
		scanErr := result.Scan(&id,
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
		if scanErr != nil {
			return nil, scanErr
		}
		animeRo.Name = name.String
		animeRo.RussuanName = russianName.String
		animeRo.URL = a.Config.ShikimoriURL + animeURL.String
		animeRo.PosterURL = a.Config.ShikimoriURL + posterURL.String
		return &animeRo, nil
	} else {
		return nil, nil
	}
}

//AnimeQueryBuilder struct
type AnimeQueryBuilder struct {
	Limit      int32
	Offset     int32
	Status     string
	Kind       string
	Phrase     string
	Order      string
	Score      int32
	GenreIds   []int64
	StudioIds  []int64
	Duration   string
	Rating     string
	Franchise  string
	Ids        []int64
	ExcludeIds []int64
	SQLQuery   strings.Builder
	CountOnly  bool
	RowNumber  int64
}

func (aqb *AnimeQueryBuilder) AddExcludeId(excludeId int64) {
	aqb.ExcludeIds = append(aqb.ExcludeIds, excludeId)
}

func (aqb *AnimeQueryBuilder) AddId(id int64) {
	aqb.Ids = append(aqb.Ids, id)
}

func (aqb *AnimeQueryBuilder) SetFranchise(franchise string) {
	aqb.Franchise = franchise
}

func (aqb *AnimeQueryBuilder) SetRating(rating string) {
	aqb.Rating = rating
}

func (aqb *AnimeQueryBuilder) SetDuration(duration string) {
	aqb.Duration = duration
}

func (aqb *AnimeQueryBuilder) AddStudioID(studioID int64) {
	aqb.StudioIds = append(aqb.StudioIds, studioID)
}

func (aqb *AnimeQueryBuilder) AddGenreId(genreID int64) {
	aqb.GenreIds = append(aqb.GenreIds)
}

func (aqb *AnimeQueryBuilder) SetScore(score int32) {
	aqb.Score = score
}

func (aqb *AnimeQueryBuilder) SetOrder(order string) {
	aqb.Order = order
}

func (aqb *AnimeQueryBuilder) SetPhrase(phrase string) {
	aqb.Phrase = phrase
}

func (aqb *AnimeQueryBuilder) SetKind(kind string) {
	aqb.Kind = kind
}

func (aqb *AnimeQueryBuilder) SetStatus(status string) {
	aqb.Status = status
}

//SetLimit func
func (aqb *AnimeQueryBuilder) SetLimit(limit int32) {
	aqb.Limit = limit
}

//SetOffset func
func (aqb *AnimeQueryBuilder) SetOffset(offset int32) {
	aqb.Offset = offset
}

func (aqb *AnimeQueryBuilder) SetCountOnly(countOnly bool) {
	aqb.CountOnly = countOnly
}

func (aqb *AnimeQueryBuilder) SetRowNumber(rowNumber int64) {
	aqb.RowNumber = rowNumber
}

//Build func
func (aqb *AnimeQueryBuilder) Build() (string, []interface{}) {
	if aqb.RowNumber > 0 && !aqb.CountOnly {
		aqb.SQLQuery.WriteString("SELECT ")
		aqb.SQLQuery.WriteString("query.anime_internal_id,")
		aqb.SQLQuery.WriteString("query.name,")
		aqb.SQLQuery.WriteString("query.anime_external_id,")
		aqb.SQLQuery.WriteString("query.russian,")
		aqb.SQLQuery.WriteString("query.amine_url,")
		aqb.SQLQuery.WriteString("query.kind,")
		aqb.SQLQuery.WriteString("query.anime_status,")
		aqb.SQLQuery.WriteString("query.epizodes,")
		aqb.SQLQuery.WriteString("query.epizodes_aired,")
		aqb.SQLQuery.WriteString("query.aired_on,")
		aqb.SQLQuery.WriteString("query.released_on,")
		aqb.SQLQuery.WriteString("query.poster_url,")
		aqb.SQLQuery.WriteString("query.score,")
		aqb.SQLQuery.WriteString("query.duration,")
		aqb.SQLQuery.WriteString("query.rating,")
		aqb.SQLQuery.WriteString("query.franchase ")
		aqb.SQLQuery.WriteString("FROM (")
	}
	if aqb.RowNumber == 0 && aqb.CountOnly {
		aqb.SQLQuery.WriteString("SELECT COUNT(*) FROM (")
	}
	aqb.SQLQuery.WriteString("SELECT ")
	aqb.SQLQuery.WriteString("row_number() over(),")
	aqb.SQLQuery.WriteString("animes.anime_internal_id,")
	aqb.SQLQuery.WriteString("animes.name,")
	aqb.SQLQuery.WriteString("animes.anime_external_id,")
	aqb.SQLQuery.WriteString("animes.russian,")
	aqb.SQLQuery.WriteString("animes.amine_url,")
	aqb.SQLQuery.WriteString("animes.kind,")
	aqb.SQLQuery.WriteString("animes.anime_status,")
	aqb.SQLQuery.WriteString("animes.epizodes,")
	aqb.SQLQuery.WriteString("animes.epizodes_aired,")
	aqb.SQLQuery.WriteString("animes.aired_on,")
	aqb.SQLQuery.WriteString("animes.released_on,")
	aqb.SQLQuery.WriteString("animes.poster_url,")
	aqb.SQLQuery.WriteString("animes.score,")
	aqb.SQLQuery.WriteString("animes.duration,")
	aqb.SQLQuery.WriteString("animes.rating,")
	aqb.SQLQuery.WriteString("animes.franchase ")
	aqb.SQLQuery.WriteString("FROM ")
	aqb.SQLQuery.WriteString("(")
	aqb.SQLQuery.WriteString("SELECT ")
	aqb.SQLQuery.WriteString("anime.id AS anime_internal_id,")
	aqb.SQLQuery.WriteString("anime.name,")
	aqb.SQLQuery.WriteString("anime.external_id as anime_external_id,")
	aqb.SQLQuery.WriteString("anime.russian,")
	aqb.SQLQuery.WriteString("anime.amine_url,")
	aqb.SQLQuery.WriteString("anime.kind,")
	aqb.SQLQuery.WriteString("anime.anime_status,")
	aqb.SQLQuery.WriteString("anime.epizodes,")
	aqb.SQLQuery.WriteString("anime.epizodes_aired,")
	aqb.SQLQuery.WriteString("anime.aired_on,")
	aqb.SQLQuery.WriteString("anime.released_on,")
	aqb.SQLQuery.WriteString("anime.poster_url,")
	aqb.SQLQuery.WriteString("anime.score,")
	aqb.SQLQuery.WriteString("anime.duration,")
	aqb.SQLQuery.WriteString("anime.rating,")
	aqb.SQLQuery.WriteString("anime.franchase ")
	countOfParameter := 0
	args := make([]interface{}, 0)
	if len(aqb.GenreIds) > 0 {
		aqb.SQLQuery.WriteString(", genre.external_id as genre_external_id")
	}
	if len(aqb.StudioIds) > 0 {
		aqb.SQLQuery.WriteString(", studio.external_id as studio_external_id")
	}
	if len(aqb.Phrase) > 0 {
		countOfParameter++
		aqb.SQLQuery.WriteString(", to_tsvector(anime.russian) as russian_tsvector, to_tsvector(anime.name) as english_tsvector, phraseto_tsquery($")
		aqb.SQLQuery.WriteString(strconv.Itoa(countOfParameter))
		aqb.SQLQuery.WriteString(") as ts_query")
		args = append(args, aqb.Phrase)
	}
	aqb.SQLQuery.WriteString(" FROM anime")
	if len(aqb.GenreIds) > 0 {
		aqb.SQLQuery.WriteString(" JOIN anime_genre ON anime.id = anime_genre.anime_id")
		aqb.SQLQuery.WriteString(" JOIN genre ON genre.id = anime_genre.genre_id")
	}
	if len(aqb.StudioIds) > 0 {
		aqb.SQLQuery.WriteString(" JOIN anime_studio ON anime.id = anime_studio.anime_id")
		aqb.SQLQuery.WriteString(" JOIN studio ON studio.id = anime_studio.studio_id")
	}
	aqb.SQLQuery.WriteString(") as animes")
	aqb.SQLQuery.WriteString(" WHERE 1=1")
	if len(aqb.Phrase) > 0 {
		aqb.SQLQuery.WriteString(" AND (animes.russian_tsvector @@ animes.ts_query OR animes.english_tsvector @@ animes.ts_query)")
	}
	if len(aqb.GenreIds) > 0 {
		aqb.SQLQuery.WriteString(" AND genre_external_id IN (")
		for ind, genreExternalID := range aqb.GenreIds {
			countOfParameter++
			args = append(args, genreExternalID)
			aqb.SQLQuery.WriteString("$")
			aqb.SQLQuery.WriteString(strconv.Itoa(countOfParameter))
			if ind < len(aqb.GenreIds)-1 {
				aqb.SQLQuery.WriteString(",")
			}
		}
		aqb.SQLQuery.WriteString(")")
	}
	if len(aqb.StudioIds) > 0 {
		aqb.SQLQuery.WriteString(" AND studio_external_id IN (")
		for ind, studioExternalID := range aqb.StudioIds {
			countOfParameter++
			args = append(args, studioExternalID)
			aqb.SQLQuery.WriteString("$")
			aqb.SQLQuery.WriteString(strconv.Itoa(countOfParameter))
			if ind < len(aqb.StudioIds)-1 {
				aqb.SQLQuery.WriteString(",")
			}
		}
		aqb.SQLQuery.WriteString(")")
	}
	if len(aqb.Status) > 0 {
		countOfParameter++
		aqb.SQLQuery.WriteString(" AND animes.anime_status = $")
		aqb.SQLQuery.WriteString(strconv.Itoa(countOfParameter))
		args = append(args, aqb.Status)
	}
	if len(aqb.Kind) > 0 {
		var kinds = [...]string{"tv", "movie", "ova", "ona", "special", "music", "tv_13", "tv_24", "tv_48"}
		for _, s := range kinds {
			if strings.Compare(s, aqb.Kind) == 0 {
				countOfParameter++
				aqb.SQLQuery.WriteString(" AND animes.kind = $")
				aqb.SQLQuery.WriteString(strconv.Itoa(countOfParameter))
				args = append(args, aqb.Kind)
				break
			}
		}
	}
	if len(aqb.Ids) > 0 {
		aqb.SQLQuery.WriteString(" AND anime_external_id IN (")
		for ind, id := range aqb.Ids {
			countOfParameter++
			args = append(args, id)
			aqb.SQLQuery.WriteString("$")
			aqb.SQLQuery.WriteString(strconv.Itoa(countOfParameter))
			if ind < len(aqb.Ids)-1 {
				aqb.SQLQuery.WriteString(",")
			}
		}
		aqb.SQLQuery.WriteString(")")
	}
	if len(aqb.ExcludeIds) > 0 {
		aqb.SQLQuery.WriteString(" AND anime_external_id NOT IN (")
		aqb.SQLQuery.WriteString(strconv.Itoa(countOfParameter))
		for ind, excludeID := range aqb.ExcludeIds {
			countOfParameter++
			args = append(args, excludeID)
			aqb.SQLQuery.WriteString("$")
			aqb.SQLQuery.WriteString(strconv.Itoa(countOfParameter))
			if ind < len(aqb.ExcludeIds)-1 {
				aqb.SQLQuery.WriteString(",")
			}
		}
		aqb.SQLQuery.WriteString(")")
	}
	if len(aqb.Duration) > 0 {
		switch aqb.Duration {
		case "S":
			{
				aqb.SQLQuery.WriteString(" AND animes.duration < 10")
			}
		case "D":
			{
				aqb.SQLQuery.WriteString(" AND animes.duration < 30")
			}
		case "F":
			{
				aqb.SQLQuery.WriteString(" AND animes.duration >= 30")
			}
		}
	}
	if len(aqb.Franchise) > 0 {
		countOfParameter++
		aqb.SQLQuery.WriteString(" AND animes.franchase = $")
		aqb.SQLQuery.WriteString(strconv.Itoa(countOfParameter))
		args = append(args, aqb.Franchise)
	}
	if len(aqb.Rating) > 0 {
		var ratings = [...]string{"none", "g", "pg", "pg_13", "r", "r_plus", "rx"}
		for _, r := range ratings {
			if strings.Compare(r, aqb.Rating) == 0 {
				countOfParameter++
				aqb.SQLQuery.WriteString(" AND animes.rating = $")
				aqb.SQLQuery.WriteString(strconv.Itoa(countOfParameter))
				args = append(args, aqb.Rating)
				break
			}
		}
	}
	if aqb.Score > 0.0 {
		//need to validate score
		countOfParameter++
		aqb.SQLQuery.WriteString(" AND animes.score >= $")
		aqb.SQLQuery.WriteString(strconv.Itoa(countOfParameter))
		args = append(args, aqb.Score)
	}
	if len(aqb.Order) > 0 {
		aqb.SQLQuery.WriteString(" ORDER BY ")
		switch aqb.Order {
		case "id":
			{
				countOfParameter++
				aqb.SQLQuery.WriteString("$")
				aqb.SQLQuery.WriteString(strconv.Itoa(countOfParameter))
				args = append(args, "anime_external_id")
			}
		case "kind":
			{
				countOfParameter++
				aqb.SQLQuery.WriteString("$")
				aqb.SQLQuery.WriteString(strconv.Itoa(countOfParameter))
				args = append(args, "animes.kind")
			}
		case "name":
			{
				countOfParameter++
				aqb.SQLQuery.WriteString("$")
				aqb.SQLQuery.WriteString(strconv.Itoa(countOfParameter))
				args = append(args, "animes.name")
			}
		case "aired_on":
			{
				countOfParameter++
				aqb.SQLQuery.WriteString("$")
				aqb.SQLQuery.WriteString(strconv.Itoa(countOfParameter))
				args = append(args, "animes.aired_on")
			}
		case "episodes":
			{
				countOfParameter++
				aqb.SQLQuery.WriteString("$")
				aqb.SQLQuery.WriteString(strconv.Itoa(countOfParameter))
				args = append(args, "animes.epizodes")
			}
		case "status":
			{
				countOfParameter++
				aqb.SQLQuery.WriteString("$")
				aqb.SQLQuery.WriteString(strconv.Itoa(countOfParameter))
				args = append(args, "animes.status")
			}
		case "relevance":
			{
				if len(aqb.Phrase) > 0 {
					countOfParameter++
					aqb.SQLQuery.WriteString("$")
					aqb.SQLQuery.WriteString(strconv.Itoa(countOfParameter))
					args = append(args, "get_rank(animes.russian_tsvector, animes.english_tsvector, animes.ts_query) DESC")
				}
			}
		}
	}
	if aqb.CountOnly && aqb.RowNumber == 0 {
		aqb.SQLQuery.WriteString(") as query")
	} else if !aqb.CountOnly && aqb.RowNumber > 0 {
		countOfParameter++
		aqb.SQLQuery.WriteString(") as query where query.row_number = $")
		aqb.SQLQuery.WriteString(strconv.Itoa(countOfParameter))
		args = append(args, aqb.RowNumber)
	} else {
		if aqb.Limit > 0 {
			countOfParameter++
			aqb.SQLQuery.WriteString(" LIMIT $" + strconv.Itoa(countOfParameter))
			args = append(args, aqb.Limit)
		} else {
			aqb.SQLQuery.WriteString(" LIMIT 50")
		}
		if aqb.Offset > 0 {
			countOfParameter++
			aqb.SQLQuery.WriteString(" OFFSET $")
			aqb.SQLQuery.WriteString(strconv.Itoa(countOfParameter))
			args = append(args, aqb.Offset)
		}
	}
	return aqb.SQLQuery.String(), args
}

//AnimeRO is rest object
type AnimeRO struct {
	Name        string `json:"name"`
	RussuanName string `json:"russian_name"`
	URL         string `json:"url"`
	PosterURL   string `json:"poster_url"`
	Ld          int
}
