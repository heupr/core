package backtest

import (
	"bytes"
	"fmt"
	"github.com/google/go-github/github"
	"net/http"
	"net/http/httptest"
	"testing"
)

var heuprTestIssue = `{"action":"opened","issue":{"url":"https://api.github.com/repos/heupr/test/issues/62","repository_url":"https://api.github.com/repos/heupr/test","labels_url":"https://api.github.com/repos/heupr/test/issues/62/labels{/name}","comments_url":"https://api.github.com/repos/heupr/test/issues/62/comments","events_url":"https://api.github.com/repos/heupr/test/issues/62/events","html_url":"https://github.com/heupr/test/issues/62","id":211988771,"number":62,"title":"Darth Test","user":{"login":"taylormike","id":15882362,"avatar_url":"https://avatars3.githubusercontent.com/u/15882362?v=3","gravatar_id":"","url":"https://api.github.com/users/taylormike","html_url":"https://github.com/taylormike","followers_url":"https://api.github.com/users/taylormike/followers","following_url":"https://api.github.com/users/taylormike/following{/other_user}","gists_url":"https://api.github.com/users/taylormike/gists{/gist_id}","starred_url":"https://api.github.com/users/taylormike/starred{/owner}{/repo}","subscriptions_url":"https://api.github.com/users/taylormike/subscriptions","organizations_url":"https://api.github.com/users/taylormike/orgs","repos_url":"https://api.github.com/users/taylormike/repos","events_url":"https://api.github.com/users/taylormike/events{/privacy}","received_events_url":"https://api.github.com/users/taylormike/received_events","type":"User","site_admin":false},"labels":[],"state":"open","locked":false,"assignee":null,"assignees":[],"milestone":null,"comments":0,"created_at":"2017-03-05T22:09:20Z","updated_at":"2017-03-05T22:09:20Z","closed_at":null,"body":"Darth "},"repository":{"id":81689981,"name":"test","full_name":"heupr/test","owner":{"login":"heupr","id":20547820,"avatar_url":"https://avatars1.githubusercontent.com/u/20547820?v=3","gravatar_id":"","url":"https://api.github.com/users/heupr","html_url":"https://github.com/heupr","followers_url":"https://api.github.com/users/heupr/followers","following_url":"https://api.github.com/users/heupr/following{/other_user}","gists_url":"https://api.github.com/users/heupr/gists{/gist_id}","starred_url":"https://api.github.com/users/heupr/starred{/owner}{/repo}","subscriptions_url":"https://api.github.com/users/heupr/subscriptions","organizations_url":"https://api.github.com/users/heupr/orgs","repos_url":"https://api.github.com/users/heupr/repos","events_url":"https://api.github.com/users/heupr/events{/privacy}","received_events_url":"https://api.github.com/users/heupr/received_events","type":"Organization","site_admin":false},"private":true,"html_url":"https://github.com/heupr/test","description":null,"fork":false,"url":"https://api.github.com/repos/heupr/test","forks_url":"https://api.github.com/repos/heupr/test/forks","keys_url":"https://api.github.com/repos/heupr/test/keys{/key_id}","collaborators_url":"https://api.github.com/repos/heupr/test/collaborators{/collaborator}","teams_url":"https://api.github.com/repos/heupr/test/teams","hooks_url":"https://api.github.com/repos/heupr/test/hooks","issue_events_url":"https://api.github.com/repos/heupr/test/issues/events{/number}","events_url":"https://api.github.com/repos/heupr/test/events","assignees_url":"https://api.github.com/repos/heupr/test/assignees{/user}","branches_url":"https://api.github.com/repos/heupr/test/branches{/branch}","tags_url":"https://api.github.com/repos/heupr/test/tags","blobs_url":"https://api.github.com/repos/heupr/test/git/blobs{/sha}","git_tags_url":"https://api.github.com/repos/heupr/test/git/tags{/sha}","git_refs_url":"https://api.github.com/repos/heupr/test/git/refs{/sha}","trees_url":"https://api.github.com/repos/heupr/test/git/trees{/sha}","statuses_url":"https://api.github.com/repos/heupr/test/statuses/{sha}","languages_url":"https://api.github.com/repos/heupr/test/languages","stargazers_url":"https://api.github.com/repos/heupr/test/stargazers","contributors_url":"https://api.github.com/repos/heupr/test/contributors","subscribers_url":"https://api.github.com/repos/heupr/test/subscribers","subscription_url":"https://api.github.com/repos/heupr/test/subscription","commits_url":"https://api.github.com/repos/heupr/test/commits{/sha}","git_commits_url":"https://api.github.com/repos/heupr/test/git/commits{/sha}","comments_url":"https://api.github.com/repos/heupr/test/comments{/number}","issue_comment_url":"https://api.github.com/repos/heupr/test/issues/comments{/number}","contents_url":"https://api.github.com/repos/heupr/test/contents/{+path}","compare_url":"https://api.github.com/repos/heupr/test/compare/{base}...{head}","merges_url":"https://api.github.com/repos/heupr/test/merges","archive_url":"https://api.github.com/repos/heupr/test/{archive_format}{/ref}","downloads_url":"https://api.github.com/repos/heupr/test/downloads","issues_url":"https://api.github.com/repos/heupr/test/issues{/number}","pulls_url":"https://api.github.com/repos/heupr/test/pulls{/number}","milestones_url":"https://api.github.com/repos/heupr/test/milestones{/number}","notifications_url":"https://api.github.com/repos/heupr/test/notifications{?since,all,participating}","labels_url":"https://api.github.com/repos/heupr/test/labels{/name}","releases_url":"https://api.github.com/repos/heupr/test/releases{/id}","deployments_url":"https://api.github.com/repos/heupr/test/deployments","created_at":"2017-02-11T23:31:50Z","updated_at":"2017-02-12T16:42:55Z","pushed_at":"2017-02-11T23:31:51Z","git_url":"git://github.com/heupr/test.git","ssh_url":"git@github.com:heupr/test.git","clone_url":"https://github.com/heupr/test.git","svn_url":"https://github.com/heupr/test","homepage":null,"size":0,"stargazers_count":0,"watchers_count":0,"language":null,"has_issues":true,"has_downloads":true,"has_wiki":true,"has_pages":false,"forks_count":0,"mirror_url":null,"open_issues_count":47,"forks":0,"open_issues":47,"watchers":0,"default_branch":"master"},"organization":{"login":"heupr","id":20547820,"url":"https://api.github.com/orgs/heupr","repos_url":"https://api.github.com/orgs/heupr/repos","events_url":"https://api.github.com/orgs/heupr/events","hooks_url":"https://api.github.com/orgs/heupr/hooks","issues_url":"https://api.github.com/orgs/heupr/issues","members_url":"https://api.github.com/orgs/heupr/members{/member}","public_members_url":"https://api.github.com/orgs/heupr/public_members{/member}","avatar_url":"https://avatars1.githubusercontent.com/u/20547820?v=3","description":"Machine learning-powered contributor integration"},"sender":{"login":"taylormike","id":15882362,"avatar_url":"https://avatars3.githubusercontent.com/u/15882362?v=3","gravatar_id":"","url":"https://api.github.com/users/taylormike","html_url":"https://github.com/taylormike","followers_url":"https://api.github.com/users/taylormike/followers","following_url":"https://api.github.com/users/taylormike/following{/other_user}","gists_url":"https://api.github.com/users/taylormike/gists{/gist_id}","starred_url":"https://api.github.com/users/taylormike/starred{/owner}{/repo}","subscriptions_url":"https://api.github.com/users/taylormike/subscriptions","organizations_url":"https://api.github.com/users/taylormike/orgs","repos_url":"https://api.github.com/users/taylormike/repos","events_url":"https://api.github.com/users/taylormike/events{/privacy}","received_events_url":"https://api.github.com/users/taylormike/received_events","type":"User","site_admin":false}}`
var githubArchiveIssue = []byte(`{"id":"4570706617","type":"IssuesEvent","actor":{"id":2444224,"login":"DirtyHairy","display_login":"DirtyHairy","gravatar_id":"","url":"https://api.github.com/users/DirtyHairy","avatar_url":"https://avatars.githubusercontent.com/u/2444224?"},"repo":{"id":23061486,"name":"6502ts/6502.ts","url":"https://api.github.com/repos/6502ts/6502.ts"},"payload":{"action":"opened","issue":{"url":"https://api.github.com/repos/6502ts/6502.ts/issues/25","repository_url":"https://api.github.com/repos/6502ts/6502.ts","labels_url":"https://api.github.com/repos/6502ts/6502.ts/issues/25/labels{/name}","comments_url":"https://api.github.com/repos/6502ts/6502.ts/issues/25/comments","events_url":"https://api.github.com/repos/6502ts/6502.ts/issues/25/events","html_url":"https://github.com/6502ts/6502.ts/issues/25","id":177288610,"number":25,"title":"Edge compatibility testing and fixing","user":{"login":"DirtyHairy","id":2444224,"avatar_url":"https://avatars.githubusercontent.com/u/2444224?v=3","gravatar_id":"","url":"https://api.github.com/users/DirtyHairy","html_url":"https://github.com/DirtyHairy","followers_url":"https://api.github.com/users/DirtyHairy/followers","following_url":"https://api.github.com/users/DirtyHairy/following{/other_user}","gists_url":"https://api.github.com/users/DirtyHairy/gists{/gist_id}","starred_url":"https://api.github.com/users/DirtyHairy/starred{/owner}{/repo}","subscriptions_url":"https://api.github.com/users/DirtyHairy/subscriptions","organizations_url":"https://api.github.com/users/DirtyHairy/orgs","repos_url":"https://api.github.com/users/DirtyHairy/repos","events_url":"https://api.github.com/users/DirtyHairy/events{/privacy}","received_events_url":"https://api.github.com/users/DirtyHairy/received_events","type":"User","site_admin":false},"labels":[{"url":"https://api.github.com/repos/6502ts/6502.ts/labels/enhancement","name":"enhancement","color":"84b6eb"},{"url":"https://api.github.com/repos/6502ts/6502.ts/labels/help%20wanted","name":"help wanted","color":"159818"}],"state":"open","locked":false,"assignee":null,"assignees":[],"milestone":null,"comments":0,"created_at":"2016-09-15T21:00:06Z","updated_at":"2016-09-15T21:00:06Z","closed_at":null,"body":"Does it run in edge?"}},"public":true,"created_at":"2016-09-15T21:00:07Z","org":{"id":22204843,"login":"6502ts","gravatar_id":"","url":"https://api.github.com/orgs/6502ts","avatar_url":"https://avatars.githubusercontent.com/u/22204843?"}}`)

