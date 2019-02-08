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
				db.Exec("INSERT INTO anime (external_id, name) VALUES ($1, $2)", (*animes)[i].ID, (*animes)[i].Name)
			}
			page++
		}
		w.Write(nil)
	})
	http.Handle("/", router)
	listenandserveErr := http.ListenAndServe(":10045", nil)
	if listenandserveErr != nil {
		panic(err)
	}
}
