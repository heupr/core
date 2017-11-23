package ingestor

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/google/go-github/github"
)

func Test_NewClient(t *testing.T) {
	// NOTE: I'd ultimately like this to be truly "table driven" but I can't
	// think of a good way to assert the NewClient results.
	tests := []struct {
		integration struct {
			appID          int
			installationID int
		}
		output interface{}
	}{
		{struct {
			appID          int
			installationID int
		}{1, 2}, nil},
		{struct {
			appID          int
			installationID int
		}{2, 2}, nil},
		{struct {
			appID          int
			installationID int
		}{3, 3}, nil},
	}

	owner := "bomarr-order"
	repo := "bt-16-perimeter-droid"
	// This is a dummy GitHub server to return test objects for issues, pulls,
	// and repos.
	mux := http.NewServeMux()
	mux.HandleFunc(fmt.Sprintf("/repos/%v/%v/issues", owner, repo), func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `[{"id":123,"number":456}]`)
	})
	mux.HandleFunc(fmt.Sprintf("/repos/%v/%v/pulls", owner, repo), func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `[{"id":321,"number":456,"base":{"repo":{"id":890}}}]`)
	})
	mux.HandleFunc(fmt.Sprintf("/repos/%v/%v", owner, repo), func(w http.ResponseWriter, r *http.Request) {
		obj := fmt.Sprintf(`{"id":213,"name":"%v","owner":{"login":"%v"}}`, repo, owner)
		fmt.Fprint(w, obj)
	})
	server := httptest.NewServer(mux)
	testURL, _ := url.Parse(server.URL)

	NewClient = func(appID int, installationID int) *github.Client {
		c := github.NewClient(nil)
		c.BaseURL = testURL
		c.UploadURL = testURL
		return c
	}

	for _, test := range tests {
		c := NewClient(test.integration.appID, test.integration.installationID)
		if c == nil {
			t.Errorf("failed to make client")
		}
	}
}
