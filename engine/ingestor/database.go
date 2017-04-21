package ingestor

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/google/go-github/github"
	"io"
)

type Event struct {
	Type    string             `json:"type"`
	Repo    github.Repository  `json:"repo"`
	Payload github.IssuesEvent `json:"payload"`
}

type Value interface{}

type Database struct {
	db         *sql.DB
	BufferPool Pool
}

func (d *Database) Open() {
	mysql, err := sql.Open("mysql", "root@/heupr?interpolateParams=true")
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	d.db = mysql
}

func (d *Database) Close() {
	d.db.Close()
}

func (d *Database) EnableRepo(repoId int) {
	var buffer bytes.Buffer
	archRepoInsert := "INSERT INTO arch_repos(repository_id, enabled) VALUES"
	valuesFmt := "(?,?)"

	buffer.WriteString(archRepoInsert)
	buffer.WriteString(valuesFmt)
	_, err := d.db.Exec(buffer.String(), repoId, true)
	fmt.Println(err)
}

//TODO: Use LOAD DATA INFILE
func (d *Database) BulkInsertBacktestEvents(events []Event) {
	var buffer bytes.Buffer
	eventsInsert := "INSERT INTO backtest_events(repo_id, repo_name, payload) VALUES"
	eventsValuesFmt := "(?,?,?)"
	numValues := 3

	buffer.WriteString(eventsInsert)
	delimeter := ""
	values := make([]interface{}, len(events)*numValues)
	for i := 0; i < len(events); i++ {
		buffer.WriteString(delimeter)
		buffer.WriteString(eventsValuesFmt)
		offset := i * numValues

		values[offset+0] = events[i].Repo.ID
		values[offset+1] = events[i].Repo.Name
		payload, _ := json.Marshal(events[i])
		values[offset+2] = stripCtlAndExtFromBytes(payload)
		delimeter = ","
	}
	_, err := d.db.Exec(buffer.String(), values...)
	if err != nil {
		fmt.Println(err)
	}
}

func (d *Database) ReadBacktestEvents(repo string) ([]Event, error) {
	events := []Event{}
	var payload []byte
	results, err := d.db.Query("select payload from backtest_events where repo_name=?", repo)
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
		json.Unmarshal(payload, &event)
		events = append(events, event)
	}
	err = results.Err()
	if err != nil {
		return nil, err
	}
	return events, nil
}

func (d *Database) InsertIssue(issue github.Issue) {
	var buffer bytes.Buffer
	eventsInsert := "INSERT INTO github_events(repo_id,issues_id,number,payload,is_pr,is_closed) VALUES"
	eventsValuesFmt := "(?,?,?,?,0,?)"
	numValues := 5

	buffer.WriteString(eventsInsert)
	buffer.WriteString(eventsValuesFmt)
	values := make([]interface{}, numValues)
	values[0] = *issue.Repository.ID
	values[1] = issue.ID
	values[2] = issue.Number
	payload, _ := json.Marshal(issue)
	values[3] = stripCtlAndExtFromBytes(payload)
	if issue.ClosedAt == nil {
		values[4] = false
	} else {
		values[4] = true
	}
	_, err := d.db.Exec(buffer.String(), values...)
	if err != nil {
		fmt.Println(err)
	}
}

func (d *Database) BulkInsertIssues(issues []*github.Issue) {
	buffer := d.BufferPool.Get()

	for i := 0; i < len(issues); i++ {
		buffer.AppendInt(int64(*issues[i].Repository.ID))
		buffer.AppendByte('~')
		buffer.AppendInt(int64(*issues[i].ID))
		buffer.AppendByte('~')
		buffer.AppendInt(int64(*issues[i].Number))
		buffer.AppendByte('~')
		payload, _ := json.Marshal(*issues[i])
		_, _ = buffer.Write(escapeBytesBackslash(stripCtlAndExtFromBytes(payload)))
		buffer.AppendByte('~')
		buffer.AppendInt(0)
		buffer.AppendByte('~')
		if issues[i].ClosedAt == nil {
			buffer.AppendInt(1)
		} else {
			buffer.AppendInt(0)
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
	_, err := d.db.Exec("LOAD DATA LOCAL INFILE 'Reader::data' INTO TABLE github_events FIELDS TERMINATED BY '~' LINES TERMINATED BY '\n' (repo_id,issues_id,number,payload,is_pr,is_closed)")
	if err != nil {
		fmt.Println(err)
	}
}

func (d *Database) InsertPullRequest(pull github.PullRequest) {
	var buffer bytes.Buffer
	eventsInsert := "INSERT INTO github_events(payload,is_pr,is_closed) VALUES"
	eventsValuesFmt := "(?,1,?)"
	numValues := 2

	buffer.WriteString(eventsInsert)
	buffer.WriteString(eventsValuesFmt)
	values := make([]interface{}, numValues)
	payload, _ := json.Marshal(pull)
	values[0] = stripCtlAndExtFromBytes(payload)
	if pull.ClosedAt == nil {
		values[1] = false
	} else {
		values[1] = true
	}
	_, err := d.db.Exec(buffer.String(), values...)
	fmt.Println(err)
}

func (d *Database) BulkInsertPullRequests(pulls []*github.PullRequest) {
	buffer := d.BufferPool.Get()

	for i := 0; i < len(pulls); i++ {
		buffer.AppendInt(int64(*pulls[i].Base.Repo.ID))
		buffer.AppendByte('~')
		buffer.AppendInt(int64(*pulls[i].ID))
		buffer.AppendByte('~')
		buffer.AppendInt(int64(*pulls[i].Number))
		buffer.AppendByte('~')
		payload, _ := json.Marshal(*pulls[i])
		_, _ = buffer.Write(escapeBytesBackslash(stripCtlAndExtFromBytes(payload)))
		buffer.AppendByte('~')
		buffer.AppendInt(1)
		buffer.AppendByte('~')
		if pulls[i].ClosedAt == nil {
			buffer.AppendInt(1)
		} else {
			buffer.AppendInt(0)
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
	_, err := d.db.Exec("LOAD DATA LOCAL INFILE 'Reader::data' INTO TABLE github_events FIELDS TERMINATED BY '~' LINES TERMINATED BY '\n' (repo_id,issues_id,number,payload,is_pr,is_closed)")
	if err != nil {
		fmt.Println(err)
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
	return b[:bl]
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
	return buf[:pos]
}
