package rest

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/HDIOES/cpa-backend/models"
	"github.com/pkg/errors"
)

//CreateNewHandler struct
type CreateNewHandler struct {
	Dao *models.NewDAO
}

func (cnh *CreateNewHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	body, ioErr := ioutil.ReadAll(r.Body)
	if ioErr != nil {
		HandleErr(errors.Wrap(ioErr, ""), w, 400, "Bad request")
	}
	new := &NewRo{}
	if unmErr := json.Unmarshal(body, new); unmErr != nil {
		HandleErr(errors.Wrap(unmErr, ""), w, 400, "Bad request")
	}
	newDto := models.NewDTO{ID: new.ID, Name: new.Name, Body: new.Body}
	if newID, createErr := cnh.Dao.Create(newDto); createErr != nil {
		HandleErr(errors.Wrap(createErr, ""), w, 400, "Error")
	} else {
		if err := ReturnResponseAsJSON(w, CreateNewResponse{ID: &newID}, 200); err != nil {
			HandleErr(errors.Wrap(createErr, ""), w, 500, "Error")
		}
	}
}

//FindNewHandler struct
type FindNewHandler struct {
	Dao *models.NewDAO
}

func (fnh *FindNewHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars, parseErr := url.ParseQuery(r.URL.RawQuery)
	if parseErr != nil {
		HandleErr(parseErr, w, 400, "Error")
	}
	if id, idOk := vars["id"]; idOk {
		if idInt64, parseErr := strconv.ParseInt(id[0], 10, 64); parseErr != nil {
			HandleErr(parseErr, w, 400, "id not valid")
		} else {
			if newDTO, findErr := fnh.Dao.Find(idInt64); findErr != nil {
				HandleErr(findErr, w, 400, "Error")
			} else {
				ro := NewRo{ID: newDTO.ID, Name: newDTO.Name, Body: newDTO.Body}
				ReturnResponseAsJSON(w, &ro, 200)
			}
		}
	}
}

//CreateNewResponse struct
type CreateNewResponse struct {
	ID *int64 `json:"id"`
}

//NewRo struct
type NewRo struct {
	ID   *int64
	Name *string
	Body *string
}
