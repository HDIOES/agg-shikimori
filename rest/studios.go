package rest

import (
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/HDIOES/cpa-backend/models"
)

//StudioHandler struct
type StudioHandler struct {
	Dao *models.StudioDAO
}

func (g *StudioHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars, parseErr := url.ParseQuery(r.URL.RawQuery)
	if parseErr != nil {
		log.Println(parseErr)
	}
	studioSQLBuilder := models.StudioQueryBuilder{}
	if limit, limitOk := vars["limit"]; limitOk {
		limitInt64, parseErr := strconv.ParseInt(limit[0], 10, 32)
		if parseErr != nil {
			HandleErr(parseErr, w, 400, "Not valid limit")
			return
		}
		studioSQLBuilder.SetOffset(int32(limitInt64))
	}
	if offset, offsetOk := vars["offset"]; offsetOk {
		offsetInt64, parseErr := strconv.ParseInt(offset[0], 10, 32)
		if parseErr != nil {
			HandleErr(parseErr, w, 400, "Not valid offset")
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
		ReturnResponseAsJSON(w, studios, 200)
	}
}

//StudioRo struct
type StudioRo struct {
	ID           *int64  `json:"id"`
	Name         *string `json:"name"`
	FilteredName *string `json:"filtered_name"`
}
