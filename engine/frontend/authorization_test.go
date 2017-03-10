package frontend

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_mainHandle(t *testing.T) {
	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/handler-test", nil)
	if err != nil {
		t.Error(
			"Failure generating testing request",
			"\n", err,
		)
	}
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(rec, req)
	if status := rec.Code; status != http.StatusOK {
		t.Errorf("Handler returning incorrect status code; returning %v", status)
	}
}

func Test_githubLoginHandle(t *testing.T) {
	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/handler-test", nil)
	if err != nil {
		t.Error(
			"Failure generating testing request",
			"\n", err,
		)
	}
	handler := http.HandlerFunc(githubLoginHandle)
	handler.ServeHTTP(rec, req)
	if status := rec.Code; status != http.StatusTemporaryRedirect {
		t.Errorf("Handler returning incorrect status code; returning %v", status)
	}
}
