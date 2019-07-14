package test

import (
	"database/sql"
	"io/ioutil"
	"log"
	"strconv"
	"testing"
	"time"

	"github.com/HDIOES/cpa-backend/integration"
	"github.com/HDIOES/cpa-backend/util"
	_ "github.com/lib/pq"
	"github.com/tkanos/gonfig"

	gock "gopkg.in/h2non/gock.v1"
)

func preTest() (*sql.DB, *util.Configuration) {
	//create config
	configuration := util.Configuration{}
	gonfigErr := gonfig.GetConf("../configuration.json", &configuration)
	if gonfigErr != nil {
		panic(gonfigErr)
	}
	//create db
	db, err := sql.Open("postgres", configuration.DatabaseURL)
	if err != nil {
		panic(err)
	}
	db.SetMaxIdleConns(configuration.MaxIdleConnections)
	db.SetMaxOpenConns(configuration.MaxOpenConnections)
	timeout := strconv.Itoa(configuration.ConnectionTimeout) + "s"
	timeoutDuration, durationErr := time.ParseDuration(timeout)
	if durationErr != nil {
		log.Println("Error parsing of timeout parameter")
		panic(durationErr)
	} else {
		db.SetConnMaxLifetime(timeoutDuration)
	}
	//clear db
	clearDb(db)
	return db, &configuration
}

func clearDb(db *sql.DB) (dbErr error) {
	tx, txErr := db.Begin()
	if txErr != nil {
		log.Println("Transaction start failed")
		return txErr
	}
	defer func(tx *sql.Tx) {
		if dbErr != nil {
			if err := tx.Rollback(); err != nil {
				log.Println("Transaction cannot rollbaked")
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
		log.Println("Transaction cannot be commited")
		dbErr = txCommitErr
		return dbErr
	}
	return nil
}

func postTest(db *sql.DB) {
	db.Close()
}

//TestSimple function
func TestSimple(t *testing.T) {
	defer gock.Off()
	db, config := preTest()
	defer postTest(db)
	shikimoriJob := &integration.ShikimoriJob{Db: db, Config: *config}

	animesData, err := ioutil.ReadFile("mock/shikimori_animes_success.json")
	if err != nil {
		//test error
	}
	gock.New(config.ShikimoriURL).
		Get(config.ShikimoriAnimeSearchURL).
		MatchParam("page", "1").
		MatchParam("limit", "50").
		Reply(200).
		JSON(animesData)

	genresData, err := ioutil.ReadFile("mock/shikimori_genres_success.json")
	if err != nil {
		//test error
	}
	gock.New(config.ShikimoriURL).
		Get(config.ShikimoriGenreURL).
		Reply(200).
		JSON(genresData)

	studiosData, err := ioutil.ReadFile("mock/shikimori_studios_success.json")
	if err != nil {
		//test error
	}
	gock.New(config.ShikimoriURL).
		Get(config.ShikimoriStudioURL).
		Reply(200).
		JSON(studiosData)

	oneAnimeData, err := ioutil.ReadFile("mock/one_anime_shikimori_success.json")
	if err != nil {
		//test error
	}
	gock.New(config.ShikimoriURL).
		Get(config.ShikimoriAnimeSearchURL+"/").
		PathParam("animes", "5114").
		Reply(200).
		JSON(oneAnimeData)

	shikimoriJob.Run()
}
