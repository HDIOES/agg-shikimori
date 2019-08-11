package test

import (
	"io/ioutil"
	"testing"

	"github.com/HDIOES/cpa-backend/models"

	"github.com/HDIOES/cpa-backend/integration"
	"github.com/HDIOES/cpa-backend/rest/util"
	_ "github.com/lib/pq"

	gock "gopkg.in/h2non/gock.v1"
)

//TestSimple function
func TestShikimoriJobSuccess(t *testing.T) {
	diContainer.Invoke(func(configuration *util.Configuration, job *integration.ShikimoriJob, newDao *models.NewDAO, animeDao *models.AnimeDAO, genreDao *models.GenreDAO, studioDao *models.StudioDAO) {
		clearDb(newDao, animeDao, genreDao, studioDao)
		defer gock.Off()
		animesData, err := ioutil.ReadFile("mock/shikimori_animes_success.json")
		if err != nil {
			t.Fatal(err)
		}
		gock.New(configuration.ShikimoriURL).
			Get(configuration.ShikimoriAnimeSearchURL).
			MatchParam("page", "1").
			MatchParam("limit", "50").
			Reply(200).
			JSON(animesData)

		genresData, err := ioutil.ReadFile("mock/shikimori_genres_success.json")
		if err != nil {
			t.Fatal(err)
		}
		gock.New(configuration.ShikimoriURL).
			Get(configuration.ShikimoriGenreURL).
			Reply(200).
			JSON(genresData)

		studiosData, err := ioutil.ReadFile("mock/shikimori_studios_success.json")
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

		job.Run()
	})
}
