package util

import (
	"fmt"
	"log"

	"github.com/pkg/errors"
)

//Configuration struct
type Configuration struct {
	DatabaseURL             string `json:"databaseUrl"`
	MaxOpenConnections      int    `json:"maxOpenConnections"`
	MaxIdleConnections      int    `json:"maxIdleConnections"`
	ConnectionTimeout       int    `json:"connectionTimeout"`
	Port                    int    `json:"port"`
	ShikimoriURL            string `json:"shikimori_url"`
	ShikimoriAnimeSearchURL string `json:"shikimori_anime_search_url"`
	ShikimoriGenreURL       string `json:"shikimori_genre_url"`
	ShikimoriStudioURL      string `json:"shikimori_studio_url"`
}

//StackTracer struct
type StackTracer interface {
	StackTrace() errors.StackTrace
}

//HandleError func
func HandleError(err error) {
	if err, ok := err.(StackTracer); ok {
		for _, f := range err.StackTrace() {
			fmt.Printf("%+s:%d\n", f, f)
		}
	} else {
		log.Println("Unknown error: ", err)
	}
}
