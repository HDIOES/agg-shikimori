package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"

	"github.com/HDIOES/agg-shikimori/models"
	"github.com/HDIOES/agg-shikimori/rest"
	"github.com/stretchr/testify/assert"

	_ "github.com/lib/pq"
)

func TestCreateNew_success(t *testing.T) {
	diContainer.Invoke(func(router *mux.Router, newDao *models.NewDAO, animeDao *models.AnimeDAO, genreDao *models.GenreDAO, studioDao *models.StudioDAO) {
		if err := clearDb(newDao, animeDao, genreDao, studioDao); err != nil {
			markAsFailAndAbortNow(t, errors.Wrap(err, ""))
		}
		Name := "Ужасная статья"
		Body := "Ужасная статья?"
		requestBody := rest.NewRo{Name: &Name, Body: &Body}
		reader, readErr := GetJSONBodyReader(requestBody)
		if readErr != nil {
			markAsFailAndAbortNow(t, errors.Wrap(readErr, ""))
		}
		//create request
		request, err := http.NewRequest("POST", "/api/news", reader)
		if err != nil {
			markAsFailAndAbortNow(t, errors.Wrap(err, ""))
		}
		recorder := executeRequest(request, router)
		//asserts
		assert.Equal(t, 200, recorder.Code)
	})
}

//GetJsonBodyReader function
func GetJSONBodyReader(body interface{}) (*bytes.Reader, error) {
	data, dataErr := json.Marshal(body)
	if dataErr != nil {
		return nil, errors.Wrap(dataErr, "")
	}
	reader := bytes.NewReader(data)
	return reader, nil
}
