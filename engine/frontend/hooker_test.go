package frontend

import (
	"coralreefci/models"
	"coralreefci/models/bhattacharya"
	"encoding/json"
	"fmt"
	"github.com/google/go-github/github"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestNewHook(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	// go server.Start()
	defer server.Close()

	client := github.NewClient(nil)
	url, _ := url.Parse(server.URL)
	client.BaseURL = url
	client.UploadURL = url

	mods := make(map[int]models.Model)
	mods[0] = models.Model{Algorithm: &bhattacharya.NBModel{}}

	testServer := HeuprServer{
		Models: mods,
	}
	mux.HandleFunc("/repos/nihilus/hunger/hooks", func(w http.ResponseWriter, r *http.Request) {
		v := new(github.Hook)
		json.NewDecoder(r.Body).Decode(v)
		fmt.Fprint(w, `{"id":1}`)
	})

	login := "nihilus"
	user := &github.User{Login: &login}
	name := "hunger"
	id := 1
	testRepo := github.Repository{
		Name:  &name,
		Owner: user,
		ID:    &id,
	}

	defer testServer.closeDB()
	err := testServer.openDB()
	if err != nil {
		t.Error(err) // TODO: Flesh out message
	}

	err = testServer.NewHook(&testRepo, client)
	if err != nil {
		t.Error(err) // TODO: Flesh out message
	}
	// server.Close()
	// fmt.Println("END OF TEST")
}

// mux.HandleFunc("/repos/nihilus/hunger/issues", func(w http.ResponseWriter, r *http.Request) {
//     name1 := "sion"
//     user1 := github.User{Login: &name1}
//     body1 := "I am pain"
//     issue1 := &github.Issue{
//         User: &user1,
//         Body: &body1,
//     }
//     name2 := "treya"
//     user2 := github.User{Login: &name2}
//     body2 := "I am betrayl"
//     issue2 := &github.Issue{
//         User: &user2,
//         Body: &body2,
//     }
//     list := []*github.Issue{issue1, issue2}
//     var issues struct {
//         Issues []github.Issue
//     }
//     json.NewDecoder(list).Decode(&issues)
// })
// mux.HandleFunc("/repos/nihilus/hunger/pulls", func(w http.ResponseWriter, r *http.Request) {})
