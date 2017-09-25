package ingestor

import (
	"net/http"

	"github.com/google/go-github/github"
	"go.uber.org/zap"

	"core/utils"
)

const secretKey = "figrin-dan-and-the-modal-nodes"
var Workload = make(chan interface{}, 100)

func collectorHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		eventType := r.Header.Get("X-Github-Event")
		if eventType != "issues" && eventType != "pull_request" {
			utils.AppLog.Warn("Ignoring event", zap.String("EventType", eventType))
			return
		}
		payload, err := github.ValidatePayload(r, []byte(secretKey))
		if err != nil {
			utils.AppLog.Error("could not validate secret: ", zap.Error(err))
			return
		}
		event, err := github.ParseWebHook(github.WebHookType(r), payload)
		if err != nil {
			utils.AppLog.Error("could not parse webhook", zap.Error(err))
			return
		}
		switch v := event.(type) {
		case *github.IssuesEvent:
			Workload <- *v
		case *github.PullRequestEvent:
			Workload <- *v
		default:
			utils.AppLog.Error("Unknown", zap.ByteString("GithubEvent", payload))
		}
	})
}
