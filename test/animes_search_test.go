package test

import (
	"testing"

	_ "github.com/lib/pq"
	gock "gopkg.in/h2non/gock.v1"
)

func initTestDataContainer() {
	//start up postgresql test container
}

func TestSearchAnimesSuccess(t *testing.T) {
	defer gock.Off()
	initTestDataContainer()

}
