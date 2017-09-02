package ingestor

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"

	"core/pipeline/frontend"
)

func Test_NewClient(t *testing.T) {
	// NOTE: I'd ultimately like this to be truly "table driven" but I can't
	// think of a good way to assert the NewClient results.
	tests := []struct {
		token  oauth2.Token
		output interface{}
	}{
		{oauth2.Token{AccessToken: ""}, nil},
		{oauth2.Token{AccessToken: "test-token"}, nil},
		{oauth2.Token{AccessToken: "039f5f2f98a87f46abef10170866ed8ecf3b5b2d"}, nil},
	}

	for _, test := range tests {
		c := NewClient(test.token)
		if c == nil {
			t.Errorf("failed to make client")
		}
	}
}

func Test_activateHandler(t *testing.T) {
	owner := "bomarr-order"
	repo := "bt-16-perimeter-droid"

	mux := http.NewServeMux()
	mux.HandleFunc(fmt.Sprintf("/repos/%v/%v/issues", owner, repo), func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `[{"id":123,"number":456}]`) // Issues dummy service.
	})
	mux.HandleFunc(fmt.Sprintf("/repos/%v/%v/pulls", owner, repo), func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `[{"id":321,"number":456,"base":{"repo":{"id":890}}}]`) // Pulls dummy service.
	})
	mux.HandleFunc(fmt.Sprintf("/repos/%v/%v", owner, repo), func(w http.ResponseWriter, r *http.Request) {
		obj := fmt.Sprintf(`{"id":213,"name":"%v","owner":{"login":"%v"}}`, repo, owner)
		fmt.Fprint(w, obj) // Repos dummy service.
	})
	server := httptest.NewServer(mux)
	testURL, _ := url.Parse(server.URL)

	NewClient = func(t oauth2.Token) *github.Client {
		c := github.NewClient(nil)
		c.BaseURL = testURL
		c.UploadURL = testURL
		return c
	}

	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/handler-test", nil)
	if err != nil {
		t.Errorf("failure generating testing request: %v", err)
	}
	req.Form = url.Values{}
	req.Form.Set("state", frontend.BackendSecret)
	req.Form.Set("repos", fmt.Sprintf("{testID, %v, %v}", owner, repo))
	req.Form.Set("token", "fake-test-token")
	req.ParseForm()

	testIS := IngestorServer{}
	handler := http.HandlerFunc(testIS.activateHandler)
	handler.ServeHTTP(rec, req)
	if received := rec.Code; received != http.StatusOK {
		t.Errorf("handler returning incorrect status code; received %v, expected %v", received, http.StatusOK)
	}
}
