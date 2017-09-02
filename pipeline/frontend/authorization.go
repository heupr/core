package frontend

import (
	"net/http"

	"golang.org/x/oauth2"
	ghoa "golang.org/x/oauth2/github"
)

var oaConfig = &oauth2.Config{
	// NOTE: Both fields will be available after registering Heupr w/ GitHub.
	ClientID:     "CLIENT_ID",
	ClientSecret: "CLIENT_SECRET",
	Endpoint:     ghoa.Endpoint,
	Scopes:       []string{"admin:repo_hook", "repo:status", "public_repo"},
}

const oaState = "the-force-shall-set-me-free"

// NOTE: I'm not sure if "../website/" or "website/" is correct - the first
// worked in testing. It may just depend on where the startup is called from.
var mainHandler = http.StripPrefix("/", http.FileServer(http.Dir("../website/")))

func githubLoginHandler(w http.ResponseWriter, r *http.Request) {
	url := oaConfig.AuthCodeURL(oaState, oauth2.AccessTypeOnline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}
