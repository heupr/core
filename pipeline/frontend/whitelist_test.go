package frontend

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/google/go-github/github"
)

func TestAutomaticWhitelist(t *testing.T) {
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
	testToken := []byte("RIGHT!")

	databaseName = name

	err = testFS.AutomaticWhitelist(*testRepo, testToken)
	if err != nil {
		t.Errorf("automatic whitelist: %v", err)
	}
}
