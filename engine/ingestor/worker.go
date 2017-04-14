package ingestor

import (
	"fmt"
	"coralreefci/engine/gateway/conflation"
	"github.com/google/go-github/github"
)

type Worker struct {
	ID     int
	Work   chan github.IssuesEvent
	Queue  chan chan github.IssuesEvent
	Quit   chan bool
}

func NewWorker(id int, queue chan chan github.IssuesEvent) Worker {
	return Worker{
		ID:     id,
		Work:   make(chan github.IssuesEvent),
		Queue:  queue,
		Quit:   make(chan bool),
	}
}

func (w *Worker) Start() {
	go func() {
		for {
			w.Queue <- w.Work
			select {
			case issuesEvent := <-w.Work:
				if issuesEvent.Issue.ClosedAt != nil {
					expandedIssue := conflation.ExpandedIssue{Issue: conflation.CRIssue{*issuesEvent.Issue, []int{}, []conflation.CRPullRequest{}}}
					fmt.Println("ID ", *expandedIssue.Issue.ID)
				} else {
					expandedIssue := conflation.ExpandedIssue{Issue: conflation.CRIssue{*issuesEvent.Issue, []int{}, []conflation.CRPullRequest{}}}
          fmt.Println("ID ", *expandedIssue.Issue.ID)
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
