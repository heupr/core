package frontend

import (
	"net/http"

	"golang.org/x/oauth2"
	ghoa "golang.org/x/oauth2/github"
)

var oaConfig = &oauth2.Config{
	ClientID:     "5ffc021b1fe3702c6888",
	ClientSecret: "42edf1ab560dce313ff3e27dd7b94f58e41df3e7",
	Endpoint:     ghoa.Endpoint,
	Scopes:       []string{"admin:repo_hook", "public_repo"},
}

const oaState = "the-force-shall-set-me-free"

// I'm not sure if "../website/" or "website/" is correct - the first
// worked in testing. It may just depend on where the startup is called from.
var mainHandler = http.StripPrefix("/", http.FileServer(http.Dir("../website/")))

func githubLoginHandler(w http.ResponseWriter, r *http.Request) {
	url := oaConfig.AuthCodeURL(oaState, oauth2.AccessTypeOnline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}
