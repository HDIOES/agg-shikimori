package rest

import (
	"database/sql"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/HDIOES/cpa-backend/models"
	"github.com/HDIOES/cpa-backend/rest/util"
	"github.com/gorilla/mux"
)

//CreateSearchAnimeHandler function
func CreateSearchAnimeHandler(db *sql.DB, router *mux.Router, config *util.Configuration) http.Handler {
	animeDao := models.AnimeDAO{Db: db}
	searchAnimeHandler := &SearchAnimeHandler{Dao: &animeDao}
	return searchAnimeHandler
}

//SearchAnimeHandler struct
type SearchAnimeHandler struct {
	Dao *models.AnimeDAO
}

func (as *SearchAnimeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
	if limit, limitOk := vars["limit"]; limitOk {
		limitInt64, parseErr := strconv.ParseInt(limit[0], 10, 32)
		if parseErr != nil {
			HandleErr(parseErr, w, 400, "Limit not valid")
			return
		}
		animeSQLBuilder.SetLimit(int32(limitInt64))
	}
	if offset, offsetOk := vars["offset"]; offsetOk {
		offsetInt64, parseErr := strconv.ParseInt(offset[0], 10, 32)
		if parseErr != nil {
			HandleErr(parseErr, w, 400, "Offset not valid")
			return
		}
		animeSQLBuilder.SetLimit(int32(offsetInt64))
	}
	animeDtos, err := as.Dao.FindByFilter(animeSQLBuilder)
	if err != nil {
		HandleErr(err, w, 400, "Error")
		return
	}
	animeRos := []AnimeRO{}
	for _, animeDto := range animeDtos {
		animeRo := AnimeRO{Name: animeDto.Name, RussuanName: animeDto.Russian, URL: animeDto.AnimeURL, PosterURL: animeDto.PosterURL}
		animeRos = append(animeRos, animeRo)
	}
	ReturnResponseAsJSON(w, animeRos, 200)
}

//AnimeRO is rest object
type AnimeRO struct {
	Name        string `json:"name"`
	RussuanName string `json:"russian_name"`
	URL         string `json:"url"`
	PosterURL   string `json:"poster_url"`
}
