package test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/HDIOES/cpa-backend/models"
	"github.com/pkg/errors"

	"github.com/HDIOES/cpa-backend/integration"
	"github.com/HDIOES/cpa-backend/rest/util"
	_ "github.com/lib/pq"

	gock "gopkg.in/h2non/gock.v1"
)

//TestSimple function
func TestShikimoriJobSuccess(t *testing.T) {
	diContainer.Invoke(func(configuration *util.Configuration, job *integration.ShikimoriJob, newDao *models.NewDAO, animeDao *models.AnimeDAO, genreDao *models.GenreDAO, studioDao *models.StudioDAO) {
		if err := clearDb(newDao, animeDao, genreDao, studioDao); err != nil {
			markAsFailAndAbortNow(t, errors.Wrap(err, ""))
		}
		defer gock.Off()

		//mock animes
		animesData, err := ioutil.ReadFile("mock/shikimori_animes_success.json")
		if err != nil {
			markAsFailAndAbortNow(t, errors.Wrap(err, ""))
		}
		gock.New(configuration.ShikimoriURL).
			Get(configuration.ShikimoriAnimeSearchURL).
			MatchParam("page", "1").
			MatchParam("limit", "50").
			Reply(200).
			JSON(animesData)

		//mock genres
		genresData, err := ioutil.ReadFile("mock/shikimori_genres_success.json")
		if err != nil {
			markAsFailAndAbortNow(t, errors.Wrap(err, ""))
		}
		gock.New(configuration.ShikimoriURL).
			Get(configuration.ShikimoriGenreURL).
			Reply(200).
			JSON(genresData)

		//mock studios
		studiosData, err := ioutil.ReadFile("mock/shikimori_studios_success.json")
		if err != nil {
			markAsFailAndAbortNow(t, errors.Wrap(err, ""))
		}
		gock.New(configuration.ShikimoriURL).
			Get(configuration.ShikimoriStudioURL).
			Reply(200).
			JSON(studiosData)

		//mock one anime
		oneAnimeData, err := ioutil.ReadFile("mock/one_anime_shikimori_success.json")
		if err != nil {
			markAsFailAndAbortNow(t, errors.Wrap(err, ""))
		}
		gock.New(configuration.ShikimoriURL).
			Get(configuration.ShikimoriAnimeSearchURL+"/").
			PathParam("animes", "5114").
			Reply(200).
			JSON(oneAnimeData)

		job.Run()

		//asserts
		anime := integration.Anime{}
		genres := []integration.Genre{}
		studios := []integration.Studio{}
		if unmarshalAnimeErr := json.Unmarshal(oneAnimeData, &anime); unmarshalAnimeErr != nil {
			markAsFailAndAbortNow(t, errors.Wrap(unmarshalAnimeErr, ""))
		}
		if unmarshalGenresErr := json.Unmarshal(genresData, &genres); unmarshalGenresErr != nil {
			markAsFailAndAbortNow(t, errors.Wrap(unmarshalGenresErr, ""))
		}
		if unmarshalStudioErr := json.Unmarshal(studiosData, &studios); unmarshalStudioErr != nil {
			markAsFailAndAbortNow(t, errors.Wrap(unmarshalStudioErr, ""))
		}
		for _, s := range studios {
			studioDto, studioDtoErr := studioDao.FindByExternalID(strconv.FormatInt(*s.ID, 10))
			if studioDtoErr != nil {
				markAsFailAndAbortNow(t, errors.Wrap(studioDtoErr, ""))
			}
			abortIfFail(t, assert.Equal(t, strconv.FormatInt(*s.ID, 10), studioDto.ExternalID))
			abortIfFail(t, EqualStringValues(t, s.FilteredName, studioDto.FilteredStudioName))
			abortIfFail(t, EqualStringValues(t, s.Image, studioDto.ImageURL))
			abortIfFail(t, EqualStringValues(t, s.Name, studioDto.Name))
			abortIfFail(t, EqualBoolValues(t, s.Real, studioDto.IsReal))
		}
		for _, g := range genres {
			genreDto, genreDtoErr := genreDao.FindByExternalID(strconv.FormatInt(*g.ID, 10))
			if genreDtoErr != nil {
				markAsFailAndAbortNow(t, errors.Wrap(genreDtoErr, ""))
			}
			abortIfFail(t, assert.Equal(t, strconv.FormatInt(*g.ID, 10), genreDto.ExternalID))
			abortIfFail(t, EqualStringValues(t, g.Kind, genreDto.Kind))
			abortIfFail(t, EqualStringValues(t, g.Name, genreDto.Name))
			abortIfFail(t, EqualStringValues(t, g.Russian, genreDto.Russian))
		}
		animeDto, animeDtoErr := animeDao.FindByExternalID(strconv.FormatInt(*anime.ID, 10))
		if animeDtoErr != nil {
			markAsFailAndAbortNow(t, errors.Wrap(animeDtoErr, ""))
		}
		abortIfFail(t, assert.Equal(t, strconv.FormatInt(*anime.ID, 10), animeDto.ExternalID))
		abortIfFail(t, EqualStringValues(t, anime.Name, animeDto.Name))
		abortIfFail(t, EqualStringValues(t, anime.Russian, animeDto.Russian))
		abortIfFail(t, EqualStringValues(t, anime.URL, animeDto.AnimeURL))
		abortIfFail(t, EqualStringValues(t, anime.Kind, animeDto.Kind))
		abortIfFail(t, EqualStringValues(t, anime.Status, animeDto.Status))
		abortIfFail(t, EqualInt64Values(t, anime.Episodes, animeDto.Epizodes))
		abortIfFail(t, EqualInt64Values(t, anime.EpisodesAired, animeDto.EpizodesAired))
		//ReleasedOn
		//AiredOn
		abortIfFail(t, EqualStringValues(t, anime.Image.Original, animeDto.PosterURL))
		abortIfFail(t, assert.NotNil(t, animeDto.Score))
		dbScore := fmt.Sprintf("%.2f", *animeDto.Score)
		abortIfFail(t, EqualStringValues(t, anime.Score, &dbScore))
		abortIfFail(t, assert.NotNil(t, animeDto.Duration))
		dbDuration := int64(*animeDto.Duration)
		abortIfFail(t, EqualInt64Values(t, anime.Duration, &dbDuration))
		abortIfFail(t, EqualStringValues(t, anime.Rating, animeDto.Rating))
		abortIfFail(t, EqualStringValues(t, anime.Franchise, animeDto.Franchise))
		abortIfFail(t, assert.NotNil(t, animeDto.Processed))
		processed := true
		abortIfFail(t, EqualBoolValues(t, &processed, animeDto.Processed))
	})
}

