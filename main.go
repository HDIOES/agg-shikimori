package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"

	"strconv"

	"time"

	"github.com/robfig/cron"
	"github.com/tkanos/gonfig"
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
	http.Handle("/", router)
	listenandserveErr := http.ListenAndServe(":10045", nil)
	if listenandserveErr != nil {
		panic(err)
	}
}

type ShikimoriJob struct {
	db *sql.DB
}

func (sj *ShikimoriJob) Run() {
	client := &http.Client{}
	animes := &[]Anime{}
	page := 1
	for len(*animes) == 50 || page == 1 {
		resp, err := client.Get("https://shikimori.org/api/animes?page=" + strconv.Itoa(page) + "&limit=50")
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		json.Unmarshal(body, animes)
		for i := 0; i < len(*animes); i++ {
			sj.db.Exec("INSERT INTO anime (external_id, name) VALUES ($1, $2)", (*animes)[i].ID, (*animes)[i].Name)
		}
		page++
	}
}
