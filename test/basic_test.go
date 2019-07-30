package test

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/gorilla/mux"

	"github.com/HDIOES/cpa-backend/rest"
	"github.com/HDIOES/cpa-backend/rest/util"
	"github.com/ory/dockertest"
	migrate "github.com/rubenv/sql-migrate"
	"github.com/tkanos/gonfig"
)

//TODO remove global vars!!!!
var db *sql.DB
var configuration *util.Configuration
var testContainer *dockertest.Resource
var router *mux.Router

func TestMain(m *testing.M) {
	//prepare test database, test configuration and test router
	os.Exit(wrapperTestMain(m))
}

func wrapperTestMain(m *testing.M) int {
	if preparedConfiguration, err := prepareConfiguration(); err != nil {
		panic(err)
	} else {
		configuration = preparedConfiguration
	}
	if preparedDb, preparedTestContainer, err := prepareDb(); err != nil {
		panic(err)
	} else {
		defer preparedTestContainer.Close()
		defer preparedDb.Close()
		testContainer = preparedTestContainer
		db = preparedDb
	}
	router = rest.ConfigureRouter(db, configuration)
	defer log.Print("Stopping test container")
	return m.Run()
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	return rr
}

func insertAnimeToDatabase(t *testing.T,
	externalID, animeName, russian, animeURL, kind, animeStatus string,
	epizodes, epizodesAired int,
	airedOn, releasedOn time.Time,
	posterURL string,
	processed bool,
	externalStudioID, externalGenreID string) error {
	tx, beginErr := db.Begin()
	if beginErr != nil {
		return beginErr
	}
	//insert anime
	_, insertAnimeErr := tx.Exec("INSERT INTO anime (external_id, name, russian, amine_url, kind, anime_status, epizodes, epizodes_aired, aired_on, released_on, poster_url, processed, lastmodifytime) "+
		"VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, now())",
		externalID,
		animeName,
		russian,
		animeURL,
		kind,
		animeStatus,
		epizodes,
		epizodesAired,
		airedOn.Format("2006-01-02"),
		releasedOn.Format("2006-01-02"),
		posterURL,
		processed)
	if insertAnimeErr != nil {
		return rollbackTransaction(tx, insertAnimeErr)
	}
	//insert anime_studio
	_, insertAnimeStudioErr := tx.Exec("INSERT INTO anime_studio (anime_id, studio_id)"+
		" SELECT anime.id, studio.id FROM anime JOIN studio ON anime.external_id = $1 AND studio.external_id = $2", externalID, externalStudioID)
	if insertAnimeStudioErr != nil {
		return rollbackTransaction(tx, insertAnimeStudioErr)
	}
	//insert anime_genre
	_, insertAnimeGenreErr := tx.Exec("INSERT INTO anime_genre (anime_id, genre_id) "+
		" SELECT anime.id, genre.id FROM anime JOIN genre ON anime.external_id = $1 AND genre.external_id = $2", externalID, externalGenreID)
	if insertAnimeGenreErr != nil {
		return rollbackTransaction(tx, insertAnimeGenreErr)
	}
	if commitErr := tx.Commit(); commitErr != nil {
		return rollbackTransaction(tx, commitErr)
	}
	return nil
}

func insertStudioToDatabase(t *testing.T, externalID, studioName, filteredStudioName string, isReal bool, imageURL string) error {
	tx, beginErr := db.Begin()
	if beginErr != nil {
		return beginErr
	}
	_, txErr := tx.Exec("INSERT INTO studio (external_id, studio_name, filtered_studio_name, is_real, image_url) VALUES ($1, $2, $3, $4, $5)",
		externalID,
		studioName,
		filteredStudioName,
		isReal,
		imageURL)
	if txErr != nil {
		return rollbackTransaction(tx, txErr)
	}
	if commitErr := tx.Commit(); commitErr != nil {
		return rollbackTransaction(tx, commitErr)
	}
	return nil
}

func insertGenreToDatabase(t *testing.T, externalID, genreName, russian, kind string) error {
	tx, beginErr := db.Begin()
	if beginErr != nil {
		return beginErr
	}
	_, txErr :=
		tx.Exec("INSERT INTO genre (external_id, genre_name, russian, kind) VALUES ($1, $2, $3, $4)", externalID, genreName, russian, kind)
	if txErr != nil {
		return rollbackTransaction(tx, txErr)
	}
	if commitErr := tx.Commit(); commitErr != nil {
		return rollbackTransaction(tx, commitErr)
	}
	return nil
}

func rollbackTransaction(tx *sql.Tx, err error) error {
	if rollbackErr := tx.Rollback(); rollbackErr != nil {
		return rollbackErr
	}
	return err
}

//prepareConfiguration function returns config data from file named "configuration.json"
func prepareConfiguration() (configuration *util.Configuration, err error) {
	configuration = &util.Configuration{}
	gonfigErr := gonfig.GetConf("../configuration.json", configuration)
	if gonfigErr != nil {
		return nil, gonfigErr
	}
	return configuration, nil

}

//prepareDb function prepares data container for using in local testing
//or uses already prepared data container in case of using docker-compose testing
func prepareDb() (*sql.DB, *dockertest.Resource, error) {
	//start up new test data container
	pool, err := dockertest.NewPool("")
	if err != nil {
		return nil, nil, err
	}
	resource, rErr := pool.Run(dbType, dbVersion, []string{
		dbUserVar,
		dbPasswordVar,
		dbNameVar})
	log.Print("Starting test container...")
	time.Sleep(10 * time.Second)
	if rErr != nil {
		defer resource.Close()
		return nil, nil, rErr
	}
	log.Print("Port - " + resource.GetPort(dbPortMapping))
	databaseURL := fmt.Sprintf(dbURLTemplate, resource.GetPort(dbPortMapping))
	//create db
	preparedDB, err := sql.Open(dbType, databaseURL)
	if err != nil {
		defer testContainer.Close()
		defer preparedDB.Close()
		return nil, nil, err
	}
	preparedDB.SetMaxIdleConns(configuration.MaxIdleConnections)
	preparedDB.SetMaxOpenConns(configuration.MaxOpenConnections)
	timeout := strconv.Itoa(configuration.ConnectionTimeout) + "s"
	timeoutDuration, durationErr := time.ParseDuration(timeout)
	if durationErr != nil {
		return nil, nil, durationErr
	}
	preparedDB.SetConnMaxLifetime(timeoutDuration)
	err = applyMigrations(preparedDB)
	if err != nil {
		return nil, nil, err
	}
	return preparedDB, resource, nil
}

func clearDb(db *sql.DB) error {
	tx, txErr := db.Begin()
	if txErr != nil {
		return txErr
	}
	if _, err := tx.Exec("DELETE FROM ANIME_STUDIO"); err != nil {
		return rollbackTransaction(tx, err)
	}
	if _, err := tx.Exec("DELETE FROM STUDIO"); err != nil {
		return rollbackTransaction(tx, err)
	}
	if _, err := tx.Exec("DELETE FROM ANIME_GENRE"); err != nil {
		return rollbackTransaction(tx, err)
	}
	if _, err := tx.Exec("DELETE FROM ANIME"); err != nil {
		return rollbackTransaction(tx, err)
	}
	if txCommitErr := tx.Commit(); txCommitErr != nil {
		return rollbackTransaction(tx, txCommitErr)
	}
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
