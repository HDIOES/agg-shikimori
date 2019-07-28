package test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	gock "gopkg.in/h2non/gock.v1"

	"github.com/HDIOES/cpa-backend/integration"
	_ "github.com/lib/pq"
)

func TestSearchAnimesSuccess(t *testing.T) {
	//fill database
	fillDatabaseByAnimes(t,
		"mock/shikimori_animes_success.json",
		"mock/shikimori_genres_success.json",
		"mock/shikimori_studios_success.json")
	//create request
	request, _ := http.NewRequest("GET", "/api/animes/search", nil)
	executeRequest(request)
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	return rr
}

func fillDatabaseByAnimes(t *testing.T, animesFilePath, genresFilePath, studiosFilePath string) {
	clearDb(db, t)
	defer gock.Off()
	shikimoriJob := &integration.ShikimoriJob{Db: db, Config: configuration}

	animesData, err := ioutil.ReadFile(animesFilePath)
	if err != nil {
		t.Fatal(err)
	}
	gock.New(configuration.ShikimoriURL).
		Get(configuration.ShikimoriAnimeSearchURL).
		MatchParam("page", "1").
		MatchParam("limit", "50").
		Reply(200).
		JSON(animesData)

	genresData, err := ioutil.ReadFile(genresFilePath)
	if err != nil {
		t.Fatal(err)
	}
	gock.New(configuration.ShikimoriURL).
		Get(configuration.ShikimoriGenreURL).
		Reply(200).
		JSON(genresData)

	studiosData, err := ioutil.ReadFile(studiosFilePath)
	if err != nil {
		t.Fatal(err)
	}
	gock.New(configuration.ShikimoriURL).
		Get(configuration.ShikimoriStudioURL).
		Reply(200).
		JSON(studiosData)

	oneAnimeData, err := ioutil.ReadFile("mock/one_anime_shikimori_success.json")
	if err != nil {
		t.Fatal(err)
	}
	gock.New(configuration.ShikimoriURL).
		Get(configuration.ShikimoriAnimeSearchURL+"/").
		PathParam("animes", "5114").
		Reply(200).
		JSON(oneAnimeData)

	shikimoriJob.Run()
}
