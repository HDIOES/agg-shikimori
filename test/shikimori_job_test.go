package test

import (
	"database/sql"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/HDIOES/cpa-backend/integration"
	"github.com/HDIOES/cpa-backend/util"
	_ "github.com/lib/pq"
	"github.com/tkanos/gonfig"

	gock "gopkg.in/h2non/gock.v1"
)

var db *sql.DB
var configuration *util.Configuration

func TestMain(m *testing.M) {
	db, configuration = prepareDbAndConfiguration(m)
	code := m.Run()
	os.Exit(code)
}

func prepareDbAndConfiguration(m *testing.M) (*sql.DB, *util.Configuration) {
	//create config
	configuration := util.Configuration{}
	gonfigErr := gonfig.GetConf("../configuration.json", &configuration)
	if gonfigErr != nil {
		log.Fatal(gonfigErr)
	}
	//create db
	db, err := sql.Open("postgres", configuration.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	db.SetMaxIdleConns(configuration.MaxIdleConnections)
	db.SetMaxOpenConns(configuration.MaxOpenConnections)
	timeout := strconv.Itoa(configuration.ConnectionTimeout) + "s"
	timeoutDuration, durationErr := time.ParseDuration(timeout)
	if durationErr != nil {
		log.Fatal(durationErr)
	} else {
		db.SetConnMaxLifetime(timeoutDuration)
	}
	return db, &configuration
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

func postTest(db *sql.DB, t *testing.T) {
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestShikimoriJobSuccess1(t *testing.T) {
}

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
