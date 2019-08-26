package main

import (
	"log"
	"os"

	"github.com/gorilla/mux"

	"net/http"

	_ "github.com/lib/pq"

	"strconv"

	"github.com/robfig/cron"

	"github.com/gorilla/handlers"

	"github.com/HDIOES/cpa-backend/di"
	"github.com/HDIOES/cpa-backend/integration"
	"github.com/HDIOES/cpa-backend/rest/util"
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
		listenandserveErr := http.ListenAndServe(":"+strconv.Itoa(configuration.Port), handlers.CombinedLoggingHandler(os.Stdout, router))
		if listenandserveErr != nil {
			panic(listenandserveErr)
		}
	})
}
