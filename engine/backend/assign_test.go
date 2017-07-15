package backend

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/google/go-github/github"
)

func TestAssignContributor(t *testing.T) {
	issueNumber := 63
	repoOwner := "heupr"
	repoID := 3
	repoName := "test"
	assignee := "forstmeier"
	testURL := fmt.Sprintf("/repos/%v/%v/issues/%v/assignees", repoOwner, repoName, issueNumber)

	mux := http.NewServeMux()
	mux.HandleFunc(testURL, func(w http.ResponseWriter, r *http.Request) {})
	server := httptest.NewServer(mux)
	client := github.NewClient(nil)
	url, _ := url.Parse(server.URL)
	client.BaseURL = url
	client.UploadURL = url

	testBS := new(BackendServer)
	testBS.Repos = new(ActiveRepos)
	testBS.Repos.Actives = make(map[int]*ArchRepo)
	testBS.Repos.Actives[repoID] = new(ArchRepo)
	testBS.Repos.Actives[repoID].Client = client

	user := github.User{Login: &repoOwner}
	repo := github.Repository{ID: &repoID, Owner: &user, Name: &repoName}
	issue := github.Issue{Number: &issueNumber}
	issuesEvent := github.IssuesEvent{Issue: &issue, Repo: &repo}

	if err := testBS.AssignContributor(assignee, issuesEvent); err != nil {
		t.Error(err)
	}
}
