package gateway

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/google/go-github/github"
)

func TestGateway(t *testing.T) {
	owner := "darth-krayt"
	repo := "one-sith"

	mux := http.NewServeMux()
	mux.HandleFunc(fmt.Sprintf("/repos/%v/%v/issues", owner, repo), func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `[{"id":123}]`)
	})
	mux.HandleFunc(fmt.Sprintf("/repos/%v/%v/pulls", owner, repo), func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `[{"id":321}]`)
	})
	server := httptest.NewServer(mux)
	url, _ := url.Parse(server.URL + "/")

	client := github.NewClient(nil)
	client.BaseURL = url
	client.UploadURL = url

	testGateway := Gateway{
		Client:      client,
		UnitTesting: true,
	}

	t.Run("issues", func(t *testing.T) {
		_, err := testGateway.getIssues(owner, repo, "closed")
		if err != nil {
			t.Errorf("failed retrieving issues: %v", err)
		}
	})
	t.Run("pulls", func(t *testing.T) {
		_, err := testGateway.getPulls(owner, repo, "closed")
		if err != nil {
			t.Errorf("failed retrieving pulls: %v", err)
		}
	})
}