func TestShikimoriJob_shikimoriGenresError(t *testing.T) {
	diContainer.Invoke(func(configuration *util.Configuration, job *integration.ShikimoriJob, newDao *models.NewDAO, animeDao *models.AnimeDAO, genreDao *models.GenreDAO, studioDao *models.StudioDAO) {
		if err := clearDb(newDao, animeDao, genreDao, studioDao); err != nil {
			markAsFailAndAbortNow(t, errors.Wrap(err, ""))
		}
		defer gock.Off()
		//mark shikimori genres methods as failed
		gock.New(configuration.ShikimoriURL).
			Get(configuration.ShikimoriGenreURL).
			Reply(500)

		//mock shikimori studios
		studiosData, err := ioutil.ReadFile("mock/shikimori_studios_success.json")
		if err != nil {
			markAsFailAndAbortNow(t, errors.Wrap(err, ""))
		}
		gock.New(configuration.ShikimoriURL).
			Get(configuration.ShikimoriStudioURL).
			Reply(200).
			JSON(studiosData)
		job.Run()
		//asserts
		studios := []integration.Studio{}
		if unmarshalStudioErr := json.Unmarshal(studiosData, &studios); unmarshalStudioErr != nil {
			markAsFailAndAbortNow(t, errors.Wrap(unmarshalStudioErr, ""))
		}
		for _, s := range studios {
			studioDto, studioDtoErr := studioDao.FindByExternalID(strconv.FormatInt(*s.ID, 10))
			if studioDtoErr != nil {
				markAsFailAndAbortNow(t, errors.Wrap(studioDtoErr, ""))
			}
			abortIfFail(t, assert.Equal(t, strconv.FormatInt(*s.ID, 10), studioDto.ExternalID))
			abortIfFail(t, EqualStringValues(t, s.FilteredName, studioDto.FilteredStudioName))
			abortIfFail(t, EqualStringValues(t, s.Image, studioDto.ImageURL))
			abortIfFail(t, EqualStringValues(t, s.Name, studioDto.Name))
			abortIfFail(t, EqualBoolValues(t, s.Real, studioDto.IsReal))
		}

		//genre and anime tables should be empty in database
		genreDtos, genreDtoErr := genreDao.FindByFilter(models.GenreQueryBuilder{})
		if genreDtoErr != nil {
			markAsFailAndAbortNow(t, errors.Wrap(genreDtoErr, ""))
		}
		abortIfFail(t, assert.Equal(t, 0, len(genreDtos)))
		animeDtos, animeDtosErr := animeDao.FindByFilter(models.AnimeQueryBuilder{})
		if animeDtosErr != nil {
			markAsFailAndAbortNow(t, errors.Wrap(animeDtosErr, ""))
		}
		abortIfFail(t, assert.Equal(t, 0, len(animeDtos)))
	})
}

