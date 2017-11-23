package ingestor

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/go-github/github"
	"go.uber.org/zap"

	"core/pipeline/gateway"
	"core/utils"
)

type ActivationParams struct {
	InstallationEvent HeuprInstallationEvent `json:"installation_event,omitempty"`
	Limit             time.Time              `json:"limit,omitempty"`
}

type RepoInitializer struct {
	Database   DataAccess
	HttpClient http.Client
}

func (r *RepoInitializer) AddRepo(authRepo AuthenticatedRepo) {
	newGateway := gateway.Gateway{
		Client:      authRepo.Client,
		UnitTesting: false,
	}
	githubIssues, err := newGateway.GetClosedIssues(*authRepo.Repo.Owner.Login, *authRepo.Repo.Name)
	if err != nil {
		// TODO: Proper error handling should be evaluated for this method;
		// possibly adjust to return an error variable.
		utils.AppLog.Error(err.Error())
	}
	// The Repo struct needs to be added to the Issue struct body - this is
	// possibly a bug in the GitHub API.
	for i := 0; i < len(githubIssues); i++ {
		githubIssues[i].Repository = authRepo.Repo
	}
	githubPulls, err := newGateway.GetClosedPulls(*authRepo.Repo.Owner.Login, *authRepo.Repo.Name)
	if err != nil {
		utils.AppLog.Error(err.Error())
	}
	r.Database.BulkInsertIssuesPullRequests(githubIssues, githubPulls)
}

func (r *RepoInitializer) RepositoryIntegrationExists(repoID int, appID int, installationID int) bool {
	_, err := r.Database.ReadIntegrationByRepoID(repoID)
	switch {
	case err == sql.ErrNoRows:
		return false
	case err != nil:
		utils.AppLog.Error("integration read error", zap.Error(err))
		return false
	default:
		return true
	}
}

func (r *RepoInitializer) AddRepositoryIntegration(repoID int, appID int, installationID int) {
	r.Database.InsertRepositoryIntegration(repoID, appID, installationID)
}

func (r *RepoInitializer) RemoveRepositoryIntegration(repoID int, appID int, installationID int) {
	r.Database.DeleteRepositoryIntegration(repoID, appID, installationID)
}

func (r *RepoInitializer) ObliterateIntegration(appID int, installationID int) {
	r.Database.ObliterateIntegration(appID, installationID)
}

func (r *RepoInitializer) RaiseRepositoryWelcomeIssue(authRepo AuthenticatedRepo, assignee string) {
	welcomeScreen := &github.IssueRequest{Title: github.String(WelcomeTitle), Body: github.String(fmt.Sprintf(WelcomeBody, time.Now().Format(time.RFC822))), Assignees: &[]string{assignee}}
	authRepo.Client.Issues.Create(context.Background(), *authRepo.Repo.Owner.Login, *authRepo.Repo.Name, welcomeScreen)
}

func (r *RepoInitializer) ActivateBackend(params ActivationParams) {
	payload, err := json.Marshal(params)
	if err != nil {
		utils.AppLog.Error("failed to marshal json", zap.Error(err))
		return
	}
	req, err := http.NewRequest("POST", utils.Config.BackendActivationEndpoint, bytes.NewBuffer(payload))
	if err != nil {
		utils.AppLog.Error("failed to create http request", zap.Error(err))
		return
	}
	resp, err := r.HttpClient.Do(req)
	if err != nil {
		utils.AppLog.Error("failed internal post", zap.Error(err))
		return
	} else {
		defer resp.Body.Close()
	}
}
