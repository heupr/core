package ingestor

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/google/go-github/github"
	"go.uber.org/zap"

	"core/utils"
)

type Event struct {
	Type    string            `json:"type"`
	Repo    github.Repository `json:"repo"`
	Action  string            `json:"action"`
	Payload interface{}       `json:"payload"`
}

// type Value interface{}

type Integration struct {
	RepoID         int64
	AppID          int
	InstallationID int
}

type EventType int

const (
	PullRequest EventType = iota
	Issue
	All
)

type EventQuery struct {
	Type EventType
	Repo string
}

type DataAccess interface {
	open()
	Close()
	continuityCheck(query string) ([][]interface{}, error)
	restartCheck(query string, repoID int64) (int, int, error)
	ReadIntegrations() ([]Integration, error)
	ReadIntegrationByRepoID(repoID int64) (*Integration, error)
	InsertIssue(issue github.Issue, action *string)
	InsertPullRequest(pull github.PullRequest, action *string)
	BulkInsertIssuesPullRequests(issues []*github.Issue, pulls []*github.PullRequest)
	InsertRepositoryIntegration(repoID int64, appID int, installationID int64)
	InsertRepositoryIntegrationSettings(settings HeuprConfigSettings)
	InsertGobLabelSettings(settings storage) error
	DeleteRepositoryIntegration(repoID int64, appID int, installationID int64)
	ObliterateIntegration(appID int, installationID int64)
}

type Database struct {
	db         *sql.DB
	BufferPool Pool
}

func (d *Database) open() {
	mysql, err := sql.Open("mysql", "root@/heupr?interpolateParams=true")
	if err != nil {
		// TODO: Implement proper error handling (not just panic).
		panic(err.Error())
	}
	d.db = mysql
}

func (d *Database) Close() {
	d.db.Close()
}

func (d *Database) continuityCheck(query string) ([][]interface{}, error) {
	results, err := d.db.Query(query)
	if err != nil {
		utils.AppLog.Error("continuity check query", zap.Error(err))
		return nil, err
	}
	defer results.Close()

	output := [][]interface{}{}
	for results.Next() {
		repoID := new(int64)
		startNum, endNum := new(int64), new(int64)
		isPull := new(bool)
		if err := results.Scan(repoID, startNum, endNum, isPull); err != nil {
			utils.AppLog.Error("continuity check row scan", zap.Error(err))
			return nil, err
		}
		output = append(output, []interface{}{repoID, startNum, endNum, isPull})
	}
	if err = results.Err(); err != nil {
		utils.AppLog.Error("continuity check next loop", zap.Error(err))
		return nil, err
	}
	return output, nil
}

func (d *Database) restartCheck(query string, repoID int64) (int, int, error) {
	results, err := d.db.Query(query, repoID, repoID)
	if err != nil {
		utils.AppLog.Error("restart check query", zap.Error(err))
		return 0, 0, err
	}

	issueNum := new(int)
	pullNum := new(int)
	for results.Next() {
		number := new(int)
		isPull := new(bool)
		if err := results.Scan(number, isPull); err != nil {
			utils.AppLog.Error("restart check scan", zap.Error(err))
			return 0, 0, err
		}
		switch *isPull {
		case false:
			*issueNum = *number
		case true:
			*pullNum = *number
		}
	}
	if err = results.Err(); err != nil {
		utils.AppLog.Error("restart check next loop", zap.Error(err))
		return 0, 0, err
	}
	return *issueNum, *pullNum, nil
}

func (d *Database) FlushBackTestTable() {
	d.db.Exec("optimize table backtest_events flush")
}

func (d *Database) EnableRepo(repoID int) {
	var buffer bytes.Buffer
	archRepoInsert := "INSERT INTO arch_repos(repository_id, enabled) VALUES"
	valuesFmt := "(?,?)"

	buffer.WriteString(archRepoInsert)
	buffer.WriteString(valuesFmt)
	result, err := d.db.Exec(buffer.String(), repoID, true)
	if err != nil {
		utils.AppLog.Error("database repo insert failure", zap.Error(err))
	} else {
		rows, _ := result.RowsAffected()
		utils.AppLog.Info("database repo insert success", zap.Int64("rows", rows))
	}
}

