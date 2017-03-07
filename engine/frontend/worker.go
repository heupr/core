package frontend

import (
	"fmt"

	"github.com/google/go-github/github"

	"coralreefci/engine/assignment"
	"coralreefci/engine/gateway/conflation"
	"coralreefci/models"
)

type Worker struct {
	ID     int
	Work   chan github.IssuesEvent
	Queue  chan chan github.IssuesEvent
	Models map[int]models.Model
	Quit   chan bool
}

func NewWorker(id int, queue chan chan github.IssuesEvent) Worker {
	return Worker{
		ID:     id,
		Work:   make(chan github.IssuesEvent),
		Queue:  queue,
		Models: make(map[int]models.Model),
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
					//TODO: Decide how we want to handle PR's
					//TODO: Call conflate (only on closed issues)
					//TODO: Remove Printlns
					expandedIssue := conflation.ExpandedIssue{Issue: conflation.CRIssue{*issuesEvent.Issue, []int{}, []conflation.CRPullRequest{}}}
					fmt.Println("ID ", *expandedIssue.Issue.ID)
					fmt.Println("URL ", *expandedIssue.Issue.URL)
					fmt.Println("Assignee ", *expandedIssue.Issue.Assignees[0].Login)
					// Call OnlineLearn Method
				} else {
					// Call Predict Method
					expandedIssue := conflation.ExpandedIssue{Issue: conflation.CRIssue{*issuesEvent.Issue, []int{}, []conflation.CRPullRequest{}}}
					fmt.Println(expandedIssue.Issue.ID)
					fmt.Println(expandedIssue.Issue.URL)
					fmt.Println(expandedIssue.Issue.Assignees)
					//fmt.Println(issue)
					assignees := w.Models[*issuesEvent.Repo.ID].Algorithm.Predict(expandedIssue)
					// NOTE: This is likely where the assignment function will be called.
					assignment.AssignContributor(assignees[0], issuesEvent, testClient())
					// HACK: using test client
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
