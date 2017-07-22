package frontend

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"net/url"

	"github.com/google/go-github/github"
	"github.com/gorilla/schema"
	"golang.org/x/oauth2"
)

const options = `
<html>
    <title>
        Heupr
    </title>
    <body>
    <form action="/github_oauth_cb" method="post">
        <p>Choose your repo(s):</p>
        {{range $idx, $repo := .Repos}}
            <p><input type="checkbox" name="Repos.{{ $idx }}.ID" value="{{ $repo.ID }}">{{ $repo.Owner.Login }}/{{ $repo.Name }}</p>
        {{end}}
    </form>
    </body>
</html>
`

const setup = `
<html>
    <title>
        Heupr
    </title>
    <body>
        <p>Awesome! Setup is complete!</p>
        <p>Issue assignments will go out in a few minutes through GitHub</p>
        <p>Return to the <a href="/">main page</a></p>
    </body>
</html>
`

const BackendSecret = "fear-is-my-ally"

var decoder = schema.NewDecoder()

type Resources struct {
	repos []*github.Repository
}

func (fs *FrontendServer) githubCallbackHandler(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("state") != oaState {
		http.Redirect(w, r, "/", http.StatusForbidden)
		return
	}

	token, err := oaConfig.Exchange(oauth2.NoContext, r.FormValue("code"))
	if err != nil {
		http.Redirect(w, r, "/", http.StatusInternalServerError)
		return
	}

	oaClient := oaConfig.Client(oauth2.NoContext, token)
	client := github.NewClient(oaClient)

	repos, err := listRepositories(client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.Method == "GET" {
		tmpl, err := template.New("options").Parse(options)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		resources := Resources{repos}
		tmpl.Execute(w, resources)
	} else {
		r.ParseForm()
		results := new(Resources)
		err := decoder.Decode(results, r.PostForm)
		if err != nil {
			fmt.Println(err)
		}
		for i := 0; i < len(results.repos); i++ {
			repo := results.repos[i]
			if err := fs.NewHook(repo, client); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			tokenByte, err := json.Marshal(token)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if err := fs.Database.Store("token", *repo.ID, tokenByte); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			// TODO: Build logic to send all repo IDs simultaneously via the POST request.
			resp, err := http.PostForm("/activate-repos", url.Values{
				"state": {BackendSecret},
				"repos": {string(*repo.ID), string(*repo.Name), string(*repo.Owner.Login)},
				"token": {string(tokenByte)},
			})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer resp.Body.Close()
		}
	}
	http.Redirect(w, r, "/setup_complete", http.StatusPermanentRedirect)
}

func completeHandle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(setup))
}