//Test Workflow / Next Immediate Steps
//1) Parse a github archive file (No code yet)
//2) Transform the github archive file into payloads (No code yet)
//3) Call HttpPost on the parsed github archive payloads (Some code is written HttpPost)

//Open/Non Trivial Problem
//4?) Inside the HeuprServer when calling Predict on the parsed payloads/ replayed issues we don't want that to hit github.... We want the calls to Predict to hit our Replay Server. This will allow us to keep score on prediction accuracy.
//Items
//The Github client in the HeuprServer will need to be plug & play
//The Replay Server will need it's own predict handler that Heupr calls into.

//Open/Non Trivial Problem
//5?) Inside the HeuprServer when calling GetIssues/GetPulls during the initial model bootstrap we don't want that to hit github.... We want the calls to hit our Replay Server or the CachedGateway.
//Items
//The Github client in the HeuprServer will need to be plug & play
//I Personally need to think more about this one

//TODO: This test is a good first step as the results can be checked in Ngrok (http://127.0.0.1:4040/inspect/http) but we should really be parsing the github archive file and passing that in.
//This test closely resembles what the replay main method will look like. (This test can probably be removed later at some point)
func TestHTTPPost(t *testing.T) {
	replayServer := CreateReplayServer()
	buf := bytes.NewBufferString(heuprTestIssue)
	replayServer.HTTPPost(buf)
}

