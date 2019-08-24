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

/*func TestSearchAnimesSuccess(t *testing.T) {
	diContainer.Invoke(func(configuration *util.Configuration, job *integration.ShikimoriJob, newDao *models.NewDAO, animeDao *models.AnimeDAO, genreDao *models.GenreDAO, studioDao *models.StudioDAO, router *mux.Router) {
		if err := clearDb(newDao, animeDao, genreDao, studioDao); err != nil {
			markAsFailAndAbortNow(t, errors.Wrap(err, ""))
		}

		prepareTestData(t, animeDao, genreDao, studioDao)
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
}*/

func TestSearchAnimes_pagingSuccess(t *testing.T) {
	diContainer.Invoke(func(configuration *util.Configuration, job *integration.ShikimoriJob, newDao *models.NewDAO, animeDao *models.AnimeDAO, genreDao *models.GenreDAO, studioDao *models.StudioDAO, router *mux.Router) {
		if err := clearDb(newDao, animeDao, genreDao, studioDao); err != nil {
			markAsFailAndAbortNow(t, errors.Wrap(err, ""))
		}
		prepareTestData(t, animeDao, genreDao, studioDao)
		//create request
		request, err := http.NewRequest("GET", "/api/animes/search?limit=2&offset=2", nil)
		if err != nil {
			markAsFailAndAbortNow(t, errors.Wrap(err, ""))
		}
		recorder := executeRequest(request, router)
		//asserts
		abortIfFail(t, assert.Equal(t, 200, recorder.Code))
		expectedAnimesRos := make([]rest.AnimeRO, 0, 2)

		animeExternalID2 := "3"
		animeName2 := "One Punch Man"
		russianAnimeName2 := "Один Удар Человек"
		animeURL2 := "/url.jpg"
		animePostreURL2 := "/url.jpg"
		animePosterURLRO2 := configuration.ShikimoriURL + animePostreURL2
		animeURLRO2 := configuration.ShikimoriURL + animeURL2
		animeRO2 := rest.AnimeRO{
			ShikiID:     animeExternalID2,
			Name:        &animeName2,
			RussuanName: &russianAnimeName2,
			URL:         &animeURLRO2,
			PosterURL:   &animePosterURLRO2,
		}
		expectedAnimesRos = append(expectedAnimesRos, animeRO2)

		animeExternalID3 := "4"
		animeName3 := "One Punch Man"
		russianAnimeName3 := "Один Удар Человек"
		animeURL3 := "/url.jpg"
		animePostreURL3 := "/url.jpg"
		animePosterURLRO3 := configuration.ShikimoriURL + animePostreURL3
		animeURLRO3 := configuration.ShikimoriURL + animeURL3
		animeRO3 := rest.AnimeRO{
			ShikiID:     animeExternalID3,
			Name:        &animeName3,
			RussuanName: &russianAnimeName3,
			URL:         &animeURLRO3,
			PosterURL:   &animePosterURLRO3,
		}
		expectedAnimesRos = append(expectedAnimesRos, animeRO3)

		//get actual data
		actualJSONResponseBody := recorder.Body.String()
		expectedJSONResponseBodyBytes, marshalErr := json.Marshal(&expectedAnimesRos)
		if marshalErr != nil {
			markAsFailAndAbortNow(t, errors.Wrap(marshalErr, ""))
		}
		abortIfFail(t, assert.JSONEq(t, string(expectedJSONResponseBodyBytes), actualJSONResponseBody))
	})
}

func TestSearchAnimes_byStatusSuccess(t *testing.T) {
	diContainer.Invoke(func(configuration *util.Configuration, job *integration.ShikimoriJob, newDao *models.NewDAO, animeDao *models.AnimeDAO, genreDao *models.GenreDAO, studioDao *models.StudioDAO, router *mux.Router) {
		if err := clearDb(newDao, animeDao, genreDao, studioDao); err != nil {
			markAsFailAndAbortNow(t, errors.Wrap(err, ""))
		}
		prepareTestData(t, animeDao, genreDao, studioDao)
		//create request
		request, err := http.NewRequest("GET", "/api/animes/search?status=anons", nil)
		if err != nil {
			markAsFailAndAbortNow(t, errors.Wrap(err, ""))
		}
		recorder := executeRequest(request, router)
		//asserts
		abortIfFail(t, assert.Equal(t, 200, recorder.Code))
		expectedAnimesRos := make([]rest.AnimeRO, 0, 1)

		animeExternalID8 := "8"
		animeName8 := "One Punch Man"
		russianAnimeName8 := "Один Удар Человек"
		animeURL8 := "/url.jpg"
		animePostreURL8 := "/url.jpg"
		animePosterURLRO8 := configuration.ShikimoriURL + animePostreURL8
		animeURLRO8 := configuration.ShikimoriURL + animeURL8
		animeRO8 := rest.AnimeRO{
			ShikiID:     animeExternalID8,
			Name:        &animeName8,
			RussuanName: &russianAnimeName8,
			URL:         &animeURLRO8,
			PosterURL:   &animePosterURLRO8,
		}
		expectedAnimesRos = append(expectedAnimesRos, animeRO8)

		//get actual data
		actualJSONResponseBody := recorder.Body.String()
		expectedJSONResponseBodyBytes, marshalErr := json.Marshal(&expectedAnimesRos)
		if marshalErr != nil {
			markAsFailAndAbortNow(t, errors.Wrap(marshalErr, ""))
		}
		abortIfFail(t, assert.JSONEq(t, string(expectedJSONResponseBodyBytes), actualJSONResponseBody))
	})
}

