package frontend

import (
	"context"
	"encoding/gob"
	"fmt" // TEMPORARY
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/google/go-github/github"
	"github.com/gorilla/schema"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	ghoa "golang.org/x/oauth2/github"

	"core/utils"
)

func slackErr(msg string, err error) {
	if PROD {
		utils.SlackLog.Error(msg, zap.Error(err))
	}
}

func slackMsg(msg string) {
	if PROD {
		utils.SlackLog.Info(msg)
	}
}

func httpRedirect(w http.ResponseWriter, r *http.Request) {
	if PROD {
		http.Redirect(w, r, "https://heupr.io", http.StatusMovedPermanently)
	} else {
		http.Redirect(w, r, "https://127.0.0.1:8081", http.StatusMovedPermanently)
	}
}

func staticHandler(filepath string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, err := ioutil.ReadFile(filepath)
		if err != nil {
			slackErr("Error generating landing page", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	})
}

var (
	oauthConfig = &oauth2.Config{
		// NOTE: These will need to be added for production.
		ClientID:     "",
		ClientSecret: "",
		Scopes:       []string{""},
		Endpoint:     ghoa.Endpoint,
	}
	oauthState = "tenebrous-plagueis-sidious-maul-tyrannus-vader"
)

func loginHandler(w http.ResponseWriter, r *http.Request) {
	url := oauthConfig.AuthCodeURL(oauthState, oauth2.AccessTypeOnline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

var newClient = func(code string) (*github.Client, error) {
	token, err := oauthConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		return nil, err
	}
	client := github.NewClient(oauthConfig.Client(oauth2.NoContext, token))
	return client, nil
}

type label struct {
	Name     string
	Selected bool
}

type storage struct {
	Name    string // FullName for the given repo.
	Buckets map[string][]label
}

func updateStorage(s *storage, labels []string) {
	for bcktName, bcktLabels := range s.Buckets {
		updated := []label{}
		for i := range labels {
			label := label{Name: labels[i]}
			for j := range bcktLabels {
				if labels[i] == bcktLabels[j].Name {
					label.Selected = bcktLabels[j].Selected
				}
			}
			updated = append(updated, label)
		}
		s.Buckets[bcktName] = updated
	}
}

func reposHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "bad request method", http.StatusBadRequest)
		return
	}
	if oauthState != r.FormValue("state") {
		http.Error(w, "authorization error", http.StatusUnauthorized)
		return
	}
	code := r.FormValue("code")
	client, err := newClient(code)
	if err != nil {
		utils.AppLog.Error(
			"failure creating frontend client",
			zap.Error(err),
		)
		http.Error(w, "client failure", http.StatusInternalServerError)
		return
	}

	opts := &github.ListOptions{PerPage: 100}
	repos := make(map[int]string) // NOTE: NEEDED FOR RENDERING AS IS
	ctx := context.Background()
	for {
		repo, resp, err := client.Apps.ListUserRepos(ctx, 5535, opts)
		if err != nil {
			utils.AppLog.Error("error collecting user repos", zap.Error(err))
			http.Error(w, "error collecting user repos", http.StatusInternalServerError)
			return
		}
		for i := range repo {
			repos[*repo[i].ID] = *repo[i].FullName
		}

		if resp.NextPage == 0 {
			break
		} else {
			opts.Page = resp.NextPage
		}
	}

	opts = &github.ListOptions{PerPage: 100}
	labels := make(map[int][]string)
	for key, value := range repos {
		name := strings.Split(value, "/")
		for {
			l, resp, err := client.Issues.ListLabels(ctx, name[0], name[1], opts)
			if err != nil {
				utils.AppLog.Error("error collecting repo labels", zap.Error(err))
				http.Error(w, "error collecting repo labels", http.StatusInternalServerError)
				return
			}
			for i := range l {
				labels[key] = append(labels[key], *l[i].Name)
			}

			if resp.NextPage == 0 {
				break
			} else {
				opts.Page = resp.NextPage
			}
		}
	}

	for id, name := range repos {
		filename := strconv.Itoa(id) + ".gob"
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			file, err := os.Create(filename)
			defer file.Close()
			if err != nil {
				utils.AppLog.Error("error creating storage file", zap.Error(err))
				http.Error(w, "error creating storage file", http.StatusInternalServerError)
				return
			}

			s := storage{
				Name:    name,
				Buckets: make(map[string][]label),
			}

			for _, l := range labels[id] {
				s.Buckets[""] = append(s.Buckets[""], label{Name: l})
			}

			encoder := gob.NewEncoder(file)
			if err := encoder.Encode(s); err != nil {
				utils.AppLog.Error("error encoding info to new file", zap.Error(err))
				http.Error(w, "error encoding info to new file", http.StatusInternalServerError)
				return
			}
		} else {
			file, err := os.Open(filename)
			defer file.Close()
			if err != nil {
				http.Error(w, "error opening storage file", http.StatusInternalServerError)
				return
			}
			decoder := gob.NewDecoder(file)
			s := storage{}
			decoder.Decode(&s)

			updateStorage(&s, labels[id])

			encoder := gob.NewEncoder(file)
			if err := encoder.Encode(s); err != nil {
				utils.AppLog.Error("error re-encoding info to file", zap.Error(err))
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}

	// NOTE: Possibly change to an anonymous struct.
	input := struct {
		Repos map[int]string
	}{
		Repos: repos,
	}

	t, err := template.ParseFiles("website2/repos.html")
	if err != nil {
		slackErr("Repos selection page", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.Execute(w, input)
}

func consoleHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()
		if r.Form["state"][0] != oauthState {
			http.Error(w, "authorization error", http.StatusUnauthorized)
			return
		}

		repo := r.Form["repo-selection"][0]
		if repo == "" {
			http.Error(w, "request erro", http.StatusBadRequest)
			return
		}

		file := ""
		err := filepath.Walk(
			"./",
			func(path string, info os.FileInfo, err error) error {
				name := strings.TrimSuffix(path, filepath.Ext(path))
				if filepath.Ext(path) == ".gob" && name == repo {
					file = path
				}
				return nil
			})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		f, err := os.Open(file)
		if err != nil {
			fmt.Println("HERE") // TEMPORARY
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		decoder := gob.NewDecoder(f)
		s := storage{}
		decoder.Decode(&s)

		t, err := template.ParseFiles("website2/console.html")
		if err != nil {
			slackErr("Settings console page", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		t.Execute(w, s)
	} else {
		http.Error(w, "error loading console", http.StatusBadRequest)
		return
	}
}

func setupCompleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()

		data, err := ioutil.ReadFile("website2/setup-complete.html")
		if err != nil {
			slackErr("Error generating setup complete page", err)
			http.Error(w, "/", http.StatusInternalServerError)
			return
		}
		utils.AppLog.Info("Completed user signed up")
		slackMsg("Completed user signed up")
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	} else {
		http.Error(w, "error loading console", http.StatusBadRequest)
		return
	}

}

// NOTE: Depreciate this code.
var decoder = schema.NewDecoder()

// NOTE: Depreciate this code.
var mainHandler = http.StripPrefix(
	"/",
	http.FileServer(http.Dir("../website/")),
)
