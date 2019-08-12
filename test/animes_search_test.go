package test

import (
	"net/http"
	"testing"
	"time"

	"github.com/HDIOES/cpa-backend/integration"
	"github.com/HDIOES/cpa-backend/models"
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
		animeImageURL := "/url.jpg"
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
		animeID, insertAnimeErr := insertAnimeToDatabase(animeDao, externalAnimeID, &animeName, &russianAnimeName, &animeImageURL, &animeKind, &animeStatus, &animeEpizodes, &animeEpizodesAired,
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
		assert.Equal(t, 200, recorder.Code)
	})
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
