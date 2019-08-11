package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gorilla/mux"

	"github.com/HDIOES/cpa-backend/models"
	"github.com/HDIOES/cpa-backend/rest"
	"github.com/stretchr/testify/assert"

	_ "github.com/lib/pq"
)

func TestCreateNew_success(t *testing.T) {
	diContainer.Invoke(func(router *mux.Router, newDao *models.NewDAO, animeDao *models.AnimeDAO, genreDao *models.GenreDAO, studioDao *models.StudioDAO) {
		clearDb(newDao, animeDao, genreDao, studioDao)
		Name := "Ужасная статья"
		Body := "Ужасная статья?"
		requestBody := rest.NewRo{Name: &Name, Body: &Body}
		reader, readErr := GetJSONBodyReader(requestBody)
		if readErr != nil {
			t.Fatal(readErr)
		}
		//create request
		request, _ := http.NewRequest("POST", "/api/news", reader)
		recorder := executeRequest(request, router)
		//asserts
		assert.Equal(t, 200, recorder.Code)
	})
}

//GetJsonBodyReader function
func GetJSONBodyReader(body interface{}) (*bytes.Reader, error) {
	data, dataErr := json.Marshal(body)
	if dataErr != nil {
		return nil, dataErr
	}
	reader := bytes.NewReader(data)
	return reader, nil
}
