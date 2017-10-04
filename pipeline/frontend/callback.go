package frontend

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/boltdb/bolt"
	"github.com/google/go-github/github"
	"github.com/gorilla/schema"
	"go.uber.org/zap"
	"golang.org/x/oauth2"

	"core/utils"
)

type Repo struct {
	ID       int
	FullName string `schema:"-"`
	Selected bool   `schema:"-"`
}

type RepoForm struct {
	Name  string
	Repos []Repo `schema:"Repos"`
	Limit int    `schema:"Limit"`
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
    		<p><input type="submit" value="Submit selection(s)"></p>
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

func httpRedirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://heupr.io", http.StatusMovedPermanently)
}

func (fs *FrontendServer) githubCallbackHandler(w http.ResponseWriter, r *http.Request) {
	userErrMsg := "Apologies, we are experiencing technical difficulties. Please, let us look into the issue and try again in a few minutes"
	if r.Method == "GET" {
		if r.FormValue("state") != oaState {
			utils.AppLog.Error("incorrect callback state value - received", zap.String("state", r.FormValue("state")))
			http.Redirect(w, r, "/", http.StatusForbidden)
			return
		}
		token, err := oaConfig.Exchange(oauth2.NoContext, r.FormValue("code"))
		if err != nil {
			utils.AppLog.Error("callback token exchange failure", zap.Error(err))
			http.Redirect(w, r, "/", http.StatusInternalServerError)
			return
		}

		client := github.NewClient(oaConfig.Client(oauth2.NoContext, token))

		repos, err := listRepositories(client)
		if err != nil {
			utils.AppLog.Error("callback list user repos failure", zap.Error(err))
			http.Error(w, userErrMsg, http.StatusInternalServerError)
			return
		}

		tmpl, err := template.New("options").Parse(options)
		if err != nil {
			utils.AppLog.Error("failure parsing user options template", zap.Error(err))
			http.Error(w, userErrMsg, http.StatusInternalServerError)
			return
		}
		reposList := make([]Repo, len(repos))
		fs.state.Lock()
		for i := 0; i < len(reposList); i++ {
			reposList[i] = Repo{
				ID:       *repos[i].ID,
				FullName: *repos[i].FullName,
				Selected: false,
			}
			fs.state.Tokens[*repos[i].ID] = token
		}
		fs.state.Unlock()
		repoForm := RepoForm{Name: "Default", Repos: reposList}
		tmpl.Execute(w, repoForm)
	} else {
		r.ParseForm()
		repoForm := new(RepoForm)
		err := decoder.Decode(repoForm, r.PostForm)
		if err != nil {
			utils.AppLog.Error("failure decoding postform", zap.Error(err))
			http.Error(w, userErrMsg, http.StatusInternalServerError)
			return
		}
		for i := 0; i < len(repoForm.Repos); i++ {
			if repoForm.Repos[i].ID == 0 {
				continue
			}
			fs.state.Lock()
			token := fs.state.Tokens[repoForm.Repos[i].ID]
			fs.state.Unlock()

			if token == nil {
				utils.AppLog.Error("failed to lookup repo from shared state ", zap.Int("RepoID", repoForm.Repos[i].ID))
				http.Error(w, userErrMsg, http.StatusInternalServerError)
				return
			}
			client := github.NewClient(oaConfig.Client(oauth2.NoContext, token))

			repo, _, err := client.Repositories.GetByID(context.Background(), repoForm.Repos[i].ID)
			if err != nil {
				utils.AppLog.Error("callback get by id failed", zap.Error(err))
				return
			}
			check, err := fs.CheckWhitelist(*repo)
			if err != nil {
				utils.AppLog.Error("whitelist failure", zap.Error(err))
				http.Error(w, userErrMsg, http.StatusInternalServerError)
				return
			} else if check != "" {
				maxMsg := fmt.Sprintf("whitelist maximum reached; rejected %v", check)
				utils.SlackLog.Info(maxMsg)
				utils.AppLog.Info(maxMsg)
				http.Error(w, "Maximum allowed beta users reached. Please, send us an email if you are interested in signing up", http.StatusInternalServerError)
			}

			limit := time.Now().AddDate(0, 0, -repoForm.Limit)

			if err := fs.NewHook(repo, client); err != nil {
				utils.AppLog.Error("repo hook placement", zap.Error(err))
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			tokenByte, err := json.Marshal(token)
			if err != nil {
				utils.AppLog.Error("converting callback token", zap.Error(err))
				http.Error(w, userErrMsg, http.StatusInternalServerError)
				return
			}

			go func() {
				boltDB, err := bolt.Open(utils.Config.BoltDBPath, 0644, nil)
				if err != nil {
					errMsg := "failed opening bolt"
					utils.AppLog.Error(errMsg, zap.Error(err))
					utils.SlackLog.Error(errMsg, err)
					return
				}
				database := BoltDB{DB: boltDB}
				if err := database.Store("token", *repo.ID, tokenByte); err != nil {
					errMsg := "error storing token in bolt"
					utils.AppLog.Error(errMsg, zap.Error(err))
					utils.SlackLog.Error(errMsg, err)
					return
				}

				limitByte, err := json.Marshal(limit)
				if err != nil {
					errMsg := "error converting callback limit"
					utils.AppLog.Error(errMsg, zap.Error(err))
					utils.SlackLog.Error(errMsg, err)
					return
				}
				if err := database.Store("limit", *repo.ID, limitByte); err != nil {
					errMsg := "error storing limit in bolt"
					utils.AppLog.Error(errMsg, zap.Error(err))
					utils.SlackLog.Error(errMsg, err)
					return
				}
				boltDB.Close()
			}()

			activationParams := struct {
				Repo  github.Repository `json:"repo"`
				Token *oauth2.Token     `json:"token"`
				Limit time.Time         `json:"limit"`
			}{
				*repo,
				token,
				limit,
			}
			utils.SlackLog.Info(fmt.Sprintf("Callback signup: %v", *repo.FullName))
			payload, err := json.Marshal(activationParams)
			if err != nil {
				errMsg := "failure converting activation parameters"
				utils.AppLog.Error(errMsg, zap.Error(err))
				utils.SlackLog.Error(errMsg, err)
				continue
			}
			req, err := http.NewRequest("POST", utils.Config.ActivationServiceEndpoint, bytes.NewBuffer(payload))
			if err != nil {
				errMsg := "failed to create http request"
				utils.AppLog.Error(errMsg, zap.Error(err))
				utils.SlackLog.Error(errMsg, err)
				continue
			}
			req.Header.Set("content-type", "application/json")
			resp, err := fs.httpClient.Do(req)
			if err != nil {
				errMsg := "failed internal post call"
				utils.AppLog.Error(errMsg, zap.Error(err))
				utils.SlackLog.Error(errMsg, err)
				continue
			} else {
				defer resp.Body.Close()
			}
		}
		utils.SlackLog.Info("Completed user signed up")
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(setup))
	}
}
