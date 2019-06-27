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
)

type Configuration struct {
	DatabaseUrl        string `json:"databaseUrl"`
	MaxOpenConnections int    `json:"maxOpenConnections"`
	MaxIdleConnections int    `json:"maxIdleConnections"`
	ConnectionTimeout  int    `json:"connectionTimeout"`
	Port               int    `json:"port"`
}

func main() {
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Llongfile)
	log.Println("Application has been runned")
	log.Println("Loading configuration...")

	file, err := os.OpenFile("cpa.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open log file:", err)
	}
	log.SetOutput(file)

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
	shikimoriJob := &integration.ShikimoriJob{Db: db}
	cronRunner.AddJob("@daily", shikimoriJob)
	cronRunner.Start()

	log.Println("Job has been runned")
	pingErr := db.Ping()
	if pingErr != nil {
		panic(pingErr)
	}

	router := mux.NewRouter()

	router.Handle("/animes/random", animes.CreateAnimeHandler(db)).
		Methods("GET")

	router.Handle("/animes/search", animes.CreateSearchAnimeHandler(db, router)).
		Methods("GET")

	router.HandleFunc("/animes/job", func(w http.ResponseWriter, r *http.Request) {
		go shikimoriJob.Run()
	}).Methods("GET")

	router.Handle("/genres/search", genres.CreateGenreHandler(db)).
		Methods("GET")

	router.Handle("/studios/search", studios.CreateStudioHandler(db)).
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
