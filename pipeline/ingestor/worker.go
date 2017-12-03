package ingestor

import (
	"context"

	"github.com/google/go-github/github"
	"go.uber.org/zap"

	"core/utils"
)

type Worker struct {
	ID              int
	Database        DataAccess
	RepoInitializer *RepoInitializer
	Work            chan interface{}
	Queue           chan chan interface{}
	Quit            chan bool
}

func (w *Worker) ProcessHeuprInstallationEvent(event HeuprInstallationEvent) {
	go func(e HeuprInstallationEvent) {
		switch *e.Action {
		case "created":
			w.RepoInitializer.ActivateBackend(ActivationParams{InstallationEvent: e})
			client := NewClient(*e.HeuprInstallation.AppID, *e.HeuprInstallation.ID)
			for i := 0; i < len(e.Repositories); i++ {
				githubRepo, _, err := client.Repositories.GetByID(context.Background(), *e.Repositories[i].ID)
				if err != nil {
					utils.AppLog.Error("ingestor get by id", zap.Error(err))
					return
				}
				repo := AuthenticatedRepo{Repo: githubRepo, Client: client}
				if w.RepoInitializer.RepoIntegrationExists(*repo.Repo.ID) {
					return
				}
				go w.RepoInitializer.AddRepo(repo)
				w.RepoInitializer.AddRepoIntegration(*repo.Repo.ID, *e.HeuprInstallation.AppID, *e.HeuprInstallation.ID)
				w.RepoInitializer.RaiseRepoWelcomeIssue(repo, *event.Sender.Login)
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
			client := NewClient(*e.HeuprInstallation.AppID, *e.HeuprInstallation.ID)
			for i := 0; i < len(e.RepositoriesAdded); i++ {
				repo := AuthenticatedRepo{Repo: e.RepositoriesAdded[i], Client: client}
				if w.RepoInitializer.RepoIntegrationExists(*repo.Repo.ID) {
					return
				}
				go w.RepoInitializer.AddRepo(repo)
				w.RepoInitializer.AddRepoIntegration(*repo.Repo.ID, *e.HeuprInstallation.AppID, *e.HeuprInstallation.ID)
				w.RepoInitializer.RaiseRepoWelcomeIssue(repo, *event.Sender.Login)
			}
		case "removed":
			client := NewClient(*e.HeuprInstallation.AppID, *e.HeuprInstallation.ID)
			for i := 0; i < len(e.RepositoriesRemoved); i++ {
				repo := AuthenticatedRepo{Repo: e.RepositoriesRemoved[i], Client: client}
				if !w.RepoInitializer.RepoIntegrationExists(*repo.Repo.ID) {
					return
				}
				w.RepoInitializer.RemoveRepoIntegration(*repo.Repo.ID, *e.HeuprInstallation.AppID, *e.HeuprInstallation.ID)
			}
		}
	}(event)
}

func NewWorker(id int, db DataAccess, repoInitializer *RepoInitializer, queue chan chan interface{}) Worker {
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
					// The Action that was performed. Can be one of "assigned",
					// "unassigned", "labeled", "unlabeled", "opened",
					// "edited", "milestoned", "demilestoned", "closed", or
					// "reopened".
					v.Issue.Repository = v.Repo
					if *v.Action == "edited" && *v.Issue.User.Login == "heupr[bot]" {
						//go w.ProcessHeuprInteractionIssuesEvent(v)
						if *v.Sender.Login != "heupr[bot]" && v.Issue.Assignees != nil {
							for i := 0; i < len(v.Issue.Assignees); i++ {
								if *v.Sender.Login == *v.Issue.Assignees[i].Login {
									go w.ProcessHeuprInteractionIssuesEvent(v)
									break
								}
							}
						}
						continue
					}
					w.Database.InsertIssue(*v.Issue, v.Action)
				case github.PullRequestEvent:
					//v.PullRequest.Base.Repo = v.Repo //TODO: Confirm
					w.Database.InsertPullRequest(*v.PullRequest, v.Action)
				case github.IssueCommentEvent:
					if *v.Action == "created" && *v.Issue.User.Login == "heupr[bot]" {
						if *v.Sender.Login != "heupr[bot]" && v.Issue.Assignees != nil {
							for i := 0; i < len(v.Issue.Assignees); i++ {
								if *v.Sender.Login == *v.Issue.Assignees[i].Login {
									v.Issue.Repository = v.Repo
									go w.ProcessHeuprInteractionCommentEvent(v)
									break
								}
							}
						}
					}
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

func (w *Worker) Stop() {
	go func() {
		w.Quit <- true
	}()
}
