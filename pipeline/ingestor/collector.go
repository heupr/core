package ingestor

import (
	"encoding/json"
	"net/http"

	"github.com/google/go-github/github"

	"go.uber.org/zap"

	"core/utils"
)

// Workaround Github API limitation. This is required to wrap HeuprInstallation
type HeuprInstallationEvent struct {
	// The action that was performed. Can be either "created" or "deleted".
	Action            *string            `json:"action,omitempty"`
	Sender            *github.User       `json:"sender,omitempty"`
	HeuprInstallation *HeuprInstallation `json:"installation,omitempty"`
	Repositories      []HeuprRepository  `json:"repositories,omitempty"`
}

// Workaround Github API limitation. go-github is missing repositories field
type HeuprInstallation struct {
	ID              *int         `json:"id,omitempty"`
	Account         *github.User `json:"account,omitempty"`
	AppID           *int         `json:"app_id,omitempty"`
	AccessTokensURL *string      `json:"access_tokens_url,omitempty"`
	RepositoriesURL *string      `json:"repositories_url,omitempty"`
	HTMLURL         *string      `json:"html_url,omitempty"`
}

type HeuprRepository struct {
	ID       *int    `json:"id,omitempty"`
	Name     *string `json:"name,omitempty"`
	FullName *string `json:"full_name,omitempty"`
}

const secretKey = "figrin-dan-and-the-modal-nodes"

var Workload = make(chan interface{}, 100)

func collectorHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		eventType := r.Header.Get("X-Github-Event")
		if eventType != "issues" && eventType != "pull_request" && eventType != "installation" {
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
		case *github.InstallationEvent:
			e := &HeuprInstallationEvent{}
			err := json.Unmarshal(payload, &e)
			if err != nil {
				utils.AppLog.Error("could not parse webhook", zap.Error(err))
				return
			}
			Workload <- *e
		default:
			utils.AppLog.Error("Unknown", zap.ByteString("GithubEvent", payload))
		}
	})
}
