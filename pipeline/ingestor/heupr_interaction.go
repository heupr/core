package ingestor

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/go-github/github"
	"go.uber.org/zap"

	"core/utils"
)

type HeuprConfigSettings struct {
	Blacklist    []string
	StartTime    time.Time
	IgnoreLabels []string
	Email        string
	Twitter      string
}

func extractSettings(issue github.Issue) (blacklist []string, startTime time.Time, ignoreLabels []string, email string, twitter string, err error) {
	r := strings.Replace("ContributorBlacklist=\\{(.*?)\\}", "{", `"`, 1)
	r = strings.Replace(r, "}", `"`, 1)
	blacklistRegex := regexp.MustCompile(r)
	blacklistMatch := blacklistRegex.FindAllSubmatch([]byte(*issue.Body), -1)
	if len(blacklistMatch) > 0 {
		blacklist = strings.Split(string(blacklistMatch[len(blacklistMatch)-1][1]), ",")
	}

	r = strings.Replace("TriageStartTime=\\{(.*?)\\}", "{", `"`, 1)
	r = strings.Replace(r, "}", `"`, 1)
	startTimeRegex := regexp.MustCompile(r)
	startTimeMatch := startTimeRegex.FindAllSubmatch([]byte(*issue.Body), -1)
	if len(startTimeMatch) > 0 {
		startTime, err = time.Parse(time.RFC822, string(startTimeMatch[len(startTimeMatch)-1][1]))
	}

	r = strings.Replace("IgnoreLabels=\\{(.*?)\\}", "{", `"`, 1)
	r = strings.Replace(r, "}", `"`, 1)
	ignoreLabelsRegex := regexp.MustCompile(r)
	ignoreLabelsMatch := ignoreLabelsRegex.FindAllSubmatch([]byte(*issue.Body), -1)
	if len(ignoreLabelsMatch) > 0 {
		ignoreLabels = strings.Split(string(ignoreLabelsMatch[len(ignoreLabelsMatch)-1][1]), ",")
	}

	r = strings.Replace("Email=\\{(.*?)\\}", "{", `"`, 1)
	r = strings.Replace(r, "}", `"`, 1)
	emailRegex := regexp.MustCompile(r)
	emailMatch := emailRegex.FindAllSubmatch([]byte(*issue.Body), -1)
	if len(emailMatch) > 0 {
		email = string(emailMatch[len(emailMatch)-1][1])
	}

	r = strings.Replace("Twitter=\\{(.*?)\\}", "{", `"`, 1)
	r = strings.Replace(r, "}", `"`, 1)
	twitterRegex := regexp.MustCompile(r)
	twitterMatch := twitterRegex.FindAllSubmatch([]byte(*issue.Body), -1)
	if len(twitterMatch) > 0 {
		email = string(twitterMatch[len(twitterMatch)-1][1])
	}

	return blacklist, startTime, ignoreLabels, email, twitter, err
}

func (w *Worker) ProcessHeuprInteractionEvent(event github.IssuesEvent) {
	owner := *event.Issue.Repository.Owner.Login
	repo := *event.Issue.Repository.Name
	repoId := *event.Issue.Repository.ID
	number := *event.Issue.Number
	integration, err := w.Database.ReadIntegrationByRepoId(repoId)
	if err != nil {
		utils.AppLog.Error("Failed to process HeuprInterectionEvent.", zap.Error(err))
	}
	client := NewClient(integration.AppId, integration.InstallationId)

	blacklist, startTime, ignoreLabels, email, twitter, err := extractSettings(*event.Issue)
	err = err
	body := fmt.Sprintf(ConfirmationMessage, blacklist, startTime, ignoreLabels, email, twitter)
	comment := &github.IssueComment{Body: &body}
	_, _, err = client.Issues.CreateComment(context.Background(), owner, repo, number, comment)
	if err != nil {
		utils.AppLog.Error("Failed to process HeuprInterectionEvent.", zap.Error(err))
	}
}
