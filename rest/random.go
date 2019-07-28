package rest

import (
	"database/sql"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/HDIOES/cpa-backend/rest/util"
)

//CreateRandomAnimeHandler function receive handler for rest-method /animes/random
func CreateRandomAnimeHandler(db *sql.DB, config *util.Configuration) http.Handler {
	animeDao := AnimeDao{Db: db, Config: config}
	randomAnimeHandler := &RandomAnimeHandler{Dao: &animeDao}
	return randomAnimeHandler
}

//RandomAnimeHandler struct
type RandomAnimeHandler struct {
	Dao *AnimeDao
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
	if countOfAnimes, err := rah.Dao.GetCount(animeSQLBuilder); err != nil {
		//TODO error processing
	} else {
		animeSQLBuilder.SetCountOnly(false)
		animeSQLBuilder.SetRowNumber(rand.Int63n(countOfAnimes + 1))
	}
	if animeRO, err := rah.Dao.GetRandomAnime(animeSQLBuilder); err != nil {
		//TODO error processing
	} else {
		json.NewEncoder(w).Encode(animeRO)
	}
}
