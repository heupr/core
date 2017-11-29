package ingestor

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/google/go-github/github"
)

type repoInitializerDBStub struct {
	issues []*github.Issue
	pulls  []*github.PullRequest
}

func (r *repoInitializerDBStub) open() {}

func (r *repoInitializerDBStub) Close() {}

func (r *repoInitializerDBStub) continuityCheck(query string) ([][]interface{}, error) {
	return nil, nil
}

func (r *repoInitializerDBStub) restartCheck(query string, repoID int) (int, int, error) {
	return 0, 0, nil
}

func (r *repoInitializerDBStub) ReadIntegrations() ([]Integration, error) { return nil, nil }

func (r *repoInitializerDBStub) ReadIntegrationByRepoID(id int) (*Integration, error) { return nil, nil }

func (r *repoInitializerDBStub) InsertIssue(i github.Issue, action *string) {}

func (r *repoInitializerDBStub) InsertPullRequest(p github.PullRequest, action *string) {}

func (r *repoInitializerDBStub) BulkInsertIssuesPullRequests(i []*github.Issue, p []*github.PullRequest) {
	r.issues = i
	r.pulls = p
}

func (r *repoInitializerDBStub) InsertRepositoryIntegration(repoID, appID, installID int) {}

func (r *repoInitializerDBStub) InsertRepositoryIntegrationSettings(settings HeuprConfigSettings) {}

func (r *repoInitializerDBStub) DeleteRepositoryIntegration(repoID, appID, installID int) {}

func (r *repoInitializerDBStub) ObliterateIntegration(appID, installID int) {}

func TestAddRepo(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/repos/san-hill/banking-clan/issues", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `[{"id":1,"number":1},{"id":2,"number":2}]`)
	})
	mux.HandleFunc("repos/san-hill/banking-clan/pulls", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `[{"id":3,"number":3},{"id":4,"number":4}]`)
	})
	server := httptest.NewServer(mux)
	testURL, _ := url.Parse(server.URL + "/")

	NewClient = func(appID int, installationID int) *github.Client {
		c := github.NewClient(nil)
		c.BaseURL = testURL
		c.UploadURL = testURL
		return c
	}
	client := NewClient(1, 1)

	db := &repoInitializerDBStub{}
	testRI := RepoInitializer{
		Database: db,
	}

	testAuthRepo := AuthenticatedRepo{
		Repo: &github.Repository{
			Owner: &github.User{
				Login: github.String("san-hill"),
			},
			Name: github.String("banking-clan"),
		},
		Client: client,
	}
	testRI.AddRepo(testAuthRepo)
	if len(db.issues) != 2 && len(db.pulls) != 2 {
		t.Error("inserting incorrect number of issues/pulls")
	}
}
