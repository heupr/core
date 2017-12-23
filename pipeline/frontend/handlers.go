package frontend

import (
	"context"
	"io/ioutil"
	"net/http"

	"github.com/google/go-github/github"
	"github.com/gorilla/schema"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	ghoa "golang.org/x/oauth2/github"

	"core/utils"
)

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

// Dropdowns is a holder for information to be populated into the template.
type Dropdowns struct {
	Repos  map[int]string
	Labels map[int]string
}

func consoleHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		if oauthState != r.FormValue("state") {
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}
		code := r.FormValue("code")
		client, err := newClient(code)
		if err != nil {
			http.Redirect(w, r, "/", http.StatusInternalServerError)
			utils.AppLog.Error(
				"failure creating frontend client",
				zap.Error(err),
			)
			return
		}

		opts := &github.ListOptions{PerPage: 100}
		repoOptions := make(map[int]string)
		for {
			repos, resp, err := client.Apps.ListUserRepos(
				context.Background(),
				5535,
				opts,
			)
			if err != nil {
				http.Redirect(w, r, "/", http.StatusInternalServerError)
				utils.AppLog.Error(
					"error collecting user repos",
					zap.Error(err),
				)
			}
			for i := range repos {
				repoOptions[*repos[i].ID] = *repos[i].FullName
			}

			if resp.NextPage == 0 {
				break
			} else {
				opts.Page = resp.NextPage
			}
		}

		// TODO:
		// - Add repoOptions to a dropdown struct (repos field)
		// - Build logic for looping/catching all repo labels
		// - Place into options and onto dropdown struct (labels field)
		// - The ReadFile below will change to a template ParseFile
		// - Execute template w/ the collected dropdown struct data
		// - Make changes to the console.html file
		// - https://www.socketloop.com/tutorials/golang-populate-dropdown-with-html-template-example
		data, err := ioutil.ReadFile("website2/console.html")
		if err != nil {
			if PROD {
				utils.SlackLog.Error(
					"Error generating console page",
					zap.Error(err),
				)
			}
			http.Redirect(w, r, "/", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(string(data)))
	} else if r.Method == "POST" {
		r.ParseForm()
		// if r.Form["id"] == "repo-selection" {
		// 	// populate label dropdowns
		// } else if r.Form["id"] == "label-settings" {
		// 	// collect data received into "snapshot"
		// }
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
