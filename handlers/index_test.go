package handlers

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIndexHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "http://example.com/", nil)
	if err != nil {
		log.Fatal(err)
	}

	w := httptest.NewRecorder()
	IndexHandler(w, req)

	fmt.Printf("%d - %s", w.Code, w.Body.String())
}
