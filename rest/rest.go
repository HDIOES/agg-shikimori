package rest

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/HDIOES/cpa-backend/rest/util"
	"github.com/gorilla/mux"
)

func ConfigureRouter(db *sql.DB, configuration *util.Configuration) *mux.Router {
	router := mux.NewRouter()
	router.Handle("/api/animes/random", CreateRandomAnimeHandler(db, configuration)).
		Methods("GET")
	router.Handle("/api/animes/search", CreateSearchAnimeHandler(db, router, configuration)).
		Methods("GET")
	router.Handle("/api/genres/search", CreateGenreHandler(db)).
		Methods("GET")
	router.Handle("/api/studios/search", CreateStudioHandler(db)).
		Methods("GET")
	router.Handle("/api/news", CreateCreateNewHandler(db)).
		Methods("POST")
	router.Handle("/api/news", CreateFindNewHandler(db)).
		Methods("GET")
	router.Handle("/api/news", nil).
		Methods("DELETE")
	http.Handle("/", router)
	return router
}

func HandleErr(err error, w http.ResponseWriter, httpStatus int, errorMessage string) {
	log.Println(err)
	errorMessageBuilder := strings.Builder{}
	errorMessageBuilder.WriteString("{")
	errorMessageBuilder.WriteString("\"message\":")
	errorMessageBuilder.WriteString("\"")
	errorMessageBuilder.WriteString(errorMessage)
	errorMessageBuilder.WriteString("\"")
	errorMessageBuilder.WriteString("}")
	ReturnResponseAsJSON(w, errorMessageBuilder.String(), httpStatus)
}

func ReturnResponseAsJSON(w http.ResponseWriter, body interface{}, httpStatus int) {
	w.WriteHeader(httpStatus)
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(body)
}
