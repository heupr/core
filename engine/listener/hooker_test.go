package listener

import (
	"github.com/google/go-github/github"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestHeuprHook(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	client := github.NewClient(nil)
	url, _ := url.Parse(server.URL)
	client.BaseURL = url
	client.UploadURL = url

	defer server.Close()

	mux.HandleFunc("/repos/owner/repository/hooks", func(w http.ResponseWriter, r *http.Request) {})

	err := HeuprHook("owner", "repository", client)
	if err != nil {
		t.Error(err)
	}
}
