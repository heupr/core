package ingestor

import (
	"github.com/google/go-github/github"
)

type Worker struct {
	ID    int
	Db    Database
	Work  chan github.IssuesEvent
	Queue chan chan github.IssuesEvent
	Quit  chan bool
}

func NewWorker(id int, queue chan chan github.IssuesEvent) Worker {
	return Worker{
		ID:    id,
		Db:    Database{},
		Work:  make(chan github.IssuesEvent),
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
			case issuesEvent := <-w.Work:
				issuesEvent.Issue.Repository = issuesEvent.Repo
				w.Db.InsertIssue(*issuesEvent.Issue)
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
