package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/HDIOES/cpa-backend/models"
	"github.com/HDIOES/cpa-backend/rest"
	"github.com/gorilla/mux"

	"net/http"

	_ "github.com/lib/pq"
	migrate "github.com/rubenv/sql-migrate"
	"go.uber.org/dig"

	"strconv"

	"time"

	"github.com/robfig/cron"
	"github.com/tkanos/gonfig"

	"github.com/gorilla/handlers"

	"github.com/HDIOES/cpa-backend/integration"
	"github.com/HDIOES/cpa-backend/rest/util"
)

//CreateDI function to build new DI container
func CreateDI() *dig.Container {
	container := dig.New()
	container.Provide(func() *util.Configuration {
		log.Println("Loading configuration...")
		configPath := "configuration.json"
		if os.Getenv("CONFIG_PATH") != "" {
			configPath = os.Getenv("CONFIG_PATH")
		}
		configuration := &util.Configuration{}
		gonfigErr := gonfig.GetConf(configPath, configuration)
		if gonfigErr != nil {
			panic(gonfigErr)
		}
		return configuration
	})
	container.Provide(func(configuration *util.Configuration) *sql.DB {
		log.Println("Prepating db...")
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
		return db
	})
	container.Provide(func(db *sql.DB) (*models.AnimeDAO, *models.GenreDAO, *models.StudioDAO, *models.NewDAO) {
		return &models.AnimeDAO{Db: db}, &models.GenreDAO{Db: db}, &models.StudioDAO{Db: db}, &models.NewDAO{Db: db}
	})
	container.Provide(func(animeDao *models.AnimeDAO, genreDao *models.GenreDAO, studioDao *models.StudioDAO, configuration *util.Configuration) *integration.ShikimoriJob {
		log.Println("Prepare shikimori job")
		return &integration.ShikimoriJob{AnimeDao: animeDao, GenreDao: genreDao, StudioDao: studioDao, Config: configuration}
	})
	container.Provide(func(
		animeDao *models.AnimeDAO,
		genreDao *models.GenreDAO,
		studioDao *models.StudioDAO,
		newDao *models.NewDAO,
		configuration *util.Configuration) (
		*rest.GenreHandler,
		*rest.CreateNewHandler,
		*rest.FindNewHandler,
		*rest.RandomAnimeHandler,
		*rest.StudioHandler,
		*rest.SearchAnimeHandler) {
		genreHandler := &rest.GenreHandler{Dao: genreDao}
		createNewHandler := &rest.CreateNewHandler{Dao: newDao}
		findNewHandler := &rest.FindNewHandler{Dao: newDao}
		randomAnimeHandler := &rest.RandomAnimeHandler{Dao: animeDao}
		studioHandler := &rest.StudioHandler{Dao: studioDao}
		searchAnimeHandler := &rest.SearchAnimeHandler{Dao: animeDao}
		return genreHandler, createNewHandler, findNewHandler, randomAnimeHandler, studioHandler, searchAnimeHandler
	})
	container.Provide(func(genreHandler *rest.GenreHandler,
		createNewHandler *rest.CreateNewHandler,
		findNewHandler *rest.FindNewHandler,
		randomAnimeHandler *rest.RandomAnimeHandler,
		studioHandler *rest.StudioHandler,
		searchAnimeHandler *rest.SearchAnimeHandler) *mux.Router {
		router := mux.NewRouter()
		router.Handle("/api/animes/random", randomAnimeHandler).
			Methods("GET")
		router.Handle("/api/animes/search", searchAnimeHandler).
			Methods("GET")
		router.Handle("/api/genres/search", genreHandler).
			Methods("GET")
		router.Handle("/api/studios/search", studioHandler).
			Methods("GET")
		router.Handle("/api/news", createNewHandler).
			Methods("POST")
		router.Handle("/api/news", findNewHandler).
			Methods("GET")
		router.Handle("/api/news", nil).
			Methods("DELETE")
		http.Handle("/", router)
		return router
	})
	return container
}

func main() {
	log.Println("Application has been runned")

	di := CreateDI()

	di.Invoke(func(job *integration.ShikimoriJob) {
		cronRunner := cron.New()
		cronRunner.AddJob("@daily", job)
		cronRunner.Start()
		log.Println("Job has been runned")
	})

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	di.Invoke(func(router *mux.Router, configuration *util.Configuration) {
		listenandserveErr := http.ListenAndServe(":"+strconv.Itoa(configuration.Port), handlers.CORS(originsOk, headersOk, methodsOk)(router))
		if listenandserveErr != nil {
			panic(listenandserveErr)
		}
	})
}
