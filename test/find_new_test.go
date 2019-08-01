package test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindNew_success(t *testing.T) {
	clearDb(db)
	insertNewToDatabase(t, 1, "hello", "body")
	request, _ := http.NewRequest("GET", "/api/news?id=1", nil)
	recorder := executeRequest(request)
	//asserts
	assert.Equal(t, 200, recorder.Code)
}
