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

// TODO: Add in uniformed error handling function.
// [ ] PROD boolean (generates Slack message)
// [ ] accepts error type argument
// [ ] returns redirect w/ arguments

// TODO: Add in consts for HTTP statuses.

func httpRedirect(w http.ResponseWriter, r *http.Request) {
	if PROD {
		http.Redirect(w, r, "https://heupr.io", http.StatusMovedPermanently)
	} else {
		http.Redirect(
			w,
			r,
			"https://127.0.0.1:8081",
			http.StatusMovedPermanently,
		)
	}
}

func staticHandler(filepath string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, err := ioutil.ReadFile(filepath)
		if err != nil {
			if PROD {
				utils.SlackLog.Error(
					"Error generating landing page",
					zap.Error(err),
				)
			}
			http.Redirect(w, r, "/", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(string(data)))
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

type storage struct {
	Name       string
	Labels     []string
	Selections map[string][]string
}

// Dropdowns is a holder for information to be populated into the template.
type Dropdowns struct {
	Repos  map[int]string
	Labels map[int][]string
}

func reposHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		if oauthState != r.FormValue("state") {
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}
		code := r.FormValue("code")
		client, err := newClient(code)
		if err != nil {
			utils.AppLog.Error(
				"failure creating frontend client",
				zap.Error(err),
			)
			http.Redirect(w, r, "/", http.StatusInternalServerError)
			return
		}

		opts := &github.ListOptions{PerPage: 100}
		repos := make(map[int]string)
		ctx := context.Background()
		for {
			repo, resp, err := client.Apps.ListUserRepos(ctx, 5535, opts)
			if err != nil {
				utils.AppLog.Error(
					"error collecting user repos",
					zap.Error(err),
				)
				http.Redirect(w, r, "/", http.StatusInternalServerError)
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
				l, resp, err := client.Issues.ListLabels(
					ctx,
					name[0],
					name[1],
					opts,
				)
				if err != nil {
					utils.AppLog.Error(
						"error collecting repo labels",
						zap.Error(err),
					)
					http.Redirect(w, r, "/", http.StatusInternalServerError)
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
					utils.AppLog.Error(
						"error creating storage file",
						zap.Error(err),
					)
					http.Redirect(w, r, "/", http.StatusInternalServerError)
					return
				}

				s := storage{
					Name:   name,
					Labels: labels[id],
				}
				encoder := gob.NewEncoder(file)
				if err := encoder.Encode(s); err != nil {
					utils.AppLog.Error(
						"error encoding info to new file",
						zap.Error(err),
					)
					http.Redirect(w, r, "/", http.StatusInternalServerError)
					return
				}
			} else {
				file, err := os.Open(filename)
				defer file.Close()
				if err != nil {
					http.Redirect(w, r, "/", http.StatusInternalServerError)
					return
				}
				decoder := gob.NewDecoder(file)
				s := storage{}
				decoder.Decode(&s)
				s.Labels = labels[id]
				encoder := gob.NewEncoder(file)
				if err := encoder.Encode(s); err != nil {
					utils.AppLog.Error(
						"error re-encoding info to file",
						zap.Error(err),
					)
					http.Redirect(w, r, "/", http.StatusInternalServerError)
					return
				}
			}
		}

		// NOTE: Possibly change to an anonymous struct.
		dropdowns := Dropdowns{
			Repos: repos,
		}

		t, err := template.ParseFiles("website2/repos.html")
		if err != nil {
			if PROD {
				utils.SlackLog.Error("Repos selection page", zap.Error(err))
			}
			http.Redirect(w, r, "/", http.StatusInternalServerError)
			return
		}
		t.Execute(w, dropdowns)
	}
}

func consoleHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()
		if r.Form["state"][0] != oauthState {
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}

		repo := r.Form["repo-selection"][0]
		if repo == "" {
			http.Redirect(w, r, "/", http.StatusBadRequest)
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
			http.Redirect(w, r, "/", http.StatusInternalServerError)
			return
		}

		f, err := os.Open(file)
		if err != nil {
			fmt.Println("HERE") // TEMPORARY
			http.Redirect(w, r, "/", http.StatusInternalServerError)
			return
		}
		decoder := gob.NewDecoder(f)
		s := storage{}
		decoder.Decode(&s)

		// TODO: populate label dropdowns from GitHub storage
		// [X] retrieve GitHub repo ID
		// [X] pull in GitHub files w/ matching ID
		// [X] pull in settings files w/ matching ID
		// [X] collect data from both files
		// [X] conflate file contents in-memory
		// [ ] populate into template
		// [ ] execute template
		// TODO: collect label selections into "snapshot"

	} else {
		http.Redirect(w, r, "/", http.StatusBadRequest)
		return
	}
}

func setupCompleteHandler(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadFile("website2/setup-complete.html")
	if err != nil {
		if PROD {
			utils.SlackLog.Error(
				"Error generating setup complete page",
				zap.Error(err),
			)
		}
		http.Redirect(w, r, "/", http.StatusInternalServerError)
		return
	}
	utils.AppLog.Info("Completed user signed up")
	if PROD {
		utils.SlackLog.Info("Completed user signed up")
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(string(data)))
}

// NOTE: Depreciate this code.
var decoder = schema.NewDecoder()

// NOTE: Depreciate this code.
var mainHandler = http.StripPrefix(
	"/",
	http.FileServer(http.Dir("../website/")),
)
