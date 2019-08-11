package test

import (
	"net/http"
	"testing"
	"time"

	"github.com/HDIOES/cpa-backend/integration"
	"github.com/HDIOES/cpa-backend/models"
	"github.com/HDIOES/cpa-backend/rest/util"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"

	_ "github.com/lib/pq"
)

func TestSearchAnimesSuccess(t *testing.T) {
	diContainer.Invoke(func(configuration *util.Configuration, job *integration.ShikimoriJob, newDao *models.NewDAO, animeDao *models.AnimeDAO, genreDao *models.GenreDAO, studioDao *models.StudioDAO, router *mux.Router) {
		clearDb(newDao, animeDao, genreDao, studioDao)
		const externalGenreID string = "1"
		const externalStudioID string = "1"
		const externalAnimeID string = "1"
		//fill database
		genreID, insertGenreErr := insertGenreToDatabase(genreDao, externalGenreID, "trashcore", "трешкор", "anime")
		if insertGenreErr != nil {
			t.Fatal(insertGenreErr)
		}
		studioID, insertStudioErr := insertStudioToDatabase(studioDao, externalStudioID, "trash studio", "треш студия", true, "/url.jpg")
		if insertStudioErr != nil {
			t.Fatal(insertStudioErr)
		}
		animeID, insertAnimeErr := insertAnimeToDatabase(animeDao, externalAnimeID, "One Punch Man", "Один Удар Человек", "/url.jpg", "tv", "ongoing", 10, 5,
			time.Date(2009, 11, 17, 20, 20, 20, 0, time.UTC),
			time.Date(2009, 11, 17, 20, 20, 20, 0, time.UTC), "/url.jpg", false)
		if insertAnimeErr != nil {
			t.Fatal(insertAnimeErr)
		}
		linkAnimeAndGenre(animeDao, animeID, genreID)
		linkAnimeAndStudio(animeDao, animeID, studioID)
		//create request
		request, _ := http.NewRequest("GET", "/api/animes/search", nil)
		recorder := executeRequest(request, router)
		//asserts
		assert.Equal(t, 200, recorder.Code)
	})

}

func TestSearchAnimes_limitFail(t *testing.T) {
	diContainer.Invoke(func(router *mux.Router) {
		request, _ := http.NewRequest("GET", "/api/animes/search?limit=34df4", nil)
		recorder := executeRequest(request, router)
		assert.Equal(t, 400, recorder.Code)
	})
}

func TestSearchAnimes_offsetFail(t *testing.T) {
	diContainer.Invoke(func(router *mux.Router) {
		request, _ := http.NewRequest("GET", "/api/animes/search?offset=df44", nil)
		recorder := executeRequest(request, router)
		assert.Equal(t, 400, recorder.Code)
	})
}

func TestSearchAnimes_scoreFail(t *testing.T) {
	diContainer.Invoke(func(router *mux.Router) {
		request, _ := http.NewRequest("GET", "/api/animes/search?score=hnk", nil)
		recorder := executeRequest(request, router)
		assert.Equal(t, 400, recorder.Code)
	})
}

func TestRandom_scoreFail(t *testing.T) {
	diContainer.Invoke(func(router *mux.Router) {
		request, _ := http.NewRequest("GET", "/api/animes/random?score=hnk", nil)
		recorder := executeRequest(request, router)
		assert.Equal(t, 400, recorder.Code)
	})
}
