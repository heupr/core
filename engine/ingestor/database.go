package ingestor

import (
	"bytes"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/go-github/github"
)

type Value interface{}

type Database struct {
	db *sql.DB
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

func (d *Database) BulkInsertIssues(issues []*github.Issue) {
	var buffer bytes.Buffer
	issuesInsert := "INSERT INTO issues(issues_id,number,state,locked,title,body,user_id,comments,closed_at,created_at,updated_at,closed_by,url,html_url,pull_request_links,repository_id) VALUES"
	issuesValuesFmt := "(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
	numValues := 16

	buffer.WriteString(issuesInsert)
	delimeter := ""
	iValues := make([]interface{}, len(issues)*numValues)
	for i := 0; i < len(issues); i++ {
		buffer.WriteString(delimeter)
		buffer.WriteString(issuesValuesFmt)

		offset := i*numValues
		iValues[offset+0] = issues[i].ID
		iValues[offset+1] = issues[i].Number
		iValues[offset+2] = issues[i].State
		iValues[offset+3] = issues[i].Locked
		iValues[offset+4] = issues[i].Title
		if issues[i].Body != nil {
			iValues[offset+5] = *issues[i].Body
		} else {
			iValues[offset+5] = nil
		}
		iValues[offset+6] = issues[i].User.ID
		iValues[offset+7] = issues[i].Comments
		iValues[offset+8] = issues[i].ClosedAt
		iValues[offset+9] = issues[i].CreatedAt
		iValues[offset+10] = issues[i].UpdatedAt
		iValues[offset+11] = issues[i].ClosedBy
		iValues[offset+12] = issues[i].URL
		iValues[offset+13] = issues[i].HTMLURL
		if issues[i].PullRequestLinks != nil {
			iValues[offset+14] = issues[i].PullRequestLinks.URL
		} else {
			iValues[offset+14] = nil
		}
		if issues[i].Repository != nil {
			iValues[offset+15] = issues[i].Repository.ID
		} else {
			iValues[offset+15] = nil
		}

		delimeter = ","
	}
	result, err := d.db.Exec(buffer.String(), iValues...)
	fmt.Println(err)
	issues_fk,_ := result.LastInsertId()
	buffer.Reset()

	assignees_insert := "INSERT INTO assignees(issue_fk, user_id) VALUES"
	assignees_values_fmt := "(?,?)"

	buffer.WriteString(assignees_insert)
	delimeter = ""
	var aValues []interface{}
	for i := 0; i < len(issues); i++ {
		if issues[i].Assignees != nil {
				for j := 0; j < len(issues[i].Assignees); j++ {
					buffer.WriteString(delimeter)
					buffer.WriteString(assignees_values_fmt)
					aValues = append(aValues, issues_fk)
					aValues = append(aValues, issues[i].Assignees[j].ID)
					delimeter = ","
				}
		 }
		 issues_fk++
	}
	_, err = d.db.Exec(buffer.String(), aValues...)
	fmt.Println(err)
}


func (d *Database) InsertIssue(issue github.Issue) {

}

func (d *Database) BulkInsertPullRequest(pulls []github.PullRequest) {

}

func (d *Database) InsertPullRequest(pull github.PullRequest) {

}
