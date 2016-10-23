package github-gateway

import (
  "github.com/google/go-github/github"
)


type Conflator struct {
  issues map[int]github.Issue
  pulls map[int]github.PullRequest
}


func (c *Conflator) SetPullEvents(events []github.PullRequestEvent) {

}

//TODO: Determine if we can just use SetPullEvents? Remove if that is the case
func (c *Conflator) SetIssueEvents(events []github.IssuesEvent) {

}


func (c *Conflator) SetPullRequests(pulls []github.PullRequest) {

}
func (c *Conflator) SetIssueRequests(issues []github.Issue) {

}


func (c *Conflator) Conflate() {

}



/*
type NBClassifier struct {
	classifier *bayesian.Classifier
	assignees  []bayesian.Class
}

func (c *NBClassifier) Learn(issues []issues.Issue) {
	c.assignees = distinctAssignees(issues)
	c.classifier = bayesian.NewClassifierTfIdf(c.assignees...)
	for i := 0; i < len(issues); i++ {
		c.classifier.Learn(strings.Split(issues[i].Body, " "), bayesian.Class(issues[i].Assignee))
	}
	c.classifier.ConvertTermsFreqToTfIdf()
}
*/
