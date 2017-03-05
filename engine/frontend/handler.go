package frontend

import (
	"fmt"
	"github.com/google/go-github/github"
	"net/http"
)

const secretKey = "chalmun"

// var Workload = make(chan github.Issue, 100)
var Workload = make(chan github.IssuesEvent, 100)

func collectorHandler(repo string) http.Handler {
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

		// event := github.IssueEvent{}
		// err = json.Unmarshal(body, &event)
        // err = json.Unmarshal(payload, &event)
        // if err != nil {
		// 	fmt.Printf("Ignoring '%s' event with invalid payload", eventType)
		// 	http.Error(w, "Bad request", http.StatusBadRequest)
		// 	return
		// }
		// fmt.Printf("Handling '%s' event for %s", eventType, repo)

        issueEvent := event.(github.IssuesEvent)
        // issueEvent := event.(github.Issue)
        Workload <- issueEvent
        // fmt.Println(reflect.TypeOf(event)) // TEMPORARY
        // Workload <- *event.(&github.IssuesEvent).Issue
        // Workload <- event
        // Workload <- *event.(github.IssueEvent).Issue
	})
}
