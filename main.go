package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {

	fmt.Println("added basic code and dep tools to control dependencies")

	router := mux.NewRouter()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		db, err := sql.Open("postgres", "postgres://forna_user:12345@localhost:5432/forna")
		err1 := db.Ping()
		if err != nil {
			panic(err)
		}
		if err1 != nil {
			panic(err1)
		}
		defer db.Close()
		fmt.Fprint(w, "hello, 4na")
	})
	http.Handle("/", router)
	err := http.ListenAndServe(":10045", nil)
	if err != nil {
		panic(err)
	}
}
