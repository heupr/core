package retriever

import (
	"fmt"

	"github.com/google/go-github/github"

	"coralreefci/engine/gateway/conflation"
	"coralreefci/models"
)

type Worker struct {
	ID     int
	Work   chan github.Issue
	Queue  chan chan github.Issue
	Models map[int]models.Model
	Quit   chan bool
}

func NewWorker(id int, queue chan chan github.Issue) Worker {
	return Worker{
		ID:     id,
		Work:   make(chan github.Issue),
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
			case issue := <-w.Work:
				if issue.ClosedAt != nil {
					// TODO: Decide how we want to handle PR's
					// TODO: Call conflate (only on closed issues)
					// TODO: Remove Printlns
					// TODO: Call OnlineLearn Method
					expandedIssue := conflation.ExpandedIssue{
						Issue: conflation.CRIssue{
							issue, []int{}, []conflation.CRPullRequest{},
						},
					}
					fmt.Println(expandedIssue)
				} else {
					// TODO: Call Predict Method
					// assignees := w.Models[*issuesEvent.Repo.ID].Algorithm.Predict(expandedIssue)
					// NOTE: This is likely where the assignment function will be called.
					// assignment.AssignContributor(assignees[0], issuesEvent, testClient())
					// HACK: using test client
					expandedIssue := conflation.ExpandedIssue{
						Issue: conflation.CRIssue{
							issue, []int{}, []conflation.CRPullRequest{},
						},
					}
					fmt.Println(expandedIssue)
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
