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

    // sig := "sha1=" + hex.EncodeToString([]byte("825526bac599fa432336d57653a562142ffa94af"))

    // sig := "sha1=825526bac599fa432336d57653a562142ffa94af"

    mac := hmac.New(sha1.New, []byte(secretKey))
	mac.Write([]byte(testBody))
	// s := mac.Sum(nil)

    // h := sha1.New()
    // io.WriteString(h, secretKey)
    // io.WriteString(h)
    sig := "sha1=" + hex.EncodeToString(mac.Sum(nil))

    // TEMPORARY
    // mac := hmac.New(sha1.New, key)
    // mac.Write(message)
    // mac.Sum(nil)
    // sig := "sha1=" + hex.EncodeToString([]byte(secretKey))

    req.Header.Set("X-Hub-Signature", sig)
    // fmt.Println(req)                                 // TEMPORARY
    // fmt.Println(req.Header.Get("X-Hub-Signature"))   // TEMPORARY
    handler := collectorHandler("")
	handler.ServeHTTP(rec, req)
	if status := rec.Code; status != http.StatusOK {
		t.Errorf("Handler returning incorrect status code; returning %v", status)
	}
}

/*
// e.g. GET /api/projects?page=1&per_page=100
  req, err := http.NewRequest("GET", "/api/projects",
      // Note: url.Values is a map[string][]string
      url.Values{"page": {"1"}, "per_page": {"100"}})
  if err != nil {
      t.Fatal(err)
  }

  // Our handler might also expect an API key.
  req.Header.Set("Authorization", "Bearer abc123")
*/