func TestShikimoriJob_shikimoriStudiosError(t *testing.T) {
	diContainer.Invoke(func(configuration *util.Configuration, job *integration.ShikimoriJob, newDao *models.NewDAO, animeDao *models.AnimeDAO, genreDao *models.GenreDAO, studioDao *models.StudioDAO) {
		if err := clearDb(newDao, animeDao, genreDao, studioDao); err != nil {
			markAsFailAndAbortNow(t, errors.Wrap(err, ""))
		}
		defer gock.Off()
		//mock studios as failed
		gock.New(configuration.ShikimoriURL).
			Get(configuration.ShikimoriStudioURL).
			Reply(500)

		job.Run()
		//asserts
		//genre, studios and anime tables should be empty in database
		studioDtos, studioDtosErr := studioDao.FindByFilter(models.StudioQueryBuilder{})
		if studioDtosErr != nil {
			markAsFailAndAbortNow(t, errors.Wrap(studioDtosErr, ""))
		}
		abortIfFail(t, assert.Equal(t, 0, len(studioDtos)))
		genreDtos, genreDtoErr := genreDao.FindByFilter(models.GenreQueryBuilder{})
		if genreDtoErr != nil {
			markAsFailAndAbortNow(t, errors.Wrap(genreDtoErr, ""))
		}
		abortIfFail(t, assert.Equal(t, 0, len(genreDtos)))
		animeDtos, animeDtosErr := animeDao.FindByFilter(models.AnimeQueryBuilder{})
		if animeDtosErr != nil {
			markAsFailAndAbortNow(t, errors.Wrap(animeDtosErr, ""))
		}
		abortIfFail(t, assert.Equal(t, 0, len(animeDtos)))
	})
}

func TestShikimoriJob_shikimoriAnimesError(t *testing.T) {
	diContainer.Invoke(func(configuration *util.Configuration, job *integration.ShikimoriJob, newDao *models.NewDAO, animeDao *models.AnimeDAO, genreDao *models.GenreDAO, studioDao *models.StudioDAO) {
		if err := clearDb(newDao, animeDao, genreDao, studioDao); err != nil {
			markAsFailAndAbortNow(t, errors.Wrap(err, ""))
		}
		defer gock.Off()
		//mock genres
		genresData, err := ioutil.ReadFile("mock/shikimori_genres_success.json")
		if err != nil {
			markAsFailAndAbortNow(t, errors.Wrap(err, ""))
		}
		gock.New(configuration.ShikimoriURL).
			Get(configuration.ShikimoriGenreURL).
			Reply(200).
			JSON(genresData)

		//mock studios
		studiosData, err := ioutil.ReadFile("mock/shikimori_studios_success.json")
		if err != nil {
			markAsFailAndAbortNow(t, errors.Wrap(err, ""))
		}
		gock.New(configuration.ShikimoriURL).
			Get(configuration.ShikimoriStudioURL).
			Reply(200).
			JSON(studiosData)
		job.Run()

		//asserts
		genres := []integration.Genre{}
		studios := []integration.Studio{}
		if unmarshalGenresErr := json.Unmarshal(genresData, &genres); unmarshalGenresErr != nil {
			markAsFailAndAbortNow(t, errors.Wrap(unmarshalGenresErr, ""))
		}
		if unmarshalStudioErr := json.Unmarshal(studiosData, &studios); unmarshalStudioErr != nil {
			markAsFailAndAbortNow(t, errors.Wrap(unmarshalStudioErr, ""))
		}
		for _, s := range studios {
			studioDto, studioDtoErr := studioDao.FindByExternalID(strconv.FormatInt(*s.ID, 10))
			if studioDtoErr != nil {
				markAsFailAndAbortNow(t, errors.Wrap(studioDtoErr, ""))
			}
			abortIfFail(t, assert.Equal(t, strconv.FormatInt(*s.ID, 10), studioDto.ExternalID))
			abortIfFail(t, EqualStringValues(t, s.FilteredName, studioDto.FilteredStudioName))
			abortIfFail(t, EqualStringValues(t, s.Image, studioDto.ImageURL))
			abortIfFail(t, EqualStringValues(t, s.Name, studioDto.Name))
			abortIfFail(t, EqualBoolValues(t, s.Real, studioDto.IsReal))
		}
		for _, g := range genres {
			genreDto, genreDtoErr := genreDao.FindByExternalID(strconv.FormatInt(*g.ID, 10))
			if genreDtoErr != nil {
				markAsFailAndAbortNow(t, errors.Wrap(genreDtoErr, ""))
			}
			abortIfFail(t, assert.Equal(t, strconv.FormatInt(*g.ID, 10), genreDto.ExternalID))
			abortIfFail(t, EqualStringValues(t, g.Kind, genreDto.Kind))
			abortIfFail(t, EqualStringValues(t, g.Name, genreDto.Name))
			abortIfFail(t, EqualStringValues(t, g.Russian, genreDto.Russian))
		}
		animeDtos, animeDtosErr := animeDao.FindByFilter(models.AnimeQueryBuilder{})
		if animeDtosErr != nil {
			markAsFailAndAbortNow(t, errors.Wrap(animeDtosErr, ""))
		}
		abortIfFail(t, assert.Equal(t, 0, len(animeDtos)))
	})
}

