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
			appId          int
			installationId int
		}
		output interface{}
	}{
		{struct {
			appId          int
			installationId int
		}{1, 2}, nil},
		{struct {
			appId          int
			installationId int
		}{2, 2}, nil},
		{struct {
			appId          int
			installationId int
		}{3, 3}, nil},
	}

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

	NewClient = func(appId int, installationId int) *github.Client {
		c := github.NewClient(nil)
		c.BaseURL = testURL
		c.UploadURL = testURL
		return c
	}

	for _, test := range tests {
		c := NewClient(test.integration.appId, test.integration.installationId)
		if c == nil {
			t.Errorf("failed to make client")
		}
	}
}
