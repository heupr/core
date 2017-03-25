package frontend

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_collectorHandler(t *testing.T) {
	testBody := `{"id":1}`
	rec := httptest.NewRecorder()
	buf := bytes.NewBufferString(testBody)
	req, err := http.NewRequest("POST", "/handler-test", buf)
	if err != nil {
		t.Error(
			"Failure generating testing request",
			"\n", err,
		)
	}

	req.Header.Set("X-Github-Event", "issues")
	mac := hmac.New(sha1.New, []byte(secretKey))
	mac.Write([]byte(testBody))
	sig := "sha1=" + hex.EncodeToString(mac.Sum(nil))
	req.Header.Set("X-Hub-Signature", sig)

	handler := collectorHandler()
	handler.ServeHTTP(rec, req)
	if status := rec.Code; status != http.StatusOK {
		t.Errorf("Handler returning incorrect status code; returning %v", status)
	}
}