func TestShikimoriJob_shikimoriOneAnimeError(t *testing.T) {
	diContainer.Invoke(func(configuration *util.Configuration, job *integration.ShikimoriJob, newDao *models.NewDAO, animeDao *models.AnimeDAO, genreDao *models.GenreDAO, studioDao *models.StudioDAO) {
		if err := clearDb(newDao, animeDao, genreDao, studioDao); err != nil {
			markAsFailAndAbortNow(t, errors.Wrap(err, ""))
		}
		defer gock.Off()

		//mock animes
		animesData, err := ioutil.ReadFile("mock/shikimori_animes_success.json")
		if err != nil {
			markAsFailAndAbortNow(t, errors.Wrap(err, ""))
		}
		gock.New(configuration.ShikimoriURL).
			Get(configuration.ShikimoriAnimeSearchURL).
			MatchParam("page", "1").
			MatchParam("limit", "50").
			Reply(200).
			JSON(animesData)

		//mock genres
		genresData, err := ioutil.ReadFile("mock/shikimori_genres_success.json")
		if err != nil {
			markAsFailAndAbortNow(t, errors.Wrap(err, ""))
		}
		gock.New(configuration.ShikimoriURL).
			Get(configuration.ShikimoriGenreURL).
			Reply(200).
			JSON(genresData)

		//mock studios
		studiosData, err := ioutil.ReadFile("mock/shikimori_studios_success.json")
		if err != nil {
			markAsFailAndAbortNow(t, errors.Wrap(err, ""))
		}
		gock.New(configuration.ShikimoriURL).
			Get(configuration.ShikimoriStudioURL).
			Reply(200).
			JSON(studiosData)

		//mock one anime as failed
		gock.New(configuration.ShikimoriURL).
			Get(configuration.ShikimoriAnimeSearchURL + "/").
			Reply(500)
		job.Run()

		//asserts
		animes := []integration.Anime{}
		genres := []integration.Genre{}
		studios := []integration.Studio{}
		if unmarshalAnimesErr := json.Unmarshal(animesData, &animes); unmarshalAnimesErr != nil {
			markAsFailAndAbortNow(t, errors.Wrap(unmarshalAnimesErr, ""))
		}
		if unmarshalGenresErr := json.Unmarshal(genresData, &genres); unmarshalGenresErr != nil {
			markAsFailAndAbortNow(t, errors.Wrap(unmarshalGenresErr, ""))
		}
		if unmarshalStudioErr := json.Unmarshal(studiosData, &studios); unmarshalStudioErr != nil {
			markAsFailAndAbortNow(t, errors.Wrap(unmarshalStudioErr, ""))
		}
		for _, s := range studios {
			studioDto, studioDtoErr := studioDao.FindByExternalID(strconv.FormatInt(*s.ID, 10))
			if studioDtoErr != nil {
				markAsFailAndAbortNow(t, errors.Wrap(studioDtoErr, ""))
			}
			abortIfFail(t, assert.Equal(t, strconv.FormatInt(*s.ID, 10), studioDto.ExternalID))
			abortIfFail(t, EqualStringValues(t, s.FilteredName, studioDto.FilteredStudioName))
			abortIfFail(t, EqualStringValues(t, s.Image, studioDto.ImageURL))
			abortIfFail(t, EqualStringValues(t, s.Name, studioDto.Name))
			abortIfFail(t, EqualBoolValues(t, s.Real, studioDto.IsReal))
		}
		for _, g := range genres {
			genreDto, genreDtoErr := genreDao.FindByExternalID(strconv.FormatInt(*g.ID, 10))
			if genreDtoErr != nil {
				markAsFailAndAbortNow(t, errors.Wrap(genreDtoErr, ""))
			}
			abortIfFail(t, assert.Equal(t, strconv.FormatInt(*g.ID, 10), genreDto.ExternalID))
			abortIfFail(t, EqualStringValues(t, g.Kind, genreDto.Kind))
			abortIfFail(t, EqualStringValues(t, g.Name, genreDto.Name))
			abortIfFail(t, EqualStringValues(t, g.Russian, genreDto.Russian))
		}
		for _, a := range animes {
			animeDto, animeDtoErr := animeDao.FindByExternalID(strconv.FormatInt(*a.ID, 10))
			if animeDtoErr != nil {
				markAsFailAndAbortNow(t, errors.Wrap(animeDtoErr, ""))
			}
			abortIfFail(t, assert.Equal(t, strconv.FormatInt(*a.ID, 10), animeDto.ExternalID))
			abortIfFail(t, EqualStringValues(t, a.Name, animeDto.Name))
			abortIfFail(t, EqualStringValues(t, a.Russian, animeDto.Russian))
			abortIfFail(t, EqualStringValues(t, a.URL, animeDto.AnimeURL))
			abortIfFail(t, EqualStringValues(t, a.Kind, animeDto.Kind))
			abortIfFail(t, EqualStringValues(t, a.Status, animeDto.Status))
			abortIfFail(t, EqualInt64Values(t, a.Episodes, animeDto.Epizodes))
			abortIfFail(t, EqualInt64Values(t, a.EpisodesAired, animeDto.EpizodesAired))
			//ReleasedOn
			//AiredOn
			abortIfFail(t, EqualStringValues(t, a.Image.Original, animeDto.PosterURL))
			abortIfFail(t, assert.Nil(t, animeDto.Score))
			abortIfFail(t, assert.Nil(t, animeDto.Duration))
			abortIfFail(t, assert.Nil(t, animeDto.Rating))
			abortIfFail(t, assert.Nil(t, animeDto.Franchise))
			abortIfFail(t, assert.NotNil(t, animeDto.Processed))
			processed := false
			abortIfFail(t, EqualBoolValues(t, &processed, animeDto.Processed))
		}
	})
}