func TestArchiveIssue(t *testing.T) {
	_, err := github.ParseWebHook("issues", githubArchiveIssue)
	if err != nil {
		panic(err)
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/hook", func(w http.ResponseWriter, r *http.Request) {
		//TODO: Fix this does not have a chance to execute because the test dies to quickly
		err := validateGithubEvent(r)
		if err != nil {
			t.Error(err)
		}
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	replayServer := CreateReplayServer()
	replayServer.HTTPPost(bytes.NewBuffer(githubArchiveIssue))
}

func TestHTTPMessageIntegrity(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc(func(w http.ResponseWriter, r *http.Request) {
		//TODO: Fix this does not have a chance to execute because the test dies to quickly
		err := validateGithubEvent(r)
		if err != nil {
			t.Error(err)
		}
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	replayServer := CreateReplayServer()
	buf := bytes.NewBufferString(heuprTestIssue)
	replayServer.HTTPPost(buf)
}

func validateGithubEvent(r *http.Request) error {
	eventType := r.Header.Get("X-Github-Event")
	if eventType != "issues" {
		return fmt.Errorf("Ignoring '%v' event", eventType)
	}
	payload, err := github.ValidatePayload(r, []byte(secretKey))
	if err != nil {
		return fmt.Errorf("Could not validate payload %v", err)
	}
	_, err = github.ParseWebHook(github.WebHookType(r), payload)
	if err != nil {
		return fmt.Errorf("Could not parse webhook %v", err)
	}
	return nil
}
