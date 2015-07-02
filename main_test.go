package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHello(t *testing.T) {
	req, _ := http.NewRequest("GET", "", nil)
	w := httptest.NewRecorder()
	handler := http.HandlerFunc(hello)
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Home page didn't return %v", http.StatusOK)
	}

	if w.Body.String() != "hello, world\n" {
		t.Errorf("Home page didn't return \"hello, world\", got: \"%s\"", w.Body.String())
	}

}