func TestShikimoriJob_updateSameExistingInDatabaseGenres(t *testing.T) {
	diContainer.Invoke(func(configuration *util.Configuration, job *integration.ShikimoriJob, newDao *models.NewDAO, animeDao *models.AnimeDAO, genreDao *models.GenreDAO, studioDao *models.StudioDAO) {
		if err := clearDb(newDao, animeDao, genreDao, studioDao); err != nil {
			markAsFailAndAbortNow(t, errors.Wrap(err, ""))
		}
		defer gock.Off()
		//insert predefined test data to database
		genreDto := models.GenreDTO{}
		genreDto.ExternalID = "89"
		genre89name := "genre89"
		genreDto.Name = &genre89name
		russianGenre89 := "russianGenre89"
		genreDto.Russian = &russianGenre89
		genreKind := "genreKind"
		genreDto.Kind = &genreKind
		_, genreErr := genreDao.Create(genreDto)
		if genreErr != nil {
			markAsFailAndAbortNow(t, errors.Wrap(genreErr, ""))
		}

		//mock empty animes
		gock.New(configuration.ShikimoriURL).
			Get(configuration.ShikimoriAnimeSearchURL).
			MatchParam("page", "1").
			MatchParam("limit", "50").
			Reply(200).
			JSON("[]")

		//mock genres
		genresData, err := ioutil.ReadFile("mock/one_genre89_shikimori_success.json")
		if err != nil {
			markAsFailAndAbortNow(t, errors.Wrap(err, ""))
		}
		gock.New(configuration.ShikimoriURL).
			Get(configuration.ShikimoriGenreURL).
			Reply(200).
			JSON(genresData)

		//mock studios
		gock.New(configuration.ShikimoriURL).
			Get(configuration.ShikimoriStudioURL).
			Reply(200).
			JSON("[]")

		job.Run()

		genres := []integration.Genre{}
		if unmarshalGenresErr := json.Unmarshal(genresData, &genres); unmarshalGenresErr != nil {
			markAsFailAndAbortNow(t, errors.Wrap(unmarshalGenresErr, ""))
		}
		genreDtos, genreDtosErr := genreDao.FindByFilter(models.GenreQueryBuilder{})
		if genreDtosErr != nil {
			markAsFailAndAbortNow(t, errors.Wrap(genreDtosErr, ""))
		}
		abortIfFail(t, assert.Equal(t, 1, len(genreDtos)))

		for _, g := range genres {
			actualGenreDto, genreDtoErr := genreDao.FindByExternalID(strconv.FormatInt(*g.ID, 10))
			if genreDtoErr != nil {
				markAsFailAndAbortNow(t, errors.Wrap(genreDtoErr, ""))
			}
			abortIfFail(t, assert.Equal(t, strconv.FormatInt(*g.ID, 10), actualGenreDto.ExternalID))
			abortIfFail(t, EqualStringValues(t, g.Kind, actualGenreDto.Kind))
			abortIfFail(t, EqualStringValues(t, g.Name, actualGenreDto.Name))
			abortIfFail(t, EqualStringValues(t, g.Russian, actualGenreDto.Russian))
		}
	})
}

