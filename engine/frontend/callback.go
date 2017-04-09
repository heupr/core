package frontend

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/google/go-github/github"
	"github.com/gorilla/schema"
	"golang.org/x/oauth2"
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
    </body>
</html>
`

const setup = `
<html>
    <title>
        Heupr
    </title>
    <body>
        <p>Congratulations! Setup is complete!</p>
        <p>Issue assignments will go out in a few minutes through GitHub</p>
        <p>Return to the <a href="/">main page</a></p>
    </body>
</html>
`

var decoder = schema.NewDecoder()

type Resources struct {
	Repos []*github.Repository
}

func (h *HeuprServer) githubCallbackHandler(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("state") != oaState {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect) // TODO: Write specific redirect URL.
		return
	}

	token, err := oaConfig.Exchange(oauth2.NoContext, r.FormValue("code"))
	if err != nil {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect) // TODO: Write specific redirect URL
		return
	}

	oaClient := oaConfig.Client(oauth2.NoContext, token)
	client := github.NewClient(oaClient)

	_, _, err = client.Users.Get("")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect) // TODO: Write specific redirect URL
		return
	}

	repos, err := listRepositories(client)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect) // TODO: Write specific redirect URL
		return
		// return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//     // TODO: Check that this is the correct error status to use.
		//     http.Error(w, err.Error(), http.StatusInternalServerError)
		// })
	}

	if r.Method == "GET" {
		tmpl, err := template.New("options").Parse(options)
		if err != nil {
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect) // TODO: Write specific redirect URL
			return
		}
		resources := Resources{repos}
		tmpl.Execute(w, resources)
	} else {
		r.ParseForm()
		results := new(Resources)
		err := decoder.Decode(results, r.PostForm)
		if err != nil {
			fmt.Println(err)
		}
		// TODO: Call NewHook here.
		// - pass in: results.Repos & client variables
		// TODO: Call NewHeuprRepo here.
		// - pass in: results.Repos & client variables
        // TODO: Call AddModel here.
	}
	http.Redirect(w, r, "/setup_complete", http.StatusPermanentRedirect)
}

func completeHandle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(setup))
}
