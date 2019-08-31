package rest

import (
	"net/http"
	"net/url"
	"strconv"

	"github.com/HDIOES/cpa-backend/models"
	"github.com/pkg/errors"
)

//GenreHandler struct
type GenreHandler struct {
	Dao *models.GenreDAO
}

func (g *GenreHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	requestBody, rawQuery, headers, err := GetRequestData(r)
	if err != nil {
		HandleErr(errors.Wrap(err, ""), w, 400, "Request cannot be read")
		return
	}
	if err := LogHTTPRequest(*rawQuery, headers, requestBody); err != nil {
		HandleErr(errors.Wrap(err, ""), w, 400, "Request cannot be logged")
		return
	}
	vars, parseErr := url.ParseQuery(*rawQuery)
	if parseErr != nil {
		HandleErr(errors.Wrap(parseErr, ""), w, 400, "Url not valid")
	}
	genreSQLBuilder := models.GenreQueryBuilder{}
	if limit, limitOk := vars["limit"]; limitOk {
		limitInt64, parseErr := strconv.ParseInt(limit[0], 10, 32)
		if parseErr != nil {
			HandleErr(errors.Wrap(parseErr, ""), w, 400, "Not valid limit")
			return
		}
		genreSQLBuilder.SetLimit(int32(limitInt64))
	}
	if offset, offsetOk := vars["offset"]; offsetOk {
		offsetInt64, parseErr := strconv.ParseInt(offset[0], 10, 32)
		if parseErr != nil {
			HandleErr(errors.Wrap(parseErr, ""), w, 400, "Not valid offset")
			return
		}
		genreSQLBuilder.SetOffset(int32(offsetInt64))

	}
	genreDtos, findByFilterErr := g.Dao.FindByFilter(genreSQLBuilder)
	if findByFilterErr != nil {
		HandleErr(errors.Wrap(findByFilterErr, ""), w, 500, "Error")
		return
	}
	genres := []GenreRo{}
	for _, genreDto := range genreDtos {
		genreRo := GenreRo{}
		genreRo.ID = &genreDto.ExternalID
		genreRo.Name = genreDto.Name
		genreRo.Russian = genreDto.Russian
		genreRo.Kind = genreDto.Kind
		genres = append(genres, genreRo)
	}
	if err := ReturnResponseAsJSON(w, genres, 200); err != nil {
		HandleErr(errors.Wrap(err, ""), w, 500, "Error")
	}
}

//GenreRo struct
type GenreRo struct {
	ID      *string `json:"id"`
	Name    *string `json:"name"`
	Russian *string `json:"russian"`
	Kind    *string `json:"kind"`
}
