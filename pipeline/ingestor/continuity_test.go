package ingestor

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/google/go-github/github"
)

var tests = []struct {
	query   []interface{}
	missing []int
}{
	{[]interface{}{1, 1, 3, false}, []int{2}},
	{[]interface{}{1, 3, 6, false}, []int{4, 5}},
	{[]interface{}{1, 7, 9, true}, []int{8}},
}

type continuityDA struct{}

func (c *continuityDA) open() {}

func (c *continuityDA) Close() {}

func (c *continuityDA) continuityCheck(query string) ([][]interface{}, error) {
	testResults := [][]interface{}{}
	for i := range tests {
		testResults = append(testResults, tests[i].query)
	}
	return testResults, nil
}

func (c *continuityDA) restartCheck(query string, repoID int64) (int, int, error) {
	return 0, 0, nil
}

func (r *continuityDA) InsertGobLabelSettings(settings storage) error {
	return nil
}

func (c *continuityDA) ReadIntegrations() ([]Integration, error) { return nil, nil }

func (c *continuityDA) ReadIntegrationByRepoID(id int64) (*Integration, error) {
	return &Integration{1, 1, 1}, nil
}

func (c *continuityDA) InsertIssue(i github.Issue, action *string) {}

func (c *continuityDA) InsertPullRequest(p github.PullRequest, action *string) {}

func (c *continuityDA) BulkInsertIssuesPullRequests(i []*github.Issue, p []*github.PullRequest) {}

func (c *continuityDA) InsertRepositoryIntegration(repoID int64, appID int, installID int64) {}

func (c *continuityDA) InsertRepositoryIntegrationSettings(settings HeuprConfigSettings) {}

func (c *continuityDA) DeleteRepositoryIntegration(repoID int64, appID int, installID int64) {}

func (c *continuityDA) ObliterateIntegration(appID int, installID int64) {}

func Test_continuityCheck(t *testing.T) {
	// This is the fake GitHub server that is queried by the method. Below are
	// the handlers to return a repo, issues, and a pull, respectively.
	mux := http.NewServeMux()
	mux.HandleFunc("/repositories/1", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id":1,"name":"trade-federation","owner":{"login":"nute-gunray"}}`)
	})
	mux.HandleFunc("/repos/nute-gunray/trade-federation/issues/2", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id":2,"number":2}`)
	})
	mux.HandleFunc("/repos/nute-gunray/trade-federation/issues/4", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id":4,"number":4}`)
	})
	mux.HandleFunc("/repos/nute-gunray/trade-federation/issues/5", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id":5,"number":5}`)
	})
	mux.HandleFunc("/repos/nute-gunray/trade-federation/pulls/8", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id":8,"number":8}`)
	})
	server := httptest.NewServer(mux)
	testURL, _ := url.Parse(server.URL + "/")

	NewClient = func(appID int, installationID int) *github.Client {
		c := github.NewClient(nil)
		c.BaseURL = testURL
		c.UploadURL = testURL
		return c
	}

	testIS := IngestorServer{
		Database: &continuityDA{},
	}

	issues, pulls, err := testIS.continuityCheck()
	if err != nil {
		t.Errorf("continuity check error: %v", err)
	}

	// This is just a simple check to make sure that continuityCheck is
	// returning the same number of issues/pulls that are fed in as tests.
	issueCount := []int{}
	pullCount := []int{}
	for i := range tests {
		if tests[i].query[3] == false {
			issueCount = append(issueCount, tests[i].missing...)
		} else {
			pullCount = append(pullCount, tests[i].missing...)
		}
	}
	if len(issues) != len(issueCount) {
		t.Error("returned issues not equal to test issues")
	}
	if len(pulls) != len(pullCount) {
		t.Error("returned pull not equal to test pulls")
	}
}
