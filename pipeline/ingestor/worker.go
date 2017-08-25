package ingestor

import (
	"core/utils"
	"github.com/google/go-github/github"
	"go.uber.org/zap"
)

type Worker struct {
	ID    int
	Db    Database
	Work  chan interface{}
	Queue chan chan interface{}
	Quit  chan bool
}

func NewWorker(id int, queue chan chan interface{}) Worker {
	return Worker{
		ID:    id,
		Db:    Database{},
		Work:  make(chan interface{}),
		Queue: queue,
		Quit:  make(chan bool),
	}
}

func (w *Worker) Start() {
	//TODO: pull this out into shared state
	w.Db.Open()
	go func() {
		for {
			w.Queue <- w.Work
			select {
			case event := <-w.Work:
				switch v := event.(type) {
				case github.IssuesEvent:
					v.Issue.Repository = v.Repo
					w.Db.InsertIssue(*v.Issue)
				case github.PullRequestEvent:
					//v.PullRequest.Base.Repo = v.Repo //TODO: Confirm
					w.Db.InsertPullRequest(*v.PullRequest)
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
