package test

import (
	"database/sql"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
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
	preparedDb, preparedConfiguration, preparedTestContainer, err := prepareDbAndConfiguration(m)
	db = preparedDb
	configuration = preparedConfiguration
	testContainer = preparedTestContainer
	if err != nil {
		log.Fatal(err)
	}
	code := m.Run()
	testContainer.Close()
	os.Exit(code)
}

func postTest(db *sql.DB, t *testing.T) {
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
}

//prepareDbAndConfiguration function prepares data container for using in local testing
//or uses already prepared data container in case of using docker-compose testing
func prepareDbAndConfiguration(m *testing.M) (db *sql.DB, configuration *util.Configuration, resource *dockertest.Resource, e error) {
	//create config
	configuration = &util.Configuration{}
	gonfigErr := gonfig.GetConf("configuration-test.json", configuration)
	if gonfigErr != nil {
		return nil, nil, nil, gonfigErr
	}
	if strings.Compare(os.Args[1], "docker") == 0 {
		//use already prepared container

	} else {
		//start up new test data container
		pool, err := dockertest.NewPool("")
		if err != nil {
			return nil, nil, nil, err
		} else {
			res, rErr := pool.Run("postgres", "11.4", []string{
				"POSTGRES_USER=test_forna_user",
				"POSTGRES_PASSWORD=12345",
				"POSTGRES_DB=test_forna"})
			time.Sleep(10 * time.Second)
			resource = res
			postgresURL := "postgres://test_forna_user:12345@localhost:" + res.GetPort("5432/tcp") + "/test_forna?sslmode=disable"
			configuration.DatabaseURL = postgresURL
			if rErr != nil {
				return nil, nil, resource, rErr
			}
		}
	}
	//create db
	db, err := sql.Open("postgres", configuration.DatabaseURL)
	if err != nil {
		return nil, nil, resource, err
	}
	db.SetMaxIdleConns(configuration.MaxIdleConnections)
	db.SetMaxOpenConns(configuration.MaxOpenConnections)
	timeout := strconv.Itoa(configuration.ConnectionTimeout) + "s"
	timeoutDuration, durationErr := time.ParseDuration(timeout)
	if durationErr != nil {
		return nil, nil, resource, durationErr
	} else {
		db.SetConnMaxLifetime(timeoutDuration)
	}
	err = applyMigrations(db)
	if err != nil {
		return nil, nil, resource, err
	}
	return db, configuration, resource, nil
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
