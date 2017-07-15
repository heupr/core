package backend

import (
	"encoding/json"

	"github.com/google/go-github/github"
)

var issueID = 0

const ISSUE_QUERY = `SELECT id, repo_id, is_pull, payload FROM github_events WHERE id > ?`

type RepoData struct {
	RepoID int
	Open   []*github.Issue
	Closed []*github.Issue
	Pulls  []*github.PullRequest
}

func (m *MemSQL) Read() (map[int]*RepoData, error) {
	results, err := m.db.Query(ISSUE_QUERY, issueID)
	if err != nil {
		return nil, err
	}
	defer results.Close()

	repodata := make(map[int]*RepoData)

	for results.Next() {
		count := new(int)
		repo_id := new(int)
		is_pull := new(bool)
		payload := new(string)
		if err := results.Scan(count, repo_id, is_pull, payload); err != nil {
			return nil, err
		}

		if _, ok := repodata[*repo_id]; !ok {
			repodata[*repo_id] = new(RepoData)
			repodata[*repo_id].RepoID = *repo_id
			repodata[*repo_id].Open = []*github.Issue{}
			repodata[*repo_id].Closed = []*github.Issue{}
			repodata[*repo_id].Pulls = []*github.PullRequest{}
		}

		if *is_pull {
			pr := github.PullRequest{}
			if err := json.Unmarshal([]byte(*payload), &pr); err != nil {
				return nil, err
			}
			repodata[*repo_id].Pulls = append(repodata[*repo_id].Pulls, &pr)
		} else {
			i := github.Issue{}
			if err := json.Unmarshal([]byte(*payload), &i); err != nil {
				return nil, err
			}
			if i.ClosedAt == nil {
				repodata[*repo_id].Open = append(repodata[*repo_id].Open, &i)
			} else {
				repodata[*repo_id].Closed = append(repodata[*repo_id].Closed, &i)
			}
		}
	}
	return repodata, nil
}
