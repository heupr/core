package grades

import (
	. "coralreefci/engine/gateway"
	. "coralreefci/engine/gateway/conflation"
	"encoding/csv"
	"fmt"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"os"
	"strconv"
)

type HistoricalData struct {
	Issues map[int]github.Issue
}

func (d *HistoricalData) Download() bool {
	if _, err := os.Stat("./trainingset_corefx"); err == nil {
		return false
	} else {
		d.getIssues()
		d.write()
		return true
	}
}

func (d *HistoricalData) getIssues() {

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: "23fc398670a80700b19b1ae1587825a16aa8ce57"})
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)
	gateway := Gateway{Client: client}

	pullRequests, _ := gateway.GetPullRequests()
	issues, _ := gateway.GetIssues()

	context := &Context{}
	conflator := Conflator{Algorithm: &OneToOne{Context: context}, Context: context}
	conflator.Context.SubTasks = make(map[int][]SubTask)
	conflator.Context.Issues = make(map[int]github.Issue)

	issuesCopy := make([]github.Issue, len(issues))
	pullsCopy := make([]github.PullRequest, len(pullRequests))

	//Workaround
	for i := 0; i < len(issuesCopy); i++ {
		issuesCopy[i] = *issues[i]
	}

	for i := 0; i < len(pullsCopy); i++ {
		pullsCopy[i] = *pullRequests[i]
	}

	conflator.SetIssueRequests(issuesCopy)
	conflator.SetPullRequests(pullsCopy)

	conflator.Conflate()

	d.Issues = conflator.Context.Issues
}

func (d *HistoricalData) write() {
	// Create a csv file
	f, err := os.Create("./trainingset_corefx")
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()
	// Write Unmarshaled json data to CSV file
	w := csv.NewWriter(f)
	for _, issue := range d.Issues {
		//used Buffer because it is faster than concating strings
		//may want to replace record[] string with buffer
		//see http://golang-examples.tumblr.com/post/86169510884/fastest-string-contatenation
		var record []string

		//url column
		record = append(record, *issue.URL)

		//url id column
		//Format: https://github.com/dotnet/coreclr/pull/{number}
		record = append(record, strconv.Itoa(*issue.Number))

		//title column
		record = append(record, *issue.Title)

		if issue.Body != nil {
			record = append(record, *issue.Body)
		} else {
			record = append(record, "No description")
		}

		//username column (prediction value)
		// TEMP FIX: Write out Assignee when Assignee exists
		// TODO: remove check as this is a workaround
		if issue.Assignee == nil {
			record = append(record, *issue.User.Login)
		} else {
			record = append(record, *issue.Assignee.Login)
		}

		w.Write(record)
	}
	w.Flush()
}
