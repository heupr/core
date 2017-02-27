package frontend

import (
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	ghoa "golang.org/x/oauth2/github"
	"net/http"
)

var oaConfig = &oauth2.Config{
	// TODO: Both of these fields should ultimately be available once the
	//       Heupr app is registered on GitHub - until then, these two fields
	//       hold nil string values
	ClientID:     "",
	ClientSecret: "",
	Endpoint:     ghoa.Endpoint,
	Scopes:       []string{"admin:repo_hook", "repo:status", "public_repo"},
}

// NOTE: This is a temporary secret; it will need to change.
var oaState = "the-force-shall-set-me-free"

const htmlIndex = `<html><body>
Log in with <a href="/login">GitHub</a>
</body></html>
`

func mainHandle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(htmlIndex))
}

func githubLoginHandle(w http.ResponseWriter, r *http.Request) {
	url := oaConfig.AuthCodeURL(oaState, oauth2.AccessTypeOnline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// func githubRepoSelect() (http.Handler, *github.Repository) {
//     repo := *github.Repository{}
//     handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//
//     })
// }

func githubCallbackHandle() (http.Handler, *github.Client) {
	client := github.NewClient(nil)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		state := r.FormValue("state")
		if state != oaState {
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}

		code := r.FormValue("code")
		token, err := oaConfig.Exchange(oauth2.NoContext, code)
		if err != nil {
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}
		// TODO: The token will need to be stored in the database as well; this
		//       will require the repository ID in order to store in the proper
		//       bucket. Store with the key value of "token."

		oaClient := oaConfig.Client(oauth2.NoContext, token)
		client = github.NewClient(oaClient)

		// TODO: Add call to add token to database
		//       tokenString := token.AccessToken <- argument into the function
		// NOTE: In order to properly store the token in the database as a
		//       value, it will need a corresponding key which will likely be
		//       the target repository ID number.

		_, _, err = client.Users.Get("") // NOTE: Not sure if "=" is correct.
		if err != nil {
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	})
	return handler, client
}
