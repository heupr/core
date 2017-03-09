package frontend

import (
	"fmt"
	"net/http"

	"github.com/boltdb/bolt"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"

	"coralreefci/engine/gateway/conflation"
	"coralreefci/models"
)

type HeuprServer struct {
	Server    http.Server
	Models    map[int]models.Model
	Conflator conflation.Conflator
	Database  BoltDB
}

func (h *HeuprServer) routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", mainHandle)
	mux.HandleFunc("/login", githubLoginHandle)
	// mux.Handle("/hook", collectorHandler()) // TEMPORARILY REMOVED
	// TODO: This is a temporary work around until actual code can be built
	//       that will return the necessary repository struct.
	//       NOTE: http.HandleFunc("/select", githubRepoSelect) <- EXAMPLE
	mux.Handle("/test", h.TesthookHandler(testRepos(), testClient()))
	login := "heupr"
	user := &github.User{Login: &login}
	// NOTE: This is also a workaround for temporary testing measures.
	repo := github.Repository{Name: github.String("test"), Owner: user, ID: github.Int(81689981)}
	mux.Handle("/github_oauth_cb", h.hookHandler(&repo))
	// TODO: Implement resource handler function.
	// mux.Handle("/select_repo", h.selectRepo())
	// ^ NOTE: Needs acces to the authenticated client.
	return mux
}

func (h *HeuprServer) openDB() error {
	boltDB, err := bolt.Open("storage.db", 0644, nil)
	if err != nil {
		return err
	}
	h.Database = BoltDB{db: boltDB}
	return nil
}

func (h *HeuprServer) closeDB() {
	h.Database.db.Close()
}

func (h *HeuprServer) Start() {
	h.Server = http.Server{Addr: "127.0.0.1:8080", Handler: h.routes()}
	// TODO: Add in logging and remove print statement.
	err := h.Server.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}
}

func (h *HeuprServer) Stop() {
	//TODO: Closing the server down is a needed operation that will be
	//      implemented in the future.
}

// TODO: Rename and move into separate file.
func (h *HeuprServer) hookHandler(repo *github.Repository) http.Handler {
	handler, client := githubCallbackHandle()
	h.NewHook(repo, client)
	// TODO: This is an example of possible implementation:
	// errHandler := h.NewHook(repo, client)
	// if errHandler != nil {
	//     handler = errHandler
	// }
	return handler
	// TODO: Some sort of check to ensure the incoming traffic / payload is
	//	     what you're looking to receive (e.g. the request to setup a new
	//       webhook on the target repo - whatever struct that happens to be).
}

/*
func selectRepo() (http.Handler, *github.Repository) {
    repo := *github.Repository{}
    handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

    })
}
*/

func testClient() *github.Client {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: "634a8f39667f799a99bf2d7a852fcc5cbe412c93"})
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)
	return client
}

func testRepos() []*github.Repository {
	login := "heupr"
	user := &github.User{Login: &login}
	repo1 := github.Repository{Name: github.String("test"), Owner: user, ID: github.Int(81689981)}
	repo2 := github.Repository{Name: github.String("test2"), Owner: user, ID: github.Int(84002303)}
	return []*github.Repository{&repo1, &repo2}
}

func (h *HeuprServer) TesthookHandler(repos []*github.Repository, client *github.Client) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for i := 0; i < len(repos); i++ {
			h.NewHook(repos[i], client)
			h.AddModel(repos[i], client)
		}
	})
}
