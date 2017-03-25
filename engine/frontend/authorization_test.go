package frontend

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_mainHandler(t *testing.T) {
	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/handler-test", nil)
	if err != nil {
		t.Error(
			"Failure generating testing request",
			"\n", err,
		)
	}
	handler := http.HandlerFunc(mainHandler)
	handler.ServeHTTP(rec, req)
	if status := rec.Code; status != http.StatusOK {
		t.Errorf("Handler returning incorrect status code; returning %v", status)
	}
}

func Test_githubLoginHandler(t *testing.T) {
	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/handler-test", nil)
	if err != nil {
		t.Error(
			"Failure generating testing request",
			"\n", err,
		)
	}
	handler := http.HandlerFunc(githubLoginHandler)
	handler.ServeHTTP(rec, req)
	if status := rec.Code; status != http.StatusTemporaryRedirect {
		t.Errorf("Handler returning incorrect status code; returning %v", status)
	}
}