func TestShikimoriJob_updateSameExistingInDatabaseStudios(t *testing.T) {
	diContainer.Invoke(func(configuration *util.Configuration, job *integration.ShikimoriJob, newDao *models.NewDAO, animeDao *models.AnimeDAO, genreDao *models.GenreDAO, studioDao *models.StudioDAO) {
		if err := clearDb(newDao, animeDao, genreDao, studioDao); err != nil {
			markAsFailAndAbortNow(t, errors.Wrap(err, ""))
		}
		defer gock.Off()
		//insert predefined test data to database
		studioDto := models.StudioDTO{}
		studioDto.ExternalID = "1"

		studio1name := "studio1"
		studioDto.Name = &studio1name

		filteredStudioName1 := "filteredStudioName"
		studioDto.FilteredStudioName = &filteredStudioName1

		studioImage := "/image.jpg"
		studioDto.ImageURL = &studioImage

		isReal := false
		studioDto.IsReal = &isReal

		_, studioErr := studioDao.Create(studioDto)
		if studioErr != nil {
			markAsFailAndAbortNow(t, errors.Wrap(studioErr, ""))
		}

		//mock empty animes
		gock.New(configuration.ShikimoriURL).
			Get(configuration.ShikimoriAnimeSearchURL).
			MatchParam("page", "1").
			MatchParam("limit", "50").
			Reply(200).
			JSON("[]")

		//mock empty genres
		gock.New(configuration.ShikimoriURL).
			Get(configuration.ShikimoriGenreURL).
			Reply(200).
			JSON("[]")

		//mock one studio
		studiosData, err := ioutil.ReadFile("mock/one_studio1_shikimori_success.json")
		if err != nil {
			markAsFailAndAbortNow(t, errors.Wrap(err, ""))
		}
		gock.New(configuration.ShikimoriURL).
			Get(configuration.ShikimoriStudioURL).
			Reply(200).
			JSON(studiosData)

		job.Run()

		studios := []integration.Studio{}
		if unmarshalStudiosErr := json.Unmarshal(studiosData, &studios); unmarshalStudiosErr != nil {
			markAsFailAndAbortNow(t, errors.Wrap(unmarshalStudiosErr, ""))
		}
		studioDtos, studioDtosErr := studioDao.FindByFilter(models.StudioQueryBuilder{})
		if studioDtosErr != nil {
			markAsFailAndAbortNow(t, errors.Wrap(studioDtosErr, ""))
		}
		abortIfFail(t, assert.Equal(t, 1, len(studioDtos)))

		for _, s := range studios {
			actualStudioDto, studioDtoErr := studioDao.FindByExternalID(strconv.FormatInt(*s.ID, 10))
			if studioDtoErr != nil {
				markAsFailAndAbortNow(t, errors.Wrap(studioDtoErr, ""))
			}
			abortIfFail(t, assert.Equal(t, strconv.FormatInt(*s.ID, 10), actualStudioDto.ExternalID))
			abortIfFail(t, EqualStringValues(t, s.FilteredName, actualStudioDto.FilteredStudioName))
			abortIfFail(t, EqualStringValues(t, s.Image, actualStudioDto.ImageURL))
			abortIfFail(t, EqualStringValues(t, s.Name, actualStudioDto.Name))
			abortIfFail(t, EqualBoolValues(t, s.Real, actualStudioDto.IsReal))
		}

	})
}

