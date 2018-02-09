package backend

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func Test_activateHandler(t *testing.T) {
	id := github.Int64(1)
	owner := "bomarr-order"
	repo := "bt-16-perimeter-droid"
	activationParams := struct {
		Repo  github.Repository `json:"repo"`
		Token *oauth2.Token     `json:"token"`
	}{
		github.Repository{
			ID: id,
			Owner: &github.User{
				Login: github.String(owner),
			},
			Name:     github.String(repo),
			FullName: github.String(owner + "/" + repo),
		},
		&oauth2.Token{},
	}

	payload, err := json.Marshal(activationParams)
	if err != nil {
		t.Errorf("failure converting activation parameters: %v", err)
	}
	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/handler-test", bytes.NewBuffer(payload))
	if err != nil {
		t.Errorf("failure generating testing request: %v", err)
	}
	req.Header.Set("content-type", "application/json")

	backendServer := BackendServer{}
	backendServer.Repos = &ActiveRepos{Actives: make(map[int64]*ArchRepo)}
	backendServer.Repos.Actives[*id] = &ArchRepo{}
	handler := http.HandlerFunc(backendServer.activateHandler)
	handler.ServeHTTP(rec, req)
	if received := rec.Code; received != http.StatusOK {
		t.Errorf("handler returning incorrect status code; received %v, expected %v", received, http.StatusOK)
	}
}
