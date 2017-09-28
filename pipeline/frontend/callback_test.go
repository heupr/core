package frontend

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"golang.org/x/oauth2"
)

func Test_githubCallbackHandler(t *testing.T) {
	testFS := new(FrontendServer)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
		w.Write([]byte("access_token=039f5f2f98a87f46abef10170866ed8ecf3b5b2d&scope=user&token_type=bearer"))
	}))

	oaConfig = &oauth2.Config{
		ClientID:     "TEST_ID",
		ClientSecret: "TEST_SECRET",
		Endpoint: oauth2.Endpoint{
			AuthURL:  ts.URL + "/auth",
			TokenURL: ts.URL + "/token",
		},
		Scopes: []string{"test-scope"},
	}

	testURL := oaConfig.AuthCodeURL(oaState, oauth2.AccessTypeOffline)

	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", testURL, nil)
	if err != nil {
		t.Errorf("failure generating test request: %v", err)
	}

	handler := http.HandlerFunc(testFS.githubCallbackHandler)
	handler.ServeHTTP(rec, req)

	wantedStatus := http.StatusOK
	if receivedStatus := rec.Code; receivedStatus != wantedStatus {
		t.Errorf("handler returning incorrect status code; received %v, wanted %v", receivedStatus, wantedStatus)
	}
}
