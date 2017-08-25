package backend

import (
	"bytes"
	"encoding/json"
	"github.com/google/go-github/github"
)

var maxID = 0

const ISSUE_QUERY = `SELECT id, repo_id, is_pull, payload FROM github_events WHERE id > ?`

type RepoData struct {
	RepoID int
	Open   []*github.Issue
	Closed []*github.Issue
	Pulls  []*github.PullRequest
}

func (m *MemSQL) Read() (map[int]*RepoData, error) {
	results, err := m.db.Query(ISSUE_QUERY, maxID)
	if err != nil {
		return nil, err
	}
	defer results.Close()

	repodata := make(map[int]*RepoData)
	for results.Next() {
		id := new(int)
		repo_id := new(int)
		is_pull := new(bool)
		var payload []byte
		if err := results.Scan(id, repo_id, is_pull, &payload); err != nil {
			return nil, err
		}

		if *id > maxID {
			maxID = *id
		}
		if _, ok := repodata[*repo_id]; !ok {
			repodata[*repo_id] = new(RepoData)
			repodata[*repo_id].RepoID = *repo_id
			repodata[*repo_id].Open = []*github.Issue{}
			repodata[*repo_id].Closed = []*github.Issue{}
			repodata[*repo_id].Pulls = []*github.PullRequest{}
		}

		if *is_pull {
			var pr github.PullRequest
			decoder := json.NewDecoder(bytes.NewReader(payload))
			decoder.UseNumber()
			if err := decoder.Decode(&pr); err != nil {
				return nil, err
			}
			repodata[*repo_id].Pulls = append(repodata[*repo_id].Pulls, &pr)
		} else {
			var issue github.Issue
			decoder := json.NewDecoder(bytes.NewReader(payload))
			decoder.UseNumber()
			if err := decoder.Decode(&issue); err != nil {
				return nil, err
			}
			if issue.ClosedAt == nil {
				repodata[*repo_id].Open = append(repodata[*repo_id].Open, &issue)
			} else {
				repodata[*repo_id].Closed = append(repodata[*repo_id].Closed, &issue)
			}
		}
	}
	return repodata, nil
}
