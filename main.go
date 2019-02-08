package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {

	fmt.Println("added basic code and dep tools to control dependencies")
	db, err := sql.Open("postgres", "postgres://forna_user:12345@localhost:5432/forna")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	pingErr := db.Ping()
	if pingErr != nil {
		panic(pingErr)
	}

	router := mux.NewRouter()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "hello, 4na")
	})
	router.HandleFunc("/animes", func(w http.ResponseWriter, r *http.Request) {
		client := &http.Client{}
		resp, err := client.Get("https://shikimori.org/api/animes")
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		animes := &[]Anime{}
		json.Unmarshal(body, animes)
		db.Exec("INSERT INTO ANIME (ID, NAME) VALUES ($1, $2)", (*animes)[0].ID, (*animes)[0].Name)
		w.Write(body)
	})
	http.Handle("/", router)
	listenandserveErr := http.ListenAndServe(":10045", nil)
	if listenandserveErr != nil {
		panic(err)
	}
}
