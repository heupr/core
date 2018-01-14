package jira

import (
  "encoding/csv"
  "log"
  "os"

  "github.com/google/go-github/github"
  "github.com/andygrunwald/go-jira"
)

type Gateway struct {
  Client *jira.Client
}

func (g *Gateway) getCorrectedLables(path string) (map[string]string) {
  corrections := map[string]string{}
  file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	r := csv.NewReader(file)
  records, err := r.ReadAll()
  if err != nil {
    log.Fatal(err)
  }
  // Skip Header at Index 0
  for i := 1; i < len(records); i++ {
    corrections[records[i][0]] = records[i][1]
  }
  return corrections
}

func (g *Gateway) GetIssues(repo, correctedFile string) ([]*github.Issue, error) {
  corrections := g.getCorrectedLables(correctedFile)
  output := []*github.Issue{}
  opt := &jira.SearchOptions{StartAt: 1, MaxResults: 50, Expand: "foo"}
  err := g.Client.Issue.SearchPages("created >= 2005-08-05 AND created <= 2013-12-09 AND Project = " + repo, opt, func(issue jira.Issue) error {
    if val, ok := corrections[issue.Key]; ok {
        output = append(output, &github.Issue{Title: &issue.Fields.Summary, Body: &issue.Fields.Description, Labels: []github.Label{github.Label{Name: &val}}})
    }
		return nil
	})
  if err != nil {
    return nil, err
  }
  return output, nil
}
