package rest

import (
	"net/http"
	"net/url"
	"strconv"

	"github.com/HDIOES/cpa-backend/models"
	"github.com/pkg/errors"
)

//StudioHandler struct
type StudioHandler struct {
	Dao *models.StudioDAO
}

func (g *StudioHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	requestBody, rawQuery, headers, err := GetRequestData(r)
	if err != nil {
		HandleErr(errors.Wrap(err, ""), w, 400, "Request cannot be read")
		return
	}
	if err := LogHTTPRequest(r.URL.String(), r.Method, headers, requestBody); err != nil {
		HandleErr(errors.Wrap(err, ""), w, 400, "Request cannot be logged")
		return
	}
	vars, parseErr := url.ParseQuery(*rawQuery)
	if parseErr != nil {
		HandleErr(errors.Wrap(parseErr, ""), w, 400, "URL not valid")
	}
	studioSQLBuilder := models.StudioQueryBuilder{}
	if limit, limitOk := vars["limit"]; limitOk {
		limitInt64, parseErr := strconv.ParseInt(limit[0], 10, 32)
		if parseErr != nil {
			HandleErr(errors.Wrap(parseErr, ""), w, 400, "Not valid limit")
			return
		}
		studioSQLBuilder.SetOffset(int32(limitInt64))
	}
	if offset, offsetOk := vars["offset"]; offsetOk {
		offsetInt64, parseErr := strconv.ParseInt(offset[0], 10, 32)
		if parseErr != nil {
			HandleErr(errors.Wrap(parseErr, ""), w, 400, "Not valid offset")
			return
		}
		studioSQLBuilder.SetOffset(int32(offsetInt64))
	}
	if studiosDtos, stErr := g.Dao.FindByFilter(studioSQLBuilder); stErr != nil {
		HandleErr(stErr, w, 400, "Error")
	} else {
		studios := []StudioRo{}
		for _, dto := range studiosDtos {
			ro := StudioRo{
				ID:           &dto.ID,
				Name:         dto.Name,
				FilteredName: dto.FilteredStudioName,
			}
			studios = append(studios, ro)
		}
		if err := ReturnResponseAsJSON(w, studios, 200); err != nil {
			HandleErr(err, w, 500, "Error")
		}
	}
}

//StudioRo struct
type StudioRo struct {
	ID           *int64  `json:"id"`
	Name         *string `json:"name"`
	FilteredName *string `json:"filtered_name"`
}
