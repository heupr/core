package conflation

import (
	"fmt"
	"github.com/google/go-github/github"
	"strconv"
	"strings"
)

type OneToMany struct {
	Context *Context
}

func extractSubTasks(issue *github.Issue) []SubTask {
	rawSubTasks := strings.Split(*issue.Body, "[")
	length := len(rawSubTasks)
	subTasks := make([]SubTask, length)
	for i := 0; i < length; i++ {
		if strings.HasPrefix(rawSubTasks[i], "x]") {
			subTaskBody := rawSubTasks[i]
			subTaskBody = strings.Replace(subTaskBody, "-", "", -1)
			subTaskBody = strings.TrimSpace(subTaskBody)
			subTasks[i] = SubTask{Body: rawSubTasks[i]}
		}
	}
	return subTasks
}

// Might be a better way to do this. Once our unit testing is robust I will play around (if needed for performance)
func (c *OneToMany) extractIssueID(pull *github.PullRequest) int64 {
	fixIdx := strings.LastIndex(*pull.Body, "Fixes")
	body := string(*pull.Body)
	body = body[fixIdx:]

	issueIdx := strings.LastIndex(body, "issues/")
	body = body[issueIdx+7:]
	digit := digitRegexp.Find([]byte(body))
	s, _ := strconv.ParseInt(string(digit), 10, 32) //TODO: add error handling and logging (decide what to do if we have an error)
	return s
}

func (c *OneToMany) Conflate() {
	fmt.Println(&c.Context.Pulls)
	fmt.Println("OneToMany")
}

// 1:M Algorithm (Optimized) (We will also use this for 1:1)
// Establish reliable relation between Github Issues and pull requests using the ("reference"? event)
// Break each relation out into a seperate issue (1 checkbox/pull request)
