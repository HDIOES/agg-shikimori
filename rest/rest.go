package rest

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/HDIOES/agg-shikimori/rest/util"
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
	headers := make(map[string]string)
	headers["Content-Type"] = "application/json; charset=utf-8"
	if err := LogHTTPResponse(httpStatus, headers, body); err != nil {
		return errors.Wrap(err, "")
	}
	for key, value := range headers {
		w.Header().Set(key, value)
	}
	w.WriteHeader(httpStatus)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		return errors.Wrap(err, "")
	}
	return nil
}

func GetRequestData(r *http.Request) (requestBody []byte, rawQuery *string, headers http.Header, reqErr error) {
	if r.Body == nil {
		return nil, &r.URL.RawQuery, r.Header, nil
	}
	defer r.Body.Close()
	if requestBodyAsBytes, requestBodyErr := ioutil.ReadAll(r.Body); requestBodyErr != nil {
		return nil, nil, nil, errors.WithStack(requestBodyErr)
	} else {
		return requestBodyAsBytes, &r.URL.RawQuery, r.Header, nil
	}
}

func LogHTTPRequest(url, method string, headers http.Header, body interface{}) error {
	const logLineTemplate = "Http request: URL: %v Method: %v Headers: %v Body: %v"
	if bodyAsBytes, ok := body.([]byte); ok {
		log.Printf(logLineTemplate, url, method, headers, bodyAsBytes)
	} else if bodyAsString, ok := body.(string); ok {
		log.Printf(logLineTemplate, url, method, headers, bodyAsString)
	} else {
		bodyAsBytes, err := json.Marshal(body)
		if err != nil {
			return errors.Wrap(err, "")
		}
		log.Printf(logLineTemplate, url, headers, string(bodyAsBytes))
	}
	return nil
}

func LogHTTPResponse(httpStatus int, headers map[string]string, body interface{}) error {
	const logLineTemplate = "Http response: Status: %v Headers: %v Body: %v"
	if bodyAsBytes, ok := body.([]byte); ok {
		log.Printf(logLineTemplate, httpStatus, headers, bodyAsBytes)
	} else if bodyAsString, ok := body.(string); ok {
		log.Printf(logLineTemplate, httpStatus, headers, bodyAsString)
	} else {
		bodyAsBytes, err := json.Marshal(body)
		if err != nil {
			return errors.Wrap(err, "")
		}
		log.Printf(logLineTemplate, httpStatus, headers, string(bodyAsBytes))
	}
	return nil
}
