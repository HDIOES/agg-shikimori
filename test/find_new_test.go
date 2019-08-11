package test

import (
	"net/http"
	"strconv"
	"testing"

	"github.com/gorilla/mux"

	"github.com/HDIOES/cpa-backend/models"
	"github.com/stretchr/testify/assert"
)

func TestFindNew_success(t *testing.T) {
	diContainer.Invoke(func(router *mux.Router, newDao *models.NewDAO, animeDao *models.AnimeDAO, genreDao *models.GenreDAO, studioDao *models.StudioDAO) {
		clearDb(newDao, animeDao, genreDao, studioDao)
		newName := "hello"
		body := "body"
		id, err := insertNewToDatabase(newDao, &newName, &body)
		if err != nil {
			t.Fatal(err)
		}
		request, _ := http.NewRequest("GET", "/api/news?id="+strconv.FormatInt(id, 10), nil)
		recorder := executeRequest(request, router)
		//asserts
		assert.Equal(t, 200, recorder.Code)
	})
}
