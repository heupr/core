package frontend

import (
	"coralreefci/engine/gateway/conflation"
	"coralreefci/models"
	"fmt" // TEMPORARY
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"net/http"
)

type HeuprServer struct {
	Server    http.Server
	Models    map[int]models.Model
	Conflator conflation.Conflator
}

func (h *HeuprServer) Start() {
	mux := http.NewServeMux()
	h.Server = http.Server{}
	h.Server.Addr = "127.0.0.1:8080"
	h.Server.Handler = mux
	mux.HandleFunc("/", mainHandle)
	mux.HandleFunc("/login", githubLoginHandle)
	mux.Handle("/hook", collectorHandler(""))

	// TODO: This is a temporary work around until actual code can be built
	//       that will return the necessary repository struct.
	//       NOTE: http.HandleFunc("/select", githubRepoSelect) <- EXAMPLE
	login := "heupr"
	user := &github.User{Login: &login}
	name := "test"
	id := 1
	repo := github.Repository{
		Name:  &name,
		Owner: user,
		ID:    &id,
	}
	mux.Handle("/test", h.TesthookHandler(&repo, testClient()))
	mux.Handle("/github_oauth_cb", h.hookHandler(&repo))
	// if err := http.ListenAndServe(addr, nil); err != nil {
	// 	fmt.Println(err)
	// }
	// TODO: Add in logging and remove print statement.
	err := h.Server.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}
}

func (h *HeuprServer) Stop() {
	//TODO:
}

func (h *HeuprServer) hookHandler(repo *github.Repository) http.Handler {
	handler, client := githubCallbackHandle()
	h.NewHook(repo, client)
	return handler
	// TODO: Some sort of check to ensure the incoming traffic / payload is
	//	     what you're looking to receive (e.g. the request to setup a new
	//       webhook on the target repo - whatever struct that happens to be).
}

func testClient() *github.Client {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: "634a8f39667f799a99bf2d7a852fcc5cbe412c93"})
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)
	client.UserAgent = "heupr"
	return client
}

func (h *HeuprServer) TesthookHandler(repo *github.Repository, client *github.Client) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.NewHook(repo, client)
		h.AddModel(repo, client)
	})
}
