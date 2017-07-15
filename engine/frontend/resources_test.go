package frontend

import (
	"context"
	"testing"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func Test_listRepositories(t *testing.T) {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: "039f5f2f98a87f46abef10170866ed8ecf3b5b2d"})
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := *github.NewClient(tc)

	r, e := listRepositories(&client)

	names := []string{}

	for i := 0; i < len(r); i++ {
		names = append(names, *r[i].Name)
	}

	user, _, _ := client.Users.Get(context.Background(), "")

	if e != nil {
		t.Error(
			"\nFailing to generate repositories for target user",
			"\nUser:  ", *user.Login,
			"\nCount: ", len(r),
			"\nNames: ", names,
		)
	}
}
