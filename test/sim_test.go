package test

import (
	"database/sql"
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

//TestSimple function
func TestSimple(t *testing.T) {
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
	defer db.Close()

	shikimoriJob := &integration.ShikimoriJob{Db: db, Config: configuration}
	gock.New(configuration.ShikimoriURL).
		Get(configuration.ShikimoriAnimeSearchURL).
		MatchParam("page", "1").
		MatchParam("limit", "50").
		Reply(500)

	shikimoriJob.Run()

}
