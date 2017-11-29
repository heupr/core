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
	HTTPClient http.Client
}

func (r *RepoInitializer) AddRepo(authRepo AuthenticatedRepo) {
	newGateway := gateway.Gateway{
		Client:      authRepo.Client,
		UnitTesting: false,
	}
	issues, err := newGateway.GetClosedIssues(
		*authRepo.Repo.Owner.Login,
		*authRepo.Repo.Name,
	)
	if err != nil {
		utils.AppLog.Error("add repo get issues", zap.Error(err))
	}
	// Adding the Repo to the Issue is to cover a GitHub API deficiency.
	for i := 0; i < len(issues); i++ {
		issues[i].Repository = authRepo.Repo
	}
	pulls, err := newGateway.GetClosedPulls(
		*authRepo.Repo.Owner.Login,
		*authRepo.Repo.Name,
	)
	if err != nil {
		utils.AppLog.Error("add repo get pulls", zap.Error(err))
	}
	r.Database.BulkInsertIssuesPullRequests(issues, pulls)
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
	resp, err := r.HTTPClient.Do(req)
	defer resp.Body.Close()
	if err != nil {
		utils.AppLog.Error("failed internal post", zap.Error(err))
		return
	}
}
