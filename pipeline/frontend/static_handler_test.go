package frontend

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func newReq(method string) *http.Request {
	req, err := http.NewRequest(method, "/test-url", nil)
	if err != nil {
		fmt.Printf("failure generating test request: %v", err)
	}
	return req
}

func Test_httpRedirect(t *testing.T) {
	assert := assert.New(t)
	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(httpRedirect)
	req := newReq("GET")
	handler.ServeHTTP(rec, req)

	// This redirects from the non-PROD else check in the production code.
	assert.Equal(http.StatusMovedPermanently, rec.Code, nil)
}

func Test_render(t *testing.T) {
	assert := assert.New(t)

	// More files/scenarios can be added here as desired.
	tests := []struct {
		filepath string
		result   int
	}{
		{"", http.StatusInternalServerError},
		{"templates/home.html", http.StatusOK},
		{"templates/docs.html", http.StatusOK},
	}

	baseHTML = "templates/base.html"
	for i := range tests {
		rec := httptest.NewRecorder()
		req := newReq("GET")
		render(tests[i].filepath).ServeHTTP(rec, req)

		assert.Equal(
			tests[i].result, rec.Code,
			fmt.Sprint("filepath ", tests[i].filepath),
		)
	}
}
