package frontend

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/google/go-github/github"
	"github.com/stretchr/testify/assert"
)

var req = &http.Request{}

func init() {
	r, err := http.NewRequest("GET", "/test-url", nil)
	if err != nil {
		fmt.Printf("failure generating test request: %v", err)
	}
	*req = *r
}

func Test_httpRedirect(t *testing.T) {
	assert := assert.New(t)
	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(httpRedirect)
	handler.ServeHTTP(rec, req)

	// This redirects from the non-PROD else check in the production code.
	assert.Equal(http.StatusMovedPermanently, rec.Code, nil)
}

func Test_staticHandler(t *testing.T) {
	assert := assert.New(t)

	// More files/scenarios can be added here as desired.
	tests := []struct {
		filepath string
		result   int
	}{
		{"", http.StatusInternalServerError},
		{"website2/landing-page.html", http.StatusOK},
		{"website2/docs.html", http.StatusOK},
	}

	for i := range tests {
		rec := httptest.NewRecorder()
		staticHandler(tests[i].filepath).ServeHTTP(rec, req)

		assert.Equal(
			tests[i].result, rec.Code,
			fmt.Sprint("filepath ", tests[i].filepath),
		)
	}
}

func Test_reposHandler(t *testing.T) {
	// Dummy GitHub server to return values for ListUserInstallations.
	mux := http.NewServeMux()
	mux.HandleFunc("/user/installations/5535/repositories", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"repositories":[{"id":-65,"full_name":"contingency/chancellor"},{"id":-66,"full_name":"contingency/jedi"}]}`)
	})
	mux.HandleFunc("/repos/contingency/chancellor/labels", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `[{"name":"do-not-use"}]`)
	})
	mux.HandleFunc("/repos/contingency/jedi/labels", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `[{"name":"definitely-use"}]`)
	})

	server := httptest.NewServer(mux)
	testURL, _ := url.Parse(server.URL + "/")

	newClient = func(code string) (*github.Client, error) {
		c := github.NewClient(nil)
		c.BaseURL = testURL
		c.UploadURL = testURL
		return c, nil
	}

	assert := assert.New(t)
	handler := http.HandlerFunc(reposHandler)

	tests := []struct {
		name   string
		method string
		forms  map[string]string
		result int
	}{
		{
			"GET without state",
			"GET",
			map[string]string{"state": ""},
			http.StatusTemporaryRedirect,
		},
		{
			"passing GET",
			"GET",
			map[string]string{"state": oauthState},
			http.StatusOK,
		},
	}

	for i := range tests {
		req.Method = tests[i].method
		req.Form = url.Values{}
		for k, v := range tests[i].forms {
			req.Form.Set(k, v)
		}
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		assert.Equal(
			tests[i].result, rec.Code,
			fmt.Sprint(tests[i].name),
		)
		nums := []string{"-65", "-66"}
		for i := range nums {
			filename := nums[i] + "_github.gob"
			if _, err := os.Stat(filename); err == nil {
				os.Remove(filename)
			}
		}
	}
}

func Test_consoleHandler(t *testing.T) {
	assert := assert.New(t)
	handler := http.HandlerFunc(consoleHandler)

	tests := []struct {
		name   string
		method string
		forms  map[string]string
		result int
	}{
		{
			"reject GET method",
			"GET",
			map[string]string{"state": ""},
			http.StatusBadRequest,
		},
		{
			"POST without state",
			"POST",
			map[string]string{"state": ""},
			http.StatusTemporaryRedirect,
		},
		{
			"POST with state",
			"POST",
			map[string]string{"state": oauthState, "repo-selection": "1"},
			http.StatusOK, // TEMPORARY
		},
	}

	for i := range tests {
		req.Method = tests[i].method
		req.Form = url.Values{}
		for k, v := range tests[i].forms {
			req.Form.Set(k, v)
		}
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		assert.Equal(
			tests[i].result, rec.Code,
			fmt.Sprint(tests[i].name),
		)
		nums := []string{"-65", "-66"}
		for i := range nums {
			filename := nums[i] + "_github.gob"
			if _, err := os.Stat(filename); err == nil {
				os.Remove(filename)
			}
		}
	}
}

func Test_setupCompleteHandler(t *testing.T) {
	assert := assert.New(t)
	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(setupCompleteHandler)
	handler.ServeHTTP(rec, req)

	wanted := http.StatusOK
	received := rec.Code
	assert.Equal(
		wanted, received,
		fmt.Sprintf(
			"handler returning incorrect status code; received %v, wanted %v",
			received, wanted,
		),
	)

	setup, err := ioutil.ReadFile("website2/setup-complete.html")
	if err != nil {
		t.Errorf("Error reading from setup-complete file")
	}

	// Check that the response body is correct.
	assert.Equal(
		rec.Body.String(),
		string(setup),
		fmt.Sprintf(
			"incorrect response body\n%v",
			rec.Body.String(),
		),
	)
}
