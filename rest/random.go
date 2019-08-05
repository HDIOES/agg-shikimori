package rest

import (
	"database/sql"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/HDIOES/cpa-backend/models"
	"github.com/HDIOES/cpa-backend/rest/util"
)

//CreateRandomAnimeHandler function receive handler for rest-method /animes/random
func CreateRandomAnimeHandler(db *sql.DB, config *util.Configuration) http.Handler {
	animeDao := models.AnimeDAO{Db: db}
	randomAnimeHandler := &RandomAnimeHandler{Dao: &animeDao}
	return randomAnimeHandler
}

//RandomAnimeHandler struct
type RandomAnimeHandler struct {
	Dao *models.AnimeDAO
}

func (rah *RandomAnimeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars, parseErr := url.ParseQuery(r.URL.RawQuery)
	if parseErr != nil {
		log.Println(parseErr)
	}
	animeSQLBuilder := models.AnimeQueryBuilder{}
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
		scoreInt64, parseErr := strconv.ParseInt(score[0], 10, 32)
		if parseErr != nil {
			HandleErr(parseErr, w, 400, "Score not valid")
			return
		}
		animeSQLBuilder.SetScore(int32(scoreInt64))
	}
	if genre, genreOk := vars["genre"]; genreOk {
		for _, genreID := range strings.Split(genre[0], ",") {
			animeSQLBuilder.AddGenreID(genreID)
		}
	}
	if studio, studioOk := vars["studio"]; studioOk {
		for _, studioID := range strings.Split(studio[0], ",") {
			animeSQLBuilder.AddStudioID(studioID)
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
		for _, id := range strings.Split(ids[0], ",") {
			animeSQLBuilder.AddID(id)
		}
	}
	if excludeIds, excludeIdsOk := vars["exclude_ids"]; excludeIdsOk {
		for _, id := range strings.Split(excludeIds[0], ",") {
			animeSQLBuilder.AddExcludeID(id)
		}
	}
	animeSQLBuilder.SetCountOnly(true)
	animeSQLBuilder.SetRowNumber(0)
	countOfAnimes, err := rah.Dao.GetCount(animeSQLBuilder)
	if err != nil {
		HandleErr(parseErr, w, 400, "Internal error")
		return
	}
	animeSQLBuilder.SetCountOnly(false)
	animeSQLBuilder.SetRowNumber(rand.Int63n(countOfAnimes + 1))
	animeRO, err := rah.Dao.GetRandomAnime(animeSQLBuilder)
	if err != nil {
		HandleErr(parseErr, w, 400, "Internal error")
		return
	}
	ReturnResponseAsJSON(w, animeRO, 200)
}