func TestSearchAnimes_byKindSuccess(t *testing.T) {
	diContainer.Invoke(func(configuration *util.Configuration, job *integration.ShikimoriJob, newDao *models.NewDAO, animeDao *models.AnimeDAO, genreDao *models.GenreDAO, studioDao *models.StudioDAO, router *mux.Router) {
		if err := clearDb(newDao, animeDao, genreDao, studioDao); err != nil {
			markAsFailAndAbortNow(t, errors.Wrap(err, ""))
		}
		prepareTestData(t, animeDao, genreDao, studioDao)
		//create request
		request, err := http.NewRequest("GET", "/api/animes/search?kind=movie", nil)
		if err != nil {
			markAsFailAndAbortNow(t, errors.Wrap(err, ""))
		}
		recorder := executeRequest(request, router)
		//asserts
		abortIfFail(t, assert.Equal(t, 200, recorder.Code))
		expectedAnimesRos := make([]rest.AnimeRO, 0, 1)

		animeExternalID10 := "10"
		animeName10 := "One Punch Man"
		russianAnimeName10 := "Один Удар Человек"
		animeURL10 := "/url.jpg"
		animePostreURL10 := "/url.jpg"
		animePosterURLRO10 := configuration.ShikimoriURL + animePostreURL10
		animeURLRO10 := configuration.ShikimoriURL + animeURL10
		animeRO10 := rest.AnimeRO{
			ShikiID:     animeExternalID10,
			Name:        &animeName10,
			RussuanName: &russianAnimeName10,
			URL:         &animeURLRO10,
			PosterURL:   &animePosterURLRO10,
		}
		expectedAnimesRos = append(expectedAnimesRos, animeRO10)

		//get actual data
		actualJSONResponseBody := recorder.Body.String()
		expectedJSONResponseBodyBytes, marshalErr := json.Marshal(&expectedAnimesRos)
		if marshalErr != nil {
			markAsFailAndAbortNow(t, errors.Wrap(marshalErr, ""))
		}
		abortIfFail(t, assert.JSONEq(t, string(expectedJSONResponseBodyBytes), actualJSONResponseBody))
	})
}

