package ingestor

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/google/go-github/github"
)

type restartDA struct{}

func (r *restartDA) open() {}

func (r *restartDA) Close() {}

func (r *restartDA) continuityCheck(query string) ([][]interface{}, error) {
	return nil, nil
}

var testRestartCase = 0

func (r *restartDA) restartCheck(query string, repoID int64) (i, p int, err error) {
	switch testRestartCase {
	case 1: // Repo exists in the MemSQL database.
		i, p, err = 1, 1, nil
	case 2: // AddRepo method is called to initialize.
		i, p, err = 0, 0, nil
	}
	return
}

func (r *restartDA) ReadIntegrations() ([]Integration, error) {
	return []Integration{Integration{1, 1, 1}}, nil
}

func (r *restartDA) ReadIntegrationByRepoID(id int64) (*Integration, error) {
	return nil, nil
}

func (r *restartDA) InsertGobLabelSettings(settings storage) error {
	return nil
}

func (r *restartDA) InsertIssue(i github.Issue, action *string) {}

func (r *restartDA) InsertPullRequest(p github.PullRequest, action *string) {}

func (r *restartDA) BulkInsertIssuesPullRequests(i []*github.Issue, p []*github.PullRequest) {}

func (r *restartDA) InsertRepositoryIntegration(repoID int64, appID int, installID int64) {}

func (r *restartDA) InsertRepositoryIntegrationSettings(settings HeuprConfigSettings) {}

func (r *restartDA) DeleteRepositoryIntegration(repoID int64, appID int, installID int64) {}

func (r *restartDA) ObliterateIntegration(appID int, installID int64) {}

func TestRestart(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/repositories/1", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id":1,"name":"stalgasin-hive","owner":{"login":"poggle-the-lesser"}}`)
	})
	server := httptest.NewServer(mux)
	testURL, _ := url.Parse(server.URL + "/")

	NewClient = func(appID int, installationID int) *github.Client {
		c := github.NewClient(nil)
		c.BaseURL = testURL
		c.UploadURL = testURL
		return c
	}

	issueErr := false
	temp := issueGaps
	issueGaps = func(client *github.Client, owner, name string, dbIssueNum int) (i []*github.Issue, err error) {
		switch issueErr {
		case false:
			i, err = []*github.Issue{}, nil
		case true:
			i, err = nil, errors.New("testing error return")
		}
		return
	}

	pullGaps = func(client *github.Client, owner, name string, dbIssueNum int) ([]*github.PullRequest, error) {
		return []*github.PullRequest{}, nil
	}

	testIS := IngestorServer{
		Database: &restartDA{},
		RepoInitializer: RepoInitializer{
			Database: &restartDA{},
		},
	}

	t.Run("existing repo", func(t *testing.T) {
		testRestartCase = 1
		if err := testIS.Restart(); err != nil {
			t.Errorf("failure testing restart: %v", err)
		}
	})
	t.Run("initialize repo", func(t *testing.T) {
		testRestartCase = 2
		if err := testIS.Restart(); err != nil {
			t.Errorf("failure testing restart: %v", err)
		}
	})
	t.Run("gap error return", func(t *testing.T) {
		testRestartCase = 1
		issueErr = true
		if err := testIS.Restart(); err == nil {
			t.Errorf("failure testing restart: %v", err)
		}
	})
	issueGaps = temp
}

func Test_gapHelpers(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/repos/wat-tambor/techno-union/issues", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		switch r.Form["per_page"][0] {
		case "1":
			fmt.Fprint(w, `[{"id":3,"number":3}]`)
		default:
			fmt.Fprint(w, `[{"id":2,"number":2},{"id":3,"number":3}]`)
		}
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

	issues, err := issueGaps(client, "wat-tambor", "techno-union", 1)
	if err != nil {
		t.Errorf("error finding issue gaps: %v", err)
	}
	if len(issues) != 2 {
		t.Errorf("expected two (2) issues, received %v", len(issues))
	}
}
