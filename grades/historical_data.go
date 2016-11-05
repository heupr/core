package grades

import (
  . "coralreefci/engine/gateway/conflation"
  . "coralreefci/engine/gateway"
	"github.com/google/go-github/github"
  "golang.org/x/oauth2"
  "os"
  "encoding/csv"
  "fmt"
  "strconv"
)

type HistoricalData struct {
  Issues map[int]github.Issue
}

func (d *HistoricalData) Download() bool{
  if _, err := os.Stat("./trainingset_corefx"); err == nil {
    return false
  } else {
    d.write()
    return true
  }
}

func(d *HistoricalData) getIssues() {

  ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: "56e5fcc7bec01a3f3f797d528f08a83a5e3fec74"})
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


func(d *HistoricalData) write() {
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
    record = append(record, *issue.User.Login)

    w.Write(record)
  }
  w.Flush()
}
