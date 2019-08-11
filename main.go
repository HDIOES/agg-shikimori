package main

import (
	"log"

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

	di := di.CreateDI("configuration.json", false)

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
