package frontend

import (
	"encoding/gob"
	"fmt"
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

func Test_updateStorage(t *testing.T) {
	s := storage{
		Name: "watto/junkshop",
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
			http.StatusUnauthorized,
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
			filename := nums[i] + ".gob"
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
		{
			"breaking gob file",
			"POST",
			"rejectGob",
			map[string]string{"repo-selection": "-65"},
			http.StatusInternalServerError,
		},
		{
			"passing POST",
			"POST",
			"-65.gob",
			map[string]string{"repo-selection": "-65"},
			http.StatusOK,
		},
	}

	for i := range tests {
		f, err := os.Create(tests[i].file)
		defer os.Remove(tests[i].file)
		if err != nil {
			t.Error("failure creating test gob file")
		}
		s := storage{
			Name: "contingency/chancellor",
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

		assert.Equal(
			tests[i].result, rec.Code,
			fmt.Sprint(tests[i].name),
		)
	}
}

func Test_updateSettings(t *testing.T) {
	assert := assert.New(t)
	s := storage{
		Name: "watto/junkshop",
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

func Test_setupCompleteHandler(t *testing.T) {
	assert := assert.New(t)
	handler := http.HandlerFunc(setupCompleteHandler)

	tests := []struct {
		name   string
		method string
		file   string
		result int
	}{
		{
			"Rejected GET request",
			"GET",
			"rejectGET",
			http.StatusBadRequest,
		},
		{
			"Accepted POST request",
			"POST",
			"-65.gob",
			http.StatusOK,
		},
	}

	for i := range tests {
		f, err := os.Create(tests[i].file)
		defer os.Remove(tests[i].file)
		if err != nil {
			t.Error("failure creating test gob file")
		}
		s := storage{
			Name: "contingency/chancellor",
		}
		encoder := gob.NewEncoder(f)
		if err := encoder.Encode(s); err != nil {
			t.Error("failed encoding gob file")
		}

		req.Method = tests[i].method
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		assert.Equal(
			tests[i].result,
			rec.Code,
			tests[i].name,
		)
	}
}
