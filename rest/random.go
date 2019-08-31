package rest

import (
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/HDIOES/cpa-backend/models"
	"github.com/HDIOES/cpa-backend/rest/util"
	"github.com/pkg/errors"
)

//RandomAnimeHandler struct
type RandomAnimeHandler struct {
	Dao           *models.AnimeDAO
	Configuration *util.Configuration
}

func (rah *RandomAnimeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	requestBody, rawQuery, headers, err := GetRequestData(r)
	if err != nil {
		HandleErr(errors.Wrap(err, ""), w, 400, "Request cannot be read")
		return
	}
	if err := LogHTTPRequest(r.URL.String(), headers, requestBody); err != nil {
		HandleErr(errors.Wrap(err, ""), w, 400, "Request cannot be logged")
		return
	}
	vars, parseErr := url.ParseQuery(*rawQuery)
	if parseErr != nil {
		HandleErr(errors.Wrap(parseErr, ""), w, 400, "URL not valid")
		return
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
			HandleErr(errors.Wrap(parseErr, ""), w, 400, "Score not valid")
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
		HandleErr(errors.Wrap(err, ""), w, 400, "Internal error")
		return
	}
	animeSQLBuilder.SetCountOnly(false)
	animeSQLBuilder.SetRowNumber(rand.Int63n(countOfAnimes + 1))
	animeDto, err := rah.Dao.GetRandomAnime(animeSQLBuilder)
	if err != nil {
		HandleErr(errors.Wrap(err, ""), w, 400, "Internal error")
		return
	}
	animeRo := AnimeRO{ShikiID: animeDto.ExternalID, Name: animeDto.Name, RussuanName: animeDto.Russian}
	if animeDto.AnimeURL != nil {
		shikiURL := rah.Configuration.ShikimoriURL + *animeDto.AnimeURL
		animeRo.URL = &shikiURL
	}
	if animeDto.PosterURL != nil {
		posterURL := rah.Configuration.ShikimoriURL + *animeDto.PosterURL
		animeRo.PosterURL = &posterURL
	}
	if err := ReturnResponseAsJSON(w, animeRo, 200); err != nil {
		HandleErr(errors.Wrap(err, ""), w, 500, "Error")
	}
}
