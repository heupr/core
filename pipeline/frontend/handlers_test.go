package frontend

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"html/template"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/google/go-github/github"
	"github.com/gorilla/sessions"
	"github.com/stretchr/testify/assert"
)

var req = &http.Request{}

func init() {
	r, err := http.NewRequest("GET", "/test-url", nil)
	if err != nil {
		fmt.Printf("failure generating test request: %v", err)
	}
	*req = *r

	templatePath = ""
}

func Test_updateStorage(t *testing.T) {
	s := storage{
		FullName: "watto/junkshop",
		Buckets: map[string][]label{
			"cost": []label{
				label{
					Name:     "less-100-peggats",
					Selected: true,
				},
				label{
					Name:     "equal-100-peggats",
					Selected: true,
				},
				label{
					Name:     "junk",
					Selected: false,
				},
			},
		},
	}
	labels := []string{
		"less-100-peggats",
		"equal-100-peggats",
		"more-100-peggats",
		"junk",
	}
	updateStorage(&s, labels)
	assert := assert.New(t)
	assert.Equal(
		len(labels),
		len(s.Buckets["cost"]),
		"update storage failing",
	)

}

// NOTE: This may be refactored into the full Test_repos function.
func Test_reposHTML(t *testing.T) {
	tmpl, err := template.ParseFiles(
		"templates/base.html",
		"templates/repos.html",
	)

	repos := []struct {
		Repos map[int64]string
	}{
		{
			Repos: make(map[int64]string),
		},
		{
			Repos: map[int64]string{1: ""},
		},
		{
			Repos: map[int64]string{
				2: "fode",
				3: "beed",
			},
		},
	}

	tests := []struct {
		data    map[string]interface{}
		check   string
		result  bool
		message string
	}{
		{
			map[string]interface{}{
				"storage": repos[0],
				"domain":  "test-domain",
			},
			"<option ",
			false,
			"expected no <option> to be populated",
		},
		{
			map[string]interface{}{
				"storage": repos[1],
				"domain":  "test-domain",
			},
			`<option value="1"></option>`,
			true,
			"expected no value to populate in option",
		},
		{
			map[string]interface{}{
				"storage": repos[2],
				"domain":  "test-domain",
			},
			`<option value="2">fode</option>`,
			true,
			"expected proper value/text for option",
		},
	}

	for i := range tests {
		buf := new(bytes.Buffer)
		err = tmpl.Execute(buf, tests[i].data)
		if err != nil {
			t.Error(err)
		}
		str := buf.String()

		if strings.Contains(str, tests[i].check) != tests[i].result {
			t.Error("repos template error:", tests[i].message)
		}
	}
}

func Test_repos(t *testing.T) {
	// Dummy GitHub server to return values for ListUserInstallations.
	mux := http.NewServeMux()
	mux.HandleFunc("/user/installations/1234/repositories", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"repositories":[{"id":-65,"full_name":"contingency/chancellor"},{"id":-66,"full_name":"contingency/jedi"}]}`)
	})
	mux.HandleFunc("/repos/contingency/chancellor/labels", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `[{"name":"do-not-use"}]`)
	})
	mux.HandleFunc("/repos/contingency/jedi/labels", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `[{"name":"definitely-use"}]`)
	})
	mux.HandleFunc("/user/installations", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"installations":[{"id":1234,"app_id":6807}]}`)
	})

	server := httptest.NewServer(mux)
	testURL, _ := url.Parse(server.URL + "/")

	newUserToServerClient = func(code string) (*github.Client, error) {
		c := github.NewClient(nil)
		c.BaseURL = testURL
		c.UploadURL = testURL
		return c, nil
	}

	newServerToServerClient = func(appId int, installationId int64) (*github.Client, error) {
		c := github.NewClient(nil)
		c.BaseURL = testURL
		c.UploadURL = testURL
		return c, nil
	}

	assert := assert.New(t)
	handler := http.HandlerFunc(repos)

	tests := []struct {
		name   string
		method string
		forms  map[string]string
		result int
	}{
		{
			"POST which gets rejected",
			"POST",
			map[string]string{"state": ""},
			http.StatusBadRequest,
		},
		{
			"GET without state",
			"GET",
			map[string]string{"state": ""},
			http.StatusForbidden,
		},
		{
			"passing GET",
			"GET",
			map[string]string{"state": "something"},
			http.StatusOK,
		},
	}

	gobPath = "cmd/gob/"

	for i := range tests {
		req.Method = tests[i].method
		req.Form = url.Values{}

		for k, v := range tests[i].forms {
			req.Form.Set(k, v)
		}

		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		assert.Equal(tests[i].result, rec.Code, fmt.Sprint(tests[i].name))
		nums := []string{"-65", "-66"}
		for i := range nums {
			filename := "cmd/gob/" + nums[i] + ".gob"
			if _, err := os.Stat(filename); err == nil {
				os.Remove(filename)
			}
		}
	}
}

