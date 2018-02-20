package backend

import (
	"bytes"
	"core/utils"
	"database/sql"
	"encoding/json"
	"reflect"
	"strings"
	"time"

	"github.com/google/go-github/github"
)

var maxID = 0

type RepoData struct {
	RepoID              int64
	Open                []*github.Issue
	Closed              []*github.Issue
	Pulls               []*github.PullRequest
	AssigneeAllocations map[string]int
	EligibleAssignees   map[string]int
	Settings            HeuprConfigSettings
}

type HeuprConfigSettings struct {
	Integration            Integration
	EnableIssueAssignments bool
	EnableLabeler          bool
	Bug                    *string
	Improvement            *string
	Feature                *string
	IgnoreUsers            map[string]bool
	StartTime              time.Time
	IgnoreLabels           map[string]bool
	Email                  string
	Twitter                string
}

func (m *MemSQL) Read() (map[int64]*RepoData, error) {
	// Current state of the Issue object (equivalent to any GitHub Event)
	ISSUE_QUERY := `
    SELECT g.id, g.repo_id, g.is_pull, g.payload
    FROM github_events g
    JOIN (
        SELECT max(id) id
        FROM github_events
        WHERE id > ?
        GROUP BY repo_id, issues_id, number
    ) T
    ON T.id = g.id AND g.action IN ('opened', 'closed')
    `

	results, err := m.db.Query(ISSUE_QUERY, maxID)
	if err != nil {
		return nil, err
	}
	defer results.Close()

	repodata := make(map[int64]*RepoData)
	for results.Next() {
		id := new(int)
		repo_id := new(int64)
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
			repodata[*repo_id].RepoID = int64(*repo_id)
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
	intKeys := make([]int64, len(keys))
	for i := 0; i < len(keys); i++ {
		interfaceKeys[i] = keys[i].Interface()
		intKeys[i] = int64(keys[i].Int())
	}
	allocations, err := m.ReadAssigneeAllocations(interfaceKeys)
	if err != nil {
		utils.AppLog.Error("database read failure - ReadAssigneeAllocations()")
		return nil, err
	}
	eligibleAssignees, err := m.ReadEligibleAssignees(interfaceKeys)
	if err != nil {
		utils.AppLog.Error("database read failure - ReadEligibleAssignees()")
		return nil, err
	}
	settings, err := m.ReadHeuprConfigSettings(interfaceKeys)
	if err != nil {
		utils.AppLog.Error("database read failure - ReadHeuprConfigSettings()")
		return nil, err
	}
	for i := 0; i < len(intKeys); i++ {
		repoID := intKeys[i]
		if _, ok := allocations[repoID]; !ok {
			allocations[repoID] = make(map[string]int)
		}
		if _, ok := eligibleAssignees[repoID]; !ok {
			eligibleAssignees[repoID] = make(map[string]int)
		}
		if _, ok := settings[repoID]; !ok {
			settings[repoID] = HeuprConfigSettings{StartTime: time.Now(), IgnoreLabels: make(map[string]bool), IgnoreUsers: make(map[string]bool)}
		}
		repodata[repoID].AssigneeAllocations = allocations[repoID]
		repodata[repoID].EligibleAssignees = eligibleAssignees[repoID]
		repodata[repoID].Settings = settings[repoID]
	}
	return repodata, nil
}

func (m *MemSQL) ReadHeuprConfigSettingsByRepoID(repoID int64) (HeuprConfigSettings, error) {
	settingsMap, err := m.ReadHeuprConfigSettings([]interface{}{repoID})
	if err != nil {
		//MVP Workaround. This helps avoid inadvertently hiding the original error.
		panic(err)
	}
	if _, ok := settingsMap[repoID]; !ok {
		settings := HeuprConfigSettings{
			StartTime:    time.Now(),
			IgnoreLabels: make(map[string]bool),
			IgnoreUsers:  make(map[string]bool),
		}
		settingsMap[repoID] = settings
	}
	return settingsMap[repoID], nil
}

func (m *MemSQL) ReadHeuprConfigSettings(repos []interface{}) (map[int64]HeuprConfigSettings, error) {
	if len(repos) == 0 {
		return nil, nil
	}
	settings := make(map[int64]HeuprConfigSettings)

	integrationSettingsQuery := `
	SELECT g.repo_id, g.start_time, g.email, g.twitter
	FROM integrations_settings g
	JOIN (
		SELECT MAX(id) id
		from integrations_settings
		WHERE repo_id IN (?` + strings.Repeat(",?", len(repos)-1) + `)
	) T
	on T.id = g.id
	`
	results, err := m.db.Query(integrationSettingsQuery, repos...)
	if err != nil {
		return nil, err
	}
	defer results.Close()

	for results.Next() {
		config := HeuprConfigSettings{
			IgnoreLabels: make(map[string]bool),
			IgnoreUsers:  make(map[string]bool),
		}
		repo_id := new(int64)
		if err := results.Scan(repo_id, &config.StartTime, &config.Email, &config.Twitter); err != nil {
			return nil, err
		}
		settings[*repo_id] = config
	}
	time.Sleep(1 * time.Second)

	integrationSettingsIgnoreUsersQuery := `
	SELECT g.repo_id, lk.user
	FROM integrations_settings g
	JOIN (
		SELECT MAX(id) id
		from integrations_settings
		WHERE repo_id IN (?` + strings.Repeat(",?", len(repos)-1) + `)
	) T
	on T.id = g.id
	JOIN integrations_settings_ignoreusers_lk lk
	on lk.integrations_settings_fk = g.id
	`
	results, err = m.db.Query(integrationSettingsIgnoreUsersQuery, repos...)
	if err != nil {
		return nil, err
	}
	defer results.Close()

	for results.Next() {
		repo_id := new(int64)
		user := new(string)
		if err := results.Scan(repo_id, user); err != nil {
			return nil, err
		}
		settings[*repo_id].IgnoreUsers[*user] = true
	}

	integrationSettingsIgnoreLabelsQuery := `
	SELECT g.repo_id, lk.label
	FROM integrations_settings g
	JOIN (
		SELECT MAX(id) id
		from integrations_settings
		WHERE repo_id IN (?` + strings.Repeat(",?", len(repos)-1) + `)
	) T
	on T.id = g.id
	JOIN integrations_settings_ignorelabels_lk lk
	on lk.integrations_settings_fk = g.id
	`
	results, err = m.db.Query(integrationSettingsIgnoreLabelsQuery, repos...)
	if err != nil {
		return nil, err
	}
	defer results.Close()

	for results.Next() {
		repo_id := new(int64)
		label := new(string)
		if err := results.Scan(repo_id, label); err != nil {
			return nil, err
		}
		settings[*repo_id].IgnoreLabels[*label] = true
	}

	//We need to pull in the latest integrations_settings_labels_bif_lk key because of a 1:M relationship (due to insert ordering)
	integrationSettingsLabelsQuery := `
    SELECT settings.repo_id, bif.bug, bif.improvement, bif.feature
    FROM integrations_settings settings
    JOIN (
        SELECT MAX(id) id
        FROM integrations_settings
        WHERE repo_id IN (?` + strings.Repeat(",?", len(repos)-1) + `)
    ) t
		ON t.id = settings.id
    JOIN integrations_settings_labels_bif_lk bif
    ON bif.integrations_settings_fk = settings.id
		JOIN (
			 SELECT MAX(id) id
			 from integrations_settings_labels_bif_lk
			 group by integrations_settings_fk
		) t2
		ON t2.id = bif.id
    `

	results, err = m.db.Query(integrationSettingsLabelsQuery, repos...)
	if err != nil {
		return nil, err
	}
	defer results.Close()

	for results.Next() {
		repoID := new(int64)
		var bug sql.NullString
		var improvement sql.NullString
		var feature sql.NullString
		if err := results.Scan(repoID, &bug, &improvement, &feature); err != nil {
			return nil, err
		}
		if config, ok := settings[*repoID]; ok {
			if bug.Valid {
				config.Bug = &bug.String
			}
			if improvement.Valid {
				config.Improvement = &improvement.String
			}
			if feature.Valid {
				config.Feature = &feature.String
			}
			//TODO: TEST THIS!!!!!
			settings[*repoID] = config
		}
	}

	return settings, nil
}

func (m *MemSQL) ReadAssigneeAllocations(repos []interface{}) (map[int64]map[string]int, error) {
	if len(repos) == 0 {
		return nil, nil
	}

	// Identifies how many Issues are assigned to the contributors on a given repo
	ASSIGNEE_ALLOCATIONS_QUERY := `
    SELECT T2.repo_id, lk.assignee, COUNT(*) AS cnt
    FROM (
        SELECT g.id, g.repo_id
        FROM github_event_assignees g
        JOIN (
            SELECT MAX(id) id
            FROM github_event_assignees
            WHERE repo_id IN (?` + strings.Repeat(",?", len(repos)-1) + `)
            GROUP BY repo_id, issues_id, number
        ) T
        ON T.id = g.id AND g.is_closed = false
    ) T2
    JOIN github_event_assignees_lk lk
    ON lk.github_event_assignees_fk = T2.id AND lk.assignee IS NOT NULL
    GROUP BY T2.repo_id, lk.assignee
    `

	results, err := m.db.Query(ASSIGNEE_ALLOCATIONS_QUERY, repos...)
	if err != nil {
		return nil, err
	}
	defer results.Close()

	allocations := make(map[int64]map[string]int)
	for results.Next() {
		repo_id := new(int64)
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

func (m *MemSQL) ReadEligibleAssignees(repos []interface{}) (map[int64]map[string]int, error) {
	//TODO: Include Merged PullRequest Users.
	//TODO: Include users with a status of contributor in the repo
	//TODO: Add a whitelist
	if len(repos) == 0 {
		return nil, nil
	}

	// Finds which Contributors in a Repository have been active in the past six months.
	RECENT_ASSIGNEES_QUERY := `
    SELECT DISTINCT T3.repo_id, T3.assignee
    FROM github_events
    JOIN (
        SELECT T2.repo_id, T2.issues_id, lk.assignee
        FROM (
            SELECT g.id, g.repo_id, g.issues_id
            FROM github_event_assignees g
            JOIN (
                SELECT MAX(id) id
                FROM github_event_assignees
                WHERE is_closed = true AND repo_id IN (?` + strings.Repeat(",?", len(repos)-1) + `)
                GROUP BY repo_id, issues_id, number
            ) T
            ON T.id = g.id
        ) T2
        JOIN github_event_assignees_lk lk
        ON lk.github_event_assignees_fk = T2.id AND lk.assignee IS NOT NULL
    ) T3
    ON T3.issues_id = github_events.issues_id
    WHERE closed_at > DATE_SUB(curdate(), INTERVAL 6 MONTH)
    `

	results, err := m.db.Query(RECENT_ASSIGNEES_QUERY, repos...)
	if err != nil {
		return nil, err
	}
	defer results.Close()

	recentAssignees := make(map[int64]map[string]int)
	for results.Next() {
		repo_id := new(int64)
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
