package main

import (
	"database/sql"

	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"

	"strconv"

	"time"

	"github.com/robfig/cron"
	"github.com/tkanos/gonfig"

	"github.com/gorilla/handlers"

	"github.com/HDIOES/cpa-backend/integration"
	animes "github.com/HDIOES/cpa-backend/rest/animes"
)

type Configuration struct {
	DatabaseUrl        string `json:"databaseUrl"`
	MaxOpenConnections int    `json:"maxOpenConnections"`
	MaxIdleConnections int    `json:"maxIdleConnections"`
	ConnectionTimeout  int    `json:"connectionTimeout"`
	Port               int    `json:"port"`
}

func main() {
	fmt.Println("Application has been runned")
	fmt.Println("Loading configuration...")
	configuration := Configuration{}
	gonfigErr := gonfig.GetConf("configuration.json", &configuration)
	if gonfigErr != nil {
		panic(gonfigErr)
	}
	db, err := sql.Open("postgres", configuration.DatabaseUrl)
	if err != nil {
		panic(err)
	}
	db.SetMaxIdleConns(configuration.MaxIdleConnections)
	db.SetMaxOpenConns(configuration.MaxOpenConnections)
	timeout := strconv.Itoa(configuration.ConnectionTimeout) + "s"
	timeoutDuration, durationErr := time.ParseDuration(timeout)
	if durationErr != nil {
		fmt.Println("Error parsing of timeout parameter")
		panic(durationErr)
	} else {
		db.SetConnMaxLifetime(timeoutDuration)
	}

	fmt.Println("Configuration has been loaded")
	defer db.Close()
	fmt.Println("Job running...")
	cronRunner := cron.New()
	shikimoriJob := &integration.ShikimoriJob{Db: db}
	cronRunner.AddJob("@daily", shikimoriJob)
	cronRunner.Start()

	fmt.Println("Job has been runned")
	pingErr := db.Ping()
	if pingErr != nil {
		panic(pingErr)
	}

	router := mux.NewRouter()

	router.HandleFunc("/animes", func(w http.ResponseWriter, r *http.Request) {
		shikimoriJob.Run()
	}).Methods("GET")

	router.Handle("/animes/random", animes.CreateAnimeHandler(db)).
		Methods("GET")

	router.Handle("/animes/search", animes.CreateSearchAnimeHandler(db, router)).
		Methods("GET")

	http.Handle("/", router)
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	listenandserveErr := http.ListenAndServe(":"+strconv.Itoa(configuration.Port), handlers.CORS(originsOk, headersOk, methodsOk)(router))
	if listenandserveErr != nil {
		panic(err)
	}
}
