package pipeline

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_activationServerHandler(t *testing.T) {
	http.HandleFunc("/activate-ingestor-backend", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	go http.ListenAndServe(":8020", nil)
	go http.ListenAndServe(":8030", nil)

	testAS := new(ActivationServer)

	rec := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/handler-test", nil)
	if err != nil {
		t.Errorf("failure generating activation service test request: %v", err)
	}
	req.Body = ioutil.NopCloser(bytes.NewBufferString("Chalmun's Spaceport Cantina"))

	handler := http.HandlerFunc(testAS.activationServerHandler)
	handler.ServeHTTP(rec, req)

	statusOK := http.StatusOK
	if received := rec.Code; received != statusOK {
		t.Errorf("handler returning incorrect status code; received %v, expected %v", received, statusOK)
	}
}
