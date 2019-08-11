package test

import (
	"database/sql"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/HDIOES/cpa-backend/models"

	"github.com/gorilla/mux"

	"go.uber.org/dig"

	"github.com/HDIOES/cpa-backend/di"
	"github.com/ory/dockertest"
	migrate "github.com/rubenv/sql-migrate"
)

var diContainer *dig.Container

func init() {
	diContainer = di.CreateDI(true)
}

func TestMain(m *testing.M) {
	//prepare test database, test configuration and test router
	os.Exit(wrapperTestMain(m))
}

func wrapperTestMain(m *testing.M) int {
	defer diContainer.Invoke(func(db *sql.DB, testContainer *dockertest.Resource) {
		db.Close()
		testContainer.Close()
	})
	defer log.Print("Stopping test container")
	return m.Run()
}

func executeRequest(req *http.Request, router *mux.Router) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	return rr
}

func linkAnimeAndStudio(animeDao *models.AnimeDAO, animeID int64, studioID int64) error {
	return animeDao.LinkAnimeAndStudio(animeID, studioID)
}

func linkAnimeAndGenre(animeDao *models.AnimeDAO, animeID int64, genreID int64) error {
	return animeDao.LinkAnimeAndGenre(animeID, genreID)
}

func insertAnimeToDatabase(
	animeDao *models.AnimeDAO,
	externalID,
	animeName,
	russian,
	animeURL,
	kind,
	animeStatus string,
	epizodes, epizodesAired int64,
	airedOn, releasedOn time.Time,
	posterURL string,
	processed bool) (int64, error) {
	animeDto := models.AnimeDTO{
		Name:          animeName,
		ExternalID:    externalID,
		Russian:       russian,
		AnimeURL:      animeURL,
		Kind:          kind,
		Status:        animeStatus,
		Epizodes:      epizodes,
		EpizodesAired: epizodesAired,
		AiredOn:       airedOn,
		ReleasedOn:    releasedOn,
		PosterURL:     posterURL,
	}
	id, err := animeDao.Create(animeDto)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func insertStudioToDatabase(studioDao *models.StudioDAO,
	externalID,
	studioName,
	filteredStudioName string,
	isReal bool,
	imageURL string) (int64, error) {
	id, err := studioDao.Create(models.StudioDTO{
		Name:               studioName,
		ExternalID:         externalID,
		FilteredStudioName: filteredStudioName,
		IsReal:             isReal,
		ImageURL:           imageURL,
	})
	if err != nil {
		return 0, err
	}
	return id, nil
}

func insertGenreToDatabase(genreDao *models.GenreDAO, externalID, genreName, russian, kind string) (int64, error) {
	id, err := genreDao.Create(models.GenreDTO{
		ExternalID: externalID,
		Name:       genreName,
		Russian:    russian,
		Kind:       kind,
	})
	if err != nil {
		return 0, err
	}
	return id, nil
}

func insertNewToDatabase(newDao *models.NewDAO, name string, body string) (int64, error) {
	id, err := newDao.Create(models.NewDTO{
		Name: name,
		Body: body,
	})
	if err != nil {
		return 0, err
	}
	return id, nil
}

func clearDb(newDao *models.NewDAO, animeDao *models.AnimeDAO, genreDao *models.GenreDAO, studioDao *models.StudioDAO) error {
	newDao.DeleteAll()
	genreDao.DeleteAll()
	studioDao.DeleteAll()
	animeDao.DeleteAll()
	return nil
}

func applyMigrations(db *sql.DB) error {
	migrations := &migrate.FileMigrationSource{
		Dir: "../migrations",
	}
	if n, err := migrate.Exec(db, "postgres", migrations, migrate.Up); err == nil {
		log.Printf("Applied %d migrations!\n", n)
	} else {
		return err
	}
	return nil
}