func (d *Database) InsertGobLabelSettings(settings storage) error {
	var settingsID int64
	utils.AppLog.Info("SELECT max(id) id FROM integrations_settings WHERE repo_id =", zap.Int64("RepoID", settings.RepoID))
	err := d.db.QueryRow("SELECT max(id) id FROM integrations_settings WHERE repo_id = ?", settings.RepoID).Scan(&settingsID)
	if err != nil {
		if err == sql.ErrNoRows {
			//TODO: Automatically handle this case
			utils.AppLog.Error("0 records found. republish signup webook", zap.Error(err))
		}
		utils.AppLog.Error("database read failure - integrations_settings", zap.Error(err))
		return err
	}

	var buffer bytes.Buffer
	settingsInsert := "INSERT INTO integrations_settings_labels_bif_lk(integrations_settings_fk, repo_id, bug, feature, improvement) VALUES"
	valuesFmt := "(?,?,?,?,?)"

	buffer.WriteString(settingsInsert)
	buffer.WriteString(valuesFmt)

	var bug *string
	if bugBucket, ok := settings.Buckets["typebug"]; ok {
		for i := 0; i < len(bugBucket); i++ {
			if bugBucket[i].Selected {
				bug = &bugBucket[i].Name
			}
		}
	}
	var feature *string
	if featureBucket, ok := settings.Buckets["typefeature"]; ok {
		for i := 0; i < len(featureBucket); i++ {
			if featureBucket[i].Selected {
				feature = &featureBucket[i].Name
			}
		}
	}
	var improvement *string
	if improvementBucket, ok := settings.Buckets["typeimprovement"]; ok {
		for i := 0; i < len(improvementBucket); i++ {
			if improvementBucket[i].Selected {
				improvement = &improvementBucket[i].Name
			}
		}
	}
	result, err := d.db.Exec(buffer.String(), settingsID, settings.RepoID, bug, feature, improvement)
	if err != nil {
		utils.AppLog.Error("Database Insert Failure", zap.Error(err))
		return err
	} else {
		rows, _ := result.RowsAffected()
		utils.AppLog.Debug("Database Insert Success", zap.Int64("Rows", rows))
	}
	buffer.Reset()

	return nil
}

func (d *Database) InsertRepositoryIntegrationSettings(settings HeuprConfigSettings) {
	var settingsID int64
	var buffer bytes.Buffer
	settingsInsert := "INSERT INTO integrations_settings(repo_id, start_time, email, twitter) VALUES"
	valuesFmt := "(?,?,?,?)"

	buffer.WriteString(settingsInsert)
	buffer.WriteString(valuesFmt)
	result, err := d.db.Exec(buffer.String(), settings.Integration.RepoID, settings.StartTime, settings.Email, settings.Twitter)
	if err != nil {
		utils.AppLog.Error("Database Insert Failure", zap.Error(err))
		return
	} else {
		rows, _ := result.RowsAffected()
		settingsID, _ = result.LastInsertId()
		utils.AppLog.Info("Database Insert Success", zap.Int64("Rows", rows))
	}
	buffer.Reset()

	if settings.IgnoreUsers != nil && len(settings.IgnoreUsers) > 0 {
		settingsIgnoreUsersLookupInsert := "INSERT INTO integrations_settings_ignoreusers_lk(integrations_settings_fk, user) VALUES"
		settingsIgnoreUsersLookupValuesFmt := "(?,?)"
		settingsIgnoreUsersLookupNumValues := 2 * len(settings.IgnoreUsers)
		buffer.WriteString(settingsIgnoreUsersLookupInsert)
		values := make([]interface{}, settingsIgnoreUsersLookupNumValues)
		delimeter := ""
		for i := 0; i < len(settings.IgnoreUsers); i++ {
			buffer.WriteString(delimeter)
			buffer.WriteString(settingsIgnoreUsersLookupValuesFmt)
			values[i+i+0] = settingsID
			values[i+i+1] = settings.IgnoreUsers[i]
			delimeter = ","
		}
		result, err = d.db.Exec(buffer.String(), values...)
		if err != nil {
			utils.AppLog.Error("Database Insert Failure", zap.Error(err))
		} else {
			rows, _ := result.RowsAffected()
			utils.AppLog.Info("Database Insert Success", zap.Int64("Rows", rows))
		}
	}
	buffer.Reset()

	if settings.IgnoreLabels != nil && len(settings.IgnoreLabels) > 0 {
		settingsIgnoreLabelsLookupInsert := "INSERT INTO integrations_settings_ignorelabels_lk(integrations_settings_fk, label) VALUES"
		settingsIgnoreLabelsValuesFmt := "(?,?)"
		settingsIgnoreLabelsNumValues := 2 * len(settings.IgnoreLabels)
		buffer.WriteString(settingsIgnoreLabelsLookupInsert)
		values := make([]interface{}, settingsIgnoreLabelsNumValues)
		delimeter := ""
		for i := 0; i < len(settings.IgnoreLabels); i++ {
			buffer.WriteString(delimeter)
			buffer.WriteString(settingsIgnoreLabelsValuesFmt)
			values[i+i+0] = settingsID
			values[i+i+1] = settings.IgnoreLabels[i]
			delimeter = ","
		}
		result, err = d.db.Exec(buffer.String(), values...)
		if err != nil {
			utils.AppLog.Error("Database Insert Failure", zap.Error(err))
		} else {
			rows, _ := result.RowsAffected()
			utils.AppLog.Info("Database Insert Success", zap.Int64("Rows", rows))
		}
	}
	buffer.Reset()
}

