package rest

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/HDIOES/cpa-backend/models"
)

//CreateCreateNewHandler func
func CreateCreateNewHandler(db *sql.DB) http.Handler {
	dao := models.NewDAO{Db: db}
	createNewHandler := &CreateNewHandler{dao: &dao}
	return createNewHandler
}

//CreateNewHandler struct
type CreateNewHandler struct {
	dao *models.NewDAO
}

func (cnh *CreateNewHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	body, ioErr := ioutil.ReadAll(r.Body)
	if ioErr != nil {
		HandleErr(ioErr, w, 400, "Bad request")
	}
	new := &NewRo{}
	if unmErr := json.Unmarshal(body, new); unmErr != nil {
		HandleErr(unmErr, w, 400, "Bad request")
	}
	newDto := models.NewDTO{ID: new.ID, Name: new.Name, Body: new.Body}
	if newID, createErr := cnh.dao.Create(newDto); createErr != nil {
		HandleErr(createErr, w, 400, "Error")
	} else {
		ReturnResponseAsJSON(w, CreateNewResponse{ID: newID}, 200)
	}
}

//FindNewHandler struct
type FindNewHandler struct {
	dao *models.NewDAO
}

//CreateFindNewHandler function
func CreateFindNewHandler(db *sql.DB) http.Handler {
	dao := models.NewDAO{Db: db}
	findNewHandler := &FindNewHandler{dao: &dao}
	return findNewHandler
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
			if newDTO, findErr := fnh.dao.Find(idInt64); findErr != nil {
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
	ID int64 `json:"id"`
}

//NewRo struct
type NewRo struct {
	ID   int64
	Name string
	Body string
}
