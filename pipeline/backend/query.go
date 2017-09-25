package backend

import (
	"bytes"
	"core/utils"
	"encoding/json"
	"github.com/google/go-github/github"
	"reflect"
	"strings"
)

var maxID = 0

//const ISSUE_QUERY = `SELECT id, repo_id, is_pull, payload FROM github_events WHERE id > ?`
const ISSUE_QUERY = `select g.id, g.repo_id, g.is_pull, g.payload from github_events g
join (
SELECT max(id) id
FROM github_events
WHERE id > ?
group by repo_id, issues_id, number
) T
on T.id = g.id ` //MVP Workaround.

type RepoData struct {
	RepoID              int
	Open                []*github.Issue
	Closed              []*github.Issue
	Pulls               []*github.PullRequest
	AssigneeAllocations map[string]int
	EligibleAssignees   map[string]int //Assignees Active in the Past Year + Whitelist
}

func (m *MemSQL) ReadAssigneeAllocations(repos []interface{}) (map[int]map[string]int, error) {
	if len(repos) == 0 {
		return nil, nil
	}
	/*
		select T2.repo_id, lk.assignee, count(*) as cnt
		from (
			select g.id, g.repo_id from github_event_assignees g
			join (
				SELECT max(id) id
				FROM github_event_assignees
		    where repo_id in (?` + strings.Repeat(",?", len(repos)-1) + ") " +
			 `group by repo_id, issues_id, number
			 ) T on T.id = g.id
			 where g.is_closed = false
		) T2
		join github_event_assignees_lk lk on lk.github_event_assignees_fk = T2.id and lk.assignee is not null
		` */
	ASSIGNEE_ALLOCATIONS_QUERY := "select T2.repo_id, lk.assignee, count(*) as cnt from ( select g.id, g.repo_id from github_event_assignees g join (SELECT max(id) id FROM github_event_assignees where repo_id in (?" + strings.Repeat(",?", len(repos)-1) + ") " + "group by repo_id, issues_id, number) T on T.id = g.id where g.is_closed = false) T2 join github_event_assignees_lk lk on lk.github_event_assignees_fk = T2.id and lk.assignee is not null"

	results, err := m.db.Query(ASSIGNEE_ALLOCATIONS_QUERY, repos...)
	if err != nil {
		return nil, err
	}
	defer results.Close()

	allocations := make(map[int]map[string]int)
	for results.Next() {
		repo_id := new(int)
		assignee := new(string)
		count := new(int)
		if err := results.Scan(repo_id, assignee, count); err != nil {
			return nil, err
		}
		if _, ok := allocations[*repo_id]; !ok {
			allocations[*repo_id] = make(map[string]int)
		}
		repoAllocations := allocations[*repo_id]
		repoAllocations[*assignee] = *count
	}
	return allocations, nil
}

func (m *MemSQL) ReadEligibleAssignees(repos []interface{}) (map[int]map[string]int, error) {
	//TODO: Include Merged PullRequest Users.
	//TODO: Include users with a status of contributor in the repo
	//TODO: Add a whitelist
	if len(repos) == 0 {
		return nil, nil
	}
	/*
	   select distinct T3.repo_id, T3.assignee from github_events
	   JOIN (
	   select T2.repo_id, T2.issues_id, lk.assignee
	   from (
	   	select g.id, g.repo_id, g.issues_id from github_event_assignees g
	   	join (
	   		SELECT max(id) id
	   		FROM github_event_assignees
	   		where is_closed = true` + " and repo_id in (?" + strings.Repeat(",?", len(repos)-1) + ")" + `
	   		 group by repo_id, issues_id, number
	   	) T
	   	on T.id = g.id
	   ) T2
	   join github_event_assignees_lk lk on lk.github_event_assignees_fk = T2.id and lk.assignee is not null
	   ) T3
	   on T3.issues_id = github_events.issues_id
	   where closed_at > DATE_SUB(curdate(), INTERVAL 1 YEAR)
	*/
	RECENT_ASSIGNEES_QUERY := "select distinct T3.repo_id, T3.assignee from github_events JOIN (select T2.repo_id, T2.issues_id, lk.assignee from (select g.id, g.repo_id, g.issues_id from github_event_assignees g join (SELECT max(id) id FROM github_event_assignees where is_closed = true" + " and repo_id in (?" + strings.Repeat(",?", len(repos)-1) + ")" + "group by repo_id, issues_id, number) T on T.id = g.id) T2 join github_event_assignees_lk lk on lk.github_event_assignees_fk = T2.id and lk.assignee is not null) T3 on T3.issues_id = github_events.issues_id where closed_at > DATE_SUB(curdate(), INTERVAL 1 YEAR)"

	results, err := m.db.Query(RECENT_ASSIGNEES_QUERY, repos...)
	if err != nil {
		return nil, err
	}
	defer results.Close()

	recentAssignees := make(map[int]map[string]int)
	for results.Next() {
		repo_id := new(int)
		assignee := new(string)
		if err := results.Scan(repo_id, assignee); err != nil {
			return nil, err
		}
		if _, ok := recentAssignees[*repo_id]; !ok {
			recentAssignees[*repo_id] = make(map[string]int)
		}
		repoAssignees := recentAssignees[*repo_id]
		repoAssignees[*assignee] = 10
	}
	return recentAssignees, nil
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
	keys := reflect.ValueOf(repodata).MapKeys()
	interfaceKeys := make([]interface{}, len(keys))
	intKeys := make([]int, len(keys))
	for i := 0; i < len(keys); i++ {
		interfaceKeys[i] = keys[i].Interface()
		intKeys[i] = int(keys[i].Int())
	}
	allocations, err := m.ReadAssigneeAllocations(interfaceKeys)
	if err != nil {
		utils.AppLog.Error("Database read failure. ReadAssigneeAllocations()")
		return nil, err
	}
	eligibleAssignees, err := m.ReadEligibleAssignees(interfaceKeys)
	if err != nil {
		utils.AppLog.Error("Database read failure. ReadEligibleAssignees()")
		return nil, err
	}
	for i := 0; i < len(intKeys); i++ {
		repoID := intKeys[i]
		repodata[repoID].AssigneeAllocations = allocations[repoID]
		repodata[repoID].EligibleAssignees = eligibleAssignees[repoID]
	}
	return repodata, nil
}
