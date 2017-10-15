package ingestor

import (
	"context"

	"net/http"

	"github.com/bradleyfalzon/ghinstallation"
	"github.com/google/go-github/github"

	"go.uber.org/zap"

	"core/utils"
)

type Worker struct {
	ID              int
	Database        *Database
	RepoInitializer *RepoInitializer
	Work            chan interface{}
	Queue           chan chan interface{}
	Quit            chan bool
}

func NewWorker(id int, db *Database, repoInitializer *RepoInitializer, queue chan chan interface{}) Worker {
	return Worker{
		ID:              id,
		Database:        db,
		RepoInitializer: repoInitializer,
		Work:            make(chan interface{}),
		Queue:           queue,
		Quit:            make(chan bool),
	}
}

func (w *Worker) Start() {
	go func() {
		for {
			w.Queue <- w.Work
			select {
			case event := <-w.Work:
				switch v := event.(type) {
				case github.IssuesEvent:
					//The Action that was performed. Can be one of "assigned", "unassigned", "labeled", "unlabeled", "opened", "edited", "milestoned", "demilestoned", "closed", or "reopened".
					v.Issue.Repository = v.Repo
					w.Database.InsertIssue(*v.Issue, v.Action)
				case github.PullRequestEvent:
					//v.PullRequest.Base.Repo = v.Repo //TODO: Confirm
					w.Database.InsertPullRequest(*v.PullRequest, v.Action)
				case HeuprInstallationEvent:
					w.ProcessHeuprInstallationEvent(v)
				case HeuprInstallationRepositoriesEvent:
					w.ProcessHeuprInstallationRepositoriesEvent(v)
				default:
					utils.AppLog.Error("Unknown", zap.Any("GithubEvent", v))
				}
			case <-w.Quit:
				return
			}
		}
	}()
}

func (w *Worker) ProcessHeuprInstallationEvent(event HeuprInstallationEvent) {
	go func(e HeuprInstallationEvent) {
		switch *e.Action {
		case "created":
			w.RepoInitializer.ActivateBackend(ActivationParams{InstallationEvent: e})
			// Wrap the shared transport for use with the Github Installation.
			itr, err := ghinstallation.NewKeyFromFile(http.DefaultTransport, *e.HeuprInstallation.AppID, *e.HeuprInstallation.ID, "heupr.2017-10-04.private-key.pem")
			if err != nil {
				utils.AppLog.Error("could not obtain github installation key", zap.Error(err))
				return
			}
			// Use installation transport with client.
			client := github.NewClient(&http.Client{Transport: itr})
			for i := 0; i < len(e.Repositories); i++ {
				githubRepo, _, err := client.Repositories.GetByID(context.Background(), *e.Repositories[i].ID)
				if err != nil {
					utils.AppLog.Error("ingestor get by id", zap.Error(err))
					return
				}
				repo := AuthenticatedRepo{Repo: githubRepo, Client: client}
				if w.RepoInitializer.RepositoryIntegrationExists(*repo.Repo.ID, *e.HeuprInstallation.AppID, *e.HeuprInstallation.ID) {
					return
				}
				go w.RepoInitializer.AddRepo(repo)
				go w.RepoInitializer.AddRepositoryIntegration(*repo.Repo.ID, *e.HeuprInstallation.AppID, *e.HeuprInstallation.ID)
			}
		case "deleted":
			w.RepoInitializer.ObliterateIntegration(*e.HeuprInstallation.AppID, *e.HeuprInstallation.ID)
		}
	}(event)
}

func (w *Worker) ProcessHeuprInstallationRepositoriesEvent(event HeuprInstallationRepositoriesEvent) {
	go func(e HeuprInstallationRepositoriesEvent) {
		switch *e.Action {
		case "added":
			// Wrap the shared transport for use with the Github Installation.
			itr, err := ghinstallation.NewKeyFromFile(http.DefaultTransport, *e.HeuprInstallation.AppID, *e.HeuprInstallation.ID, "heupr.2017-10-04.private-key.pem")
			if err != nil {
				utils.AppLog.Error("could not obtain github installation key", zap.Error(err))
				return
			}
			// Use installation transport with client.
			client := github.NewClient(&http.Client{Transport: itr})
			for i := 0; i < len(e.RepositoriesAdded); i++ {
				repo := AuthenticatedRepo{Repo: e.RepositoriesAdded[i], Client: client}
				if w.RepoInitializer.RepositoryIntegrationExists(*repo.Repo.ID, *e.HeuprInstallation.AppID, *e.HeuprInstallation.ID) {
					return
				}
				go w.RepoInitializer.AddRepo(repo)
				go w.RepoInitializer.AddRepositoryIntegration(*repo.Repo.ID, *e.HeuprInstallation.AppID, *e.HeuprInstallation.ID)
			}
		case "removed":
			// Wrap the shared transport for use with the Github Installation.
			itr, err := ghinstallation.NewKeyFromFile(http.DefaultTransport, *e.HeuprInstallation.AppID, *e.HeuprInstallation.ID, "heupr.2017-10-04.private-key.pem")
			if err != nil {
				utils.AppLog.Error("could not obtain github installation key", zap.Error(err))
				return
			}
			// Use installation transport with client.
			client := github.NewClient(&http.Client{Transport: itr})
			for i := 0; i < len(e.RepositoriesRemoved); i++ {
				repo := AuthenticatedRepo{Repo: e.RepositoriesRemoved[i], Client: client}
				if !w.RepoInitializer.RepositoryIntegrationExists(*repo.Repo.ID, *e.HeuprInstallation.AppID, *e.HeuprInstallation.ID) {
					return
				}
				go w.RepoInitializer.RemoveRepositoryIntegration(*repo.Repo.ID, *e.HeuprInstallation.AppID, *e.HeuprInstallation.ID)
			}
		}
	}(event)
}

func (w *Worker) Stop() {
	go func() {
		w.Quit <- true
	}()
}
