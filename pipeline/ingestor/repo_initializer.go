package ingestor

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"go.uber.org/zap"

	"core/pipeline/gateway"
	"core/utils"
)

type ActivationParams struct {
	InstallationEvent HeuprInstallationEvent `json:"installation_event,omitempty"`
	Limit             time.Time              `json:"limit,omitempty"`
}

type RepoInitializer struct {
	Database   *Database
	HttpClient http.Client
}

func (r *RepoInitializer) AddRepo(authRepo AuthenticatedRepo) {
	newGateway := gateway.Gateway{
		Client:      authRepo.Client,
		UnitTesting: false,
	}
	githubIssues, err := newGateway.GetIssues(*authRepo.Repo.Owner.Login, *authRepo.Repo.Name)
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
	githubPulls, err := newGateway.GetPullRequests(*authRepo.Repo.Owner.Login, *authRepo.Repo.Name)
	if err != nil {
		utils.AppLog.Error(err.Error())
	}
	r.Database.BulkInsertIssuesPullRequests(githubIssues, githubPulls)
}

func (r *RepoInitializer) AddRepositoryIntegration(repoId int, appId int, installationId int) {
	r.Database.InsertIntegration(repoId, appId, installationId)
}

func (r *RepoInitializer) RemoveRepositoryIntegration(repoId int, appId int, installationId int) {
	//TODO:
}

func (r *RepoInitializer) ObliterateIntegration(repoId int, appId int, installationId int) {
	//TODO:
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
