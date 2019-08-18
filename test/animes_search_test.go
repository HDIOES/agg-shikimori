package test

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/HDIOES/cpa-backend/integration"
	"github.com/HDIOES/cpa-backend/models"
	"github.com/HDIOES/cpa-backend/rest"
	"github.com/HDIOES/cpa-backend/rest/util"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	_ "github.com/lib/pq"
)

func TestSearchAnimesSuccess(t *testing.T) {
	diContainer.Invoke(func(configuration *util.Configuration, job *integration.ShikimoriJob, newDao *models.NewDAO, animeDao *models.AnimeDAO, genreDao *models.GenreDAO, studioDao *models.StudioDAO, router *mux.Router) {
		if err := clearDb(newDao, animeDao, genreDao, studioDao); err != nil {
			markAsFailAndAbortNow(t, errors.Wrap(err, ""))
		}
		externalGenreID := "1"
		externalStudioID := "1"
		externalAnimeID := "1"

		genreName := "trashcore"
		genreRussianName := "трешкор"
		genreKind := "tv"

		studioName := "trash studio"
		russianStudioName := "треш студия"
		isReal := false
		imageURL := "/url.jpg"

		animeName := "One Punch Man"
		russianAnimeName := "Один Удар Человек"
		animeURL := "/url.jpg"
		animeKind := "tv"
		animeStatus := "ongoing"
		animePostreURL := "/url.jpg"
		var animeEpizodes int64 = 12
		var animeEpizodesAired int64 = 6
		airedOn := time.Date(2009, 11, 17, 20, 20, 20, 0, time.UTC)
		releasedOn := time.Date(2009, 11, 17, 20, 20, 20, 0, time.UTC)

		animeProcessed := false

		//fill database
		genreID, insertGenreErr := insertGenreToDatabase(genreDao, externalGenreID, &genreName, &genreRussianName, &genreKind)
		if insertGenreErr != nil {
			markAsFailAndAbortNow(t, errors.Wrap(insertGenreErr, ""))
		}
		studioID, insertStudioErr := insertStudioToDatabase(studioDao, externalStudioID, &studioName, &russianStudioName, &isReal, &imageURL)
		if insertStudioErr != nil {
			markAsFailAndAbortNow(t, errors.Wrap(insertStudioErr, ""))
		}
		animeID, insertAnimeErr := insertAnimeToDatabase(animeDao, externalAnimeID, &animeName, &russianAnimeName, &animeURL, &animeKind, &animeStatus, &animeEpizodes, &animeEpizodesAired,
			&airedOn,
			&releasedOn, &animePostreURL, &animeProcessed)
		if insertAnimeErr != nil {
			markAsFailAndAbortNow(t, errors.Wrap(insertAnimeErr, ""))
		}
		if linkAnimeAndGenreErr := linkAnimeAndGenre(animeDao, animeID, genreID); linkAnimeAndGenreErr != nil {
			markAsFailAndAbortNow(t, errors.Wrap(linkAnimeAndGenreErr, ""))
		}
		if linkAnimeAndStudioErr := linkAnimeAndStudio(animeDao, animeID, studioID); linkAnimeAndStudioErr != nil {
			markAsFailAndAbortNow(t, errors.Wrap(linkAnimeAndStudioErr, ""))
		}
		//create request
		request, err := http.NewRequest("GET", "/api/animes/search", nil)
		if err != nil {
			markAsFailAndAbortNow(t, errors.Wrap(err, ""))
		}
		recorder := executeRequest(request, router)
		//asserts
		abortIfFail(t, assert.Equal(t, 200, recorder.Code))
		//get actual data
		actualJSONResponseBody := recorder.Body.String()
		//form expected data
		animesRos := make([]rest.AnimeRO, 0, 1)
		animePosterURLRO := configuration.ShikimoriURL + animePostreURL
		animeURLRO := configuration.ShikimoriURL + animeURL
		animeRO := rest.AnimeRO{
			Name:        &animeName,
			RussuanName: &russianAnimeName,
			URL:         &animeURLRO,
			PosterURL:   &animePosterURLRO,
		}
		animesRos = append(animesRos, animeRO)
		expectedJSONResponseBodyBytes, marshalErr := json.Marshal(&animesRos)
		if marshalErr != nil {
			markAsFailAndAbortNow(t, errors.Wrap(marshalErr, ""))
		}
		abortIfFail(t, assert.JSONEq(t, string(expectedJSONResponseBodyBytes), actualJSONResponseBody))
	})
}

func TestSearchAnimes_pagingSuccess(t *testing.T) {

}

func TestSearchAnimes_byStatusSuccess(t *testing.T) {
}

func TestSearchAnimes_byKindSuccess(t *testing.T) {
}

func TestSearchAnimes_byOrderSuccess(t *testing.T) {
}

func TestSearchAnimes_byScoreSuccess(t *testing.T) {
}

func TestSearchAnimes_byGenresIdsSuccess(t *testing.T) {
}

func TestSearchAnimes_byStudioIdsSuccess(t *testing.T) {
}

func TestSearchAnimes_byDurationSuccess(t *testing.T) {
}

func TestSearchAnimes_byRatingSuccess(t *testing.T) {
}

func TestSearchAnimes_byFranchiseSuccess(t *testing.T) {
}

func TestSearchAnimes_byIdsSuccess(t *testing.T) {
}

func TestSearchAnimes_byExludeIdsSuccess(t *testing.T) {
}

func TestSearchAnimes_limitFail(t *testing.T) {
	diContainer.Invoke(func(router *mux.Router) {
		request, err := http.NewRequest("GET", "/api/animes/search?limit=34df4", nil)
		if err != nil {
			markAsFailAndAbortNow(t, errors.Wrap(err, ""))
		}
		recorder := executeRequest(request, router)
		assert.Equal(t, 400, recorder.Code)
	})
}

func TestSearchAnimes_offsetFail(t *testing.T) {
	diContainer.Invoke(func(router *mux.Router) {
		request, err := http.NewRequest("GET", "/api/animes/search?offset=df44", nil)
		if err != nil {
			markAsFailAndAbortNow(t, errors.Wrap(err, ""))
		}
		recorder := executeRequest(request, router)
		assert.Equal(t, 400, recorder.Code)
	})
}

func TestSearchAnimes_scoreFail(t *testing.T) {
	diContainer.Invoke(func(router *mux.Router) {
		request, err := http.NewRequest("GET", "/api/animes/search?score=hnk", nil)
		if err != nil {
			markAsFailAndAbortNow(t, err)
		}
		recorder := executeRequest(request, router)
		assert.Equal(t, 400, recorder.Code)
	})
}

func TestRandom_scoreFail(t *testing.T) {
	diContainer.Invoke(func(router *mux.Router) {
		request, err := http.NewRequest("GET", "/api/animes/random?score=hnk", nil)
		if err != nil {
			markAsFailAndAbortNow(t, err)
		}
		recorder := executeRequest(request, router)
		assert.Equal(t, 400, recorder.Code)
	})
}
