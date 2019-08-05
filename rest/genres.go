package rest

import (
	"database/sql"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/HDIOES/cpa-backend/models"
)

func CreateGenreHandler(db *sql.DB) http.Handler {
	dao := models.GenreDAO{Db: db}
	genreHandler := &GenreHandler{GenreDao: &dao}
	return genreHandler
}

type GenreHandler struct {
	GenreDao *models.GenreDAO
}

func (g *GenreHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars, parseErr := url.ParseQuery(r.URL.RawQuery)
	if parseErr != nil {
		log.Println(parseErr)
	}
	genreSQLBuilder := models.GenreQueryBuilder{}
	if limit, limitOk := vars["limit"]; limitOk {
		if limitInt64, parseErr := strconv.ParseInt(limit[0], 10, 32); parseErr != nil {
			HandleErr(parseErr, w, 400, "Not valid limit")
			return
		} else {
			genreSQLBuilder.SetOffset(int32(limitInt64))
		}
	}
	if offset, offsetOk := vars["offset"]; offsetOk {
		if offsetInt64, parseErr := strconv.ParseInt(offset[0], 10, 32); parseErr != nil {
			HandleErr(parseErr, w, 400, "Not valid offset")
			return
		} else {
			genreSQLBuilder.SetOffset(int32(offsetInt64))
		}
	}
	genreDtos, findByFilterErr := g.GenreDao.FindByFilter(genreSQLBuilder)
	if findByFilterErr != nil {
		HandleErr(findByFilterErr, w, 400, "Error")
		return
	}
	genres := []GenreRo{}
	for _, genreDto := range genreDtos {
		genreRo := GenreRo{}
		genreRo.ID = genreDto.ExternalID
		genreRo.Name = genreDto.Name
		genreRo.Russian = genreDto.Russian
		genreRo.Kind = genreDto.Kind
		genres = append(genres, genreRo)
	}
	ReturnResponseAsJSON(w, genres, 200)
}

//GenreRo struct
type GenreRo struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Russian string `json:"russian"`
	Kind    string `json:"kind"`
}
