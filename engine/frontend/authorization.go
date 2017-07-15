package frontend

import (
	"net/http"

	"golang.org/x/oauth2"
	ghoa "golang.org/x/oauth2/github"
)

var oaConfig = &oauth2.Config{
	// NOTE: Both fields will be available after registering Heupr w/ GitHub.
	ClientID:     "",
	ClientSecret: "",
	Endpoint:     ghoa.Endpoint,
	Scopes:       []string{"admin:repo_hook", "repo:status", "public_repo"},
}

const oaState = "the-force-shall-set-me-free"

var mainHandler = http.StripPrefix("/", http.FileServer(http.Dir("website/public/")))

// func mainHandler(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "text/html; charset=utf-8")
// 	w.WriteHeader(http.StatusOK)
// 	w.Write([]byte(websiteHTML))
// }

func githubLoginHandler(w http.ResponseWriter, r *http.Request) {
	url := oaConfig.AuthCodeURL(oaState, oauth2.AccessTypeOnline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}
