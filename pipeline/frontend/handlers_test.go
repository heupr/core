package frontend

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

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
	wanted := http.StatusMovedPermanently
	received := rec.Code

	assert.Equal(
		wanted, received,
		fmt.Sprintf(
			"handler returning incorrect status code; received %v, wanted %v",
			received, wanted,
		),
	)
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
			fmt.Sprint("filepath", tests[i].filepath),
		)
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
