package frontend

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/google/go-github/github"
	"github.com/gorilla/schema"
	"go.uber.org/zap"
	"golang.org/x/oauth2"

	"core/utils"
)

type Repo struct {
    ID int
    FullName string `schema:"-"`
    Selected bool `schema:"-"`
}

type RepoForm struct {
    Name string
    Repos []Repo `schema:"Repos"`
		Limit int `schema:"Limit"`
}

const options = `
<html>
    <title>
        Heupr
    </title>
    <body>
    <form action="/github_oauth_cb" method="post">
        <p>Choose your repo(s):</p>
        {{range $idx, $repo := .Repos}}
            <p><input type="checkbox" name="Repos.{{ $idx }}.ID" value="{{ $repo.ID }}">{{ $repo.FullName }}</p>
        {{end}}
        <p>Select range of issues to be assigned:</p>
        <input type="radio" name="Limit" value=5> Week<br>
        <input type="radio" name="Limit" value=30> Month<br>
        <input type="radio" name="Limit" value=365> Year<br>
        <input type="radio" name="Limit" value=100000 checked> All<br>
				<p>
            <input type="submit" value="Submit selection(s)">
        </p>
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


func (fs *FrontendServer) githubCallbackHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
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

		client := github.NewClient(oaConfig.Client(oauth2.NoContext, token))

		repos, err := listRepositories(client)
		if err != nil {
			utils.AppLog.Error("callback list user repos: ", zap.Error(err))
			http.Error(w, "Apologies, we are experiencing technical difficulties. Standby for a signup confirmation email", http.StatusInternalServerError)
			return
		}

		tmpl, err := template.New("options").Parse(options)
		if err != nil {
			utils.AppLog.Error("failure parsing user options template: ", zap.Error(err))
			http.Error(w, "Apologies, we are experiencing technical difficulties. Standby for a signup confirmation email", http.StatusInternalServerError)
			return
		}
		reposList := make([]Repo, len(repos))
		for i := 0; i < len(reposList); i++ {
			reposList[i] = Repo{ID: *repos[i].ID, FullName: *repos[i].FullName, Selected: false}
			fs.state.Tokens[*repos[i].ID] = token
		}
		repoForm := RepoForm{Name: "Default", Repos: reposList}
		tmpl.Execute(w, repoForm)
	} else {
		r.ParseForm()
		repoForm := new(RepoForm)
		err := decoder.Decode(repoForm, r.PostForm)
		if err != nil {
			utils.AppLog.Error("failure decoding postform: ", zap.Error(err))
			http.Error(w, "Apologies, we are experiencing technical difficulties. Standby for a signup confirmation email", http.StatusInternalServerError)
			return
		}
		for i := 0; i < len(repoForm.Repos); i++ {
			if repoForm.Repos[i].ID == 0 {
				continue
			}
			token := fs.state.Tokens[repoForm.Repos[i].ID]
			if token == nil {
				utils.AppLog.Error("failed to lookup repo from shared state ", zap.Int("RepoID", repoForm.Repos[i].ID))
				http.Error(w, "Apologies, we are experiencing technical difficulties. Standby for a signup confirmation email", http.StatusInternalServerError)
				return
			}
			client := github.NewClient(oaConfig.Client(oauth2.NoContext, token))

			repo, _, err := client.Repositories.GetByID(context.Background(), repoForm.Repos[i].ID)
			if err != nil {
				utils.AppLog.Error("Callback get by id", zap.Error(err))
				return
			}
			if err := fs.AutomaticWhitelist(*repo); err != nil {
				utils.AppLog.Error("whitelist failure: ", zap.Error(err))
				http.Error(w, "Maximum allowed beta users reached. Standby for a signup confirmation email", http.StatusInternalServerError)
				return
			}

			limit := time.Now().AddDate(0, 0, -repoForm.Limit)

			if err := fs.NewHook(repo, client); err != nil {
				utils.AppLog.Error("repo hook placement: ", zap.Error(err))
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			tokenByte, err := json.Marshal(token)
			if err != nil {
				utils.AppLog.Error("converting callback token: ", zap.Error(err))
				http.Error(w, "Apologies, we are experiencing technical difficulties. Standby for a signup confirmation email", http.StatusInternalServerError)
				return
			}
			if err := fs.Database.Store("token", *repo.ID, tokenByte); err != nil {
				utils.AppLog.Error("storing token in bolt: ", zap.Error(err))
				http.Error(w, "Apologies, we are experiencing technical difficulties. Standby for a signup confirmation email", http.StatusInternalServerError)
				return
			}

			limitByte, err := json.Marshal(limit)
			if err != nil {
				utils.AppLog.Error("converting callback limit: ", zap.Error(err))
				http.Error(w, "Apologies, we are experiencing technical difficulties. Standby for a signup confirmation email", http.StatusInternalServerError)
				return
			}
			if err := fs.Database.Store("limit", *repo.ID, limitByte); err != nil {
				utils.AppLog.Error("storing limit in bolt: ", zap.Error(err))
				http.Error(w, "Apologies, we are experiencing technical difficulties. Standby for a signup confirmation email", http.StatusInternalServerError)
				return
			}
			activationParams := struct {
				Repo  github.Repository `json:"repo"`
				Token *oauth2.Token     `json:"token"`
				Limit time.Time         `json:"limit"`
			}{
				*repo,
				token,
				limit,
			}
			payload, err := json.Marshal(activationParams)
			if err != nil {
				utils.AppLog.Error("failure converting activation parameters: ", zap.Error(err))
				http.Error(w, "Apologies, we are experiencing technical difficulties. Standby for a signup confirmation email", http.StatusInternalServerError)
				return
			}
			req, err := http.NewRequest("POST", utils.Config.ActivationServiceEndpoint, bytes.NewBuffer(payload))
			if err != nil {
				utils.AppLog.Error("failed to create http request:", zap.Error(err))
				http.Error(w, "Apologies, we are experiencing technical difficulties. Standby for a signup confirmation email", http.StatusInternalServerError)
				return
			}
			req.Header.Set("content-type", "application/json")
			resp, err := fs.httpClient.Do(req)
			if err != nil {
				utils.AppLog.Error("failed internal post call:", zap.Error(err))
				http.Error(w, "Apologies, we are experiencing technical difficulties. Standby for a signup confirmation email", http.StatusInternalServerError)
				return
			} else {
				defer resp.Body.Close()
			}
			utils.SlackLog.Info(fmt.Sprintf("Signup %v", *repo.FullName))
		}
		utils.SlackLog.Info("Complete handle entered - user signed up")
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(setup))
	}
}