func TestSearchAnimes_byOrderSuccess(t *testing.T) {
	diContainer.Invoke(func(configuration *util.Configuration, job *integration.ShikimoriJob, newDao *models.NewDAO, animeDao *models.AnimeDAO, genreDao *models.GenreDAO, studioDao *models.StudioDAO, router *mux.Router) {
		if err := clearDb(newDao, animeDao, genreDao, studioDao); err != nil {
			markAsFailAndAbortNow(t, errors.Wrap(err, ""))
		}
		prepareTestData(t, animeDao, genreDao, studioDao)
		//create request
		request, err := http.NewRequest("GET", "/api/animes/search?kind=movie&order=aired_on", nil)
		if err != nil {
			markAsFailAndAbortNow(t, errors.Wrap(err, ""))
		}
		recorder := executeRequest(request, router)
		//asserts
		abortIfFail(t, assert.Equal(t, 200, recorder.Code))
		expectedAnimesRos := make([]rest.AnimeRO, 0, 1)

		animeExternalID1 := "1"
		animeName1 := "One Punch Man"
		russianAnimeName1 := "Один Удар Человек"
		animeURL1 := "/url.jpg"
		animePostreURL1 := "/url.jpg"
		animePosterURLRO1 := configuration.ShikimoriURL + animePostreURL1
		animeURLRO1 := configuration.ShikimoriURL + animeURL1
		animeRO1 := rest.AnimeRO{
			ShikiID:     animeExternalID1,
			Name:        &animeName1,
			RussuanName: &russianAnimeName1,
			URL:         &animeURLRO1,
			PosterURL:   &animePosterURLRO1,
		}
		expectedAnimesRos = append(expectedAnimesRos, animeRO1)

		animeExternalID10 := "10"
		animeName10 := "One Punch Man"
		russianAnimeName10 := "Один Удар Человек"
		animeURL10 := "/url.jpg"
		animePostreURL10 := "/url.jpg"
		animePosterURLRO10 := configuration.ShikimoriURL + animePostreURL10
		animeURLRO10 := configuration.ShikimoriURL + animeURL10
		animeRO10 := rest.AnimeRO{
			ShikiID:     animeExternalID10,
			Name:        &animeName10,
			RussuanName: &russianAnimeName10,
			URL:         &animeURLRO10,
			PosterURL:   &animePosterURLRO10,
		}
		expectedAnimesRos = append(expectedAnimesRos, animeRO10)

		animeExternalID5 := "5"
		animeName5 := "One Punch Man"
		russianAnimeName5 := "Один Удар Человек"
		animeURL5 := "/url.jpg"
		animePostreURL5 := "/url.jpg"
		animePosterURLRO5 := configuration.ShikimoriURL + animePostreURL5
		animeURLRO5 := configuration.ShikimoriURL + animeURL5
		animeRO5 := rest.AnimeRO{
			ShikiID:     animeExternalID5,
			Name:        &animeName5,
			RussuanName: &russianAnimeName5,
			URL:         &animeURLRO5,
			PosterURL:   &animePosterURLRO5,
		}
		expectedAnimesRos = append(expectedAnimesRos, animeRO5)

		//get actual data
		actualJSONResponseBody := recorder.Body.String()
		expectedJSONResponseBodyBytes, marshalErr := json.Marshal(&expectedAnimesRos)
		if marshalErr != nil {
			markAsFailAndAbortNow(t, errors.Wrap(marshalErr, ""))
		}
		abortIfFail(t, assert.JSONEq(t, string(expectedJSONResponseBodyBytes), actualJSONResponseBody))
	})
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

//prepareTestData prepares test data in db includes 10 animes with different externalId, 1 genre and 1 studio
func prepareTestData(t *testing.T, animeDao *models.AnimeDAO, genreDao *models.GenreDAO, studioDao *models.StudioDAO) {

	genreExternalID := "234"
	genreName := "trashcore"
	genreRussianName := "трешкор"
	genreKind := "tv"

	studioExternalID := "345"
	studioName := "trash studio"
	studioFilteredName := "треш студия"
	studioIsReal := false
	studioImageURL := "/url.jpg"

	//insert genre to database
	genreDTO := models.GenreDTO{
		ExternalID: genreExternalID,
		Name:       &genreName,
		Russian:    &genreRussianName,
		Kind:       &genreKind,
	}
	genreID, insertGenreErr := insertGenreToDatabase(genreDao, genreDTO)
	if insertGenreErr != nil {
		markAsFailAndAbortNow(t, errors.Wrap(insertGenreErr, ""))
	}
	//insert studio to database
	studioDTO := models.StudioDTO{
		ExternalID:         studioExternalID,
		Name:               &studioName,
		FilteredStudioName: &studioFilteredName,
		IsReal:             &studioIsReal,
		ImageURL:           &studioImageURL,
	}
	studioID, insertStudioErr := insertStudioToDatabase(studioDao, studioDTO)
	if insertStudioErr != nil {
		markAsFailAndAbortNow(t, errors.Wrap(insertStudioErr, ""))
	}

	buildAnime1(t, animeDao, genreID, studioID)
	buildAnime2(t, animeDao, genreID, studioID)
	buildAnime3(t, animeDao, genreID, studioID)
	buildAnime4(t, animeDao, genreID, studioID)
	buildAnime5(t, animeDao, genreID, studioID)
	buildAnime6(t, animeDao, genreID, studioID)
	buildAnime7(t, animeDao, genreID, studioID)
	buildAnime8(t, animeDao, genreID, studioID)
	buildAnime9(t, animeDao, genreID, studioID)
	buildAnime10(t, animeDao, genreID, studioID)

}

func buildAnime1(t *testing.T, animeDao *models.AnimeDAO, genreID int64, studioID int64) {
	//insert anime 1 to database
	animeExternalID1 := "1"
	animeName1 := "One Punch Man"
	animeRussianName1 := "Один Удар Человек"
	animeURL1 := "/url.jpg"
	animeKind1 := "movie"
	animeStatus1 := "ongoing"
	var animeEpizodes1 int64 = 12
	var animeEpizodesAired1 int64 = 6
	animeAiredOn1 := time.Date(2009, 11, 17, 20, 20, 20, 0, time.UTC)
	animeReleasedOn1 := time.Date(2009, 11, 17, 20, 20, 20, 0, time.UTC)
	animePosterURL1 := "/url.jpg"
	animeScore1 := 7.12
	animeDuration1 := 10.0
	animeRating1 := "r"
	animeFranchise1 := "onepunchman"
	animeProcessed1 := false
	testAnimeDto1 := models.AnimeDTO{
		ExternalID:    animeExternalID1,
		Name:          &animeName1,
		Russian:       &animeRussianName1,
		AnimeURL:      &animeURL1,
		Kind:          &animeKind1,
		Status:        &animeStatus1,
		Epizodes:      &animeEpizodes1,
		EpizodesAired: &animeEpizodesAired1,
		AiredOn:       &animeAiredOn1,
		ReleasedOn:    &animeReleasedOn1,
		PosterURL:     &animePosterURL1,
		Score:         &animeScore1,
		Duration:      &animeDuration1,
		Rating:        &animeRating1,
		Franchise:     &animeFranchise1,
		Processed:     &animeProcessed1,
	}
	animeID1, insertAnimeErr := insertAnimeToDatabase(animeDao, testAnimeDto1)
	if insertAnimeErr != nil {
		markAsFailAndAbortNow(t, errors.Wrap(insertAnimeErr, ""))
	}
	if linkAnimeAndGenreErr := linkAnimeAndGenre(animeDao, animeID1, genreID); linkAnimeAndGenreErr != nil {
		markAsFailAndAbortNow(t, errors.Wrap(linkAnimeAndGenreErr, ""))
	}
	if linkAnimeAndStudioErr := linkAnimeAndStudio(animeDao, animeID1, studioID); linkAnimeAndStudioErr != nil {
		markAsFailAndAbortNow(t, errors.Wrap(linkAnimeAndStudioErr, ""))
	}
}

func buildAnime2(t *testing.T, animeDao *models.AnimeDAO, genreID int64, studioID int64) {
	//insert anime 2 to database
	animeExternalID2 := "2"
	animeName2 := "One Punch Man"
	animeRussianName2 := "Один Удар Человек"
	animeURL2 := "/url.jpg"
	animeKind2 := "tv"
	animeStatus2 := "ongoing"
	var animeEpizodes2 int64 = 22
	var animeEpizodesAired2 int64 = 6
	animeAiredOn2 := time.Date(2009, 11, 17, 20, 20, 20, 0, time.UTC)
	animeReleasedOn2 := time.Date(2009, 11, 17, 20, 20, 20, 0, time.UTC)
	animePosterURL2 := "/url.jpg"
	animeScore2 := 7.22
	animeDuration2 := 20.0
	animeRating2 := "r"
	animeFranchise2 := "onepunchman"
	animeProcessed2 := false
	testAnimeDto2 := models.AnimeDTO{
		ExternalID:    animeExternalID2,
		Name:          &animeName2,
		Russian:       &animeRussianName2,
		AnimeURL:      &animeURL2,
		Kind:          &animeKind2,
		Status:        &animeStatus2,
		Epizodes:      &animeEpizodes2,
		EpizodesAired: &animeEpizodesAired2,
		AiredOn:       &animeAiredOn2,
		ReleasedOn:    &animeReleasedOn2,
		PosterURL:     &animePosterURL2,
		Score:         &animeScore2,
		Duration:      &animeDuration2,
		Rating:        &animeRating2,
		Franchise:     &animeFranchise2,
		Processed:     &animeProcessed2,
	}
	animeID2, insertAnimeErr := insertAnimeToDatabase(animeDao, testAnimeDto2)
	if insertAnimeErr != nil {
		markAsFailAndAbortNow(t, errors.Wrap(insertAnimeErr, ""))
	}
	if linkAnimeAndGenreErr := linkAnimeAndGenre(animeDao, animeID2, genreID); linkAnimeAndGenreErr != nil {
		markAsFailAndAbortNow(t, errors.Wrap(linkAnimeAndGenreErr, ""))
	}
	if linkAnimeAndStudioErr := linkAnimeAndStudio(animeDao, animeID2, studioID); linkAnimeAndStudioErr != nil {
		markAsFailAndAbortNow(t, errors.Wrap(linkAnimeAndStudioErr, ""))
	}
}

func buildAnime3(t *testing.T, animeDao *models.AnimeDAO, genreID int64, studioID int64) {
	//insert anime 3 to database
	animeExternalID3 := "3"
	animeName3 := "One Punch Man"
	animeRussianName3 := "Один Удар Человек"
	animeURL3 := "/url.jpg"
	animeKind3 := "tv"
	animeStatus3 := "ongoing"
	var animeEpizodes3 int64 = 33
	var animeEpizodesAired3 int64 = 6
	animeAiredOn3 := time.Date(2009, 11, 17, 20, 20, 20, 0, time.UTC)
	animeReleasedOn3 := time.Date(2009, 11, 17, 20, 20, 20, 0, time.UTC)
	animePosterURL3 := "/url.jpg"
	animeScore3 := 7.33
	animeDuration3 := 30.0
	animeRating3 := "r"
	animeFranchise3 := "onepunchman"
	animeProcessed3 := false
	testAnimeDto3 := models.AnimeDTO{
		ExternalID:    animeExternalID3,
		Name:          &animeName3,
		Russian:       &animeRussianName3,
		AnimeURL:      &animeURL3,
		Kind:          &animeKind3,
		Status:        &animeStatus3,
		Epizodes:      &animeEpizodes3,
		EpizodesAired: &animeEpizodesAired3,
		AiredOn:       &animeAiredOn3,
		ReleasedOn:    &animeReleasedOn3,
		PosterURL:     &animePosterURL3,
		Score:         &animeScore3,
		Duration:      &animeDuration3,
		Rating:        &animeRating3,
		Franchise:     &animeFranchise3,
		Processed:     &animeProcessed3,
	}
	animeID3, insertAnimeErr := insertAnimeToDatabase(animeDao, testAnimeDto3)
	if insertAnimeErr != nil {
		markAsFailAndAbortNow(t, errors.Wrap(insertAnimeErr, ""))
	}
	if linkAnimeAndGenreErr := linkAnimeAndGenre(animeDao, animeID3, genreID); linkAnimeAndGenreErr != nil {
		markAsFailAndAbortNow(t, errors.Wrap(linkAnimeAndGenreErr, ""))
	}
	if linkAnimeAndStudioErr := linkAnimeAndStudio(animeDao, animeID3, studioID); linkAnimeAndStudioErr != nil {
		markAsFailAndAbortNow(t, errors.Wrap(linkAnimeAndStudioErr, ""))
	}
}

func buildAnime4(t *testing.T, animeDao *models.AnimeDAO, genreID int64, studioID int64) {
	//insert anime 4 to database
	animeExternalID4 := "4"
	animeName4 := "One Punch Man"
	animeRussianName4 := "Один Удар Человек"
	animeURL4 := "/url.jpg"
	animeKind4 := "tv"
	animeStatus4 := "ongoing"
	var animeEpizodes4 int64 = 44
	var animeEpizodesAired4 int64 = 6
	animeAiredOn4 := time.Date(2009, 11, 17, 20, 20, 20, 0, time.UTC)
	animeReleasedOn4 := time.Date(2009, 11, 17, 20, 20, 20, 0, time.UTC)
	animePosterURL4 := "/url.jpg"
	animeScore4 := 7.44
	animeDuration4 := 40.0
	animeRating4 := "r"
	animeFranchise4 := "onepunchman"
	animeProcessed4 := false
	testAnimeDto4 := models.AnimeDTO{
		ExternalID:    animeExternalID4,
		Name:          &animeName4,
		Russian:       &animeRussianName4,
		AnimeURL:      &animeURL4,
		Kind:          &animeKind4,
		Status:        &animeStatus4,
		Epizodes:      &animeEpizodes4,
		EpizodesAired: &animeEpizodesAired4,
		AiredOn:       &animeAiredOn4,
		ReleasedOn:    &animeReleasedOn4,
		PosterURL:     &animePosterURL4,
		Score:         &animeScore4,
		Duration:      &animeDuration4,
		Rating:        &animeRating4,
		Franchise:     &animeFranchise4,
		Processed:     &animeProcessed4,
	}
	animeID4, insertAnimeErr := insertAnimeToDatabase(animeDao, testAnimeDto4)
	if insertAnimeErr != nil {
		markAsFailAndAbortNow(t, errors.Wrap(insertAnimeErr, ""))
	}
	if linkAnimeAndGenreErr := linkAnimeAndGenre(animeDao, animeID4, genreID); linkAnimeAndGenreErr != nil {
		markAsFailAndAbortNow(t, errors.Wrap(linkAnimeAndGenreErr, ""))
	}
	if linkAnimeAndStudioErr := linkAnimeAndStudio(animeDao, animeID4, studioID); linkAnimeAndStudioErr != nil {
		markAsFailAndAbortNow(t, errors.Wrap(linkAnimeAndStudioErr, ""))
	}
}

func buildAnime5(t *testing.T, animeDao *models.AnimeDAO, genreID int64, studioID int64) {
	//insert anime 5 to database
	animeExternalID5 := "5"
	animeName5 := "One Punch Man"
	animeRussianName5 := "Один Удар Человек"
	animeURL5 := "/url.jpg"
	animeKind5 := "movie"
	animeStatus5 := "ongoing"
	var animeEpizodes5 int64 = 55
	var animeEpizodesAired5 int64 = 6
	animeAiredOn5 := time.Date(2011, 11, 17, 20, 20, 20, 0, time.UTC)
	animeReleasedOn5 := time.Date(2011, 11, 17, 20, 20, 20, 0, time.UTC)
	animePosterURL5 := "/url.jpg"
	animeScore5 := 7.55
	animeDuration5 := 50.0
	animeRating5 := "r"
	animeFranchise5 := "onepunchman"
	animeProcessed5 := false
	testAnimeDto5 := models.AnimeDTO{
		ExternalID:    animeExternalID5,
		Name:          &animeName5,
		Russian:       &animeRussianName5,
		AnimeURL:      &animeURL5,
		Kind:          &animeKind5,
		Status:        &animeStatus5,
		Epizodes:      &animeEpizodes5,
		EpizodesAired: &animeEpizodesAired5,
		AiredOn:       &animeAiredOn5,
		ReleasedOn:    &animeReleasedOn5,
		PosterURL:     &animePosterURL5,
		Score:         &animeScore5,
		Duration:      &animeDuration5,
		Rating:        &animeRating5,
		Franchise:     &animeFranchise5,
		Processed:     &animeProcessed5,
	}
	animeID5, insertAnimeErr := insertAnimeToDatabase(animeDao, testAnimeDto5)
	if insertAnimeErr != nil {
		markAsFailAndAbortNow(t, errors.Wrap(insertAnimeErr, ""))
	}
	if linkAnimeAndGenreErr := linkAnimeAndGenre(animeDao, animeID5, genreID); linkAnimeAndGenreErr != nil {
		markAsFailAndAbortNow(t, errors.Wrap(linkAnimeAndGenreErr, ""))
	}
	if linkAnimeAndStudioErr := linkAnimeAndStudio(animeDao, animeID5, studioID); linkAnimeAndStudioErr != nil {
		markAsFailAndAbortNow(t, errors.Wrap(linkAnimeAndStudioErr, ""))
	}
}

func buildAnime6(t *testing.T, animeDao *models.AnimeDAO, genreID int64, studioID int64) {
	//insert anime 6 to database
	animeExternalID6 := "6"
	animeName6 := "One Punch Man"
	animeRussianName6 := "Один Удар Человек"
	animeURL6 := "/url.jpg"
	animeKind6 := "tv"
	animeStatus6 := "ongoing"
	var animeEpizodes6 int64 = 6
	var animeEpizodesAired6 int64 = 6
	animeAiredOn6 := time.Date(2009, 11, 17, 20, 20, 20, 0, time.UTC)
	animeReleasedOn6 := time.Date(2009, 11, 17, 20, 20, 20, 0, time.UTC)
	animePosterURL6 := "/url.jpg"
	animeScore6 := 7.66
	animeDuration6 := 20.0
	animeRating6 := "r"
	animeFranchise6 := "onepunchman"
	animeProcessed6 := false
	testAnimeDto6 := models.AnimeDTO{
		ExternalID:    animeExternalID6,
		Name:          &animeName6,
		Russian:       &animeRussianName6,
		AnimeURL:      &animeURL6,
		Kind:          &animeKind6,
		Status:        &animeStatus6,
		Epizodes:      &animeEpizodes6,
		EpizodesAired: &animeEpizodesAired6,
		AiredOn:       &animeAiredOn6,
		ReleasedOn:    &animeReleasedOn6,
		PosterURL:     &animePosterURL6,
		Score:         &animeScore6,
		Duration:      &animeDuration6,
		Rating:        &animeRating6,
		Franchise:     &animeFranchise6,
		Processed:     &animeProcessed6,
	}
	animeID6, insertAnimeErr := insertAnimeToDatabase(animeDao, testAnimeDto6)
	if insertAnimeErr != nil {
		markAsFailAndAbortNow(t, errors.Wrap(insertAnimeErr, ""))
	}
	if linkAnimeAndGenreErr := linkAnimeAndGenre(animeDao, animeID6, genreID); linkAnimeAndGenreErr != nil {
		markAsFailAndAbortNow(t, errors.Wrap(linkAnimeAndGenreErr, ""))
	}
	if linkAnimeAndStudioErr := linkAnimeAndStudio(animeDao, animeID6, studioID); linkAnimeAndStudioErr != nil {
		markAsFailAndAbortNow(t, errors.Wrap(linkAnimeAndStudioErr, ""))
	}
}

func buildAnime7(t *testing.T, animeDao *models.AnimeDAO, genreID int64, studioID int64) {
	//insert anime 7 to database
	animeExternalID7 := "7"
	animeName7 := "One Punch Man"
	animeRussianName7 := "Один Удар Человек"
	animeURL7 := "/url.jpg"
	animeKind7 := "tv"
	animeStatus7 := "ongoing"
	var animeEpizodes7 int64 = 7
	var animeEpizodesAired7 int64 = 7
	animeAiredOn7 := time.Date(2009, 11, 17, 20, 20, 20, 0, time.UTC)
	animeReleasedOn7 := time.Date(2009, 11, 17, 20, 20, 20, 0, time.UTC)
	animePosterURL7 := "/url.jpg"
	animeScore7 := 7.77
	animeDuration7 := 20.0
	animeRating7 := "r"
	animeFranchise7 := "onepunchman"
	animeProcessed7 := false
	testAnimeDto7 := models.AnimeDTO{
		ExternalID:    animeExternalID7,
		Name:          &animeName7,
		Russian:       &animeRussianName7,
		AnimeURL:      &animeURL7,
		Kind:          &animeKind7,
		Status:        &animeStatus7,
		Epizodes:      &animeEpizodes7,
		EpizodesAired: &animeEpizodesAired7,
		AiredOn:       &animeAiredOn7,
		ReleasedOn:    &animeReleasedOn7,
		PosterURL:     &animePosterURL7,
		Score:         &animeScore7,
		Duration:      &animeDuration7,
		Rating:        &animeRating7,
		Franchise:     &animeFranchise7,
		Processed:     &animeProcessed7,
	}
	animeID7, insertAnimeErr := insertAnimeToDatabase(animeDao, testAnimeDto7)
	if insertAnimeErr != nil {
		markAsFailAndAbortNow(t, errors.Wrap(insertAnimeErr, ""))
	}
	if linkAnimeAndGenreErr := linkAnimeAndGenre(animeDao, animeID7, genreID); linkAnimeAndGenreErr != nil {
		markAsFailAndAbortNow(t, errors.Wrap(linkAnimeAndGenreErr, ""))
	}
	if linkAnimeAndStudioErr := linkAnimeAndStudio(animeDao, animeID7, studioID); linkAnimeAndStudioErr != nil {
		markAsFailAndAbortNow(t, errors.Wrap(linkAnimeAndStudioErr, ""))
	}
}

func buildAnime8(t *testing.T, animeDao *models.AnimeDAO, genreID int64, studioID int64) {
	//insert anime 8 to database
	animeExternalID8 := "8"
	animeName8 := "One Punch Man"
	animeRussianName8 := "Один Удар Человек"
	animeURL8 := "/url.jpg"
	animeKind8 := "tv"
	animeStatus8 := "anons"
	var animeEpizodes8 int64 = 8
	var animeEpizodesAired8 int64 = 8
	animeAiredOn8 := time.Date(2009, 11, 17, 20, 20, 20, 0, time.UTC)
	animeReleasedOn8 := time.Date(2009, 11, 17, 20, 20, 20, 0, time.UTC)
	animePosterURL8 := "/url.jpg"
	animeScore8 := 7.12
	animeDuration8 := 20.0
	animeRating8 := "r"
	animeFranchise8 := "onepunchman"
	animeProcessed8 := false
	testAnimeDto8 := models.AnimeDTO{
		ExternalID:    animeExternalID8,
		Name:          &animeName8,
		Russian:       &animeRussianName8,
		AnimeURL:      &animeURL8,
		Kind:          &animeKind8,
		Status:        &animeStatus8,
		Epizodes:      &animeEpizodes8,
		EpizodesAired: &animeEpizodesAired8,
		AiredOn:       &animeAiredOn8,
		ReleasedOn:    &animeReleasedOn8,
		PosterURL:     &animePosterURL8,
		Score:         &animeScore8,
		Duration:      &animeDuration8,
		Rating:        &animeRating8,
		Franchise:     &animeFranchise8,
		Processed:     &animeProcessed8,
	}
	animeID8, insertAnimeErr := insertAnimeToDatabase(animeDao, testAnimeDto8)
	if insertAnimeErr != nil {
		markAsFailAndAbortNow(t, errors.Wrap(insertAnimeErr, ""))
	}
	if linkAnimeAndGenreErr := linkAnimeAndGenre(animeDao, animeID8, genreID); linkAnimeAndGenreErr != nil {
		markAsFailAndAbortNow(t, errors.Wrap(linkAnimeAndGenreErr, ""))
	}
	if linkAnimeAndStudioErr := linkAnimeAndStudio(animeDao, animeID8, studioID); linkAnimeAndStudioErr != nil {
		markAsFailAndAbortNow(t, errors.Wrap(linkAnimeAndStudioErr, ""))
	}
}

func buildAnime9(t *testing.T, animeDao *models.AnimeDAO, genreID int64, studioID int64) {
	//insert anime 9 to database
	animeExternalID9 := "9"
	animeName9 := "One Punch Man"
	animeRussianName9 := "Один Удар Человек"
	animeURL9 := "/url.jpg"
	animeKind9 := "tv"
	animeStatus9 := "ongoing"
	var animeEpizodes9 int64 = 9
	var animeEpizodesAired9 int64 = 9
	animeAiredOn9 := time.Date(2009, 11, 17, 20, 20, 20, 0, time.UTC)
	animeReleasedOn9 := time.Date(2009, 11, 17, 20, 20, 20, 0, time.UTC)
	animePosterURL9 := "/url.jpg"
	animeScore9 := 7.12
	animeDuration9 := 20.0
	animeRating9 := "r"
	animeFranchise9 := "onepunchman"
	animeProcessed9 := false
	testAnimeDto9 := models.AnimeDTO{
		ExternalID:    animeExternalID9,
		Name:          &animeName9,
		Russian:       &animeRussianName9,
		AnimeURL:      &animeURL9,
		Kind:          &animeKind9,
		Status:        &animeStatus9,
		Epizodes:      &animeEpizodes9,
		EpizodesAired: &animeEpizodesAired9,
		AiredOn:       &animeAiredOn9,
		ReleasedOn:    &animeReleasedOn9,
		PosterURL:     &animePosterURL9,
		Score:         &animeScore9,
		Duration:      &animeDuration9,
		Rating:        &animeRating9,
		Franchise:     &animeFranchise9,
		Processed:     &animeProcessed9,
	}
	animeID9, insertAnimeErr := insertAnimeToDatabase(animeDao, testAnimeDto9)
	if insertAnimeErr != nil {
		markAsFailAndAbortNow(t, errors.Wrap(insertAnimeErr, ""))
	}
	if linkAnimeAndGenreErr := linkAnimeAndGenre(animeDao, animeID9, genreID); linkAnimeAndGenreErr != nil {
		markAsFailAndAbortNow(t, errors.Wrap(linkAnimeAndGenreErr, ""))
	}
	if linkAnimeAndStudioErr := linkAnimeAndStudio(animeDao, animeID9, studioID); linkAnimeAndStudioErr != nil {
		markAsFailAndAbortNow(t, errors.Wrap(linkAnimeAndStudioErr, ""))
	}
}

func buildAnime10(t *testing.T, animeDao *models.AnimeDAO, genreID int64, studioID int64) {
	//insert anime 10 to database
	animeExternalID10 := "10"
	animeName10 := "One Punch Man"
	animeRussianName10 := "Один Удар Человек"
	animeURL10 := "/url.jpg"
	animeKind10 := "movie"
	animeStatus10 := "ongoing"
	var animeEpizodes10 int64 = 10
	var animeEpizodesAired10 int64 = 10
	animeAiredOn10 := time.Date(2010, 11, 17, 20, 20, 20, 0, time.UTC)
	animeReleasedOn10 := time.Date(2010, 11, 17, 20, 20, 20, 0, time.UTC)
	animePosterURL10 := "/url.jpg"
	animeScore10 := 7.12
	animeDuration10 := 20.0
	animeRating10 := "r"
	animeFranchise10 := "onepunchman"
	animeProcessed10 := false
	testAnimeDto10 := models.AnimeDTO{
		ExternalID:    animeExternalID10,
		Name:          &animeName10,
		Russian:       &animeRussianName10,
		AnimeURL:      &animeURL10,
		Kind:          &animeKind10,
		Status:        &animeStatus10,
		Epizodes:      &animeEpizodes10,
		EpizodesAired: &animeEpizodesAired10,
		AiredOn:       &animeAiredOn10,
		ReleasedOn:    &animeReleasedOn10,
		PosterURL:     &animePosterURL10,
		Score:         &animeScore10,
		Duration:      &animeDuration10,
		Rating:        &animeRating10,
		Franchise:     &animeFranchise10,
		Processed:     &animeProcessed10,
	}
	animeID10, insertAnimeErr := insertAnimeToDatabase(animeDao, testAnimeDto10)
	if insertAnimeErr != nil {
		markAsFailAndAbortNow(t, errors.Wrap(insertAnimeErr, ""))
	}
	if linkAnimeAndGenreErr := linkAnimeAndGenre(animeDao, animeID10, genreID); linkAnimeAndGenreErr != nil {
		markAsFailAndAbortNow(t, errors.Wrap(linkAnimeAndGenreErr, ""))
	}
	if linkAnimeAndStudioErr := linkAnimeAndStudio(animeDao, animeID10, studioID); linkAnimeAndStudioErr != nil {
		markAsFailAndAbortNow(t, errors.Wrap(linkAnimeAndStudioErr, ""))
	}
}
