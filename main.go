package main

import (
	"log"
	"os"

	"github.com/gorilla/mux"

	"net/http"

	_ "github.com/lib/pq"

	"strconv"

	"github.com/robfig/cron"

	"github.com/HDIOES/agg-shikimori/di"
	"github.com/HDIOES/agg-shikimori/integration"
	"github.com/HDIOES/agg-shikimori/rest/util"
)

//CreateDI function to build new DI container
func main() {
	log.Println("Application has been runned")
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "configuration.json"
	}
	di := di.CreateDI(configPath, "migrations", false)
	di.Invoke(func(job *integration.ShikimoriJob) {
		cronRunner := cron.New()
		cronRunner.AddJob("@daily", job)
		cronRunner.Start()
		log.Println("Job has been runned")
	})

	di.Invoke(func(router *mux.Router, configuration *util.Configuration) {
		listenandserveErr := http.ListenAndServe(":"+strconv.Itoa(configuration.Port), router)
		if listenandserveErr != nil {
			panic(listenandserveErr)
		}
	})
}
