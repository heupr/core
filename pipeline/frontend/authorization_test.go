package frontend

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

/*
func Test_mainHandler(t *testing.T) {
	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/handler-test", nil)
	if err != nil {
		t.Errorf("Failure generating testing request: %v", err)
	}
	handler := mainHandler
	handler.ServeHTTP(rec, req)
	if status := rec.Code; status != http.StatusOK {
		t.Errorf("Handler returning incorrect status code; returning %v", status)
	}
}
*/

func Test_githubLoginHandler(t *testing.T) {
	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/handler-test", nil)
	if err != nil {
		t.Errorf("Failure generating testing request: %v", err)
	}
	handler := http.HandlerFunc(githubLoginHandler)
	handler.ServeHTTP(rec, req)
	if status := rec.Code; status != http.StatusTemporaryRedirect {
		t.Errorf("Handler returning incorrect status code; returning %v", status)
	}
}
