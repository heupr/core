package ingestor

import (
  "fmt"
  "github.com/google/go-github/github"
  "database/sql"
_ "github.com/go-sql-driver/mysql"
)

type Database struct {
  db *sql.DB
}

func (d *Database) Open() {
  mysql, err := sql.Open("mysql", "root@/heupr?interpolateParams=true")
    if err != nil {
        panic(err.Error())  // Just for example purpose. You should use proper error handling instead of panic
    }
  d.db = mysql
}

func (d *Database) Close() {
  d.db.Close()
}

func (d *Database) BulkInsertIssues(issues []github.Issue) {

}

func (d *Database) InsertIssue(issue github.Issue) {
  fmt.Println(issue)
  //query := fmt.Sprintf("INSERT INTO user_permissions (user_id) VALUES(%v)", issue.User.ID)
  query := fmt.Sprintf("INSERT INTO Issues (issues_id, number, state, locked, title, body, user_id, lables, assignee, comments, closed_at, created_at, updated_at, closed_by, url, html_url, milestone, pull_request_links, repository, reactions, assignees, text_matches) VALUES(%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v)",*issue.ID, *issue.Number, *issue.State, false, *issue.Title, *issue.Body, -1, -1, -1, -1, *issue.ClosedAt, *issue.CreatedAt, *issue.UpdatedAt, -1, *issue.URL, *issue.HTMLURL, -1, -1, -1, -1, -1, -1)
  //query := fmt.Sprintf("INSERT INTO issues VALUES(%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v)",*issue.ID, *issue.Number, *issue.State, false, *issue.Title, *issue.Body, *issue.User.ID, -1, *issue.Assignee.ID, *issue.Comments, *issue.ClosedAt, *issue.CreatedAt, *issue.UpdatedAt, -1, *issue.URL, *issue.HTMLURL, *issue.Milestone.ID, -1, *issue.Repository.ID, -1, -1, -1)
  result, err := d.db.Exec(query)
  fmt.Println(result)
  fmt.Println(err)
	//fmt.Printf("Rows affected %v", result.RowsAffected())
}

func (d *Database) BulkInsertPullRequest(pulls []github.PullRequest) {

}

func (d *Database) InsertPullRequest(pull github.PullRequest) {

}
