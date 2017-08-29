package ingestor

import (
	// "net/http"
	// "net/http/httptest"
	// "net/url"
	"testing"

	// "github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func Test_makeClient(t *testing.T) {
	// NOTE: I'd ultimately like this to be truly "table driven" but I can't
	// think of a good way to assert the makeClient results.
	tests := []struct {
		token  oauth2.Token
		output interface{}
	}{
		{oauth2.Token{AccessToken: ""}, nil},
		{oauth2.Token{AccessToken: "test-token"}, nil},
		{oauth2.Token{AccessToken: "039f5f2f98a87f46abef10170866ed8ecf3b5b2d"}, nil},
	}

	for _, test := range tests {
		c := makeClient(test.token)
		if c == nil {
			t.Errorf("failed to make client")
		}
	}
}

/*
func Test_activateHandler(t *testing.T) {
	testIS := IngestorServer{}
	fetchGitHub = func(owner, name string, client github.Client) ([]*github.Issue, []*github.PullRequest, *github.Repository, error) {
		i := []*github.Issue{}
		p := []*github.PullRequest{}
		r := &github.Repository{}
		return i, p, r, nil
	}

	rec := httptest.NewRecorder()

	req, err := http.NewRequest("GET", "/handler-test", nil)
	if err != nil {
		t.Errorf("failure generating testing request: %v", err)
	}
	req.Form = url.Values{}
	req.Form.Set("state", frontend.BackendSecret)
	req.Form.Set("repos", "testID,testOwner,testRepo")
	req.Form.Set("token", "fake-test-token")

	handler := http.HandlerFunc(testIS.activateHandler)
	handler.ServeHTTP(rec, req)
	if received := rec.Code; received != http.StatusOK {
		t.Errorf("handler returning incorrect status code; received %v, expected %v", received, http.StatusOK)
	}

}
*/
