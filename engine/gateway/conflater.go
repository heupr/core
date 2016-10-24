package gateway

import (
  "github.com/google/go-github/github"
)

type Conflator struct {
  Issues map[int]github.Issue
  Pulls map[int]github.PullRequest
  Options ConflationOptions
}

type ConflationOptions struct {
  OneToMany bool
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

// 1:1 Algorithm (Naive) (We may need to exclude 1:M issues)
// (Ideal? Approach 1)  We should be able to just use the closed indicator in corefx/pulls/12923
// Query https://api.github.com/repos/dotnet/corefx/pulls/12923 (pull request)
// Query https://api.github.com/repos/dotnet/corefx/issues/12886 (issue)
//
// (Alternative Approach 2) We can also use the event id
// (step 1)https://api.github.com/repos/dotnet/corefx/issues/12886/events
//                "id": 832840421,
//                "url": "https://api.github.com/repos/dotnet/corefx/issues/events/832840421",
//                "actor": {
//                "login": "stephentoub",
// (optional step 2) https://api.github.com/repos/dotnet/corefx/issues/events/832840421
// Next steps: Implement Approach 1
// Test Approach 1 with unit Test
func (c *Conflator) otoAlgorithm() {

}


// 1:M Algorithm (Optimized) (We will also use this for 1:1)
// Establish reliable relation between Github Issues and pull requests using the ("reference"? event)
// Break each relation out into a seperate issue (1 checkbox/pull request)
func (c *Conflator) otmAlgorithm() {

}

func (c *Conflator) Conflate() {
  if c.Options.OneToMany == true {
    otmAlgorithm()
  } else {
    otoAlgorithm()
  }
}

// Just a Thought
// Handle or Ignore? pull requests that were intentionally not merged (Tossing graph: should that be a thumbs down for a developer in our model? )
//       "merged": false,
//       "mergeable": false,
// example: https://api.github.com/repos/dotnet/corefx/pulls/9460
// Next steps: Ensure merged flag is populated
// Possibly create a "thin ice" list for our model's tossing graph
