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
)

//CreateRandomAnimeHandler function receive handler for rest-method /animes/random
func CreateRandomAnimeHandler(db *sql.DB, config util.Configuration) http.Handler {
	randomAnimeHandler := &RandomAnimeHandler{Db: db, Config: config}
	return randomAnimeHandler
}

//RandomAnimeHandler struct
type RandomAnimeHandler struct {
	Db     *sql.DB
	Config util.Configuration
}

func (rah *RandomAnimeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
	animeSQLBuilder.SetCountOnly(true)
	animeSQLBuilder.SetRowNumber(0)
	countOfAnimes := rah.getCount(animeSQLBuilder)
	animeSQLBuilder.SetCountOnly(false)
	animeSQLBuilder.SetRowNumber(countOfAnimes + 1)
	animeRO := rah.getRandomAnime(animeSQLBuilder)
	json.NewEncoder(w).Encode(animeRO)
}

func (rah *RandomAnimeHandler) getCount(sqlBuilder AnimeQueryBuilder) int64 {
	sqlQuery, args := sqlBuilder.Build()
	result, queryErr := rah.Db.Query(sqlQuery, args...)
	if queryErr != nil {
		log.Println(queryErr)
		panic(queryErr)
	}
	defer result.Close()
	if result.Next() {
		var count sql.NullInt64
		result.Scan(&count)
		return count.Int64
	}
	return 0
}

func (rah *RandomAnimeHandler) getRandomAnime(sqlBuilder AnimeQueryBuilder) *AnimeRO {
	sqlQuery, args := sqlBuilder.Build()
	result, queryErr := rah.Db.Query(sqlQuery, args...)
	if queryErr != nil {
		log.Println(queryErr)
		panic(queryErr)
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
		animeRo.URL = rah.Config.ShikimoriURL + animeURL.String
		animeRo.PosterURL = rah.Config.ShikimoriURL + posterURL.String
		return &animeRo
	} else {
		return nil
	}
}
