package main

import (
	"database/sql"
	"log"
	"os"

	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	migrate "github.com/rubenv/sql-migrate"

	"strconv"

	"time"

	"github.com/robfig/cron"
	"github.com/tkanos/gonfig"

	"github.com/gorilla/handlers"

	"github.com/HDIOES/cpa-backend/integration"
	"github.com/HDIOES/cpa-backend/rest/animes"
	"github.com/HDIOES/cpa-backend/rest/genres"
	"github.com/HDIOES/cpa-backend/rest/studios"
	"github.com/HDIOES/cpa-backend/util"
)

func main() {
	log.Println("Application has been runned")
	log.Println("Loading configuration...")
	configPath := "configuration.json"
	if len(os.Args) > 1 {
		configPath = os.Args[1]
	}

	configuration := &util.Configuration{}
	gonfigErr := gonfig.GetConf(configPath, configuration)
	if gonfigErr != nil {
		panic(gonfigErr)
	}
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

	migrations := &migrate.FileMigrationSource{
		Dir: "migrations",
	}

	if n, err := migrate.Exec(db, "postgres", migrations, migrate.Up); err == nil {
		log.Printf("Applied %d migrations!\n", n)
	} else {
		panic(err)
	}

	log.Println("Configuration has been loaded")
	defer db.Close()
	log.Println("Job running...")
	cronRunner := cron.New()
	shikimoriJob := &integration.ShikimoriJob{Db: db, Config: configuration}
	cronRunner.AddJob("@daily", shikimoriJob)
	cronRunner.Start()

	log.Println("Job has been runned")
	pingErr := db.Ping()
	if pingErr != nil {
		panic(pingErr)
	}

	router := mux.NewRouter()

	router.Handle("/api/animes/random", animes.CreateRandomAnimeHandler(db, configuration)).
		Methods("GET")

	router.Handle("/api/animes/search", animes.CreateSearchAnimeHandler(db, router, configuration)).
		Methods("GET")

	router.HandleFunc("/api/animes/job", func(w http.ResponseWriter, r *http.Request) {
		go shikimoriJob.Run()
	}).Methods("GET")

	router.Handle("/api/genres/search", genres.CreateGenreHandler(db)).
		Methods("GET")

	router.Handle("/api/studios/search", studios.CreateStudioHandler(db)).
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