func Test_console(t *testing.T) {
	assert := assert.New(t)
	handler := http.HandlerFunc(console)

	tests := []struct {
		name   string
		method string
		file   string
		forms  map[string]string
		result int
	}{
		{
			"reject GET method",
			"GET",
			"rejectGET",
			make(map[string]string),
			http.StatusBadRequest,
		},
		// NOTE: Temporarily commented out due to hard-coded temporary values
		// in the production code.
		// {
		// 	"breaking gob file",
		// 	"GET",
		// 	"rejectGob",
		// 	map[string]string{"state": oauthState, "repo-selection": "-65"},
		// 	http.StatusInternalServerError,
		// },
		{
			"passing POST method",
			"POST",
			"-65.gob",
			map[string]string{"repo-selection": "-65"},
			http.StatusOK,
		},
	}

	for i := range tests {
		f, err := os.Create("cmd/gob/" + tests[i].file)
		if err != nil {
			t.Error(err)
		}
		defer os.Remove("cmd/gob/" + tests[i].file)
		s := storage{
			FullName: "contingency/chancellor",
		}
		encoder := gob.NewEncoder(f)
		if err := encoder.Encode(s); err != nil {
			t.Error(err)
		}

		req.Method = tests[i].method
		req.Form = url.Values{}
		if len(tests[i].forms) > 0 {
			for k, v := range tests[i].forms {
				req.Form.Set(k, v)
			}
		}
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		assert.Equal(tests[i].result, rec.Code, fmt.Sprint(tests[i].name))
	}
}

func Test_updateSettings(t *testing.T) {
	assert := assert.New(t)
	s := storage{
		FullName: "watto/junkshop",
		Buckets: map[string][]label{
			"cost": []label{
				label{
					Name:     "less-100-peggats",
					Selected: true,
				},
				label{
					Name:     "junk",
					Selected: true,
				},
				label{
					Name:     "equal-100-peggats",
					Selected: true,
				},
				label{
					Name:     "more-100-peggats",
					Selected: false,
				},
			},
		},
	}
	labels := map[string][]string{
		"cost": []string{
			"less-100-peggats",
			"equal-100-peggats",
			"more-100-peggats",
		},
	}
	updateSettings(&s, labels)

	count := func(s storage) (cnt int) {
		for i := range s.Buckets["cost"] {
			if s.Buckets["cost"][i].Selected {
				cnt++
			}
		}
		return
	}(s)

	assert.Equal(
		(len(labels["cost"])),
		count,
		"trying this out",
	)
}

func Test_complete(t *testing.T) {
	store = sessions.NewCookieStore([]byte("test-sesions"))
	assert := assert.New(t)
	handler := http.HandlerFunc(complete)

	tests := []struct {
		name   string
		method string
		repoID string
		file   string
		result int
	}{
		{
			"Rejected GET request",
			"GET",
			"",
			"rejectGET",
			http.StatusBadRequest,
		},
		{
			"Accepted POST request",
			"POST",
			"-65",
			"-65.gob",
			http.StatusOK,
		},
	}

	tmpl, err := template.ParseFiles(
		"templates/base.html",
		"templates/complete.html",
	)

	// NOTE: This should be pulled out into a separate test function.
	buf := new(bytes.Buffer)
	data := map[string]interface{}{"domain": "test-domain"}
	err = tmpl.Execute(buf, data)
	if err != nil {
		t.Error(err)
	}

	str := buf.String()

	if strings.Contains(str, "{{") || strings.Contains(str, "}}") {
		t.Error(err)
	}

	for i := range tests {
		f, err := os.Create("cmd/gob/" + tests[i].file)
		if err != nil {
			t.Error(err)
		}
		defer os.Remove("cmd/gob/" + tests[i].file)

		s := storage{
			FullName: "contingency/chancellor",
		}
		encoder := gob.NewEncoder(f)
		if err := encoder.Encode(s); err != nil {
			t.Error(err)
		}

		req.Method = tests[i].method
		rec := httptest.NewRecorder()

		session, err := store.Get(req, sessionName)
		if err != nil {
			t.Error(err)
		}
		session.Values["repoID"] = tests[i].repoID

		handler.ServeHTTP(rec, req)

		assert.Equal(tests[i].result, rec.Code, tests[i].name)

	}
}
