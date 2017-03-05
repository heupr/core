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
	issuesInsert := "INSERT INTO issues(issues_id,number,state,locked,title,body,user_id,lables,assignees_id) VALUES"
	issuesValuesFmt := "(?,?,?,?,?,?,?,?,?)"
	numValues := 9

	//assignees_insert := "INSERT INTO assignees(issue_fk, user_id) VALUES"
	//assignees_values_fmt := "(?,?)"
	buffer.WriteString(issuesInsert)
	delimeter := ""
	values := make([]interface{}, len(issues)*numValues)
	for i := 0; i < len(issues); i++ {
		buffer.WriteString(delimeter)
		buffer.WriteString(issuesValuesFmt)

		values[0+i] = issues[i].ID
		values[1+i] = issues[i].Number
		values[2+i] = issues[i].State
		values[3+i] = false
		values[4+i] = issues[i].Title
		values[5+i] = issues[i].Body
		values[6+i] = issues[i].User.ID
		values[7+i] = -1
		values[8+i] = -1
		delimeter = ","
	}
	_, err := d.db.Exec(buffer.String(), values...)
	fmt.Println(err)
}

func (d *Database) InsertIssue(issue github.Issue) {

}

func (d *Database) BulkInsertPullRequest(pulls []github.PullRequest) {

}

func (d *Database) InsertPullRequest(pull github.PullRequest) {

}
