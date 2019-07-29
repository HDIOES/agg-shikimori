package test

import (
	"log"
	"net/http"
	"testing"
	"time"

	_ "github.com/lib/pq"
)

func TestSearchAnimesSuccess(t *testing.T) {
	//fill database
	genreIDErr := insertGenreToDatabase(t, "3", "genre1", "трешкор", "anime")
	if genreIDErr != nil {
		log.Fatal(genreIDErr)
	}
	studioIDErr := insertStudioToDatabase(t, "4", "studio", "studio", true, "/url.jpg")
	if studioIDErr != nil {
		log.Fatal(studioIDErr)
	}
	insertAnimeToDatabase(t, "123", "anime", "аниме", "/url.jpg", "tv", "ongoing", 10, 5,
		time.Date(2009, 11, 17, 20, 20, 20, 0, time.UTC),
		time.Date(2009, 11, 17, 20, 20, 20, 0, time.UTC), "/url.jpg", false, "4", "3")
	//create request
	request, _ := http.NewRequest("GET", "/api/animes/search", nil)
	executeRequest(request)
}
