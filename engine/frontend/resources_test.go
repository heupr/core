package frontend

import (
	"testing"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func Test_listRepositories(t *testing.T) {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: "552cfadac27c94e91ce960c36cc3a1ee15fb134a"})
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := *github.NewClient(tc)

	r, e := listRepositories(&client)

	names := []string{}

	for i := 0; i < len(r); i++ {
		names = append(names, *r[i].Name)
	}

	user, _, _ := client.Users.Get("")

	if e != nil {
		t.Error(
			"\nFailing to generate repositories for target user",
			"\nUser:  ", *user.Login,
			"\nCount: ", len(r),
			"\nNames: ", names,
		)
	}
}
