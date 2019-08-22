package rest

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/HDIOES/cpa-backend/rest/util"
	"github.com/pkg/errors"
)

//HandleErr function
func HandleErr(err error, w http.ResponseWriter, httpStatus int, errorMessage string) error {
	util.HandleError(err)
	errorMessageBuilder := strings.Builder{}
	errorMessageBuilder.WriteString("{")
	errorMessageBuilder.WriteString("\"message\":")
	errorMessageBuilder.WriteString("\"")
	errorMessageBuilder.WriteString(errorMessage)
	errorMessageBuilder.WriteString("\"")
	errorMessageBuilder.WriteString("}")
	return ReturnResponseAsJSON(w, errorMessageBuilder.String(), httpStatus)
}

//ReturnResponseAsJSON function
func ReturnResponseAsJSON(w http.ResponseWriter, body interface{}, httpStatus int) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(httpStatus)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		return errors.Wrap(err, "")
	}
	return nil
}