func TestShikimoriJob_updateSameExistingInDatabaseAnimes(t *testing.T) {
	diContainer.Invoke(func(configuration *util.Configuration, job *integration.ShikimoriJob, newDao *models.NewDAO, animeDao *models.AnimeDAO, genreDao *models.GenreDAO, studioDao *models.StudioDAO) {
		if err := clearDb(newDao, animeDao, genreDao, studioDao); err != nil {
			markAsFailAndAbortNow(t, errors.Wrap(err, ""))
		}
		defer gock.Off()
		//insert genre 1 to database
		genreExternalID1 := "234"
		genreName1 := "trashcore"
		genreRussianName1 := "трешкор"
		genreKind1 := "tv"
		genreDTO1 := models.GenreDTO{
			ExternalID: genreExternalID1,
			Name:       &genreName1,
			Russian:    &genreRussianName1,
			Kind:       &genreKind1,
		}
		genreID1, insertGenreErr := insertGenreToDatabase(genreDao, genreDTO1)
		if insertGenreErr != nil {
			markAsFailAndAbortNow(t, errors.Wrap(insertGenreErr, ""))
		}
		//insert studio 1 to database
		studioExternalID1 := "234"
		studioName1 := "trash studio"
		studioFilteredName1 := "треш студия"
		studioIsReal1 := false
		studioImageURL1 := "/url.jpg"
		studioDTO1 := models.StudioDTO{
			ExternalID:         studioExternalID1,
			Name:               &studioName1,
			FilteredStudioName: &studioFilteredName1,
			IsReal:             &studioIsReal1,
			ImageURL:           &studioImageURL1,
		}
		studioID1, insertStudioErr := insertStudioToDatabase(studioDao, studioDTO1)
		if insertStudioErr != nil {
			markAsFailAndAbortNow(t, errors.Wrap(insertStudioErr, ""))
		}
		buildAnime5114(t, animeDao, studioID1, genreID1)

		//mock animes
		animesData, err := ioutil.ReadFile("mock/shikimori_animes_success.json")
		if err != nil {
			markAsFailAndAbortNow(t, errors.Wrap(err, ""))
		}
		gock.New(configuration.ShikimoriURL).
			Get(configuration.ShikimoriAnimeSearchURL).
			MatchParam("page", "1").
			MatchParam("limit", "50").
			Reply(200).
			JSON(animesData)

		//mock genres
		genresData, err := ioutil.ReadFile("mock/shikimori_genres_success.json")
		if err != nil {
			markAsFailAndAbortNow(t, errors.Wrap(err, ""))
		}
		gock.New(configuration.ShikimoriURL).
			Get(configuration.ShikimoriGenreURL).
			Reply(200).
			JSON(genresData)

		//mock studios
		studiosData, err := ioutil.ReadFile("mock/shikimori_studios_success.json")
		if err != nil {
			markAsFailAndAbortNow(t, errors.Wrap(err, ""))
		}
		gock.New(configuration.ShikimoriURL).
			Get(configuration.ShikimoriStudioURL).
			Reply(200).
			JSON(studiosData)

		//mock one anime
		oneAnimeData, err := ioutil.ReadFile("mock/one_anime_shikimori_success.json")
		if err != nil {
			markAsFailAndAbortNow(t, errors.Wrap(err, ""))
		}
		gock.New(configuration.ShikimoriURL).
			Get(configuration.ShikimoriAnimeSearchURL+"/").
			PathParam("animes", "5114").
			Reply(200).
			JSON(oneAnimeData)

		job.Run()

		//asserts
		anime := integration.Anime{}
		genres := []integration.Genre{}
		studios := []integration.Studio{}
		if unmarshalAnimeErr := json.Unmarshal(oneAnimeData, &anime); unmarshalAnimeErr != nil {
			markAsFailAndAbortNow(t, errors.Wrap(unmarshalAnimeErr, ""))
		}
		if unmarshalGenresErr := json.Unmarshal(genresData, &genres); unmarshalGenresErr != nil {
			markAsFailAndAbortNow(t, errors.Wrap(unmarshalGenresErr, ""))
		}
		if unmarshalStudioErr := json.Unmarshal(studiosData, &studios); unmarshalStudioErr != nil {
			markAsFailAndAbortNow(t, errors.Wrap(unmarshalStudioErr, ""))
		}
		for _, s := range studios {
			studioDto, studioDtoErr := studioDao.FindByExternalID(strconv.FormatInt(*s.ID, 10))
			if studioDtoErr != nil {
				markAsFailAndAbortNow(t, errors.Wrap(studioDtoErr, ""))
			}
			abortIfFail(t, assert.Equal(t, strconv.FormatInt(*s.ID, 10), studioDto.ExternalID))
			abortIfFail(t, EqualStringValues(t, s.FilteredName, studioDto.FilteredStudioName))
			abortIfFail(t, EqualStringValues(t, s.Image, studioDto.ImageURL))
			abortIfFail(t, EqualStringValues(t, s.Name, studioDto.Name))
			abortIfFail(t, EqualBoolValues(t, s.Real, studioDto.IsReal))
		}
		for _, g := range genres {
			genreDto, genreDtoErr := genreDao.FindByExternalID(strconv.FormatInt(*g.ID, 10))
			if genreDtoErr != nil {
				markAsFailAndAbortNow(t, errors.Wrap(genreDtoErr, ""))
			}
			abortIfFail(t, assert.Equal(t, strconv.FormatInt(*g.ID, 10), genreDto.ExternalID))
			abortIfFail(t, EqualStringValues(t, g.Kind, genreDto.Kind))
			abortIfFail(t, EqualStringValues(t, g.Name, genreDto.Name))
			abortIfFail(t, EqualStringValues(t, g.Russian, genreDto.Russian))
		}
		animeDto, animeDtoErr := animeDao.FindByExternalID(strconv.FormatInt(*anime.ID, 10))
		if animeDtoErr != nil {
			markAsFailAndAbortNow(t, errors.Wrap(animeDtoErr, ""))
		}
		abortIfFail(t, assert.Equal(t, strconv.FormatInt(*anime.ID, 10), animeDto.ExternalID))
		abortIfFail(t, EqualStringValues(t, anime.Name, animeDto.Name))
		abortIfFail(t, EqualStringValues(t, anime.Russian, animeDto.Russian))
		abortIfFail(t, EqualStringValues(t, anime.URL, animeDto.AnimeURL))
		abortIfFail(t, EqualStringValues(t, anime.Kind, animeDto.Kind))
		abortIfFail(t, EqualStringValues(t, anime.Status, animeDto.Status))
		abortIfFail(t, EqualInt64Values(t, anime.Episodes, animeDto.Epizodes))
		abortIfFail(t, EqualInt64Values(t, anime.EpisodesAired, animeDto.EpizodesAired))
		//ReleasedOn
		//AiredOn
		abortIfFail(t, EqualStringValues(t, anime.Image.Original, animeDto.PosterURL))
		abortIfFail(t, assert.NotNil(t, animeDto.Score))
		dbScore := fmt.Sprintf("%.2f", *animeDto.Score)
		abortIfFail(t, EqualStringValues(t, anime.Score, &dbScore))
		abortIfFail(t, assert.NotNil(t, animeDto.Duration))
		dbDuration := int64(*animeDto.Duration)
		abortIfFail(t, EqualInt64Values(t, anime.Duration, &dbDuration))
		abortIfFail(t, EqualStringValues(t, anime.Rating, animeDto.Rating))
		abortIfFail(t, EqualStringValues(t, anime.Franchise, animeDto.Franchise))
		abortIfFail(t, assert.NotNil(t, animeDto.Processed))
		processed := true
		abortIfFail(t, EqualBoolValues(t, &processed, animeDto.Processed))

	})
}