func (d *Database) InsertRepositoryIntegration(repoID int64, appID int, installationID int64) {
	var buffer bytes.Buffer
	integrationsInsert := "INSERT INTO integrations(repo_id, app_id, installation_id) VALUES"
	valuesFmt := "(?,?,?)"

	buffer.WriteString(integrationsInsert)
	buffer.WriteString(valuesFmt)
	result, err := d.db.Exec(buffer.String(), repoID, appID, installationID)
	if err != nil {
		utils.AppLog.Error("database integration insert failure", zap.Error(err))
	} else {
		rows, _ := result.RowsAffected()
		utils.AppLog.Info("database integration insert success", zap.Int64("rows", rows))
	}
}

func (d *Database) DeleteRepositoryIntegration(repoID int64, appID int, installationID int64) {
	result, err := d.db.Exec("DELETE FROM integrations where repo_id = ? and app_id = ? and installation_id = ?", repoID, appID, installationID)
	if err != nil {
		utils.AppLog.Error("database integration delete failure", zap.Error(err))
	} else {
		rows, _ := result.RowsAffected()
		utils.AppLog.Info("database integration delete success", zap.Int64("rows", rows))
	}
}

func (d *Database) ObliterateIntegration(appID int, installationID int64) {
	result, err := d.db.Exec("DELETE FROM integrations where app_id = ? and installation_id = ?", appID, installationID)
	if err != nil {
		utils.AppLog.Error("database integration obliterate failure", zap.Error(err))
	} else {
		rows, _ := result.RowsAffected()
		utils.AppLog.Info("database integration obliterate success", zap.Int64("rows", rows))
	}
}

func (d *Database) ReadIntegrations() ([]Integration, error) {
	integrations := []Integration{}
	results, err := d.db.Query("SELECT repo_id, app_id, installation_id FROM integrations")
	if err != nil {
		return nil, err
	}

	defer results.Close()
	for results.Next() {
		integration := Integration{}
		err := results.Scan(&integration.RepoID, &integration.AppID, &integration.InstallationID)
		if err != nil {
			return nil, err
		}
		integrations = append(integrations, integration)
		err = results.Err()
		if err != nil {
			return nil, err
		}
	}
	return integrations, nil
}

func (d *Database) ReadIntegrationByRepoID(repoID int64) (*Integration, error) {
	integration := new(Integration)
	err := d.db.QueryRow("SELECT repo_id, app_id, installation_id FROM integrations WHERE repo_id = ?", repoID).Scan(&integration.RepoID, &integration.AppID, &integration.InstallationID)
	if err != nil {
		if err != sql.ErrNoRows {
			utils.AppLog.Error("database read failure - ReadIntegrationByRepoID()", zap.Error(err))
		}
		return nil, err
	}
	return integration, nil
}

