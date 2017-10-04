package frontend

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/google/go-github/github"
)

func Test_checkTOML(t *testing.T) {
	owner := "darth-krayt"
	repo := "sith-holocron"
	users := []User{
		User{Owner: owner, Repo: repo},
	}
	output := checkTOML(users, owner, repo)
	expected := true
	if output != expected {
		t.Errorf("Returning incorrect status: expected %v, received %v", expected, output)
	}
}

func TestCheckWhitelist(t *testing.T) {
	name := "test-whitelist.db"
	file, err := ioutil.TempFile("", name)
	if err != nil {
		t.Errorf("generate whitelist test file %v", err)
	}
	file.Close()
	defer os.Remove(name)

	testFS := new(FrontendServer)

	n, l := "the-colonel", "chapmang"
	testRepo := &github.Repository{
		Name: &n,
		Owner: &github.User{
			Login: &l,
		},
	}

	databaseName = name

	_, err = testFS.CheckWhitelist(*testRepo)
	if err != nil {
		t.Errorf("automatic whitelist: %v", err)
	}
}
