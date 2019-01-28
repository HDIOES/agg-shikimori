package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	fmt.Println("added basic code and dep tools to control dependencies")

	router := mux.NewRouter()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "hello, 4na")
	})
	http.Handle("/", router)
	err := http.ListenAndServe(":10045", nil)
	if err != nil {
		panic(err)
	}
}
