package frontend

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt" // TEMPORARY
	"github.com/google/go-github/github"
	"io/ioutil"
	"net/http"
	"strings"
)

var Workload = make(chan github.Issue, 100)

func collectorHandler() http.Handler {
	// NOTE: Temporarily removed the "secret" argument - eventually implement
	//       for security purposes.
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//TODO: Pass in secret
		secret := "chalmun's-spaceport-cantina"
		eventType := r.Header.Get("X-Github-Event")
		if eventType != "issues" {
			fmt.Printf("Ignoring '%s' event", eventType)
			return
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if secret != "" {
			ok := false
			for _, sig := range strings.Fields(r.Header.Get("X-Hub-Signature")) {
				if !strings.HasPrefix(sig, "sha1=") {
					continue
				}
				sig = strings.TrimPrefix(sig, "sha1=")
				mac := hmac.New(sha1.New, []byte(secret))
				mac.Write(body)
				if sig == hex.EncodeToString(mac.Sum(nil)) {
					ok = true
					break
				}
			}
			if !ok {
				fmt.Printf("Ignoring '%s' event with incorrect signature", eventType)
				return
			}
		}

		event := github.IssueEvent{}
		err = json.Unmarshal(body, &event)
		if err != nil {
			fmt.Printf("Ignoring '%s' event with invalid payload", eventType)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		//TODO: Get Repo Name from Issue object
		repo := "TEST_REPO_NAME"
		fmt.Printf("Handling '%s' event for %s", eventType, repo)

		Workload <- *event.Issue
	})
}