func (d *Database) BulkInsertBacktestEvents(events []*Event) {
	buffer := d.BufferPool.Get()
	for i := 0; i < len(events); i++ {
		buffer.AppendInt(int64(*events[i].Repo.ID))
		buffer.AppendByte('~')
		buffer.AppendString(*events[i].Repo.Name)
		buffer.AppendByte('~')
		if events[i].Action == "closed" {
			buffer.AppendInt(1)
		} else {
			buffer.AppendInt(0)
		}
		buffer.AppendByte('~')
		if events[i].Type == "PullRequestEvent" {
			buffer.AppendInt(1)
		} else {
			buffer.AppendInt(0)
		}
		buffer.AppendByte('~')
		payload, _ := json.Marshal(events[i])
		_, _ = buffer.Write(escapeBytesBackslash(stripCtlAndExtFromBytes(payload)))
		buffer.AppendByte('\n')
	}

	sqlBuffer := bytes.NewBuffer(buffer.Bytes())
	buffer.Reset()
	buffer.Free()

	mysql.RegisterReaderHandler("data", func() io.Reader {
		return sqlBuffer
	})
	defer mysql.DeregisterReaderHandler("data")
	result, err := d.db.Exec("LOAD DATA LOCAL INFILE 'Reader::data' INTO TABLE backtest_events FIELDS TERMINATED BY '~' LINES TERMINATED BY '\n' (repo_id,repo_name,is_closed,is_pull,payload)")
	if err != nil {
		utils.AppLog.Error("Database Insert Failure", zap.Error(err))
	} else {
		rows, _ := result.RowsAffected()
		utils.AppLog.Info("database insert success", zap.Int64("rows", rows))
	}
	sqlBuffer.Reset()
}

