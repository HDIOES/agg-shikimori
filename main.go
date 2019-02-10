package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"

	"strconv"

	"time"

	"github.com/robfig/cron"
	"github.com/tkanos/gonfig"

	"math/rand"
)

type Configuration struct {
	DatabaseUrl        string `json:"databaseUrl"`
	MaxOpenConnections int    `json:"maxOpenConnections"`
	MaxIdleConnections int    `json:"maxIdleConnections"`
	ConnectionTimeout  int    `json:"connectionTimeout"`
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
	shikimoriJob := &ShikimoriJob{db: db}
	cronRunner.AddJob("@daily", shikimoriJob)
	cronRunner.Start()

	fmt.Println("Job has been runned")
	pingErr := db.Ping()
	if pingErr != nil {
		panic(pingErr)
	}

	router := mux.NewRouter()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "hello, 4na")
	})
	router.HandleFunc("/animes", func(w http.ResponseWriter, r *http.Request) {
		shikimoriJob.Run()
	})
	router.HandleFunc("/animes/random", func(w http.ResponseWriter, r *http.Request) {
		rows, queryErr := db.Query("SELECT COUNT(*) FROM anime")
		if queryErr != nil {
			fmt.Println(queryErr)
		}
		defer rows.Close()
		var count sql.NullInt64
		if rows.Next() {
			err := rows.Scan(&count)
			if err != nil {
				fmt.Println(err)
			}
		}
		randowRowNumber := rand.Int63n(count.Int64) + 1
		animeRows, animeRowsErr := db.Query("select russian, amine_url from (select row_number() over(), russian, amine_url from anime) as query where query.row_number = $1", randowRowNumber)
		if animeRowsErr != nil {
			fmt.Println(animeRowsErr)
		}
		defer animeRows.Close()
		animeRo := &AnimeRO{}
		if animeRows.Next() {
			var russianName sql.NullString
			var animeURL sql.NullString
			animeRows.Scan(&russianName, &animeURL)
			animeRo.Name = russianName.String
			animeRo.URL = "https://shikimori.ru" + animeURL.String
		}
		json.NewEncoder(w).Encode(animeRo)
	})
	http.Handle("/", router)
	listenandserveErr := http.ListenAndServe(":10045", nil)
	if listenandserveErr != nil {
		panic(err)
	}
}

//AnimeRO is rest object
type AnimeRO struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
