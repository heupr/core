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
  Integration  Integration
	IgnoreUsers  []string
	StartTime    time.Time
	IgnoreLabels []string
	Email        string
	Twitter      string
}

func extractSettings(issue github.Issue) (ignoreUsers []string, startTime time.Time, ignoreLabels []string, email string, twitter string, err error) {
	r := strings.Replace("IgnoreUsers=\\{(.*?)\\}", "{", `"`, 1)
	r = strings.Replace(r, "}", `"`, 1)
	ignoreUsersRegex := regexp.MustCompile(r)
	ignoreUsersMatch := ignoreUsersRegex.FindAllSubmatch([]byte(*issue.Body), -1)
	if len(ignoreUsersMatch) > 0 {
		ignoreUsers = strings.Split(string(ignoreUsersMatch[len(ignoreUsersMatch)-1][1]), ",")
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
		twitter = string(twitterMatch[len(twitterMatch)-1][1])
	}

	return ignoreUsers, startTime, ignoreLabels, email, twitter, err
}

func (w *Worker) ProcessHeuprInteractionCommentEvent(event github.IssueCommentEvent) {
  owner := *event.Repo.Owner.Login
  repo := *event.Repo.Name
  repoId := *event.Repo.ID
  number := *event.Issue.Number
  integration, err := w.Database.ReadIntegrationByRepoId(repoId)
  if err != nil {
    utils.AppLog.Error("Failed to process HeuprInterectionEvent.", zap.Error(err))
  }
  client := NewClient(integration.AppId, integration.InstallationId)

  if strings.Contains(*event.Comment.Body, "no") || strings.Contains(*event.Comment.Body, "No") {
    body := fmt.Sprintf(HoldOnMessage, *event.Sender.Login)
    comment := &github.IssueComment{Body: &body}
    _, _, err = client.Issues.CreateComment(context.Background(), owner, repo, number, comment)
    if err != nil {
      utils.AppLog.Error("Failed to process HeuprInterectionEvent.", zap.Error(err))
    }
    return
  }

  if !strings.Contains(*event.Comment.Body, "yes") && !strings.Contains(*event.Comment.Body, "Yes") {
    return
  }

  //Duplicate Logic. It gets the job done for validation & settings extraction)
  ignoreUsers, startTime, ignoreLabels, email, twitter, err := extractSettings(*event.Issue)
  var body string
  if err != nil {
    body = fmt.Sprintf(ConfirmationErrMessage, *event.Sender.Login, err, ignoreUsers, startTime, ignoreLabels, email, twitter)
  } else {
    body = fmt.Sprintf(AppliedSettingsMessage, *event.Sender.Login, ignoreUsers, startTime, ignoreLabels, email, twitter)
  }
  comment := &github.IssueComment{Body: &body}
  _, _, err = client.Issues.CreateComment(context.Background(), owner, repo, number, comment)
  if err != nil {
    utils.AppLog.Error("Failed to process HeuprInterectionEvent.", zap.Error(err))
  }

  settings := HeuprConfigSettings{Integration: *integration, IgnoreUsers: ignoreUsers, StartTime: startTime, IgnoreLabels: ignoreLabels, Email: email, Twitter: twitter}
  w.Database.InsertRepositoryIntegrationSettings(settings)
  //Workaround: This causes the backend to kick in and pull in the latest settings.
	action := "opened"
	w.Database.InsertIssue(*event.Issue, &action)
}

func (w *Worker) ProcessHeuprInteractionIssuesEvent(event github.IssuesEvent) {
  //TODO: Add user validation
	owner := *event.Issue.Repository.Owner.Login
	repo := *event.Issue.Repository.Name
	repoId := *event.Issue.Repository.ID
	number := *event.Issue.Number
	integration, err := w.Database.ReadIntegrationByRepoId(repoId)
	if err != nil {
		utils.AppLog.Error("Failed to process HeuprInterectionEvent.", zap.Error(err))
	}
	client := NewClient(integration.AppId, integration.InstallationId)

	ignoreUsers, startTime, ignoreLabels, email, twitter, err := extractSettings(*event.Issue)
  var body string
  if err == nil {
    body = fmt.Sprintf(ConfirmationMessage, *event.Sender.Login, ignoreUsers, startTime, ignoreLabels, email, twitter)
  } else {
    body = fmt.Sprintf(ConfirmationErrMessage, *event.Sender.Login, err, ignoreUsers, startTime, ignoreLabels, email, twitter)
  }
	comment := &github.IssueComment{Body: &body}
	_, _, err = client.Issues.CreateComment(context.Background(), owner, repo, number, comment)
	if err != nil {
		utils.AppLog.Error("Failed to process HeuprInterectionEvent.", zap.Error(err))
	}
}
