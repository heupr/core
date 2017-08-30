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
	url, _ := url.Parse(server.URL)

	client := github.NewClient(nil)
	client.BaseURL = url
	client.UploadURL = url

	testGateway := Gateway{
		Client:      client,
		UnitTesting: false,
	}

	t.Run("issues", func(t *testing.T) {
		_, err := testGateway.GetIssues(owner, repo)
		if err != nil {
			t.Errorf("failued retrieving issues, %v", err)
		}
	})
	t.Run("pull request", func(t *testing.T) {
		_, err := testGateway.GetPullRequests(owner, repo)
		if err != nil {
			t.Errorf("failued retrieving issues, %v", err)
		}
	})
}

func TestCachedGateway(t *testing.T) {
	client := github.NewClient(nil)
	gateway := CachedGateway{Gateway: &Gateway{Client: client, UnitTesting: true}, DiskCache: &DiskCache{}}

	pulls, _ := gateway.GetPullRequests("dotnet", "corefx")
	issues, _ := gateway.GetIssues("dotnet", "corefx")

	if pulls == nil {
		t.Error("failed cached gateway pull request fetch")
	}

	if issues == nil {
		t.Error("failed cached gateway issue fetch")
	}
}
