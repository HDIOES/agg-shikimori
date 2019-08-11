package rest

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func HandleErr(err error, w http.ResponseWriter, httpStatus int, errorMessage string) {
	log.Println(err)
	errorMessageBuilder := strings.Builder{}
	errorMessageBuilder.WriteString("{")
	errorMessageBuilder.WriteString("\"message\":")
	errorMessageBuilder.WriteString("\"")
	errorMessageBuilder.WriteString(errorMessage)
	errorMessageBuilder.WriteString("\"")
	errorMessageBuilder.WriteString("}")
	ReturnResponseAsJSON(w, errorMessageBuilder.String(), httpStatus)
}

func ReturnResponseAsJSON(w http.ResponseWriter, body interface{}, httpStatus int) {
	w.WriteHeader(httpStatus)
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(body)
}
