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

func (d *Database) EnableRepo(repoId int) {
	var buffer bytes.Buffer
	archRepoInsert := "INSERT INTO arch_repos(repository_id, enabled) VALUES"
	valuesFmt := "(?,?)"

	buffer.WriteString(archRepoInsert)
	buffer.WriteString(valuesFmt)
	_, err := d.db.Exec(buffer.String(), repoId, true)
	fmt.Println(err)
}

func (d *Database) BulkInsertIssues(issues []*github.Issue, repoId int) {
	var buffer bytes.Buffer
	issuesInsert := "INSERT INTO issues(issues_id,number,state,locked,title,body,user_id,comments,closed_at,created_at,updated_at,closed_by,url,html_url,milestone,pull_request_links,repository_id) VALUES"
	issuesValuesFmt := "(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
	numValues := 17

	buffer.WriteString(issuesInsert)
	delimeter := ""
	iValues := make([]interface{}, len(issues)*numValues)
	for i := 0; i < len(issues); i++ {
		buffer.WriteString(delimeter)
		buffer.WriteString(issuesValuesFmt)

		offset := i * numValues
		iValues[offset+0] = issues[i].ID
		iValues[offset+1] = issues[i].Number
		iValues[offset+2] = issues[i].State
		iValues[offset+3] = issues[i].Locked
		iValues[offset+4] = issues[i].Title
		if issues[i].Body != nil {
			iValues[offset+5] = issues[i].Body
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
		if issues[i].Milestone != nil {
			iValues[offset+14] = issues[i].Milestone.DueOn
		} else {
			iValues[offset+14] = nil
		}

		if issues[i].PullRequestLinks != nil {
			iValues[offset+15] = issues[i].PullRequestLinks.URL
		} else {
			iValues[offset+15] = nil
		}
		iValues[offset+16] = repoId

		delimeter = ","
	}
	result, err := d.db.Exec(buffer.String(), iValues...)
	fmt.Println(err)
	issues_fk, _ := result.LastInsertId()
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

func (d *Database) BulkInsertPullRequests(pulls []*github.PullRequest, repoId int) {
	var buffer bytes.Buffer
	pullsInsert := "INSERT INTO pull_requests(pr_id,number,state,title,body,created_at,updated_at,closed_at,merged_at,user_id,merged,mergable,merged_by_user_id,comments,commits,additions,deletions,changed_files,url,html_url,issue_url,statuses_url,diff_url,patch_url,review_comments_url,review_comment_url,assignee_user_id,milestone,maintainer_can_modify,head_pr_branch_id,base_pr_branch_id,repository_id) VALUES"
	pullsValuesFmt := "(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
	numValues := 32

	buffer.WriteString(pullsInsert)
	delimeter := ""
	pValues := make([]interface{}, len(pulls)*numValues)
	for i := 0; i < len(pulls); i++ {
		buffer.WriteString(delimeter)
		buffer.WriteString(pullsValuesFmt)

		offset := i * numValues
		pValues[offset+0] = pulls[i].ID
		pValues[offset+1] = pulls[i].Number
		pValues[offset+2] = pulls[i].State
		pValues[offset+3] = pulls[i].Title
		if pulls[i].Body != nil {
			pValues[offset+4] = pulls[i].Body
		} else {
			pValues[offset+4] = nil
		}
		pValues[offset+5] = pulls[i].CreatedAt
		pValues[offset+6] = pulls[i].UpdatedAt
		pValues[offset+7] = pulls[i].ClosedAt
		pValues[offset+8] = pulls[i].MergedAt
		if pulls[i].User != nil {
			pValues[offset+9] = *pulls[i].User.ID
		} else {
			pValues[offset+9] = nil
		}
		pValues[offset+10] = pulls[i].Merged
		pValues[offset+11] = pulls[i].Mergeable
		pValues[offset+12] = pulls[i].MergedBy
		pValues[offset+13] = pulls[i].Comments
		pValues[offset+14] = pulls[i].Commits
		pValues[offset+15] = pulls[i].Additions
		pValues[offset+16] = pulls[i].Deletions
		pValues[offset+17] = pulls[i].ChangedFiles
		pValues[offset+18] = pulls[i].URL
		pValues[offset+19] = pulls[i].HTMLURL
		pValues[offset+20] = pulls[i].IssueURL
		pValues[offset+21] = pulls[i].StatusesURL
		pValues[offset+22] = pulls[i].DiffURL
		pValues[offset+23] = pulls[i].PatchURL
		pValues[offset+24] = pulls[i].ReviewCommentsURL
		pValues[offset+25] = pulls[i].ReviewCommentURL
		if pulls[i].Assignee != nil {
			pValues[offset+26] = pulls[i].Assignee.ID
		} else {
			pValues[offset+26] = nil
		}
		if pulls[i].Milestone != nil {
			pValues[offset+27] = pulls[i].Milestone.DueOn
		} else {
			pValues[offset+27] = nil
		}
		pValues[offset+28] = pulls[i].MaintainerCanModify
		pValues[offset+29] = -1
		pValues[offset+30] = -1
		pValues[offset+31] = repoId

		delimeter = ","
	}
	_, err := d.db.Exec(buffer.String(), pValues...)
	fmt.Println(err)
	buffer.Reset()
}

func (d *Database) InsertPullRequest(pull github.PullRequest) {

}
