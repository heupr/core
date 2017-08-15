package frontend

import (
	"encoding/json"
	"html/template"
	"net/http"
	"net/url"

	"github.com/google/go-github/github"
	"github.com/gorilla/schema"
	"go.uber.org/zap"
	"golang.org/x/oauth2"

	"coralreefci/utils"
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
        <p><a href="/setup_complete">Submit selection(s)</a></p>
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
	Repos []*github.Repository
}

func (fs *FrontendServer) githubCallbackHandler(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("state") != oaState {
		utils.AppLog.Error("incorrect callback state value")
		http.Redirect(w, r, "/", http.StatusForbidden)
		return
	}

	token, err := oaConfig.Exchange(oauth2.NoContext, r.FormValue("code"))
	if err != nil {
		utils.AppLog.Error("callback token exchange: ", zap.Error(err))
		http.Redirect(w, r, "/", http.StatusInternalServerError)
		return
	}

	oaClient := oaConfig.Client(oauth2.NoContext, token)
	client := github.NewClient(oaClient)

	repos, err := listRepositories(client)
	if err != nil {
		utils.AppLog.Error("callback list user repos: ", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.Method == "GET" {
		tmpl, err := template.New("options").Parse(options)
		if err != nil {
			utils.AppLog.Error("failure parsing user options template: ", zap.Error(err))
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
			utils.AppLog.Error("failure decoding postform: ", zap.Error(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		for i := 0; i < len(results.Repos); i++ {
			repo := results.Repos[i]
			if err := fs.NewHook(repo, client); err != nil {
				utils.AppLog.Error("repo hook placement: ", zap.Error(err))
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			tokenByte, err := json.Marshal(token)
			if err != nil {
				utils.AppLog.Error("failure converting callback token: ", zap.Error(err))
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if err := fs.Database.Store("token", *repo.ID, tokenByte); err != nil {
				utils.AppLog.Error("failure storing token in bolt: ", zap.Error(err))
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			resp, err := http.PostForm("/activate-repos", url.Values{
				"state": {BackendSecret},
				"repos": {string(*repo.ID), string(*repo.Name), string(*repo.Owner.Login)},
				"token": {string(tokenByte)},
			})
			if err != nil {
				utils.AppLog.Error("error posting to activation service: ", zap.Error(err))
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer resp.Body.Close()
		}
	}
}

func completeHandle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(setup))
}