func (d *Database) ReadBacktestEvents(params EventQuery) ([]Event, error) {
	events := []Event{}
	var payload []byte
	var results *sql.Rows
	var err error
	switch t := params.Type; t {
	case PullRequest:
		results, err = d.db.Query("select payload from backtest_events where repo_name=? and is_pull=? and is_closed=?", params.Repo, 1, 1)
	case Issue:
		results, err = d.db.Query("select payload from backtest_events where repo_name=? and is_pull=? and is_closed=?", params.Repo, 0, 1)
	default:
		results, err = d.db.Query("select payload from backtest_events where repo_name=? and is_closed=?", params.Repo, 1)
	}
	if err != nil {
		return nil, err
	}
	defer results.Close()
	for results.Next() {
		var event Event
		err := results.Scan(&payload)
		if err != nil {
			return nil, err
		}
		decoder := json.NewDecoder(bytes.NewReader(payload))
		decoder.UseNumber()
		if err := decoder.Decode(&event); err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	err = results.Err()
	if err != nil {
		return nil, err
	}
	return events, nil
}

func (d *Database) ReadBacktestRepos() ([]github.Repository, error) {
	repos := []github.Repository{}

	results, err := d.db.Query(`select T.repo_name, T.repo_id
	from
	(
		select count(*) cnt, repo_name, repo_id from backtest_events where is_pull = 0 and is_closed = 1 and repo_name != 'chrsmith/google-api-java-client'
		group by repo_name
	) T
	order by T.cnt desc LIMIT 10
    `)
	if err != nil {
		return nil, err
	}
	defer results.Close()

	for results.Next() {
		repo_name := new(string)
		repo_id := new(int64)
		if err := results.Scan(repo_name, repo_id); err != nil {
			return nil, err
		}
		r := strings.Split(*repo_name, "/")
		repos = append(repos, github.Repository{ID: repo_id, Name: github.String(r[1]), Organization: &github.Organization{Name: github.String(r[0])}})
	}

	err = results.Err()
	if err != nil {
		return nil, err
	}
	return repos, nil
}

func (d *Database) ReadIssuesTest() ([]github.Issue, error) {
	events := []github.Issue{}
	var payload []byte
	var results *sql.Rows
	var err error
	results, err = d.db.Query("select payload from github_events where is_pull=0")
	if err != nil {
		return nil, err
	}
	defer results.Close()
	for results.Next() {
		var event github.Issue
		err := results.Scan(&payload)
		if err != nil {
			return nil, err
		}
		decoder := json.NewDecoder(bytes.NewReader(payload))
		decoder.UseNumber()
		if err := decoder.Decode(&event); err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	err = results.Err()
	if err != nil {
		return nil, err
	}
	return events, nil
}

func (d *Database) ReadPullRequestTest() ([]github.PullRequest, error) {
	events := []github.PullRequest{}
	var payload []byte
	var results *sql.Rows
	var err error
	results, err = d.db.Query("select payload from github_events where is_pull=1")
	if err != nil {
		return nil, err
	}
	defer results.Close()
	for results.Next() {
		var event github.PullRequest
		err := results.Scan(&payload)
		if err != nil {
			return nil, err
		}
		decoder := json.NewDecoder(bytes.NewReader(payload))
		decoder.UseNumber()
		if err := decoder.Decode(&event); err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	err = results.Err()
	if err != nil {
		return nil, err
	}
	return events, nil
}

func (d *Database) InsertIssue(issue github.Issue, action *string) {
	d.LogIssueAssignees(issue)

	var buffer bytes.Buffer
	eventsInsert := "INSERT INTO github_events(repo_id,issues_id,number,action,payload,is_pull,is_closed,closed_at) VALUES"
	eventsValuesFmt := "(?,?,?,?,?,0,?,?)"
	numValues := 7

	buffer.WriteString(eventsInsert)
	buffer.WriteString(eventsValuesFmt)
	values := make([]interface{}, numValues)
	values[0] = *issue.Repository.ID
	values[1] = issue.ID
	values[2] = issue.Number
	values[3] = action
	payload, _ := json.Marshal(issue)
	values[4] = stripCtlAndExtFromBytes(payload)
	if issue.ClosedAt == nil {
		values[5] = false
	} else {
		values[5] = true
	}
	values[6] = issue.ClosedAt
	result, err := d.db.Exec(buffer.String(), values...)
	if err != nil {
		utils.AppLog.Error("Database Insert Failure", zap.Error(err))
	} else {
		rows, _ := result.RowsAffected()
		utils.AppLog.Debug("Database Insert Success", zap.Int64("Rows", rows))
	}
}

func (d *Database) BulkInsertIssuesPullRequests(issues []*github.Issue, pulls []*github.PullRequest) {
	buffer := d.BufferPool.Get()

	for i := 0; i < len(issues); i++ {
		d.LogIssueAssignees(*issues[i])

		buffer.AppendInt(int64(*issues[i].Repository.ID))
		buffer.AppendByte('~')
		buffer.AppendInt(int64(*issues[i].ID))
		buffer.AppendByte('~')
		buffer.AppendInt(int64(*issues[i].Number))
		buffer.AppendByte('~')
		if issues[i].ClosedAt == nil {
			buffer.AppendString("opened")
			buffer.AppendByte('~')
		} else {
			buffer.AppendString("closed")
			buffer.AppendByte('~')
		}
		payload, _ := json.Marshal(*issues[i])
		_, _ = buffer.Write(escapeBytesBackslash(stripCtlAndExtFromBytes(payload)))
		buffer.AppendByte('~')
		buffer.AppendInt(0)
		buffer.AppendByte('~')
		if issues[i].ClosedAt == nil {
			buffer.AppendInt(0)
			buffer.AppendByte('~')
		} else {
			buffer.AppendInt(1)
			buffer.AppendByte('~')
			buffer.Write([]byte(issues[i].ClosedAt.Format(time.RFC3339Nano)))
		}
		buffer.AppendByte('\n')
	}

	for i := 0; i < len(pulls); i++ {
		if pulls[i].Merged != nil && *pulls[i].Merged == true {
			d.LogMergedPullRequestAssignees(*pulls[i])
		}

		buffer.AppendInt(int64(*pulls[i].Base.Repo.ID))
		buffer.AppendByte('~')
		buffer.AppendInt(int64(*pulls[i].ID))
		buffer.AppendByte('~')
		buffer.AppendInt(int64(*pulls[i].Number))
		buffer.AppendByte('~')
		if pulls[i].ClosedAt == nil {
			buffer.AppendString("opened")
			buffer.AppendByte('~')
		} else {
			buffer.AppendString("closed")
			buffer.AppendByte('~')
		}
		payload, _ := json.Marshal(*pulls[i])
		_, _ = buffer.Write(escapeBytesBackslash(stripCtlAndExtFromBytes(payload)))
		buffer.AppendByte('~')
		buffer.AppendInt(1)
		buffer.AppendByte('~')
		if pulls[i].ClosedAt == nil {
			buffer.AppendInt(0)
			buffer.AppendByte('~')
		} else {
			buffer.AppendInt(1)
			buffer.AppendByte('~')
			buffer.Write([]byte(pulls[i].ClosedAt.Format(time.RFC3339Nano)))
		}
		buffer.AppendByte('\n')
	}

	issues = nil //PERF: Mark for garbage collection
	pulls = nil  //PERF: Mark for garbage collection
	sqlBuffer := bytes.NewBuffer(buffer.Bytes())
	buffer.Reset()
	buffer.Free()
	mysql.RegisterReaderHandler("data", func() io.Reader {
		return sqlBuffer
	})
	defer mysql.DeregisterReaderHandler("data")
	result, err := d.db.Exec("LOAD DATA LOCAL INFILE 'Reader::data' INTO TABLE github_events FIELDS TERMINATED BY '~' LINES TERMINATED BY '\n' (repo_id,issues_id,number,action,payload,is_pull,is_closed,closed_at)")
	if err != nil {
		utils.AppLog.Error("Database Insert Failure", zap.Error(err))
	} else {
		rows, _ := result.RowsAffected()
		utils.AppLog.Info("Database Insert Success", zap.Int64("Rows", rows))
	}
}

func (d *Database) BulkInsertIssues(issues []*github.Issue) {
	buffer := d.BufferPool.Get()

	for i := 0; i < len(issues); i++ {
		d.LogIssueAssignees(*issues[i])

		buffer.AppendInt(int64(*issues[i].Repository.ID))
		buffer.AppendByte('~')
		buffer.AppendInt(int64(*issues[i].ID))
		buffer.AppendByte('~')
		buffer.AppendInt(int64(*issues[i].Number))
		buffer.AppendByte('~')
		if issues[i].ClosedAt == nil {
			buffer.AppendString("opened")
			buffer.AppendByte('~')
		} else {
			buffer.AppendString("closed")
			buffer.AppendByte('~')
		}
		payload, _ := json.Marshal(*issues[i])
		_, _ = buffer.Write(escapeBytesBackslash(stripCtlAndExtFromBytes(payload)))
		buffer.AppendByte('~')
		buffer.AppendInt(0)
		buffer.AppendByte('~')
		if issues[i].ClosedAt == nil {
			buffer.AppendInt(0)
			buffer.AppendByte('~')
		} else {
			buffer.AppendInt(1)
			buffer.AppendByte('~')
			buffer.Write([]byte(issues[i].ClosedAt.Format(time.RFC3339Nano)))
		}
		buffer.AppendByte('\n')
	}

	issues = nil //PERF: Mark for garbage collection
	sqlBuffer := bytes.NewBuffer(buffer.Bytes())
	buffer.Reset()
	buffer.Free()
	mysql.RegisterReaderHandler("data", func() io.Reader {
		return sqlBuffer
	})
	defer mysql.DeregisterReaderHandler("data")
	result, err := d.db.Exec("LOAD DATA LOCAL INFILE 'Reader::data' INTO TABLE github_events FIELDS TERMINATED BY '~' LINES TERMINATED BY '\n' (repo_id,issues_id,number,action,payload,is_pull,is_closed,closed_at)")
	if err != nil {
		utils.AppLog.Error("Database Insert Failure", zap.Error(err))
	} else {
		rows, _ := result.RowsAffected()
		utils.AppLog.Info("Database Insert Success", zap.Int64("Rows", rows))
	}
}

func (d *Database) InsertPullRequest(pull github.PullRequest, action *string) {
	if pull.Merged != nil && *pull.Merged == true {
		d.LogMergedPullRequestAssignees(pull)
	}
	var buffer bytes.Buffer
	eventsInsert := "INSERT INTO github_events(repo_id,issues_id,action,number,payload,is_pull,is_closed,closed_at) VALUES"
	eventsValuesFmt := "(?,?,?,?,?,1,?)"
	numValues := 7

	buffer.WriteString(eventsInsert)
	buffer.WriteString(eventsValuesFmt)
	values := make([]interface{}, numValues)
	values[0] = pull.Base.Repo.ID
	values[1] = pull.ID
	values[2] = pull.Number
	values[3] = action
	payload, _ := json.Marshal(pull)
	values[4] = stripCtlAndExtFromBytes(payload)
	if pull.ClosedAt == nil {
		values[5] = false
	} else {
		values[5] = true
	}
	values[6] = pull.ClosedAt
	result, err := d.db.Exec(buffer.String(), values...)
	if err != nil {
		utils.AppLog.Error("Database Insert Failure", zap.Error(err))
	} else {
		rows, _ := result.RowsAffected()
		utils.AppLog.Debug("Database Insert Success", zap.Int64("Rows", rows))
	}
}

func (d *Database) BulkInsertPullRequests(pulls []*github.PullRequest) {
	buffer := d.BufferPool.Get()

	for i := 0; i < len(pulls); i++ {
		if pulls[i].Merged != nil && *pulls[i].Merged == true {
			d.LogMergedPullRequestAssignees(*pulls[i])
		}

		buffer.AppendInt(int64(*pulls[i].Base.Repo.ID))
		buffer.AppendByte('~')
		buffer.AppendInt(int64(*pulls[i].ID))
		buffer.AppendByte('~')
		buffer.AppendInt(int64(*pulls[i].Number))
		buffer.AppendByte('~')
		if pulls[i].ClosedAt == nil {
			buffer.AppendString("opened")
			buffer.AppendByte('~')
		} else {
			buffer.AppendString("closed")
			buffer.AppendByte('~')
		}
		payload, _ := json.Marshal(*pulls[i])
		_, _ = buffer.Write(escapeBytesBackslash(stripCtlAndExtFromBytes(payload)))
		buffer.AppendByte('~')
		buffer.AppendInt(1)
		buffer.AppendByte('~')
		if pulls[i].ClosedAt == nil {
			buffer.AppendInt(0)
			buffer.AppendByte('~')
		} else {
			buffer.AppendInt(1)
			buffer.AppendByte('~')
			buffer.Write([]byte(pulls[i].ClosedAt.Format(time.RFC3339Nano)))
		}
		buffer.AppendByte('\n')
	}

	pulls = nil //PERF: Mark for garbage collection
	sqlBuffer := bytes.NewBuffer(buffer.Bytes())
	buffer.Reset()
	buffer.Free()
	mysql.RegisterReaderHandler("data", func() io.Reader {
		return sqlBuffer
	})
	defer mysql.DeregisterReaderHandler("data")
	result, err := d.db.Exec("LOAD DATA LOCAL INFILE 'Reader::data' INTO TABLE github_events FIELDS TERMINATED BY '~' LINES TERMINATED BY '\n' (repo_id,issues_id,number,action,payload,is_pull,is_closed,closed_at)")
	if err != nil {
		utils.AppLog.Error("Database Insert Failure", zap.Error(err))
	} else {
		rows, _ := result.RowsAffected()
		utils.AppLog.Info("Database Insert Success", zap.Int64("Rows", rows))
	}
}

func (d *Database) LogIssueAssignees(issue github.Issue) {
	var assigneesID int64
	var buffer bytes.Buffer
	issueAssigneesInsert := "INSERT INTO github_event_assignees(repo_id,issues_id,number,is_closed,is_pull) VALUES"
	issueAssigneesValuesFmt := "(?,?,?,?,0)"
	issueAssigneesNumValues := 4

	buffer.WriteString(issueAssigneesInsert)
	buffer.WriteString(issueAssigneesValuesFmt)
	values := make([]interface{}, issueAssigneesNumValues)
	values[0] = issue.Repository.ID
	values[1] = issue.ID
	values[2] = issue.Number
	if issue.ClosedAt == nil {
		values[3] = false
	} else {
		values[3] = true
	}

	result, err := d.db.Exec(buffer.String(), values...)
	if err != nil {
		utils.AppLog.Error("Database Insert Failure", zap.Error(err))
		return
	} else {
		rows, _ := result.RowsAffected()
		assigneesID, _ = result.LastInsertId()
		utils.AppLog.Debug("Database Insert Success", zap.Int64("Rows", rows))
	}
	buffer.Reset()

	issueAssigneesLookupInsert := "INSERT INTO github_event_assignees_lk(github_event_assignees_fk, assignee) VALUES"
	issueAssigneesLookupValuesFmt := "(?,?)"
	issueAssigneesLookupNumValues := 2
	if issue.Assignees != nil && len(issue.Assignees) > 0 {
		issueAssigneesLookupNumValues = 2 * len(issue.Assignees)
	}

	buffer.WriteString(issueAssigneesLookupInsert)
	values = make([]interface{}, issueAssigneesLookupNumValues)
	if issue.Assignees != nil && len(issue.Assignees) > 0 {
		delimeter := ""
		for i := 0; i < len(issue.Assignees); i++ {
			buffer.WriteString(delimeter)
			buffer.WriteString(issueAssigneesLookupValuesFmt)
			values[i+i+0] = assigneesID
			values[i+i+1] = issue.Assignees[i].Login
			delimeter = ","
		}
	} else {
		values[0] = assigneesID
		if issue.Assignee != nil {
			values[1] = issue.Assignee.Login
		}
		buffer.WriteString(issueAssigneesLookupValuesFmt)
	}

	result, err = d.db.Exec(buffer.String(), values...)
	if err != nil {
		utils.AppLog.Error("Database Insert Failure", zap.Error(err))
	} else {
		rows, _ := result.RowsAffected()
		utils.AppLog.Debug("Database Insert Success", zap.Int64("Rows", rows))
	}
}

func (d *Database) LogMergedPullRequestAssignees(pull github.PullRequest) {
	var assigneesID int64
	var buffer bytes.Buffer
	pullAssigneesInsert := "INSERT INTO github_event_assignees(repo_id,issues_id,number,is_closed, is_pull) VALUES"
	pullAssigneesValuesFmt := "(?,?,?,?,1)"
	pullAssigneesNumValues := 4

	buffer.WriteString(pullAssigneesInsert)
	buffer.WriteString(pullAssigneesValuesFmt)
	values := make([]interface{}, pullAssigneesNumValues)
	values[0] = pull.Base.Repo.ID
	values[1] = pull.ID
	values[2] = pull.Number
	if pull.ClosedAt == nil {
		values[3] = false
	} else {
		values[3] = true
	}

	result, err := d.db.Exec(buffer.String(), values...)
	if err != nil {
		utils.AppLog.Error("Database Insert Failure", zap.Error(err))
		return
	} else {
		rows, _ := result.RowsAffected()
		assigneesID, _ = result.LastInsertId()
		utils.AppLog.Debug("Database Insert Success", zap.Int64("Rows", rows))
	}
	buffer.Reset()

	pullAssigneesLookupInsert := "INSERT INTO github_event_assignees_lk(github_event_assignees_fk, assignee) VALUES"
	pullAssigneesLookupValuesFmt := "(?,?)"
	pullAssigneesLookupNumValues := 2

	buffer.WriteString(pullAssigneesLookupInsert)
	values = make([]interface{}, pullAssigneesLookupNumValues)
	values[0] = assigneesID
	if pull.User != nil {
		values[1] = pull.User.Login
	}
	buffer.WriteString(pullAssigneesLookupValuesFmt)

	result, err = d.db.Exec(buffer.String(), values...)
	if err != nil {
		utils.AppLog.Error("Database Insert Failure", zap.Error(err))
	} else {
		rows, _ := result.RowsAffected()
		utils.AppLog.Debug("Database Insert Success", zap.Int64("Rows", rows))
	}
}

func stripCtlAndExtFromBytes(str []byte) []byte {
	b := make([]byte, len(str))
	var bl int
	for i := 0; i < len(str); i++ {
		c := str[i]
		if c >= 32 && c < 127 {
			b[bl] = c
			bl++
		}
	}
	//return b[:bl]
	str = b[:bl] //PERF
	return str
}

func escapeString(sql string) string {
	dest := make([]byte, 0, 2*len(sql))
	var escape byte
	for i := 0; i < len(sql); i++ {
		c := sql[i]
		escape = 0
		switch c {
		case '\\':
			escape = '\\'
			break
		case '\'':
			escape = '\''
			break
		}
		if escape != 0 {
			dest = append(dest, '\\', escape)
		} else {
			dest = append(dest, c)
		}
	}
	return string(dest)
}

func escapeBytesQuotes(v []byte) []byte {
	buf := make([]byte, 2*len(v))
	pos := 0
	for _, c := range v {
		if c == '\'' {
			buf[pos] = '\''
			buf[pos+1] = '\''
			pos += 2
		} else {
			buf[pos] = c
			pos++
		}
	}
	return buf[:pos]
}

func escapeBytesBackslash(v []byte) []byte {
	buf := make([]byte, 2*len(v))
	pos := 0
	for i := 0; i < len(v); i++ {
		switch v[i] {
		case '\x00':
			buf[pos] = '\\'
			buf[pos+1] = '0'
			pos += 2
		case '\n':
			buf[pos] = '\\'
			buf[pos+1] = 'n'
			pos += 2
		case '\r':
			buf[pos] = '\\'
			buf[pos+1] = 'r'
			pos += 2
		case '\x1a':
			buf[pos] = '\\'
			buf[pos+1] = 'Z'
			pos += 2
		case '\'':
			buf[pos] = '\\'
			buf[pos+1] = '\''
			pos += 2
		case '"':
			buf[pos] = '\\'
			buf[pos+1] = '"'
			pos += 2
		case '\\':
			buf[pos] = '\\'
			buf[pos+1] = '\\'
			pos += 2
		case '~': //sql delimeter
			continue
		default:
			buf[pos] = v[i]
			pos++
		}
	}
	//return buf[:pos]
	v = buf[:pos] //PERF
	return v
}
