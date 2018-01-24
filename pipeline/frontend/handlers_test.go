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
	"path/filepath"
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

	newServerToServerClient = func(appId, installationId int) (*github.Client, error) {
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

		assert.Equal(tests[i].result, rec.Code, fmt.Sprint(tests[i].name))
		nums := []string{"-65", "-66"}
		for i := range nums {
			filename := "gob/" + nums[i] + ".gob"
			if _, err := os.Stat(filename); err == nil {
				os.Remove(filename)
			}
		}
	}
}

func Test_generateWalkFunc(t *testing.T) {
	assert := assert.New(t)

	file := "-66.gob"
	found := ""
	_, err := os.Create(file)
	if err != nil {
		t.Errorf("create file error", err)
	}
	defer os.Remove(file)

	wkFn := generateWalkFunc(&found, "-66")
	filepath.Walk(".", wkFn)

	assert.Equal(file, found, "generate walk function files not matching")
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
			http.StatusForbidden,
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
		f, err := os.Create(tests[i].file)
		if err != nil {
			t.Error("failure creating test gob file")
		}
		defer os.Remove(tests[i].file)
		s := storage{
			FullName: "contingency/chancellor",
		}
		encoder := gob.NewEncoder(f)
		if err := encoder.Encode(s); err != nil {
			t.Error("failed encoding gob file")
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

	buf := new(bytes.Buffer)
	err = tmpl.Execute(buf, "")
	if err != nil {
		t.Error("failure executing complete template")
	}

	str := buf.String()

	if ok := strings.Contains(str, "<title>Heupr</title>"); !ok {
		t.Error("base template html missing value")
	}
	if ok := strings.Contains(str, `{{ template "body" . }}`); ok {
		t.Error("error executing nested complete template")
	}
	if ok := strings.Contains(str, "<h2>Awesome! Setup is complete!</h2>"); !ok {
		t.Error("complete html missing value")
	}

	for i := range tests {
		f, err := os.Create(tests[i].file)
		defer os.Remove(tests[i].file)

		if err != nil {
			t.Error("failure creating test gob file")
		}
		s := storage{
			FullName: "contingency/chancellor",
		}
		encoder := gob.NewEncoder(f)
		if err := encoder.Encode(s); err != nil {
			t.Error("failed encoding gob file")
		}

		req.Method = tests[i].method
		rec := httptest.NewRecorder()

		session, err := store.Get(req, sessionName)
		if err != nil {
			t.Error("failure creating test session")
		}
		session.Values["repoID"] = tests[i].repoID

		handler.ServeHTTP(rec, req)

		assert.Equal(tests[i].result, rec.Code, tests[i].name)

	}
}
