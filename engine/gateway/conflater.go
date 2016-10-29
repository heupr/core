package gateway

import (
  "github.com/google/go-github/github"
  "fmt"
  "strings"
  "strconv"
  "regexp"
)

var digitRegexp = regexp.MustCompile("[0-9]+")

type SubTask struct {
  Assignee string
  Body string
}

type Conflator struct {
  Issues map[int]github.Issue
  SubTasks map[int][]SubTask
  Pulls []github.PullRequest
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
  c.Pulls = pulls
}

func (c *Conflator) SetIssueRequests(issues []github.Issue) {
  for i := 0; i < len(issues); i++ {
    issueNumber := *issues[i].Number
    c.Issues[issueNumber] = issues[i]
    if c.Options.OneToMany == true {
      c.SubTasks[issueNumber] = extractSubTasks(&issues[i])
    }
  }
}

// Might be a better way to do this. Once our unit testing is robust I will play around (if needed for performance)
func extractIssueId(pull *github.PullRequest) int64 {
	fixIdx := strings.LastIndex(*pull.Body, "Fixes" )
  body := string(*pull.Body)
  body = body[fixIdx:]

  issueIdx := strings.LastIndex(body, "issues/")
  body = body[issueIdx+7:]
  digit := digitRegexp.Find([]byte(body))
  s, _ := strconv.ParseInt(string(digit), 10, 32) //TODO: add error handling and logging (decide what to do if we have an error)
  return s
}


//TODO: flesh this out
func (c *Conflator) linkPullRequestsToIssues() {
  fmt.Println(extractIssueId(&c.Pulls[0]))
}

func extractSubTasks(issue *github.Issue) []SubTask {
  rawSubTasks := strings.Split(*issue.Body, "[")
  length := len(rawSubTasks)
  subTasks := make([]SubTask, length)
  for i := 0; i < length; i++ {
    if strings.HasPrefix(rawSubTasks[i], "x]") {
      fmt.Println(rawSubTasks[i])
      subTasks[i] = SubTask{ Body: rawSubTasks[i]}
    }
  }
  return subTasks
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
  c.linkPullRequestsToIssues()
}

func (c *Conflator) Conflate() {
  if c.Options.OneToMany == true {
    c.otmAlgorithm()
  } else {
    c.otoAlgorithm()
  }
}

// Just a Thought
// Handle or Ignore? pull requests that were intentionally not merged (Tossing graph: should that be a thumbs down for a developer in our model? )
//       "merged": false,
//       "mergeable": false,
// example: https://api.github.com/repos/dotnet/corefx/pulls/9460
// Next steps: Ensure merged flag is populated
// Possibly create a "thin ice" list for our model's tossing graph
