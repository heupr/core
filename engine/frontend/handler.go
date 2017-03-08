package frontend

import (
	"fmt"
	"net/http"

	"github.com/google/go-github/github"
)

var Workload = make(chan github.IssuesEvent, 100)

func collectorHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		eventType := r.Header.Get("X-Github-Event")
		if eventType != "issues" {
			fmt.Printf("Ignoring '%v' event", eventType)
			return
		}
		payload, err := github.ValidatePayload(r, []byte(secretKey))
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		event, err := github.ParseWebHook(github.WebHookType(r), payload)

		if err != nil {
			fmt.Printf("Could not parse webhook %v", err)
			return
		}
		Workload <- *event.(*github.IssuesEvent)
	})
}
