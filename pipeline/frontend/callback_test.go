package frontend

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_httpRedirect(t *testing.T) {
	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/test-url", nil)
	if err != nil {
		t.Errorf("failure generating test request: %v", err)
	}

	handler := http.HandlerFunc(httpRedirect)
	handler.ServeHTTP(rec, req)

	// This redirects from the non-PROD else check in the production code.
	wantedStatus := http.StatusMovedPermanently
	if receivedStatus := rec.Code; receivedStatus != wantedStatus {
		t.Errorf(
			"handler returning incorrect status code; received %v, wanted %v",
			receivedStatus, wantedStatus,
		)
	}

}

func Test_setupCompleteHandler(t *testing.T) {
	testFS := new(FrontendServer)

	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/test-url", nil)
	if err != nil {
		t.Errorf("failure generating test request: %v", err)
	}

	handler := http.HandlerFunc(testFS.setupCompleteHandler)
	handler.ServeHTTP(rec, req)

	wantedStatus := http.StatusOK
	if receivedStatus := rec.Code; receivedStatus != wantedStatus {
		t.Errorf(
			"handler returning incorrect status code; received %v, wanted %v",
			receivedStatus, wantedStatus,
		)
	}

	// Check that the response body is correct.
	if rec.Body.String() != setup {
		t.Errorf(
			"incorrect response body\n%v",
			rec.Body.String(),
		)
	}
}
