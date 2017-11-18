package frontend

import (
	"fmt"
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

	// Check that the response body is correct.
	assert.Equal(
		rec.Body.String(),
		setup,
		fmt.Sprintf(
			"incorrect response body\n%v",
			rec.Body.String(),
		),
	)
}
