package test

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	_ "github.com/lib/pq"
)

func TestSearchAnimesSuccess(t *testing.T) {
	clearDb(db)
	const externalGenreID string = "1"
	const externalStudioID string = "1"
	const externalAnimeID string = "1"
	//fill database
	insertGenreErr := insertGenreToDatabase(t, externalGenreID, "trashcore", "трешкор", "anime")
	if insertGenreErr != nil {
		t.Fatal(insertGenreErr)
	}
	insertStudioErr := insertStudioToDatabase(t, externalStudioID, "trash studio", "треш студия", true, "/url.jpg")
	if insertStudioErr != nil {
		t.Fatal(insertStudioErr)
	}
	insertAnimeErr := insertAnimeToDatabase(t, externalAnimeID, "One Punch Man", "Один Удар Человек", "/url.jpg", "tv", "ongoing", 10, 5,
		time.Date(2009, 11, 17, 20, 20, 20, 0, time.UTC),
		time.Date(2009, 11, 17, 20, 20, 20, 0, time.UTC), "/url.jpg", false, externalStudioID, externalGenreID)
	if insertAnimeErr != nil {
		t.Fatal(insertAnimeErr)
	}
	//create request
	request, _ := http.NewRequest("GET", "/api/animes/search", nil)
	recorder := executeRequest(request)
	//asserts
	assert.Equal(t, 200, recorder.Code)

}

func TestSearchAnimes_limitFail(t *testing.T) {
	request, _ := http.NewRequest("GET", "/api/animes/search?limit=34df4", nil)
	recorder := executeRequest(request)
	assert.Equal(t, 400, recorder.Code)
}

func TestSearchAnimes_offsetFail(t *testing.T) {
	request, _ := http.NewRequest("GET", "/api/animes/search?offset=df44", nil)
	recorder := executeRequest(request)
	assert.Equal(t, 400, recorder.Code)
}

func TestSearchAnimes_scoreFail(t *testing.T) {
	request, _ := http.NewRequest("GET", "/api/animes/search?score=hnk", nil)
	recorder := executeRequest(request)
	assert.Equal(t, 400, recorder.Code)
}

func TestRandom_scoreFail(t *testing.T) {
	request, _ := http.NewRequest("GET", "/api/animes/random?score=hnk", nil)
	recorder := executeRequest(request)
	assert.Equal(t, 400, recorder.Code)
}
