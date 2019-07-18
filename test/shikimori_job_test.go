package test

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/HDIOES/cpa-backend/integration"
	"github.com/HDIOES/cpa-backend/util"
	_ "github.com/lib/pq"
	migrate "github.com/rubenv/sql-migrate"
	"github.com/tkanos/gonfig"

	"github.com/ory/dockertest"
	gock "gopkg.in/h2non/gock.v1"
)

//TODO remove global vars!!! it's not good practice
var db *sql.DB
var configuration *util.Configuration
var testContainer *dockertest.Resource

//TestSimple function
func TestShikimoriJobSuccess(t *testing.T) {
	clearDb(db, t)
	defer gock.Off()
	defer postTest(db, t)
	shikimoriJob := &integration.ShikimoriJob{Db: db, Config: configuration}

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

	shikimoriJob.Run()
}

func TestMain(m *testing.M) {
	errors := []error{}
	if preparedConfiguration, err := prepareConfiguration(); err != nil {
		errors = append(errors, err)
	} else {
		configuration = preparedConfiguration
	}
	if preparedDb, preparedTestContainer, err := prepareDb(m); err != nil {
		errors = append(errors, err)
	} else {
		testContainer = preparedTestContainer
		db = preparedDb
	}
	code := m.Run()
	log.Print("Stopping test container")
	testContainer.Close()
	log.Fatal(errors)
	os.Exit(code)
}

func postTest(db *sql.DB, t *testing.T) {
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
}

//prepareConfiguration function returns config data from file named "configuration.json"
func prepareConfiguration() (configuration *util.Configuration, err error) {
	configuration = &util.Configuration{}
	gonfigErr := gonfig.GetConf("../configuration.json", configuration)
	if gonfigErr != nil {
		return nil, gonfigErr
	} else {
		return
	}
}

//prepareDb function prepares data container for using in local testing
//or uses already prepared data container in case of using docker-compose testing
func prepareDb(m *testing.M) (*sql.DB, *dockertest.Resource, error) {
	//start up new test data container
	pool, err := dockertest.NewPool("")
	if err != nil {
		return nil, nil, err
	} else {
		resource, rErr := pool.Run(dbType, dbVersion, []string{
			dbUserVar,
			dbPasswordVar,
			dbNameVar})
		log.Print("Starting test container...")
		time.Sleep(10 * time.Second)
		if rErr != nil {
			return nil, nil, rErr
		}
		databaseURL := fmt.Sprintf(dbURLTemplate, resource.GetPort(dbPortMapping))
		//create db
		preparedDB, err := sql.Open(dbType, databaseURL)
		if err != nil {
			return nil, nil, err
		}
		preparedDB.SetMaxIdleConns(configuration.MaxIdleConnections)
		preparedDB.SetMaxOpenConns(configuration.MaxOpenConnections)
		timeout := strconv.Itoa(configuration.ConnectionTimeout) + "s"
		timeoutDuration, durationErr := time.ParseDuration(timeout)
		if durationErr != nil {
			return nil, nil, durationErr
		} else {
			preparedDB.SetConnMaxLifetime(timeoutDuration)
		}
		err = applyMigrations(preparedDB)
		if err != nil {
			return preparedDB, nil, err
		}
		return preparedDB, resource, nil
	}
}

func clearDb(db *sql.DB, t *testing.T) (dbErr error) {
	tx, txErr := db.Begin()
	if txErr != nil {
		t.Log(txErr)
		return txErr
	}
	defer func(tx *sql.Tx) {
		if dbErr != nil {
			t.Log(dbErr)
			if err := tx.Rollback(); err != nil {
				t.Log(err)
			}
		}
	}(tx)
	if _, err := tx.Exec("DELETE FROM ANIME_STUDIO"); err != nil {
		dbErr = err
		return dbErr
	}
	if _, err := tx.Exec("DELETE FROM STUDIO"); err != nil {
		dbErr = err
		return dbErr
	}
	if _, err := tx.Exec("DELETE FROM ANIME_GENRE"); err != nil {
		dbErr = err
		return dbErr
	}
	if _, err := tx.Exec("DELETE FROM ANIME"); err != nil {
		dbErr = err
		return dbErr
	}
	if txCommitErr := tx.Commit(); txCommitErr != nil {
		dbErr = txCommitErr
		return dbErr
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
